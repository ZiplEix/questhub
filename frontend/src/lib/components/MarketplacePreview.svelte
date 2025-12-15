<script lang="ts">
    import { Store } from "lucide-svelte";
</script>

<section class="max-w-7xl mx-auto px-4 md:px-8 mb-12">
    <div
        class="bg-stone-50 rounded-3xl p-12 md:p-16 text-center border border-stone-100 relative overflow-hidden group"
    >
        <!-- Decorative background elements -->
        <div
            class="absolute top-0 right-0 w-64 h-64 bg-burnt-orange/5 rounded-full blur-3xl transform translate-x-1/2 -translate-y-1/2 group-hover:bg-burnt-orange/10 transition-all duration-700"
        ></div>

        <div class="relative z-10 max-w-2xl mx-auto">
            <h2
                class="font-display font-bold text-3xl md:text-4xl text-dark-gray mb-6"
            >
                Partagez avec la Communauté
            </h2>
            <p class="text-lg text-dark-gray/60 mb-8 leading-relaxed">
                Accédez à une bibliothèque grandissante de personnages, monstres
                et objets créés par d'autres Maîtres du Jeu. Importez-les
                directement dans vos parties en un clic.
            </p>

            <a
                href="/community-market"
                class="inline-flex items-center gap-2 bg-dark-gray text-white px-8 py-4 rounded-2xl font-bold shadow-lg hover:bg-burnt-orange transition-all hover:-translate-y-1"
            >
                <Store size={20} />
                Explorer le Marché
            </a>
        </div>
    </div>
</section>
