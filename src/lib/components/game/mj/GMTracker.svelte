<script lang="ts">
    import { activeBoardStore } from "$lib/websocket";
    import { supabase } from "$lib/supabaseClient";
    import {
        fetchBoardTokens,
        updateBoardCombatState,
        updateTokenCombatState,
        resetAllTokensCombatState,
        type BoardToken
    } from "$lib/api/board";
    import { updateCharacterHP, updateCharacterConditions, parseConditions } from "$lib/api/character";
    import { parseDiceAndMath } from "$lib/utils/diceParser";
    import { sendMessage } from "$lib/chat";
    import { page } from "$app/state";
    import {
        Play,
        Trash2,
        Eye,
        EyeOff,
        Plus,
        RotateCcw,
        ChevronRight,
        ChevronLeft,
        Shield,
        Heart,
        Activity,
        Flame,
        Zap,
        Skull,
        Moon,
        Ghost,
        X
    } from "lucide-svelte";

    const gameId = page.params.id || "";

    let tokens = $state<BoardToken[]>([]);
    let loading = $state(false);
    let hpInputs = $state<Record<string, string>>({});
    let openPopoverTokenId = $state<string | null>(null);
    let customConditionText = $state("");

    // Predefined status effects/conditions
    const AVAILABLE_CONDITIONS = [
        { id: "poisoned", label: "Empoisonné", color: "bg-emerald-500 text-white", icon: Skull },
        { id: "stunned", label: "Étourdi", color: "bg-purple-500 text-white", icon: Zap },
        { id: "prone", label: "À terre", color: "bg-amber-500 text-stone-900", icon: RotateCcw },
        { id: "blessed", label: "Béni", color: "bg-yellow-400 text-stone-900", icon: Flame },
        { id: "invisible", label: "Invisible", color: "bg-indigo-500 text-white", icon: Ghost },
        { id: "unconscious", label: "Inconscient", color: "bg-rose-700 text-white", icon: Moon }
    ];

    const activeBoard = $derived($activeBoardStore);
    const activeBoardId = $derived(activeBoard?.id || "");
    const combatActive = $derived(activeBoard?.combat_active || false);
    const combatRound = $derived(activeBoard?.combat_round || 1);
    const activeTokenId = $derived(activeBoard?.combat_active_token_id || null);

    // Derived list of tokens participating in combat, sorted by initiative
    const combatants = $derived(
        tokens
            .filter(t => t.in_combat)
            .sort((a, b) => {
                if (b.initiative !== a.initiative) {
                    return b.initiative - a.initiative;
                }
                // Tie breaker by name
                return (a.character?.name || "").localeCompare(b.character?.name || "");
            })
    );

    // Derived list of tokens NOT participating in combat
    const nonCombatants = $derived(
        tokens.filter(t => !t.in_combat)
    );

    async function loadTokens() {
        if (!activeBoardId) return;
        try {
            loading = true;
            tokens = await fetchBoardTokens(activeBoardId);
            // Initialize HP inputs
            tokens.forEach(t => {
                if (t.character && hpInputs[t.id] === undefined) {
                    hpInputs[t.id] = String(t.character.current_hp);
                }
            });
        } catch (e) {
            console.error("Error loading board tokens for combat tracker:", e);
        } finally {
            loading = false;
        }
    }

    // Realtime listeners
    let tokensChannel: any = null;
    let charactersChannel: any = null;

    $effect(() => {
        if (activeBoardId) {
            loadTokens();

            // Listen to board tokens changes
            if (tokensChannel) supabase.removeChannel(tokensChannel);
            tokensChannel = supabase.channel(`combat_tokens:${activeBoardId}`)
                .on(
                    'postgres_changes',
                    { event: '*', schema: 'public', table: 'board_tokens', filter: `board_id=eq.${activeBoardId}` },
                    () => {
                        loadTokens();
                    }
                )
                .subscribe();

            // Listen to characters changes in this game
            if (charactersChannel) supabase.removeChannel(charactersChannel);
            charactersChannel = supabase.channel(`combat_characters:${gameId}`)
                .on(
                    'postgres_changes',
                    { event: 'UPDATE', schema: 'public', table: 'characters' },
                    (payload) => {
                        const updatedChar = payload.new;
                        tokens = tokens.map(t => {
                            if (t.character_id === updatedChar.id) {
                                return {
                                    ...t,
                                    character: t.character ? {
                                        ...t.character,
                                        ...updatedChar,
                                        conditions: parseConditions(updatedChar.conditions)
                                    } : undefined
                                };
                            }
                            return t;
                        });
                        // Update inputs if not focused
                        tokens.forEach(t => {
                            if (t.character_id === updatedChar.id) {
                                const inputEl = document.getElementById(`hp-input-${t.id}`);
                                if (inputEl !== document.activeElement) {
                                    hpInputs[t.id] = String(updatedChar.current_hp);
                                }
                            }
                        });
                    }
                )
                .subscribe();
        } else {
            tokens = [];
        }

        return () => {
            if (tokensChannel) supabase.removeChannel(tokensChannel);
            if (charactersChannel) supabase.removeChannel(charactersChannel);
        };
    });

    // Toggle participant status in combat
    async function toggleCombatant(token: BoardToken) {
        try {
            const newInCombat = !token.in_combat;
            let initVal = token.initiative;

            // Roll initiative automatically for NPCs/monsters when added
            if (newInCombat && token.character?.is_npc) {
                const mod = token.character.initiative || 0;
                initVal = Math.floor(Math.random() * 20) + 1 + mod;
            }

            await updateTokenCombatState(token.id, {
                in_combat: newInCombat,
                initiative: initVal
            });
        } catch (e) {
            console.error("Failed to toggle combatant:", e);
        }
    }

    // Toggle visibility of a token
    async function toggleVisibility(token: BoardToken) {
        try {
            await updateTokenCombatState(token.id, {
                is_hidden: !token.is_hidden
            });
        } catch (e) {
            console.error("Failed to toggle visibility:", e);
        }
    }

    // Update initiative score
    async function changeInitiative(tokenId: string, score: number) {
        try {
            await updateTokenCombatState(tokenId, { initiative: score });
        } catch (e) {
            console.error("Failed to update initiative:", e);
        }
    }

    // Cycle faction: ally -> enemy -> neutral -> ally
    async function cycleFaction(token: BoardToken) {
        const factions: ('ally' | 'enemy' | 'neutral')[] = ['ally', 'enemy', 'neutral'];
        const nextIdx = (factions.indexOf(token.faction) + 1) % factions.length;
        const newFaction = factions[nextIdx];
        try {
            await updateTokenCombatState(token.id, { faction: newFaction });
        } catch (e) {
            console.error("Failed to update token faction:", e);
        }
    }

    // Smart HP Modifier function
    async function handleHpSubmit(token: BoardToken, event: KeyboardEvent) {
        if (event.key !== "Enter") return;

        const inputStr = hpInputs[token.id];
        if (!token.character) return;

        const currentHp = token.character.current_hp;
        const maxHp = token.character.max_hp;
        const charName = token.character.name;

        const parseResult = parseDiceAndMath(inputStr);
        if (!parseResult) {
            // Revert on failure
            hpInputs[token.id] = String(currentHp);
            return;
        }

        let newHp = currentHp;
        let actionMsg = "";

        if (parseResult.type === "add") {
            newHp = Math.min(currentHp + parseResult.evaluatedTotal, maxHp);
            actionMsg = `💖 **${charName}** : Soigne \`${parseResult.logText}\` PV (Nouveau total : **${newHp}/${maxHp}** PV)`;
        } else if (parseResult.type === "subtract") {
            newHp = Math.max(currentHp - parseResult.evaluatedTotal, 0);
            actionMsg = `⚔️ **${charName}** : Soustrait \`${parseResult.logText}\` PV (Nouveau total : **${newHp}/${maxHp}** PV)`;
        } else {
            newHp = Math.max(Math.min(parseResult.evaluatedTotal, maxHp), 0);
            actionMsg = `⚙️ **${charName}** : PV définis à \`${parseResult.logText}\` (Nouveau total : **${newHp}/${maxHp}** PV)`;
        }

        try {
            // Update HP in DB
            await updateCharacterHP(token.character_id, newHp);
            hpInputs[token.id] = String(newHp);

            // Log event to chat
            await sendMessage({
                game_id: gameId,
                content: actionMsg,
                type: "EVENT",
                sender_name: "Combat"
            });
        } catch (e) {
            console.error("Failed to update HP:", e);
            hpInputs[token.id] = String(currentHp);
        }
    }

    // Add / Remove conditions
    async function toggleCondition(token: BoardToken, conditionId: string) {
        if (!token.character) return;
        const activeConditions = parseConditions(token.character.conditions);
        let newConditions = [];

        if (activeConditions.includes(conditionId)) {
            newConditions = activeConditions.filter(c => c !== conditionId);
        } else {
            newConditions = [...activeConditions, conditionId];
        }

        try {
            await updateCharacterConditions(token.character_id, newConditions);
        } catch (e) {
            console.error("Failed to update conditions:", e);
        }
    }

    async function addCustomCondition(token: BoardToken) {
        if (!token.character) return;
        const trimmed = customConditionText.trim();
        if (!trimmed) {
            openPopoverTokenId = null;
            return;
        }

        const activeConditions = parseConditions(token.character.conditions);
        if (!activeConditions.includes(trimmed)) {
            try {
                await updateCharacterConditions(token.character_id, [...activeConditions, trimmed]);
                customConditionText = "";
            } catch (e) {
                console.error("Failed to add custom condition:", e);
            }
        }
        openPopoverTokenId = null;
    }

    // Start combat session
    async function startCombat() {
        if (combatants.length === 0) return;
        try {
            const firstTokenId = combatants[0].id;
            await updateBoardCombatState(activeBoardId, {
                combat_active: true,
                combat_round: 1,
                combat_active_token_id: firstTokenId
            });

            await sendMessage({
                game_id: gameId,
                content: `⚔️ **Le combat commence !** Tour de **${combatants[0].character?.name}** (Round 1)`,
                type: "EVENT",
                sender_name: "Combat"
            });
        } catch (e) {
            console.error("Failed to start combat:", e);
        }
    }

    // Next turn
    async function nextTurn() {
        if (combatants.length === 0) return;
        const currentIndex = combatants.findIndex(c => c.id === activeTokenId);
        let nextIndex = currentIndex + 1;
        let newRound = combatRound;

        if (nextIndex >= combatants.length || currentIndex === -1) {
            nextIndex = 0;
            newRound += 1;
        }

        const nextToken = combatants[nextIndex];
        try {
            await updateBoardCombatState(activeBoardId, {
                combat_round: newRound,
                combat_active_token_id: nextToken.id
            });

            await sendMessage({
                game_id: gameId,
                content: `➡️ **Round ${newRound}** : Tour de **${nextToken.character?.name}**`,
                type: "EVENT",
                sender_name: "Combat"
            });
        } catch (e) {
            console.error("Failed to advance turn:", e);
        }
    }

    // Previous turn
    async function previousTurn() {
        if (combatants.length === 0) return;
        const currentIndex = combatants.findIndex(c => c.id === activeTokenId);
        let prevIndex = currentIndex - 1;
        let newRound = combatRound;

        if (prevIndex < 0) {
            prevIndex = combatants.length - 1;
            newRound = Math.max(1, combatRound - 1);
        }

        const prevToken = combatants[prevIndex];
        try {
            await updateBoardCombatState(activeBoardId, {
                combat_round: newRound,
                combat_active_token_id: prevToken.id
            });
        } catch (e) {
            console.error("Failed to step back turn:", e);
        }
    }

    // Stop combat session
    async function stopCombat() {
        try {
            await updateBoardCombatState(activeBoardId, {
                combat_active: false,
                combat_round: 1,
                combat_active_token_id: null
            });
            await resetAllTokensCombatState(activeBoardId);

            await sendMessage({
                game_id: gameId,
                content: `🏁 **Le combat est terminé !**`,
                type: "EVENT",
                sender_name: "Combat"
            });
        } catch (e) {
            console.error("Failed to stop combat:", e);
        }
    }

    // Toggle monster stats visibility option
    async function toggleHideMonsterStats() {
        try {
            await updateBoardCombatState(activeBoardId, {
                hide_monster_stats: !activeBoard.hide_monster_stats
            });
        } catch (e) {
            console.error("Failed to toggle hide monster stats:", e);
        }
    }
</script>

<div class="h-full flex flex-col bg-stone-900 text-stone-200 border-r border-stone-800 font-sans">
    <!-- Header -->
    <div class="p-4 border-b border-stone-800 flex items-center justify-between bg-stone-950">
        <div class="flex items-center gap-2">
            <Activity size={18} class="text-burnt-orange animate-pulse" />
            <span class="font-bold text-xs tracking-wider uppercase text-stone-400">
                Combat & Initiative
            </span>
        </div>
        {#if combatActive}
            <span class="px-2 py-0.5 rounded text-[10px] font-bold bg-burnt-orange text-white uppercase tracking-wide">
                En Cours
            </span>
        {/if}
    </div>

    <!-- Active State: IN COMBAT -->
    {#if combatActive}
        <div class="p-3 bg-stone-950/50 border-b border-stone-800 flex items-center justify-between">
            <div class="flex flex-col">
                <span class="text-xs text-stone-500 font-bold uppercase tracking-wider">Round</span>
                <span class="text-xl font-black text-burnt-orange leading-none">{combatRound}</span>
            </div>
            <div class="flex items-center gap-1.5">
                <button
                    onclick={previousTurn}
                    class="p-2 rounded-lg bg-stone-800 text-stone-300 hover:bg-stone-700 hover:text-white transition-colors cursor-pointer"
                    title="Tour précédent"
                >
                    <ChevronLeft size={16} />
                </button>
                <button
                    onclick={nextTurn}
                    class="px-4 py-2 rounded-lg bg-burnt-orange text-white font-bold text-xs hover:bg-burnt-orange-dark transition-colors flex items-center gap-1.5 shadow-md shadow-burnt-orange/10 cursor-pointer"
                    title="Tour suivant"
                >
                    <span>Tour Suivant</span>
                    <ChevronRight size={14} />
                </button>
                <button
                    onclick={stopCombat}
                    class="p-2 rounded-lg bg-red-950/60 text-red-400 border border-red-900/50 hover:bg-red-900 hover:text-white transition-colors ml-2 cursor-pointer"
                    title="Terminer le combat"
                >
                    <X size={16} />
                </button>
            </div>
        </div>

        <!-- Combatant list -->
        <div class="flex-1 overflow-y-auto p-4 space-y-3">
            {#each combatants as combatant, index (combatant.id)}
                {@const isCurrentTurn = combatant.id === activeTokenId}
                {@const char = combatant.character}
                <div
                    class="group relative rounded-xl border p-4 transition-all duration-300 flex flex-col gap-3.5
                    {isCurrentTurn 
                        ? 'bg-stone-800 border-burnt-orange/50 shadow-md shadow-burnt-orange/5' 
                        : 'bg-stone-900/50 border-stone-800/80 hover:border-stone-700'}
                    {openPopoverTokenId === combatant.id ? 'z-40' : 'z-10'}"
                >
                    <!-- Header Row -->
                    <div class="flex items-center justify-between gap-3">
                        <div class="flex items-center gap-3 min-w-0">
                            <!-- Avatar / Initiative circle -->
                            <div class="relative shrink-0">
                                {#if char?.avatar_url}
                                    <img src={char.avatar_url} alt="" class="w-10 h-10 rounded-full border border-stone-700 object-cover shadow-sm" />
                                {:else}
                                    <div class="w-10 h-10 rounded-full bg-stone-800 border border-stone-700 flex items-center justify-center font-bold text-sm text-stone-300">
                                        {(char?.name || "?").charAt(0).toUpperCase()}
                                    </div>
                                {/if}
                                {#if isCurrentTurn}
                                    <span class="absolute -top-1 -right-1 w-2.5 h-2.5 rounded-full bg-burnt-orange border-2 border-stone-900 animate-ping"></span>
                                {/if}
                            </div>

                            <div class="flex flex-col min-w-0 leading-normal">
                                <span class="font-bold text-sm text-stone-100 truncate group-hover:text-white transition-colors {char && char.current_hp <= 0 ? 'line-through text-stone-500' : ''}">
                                    {char?.name || "Sans nom"}
                                </span>
                                <div class="flex items-center gap-1.5 mt-0.5">
                                    <span class="px-1.5 py-0.2 rounded bg-stone-800/80 border border-stone-700 text-[8px] font-bold text-stone-400 uppercase tracking-wider">
                                        {char?.is_npc ? "PNJ" : "PJ"}
                                    </span>
                                    <button
                                        onclick={() => cycleFaction(combatant)}
                                        class="px-1.5 py-0.5 rounded text-[8px] font-black uppercase transition-colors tracking-wider cursor-pointer
                                        {combatant.faction === 'ally' ? 'bg-emerald-950/60 text-emerald-400 border border-emerald-900/40' : ''}
                                        {combatant.faction === 'enemy' ? 'bg-rose-950/60 text-rose-400 border border-rose-900/40' : ''}
                                        {combatant.faction === 'neutral' ? 'bg-amber-950/60 text-amber-400 border border-amber-900/40' : ''}"
                                        title="Changer de camp (Allié / Ennemi / Neutre)"
                                    >
                                        {combatant.faction === 'ally' ? 'Allié' : ''}
                                        {combatant.faction === 'enemy' ? 'Ennemi' : ''}
                                        {combatant.faction === 'neutral' ? 'Neutre' : ''}
                                    </button>
                                </div>
                            </div>
                        </div>

                        <!-- Top Right actions (Visibility & Exclude) -->
                        <div class="flex items-center gap-1 shrink-0">
                            <!-- Visibility toggle -->
                            <button
                                onclick={() => toggleVisibility(combatant)}
                                class="p-1 rounded text-stone-500 hover:text-stone-300 hover:bg-stone-800/60 transition-all cursor-pointer"
                                title={combatant.is_hidden ? "Masqué pour les joueurs" : "Visible pour les joueurs"}
                            >
                                {#if combatant.is_hidden}
                                    <EyeOff size={14} class="text-rose-500" />
                                {:else}
                                    <Eye size={14} class="text-emerald-500" />
                                	{/if}
                            </button>

                            <!-- Exclude -->
                            <button
                                onclick={() => toggleCombatant(combatant)}
                                class="p-1 rounded text-stone-500 hover:text-red-400 hover:bg-stone-800/60 transition-all cursor-pointer"
                                title="Retirer du combat"
                            >
                                <Trash2 size={14} />
                            </button>
                        </div>
                    </div>

                    <!-- Health Bar Row -->
                    {#if char}
                        {@const hpPercent = Math.max(0, Math.min(100, (char.current_hp / char.max_hp) * 100))}
                        <div class="w-full bg-stone-950 h-1.5 rounded-full overflow-hidden border border-stone-800/80 shadow-inner">
                            <div 
                                class="h-full rounded-full transition-all duration-300 
                                {hpPercent > 50 ? 'bg-emerald-500/90' : hpPercent > 20 ? 'bg-amber-500/90' : 'bg-rose-500/90'}" 
                                style="width: {hpPercent}%"
                            ></div>
                        </div>
                    {/if}

                    <!-- Stats & Inputs Row -->
                    <div class="flex flex-wrap items-center justify-between gap-3 pt-1">
                        <!-- CA & VIT Badges -->
                        <div class="flex items-center gap-1.5 text-[10px] text-stone-400 font-medium">
                            {#if char}
                                <span class="flex items-center gap-1 bg-stone-950/60 border border-stone-800 px-2 py-1 rounded-md" title="Classe d'armure">
                                    <Shield size={11} class="text-stone-400" />
                                    <span>CA: <span class="font-bold text-stone-200">{char.armor_class || 10}</span></span>
                                </span>
                                <span class="flex items-center gap-1 bg-stone-950/40 border border-stone-800 px-2 py-1 rounded-md" title="Vitesse de déplacement">
                                    <RotateCcw size={11} class="text-stone-400 rotate-90" />
                                    <span>VIT: <span class="font-bold text-stone-200">{char.speed || 30}m</span></span>
                                </span>
                            {/if}
                        </div>

                        <!-- HP & Initiative inputs -->
                        <div class="flex items-center gap-2 shrink-0">
                            <!-- HP input -->
                            {#if char}
                                <div class="flex items-center gap-1 bg-stone-950 border border-stone-800 rounded-lg px-2.5 py-1 shadow-inner" title="PV: Entrez la valeur (ex: 20) ou le modificateur (ex: -5, +2d6)">
                                    <Heart size={11} class="text-rose-500 fill-rose-500/10 shrink-0" />
                                    <input
                                        type="text"
                                        id="hp-input-{combatant.id}"
                                        bind:value={hpInputs[combatant.id]}
                                        onkeypress={(e) => handleHpSubmit(combatant, e)}
                                        class="w-12 bg-transparent text-center font-black text-xs border-none outline-none focus:text-burnt-orange text-stone-200"
                                    />
                                    <span class="text-[10px] text-stone-600 font-bold">/</span>
                                    <span class="text-[10px] text-stone-400 font-bold pr-0.5">{char.max_hp}</span>
                                </div>
                            {/if}

                            <!-- Initiative input -->
                            <div class="flex items-center gap-1 bg-stone-950 border border-stone-800 rounded-lg px-2 py-1 shadow-inner">
                                <span class="text-[9px] font-bold text-stone-500 uppercase tracking-wide">Init</span>
                                <input
                                    type="number"
                                    value={combatant.initiative}
                                    onchange={(e) => changeInitiative(combatant.id, parseInt(e.currentTarget.value, 10))}
                                    class="w-8 bg-transparent font-black text-center text-xs border-none outline-none focus:text-burnt-orange text-stone-200"
                                />
                            </div>
                        </div>
                    </div>

                    <!-- Conditions Row -->
                    <div class="flex flex-col gap-2 pt-2 border-t border-stone-800/50">
                        <div class="flex items-center justify-between gap-2">
                            <span class="text-[10px] font-bold text-stone-400 uppercase tracking-wider">États & Effets</span>
                            
                            <!-- Conditions popover trigger -->
                            <div class="relative">
                                <button
                                    onclick={() => openPopoverTokenId = openPopoverTokenId === combatant.id ? null : combatant.id}
                                    class="px-2 py-0.5 rounded text-[10px] font-bold bg-stone-800 hover:bg-stone-700 text-stone-300 border border-stone-700 hover:text-white transition-colors cursor-pointer flex items-center gap-1"
                                >
                                    <Plus size={10} />
                                    <span>Gérer</span>
                                </button>

                                {#if openPopoverTokenId === combatant.id}
                                    {@const isLowerHalf = index >= Math.ceil(combatants.length / 2)}
                                    <!-- Popover -->
                                    <div class="absolute right-0 {isLowerHalf ? 'bottom-full mb-2' : 'top-full mt-2'} w-52 bg-stone-950 border border-stone-700/80 rounded-xl p-3 shadow-2xl shadow-black/80 z-50 flex flex-col gap-2.5">
                                        <div class="flex justify-between items-center pb-1.5 border-b border-stone-800">
                                            <span class="text-[10px] font-bold text-stone-400 uppercase tracking-wider">Gérer les états</span>
                                            <button onclick={() => openPopoverTokenId = null} class="text-stone-400 hover:text-stone-100 cursor-pointer p-0.5 rounded hover:bg-stone-800 transition-colors">
                                                <X size={12} />
                                            </button>
                                        </div>

                                        <div class="grid grid-cols-1 gap-1.5">
                                            {#each AVAILABLE_CONDITIONS as cond}
                                                {@const conds = parseConditions(char?.conditions)}
                                                {@const active = conds.includes(cond.id)}
                                                <button
                                                    onclick={() => toggleCondition(combatant, cond.id)}
                                                    class="flex items-center justify-between p-1.5 px-2.5 rounded-lg text-left text-xs font-semibold transition-colors cursor-pointer
                                                    {active ? 'bg-burnt-orange/15 text-burnt-orange border border-burnt-orange/30' : 'hover:bg-stone-900 text-stone-300 border border-transparent'}"
                                                >
                                                    <div class="flex items-center gap-2">
                                                        <cond.icon size={12} class={active ? "text-burnt-orange" : "text-stone-500"} />
                                                        <span>{cond.label}</span>
                                                    </div>
                                                    {#if active}
                                                        <span class="w-1.5 h-1.5 rounded-full bg-burnt-orange"></span>
                                                    {/if}
                                                </button>
                                            {/each}
                                        </div>

                                        <!-- Custom condition -->
                                        <div class="pt-2 border-t border-stone-800 flex flex-col gap-1">
                                            <span class="text-[9px] text-stone-500 font-bold uppercase">Condition Libre</span>
                                            <div class="flex gap-1">
                                                <input
                                                    type="text"
                                                    placeholder="ex: Effrayé"
                                                    bind:value={customConditionText}
                                                    class="flex-1 bg-stone-900 border border-stone-800 rounded px-1.5 py-0.5 text-[10px] text-white focus:outline-none focus:border-burnt-orange"
                                                />
                                                <button
                                                    onclick={() => addCustomCondition(combatant)}
                                                    class="px-1.5 py-0.5 rounded bg-burnt-orange text-white text-[10px] font-bold hover:bg-burnt-orange-dark cursor-pointer"
                                                >
                                                    OK
                                                </button>
                                            </div>
                                        </div>
                                    </div>
                                {/if}
                            </div>
                        </div>

                        <!-- Active conditions list -->
                        {#if char}
                            {@const conds = parseConditions(char.conditions)}
                            {#if conds.length > 0}
                                <div class="flex flex-wrap gap-1.5 mt-0.5">
                                    {#each conds as condId}
                                        {@const predefined = AVAILABLE_CONDITIONS.find(c => c.id === condId)}
                                        {#if predefined}
                                            <span class="px-2 py-0.5 rounded-full text-[9px] font-bold flex items-center gap-1 shadow-sm border {predefined.color}">
                                                <predefined.icon size={10} />
                                                {predefined.label}
                                            </span>
                                        {:else}
                                            <!-- Custom condition tag -->
                                            <span class="px-2 py-0.5 rounded-full text-[9px] font-bold flex items-center gap-1 bg-stone-800 text-stone-300 border border-stone-700 shadow-sm">
                                                <Activity size={10} class="text-stone-400" />
                                                <span>{condId}</span>
                                                <button 
                                                    onclick={() => toggleCondition(combatant, condId)}
                                                    class="ml-0.5 hover:text-red-400 transition-colors cursor-pointer"
                                                    title="Retirer la condition"
                                                >
                                                    <X size={8} />
                                                </button>
                                            </span>
                                        {/if}
                                    {/each}
                                </div>
                            {:else}
                                <span class="text-[10px] text-stone-600 italic">Aucune condition active</span>
                            {/if}
                        {/if}
                    </div>
                </div>
            {/each}
        </div>
    {:else}
        <!-- Setup / Inactive State -->
        <div class="p-4 bg-stone-950/20 border-b border-stone-800 flex flex-col gap-2">
            <span class="text-xs font-bold text-stone-400">Préparation du combat</span>
            <p class="text-[11px] text-stone-500">
                Sélectionnez les jetons de la carte à inclure dans le combat. Les initiatives des monstres seront tirées automatiquement.
            </p>
            <div class="flex items-center gap-2 pt-1">
                <button
                    onclick={startCombat}
                    disabled={combatants.length === 0}
                    class="flex-1 py-2 rounded-lg bg-burnt-orange text-white font-bold text-xs hover:bg-burnt-orange-dark transition-colors flex items-center justify-center gap-2 disabled:opacity-40 disabled:cursor-not-allowed shadow-md"
                >
                    <Play size={12} fill="currentColor" />
                    Commencer le Combat ({combatants.length})
                </button>
            </div>
            <!-- Toggle Hide Monster Stats -->
            <label class="flex items-center gap-2 pt-2 border-t border-stone-850 cursor-pointer select-none">
                <input
                    type="checkbox"
                    checked={activeBoard?.hide_monster_stats || false}
                    onchange={toggleHideMonsterStats}
                    class="rounded border-stone-700 bg-stone-900 text-burnt-orange focus:ring-0 focus:ring-offset-0 w-3.5 h-3.5"
                />
                <span class="text-[10px] text-stone-400 font-medium">Masquer complètement la santé des monstres</span>
            </label>
        </div>

        <!-- Scrollable list of board tokens to check -->
        <div class="flex-1 overflow-y-auto p-3 space-y-3">
            <!-- Active participants -->
            <div>
                <span class="text-[10px] text-stone-500 font-bold uppercase tracking-wider block mb-2">Combattants ({combatants.length})</span>
                {#if combatants.length === 0}
                    <div class="p-4 rounded-xl border border-dashed border-stone-800 text-center text-stone-600 text-xs">
                        Aucun participant sélectionné.
                    </div>
                {:else}
                    <div class="space-y-2.5">
                        {#each combatants as token (token.id)}
                            <div class="flex items-center justify-between p-3.5 rounded-xl bg-stone-850/40 border border-stone-800 gap-3">
                                <div class="flex items-center gap-3 min-w-0 flex-1">
                                    <input
                                        type="checkbox"
                                        checked={true}
                                        onchange={() => toggleCombatant(token)}
                                        class="rounded border-stone-700 bg-stone-900 text-burnt-orange focus:ring-0 focus:ring-offset-0 cursor-pointer w-4 h-4"
                                    />
                                    {#if token.character?.avatar_url}
                                        <img src={token.character.avatar_url} alt="" class="w-8 h-8 rounded-full object-cover shrink-0" />
                                    {:else}
                                        <div class="w-8 h-8 rounded-full bg-stone-800 flex items-center justify-center font-bold text-xs shrink-0 text-stone-300">
                                            {(token.character?.name || "?").charAt(0).toUpperCase()}
                                        </div>
                                    {/if}
                                    <div class="flex flex-col min-w-0 flex-1 gap-0.5">
                                        <span class="text-xs text-stone-300 font-bold truncate">
                                            {token.character?.name}
                                        </span>
                                        <div>
                                            <button
                                                onclick={() => cycleFaction(token)}
                                                class="px-1.5 py-0.5 rounded text-[8px] font-black uppercase transition-colors inline-block tracking-wider
                                                {token.faction === 'ally' ? 'bg-emerald-950/60 text-emerald-400 border border-emerald-900/40' : ''}
                                                {token.faction === 'enemy' ? 'bg-rose-950/60 text-rose-400 border border-rose-900/40' : ''}
                                                {token.faction === 'neutral' ? 'bg-amber-950/60 text-amber-400 border border-amber-900/40' : ''}"
                                                title="Changer de camp"
                                            >
                                                {token.faction === 'ally' ? 'Allié' : ''}
                                                {token.faction === 'enemy' ? 'Ennemi' : ''}
                                                {token.faction === 'neutral' ? 'Neutre' : ''}
                                            </button>
                                        </div>
                                    </div>
                                </div>
                                <div class="flex items-center gap-2 shrink-0">
                                    <span class="text-[9px] text-stone-500 font-bold uppercase tracking-wider">Init</span>
                                    <input
                                        type="number"
                                        value={token.initiative}
                                        onchange={(e) => changeInitiative(token.id, parseInt(e.currentTarget.value, 10))}
                                        class="w-10 bg-stone-950 border border-stone-800 rounded font-black text-center text-xs py-1 text-stone-200 focus:outline-none focus:border-burnt-orange"
                                    />
                                </div>
                            </div>
                        {/each}
                    </div>
                {/if}
            </div>

            <!-- Non-participants -->
            {#if nonCombatants.length > 0}
                <div>
                    <span class="text-[10px] text-stone-500 font-bold uppercase tracking-wider block mb-2">Autres jetons sur la carte ({nonCombatants.length})</span>
                    <div class="space-y-2.5">
                        {#each nonCombatants as token (token.id)}
                            <div
                                class="w-full flex items-center justify-between p-3.5 rounded-xl bg-stone-900/20 border border-stone-850 hover:bg-stone-850/30 hover:border-stone-800 transition-all group gap-3"
                            >
                                <button
                                    onclick={() => toggleCombatant(token)}
                                    class="flex items-center gap-3 min-w-0 text-left flex-1"
                                >
                                    <input
                                        type="checkbox"
                                        checked={false}
                                        class="rounded border-stone-700 bg-stone-900 pointer-events-none w-4 h-4"
                                    />
                                    {#if token.character?.avatar_url}
                                        <img src={token.character.avatar_url} alt="" class="w-8 h-8 rounded-full object-cover shrink-0 shadow-sm" />
                                    {:else}
                                        <div class="w-8 h-8 rounded-full bg-stone-800 flex items-center justify-center font-bold text-xs shrink-0 text-stone-400">
                                            {(token.character?.name || "?").charAt(0).toUpperCase()}
                                        </div>
                                    {/if}
                                    <div class="flex flex-col min-w-0">
                                        <span class="text-xs text-stone-350 group-hover:text-stone-100 transition-colors font-bold truncate">
                                            {token.character?.name}
                                        </span>
                                        <span class="text-[9px] text-stone-550 font-bold uppercase tracking-wider mt-0.5">
                                            {token.character?.is_npc ? "Monstre / PNJ" : "PJ"}
                                        </span>
                                    </div>
                                </button>
                                
                                <button
                                    onclick={() => cycleFaction(token)}
                                    class="px-1.5 py-0.5 rounded text-[8px] font-black uppercase transition-colors shrink-0 ml-2 tracking-wider
                                    {token.faction === 'ally' ? 'bg-emerald-950/60 text-emerald-400 border border-emerald-900/40' : ''}
                                    {token.faction === 'enemy' ? 'bg-rose-950/60 text-rose-400 border border-rose-900/40' : ''}
                                    {token.faction === 'neutral' ? 'bg-amber-950/60 text-amber-400 border border-amber-900/40' : ''}"
                                    title="Changer de camp"
                                >
                                    {token.faction === 'ally' ? 'Allié' : ''}
                                    {token.faction === 'enemy' ? 'Ennemi' : ''}
                                    {token.faction === 'neutral' ? 'Neutre' : ''}
                                </button>
                            </div>
                        {/each}
                    </div>
                </div>
            {/if}
        </div>
    {/if}
</div>
