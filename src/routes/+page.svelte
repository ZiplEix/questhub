<script lang="ts">
    import Header from "$lib/components/Header.svelte";
    import Hero from "$lib/components/Hero.svelte";
    import FeatureCard from "$lib/components/FeatureCard.svelte";
    import PreviewSection from "$lib/components/PreviewSection.svelte";
    import MarketplacePreview from "$lib/components/MarketplacePreview.svelte";
    import { Sword, Backpack, Book, BookOpen, TriangleAlert, X } from "lucide-svelte";

    let betaBannerOpen = $state(true);
</script>

<svelte:head>
    <title>QuestHub — Table de jeu virtuelle (VTT) & Compagnon de JDR</title>
    <meta name="description" content="QuestHub est une table de jeu virtuelle (VTT) moderne et simplifiée. Gérez vos combats en temps réel, vos feuilles de personnages, votre lore et partagez vos créations." />
    <meta property="og:type" content="website" />
    <meta property="og:title" content="QuestHub — Table de jeu virtuelle (VTT) & Compagnon de JDR" />
    <meta property="og:description" content="Gérez vos campagnes de jeu de rôle sur table en toute simplicité. Combats fluides en temps réel, inventaires partagés, bestiaire et lore interactif." />
    <meta property="og:image" content="https://questhub.fr/example/gm_dashboard.png" />
    <meta property="twitter:card" content="summary_large_image" />
    <meta property="twitter:title" content="QuestHub — Table de jeu virtuelle (VTT) & Compagnon de JDR" />
    <meta property="twitter:description" content="Gérez vos campagnes de jeu de rôle sur table en toute simplicité. Combats fluides en temps réel, inventaires partagés, bestiaire et lore interactif." />
    <meta property="twitter:image" content="https://questhub.fr/example/gm_dashboard.png" />
</svelte:head>


{#if betaBannerOpen}
    <div class="w-full bg-amber-50 border-b border-amber-200 px-4 py-2.5 flex items-center justify-center gap-3 text-sm text-amber-800">
        <TriangleAlert size={16} class="text-amber-500 shrink-0" />
        <span>
            <strong>Version bêta</strong> — QuestHub est en cours de développement. Nous ne pouvons pas garantir l'intégrité des données stockées sur le site.
        </span>
        <button
            onclick={() => (betaBannerOpen = false)}
            class="ml-2 p-1 text-amber-500 hover:text-amber-700 hover:bg-amber-100 rounded-full transition-all cursor-pointer shrink-0"
            aria-label="Fermer"
        >
            <X size={14} />
        </button>
    </div>
{/if}

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
