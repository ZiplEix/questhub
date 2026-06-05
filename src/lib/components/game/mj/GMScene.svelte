<script lang="ts">
    import Scene from "../Scene.svelte";
    import ContextMenu from "./ContextMenu.svelte";
    import { MousePointer2, CloudFog, Ruler } from "lucide-svelte";

    let activeTool = $state<"pointer" | "fog" | "measure">("pointer");

    // Context Menu State
    let contextMenu = $state<{ x: number; y: number; visible: boolean }>({
        x: 0,
        y: 0,
        visible: false,
    });

    function handleContextMenu(e: MouseEvent) {
        e.preventDefault();
        contextMenu = {
            x: e.clientX,
            y: e.clientY,
            visible: true,
        };
    }

    function closeContextMenu() {
        contextMenu.visible = false;
    }
</script>

<div
    class="relative w-full h-full bg-stone-900"
    oncontextmenu={handleContextMenu}
    role="application"
>
    <!-- Toolbar -->
    <div
        class="absolute top-4 left-1/2 -translate-x-1/2 z-20 flex gap-1 bg-black/70 backdrop-blur-md p-1.5 rounded-full border border-white/10 shadow-xl"
    >
        <button
            class="p-2 rounded-full transition-colors {activeTool === 'pointer'
                ? 'bg-white text-black'
                : 'text-stone-400 hover:text-white hover:bg-white/10'}"
            onclick={() => (activeTool = "pointer")}
            title="Pointeur (Ping)"
        >
            <MousePointer2 size={18} />
        </button>
        <button
            class="p-2 rounded-full transition-colors {activeTool === 'fog'
                ? 'bg-white text-black'
                : 'text-stone-400 hover:text-white hover:bg-white/10'}"
            onclick={() => (activeTool = "fog")}
            title="Brouillard de guerre"
        >
            <CloudFog size={18} />
        </button>
        <button
            class="p-2 rounded-full transition-colors {activeTool === 'measure'
                ? 'bg-white text-black'
                : 'text-stone-400 hover:text-white hover:bg-white/10'}"
            onclick={() => (activeTool = "measure")}
            title="Mesure"
        >
            <Ruler size={18} />
        </button>
    </div>

    <!-- Scene Wrapper -->
    <Scene />

    <!-- Context Menu -->
    {#if contextMenu.visible}
        <ContextMenu
            x={contextMenu.x}
            y={contextMenu.y}
            onClose={closeContextMenu}
        />
    {/if}
</div>
