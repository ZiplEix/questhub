<script lang="ts">
    import { onMount } from "svelte";
    import { page } from "$app/state";
    import {
        Search,
        Users,
        ArrowLeft,
        Pencil,
        Shield,
        Heart,
        UserCheck,
        UserX,
        MapPin
    } from "lucide-svelte";
    import { fetchCharacters, fetchPlayers, type Character, type Player } from "$lib/api";
    import { activeBoardStore } from "$lib/websocket";
    import { addBoardToken } from "$lib/api/board";
    import CharacterSheet from "../pupitre/CharacterSheet.svelte";
    import { supabase } from "$lib/supabaseClient";
    import { parseConditions } from "$lib/api/character";

    let characters = $state<Character[]>([]);
    let players = $state<Player[]>([]);
    let selectedCharacter = $state<Character | null>(null);
    let searchQuery = $state("");
    let activeFilter = $state<"all" | "pj" | "pnj">("all");
    let loading = $state(true);

    const gameId = $derived(page.params.id || "");

    async function loadData() {
        if (!gameId) return;
        try {
            loading = true;
            const [charData, playerData] = await Promise.all([
                fetchCharacters(gameId),
                fetchPlayers(gameId)
            ]);
            characters = charData;
            players = playerData;
        } catch (error) {
            console.error("Failed to load characters in sidebar:", error);
        } finally {
            loading = false;
        }
    }

    onMount(() => {
        loadData();
    });

    // Subscribe to character updates to reflect HP/conditions in real-time
    let charactersChannel: any = null;
    $effect(() => {
        if (gameId) {
            if (charactersChannel) supabase.removeChannel(charactersChannel);
            charactersChannel = supabase.channel(`sidebar_characters:${gameId}`)
                .on(
                    'postgres_changes',
                    { event: 'UPDATE', schema: 'public', table: 'characters' },
                    (payload) => {
                        const updatedChar = payload.new;
                        // Map updated character properties into the characters list
                        characters = characters.map(c => {
                            if (c.id === updatedChar.id) {
                                return {
                                    ...c,
                                    ...updatedChar,
                                    conditions: parseConditions(updatedChar.conditions)
                                };
                            }
                            return c;
                        });
                        // Sync selectedCharacter if it's the one that was updated
                        if (selectedCharacter && selectedCharacter.id === updatedChar.id) {
                            selectedCharacter = {
                                ...selectedCharacter,
                                ...updatedChar,
                                conditions: parseConditions(updatedChar.conditions)
                            };
                        }
                    }
                )
                .subscribe();
        }
        return () => {
            if (charactersChannel) supabase.removeChannel(charactersChannel);
        };
    });

    // Filtering logic
    let filteredCharacters = $derived(
        characters.filter((char) => {
            const matchesSearch = char.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
                (char.race && char.race.toLowerCase().includes(searchQuery.toLowerCase()));
            
            if (!matchesSearch) return false;

            if (activeFilter === "pj") {
                return !char.is_npc;
            } else if (activeFilter === "pnj") {
                return char.is_npc;
            }
            return true;
        })
    );

    function getPlayerName(char: Character): string | null {
        if (char.is_npc || !char.user_id) return null;
        const player = players.find(p => p.user_id === char.user_id);
        return player ? player.name : "Joueur inconnu";
    }

    function selectCharacter(char: Character) {
        selectedCharacter = char;
    }

    function deselectCharacter() {
        selectedCharacter = null;
    }

    async function placeToken(charId: string, event: Event) {
        event.stopPropagation();
        if (!$activeBoardStore) {
            alert("Aucun plateau n'est actif.");
            return;
        }
        try {
            await addBoardToken($activeBoardStore.id, gameId, charId, 0.5, 0.5);
        } catch (error: any) {
            console.error("Failed to place token:", error);
            alert(error.message || "Erreur lors du placement du personnage.");
        }
    }
</script>

<div class="h-full flex flex-col bg-stone-50">
    {#if selectedCharacter}
        <!-- Character Sheet View -->
        <div class="flex-1 flex flex-col min-h-0 bg-white">
            <!-- Sheet Header -->
            <div class="flex items-center justify-between px-4 py-3 border-b border-stone-100 bg-white">
                <button
                    onclick={deselectCharacter}
                    class="flex items-center gap-1.5 text-stone-500 hover:text-dark-gray transition-colors font-bold text-xs"
                >
                    <ArrowLeft size={16} />
                    <span>Retour</span>
                </button>

                <div class="flex items-center gap-2">
                    <span class="text-xs px-2 py-0.5 rounded-full font-semibold border
                        {selectedCharacter.is_npc 
                            ? 'bg-purple-50 text-purple-600 border-purple-100' 
                            : 'bg-emerald-50 text-emerald-600 border-emerald-100'}"
                    >
                        {selectedCharacter.is_npc ? 'PNJ' : 'PJ'}
                    </span>
                    <a
                        href="/table/{gameId}/gm/characters/{selectedCharacter.id}/edit"
                        class="p-1.5 bg-stone-100 hover:bg-stone-200 text-stone-600 hover:text-dark-gray rounded-lg transition-all flex items-center justify-center"
                        title="Modifier la fiche"
                    >
                        <Pencil size={14} />
                    </a>
                </div>
            </div>

            <!-- Sheet Content -->
            <div class="flex-1 overflow-y-auto">
                <CharacterSheet character={selectedCharacter} />
            </div>
        </div>
    {:else}
        <!-- Characters List View -->
        <!-- Search & Filters -->
        <div class="p-3 bg-white border-b border-stone-200 space-y-3">
            <div class="relative">
                <Search
                    class="absolute left-3 top-1/2 -translate-y-1/2 text-stone-400"
                    size={16}
                />
                <input
                    type="text"
                    bind:value={searchQuery}
                    placeholder="Rechercher un personnage..."
                    class="w-full pl-9 pr-4 py-2 bg-stone-50 border border-stone-200 rounded-lg text-sm focus:outline-none focus:border-burnt-orange transition-colors"
                />
            </div>

            <!-- Filters -->
            <div class="flex bg-stone-100 p-0.5 rounded-lg text-xs font-bold">
                <button
                    onclick={() => activeFilter = "all"}
                    class="flex-1 py-1.5 rounded-md text-center transition-all
                        {activeFilter === 'all' 
                            ? 'bg-white text-dark-gray shadow-sm' 
                            : 'text-stone-500 hover:text-dark-gray'}"
                >
                    Tous
                </button>
                <button
                    onclick={() => activeFilter = "pj"}
                    class="flex-1 py-1.5 rounded-md text-center transition-all
                        {activeFilter === 'pj' 
                            ? 'bg-white text-dark-gray shadow-sm' 
                            : 'text-stone-500 hover:text-dark-gray'}"
                >
                    Joueurs (PJ)
                </button>
                <button
                    onclick={() => activeFilter = "pnj"}
                    class="flex-1 py-1.5 rounded-md text-center transition-all
                        {activeFilter === 'pnj' 
                            ? 'bg-white text-dark-gray shadow-sm' 
                            : 'text-stone-500 hover:text-dark-gray'}"
                >
                    PNJ
                </button>
            </div>
        </div>

        <!-- List -->
        <div class="flex-1 overflow-y-auto p-2 space-y-2">
            {#if loading}
                <div class="flex justify-center py-12">
                    <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-burnt-orange"></div>
                </div>
            {:else if filteredCharacters.length === 0}
                <div class="text-center py-12 text-stone-400 flex flex-col items-center gap-2">
                    <Users size={32} class="text-stone-300" />
                    <p class="text-xs font-medium">Aucun personnage trouvé.</p>
                </div>
            {:else}
                {#each filteredCharacters as char}
                    <div
                        role="button"
                        tabindex="0"
                        onkeydown={(e) => { if (e.key === 'Enter' || e.key === ' ') selectCharacter(char); }}
                        onclick={() => selectCharacter(char)}
                        class="w-full text-left bg-white p-3 rounded-xl border border-stone-200 hover:border-stone-300 shadow-sm hover:shadow transition-all group flex flex-col gap-2 relative overflow-hidden cursor-pointer"
                    >
                        <!-- Top details -->
                        <div class="flex items-center gap-3">
                            {#if char.avatar_url}
                                <img
                                    src={char.avatar_url}
                                    alt={char.name}
                                    class="w-10 h-10 rounded-xl object-cover border border-stone-100 group-hover:scale-105 transition-transform duration-300"
                                />
                            {:else}
                                <div class="w-10 h-10 rounded-xl bg-gradient-to-br from-burnt-orange/10 to-burnt-orange/20 text-burnt-orange flex items-center justify-center font-bold text-sm">
                                    {(char.name || "?").charAt(0)}
                                </div>
                            {/if}

                            <div class="flex-1 min-w-0">
                                <div class="flex items-center justify-between">
                                    <h4 class="font-bold text-stone-800 truncate text-sm">
                                        {char.name}
                                    </h4>
                                    <div class="flex items-center gap-1.5">
                                        <button
                                            onclick={(e) => placeToken(char.id, e)}
                                            class="p-1 hover:bg-stone-100 text-stone-400 hover:text-burnt-orange rounded transition-colors"
                                            title="Placer sur le plateau"
                                        >
                                            <MapPin size={12} />
                                        </button>
                                        <span class="text-[10px] font-semibold px-1.5 py-0.5 rounded
                                            {char.is_npc 
                                                ? 'bg-purple-50 text-purple-600' 
                                                : 'bg-emerald-50 text-emerald-600'}"
                                        >
                                            {char.is_npc ? 'PNJ' : 'PJ'}
                                        </span>
                                    </div>
                                </div>
                                <div class="text-xs text-stone-400 flex items-center gap-1">
                                    <span>{char.race || "Inconnue"}</span>
                                    {#if char.sub_race}
                                        <span class="text-stone-300">•</span>
                                        <span class="truncate">{char.sub_race}</span>
                                    {/if}
                                </div>
                            </div>
                        </div>

                        <!-- Stats and HP -->
                        <div class="flex items-center justify-between pt-2 border-t border-stone-100">
                            <!-- HP -->
                            <div class="flex items-center gap-2 flex-1 max-w-[60%]">
                                <Heart size={12} class="text-red-500 fill-red-500/10" />
                                <div class="flex-1 bg-stone-100 h-1.5 rounded-full overflow-hidden border border-stone-200/50">
                                    <div 
                                        class="h-full rounded-full transition-all duration-300
                                            {char.current_hp < char.max_hp / 3 
                                                ? 'bg-red-500' 
                                                : 'bg-green-500'}"
                                        style="width: {(char.current_hp / char.max_hp) * 100}%"
                                    ></div>
                                </div>
                                <span class="text-[10px] font-mono font-bold text-stone-500">
                                    {char.current_hp}/{char.max_hp}
                                </span>
                            </div>

                            <!-- Badges CA & Speed -->
                            <div class="flex gap-2">
                                <div class="flex items-center gap-0.5 text-stone-500" title="Classe d'Armure (CA)">
                                    <Shield size={11} />
                                    <span class="text-[10px] font-bold font-mono">{char.armor_class ?? 10}</span>
                                </div>
                                {#if !char.is_npc && char.user_id}
                                    <div class="flex items-center gap-0.5 text-emerald-600" title="Assigné à un joueur">
                                        <UserCheck size={11} />
                                        <span class="text-[9px] font-bold truncate max-w-[50px]">
                                            {getPlayerName(char)}
                                        </span>
                                    </div>
                                {:else if !char.is_npc}
                                    <div class="flex items-center gap-0.5 text-stone-400" title="Non assigné">
                                        <UserX size={11} />
                                    </div>
                                {/if}
                            </div>
                        </div>
                    </div>
                {/each}
            {/if}
        </div>
    {/if}
</div>
