<script lang="ts">
    import GMLayout from "$lib/components/game/mj/GMLayout.svelte";
    import GMTracker from "$lib/components/game/mj/GMTracker.svelte";
    import GMScene from "$lib/components/game/mj/GMScene.svelte";
    import GMOmniTool from "$lib/components/game/mj/GMOmniTool.svelte";

    import { ChevronLeft, ChevronRight } from "lucide-svelte";
    import { onMount, setContext } from "svelte";
    import { page } from "$app/state";
    import { authClient } from "$lib/auth-client";
    import { fetchHistory } from "$lib/websocket";
    import { fetchGame } from "$lib/api";
    import { supabase } from "$lib/supabaseClient";

    let game = $state<any>(null);
    let isLeftPanelOpen = $state(true);
    let isRightPanelOpen = $state(true);
    let currentUserId = $state("");
    let gmCharacterId = $state("");

    setContext("tableContext", {
        get isReadOnly() {
            return game ? !game.is_active : false;
        }
    });

    let gameChannel: any = null;
    $effect(() => {
        const gameId = page.params.id;
        if (gameId) {
            if (gameChannel) supabase.removeChannel(gameChannel);
            gameChannel = supabase.channel(`game_realtime_gm:${gameId}`)
                .on(
                    'postgres_changes',
                    { event: 'UPDATE', schema: 'public', table: 'games', filter: `id=eq.${gameId}` },
                    (payload) => {
                        console.log("Realtime game update received in GM view:", payload.new);
                        game = {
                            ...game,
                            ...payload.new
                        };
                    }
                )
                .subscribe();
        }
        return () => {
            if (gameChannel) {
                supabase.removeChannel(gameChannel);
            }
        };
    });

    let leftPanelWidth = $state(380);
    let isDraggingLeft = $state(false);
    let rightPanelWidth = $state(400);
    let isDraggingRight = $state(false);

    function handleMouseDownLeft(e: MouseEvent) {
        e.preventDefault();
        isDraggingLeft = true;
    }

    function handleMouseDownRight(e: MouseEvent) {
        e.preventDefault();
        isDraggingRight = true;
    }

    function handleMouseMove(e: MouseEvent) {
        if (isDraggingLeft) {
            const newWidth = e.clientX;
            leftPanelWidth = Math.min(Math.max(newWidth, 320), 600);
        }
        if (isDraggingRight) {
            const newWidth = window.innerWidth - e.clientX;
            rightPanelWidth = Math.min(Math.max(newWidth, 320), 800);
        }
    }

    function handleMouseUp() {
        isDraggingLeft = false;
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
                game = gameData;
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
        class="shrink-0 z-20 shadow-xl relative {isDraggingLeft ? '' : 'transition-all duration-500 ease-in-out'} overflow-hidden"
        style="width: {isLeftPanelOpen
            ? leftPanelWidth + 'px'
            : '0px'}; opacity: {isLeftPanelOpen ? '1' : '0'};"
    >
        <div class="h-full" style="width: {leftPanelWidth}px;">
            <GMTracker />
        </div>
        <!-- Resize Handle -->
        <div
            onmousedown={handleMouseDownLeft}
            class="absolute top-0 right-0 w-1.5 h-full cursor-col-resize hover:bg-burnt-orange/50 active:bg-burnt-orange transition-all z-30 flex items-center justify-center group"
            role="separator"
            aria-label="Ajuster la largeur"
        >
            <div class="w-[2px] h-12 bg-stone-300/20 group-hover:bg-white/40 rounded transition-colors"></div>
        </div>
    </aside>

    <!-- CENTER ZONE: Game Table (Flexible) -->
    <main class="flex-1 relative z-0 bg-stone-900 flex flex-col">
        {#if game && !game.is_active}
            <div
                class="absolute top-6 left-1/2 -translate-x-1/2 z-[60] bg-stone-900/95 backdrop-blur-md border border-stone-850 text-stone-200 px-5 py-2.5 rounded-full shadow-2xl flex items-center gap-2.5 animate-in slide-in-from-top-4 fade-in duration-300 pointer-events-none select-none"
            >
                <span class="relative flex h-2 w-2">
                    <span class="animate-ping absolute inline-flex h-full w-full rounded-full bg-red-400 opacity-75"></span>
                    <span class="relative inline-flex rounded-full h-2 w-2 bg-red-500"></span>
                </span>
                <span class="font-bold text-sm tracking-wide uppercase"
                    >Partie Archivée (Lecture Seule)</span
                >
            </div>
        {/if}
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
