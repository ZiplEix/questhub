<script lang="ts">
    import { ArrowRight } from "lucide-svelte";
</script>

<section
    class="w-full max-w-7xl mx-auto px-4 md:px-8 py-12 md:py-20 grid md:grid-cols-2 gap-12 items-center"
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
    </div>
    <div class="relative">
        <!-- Decorative blob background -->
        <div
            class="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-[120%] h-[120%] bg-mustard-yellow/20 rounded-full blur-3xl -z-10"
        ></div>
        <img
            src="https://placehold.co/600x400/F2CC8F/3D405B?text=Hero+Illustration"
            alt="QuestHub Interface on Laptop"
            class="w-full h-auto drop-shadow-xl transform hover:scale-[1.02] transition-transform duration-500"
        />
    </div>
</section>
