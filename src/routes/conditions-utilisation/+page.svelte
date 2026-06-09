<script lang="ts">
    import Header from "$lib/components/Header.svelte";
    import { ScrollText, UserCheck, Shield, AlertTriangle, Gavel } from "lucide-svelte";
</script>

<svelte:head>
    <title>Conditions d'utilisation | QuestHub</title>
    <meta
        name="description"
        content="Conditions d'utilisation de QuestHub : règles d'accès, responsabilités, contenus publiés et conditions de suspension de compte."
    />
</svelte:head>

<Header />

<main class="w-full max-w-7xl mx-auto px-4 md:px-8 py-16 grow">
    <div class="mb-10 text-center">
        <h1 class="text-3xl font-display font-black text-burnt-orange tracking-tight flex items-center justify-center gap-3">
            <ScrollText class="text-burnt-orange" size={28} />
            Conditions d'utilisation
        </h1>
        <p class="text-sm text-stone-500 mt-2 font-medium">
            En vigueur au 9 juin 2026. Ces conditions encadrent l'utilisation de la plateforme QuestHub.
        </p>
    </div>

    <div class="bg-white rounded-2xl border border-stone-200/80 shadow-xs p-6 md:p-10 space-y-10 leading-relaxed text-stone-700 max-w-3xl mx-auto">
        <section class="space-y-3">
            <h2 class="text-lg font-bold text-dark-gray flex items-center gap-2 border-b border-stone-100 pb-2 font-display">
                <UserCheck size={18} class="text-burnt-orange" />
                1. Objet et acceptation
            </h2>
            <p class="text-sm">
                Les présentes Conditions Générales d'Utilisation (CGU) définissent les règles d'accès et d'utilisation du site
                <strong>QuestHub</strong>. En créant un compte ou en utilisant la plateforme, vous reconnaissez avoir lu,
                compris et accepté les présentes conditions.
            </p>
        </section>

        <section class="space-y-3">
            <h2 class="text-lg font-bold text-dark-gray flex items-center gap-2 border-b border-stone-100 pb-2 font-display">
                <Shield size={18} class="text-burnt-orange" />
                2. Compte utilisateur et sécurité
            </h2>
            <p class="text-sm">
                L'accès à certaines fonctionnalités nécessite un compte personnel. Vous vous engagez à fournir des informations exactes
                et à conserver la confidentialité de vos identifiants. Toute activité réalisée depuis votre compte est présumée
                effectuée par vous, sauf preuve d'un accès frauduleux.
            </p>
        </section>

        <section class="space-y-3">
            <h2 class="text-lg font-bold text-dark-gray flex items-center gap-2 border-b border-stone-100 pb-2 font-display">
                <Gavel size={18} class="text-burnt-orange" />
                3. Règles de conduite
            </h2>
            <p class="text-sm">Vous vous engagez à ne pas :</p>
            <ul class="list-disc list-inside text-sm pl-2 space-y-1.5">
                <li>publier de contenu illicite, diffamatoire, haineux, discriminatoire ou portant atteinte aux droits de tiers ;</li>
                <li>usurper l'identité d'une autre personne ou contourner les mécanismes d'authentification ;</li>
                <li>perturber le fonctionnement du service (scripts abusifs, tentatives d'intrusion, surcharge volontaire).</li>
            </ul>
        </section>

        <section class="space-y-3">
            <h2 class="text-lg font-bold text-dark-gray flex items-center gap-2 border-b border-stone-100 pb-2 font-display">
                <ScrollText size={18} class="text-burnt-orange" />
                4. Contenus publiés par les utilisateurs
            </h2>
            <p class="text-sm">
                Tout contenu que vous uploadez sur QuestHub (textes, images, fiches, notes, créations) est susceptible d'être vu par
                d'autres utilisateurs selon les espaces où vous le publiez. En conséquence, tout contenu uploadé doit être considéré
                comme public.
            </p>
            <p class="text-sm">
                Vous restez propriétaire des contenus que vous publiez (fiches, messages, notes, créations), mais vous accordez à
                QuestHub une licence non exclusive et gratuite permettant leur affichage, leur stockage et leur diffusion dans le
                cadre strict du fonctionnement de la plateforme.
            </p>
            <p class="text-sm">
                Vous garantissez disposer des droits nécessaires sur les contenus publiés. En cas de signalement légitime, l'éditeur
                peut retirer un contenu sans préavis.
            </p>
        </section>

        <section class="space-y-3">
            <h2 class="text-lg font-bold text-dark-gray flex items-center gap-2 border-b border-stone-100 pb-2 font-display">
                <AlertTriangle size={18} class="text-burnt-orange" />
                5. Suspension, limitation et responsabilité
            </h2>
            <p class="text-sm">
                L'éditeur peut suspendre ou supprimer un compte en cas de non-respect des présentes CGU, de comportement frauduleux
                ou de risque pour la sécurité de la plateforme.
            </p>
            <p class="text-sm">
                QuestHub est fourni "en l'état". Malgré le soin apporté au service, l'éditeur ne peut garantir une disponibilité
                continue ni l'absence totale d'erreurs. L'utilisateur demeure responsable de la sauvegarde de ses contenus importants.
            </p>
        </section>

        <section class="space-y-3">
            <h2 class="text-lg font-bold text-dark-gray flex items-center gap-2 border-b border-stone-100 pb-2 font-display">
                <Gavel size={18} class="text-burnt-orange" />
                6. Évolution des conditions et contact
            </h2>
            <p class="text-sm">
                Les présentes CGU peuvent être mises à jour à tout moment afin de refléter les évolutions du service ou de la
                réglementation. La version en vigueur est celle publiée sur cette page.
            </p>
            <p class="text-sm">
                Pour toute question, vous pouvez contacter l'éditeur à l'adresse
                <a href="mailto:support@questhub.fr" class="text-burnt-orange hover:underline font-semibold">support@questhub.fr</a>.
            </p>
        </section>
    </div>
</main>