<script lang="ts">
    import { Search, Plus, Skull } from "lucide-svelte";
    import { fetchMonsters } from "$lib/api";
    import { page } from "$app/state";
    import { onMount } from "svelte";



    let searchQuery = $state("");
    let monsters = $state<any[]>([]);
    let loading = $state(true);

    onMount(async () => {
        try {
            const gameId = page.params.id;
            if (gameId) {
                monsters = await fetchMonsters(gameId);
            }
        } catch (e) {
            console.error("Failed to fetch bestiary monsters:", e);
        } finally {
            loading = false;
        }
    });

    let filteredMonsters = $derived(
        monsters.filter((m) =>
            m.name.toLowerCase().includes(searchQuery.toLowerCase()),
        ),
    );
</script>

<div class="h-full flex flex-col bg-stone-50">
    <!-- Search -->
    <div class="p-3 bg-white border-b border-stone-200 space-y-2">
        <div class="relative">
            <Search
                class="absolute left-3 top-1/2 -translate-y-1/2 text-stone-400"
                size={16}
            />
            <input
                type="text"
                bind:value={searchQuery}
                placeholder="Rechercher un monstre..."
                class="w-full pl-9 pr-4 py-2 bg-stone-50 border border-stone-200 rounded-lg text-sm focus:outline-none focus:border-burnt-orange transition-colors"
            />
        </div>
    </div>

    <!-- Monster List -->
    <div class="flex-1 overflow-y-auto p-2 space-y-2">
        {#if loading}
            <div class="flex justify-center py-8">
                <div
                    class="animate-spin rounded-full h-6 w-6 border-b-2 border-burnt-orange"
                ></div>
            </div>
        {:else if filteredMonsters.length === 0}
            <div
                class="text-center py-8 text-stone-400 flex flex-col items-center gap-2"
            >
                <Skull size={24} />
                <p class="text-xs font-medium">
                    {searchQuery
                        ? "Aucun monstre trouvé."
                        : "Aucun monstre dans le bestiaire."}
                </p>
                <p class="text-[11px] text-stone-300">
                    Ajoutez des monstres depuis les paramètres de la partie.
                </p>
            </div>
        {:else}
            {#each filteredMonsters as monster}
                <div
                    class="bg-white p-3 rounded-xl border border-stone-200 shadow-sm hover:shadow-md transition-all group"
                >
                    <div class="flex justify-between items-start mb-2">
                        <div class="flex items-center gap-3">
                            {#if monster.avatar_url}
                                <img
                                    src={monster.avatar_url}
                                    alt={monster.name}
                                    class="w-8 h-8 rounded-lg object-cover"
                                />
                            {/if}
                            <div>
                                <div class="font-bold text-stone-800">
                                    {monster.name}
                                </div>
                                <div class="text-xs text-stone-500">
                                    {monster.race || "Monstre"}
                                </div>
                            </div>
                        </div>

                    </div>

                    <div
                        class="flex gap-3 text-xs font-mono text-stone-500 bg-stone-50 p-1.5 rounded-lg"
                    >
                        <span
                            ><strong class="text-stone-700">AC</strong>
                            {monster.armor_class ?? 10}</span
                        >
                        <span
                            ><strong class="text-stone-700">HP</strong>
                            {monster.max_hp}</span
                        >
                    </div>
                </div>
            {/each}
        {/if}
    </div>
</div>
