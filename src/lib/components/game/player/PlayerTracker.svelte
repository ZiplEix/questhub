<script lang="ts">
    import { activeBoardStore } from "$lib/websocket";
    import { supabase } from "$lib/supabaseClient";
    import { fetchBoardTokens, type BoardToken } from "$lib/api/board";
    import { parseConditions } from "$lib/api/character";
    import { page } from "$app/state";
    import {
        Activity,
        Heart,
        Shield,
        Flame,
        Zap,
        Skull,
        Moon,
        Ghost,
        RotateCcw
    } from "lucide-svelte";

    const gameId = page.params.id || "";

    let tokens = $state<BoardToken[]>([]);
    let loading = $state(false);

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
    const hideMonsterStats = $derived(activeBoard?.hide_monster_stats || false);

    // Filter out hidden tokens and sort by initiative
    const combatants = $derived(
        tokens
            .filter(t => t.in_combat && !t.is_hidden)
            .sort((a, b) => {
                if (b.initiative !== a.initiative) {
                    return b.initiative - a.initiative;
                }
                return (a.character?.name || "").localeCompare(b.character?.name || "");
            })
    );

    async function loadTokens() {
        if (!activeBoardId) return;
        try {
            loading = true;
            tokens = await fetchBoardTokens(activeBoardId);
        } catch (e) {
            console.error("Error loading board tokens for player combat tracker:", e);
        } finally {
            loading = false;
        }
    }

    let tokensChannel: any = null;
    let charactersChannel: any = null;

    $effect(() => {
        if (activeBoardId) {
            loadTokens();

            // Realtime sync for board tokens
            if (tokensChannel) supabase.removeChannel(tokensChannel);
            tokensChannel = supabase.channel(`player_combat_tokens:${activeBoardId}`)
                .on(
                    'postgres_changes',
                    { event: '*', schema: 'public', table: 'board_tokens', filter: `board_id=eq.${activeBoardId}` },
                    () => {
                        loadTokens();
                    }
                )
                .subscribe();

            // Realtime sync for characters
            if (charactersChannel) supabase.removeChannel(charactersChannel);
            charactersChannel = supabase.channel(`player_combat_characters:${gameId}`)
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

    // Helper to determine monster health state
    function getMonsterHealthState(current: number, max: number): { label: string; color: string; width: string } {
        if (current <= 0) return { label: "K.O.", color: "bg-rose-600", width: "0%" };
        const ratio = current / max;
        if (ratio >= 0.99) return { label: "Indemne", color: "bg-emerald-500", width: "100%" };
        if (ratio >= 0.5) return { label: "Blessé", color: "bg-yellow-500", width: `${ratio * 100}%` };
        return { label: "Gravement Blessé", color: "bg-amber-600", width: `${ratio * 100}%` };
    }
</script>

<div class="h-full flex flex-col bg-stone-50 font-sans">
    {#if !combatActive}
        <div class="flex-1 flex flex-col items-center justify-center p-6 text-center text-stone-400 gap-3">
            <Activity size={32} class="text-stone-300" />
            <p class="text-sm font-semibold text-stone-500">Aucun combat actif</p>
            <p class="text-xs text-stone-400 max-w-xs">
                Le tracker de combat s'activera automatiquement lorsque le Maître de Jeu démarrera un combat.
            </p>
        </div>
    {:else}
        <!-- Active Combat Header -->
        <div class="p-4 border-b border-stone-100 bg-white flex items-center justify-between shadow-sm">
            <div class="flex flex-col">
                <span class="text-[10px] font-bold text-stone-400 uppercase tracking-wider">Round</span>
                <span class="text-lg font-black text-dark-gray leading-none">{combatRound}</span>
            </div>
            <div class="flex items-center gap-1 bg-stone-50 border border-stone-200/60 rounded-full px-3 py-1">
                <span class="w-1.5 h-1.5 rounded-full bg-burnt-orange animate-ping"></span>
                <span class="text-[10px] font-bold text-stone-600 uppercase tracking-wider">Combat en cours</span>
            </div>
        </div>

        <!-- Participants List -->
        <div class="flex-1 overflow-y-auto p-3 space-y-2.5">
            {#each combatants as combatant (combatant.id)}
                {@const isCurrentTurn = combatant.id === activeTokenId}
                {@const char = combatant.character}
                {@const isAlly = combatant.faction === 'ally'}
                
                <div
                    class="group relative rounded-xl border p-3 transition-all duration-300 flex flex-col gap-2.5 bg-white
                    {isCurrentTurn 
                        ? 'border-burnt-orange/50 shadow-md shadow-burnt-orange/5 bg-stone-50/50' 
                        : 'border-stone-100 hover:border-stone-200/80'}"
                >
                    <!-- Current Turn Marker Line -->
                    {#if isCurrentTurn}
                        <div class="absolute left-0 top-0 bottom-0 w-1 bg-burnt-orange rounded-l-xl"></div>
                    {/if}

                    <div class="flex items-center justify-between gap-3">
                        <div class="flex items-center gap-3 min-w-0">
                            <!-- Avatar -->
                            <div class="relative shrink-0">
                                {#if char?.avatar_url}
                                    <img src={char.avatar_url} alt="" class="w-9 h-9 rounded-full border border-stone-200 object-cover" />
                                {:else}
                                    <div class="w-9 h-9 rounded-full bg-stone-100 border border-stone-200 flex items-center justify-center font-bold text-sm text-stone-500">
                                        {(char?.name || "?").charAt(0)}
                                    </div>
                                {/if}
                                {#if isCurrentTurn}
                                    <span class="absolute -top-1 -right-1 w-2.5 h-2.5 rounded-full bg-burnt-orange border-2 border-white animate-ping"></span>
                                {/if}
                            </div>

                            <div class="flex flex-col min-w-0 leading-tight">
                                <span class="font-bold text-xs text-stone-800 truncate {char && char.current_hp <= 0 ? 'line-through text-stone-400' : ''}">
                                    {char?.name || "Inconnu"}
                                </span>
                                <span class="text-[9px] font-bold tracking-wide uppercase
                                    {combatant.faction === 'ally' ? 'text-emerald-500' : ''}
                                    {combatant.faction === 'enemy' ? 'text-rose-500' : ''}
                                    {combatant.faction === 'neutral' ? 'text-amber-500' : ''}"
                                >
                                    {combatant.faction === 'ally' ? 'Allié' : ''}
                                    {combatant.faction === 'enemy' ? 'Ennemi' : ''}
                                    {combatant.faction === 'neutral' ? 'Neutre' : ''}
                                </span>
                            </div>
                        </div>

                        <!-- HP display -->
                        {#if char}
                            <div class="shrink-0 flex flex-col items-end gap-1 min-w-[80px]">
                                {#if isAlly}
                                    <!-- Player HP: Exact Numbers -->
                                    <div class="flex items-center gap-1.5 text-stone-600">
                                        <Heart size={12} class="text-rose-500 fill-rose-500/10 shrink-0" />
                                        <span class="font-black text-xs text-stone-700">{char.current_hp}</span>
                                        <span class="text-[10px] text-stone-300 font-bold">/</span>
                                        <span class="text-[10px] text-stone-400 font-bold">{char.max_hp}</span>
                                    </div>
                                {:else if !hideMonsterStats}
                                    <!-- Monster HP: Status Bar (No exact numbers) -->
                                    {@const hState = getMonsterHealthState(char.current_hp, char.max_hp)}
                                    <div class="w-20 flex flex-col items-stretch gap-1">
                                        <div class="h-1.5 w-full bg-stone-100 rounded-full overflow-hidden border border-stone-200/50">
                                            <div class="h-full {hState.color} rounded-full" style="width: {hState.width}"></div>
                                        </div>
                                        <span class="text-[8px] font-bold text-stone-400 uppercase text-right leading-none">
                                            {hState.label}
                                        </span>
                                    </div>
                                {/if}
                            </div>
                        {/if}
                    </div>

                    <!-- Bottom detail: Armor Class & Conditions -->
                    <div class="flex flex-wrap items-center justify-between gap-1 pt-1.5 border-t border-stone-100">
                        <div class="flex items-center gap-3 text-[10px] text-stone-400 font-medium">
                            {#if char && isAlly}
                                <span class="flex items-center gap-1">
                                    <Shield size={10} class="text-stone-400" />
                                    CA: <span class="font-bold text-stone-600">{char.armor_class || 10}</span>
                                </span>
                            {/if}
                        </div>
                        <div class="flex items-center gap-1.5">
                            <span class="text-[9px] text-stone-300 font-bold uppercase">Initiative: {combatant.initiative}</span>
                        </div>
                    </div>

                    <!-- Active conditions list -->
                    {#if char}
                        {@const conds = parseConditions(char.conditions)}
                        {#if conds.length > 0}
                            <div class="flex flex-wrap gap-1 mt-1">
                                {#each conds as condId}
                                    {@const predefined = AVAILABLE_CONDITIONS.find(c => c.id === condId)}
                                    {#if predefined}
                                        <span class="px-1.5 py-0.5 rounded text-[9px] font-bold flex items-center gap-1 {predefined.color}">
                                            <predefined.icon size={8} />
                                            {predefined.label}
                                        </span>
                                    {:else}
                                        <!-- Custom condition tag -->
                                        <span class="px-1.5 py-0.5 rounded text-[9px] font-bold bg-stone-100 text-stone-500 border border-stone-200 flex items-center gap-1">
                                            <Activity size={8} class="text-stone-400" />
                                            {condId}
                                        </span>
                                    {/if}
                                {/each}
                            </div>
                        {/if}
                    {/if}
                </div>
            {/each}
        </div>
    {/if}
</div>
