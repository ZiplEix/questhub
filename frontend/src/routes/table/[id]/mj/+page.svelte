<script lang="ts">
    import GMLayout from "$lib/components/game/mj/GMLayout.svelte";
    import GMTracker from "$lib/components/game/mj/GMTracker.svelte";
    import GMScene from "$lib/components/game/mj/GMScene.svelte";
    import GMOmniTool from "$lib/components/game/mj/GMOmniTool.svelte";
    import MonsterDrawer from "$lib/components/game/mj/MonsterDrawer.svelte";

    import { ChevronLeft, ChevronRight } from "lucide-svelte";

    let selectedMonster = $state<any>(null);
    let isLeftPanelOpen = $state(true);
    let isRightPanelOpen = $state(true);

    // Lifted state for entities
    let entities = $state([
        {
            id: 1,
            name: "Eldric",
            type: "player",
            hp: 24,
            maxHp: 24,
            init: 18,
            status: [],
        },
        {
            id: 2,
            name: "Gobelin Chef",
            type: "monster",
            hp: 15,
            maxHp: 20,
            init: 12,
            status: ["invisible"],
        },
        {
            id: 3,
            name: "Lyra",
            type: "player",
            hp: 18,
            maxHp: 18,
            init: 10,
            status: [],
        },
        {
            id: 4,
            name: "Gobelin 2",
            type: "monster",
            hp: 0,
            maxHp: 7,
            init: 5,
            status: ["dead"],
        },
    ]);

    function handleSelectMonster(monster: any) {
        selectedMonster = monster;
    }

    function spawnMonster(monsterTemplate: any) {
        const newMonster = {
            id: Date.now(),
            name: monsterTemplate.name,
            type: "monster",
            hp: monsterTemplate.hp,
            maxHp: monsterTemplate.hp,
            init: Math.floor(Math.random() * 20) + 1,
            status: [],
            // Extra stats for drawer
            ac: monsterTemplate.ac,
            cr: monsterTemplate.cr,
            monsterType: monsterTemplate.type,
        };
        entities.push(newMonster);
        // Sort by initiative (simple desc sort)
        entities.sort((a, b) => b.init - a.init);
    }
</script>

<GMLayout>
    <!-- LEFT COLUMN: Flow Controller -->
    <aside
        class="flex-shrink-0 z-20 shadow-xl relative transition-all duration-500 ease-in-out overflow-hidden"
        style="width: {isLeftPanelOpen
            ? '300px'
            : '0px'}; opacity: {isLeftPanelOpen ? '1' : '0'};"
    >
        <div class="w-[300px] h-full">
            <GMTracker {entities} onSelect={handleSelectMonster} />
        </div>
    </aside>

    <!-- CENTER ZONE: Game Table (Flexible) -->
    <main class="flex-1 relative z-0 bg-stone-900 flex flex-col">
        <!-- Left Toggle Button -->
        <button
            onclick={() => (isLeftPanelOpen = !isLeftPanelOpen)}
            class="absolute top-1/2 -translate-y-1/2 left-0 z-50 bg-white border border-stone-200 shadow-md p-1 rounded-r-lg hover:bg-stone-50 text-stone-500"
            aria-label={isLeftPanelOpen
                ? "Fermer le panneau gauche"
                : "Ouvrir le panneau gauche"}
        >
            {#if isLeftPanelOpen}
                <ChevronLeft size={16} />
            {:else}
                <ChevronRight size={16} />
            {/if}
        </button>

        <GMScene />

        <!-- Monster Detail Drawer (Overlay) -->
        {#if selectedMonster}
            <MonsterDrawer
                monster={selectedMonster}
                onClose={() => (selectedMonster = null)}
            />
        {/if}

        <!-- Right Toggle Button -->
        <button
            onclick={() => (isRightPanelOpen = !isRightPanelOpen)}
            class="absolute top-1/2 -translate-y-1/2 right-0 z-50 bg-white border border-stone-200 shadow-md p-1 rounded-l-lg hover:bg-stone-50 text-stone-500"
            aria-label={isRightPanelOpen
                ? "Fermer le panneau droit"
                : "Ouvrir le panneau droit"}
        >
            {#if isRightPanelOpen}
                <ChevronRight size={16} />
            {:else}
                <ChevronLeft size={16} />
            {/if}
        </button>
    </main>

    <!-- RIGHT COLUMN: Omni-Tool -->
    <aside
        class="flex-shrink-0 z-20 shadow-xl relative transition-all duration-500 ease-in-out overflow-hidden"
        style="width: {isRightPanelOpen
            ? '400px'
            : '0px'}; opacity: {isRightPanelOpen ? '1' : '0'};"
    >
        <div class="w-[400px] h-full">
            <GMOmniTool onSpawn={spawnMonster} />
        </div>
    </aside>
</GMLayout>
