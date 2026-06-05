<script lang="ts">
    import { Search, Plus, Skull } from "lucide-svelte";

    let { onSpawn } = $props();

    let searchQuery = $state("");

    let monsters = [
        { id: 1, name: "Gobelin", cr: "1/4", type: "Humanoid", hp: 7, ac: 15 },
        {
            id: 2,
            name: "Gobelin Chef",
            cr: "1",
            type: "Humanoid",
            hp: 22,
            ac: 17,
        },
        { id: 3, name: "Loup", cr: "1/4", type: "Beast", hp: 11, ac: 13 },
        { id: 4, name: "Ogre", cr: "2", type: "Giant", hp: 59, ac: 11 },
        { id: 5, name: "Squelette", cr: "1/4", type: "Undead", hp: 13, ac: 13 },
    ];

    let filteredMonsters = $derived(
        monsters.filter((m) =>
            m.name.toLowerCase().includes(searchQuery.toLowerCase()),
        ),
    );
</script>

<div class="h-full flex flex-col bg-stone-50">
    <!-- Search & Quick Add -->
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
        <button
            class="w-full py-2 flex items-center justify-center gap-2 bg-stone-100 hover:bg-stone-200 text-stone-600 font-bold text-xs rounded-lg transition-colors border border-stone-200 border-dashed"
        >
            <Plus size={14} />
            CRÉER UN MONSTRE RAPIDE
        </button>
    </div>

    <!-- Monster List -->
    <div class="flex-1 overflow-y-auto p-2 space-y-2">
        {#each filteredMonsters as monster}
            <div
                class="bg-white p-3 rounded-xl border border-stone-200 shadow-sm hover:shadow-md transition-all group"
            >
                <div class="flex justify-between items-start mb-2">
                    <div>
                        <div class="font-bold text-stone-800">
                            {monster.name}
                        </div>
                        <div class="text-xs text-stone-500">
                            {monster.type} • CR {monster.cr}
                        </div>
                    </div>
                    <div
                        class="flex gap-1 opacity-0 group-hover:opacity-100 transition-opacity"
                    >
                        <button
                            class="p-1.5 bg-burnt-orange text-white rounded-lg hover:bg-orange-600"
                            title="Ajouter au combat"
                            onclick={() => onSpawn?.(monster)}
                        >
                            <Plus size={14} />
                        </button>
                    </div>
                </div>

                <div
                    class="flex gap-3 text-xs font-mono text-stone-500 bg-stone-50 p-1.5 rounded-lg"
                >
                    <span
                        ><strong class="text-stone-700">AC</strong>
                        {monster.ac}</span
                    >
                    <span
                        ><strong class="text-stone-700">HP</strong>
                        {monster.hp}</span
                    >
                </div>
            </div>
        {/each}
    </div>
</div>
