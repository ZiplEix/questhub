<script lang="ts">
    import { Send, Dices, EyeOff, MessageSquare } from "lucide-svelte";
    import { websocketStore, sendMessage } from "$lib/websocket";
    import { page } from "$app/state";
    import { untrack } from "svelte";

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

    // Auto-scroll to bottom when messages change
    $effect(() => {
        const messages = $websocketStore.messages;
        untrack(() => {
            if (chatContainer) {
                setTimeout(() => {
                    chatContainer.scrollTop = chatContainer.scrollHeight;
                }, 0);
            }
        });
    });

    function handleSendMessage() {
        if (!newMessage.trim()) return;

        const gameId = page.params.id;
        let type = "CHAT_GLOBAL";
        let targetId = whisperTarget;

        if (isSecretRoll) {
            type = "CHAT_PRIVATE";
            targetId = currentUserId || ""; // Send to self
        } else if (whisperTarget) {
            type = "CHAT_PRIVATE";
        }

        const payload: any = {
            type,
            game_id: gameId,
            content: newMessage,
            sender_name: senderName,
        };

        if (targetId) {
            payload.target_id = targetId;
        }

        sendMessage(payload);
        newMessage = "";
    }
    function getTargetName(id: string) {
        const p = players.find(
            (p: { id: string; name: string }) => p.id === id,
        );
        return p ? p.name : "Inconnu";
    }
</script>

<div class="h-full flex flex-col bg-stone-50">
    <!-- Messages Area -->
    <div bind:this={chatContainer} class="flex-1 overflow-y-auto p-4 space-y-3">
        {#each $websocketStore.messages as msg}
            {#if msg.type === "EVENT"}
                <div class="flex justify-center my-2">
                    <div
                        class="bg-stone-200 text-stone-600 text-xs px-3 py-1 rounded-full italic border border-stone-300"
                    >
                        <span class="font-bold">{msg.sender_name}</span>
                        {msg.content}
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

                        {msg.content}
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
            <input
                type="text"
                bind:value={newMessage}
                onkeydown={(e) => e.key === "Enter" && handleSendMessage()}
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
