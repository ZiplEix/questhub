<script lang="ts">
    import PlayerLayout from "$lib/components/game/player/PlayerLayout.svelte";
    import ImmersionZone from "$lib/components/game/player/ImmersionZone.svelte";
    import PlayerDashboard from "$lib/components/game/player/PlayerDashboard.svelte";
    import { ChevronRight, ChevronLeft } from "lucide-svelte";

    let isDashboardOpen = $state(true);

    function toggleDashboard() {
        isDashboardOpen = !isDashboardOpen;
    }
</script>

<PlayerLayout>
    <div class="relative w-full h-full flex overflow-hidden">
        <!-- Left Zone: Immersion (Flexible) -->
        <div
            class="h-full relative transition-all duration-500 ease-in-out"
            style="width: {isDashboardOpen ? '65%' : '100%'}"
        >
            <ImmersionZone />

            <!-- Toggle Button (Attached to the right edge of Immersion Zone) -->
            <button
                onclick={toggleDashboard}
                class="absolute top-1/2 -translate-y-1/2 right-0 z-50 bg-white border border-stone-200 shadow-md p-1 rounded-l-lg hover:bg-stone-50 text-stone-500"
                aria-label={isDashboardOpen
                    ? "Fermer le tableau de bord"
                    : "Ouvrir le tableau de bord"}
            >
                {#if isDashboardOpen}
                    <ChevronRight size={20} />
                {:else}
                    <ChevronLeft size={20} />
                {/if}
            </button>
        </div>

        <!-- Right Zone: Dashboard (Collapsible) -->
        <div
            class="h-full relative transition-all duration-500 ease-in-out border-l border-stone-200 shadow-2xl z-40 bg-stone-50"
            style="width: {isDashboardOpen
                ? '35%'
                : '0%'}; opacity: {isDashboardOpen
                ? '1'
                : '0'}; pointer-events: {isDashboardOpen ? 'auto' : 'none'}"
        >
            <div class="w-full h-full min-w-[350px]">
                <!-- Prevent content squashing -->
                <PlayerDashboard />
            </div>
        </div>
    </div>
</PlayerLayout>
