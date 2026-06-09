<script lang="ts">
    import { ZoomIn, X } from "lucide-svelte";
    let lightboxOpen = $state(false);
</script>

<section class="w-full max-w-7xl mx-auto px-4 md:px-8 py-20">
    <div class="bg-dark-gray/5 p-4 md:p-8 rounded-3xl">
        <div class="text-center mb-12">
            <h2
                class="font-display font-bold text-3xl md:text-4xl text-dark-gray mb-4"
            >
                Une interface pensée pour les MJs
            </h2>
            <p class="text-dark-gray/70 max-w-2xl mx-auto">
                Tout ce dont vous avez besoin à portée de clic, sans encombrer
                votre écran.
            </p>
        </div>

        <!-- Browser mockup frame — cliquable -->
        <button
            onclick={() => (lightboxOpen = true)}
            class="group w-full rounded-2xl overflow-hidden border border-stone-200 shadow-2xl cursor-zoom-in relative block"
            aria-label="Voir l'interface en grand"
        >
            <!-- Browser top bar -->
            <div class="bg-stone-100 border-b border-stone-200 px-4 py-2.5 flex items-center gap-3">
                <div class="flex items-center gap-1.5">
                    <div class="w-3 h-3 rounded-full bg-red-400"></div>
                    <div class="w-3 h-3 rounded-full bg-yellow-400"></div>
                    <div class="w-3 h-3 rounded-full bg-green-400"></div>
                </div>
                <div class="flex-1 bg-white rounded-md px-3 py-1 text-[11px] text-stone-400 font-mono border border-stone-200 max-w-xs mx-auto text-center">
                    https://questhub.fr/table/…/gm
                </div>
            </div>
            <!-- Screenshot -->
            <img
                src="/example/gm_dashboard.png"
                alt="Interface du Maître du Jeu"
                class="w-full h-auto block"
            />
            <!-- Hover overlay -->
            <div class="absolute inset-0 bg-black/0 group-hover:bg-black/20 transition-all duration-300 flex items-center justify-center">
                <div class="opacity-0 group-hover:opacity-100 transition-opacity bg-white/90 rounded-full p-3 shadow-lg">
                    <ZoomIn size={24} class="text-dark-gray" />
                </div>
            </div>
        </button>
    </div>
</section>

<!-- Lightbox -->
{#if lightboxOpen}
    <div
        class="fixed inset-0 z-50 bg-black/80 backdrop-blur-sm flex items-center justify-center p-4 animate-in fade-in duration-200"
        onclick={() => (lightboxOpen = false)}
        role="dialog"
        aria-modal="true"
    >
        <button
            onclick={() => (lightboxOpen = false)}
            class="absolute top-4 right-4 p-2 text-white/70 hover:text-white hover:bg-white/10 rounded-full transition-all cursor-pointer"
            aria-label="Fermer"
        >
            <X size={28} />
        </button>

        <div
            class="w-full max-w-6xl rounded-2xl overflow-hidden shadow-2xl border border-white/10 animate-in zoom-in-95 duration-200"
            onclick={(e) => e.stopPropagation()}
        >
            <div class="bg-stone-900 border-b border-white/10 px-4 py-2.5 flex items-center gap-3">
                <div class="flex items-center gap-1.5">
                    <div class="w-3 h-3 rounded-full bg-red-400"></div>
                    <div class="w-3 h-3 rounded-full bg-yellow-400"></div>
                    <div class="w-3 h-3 rounded-full bg-green-400"></div>
                </div>
                <div class="flex-1 bg-white/10 rounded-md px-3 py-1 text-[11px] text-white/50 font-mono max-w-xs mx-auto text-center">
                    https://questhub.fr/table/…/gm
                </div>
            </div>
            <img
                src="/example/gm_dashboard.png"
                alt="Interface du Maître du Jeu"
                class="w-full h-auto block"
            />
        </div>
    </div>
{/if}
