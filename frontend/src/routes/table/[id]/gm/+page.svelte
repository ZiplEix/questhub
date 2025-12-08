<script lang="ts">
    import GMLayout from "$lib/components/game/mj/GMLayout.svelte";
    import GMTracker from "$lib/components/game/mj/GMTracker.svelte";
    import GMScene from "$lib/components/game/mj/GMScene.svelte";
    import GMOmniTool from "$lib/components/game/mj/GMOmniTool.svelte";
    import MonsterDrawer from "$lib/components/game/mj/MonsterDrawer.svelte";

    import { ChevronLeft, ChevronRight } from "lucide-svelte";
    import { onMount } from "svelte";
    import { page } from "$app/state";
    import { authClient } from "$lib/auth-client";
    import { fetchHistory } from "$lib/websocket";
    import { api } from "$lib/api";

    let selectedMonster = $state<any>(null);
    let isLeftPanelOpen = $state(true);
    let isRightPanelOpen = $state(true);
    let currentUserId = $state("");
    let gmCharacterId = $state("");

    // Lifted state for entities
    let entities = $state<any[]>([
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
        entities.sort((a: any, b: any) => b.init - a.init);
    }

    onMount(async () => {
        const gameId = page.params.id;
        try {
            const { data: tokenData } = await authClient.token();
            const token = tokenData?.token;
            if (typeof token === "string") {
                // Fetch chat history
                // @ts-ignore
                fetchHistory(gameId, token);

                // Fetch table data to get GM character ID
                const tableResponse = await api.get(`/table/${gameId}`, {
                    headers: {
                        Authorization: `Bearer ${token}`,
                    },
                });
                console.log("GM Page: Table Response", tableResponse.data);
                if (tableResponse.data.current_character_id) {
                    gmCharacterId = tableResponse.data.current_character_id;
                    console.log("GM Page: Set gmCharacterId", gmCharacterId);
                } else {
                    console.warn(
                        "GM Page: No current_character_id found in table response",
                    );
                }

                const { data: sessionData } = await authClient.getSession();
                if (sessionData?.user) {
                    currentUserId = sessionData.user.id;
                }

                // Fetch characters (players)
                const charsResponse = await api.get(
                    `/table/${gameId}/characters`,
                    {
                        headers: {
                            Authorization: `Bearer ${token}`,
                        },
                    },
                );
                const realPlayers = charsResponse.data
                    .filter((c: any) => c.user_id) // Only characters assigned to a user
                    .map((c: any) => ({
                        id: c.id,
                        name: c.name,
                        type: "player",
                        hp: c.current_hp,
                        maxHp: c.max_hp,
                        init: c.initiative,
                        status: [],
                        userId: c.user_id, // Important for chat targeting
                    }));

                // Add real players to entities, removing any existing players to avoid duplicates
                entities = [
                    ...entities.filter((e) => e.type !== "player"),
                    ...realPlayers,
                ];
            }
        } catch (e) {
            console.error(e);
        }
    });
</script>

<GMLayout>
    <!-- LEFT COLUMN: Flow Controller -->
    <aside
        class="shrink-0 z-20 shadow-xl relative transition-all duration-500 ease-in-out overflow-hidden"
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
        class="shrink-0 z-20 shadow-xl relative transition-all duration-500 ease-in-out overflow-hidden"
        style="width: {isRightPanelOpen
            ? '400px'
            : '0px'}; opacity: {isRightPanelOpen ? '1' : '0'};"
    >
        <div class="w-[400px] h-full">
            <GMOmniTool
                onSpawn={spawnMonster}
                {entities}
                {currentUserId}
                {gmCharacterId}
            />
        </div>
    </aside>
</GMLayout>
