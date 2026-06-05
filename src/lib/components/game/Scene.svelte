<script lang="ts">
    let {
        backgroundUrl = "https://placehold.co/600x400/F2CC8F/3D405B?text=Map",
    } = $props();

    let pings = $state<{ x: number; y: number; id: number }[]>([]);

    function handleSceneClick(e: MouseEvent) {
        // Long press logic could go here for pings
        // For now, just a simple click creates a ping for demo
        const rect = (e.target as HTMLElement).getBoundingClientRect();
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
    class="relative w-full h-full bg-stone-900 overflow-hidden cursor-crosshair select-none"
    onmousedown={handleSceneClick}
    role="button"
    tabindex="0"
>
    <!-- Background Layer -->
    <div
        class="absolute inset-0 bg-cover bg-center transition-opacity duration-500"
        style="background-image: url('{backgroundUrl}'); opacity: 0.8;"
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
            class="bg-black/50 backdrop-blur-md text-white px-3 py-1 rounded-full text-xs font-mono"
        >
            100%
        </div>
    </div>
</div>


