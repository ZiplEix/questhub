<script lang="ts">
    import { Image, Upload, Play } from "lucide-svelte";

    let scenes = [
        {
            id: 1,
            name: "Clairière",
            url: "https://placehold.co/600x400/2A9D8F/FFFFFF?text=Clairiere",
        },
        {
            id: 2,
            name: "Donjon - Entrée",
            url: "https://placehold.co/600x400/264653/FFFFFF?text=Donjon",
        },
        {
            id: 3,
            name: "Taverne",
            url: "https://placehold.co/600x400/E9C46A/000000?text=Taverne",
        },
    ];

    let activeSceneId = $state(1);
</script>

<div class="h-full flex flex-col bg-stone-50">
    <!-- Upload -->
    <div class="p-3 bg-white border-b border-stone-200">
        <button
            class="w-full py-2 flex items-center justify-center gap-2 bg-stone-100 hover:bg-stone-200 text-stone-600 font-bold text-xs rounded-lg transition-colors border border-stone-200 border-dashed"
        >
            <Upload size={14} />
            UPLOADER UNE IMAGE
        </button>
    </div>

    <!-- Scene List -->
    <div
        class="flex-1 overflow-y-auto p-2 grid grid-cols-2 gap-2 content-start"
    >
        {#each scenes as scene}
            <button
                class="group relative aspect-video rounded-xl overflow-hidden border-2 transition-all
                {activeSceneId === scene.id
                    ? 'border-burnt-orange ring-2 ring-burnt-orange/20'
                    : 'border-stone-200 hover:border-stone-300'}"
                onclick={() => (activeSceneId = scene.id)}
            >
                <img
                    src={scene.url}
                    alt={scene.name}
                    class="w-full h-full object-cover"
                />

                <div
                    class="absolute inset-x-0 bottom-0 bg-gradient-to-t from-black/80 to-transparent p-2 pt-6"
                >
                    <div class="text-white text-xs font-bold truncate">
                        {scene.name}
                    </div>
                </div>

                {#if activeSceneId !== scene.id}
                    <div
                        class="absolute inset-0 bg-black/40 opacity-0 group-hover:opacity-100 flex items-center justify-center transition-opacity"
                    >
                        <div
                            class="bg-white text-stone-900 px-2 py-1 rounded text-xs font-bold flex items-center gap-1"
                        >
                            <Play size={10} /> ACTIVER
                        </div>
                    </div>
                {/if}
            </button>
        {/each}
    </div>
</div>
