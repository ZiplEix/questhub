<script lang="ts">
    import { activeBoardStore } from "$lib/websocket";
    import { Compass } from "lucide-svelte";
    import { fade } from "svelte/transition";

    let {
        isGM = false
    } = $props<{
        isGM?: boolean;
    }>();

    // Derived active map url from store
    let activeMapUrl = $derived($activeBoardStore?.image_url || "");

    let pings = $state<{ x: number; y: number; id: number }[]>([]);

    function handleSceneClick(e: MouseEvent) {
        // Pings are only enabled if a map is active
        if (!activeMapUrl) return;

        const rect = (e.currentTarget as HTMLElement).getBoundingClientRect();
        const x = e.clientX - rect.left;
        const y = e.clientY - rect.top;

        const id = Date.now();
        pings = [...pings, { x, y, id }];

        setTimeout(() => {
            pings = pings.filter((p) => p.id !== id);
        }, 2000);
    }
</script>

<div
    class="relative w-full h-full bg-stone-950 overflow-hidden select-none"
    onmousedown={handleSceneClick}
    role="button"
    tabindex="0"
    style="cursor: {activeMapUrl ? 'crosshair' : 'default'};"
>
    {#if activeMapUrl}
        <!-- Background Layer -->
        <div
            in:fade={{ duration: 400 }}
            class="absolute inset-0 bg-cover bg-center"
            style="background-image: url('{activeMapUrl}'); opacity: 0.85;"
        ></div>

        <!-- Grid Overlay (Optional, could be toggled) -->
        <div
            class="absolute inset-0 pointer-events-none opacity-10"
            style="background-image: radial-gradient(circle, #000 1px, transparent 1px); background-size: 40px 40px;"
        ></div>

        <!-- Pings Layer -->
        {#each pings as ping (ping.id)}
            <div
                class="absolute w-12 h-12 -ml-6 -mt-6 border-4 border-burnt-orange rounded-full animate-ping pointer-events-none"
                style="left: {ping.x}px; top: {ping.y}px;"
            ></div>
            <div
                class="absolute w-4 h-4 -ml-2 -mt-2 bg-burnt-orange rounded-full pointer-events-none"
                style="left: {ping.x}px; top: {ping.y}px;"
            ></div>
        {/each}

        <!-- Controls Overlay (Zoom, etc - placeholder) -->
        <div class="absolute bottom-4 right-4 flex gap-2">
            <div
                class="bg-black/50 backdrop-blur-md text-white px-3 py-1 rounded-full text-xs font-mono border border-white/10"
            >
                100%
            </div>
        </div>
    {:else}
        <!-- Premium Empty State -->
        <div class="absolute inset-0 flex flex-col items-center justify-center p-6 bg-stone-900">
            <div class="max-w-md w-full text-center space-y-6 bg-stone-900/60 backdrop-blur-md border border-stone-800/80 p-8 rounded-3xl shadow-2xl animate-in fade-in zoom-in duration-300">
                <div class="relative w-20 h-20 mx-auto flex items-center justify-center bg-burnt-orange/10 rounded-full border border-burnt-orange/20 text-burnt-orange animate-pulse">
                    <Compass size={40} />
                </div>
                <div class="space-y-2">
                    <h3 class="text-xl font-bold text-white tracking-wide">
                        {#if isGM}
                            Aucun plateau actif
                        {:else}
                            En attente du MJ
                        {/if}
                    </h3>
                    <p class="text-sm text-stone-400 leading-relaxed">
                        {#if isGM}
                            Créez un plateau et ajoutez-y une carte dans le panneau latéral droit (**onglet Plateaux**), puis activez-les pour lancer l'aventure.
                        {:else}
                            Le Maître du Jeu prépare le terrain. La carte s'affichera automatiquement ici dès qu'elle sera activée.
                        {/if}
                    </p>
                </div>
            </div>
        </div>
    {/if}
</div>


