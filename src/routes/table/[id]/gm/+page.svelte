<script lang="ts">
    import GMLayout from "$lib/components/game/mj/GMLayout.svelte";
    import GMTracker from "$lib/components/game/mj/GMTracker.svelte";
    import GMScene from "$lib/components/game/mj/GMScene.svelte";
    import GMOmniTool from "$lib/components/game/mj/GMOmniTool.svelte";

    import { ChevronLeft, ChevronRight } from "lucide-svelte";
    import { onMount } from "svelte";
    import { page } from "$app/state";
    import { authClient } from "$lib/auth-client";
    import { fetchHistory } from "$lib/websocket";
    import { fetchGame } from "$lib/api";

    let isLeftPanelOpen = $state(true);
    let isRightPanelOpen = $state(true);
    let currentUserId = $state("");
    let gmCharacterId = $state("");
    let rightPanelWidth = $state(400);
    let isDraggingRight = $state(false);

    function handleMouseDownRight(e: MouseEvent) {
        e.preventDefault();
        isDraggingRight = true;
    }

    function handleMouseMove(e: MouseEvent) {
        if (isDraggingRight) {
            const newWidth = window.innerWidth - e.clientX;
            rightPanelWidth = Math.min(Math.max(newWidth, 320), 800);
        }
    }

    function handleMouseUp() {
        isDraggingRight = false;
    }

    onMount(async () => {
        const gameId = page.params.id || "";
        try {
            const { data: tokenData } = await authClient.token();
            const token = tokenData?.token;
            if (typeof token === "string") {
                // Fetch chat history
                // @ts-ignore
                fetchHistory(gameId, token);

                // Fetch table data to get GM character ID
                const gameData = await fetchGame(gameId);
                console.log("GM Page: Table Response", gameData);
                if (gameData.current_character_id) {
                    gmCharacterId = gameData.current_character_id;
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
            }
        } catch (e) {
            console.error(e);
        }
    });
</script>

<svelte:window 
    onmousemove={handleMouseMove} 
    onmouseup={handleMouseUp} 
/>

<GMLayout>
    <!-- LEFT COLUMN: Flow Controller -->
    <aside
        class="shrink-0 z-20 shadow-xl relative transition-all duration-500 ease-in-out overflow-hidden"
        style="width: {isLeftPanelOpen
            ? '300px'
            : '0px'}; opacity: {isLeftPanelOpen ? '1' : '0'};"
    >
        <div class="w-[300px] h-full">
            <GMTracker />
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
        class="shrink-0 z-20 shadow-xl relative {isDraggingRight ? '' : 'transition-all duration-500 ease-in-out'} overflow-hidden"
        style="width: {isRightPanelOpen
            ? rightPanelWidth + 'px'
            : '0px'}; opacity: {isRightPanelOpen ? '1' : '0'};"
    >
        <!-- Resize Handle -->
        <div
            onmousedown={handleMouseDownRight}
            class="absolute top-0 left-0 w-1.5 h-full cursor-col-resize hover:bg-burnt-orange/50 active:bg-burnt-orange transition-all z-30 flex items-center justify-center group"
            role="separator"
            aria-label="Ajuster la largeur"
        >
            <div class="w-[2px] h-12 bg-stone-300/20 group-hover:bg-white/40 rounded transition-colors"></div>
        </div>

        <div class="h-full" style="width: {rightPanelWidth}px;">
            <GMOmniTool
                {currentUserId}
                {gmCharacterId}
            />
        </div>
    </aside>
</GMLayout>
