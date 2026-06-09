<script lang="ts">
    import { ArrowRight, ZoomIn, X } from "lucide-svelte";
    import { authClient } from "$lib/auth-client";

    const session = authClient.useSession();

    let lightboxOpen = $state(false);
</script>

<section
    class="w-full max-w-7xl mx-auto px-4 md:px-8 py-12 md:py-20 grid md:grid-cols-2 gap-8 items-center"
>
    <div class="flex flex-col items-start gap-6">
        <h1
            class="font-display font-bold text-5xl md:text-6xl text-dark-gray leading-tight"
        >
            Jouez plus,<br />
            <span class="text-burnt-orange">calculez moins.</span>
        </h1>
        <p class="text-lg text-dark-gray/80 leading-relaxed max-w-lg">
            Le compagnon de jeu de rôle gratuit qui gère vos inventaires, vos
            dés et vos combats. Gardez l'esprit libre pour l'aventure.
        </p>

        {#if $session.data}
            <a
                href="/dashboard"
                class="group bg-burnt-orange text-white px-8 py-4 rounded-2xl font-bold text-lg shadow-lg hover:bg-opacity-90 transition-all hover:-translate-y-1 flex items-center gap-2 w-fit"
            >
                Vos tables vous attendent
                <ArrowRight
                    size={20}
                    class="group-hover:translate-x-1 transition-transform"
                />
            </a>
        {:else}
            <a
                href="/signup"
                class="group bg-burnt-orange text-white px-8 py-4 rounded-2xl font-bold text-lg shadow-lg hover:bg-opacity-90 transition-all hover:-translate-y-1 flex items-center gap-2 w-fit"
            >
                Rejoindre l'aventure
                <span class="text-sm font-normal opacity-90">(Gratuit)</span>
                <ArrowRight
                    size={20}
                    class="group-hover:translate-x-1 transition-transform"
                />
            </a>
        {/if}
    </div>
    <div class="relative flex items-center justify-center md:-mr-8 lg:-mr-16">
        <!-- Decorative blob background -->
        <div
            class="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-[110%] h-[110%] bg-mustard-yellow/25 rounded-full blur-3xl -z-10"
        ></div>

        <!-- Browser mockup frame — cliquable -->
        <button
            onclick={() => (lightboxOpen = true)}
            class="group w-full rounded-2xl overflow-hidden border border-stone-200 shadow-2xl cursor-zoom-in relative"
            style="transform: perspective(1200px) rotateY(-8deg) rotateX(3deg) scale(1.02); transform-origin: right center;"
            aria-label="Voir l'interface en grand"
        >
            <!-- Browser top bar -->
            <div class="bg-stone-100 border-b border-stone-200 px-4 py-2.5 flex items-center gap-3 shrink-0">
                <div class="flex items-center gap-1.5">
                    <div class="w-3 h-3 rounded-full bg-red-400"></div>
                    <div class="w-3 h-3 rounded-full bg-yellow-400"></div>
                    <div class="w-3 h-3 rounded-full bg-green-400"></div>
                </div>
                <div class="flex-1 bg-white rounded-md px-3 py-1 text-[11px] text-stone-400 font-mono border border-stone-200 max-w-xs mx-auto text-center">
                    https://questhub.fr/table/…
                </div>
            </div>
            <!-- Screenshot -->
            <img
                src="/example/player_view.png"
                alt="QuestHub — Vue Joueur"
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
                    https://questhub.fr/table/…
                </div>
            </div>
            <img
                src="/example/player_view.png"
                alt="QuestHub — Vue Joueur"
                class="w-full h-auto block"
            />
        </div>
    </div>
{/if}
