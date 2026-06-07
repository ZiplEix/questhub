<script lang="ts">
    import Header from "$lib/components/Header.svelte";
    import Hero from "$lib/components/Hero.svelte";
    import FeatureCard from "$lib/components/FeatureCard.svelte";
    import PreviewSection from "$lib/components/PreviewSection.svelte";
    import MarketplacePreview from "$lib/components/MarketplacePreview.svelte";
    import { Sword, Backpack, Book, BookOpen } from "lucide-svelte";
</script>

<Header />

<main class="grow">
    <Hero />

    <section class="w-full max-w-7xl mx-auto px-4 md:px-8 py-12">
        <div class="grid sm:grid-cols-2 lg:grid-cols-4 gap-8">
            <FeatureCard
                icon={Sword}
                title="Combats Fluides"
                description="Passez en mode tour par tour en un clic. Initiatives automatiques et suivi des PV en temps réel."
            />
            <FeatureCard
                icon={Backpack}
                title="Inventaires Intelligents"
                description="Fini les gommes et les ratures. Glissez-déposez le butin directement dans la poche des joueurs."
            />
            <FeatureCard
                icon={Book}
                title="Bestiaire & Notes"
                description="Créez vos monstres, gardez une trace des PNJ rencontrés et révélez vos secrets au bon moment."
            />
            <FeatureCard
                icon={BookOpen}
                title="Histoire & Lore"
                description="Rédigez le lore en Markdown, organisez-le en dossiers et révélez-le en temps réel à vos joueurs."
            />
        </div>
    </section>

    <MarketplacePreview />

    <PreviewSection />
</main>
