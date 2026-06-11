<script lang="ts">
    import Scene from "$lib/components/game/Scene.svelte";
    import { activeBoardStore } from "$lib/websocket";
    import { fetchBoardTokens, type BoardToken } from "$lib/api/board";
    import { supabase } from "$lib/supabaseClient";
    import { page } from "$app/state";
    import { X, Eye, EyeOff } from "lucide-svelte";

    let tokens = $state<BoardToken[]>([]);
    let loading = $state(true);
    let showInitiative = $state(true);

    const activeBoardId = $derived($activeBoardStore?.id || "");
    const combatActive = $derived($activeBoardStore?.combat_active || false);
    const activeTokenId = $derived($activeBoardStore?.combat_active_token_id || null);

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
        if (!activeBoardId) {
            tokens = [];
            loading = false;
            return;
        }
        try {
            loading = true;
            tokens = await fetchBoardTokens(activeBoardId);
        } catch (e) {
            console.error("Error loading board tokens for presentation:", e);
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
            tokensChannel = supabase.channel(`presentation_tokens:${activeBoardId}`)
                .on(
                    'postgres_changes',
                    { event: '*', schema: 'public', table: 'board_tokens', filter: `board_id=eq.${activeBoardId}` },
                    () => {
                        loadTokens();
                    }
                )
                .subscribe();

            // Realtime sync for character modifications
            if (charactersChannel) supabase.removeChannel(charactersChannel);
            charactersChannel = supabase.channel(`presentation_characters:${page.params.id}`)
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
                                        ...updatedChar
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
</script>

<div class="relative w-screen h-screen bg-stone-950 overflow-hidden select-none">
    <!-- Red glowing vignette for combat mode -->
    {#if combatActive}
        <div class="absolute inset-0 z-10 pointer-events-none border-red-glow"></div>
    {/if}

    <!-- The VTT Scene -->
    <div class="w-full h-full">
        <Scene isGM={false} />
    </div>

    <!-- Floating Top Bar for Initiative Order -->
    {#if combatActive && showInitiative && combatants.length > 0}
        <div class="absolute top-6 left-1/2 -translate-x-1/2 z-20 flex flex-col items-center">
            <div class="bg-stone-900/90 backdrop-blur-md px-6 py-3 rounded-2xl border border-stone-850 shadow-2xl flex items-center gap-5 animate-in slide-in-from-top-4 duration-300">
                {#each combatants as combatant (combatant.id)}
                    {@const isCurrentTurn = combatant.id === activeTokenId}
                    {@const char = combatant.character}
                    
                    <div class="flex flex-col items-center gap-1.5 relative transition-all duration-300">
                        <div class="relative">
                            {#if char?.avatar_url}
                                <img
                                    src={char.avatar_url}
                                    alt={char.name}
                                    class="w-12 h-12 rounded-full object-cover border-2 transition-all duration-300
                                    {isCurrentTurn 
                                        ? 'border-amber-400 scale-110 shadow-lg shadow-amber-500/20 ring-4 ring-amber-400/20 active-turn-pulse' 
                                        : 'border-stone-700'}"
                                />
                            {:else}
                                <div class="w-12 h-12 rounded-full bg-stone-800 border-2 flex items-center justify-center font-bold text-stone-350 transition-all duration-300
                                    {isCurrentTurn 
                                        ? 'border-amber-400 scale-110 ring-4 ring-amber-400/20' 
                                        : 'border-stone-700'}"
                                >
                                    {(char?.name || "?").charAt(0)}
                                </div>
                            {/if}
                            
                            {#if isCurrentTurn}
                                <span class="absolute -top-1 -right-1 w-4 h-4 bg-amber-400 border border-stone-900 rounded-full flex items-center justify-center text-[8px] font-black text-stone-950">
                                    ★
                                </span>
                            {/if}
                        </div>
                        <span class="text-[10px] font-medium tracking-wide truncate max-w-[75px] text-center
                            {isCurrentTurn ? 'text-amber-300 font-bold' : 'text-stone-400'}"
                        >
                            {char?.name || "Inconnu"}
                        </span>
                    </div>
                {/each}
            </div>
        </div>
    {/if}

    <!-- Floating Actions Panel (discreet controls) -->
    <div class="absolute top-4 left-4 z-20 flex gap-2">
        <button
            onclick={() => window.close()}
            class="bg-stone-900/80 hover:bg-stone-900 border border-stone-800 text-stone-400 hover:text-white p-2 rounded-xl transition-all shadow-lg backdrop-blur-sm cursor-pointer"
            title="Quitter le mode présentation"
        >
            <X size={16} />
        </button>
        
        {#if combatActive && combatants.length > 0}
            <button
                onclick={() => showInitiative = !showInitiative}
                class="bg-stone-900/80 hover:bg-stone-900 border border-stone-800 text-stone-400 hover:text-white p-2 rounded-xl transition-all shadow-lg backdrop-blur-sm cursor-pointer"
                title={showInitiative ? "Masquer l'initiative" : "Afficher l'initiative"}
            >
                {#if showInitiative}
                    <EyeOff size={16} />
                {:else}
                    <Eye size={16} />
                {/if}
            </button>
        {/if}
    </div>
</div>

<style>
    /* Internal red glowing vignette for combat */
    .border-red-glow {
        box-shadow: inset 0 0 50px rgba(239, 68, 68, 0.4);
        animation: pulse-red 2s infinite ease-in-out;
    }

    @keyframes pulse-red {
        0% {
            box-shadow: inset 0 0 50px rgba(239, 68, 68, 0.3);
        }
        50% {
            box-shadow: inset 0 0 75px rgba(239, 68, 68, 0.65);
        }
        100% {
            box-shadow: inset 0 0 50px rgba(239, 68, 68, 0.3);
        }
    }

    /* Active turn animation */
    :global(.active-turn-pulse) {
        animation: pulse-gold-ring 2s infinite ease-in-out;
    }

    @keyframes pulse-gold-ring {
        0%, 100% {
            transform: scale(1.1);
            box-shadow: 0 0 10px rgba(251, 191, 36, 0.3);
        }
        50% {
            transform: scale(1.15);
            box-shadow: 0 0 20px rgba(251, 191, 36, 0.6);
        }
    }
</style>
