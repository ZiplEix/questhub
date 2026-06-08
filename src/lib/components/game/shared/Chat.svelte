<script lang="ts">
    import { Send, EyeOff } from "lucide-svelte";
    import { websocketStore } from "$lib/websocket";
    import { sendMessage } from "$lib/chat";
    import { page } from "$app/state";
    import { untrack } from "svelte";
    import { parseDiceAndMath } from "$lib/utils/diceParser";
    import { supabase } from "$lib/supabaseClient";

    let {
        isGM = false,
        players = [],
        currentUserId = "",
        senderName = "Me",
    } = $props<{
        isGM?: boolean;
        players?: { id: string; name: string }[];
        currentUserId?: string;
        senderName?: string;
    }>();

    let newMessage = $state("");
    let isSecretRoll = $state(false);
    let whisperTarget = $state("");
    let chatContainer: HTMLDivElement;
    let inputElement: HTMLInputElement;
    let localMessages = $state<any[]>([]);

    let autocompleteIndex = $state(0);
    let hideAutocomplete = $state(false);

    let autocompleteList = $derived.by(() => {
        const text = newMessage.trimStart();
        if (!text.startsWith("/")) return [];

        const firstSpaceIdx = text.indexOf(" ");
        if (firstSpaceIdx === -1) {
            const searchCmd = text.toLowerCase();
            const commands = [
                { label: "/roll", value: "/roll ", desc: "Lancer des dés" },
                { label: "/r", value: "/r ", desc: "Lancer des dés (raccourci)" },
                { label: "/chuchotement", value: "/chuchotement ", desc: "Chuchoter à un joueur" },
                { label: "/w", value: "/w ", desc: "Chuchoter à un joueur (raccourci)" },
                { label: "/whisper", value: "/whisper ", desc: "Chuchoter à un joueur (alias)" },
                { label: "/gmroll", value: "/gmroll ", desc: "Jet secret pour le MJ" },
                { label: "/gr", value: "/gr ", desc: "Jet secret pour le MJ (raccourci)" },
                { label: "/secret", value: "/secret ", desc: "Message secret au MJ" },
                { label: "/s", value: "/s ", desc: "Message secret au MJ (raccourci)" },
                { label: "/me", value: "/me ", desc: "Décrire une action" },
                { label: "/aide", value: "/aide", desc: "Afficher l'aide" },
                { label: "/help", value: "/help", desc: "Afficher l'aide (alias)" }
            ];
            return commands.filter(c => c.label.startsWith(searchCmd));
        } else {
            const command = text.substring(0, firstSpaceIdx).toLowerCase();
            const args = text.substring(firstSpaceIdx + 1);

            if (command === "/chuchotement" || command === "/w" || command === "/whisper") {
                let searchPlayer = args;
                if (args.startsWith('"')) {
                    const closeQuoteIdx = args.indexOf('"', 1);
                    if (closeQuoteIdx !== -1) {
                        return [];
                    }
                    searchPlayer = args.substring(1);
                } else {
                    if (args.includes(" ")) {
                        return [];
                    }
                }

                const searchLower = searchPlayer.toLowerCase();
                return players
                    .filter((p: { id: string; name: string }) => p.name.toLowerCase().startsWith(searchLower))
                    .map((p: { id: string; name: string }) => {
                        const value = p.name.includes(" ") ? `"${p.name}" ` : `${p.name} `;
                        return {
                            label: p.name,
                            value: `${command} ${value}`,
                            desc: "Joueur de la campagne"
                        };
                    });
            }
        }
        return [];
    });

    $effect(() => {
        newMessage;
        untrack(() => {
            hideAutocomplete = false;
            autocompleteIndex = 0;
        });
    });

    function selectAutocomplete(item: { label: string; value: string }) {
        newMessage = item.value;
        hideAutocomplete = true;
        autocompleteIndex = 0;
        if (inputElement) {
            inputElement.focus();
        }
    }

    function handleKeyDown(e: KeyboardEvent) {
        if (autocompleteList.length > 0 && !hideAutocomplete) {
            if (e.key === "ArrowDown") {
                e.preventDefault();
                autocompleteIndex = (autocompleteIndex + 1) % autocompleteList.length;
            } else if (e.key === "ArrowUp") {
                e.preventDefault();
                autocompleteIndex = (autocompleteIndex - 1 + autocompleteList.length) % autocompleteList.length;
            } else if (e.key === "Enter" || e.key === "Tab") {
                e.preventDefault();
                selectAutocomplete(autocompleteList[autocompleteIndex]);
            } else if (e.key === "Escape") {
                e.preventDefault();
                hideAutocomplete = true;
            }
        } else if (e.key === "Enter") {
            handleSendMessage();
        }
    }

    let allMessages = $derived.by(() => {
        const dbMsgs = $websocketStore.messages || [];
        return [...dbMsgs, ...localMessages].sort((a, b) => {
            const timeA = a.created_at ? new Date(a.created_at).getTime() : (a.local_timestamp || 0);
            const timeB = b.created_at ? new Date(b.created_at).getTime() : (b.local_timestamp || 0);
            return timeA - timeB;
        });
    });

    // Auto-scroll to bottom when messages change
    $effect(() => {
        const messages = allMessages;
        untrack(() => {
            if (chatContainer) {
                setTimeout(() => {
                    chatContainer.scrollTop = chatContainer.scrollHeight;
                }, 0);
            }
        });
    });

    function addSystemMessage(content: string) {
        localMessages = [
            ...localMessages,
            {
                id: `system-${Date.now()}-${Math.random()}`,
                type: "SYSTEM",
                sender_name: "Système",
                content,
                local_timestamp: Date.now()
            }
        ];
    }

    async function handleChatCommand(rawText: string, gameId: string) {
        const spaceIndex = rawText.indexOf(" ");
        const command = (spaceIndex === -1 ? rawText : rawText.substring(0, spaceIndex)).toLowerCase();
        const args = spaceIndex === -1 ? "" : rawText.substring(spaceIndex + 1).trim();

        if (command === "/aide" || command === "/help" || command === "/?") {
            addSystemMessage("help_placeholder");
            return;
        }

        if (command === "/roll" || command === "/r" || command === "/gmroll" || command === "/gr") {
            if (!args) {
                addSystemMessage("⚠️ Commande /roll invalide. Exemple : /roll 2d6+4");
                return;
            }

            const parseResult = parseDiceAndMath(args);
            if (!parseResult) {
                addSystemMessage(`⚠️ Formule de dés invalide : "${args}"`);
                return;
            }

            const isSecretRollCmd = command === "/gmroll" || command === "/gr" || isSecretRoll;
            const gmItem = players.find((p: { id: string; name: string }) => p.name === "GM");
            const gmId = gmItem ? gmItem.id : currentUserId;

            let type = "EVENT";
            let targetId: string | undefined = undefined;

            if (isSecretRollCmd) {
                type = "CHAT_PRIVATE";
                targetId = isGM ? currentUserId : gmId;
            }

            const formattedMsg = `🎲 Jet : \`${args}\` ➔ **${parseResult.logText}**`;

            await sendMessage({
                game_id: gameId,
                content: formattedMsg,
                type,
                sender_name: senderName,
                target_id: targetId
            });
            return;
        }

        if (command === "/chuchotement" || command === "/w" || command === "/whisper") {
            if (!args) {
                addSystemMessage("⚠️ Commande /chuchotement invalide. Exemple : /w Baptiste Salut !");
                return;
            }

            let targetName = "";
            let messageContent = "";

            if (args.startsWith('"')) {
                const closeQuoteIndex = args.indexOf('"', 1);
                if (closeQuoteIndex === -1) {
                    addSystemMessage("⚠️ Commande invalide : guillemet fermant manquant pour le nom du joueur.");
                    return;
                }
                targetName = args.substring(1, closeQuoteIndex);
                messageContent = args.substring(closeQuoteIndex + 1).trim();
            } else {
                const spaceIdx = args.indexOf(" ");
                if (spaceIdx === -1) {
                    addSystemMessage("⚠️ Commande invalide : veuillez spécifier un message à envoyer.");
                    return;
                }
                targetName = args.substring(0, spaceIdx);
                messageContent = args.substring(spaceIdx + 1).trim();
            }

            if (!targetName || !messageContent) {
                addSystemMessage("⚠️ Commande /chuchotement invalide. Exemple : /w Baptiste Salut !");
                return;
            }

            const targetLower = targetName.toLowerCase();
            const matchedPlayers = players.filter((p: { id: string; name: string }) => p.name.toLowerCase().startsWith(targetLower));

            if (matchedPlayers.length === 0) {
                addSystemMessage(`⚠️ Aucun joueur trouvé commençant par "${targetName}".`);
                return;
            }
            if (matchedPlayers.length > 1) {
                addSystemMessage(`⚠️ Plusieurs joueurs correspondent à "${targetName}" : ${matchedPlayers.map((p: { id: string; name: string }) => p.name).join(", ")}.`);
                return;
            }

            const targetPlayer = matchedPlayers[0];

            await sendMessage({
                game_id: gameId,
                content: messageContent,
                type: "CHAT_PRIVATE",
                sender_name: senderName,
                target_id: targetPlayer.id
            });
            return;
        }

        if (command === "/secret" || command === "/s") {
            if (!args) {
                addSystemMessage("⚠️ Commande /secret invalide. Exemple : /s j'inspecte la serrure");
                return;
            }

            const gmItem = players.find((p: { id: string; name: string }) => p.name === "GM");
            const gmId = gmItem ? gmItem.id : currentUserId;

            await sendMessage({
                game_id: gameId,
                content: args,
                type: "CHAT_PRIVATE",
                sender_name: senderName,
                target_id: isGM ? currentUserId : gmId
            });
            return;
        }

        if (command === "/me") {
            if (!args) {
                addSystemMessage("⚠️ Commande /me invalide. Exemple : /me sourit");
                return;
            }

            await sendMessage({
                game_id: gameId,
                content: `*${args}*`,
                type: "CHAT_GLOBAL",
                sender_name: senderName
            });
            return;
        }

        addSystemMessage(`⚠️ Commande inconnue : "${command}". Tapez /aide pour voir la liste des commandes.`);
    }

    function handleSendMessage() {
        const rawText = newMessage.trim();
        if (!rawText) return;

        newMessage = "";

        const gameId = page.params.id;
        if (!gameId) return;

        if (rawText.startsWith("/")) {
            handleChatCommand(rawText, gameId);
            return;
        }

        let type = "CHAT_GLOBAL";
        let targetId = whisperTarget;

        if (isSecretRoll) {
            type = "CHAT_PRIVATE";
            targetId = currentUserId || "";
        } else if (whisperTarget) {
            type = "CHAT_PRIVATE";
        }

        const payload: any = {
            type,
            game_id: gameId,
            content: rawText,
            sender_name: senderName,
        };

        if (targetId) {
            payload.target_id = targetId;
        }

        sendMessage(payload);
    }
    function getTargetName(id: string) {
        const p = players.find(
            (p: { id: string; name: string }) => p.id === id,
        );
        return p ? p.name : "Inconnu";
    }

    function formatMessageContent(content: string): string {
        if (!content) return "";
        // Escape HTML to prevent XSS
        let escaped = content
            .replace(/&/g, "&amp;")
            .replace(/</g, "&lt;")
            .replace(/>/g, "&gt;")
            .replace(/"/g, "&quot;")
            .replace(/'/g, "&#039;");

        // Parse bold: **text** -> <strong>text</strong>
        escaped = escaped.replace(/\*\*(.*?)\*\*/g, "<strong>$1</strong>");

        // Parse italic: *text* -> <em>text</em>
        escaped = escaped.replace(/\*(.*?)\*/g, "<em>$1</em>");

        // Parse code: `text` -> <code class="bg-black/10 px-1 py-0.5 rounded font-mono text-[11px]">text</code>
        escaped = escaped.replace(/`(.*?)`/g, '<code class="bg-black/10 px-1 py-0.5 rounded font-mono text-[11px]">$1</code>');

        return escaped;
    }
</script>

<div class="h-full flex flex-col bg-stone-50">
    <!-- Messages Area -->
    <div bind:this={chatContainer} class="flex-1 overflow-y-auto p-4 space-y-3">
        {#each allMessages as msg}
            {#if msg.type === "SYSTEM"}
                <div class="flex justify-center my-2 w-full">
                    {#if msg.content === "help_placeholder"}
                        <div
                            class="bg-stone-100 text-stone-600 text-xs px-4 py-3 rounded-2xl border border-stone-200 font-medium w-full space-y-2 shadow-xs text-left"
                        >
                            <p class="font-bold text-stone-800 flex items-center gap-1.5 text-[13px]">
                                <span>📌 Commandes de chat disponibles :</span>
                            </p>
                            <div class="space-y-1.5 text-stone-600 font-mono text-[10px] leading-relaxed bg-white/50 p-2.5 rounded-xl border border-stone-150">
                                <div><strong class="text-burnt-orange font-sans text-xs">/roll [formule]</strong> (ou <strong class="text-burnt-orange font-sans text-xs">/r</strong>)</div>
                                <div class="text-stone-400 pl-3">Lance les dés (ex: /r 2d6+4)</div>
                                
                                <div><strong class="text-burnt-orange font-sans text-xs">/gmroll [formule]</strong> (ou <strong class="text-burnt-orange font-sans text-xs">/gr</strong>)</div>
                                <div class="text-stone-400 pl-3">Lance les dés secrètement pour le MJ</div>

                                <div><strong class="text-burnt-orange font-sans text-xs">/chuchotement [nom] [msg]</strong> (ou <strong class="text-burnt-orange font-sans text-xs">/w</strong>)</div>
                                <div class="text-stone-400 pl-3">Chuchote à un joueur (ex: /w Baptiste coucou)</div>

                                <div><strong class="text-burnt-orange font-sans text-xs">/secret [msg]</strong> (ou <strong class="text-burnt-orange font-sans text-xs">/s</strong>)</div>
                                <div class="text-stone-400 pl-3">Message secret visible uniquement par le MJ</div>

                                <div><strong class="text-burnt-orange font-sans text-xs">/me [action]</strong></div>
                                <div class="text-stone-400 pl-3">Décrit une action en italique (ex: /me fouille la pièce)</div>

                                <div><strong class="text-burnt-orange font-sans text-xs">/aide</strong> (ou <strong class="text-burnt-orange font-sans text-xs">/help</strong>, <strong class="text-burnt-orange font-sans text-xs">/?</strong>)</div>
                                <div class="text-stone-400 pl-3">Affiche cette aide</div>
                            </div>
                        </div>
                    {:else}
                        <div
                            class="bg-rose-50 text-rose-600 text-xs px-3.5 py-2 rounded-xl border border-rose-200 font-medium shadow-xs"
                        >
                            {msg.content}
                        </div>
                    {/if}
                </div>
            {:else if msg.type === "EVENT"}
                <div class="flex justify-center my-2">
                    <div
                        class="bg-stone-200 text-stone-600 text-xs px-3 py-1 rounded-full italic border border-stone-300"
                    >
                        <span class="font-bold">{msg.sender_name}</span>
                        {@html formatMessageContent(msg.content)}
                    </div>
                </div>
            {:else if msg.type.startsWith("CHAT")}
                <div
                    class="flex flex-col {msg.sender_id === currentUserId
                        ? 'items-end'
                        : 'items-start'}"
                >
                    <div class="flex items-baseline gap-2 mb-1">
                        <span class="text-xs font-bold text-stone-500"
                            >{msg.sender_name}</span
                        >
                        <span class="text-[10px] text-stone-400">
                            {new Date(msg.created_at).toLocaleTimeString([], {
                                hour: "2-digit",
                                minute: "2-digit",
                            })}
                        </span>
                    </div>
                    <div
                        class="max-w-[85%] p-2.5 rounded-xl text-sm shadow-sm
                        {msg.type === 'CHAT_PRIVATE'
                            ? msg.sender_id === msg.target_id
                                ? 'bg-stone-800 text-stone-300 border border-stone-700'
                                : 'bg-purple-100 text-purple-800 border border-purple-200'
                            : msg.sender_name === 'GM'
                              ? 'bg-burnt-orange text-white rounded-tr-none'
                              : 'bg-white text-stone-800 border border-stone-200 rounded-tl-none'}"
                    >
                        {#if msg.type === "CHAT_PRIVATE"}
                            {#if msg.sender_id === msg.target_id}
                                <div
                                    class="flex items-center gap-1 mb-1 text-xs font-bold uppercase opacity-70 text-stone-400"
                                >
                                    <EyeOff size={10} /> Message Secret
                                </div>
                            {:else if msg.sender_id === currentUserId}
                                <div
                                    class="flex items-center gap-1 mb-1 text-xs font-bold uppercase opacity-70"
                                >
                                    <EyeOff size={10} /> Chuchotement à {getTargetName(
                                        msg.target_id,
                                    )}
                                </div>
                            {:else}
                                <div
                                    class="flex items-center gap-1 mb-1 text-xs font-bold uppercase opacity-70"
                                >
                                    <EyeOff size={10} /> Chuchotement de {msg.sender_name}
                                </div>
                            {/if}
                        {/if}

                        {@html formatMessageContent(msg.content)}
                    </div>
                </div>
            {/if}
        {/each}
    </div>

    <!-- Input Area -->
    <div class="p-3 bg-white border-t border-stone-200">
        <!-- Tools -->
        <div class="flex gap-2 mb-2">
            {#if isGM}
                <button
                    class="flex items-center gap-1 px-2 py-1 rounded text-xs font-bold transition-colors border
                    {isSecretRoll
                        ? 'bg-stone-800 text-white border-stone-800'
                        : 'bg-stone-100 text-stone-500 border-stone-200 hover:bg-stone-200'}"
                    onclick={() => (isSecretRoll = !isSecretRoll)}
                    title="Message secret (visible uniquement par vous)"
                >
                    <EyeOff size={12} />
                    Secret
                </button>
            {/if}
            <select
                bind:value={whisperTarget}
                class="px-2 py-1 rounded text-xs font-bold bg-stone-100 text-stone-500 border border-stone-200 outline-none focus:border-burnt-orange"
            >
                <option value="">À tous</option>
                {#each players as player}
                    <option value={player.id}>{player.name}</option>
                {/each}
            </select>
        </div>

        <div class="relative">
            {#if autocompleteList.length > 0 && !hideAutocomplete}
                <div
                    class="absolute bottom-full left-0 right-0 mb-2 bg-white/95 backdrop-blur-md border border-stone-200 shadow-xl rounded-xl overflow-hidden z-50 flex flex-col divide-y divide-stone-100 max-h-60 overflow-y-auto"
                    onmousedown={(e) => e.preventDefault()}
                >
                    {#each autocompleteList as item, index}
                        <!-- svelte-ignore a11y_click_events_have_key_events -->
                        <!-- svelte-ignore a11y_no_static_element_interactions -->
                        <div
                            class="flex items-center justify-between py-2 px-3.5 cursor-pointer transition-all border-l-4
                            {index === autocompleteIndex
                                ? 'bg-burnt-orange/10 border-burnt-orange text-burnt-orange font-medium pl-2.5'
                                : 'border-transparent text-stone-700 hover:bg-stone-50'}"
                            onmousedown={(e) => {
                                e.preventDefault();
                                selectAutocomplete(item);
                            }}
                        >
                            <div class="flex items-center gap-2">
                                <span class="font-mono text-xs">{item.label}</span>
                                {#if item.desc}
                                    <span class="text-xs opacity-60 font-sans">{item.desc}</span>
                                {/if}
                            </div>
                            {#if index === autocompleteIndex}
                                <span class="text-[10px] uppercase tracking-wider font-bold opacity-80">Entrée</span>
                            {/if}
                        </div>
                    {/each}
                </div>
            {/if}
            <input
                bind:this={inputElement}
                type="text"
                bind:value={newMessage}
                onkeydown={handleKeyDown}
                onblur={() => hideAutocomplete = true}
                placeholder={isSecretRoll ? "Message secret..." : "Message..."}
                class="w-full pl-4 pr-10 py-2.5 bg-stone-50 border border-stone-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-burnt-orange/20 focus:border-burnt-orange transition-all"
            />
            <button
                onclick={handleSendMessage}
                class="absolute right-2 top-1/2 -translate-y-1/2 p-1.5 text-burnt-orange hover:bg-burnt-orange/10 rounded-lg transition-colors"
            >
                <Send size={16} />
            </button>
        </div>
    </div>
</div>
