<script lang="ts">
    import Header from "$lib/components/Header.svelte";
    import Hero from "$lib/components/Hero.svelte";
    import FeatureCard from "$lib/components/FeatureCard.svelte";
    import PreviewSection from "$lib/components/PreviewSection.svelte";
    import Footer from "$lib/components/Footer.svelte";
    import { Sword, Backpack, Book } from "lucide-svelte";
</script>

<div class="min-h-screen flex flex-col">
    <Header />

    <main class="flex-grow">
        <Hero />

        <section class="w-full max-w-7xl mx-auto px-4 md:px-8 py-12">
            <div class="grid md:grid-cols-3 gap-8">
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
            </div>
        </section>

        <PreviewSection />
    </main>

    <Footer />
</div>
