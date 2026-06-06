<script lang="ts">
    import Header from "$lib/components/Header.svelte";
    import { Scale, ShieldAlert, FileText, Landmark, KeyRound } from "lucide-svelte";
</script>

<svelte:head>
    <title>Mentions Légales | QuestHub</title>
    <meta name="description" content="Mentions légales, conditions d'utilisation et politique de protection des données de QuestHub." />
</svelte:head>

<Header />

<main class="w-full max-w-7xl mx-auto px-4 md:px-8 py-16 grow">
    <div class="mb-10 text-center">
        <h1 class="text-3xl font-display font-black text-burnt-orange tracking-tight flex items-center justify-center gap-3">
            <Scale class="text-burnt-orange" size={28} />
            Mentions Légales
        </h1>
        <p class="text-sm text-stone-500 mt-2 font-medium">
            En vigueur au 7 juin 2026. Informations légales sur l'éditeur, l'hébergeur et la protection des données de QuestHub.
        </p>
    </div>

    <div class="bg-white rounded-2xl border border-stone-200/80 shadow-xs p-6 md:p-10 space-y-10 leading-relaxed text-stone-700 max-w-3xl mx-auto">
        <!-- 1. Éditeur -->
        <section class="space-y-3">
            <h2 class="text-lg font-bold text-dark-gray flex items-center gap-2 border-b border-stone-100 pb-2 font-display">
                <Landmark size={18} class="text-burnt-orange" />
                1. Éditeur du Site
            </h2>
            <p class="text-sm">
                Le site internet <strong>QuestHub</strong> est un projet personnel développé et édité à titre non professionnel par <strong>Baptiste</strong>.
            </p>
            <p class="text-sm mt-1">
                <strong>Directeur de la publication :</strong> Baptiste. <br />
                <strong>Contact Email :</strong> <a href="mailto:support@questhub.fr" class="text-burnt-orange hover:underline font-semibold">support@questhub.fr</a>
            </p>
        </section>

        <!-- 2. Hébergeur -->
        <section class="space-y-3">
            <h2 class="text-lg font-bold text-dark-gray flex items-center gap-2 border-b border-stone-100 pb-2 font-display">
                <KeyRound size={18} class="text-burnt-orange" />
                2. Hébergement
            </h2>
            <p class="text-sm">
                L'infrastructure applicative et le serveur d'hébergement du site sont fournis par :
            </p>
            <ul class="list-disc list-inside text-sm pl-2 space-y-1.5">
                <li><strong>Hébergement Frontend :</strong> Vercel Inc., 340 S Lemon Ave #4133 Walnut, CA 91789, États-Unis.</li>
                <li><strong>Hébergement Base de Données & Stockage :</strong> Supabase Inc., 970 Topper Ln, Lafayette, CA 94549, États-Unis (serveurs localisés au sein de l'Union Européenne - Frankfurt, Allemagne).</li>
            </ul>
        </section>

        <!-- 3. Propriété Intellectuelle -->
        <section class="space-y-3">
            <h2 class="text-lg font-bold text-dark-gray flex items-center gap-2 border-b border-stone-100 pb-2 font-display">
                <FileText size={18} class="text-burnt-orange" />
                3. Propriété Intellectuelle
            </h2>
            <p class="text-sm">
                L'ensemble du contenu présent sur le site (logos, designs graphiques, codes sources, icônes) est la propriété 
                exclusive de l'éditeur ou fait l'objet d'une autorisation d'utilisation (notamment les illustrations des avatars et 
                des créatures de la communauté). 
            </p>
            <p class="text-sm">
                Toute reproduction, distribution, modification ou adaptation du contenu sans accord préalable écrit de l'éditeur 
                est strictement interdite et constitue une contrefaçon sanctionnée par le Code de la propriété intellectuelle.
            </p>
        </section>

        <!-- 4. Protection des Données (RGPD) -->
        <section class="space-y-3">
            <h2 class="text-lg font-bold text-dark-gray flex items-center gap-2 border-b border-stone-100 pb-2 font-display">
                <ShieldAlert size={18} class="text-burnt-orange" />
                4. Protection des Données Personnelles (RGPD)
            </h2>
            <p class="text-sm">
                L'éditeur s'engage à ce que la collecte et le traitement de vos données soient conformes au Règlement Général sur la 
                Protection des Données (RGPD) et à la loi Informatique et Libertés.
            </p>
            <p class="text-sm">
                <strong>Données collectées :</strong> Adresse email, nom d'utilisateur (pour l'authentification et l'attribution des jetons 
                sur les plateaux de jeu ou publications de créations).
            </p>
            <p class="text-sm">
                <strong>Vos droits :</strong> Conformément à la réglementation, vous disposez d'un droit d'accès, de rectification, 
                d'effacement (droit à l'oubli) et de portabilité de vos données personnelles. Vous pouvez exercer ce droit à tout moment 
                en nous écrivant à <a href="mailto:support@questhub.fr" class="text-burnt-orange hover:underline font-semibold">support@questhub.fr</a> ou en supprimant votre compte depuis vos paramètres.
            </p>
        </section>

        <!-- 5. Cookies -->
        <section class="space-y-3">
            <h2 class="text-lg font-bold text-dark-gray flex items-center gap-2 border-b border-stone-100 pb-2 font-display">
                <Scale size={18} class="text-burnt-orange" />
                5. Traceurs et Cookies
            </h2>
            <p class="text-sm">
                Le site n'utilise aucun traceur publicitaire ou de profilage tiers. Les seuls cookies techniques stockés sur votre 
                ordinateur sont strictement nécessaires à l'authentification et au maintien de votre session de jeu (Supabase Auth tokens). 
                À ce titre, l'obtention préalable de votre consentement n'est pas requise.
            </p>
        </section>
    </div>
</main>
