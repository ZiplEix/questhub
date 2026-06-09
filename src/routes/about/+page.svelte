<script lang="ts">
    import Header from "$lib/components/Header.svelte";
    import { 
        Sword, 
        Backpack, 
        BookOpen, 
        Users, 
        Sparkles, 
        Layout, 
        Play, 
        CircleDot,
        Info
    } from "lucide-svelte";
</script>

<svelte:head>
    <title>À Propos | QuestHub</title>
    <meta name="description" content="Découvrez QuestHub, la table de jeu virtuelle (VTT) simplifiée pour vos parties de jeu de rôle sur table." />
</svelte:head>

<Header />

<main class="w-full max-w-7xl mx-auto px-4 md:px-8 py-16 grow space-y-20">
    <!-- Hero Header -->
    <div class="text-center space-y-4">
        <span class="inline-flex items-center gap-1.5 px-3 py-1 rounded-full text-[11px] font-extrabold uppercase tracking-widest bg-burnt-orange/10 text-burnt-orange border border-burnt-orange/20">
            <Sparkles size={12} class="animate-pulse" />
            Table de Jeu Virtuelle
        </span>
        <h1 class="text-4xl md:text-5xl font-display font-black text-dark-gray tracking-tight">
            Qu'est-ce que <span class="text-burnt-orange">QuestHub</span> ?
        </h1>
        <p class="text-lg text-stone-500 max-w-2xl mx-auto font-medium leading-relaxed">
            QuestHub est une plateforme web conçue pour simplifier la gestion de vos parties de jeu de rôle sur table. 
            Pas de configurations complexes ni de téléchargements lourds : tout se passe directement dans votre navigateur, en temps réel.
        </p>
    </div>

    <!-- Core Philosophy Section -->
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-12 items-center bg-white p-8 md:p-12 rounded-3xl border border-stone-200/60 shadow-xs">
        <div class="space-y-6">
            <h2 class="text-2xl font-display font-black text-dark-gray">
                La simplicité au service de l'histoire
            </h2>
            <p class="text-stone-600 leading-relaxed text-sm md:text-base">
                Les tables de jeu virtuelles traditionnelles sont souvent des usines à gaz qui ralentissent le rythme de la partie. 
                QuestHub prend le contre-pied en proposant des outils épurés et intuitifs qui s'effacent pour laisser place à la narration :
            </p>
            <ul class="space-y-3.5">
                <li class="flex items-start gap-3 text-sm text-stone-600">
                    <div class="w-5 h-5 rounded-full bg-green-50 border border-green-200 flex items-center justify-center text-green-500 shrink-0 mt-0.5">✓</div>
                    <span><strong>Synchronisation instantanée :</strong> Vos joueurs voient les cartes se dévoiler et les jetons se déplacer sans aucun temps de latence.</span>
                </li>
                <li class="flex items-start gap-3 text-sm text-stone-600">
                    <div class="w-5 h-5 rounded-full bg-green-50 border border-green-200 flex items-center justify-center text-green-500 shrink-0 mt-0.5">✓</div>
                    <span><strong>Ergonomie moderne :</strong> Un design clair inspiré des meilleures applications modernes, adapté aux mobiles comme aux ordinateurs.</span>
                </li>
                <li class="flex items-start gap-3 text-sm text-stone-600">
                    <div class="w-5 h-5 rounded-full bg-green-50 border border-green-200 flex items-center justify-center text-green-500 shrink-0 mt-0.5">✓</div>
                    <span><strong>Zéro barrière d'entrée :</strong> Créez votre compte, lancez une table et invitez vos joueurs via un simple lien ou code d'invitation.</span>
                </li>
            </ul>
        </div>
        <div class="bg-stone-50 border border-stone-200/80 rounded-2xl p-6 md:p-8 flex flex-col justify-center space-y-6">
            <div class="space-y-2">
                <span class="text-xs font-bold text-stone-400 uppercase tracking-wider block">QuestHub en chiffres</span>
                <div class="grid grid-cols-2 gap-4">
                    <div class="bg-white p-4 rounded-xl border border-stone-200/50 shadow-2xs">
                        <p class="text-3xl font-display font-black text-burnt-orange">100%</p>
                        <p class="text-[11px] text-stone-500 font-semibold mt-1">Gratuit & Open-Source</p>
                    </div>
                    <div class="bg-white p-4 rounded-xl border border-stone-200/50 shadow-2xs">
                        <p class="text-3xl font-display font-black text-dark-gray">En ligne</p>
                        <p class="text-[11px] text-stone-500 font-semibold mt-1">Sans installation</p>
                    </div>
                </div>
            </div>
            <div class="bg-burnt-orange/5 border border-burnt-orange/15 rounded-xl p-4 flex gap-3 text-xs text-stone-600 leading-relaxed">
                <Info class="text-burnt-orange shrink-0 mt-0.5" size={16} />
                <p>
                    Que vous soyez Maître de Jeu (MJ) ou joueur, QuestHub s'adapte à vos besoins pour vous offrir une expérience fluide et sans frictions.
                </p>
            </div>
        </div>
    </div>

    <!-- Feature Grid Section -->
    <div class="space-y-10">
        <div class="text-center">
            <h2 class="text-2xl md:text-3xl font-display font-black text-dark-gray">
                Des fonctionnalités taillées pour le jeu
            </h2>
            <p class="text-sm text-stone-500 mt-2 font-medium">
                Découvrez la palette d'outils mis à votre disposition pour animer vos scénarios.
            </p>
        </div>

        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
            <!-- 1. VTT -->
            <div class="bg-white p-6 rounded-2xl border border-stone-200/80 shadow-2xs hover:shadow-xs transition-all space-y-4">
                <div class="w-12 h-12 bg-burnt-orange/15 rounded-xl flex items-center justify-center text-burnt-orange border border-burnt-orange/20 shadow-xs">
                    <Layout size={24} />
                </div>
                <h3 class="text-lg font-display font-black text-dark-gray">Plateau Tactique (VTT)</h3>
                <p class="text-stone-600 text-sm leading-relaxed">
                    Importez vos cartes de combat (Battlemaps) et positionnez vos joueurs en un clin d'œil.
                </p>
                <ul class="space-y-2 text-xs text-stone-500 font-medium pl-1">
                    <li class="flex items-center gap-2">
                        <CircleDot size={12} class="text-burnt-orange" />
                        Ajustement de la grille et calibration de l'échelle (cases/mètres).
                    </li>
                    <li class="flex items-center gap-2">
                        <CircleDot size={12} class="text-burnt-orange" />
                        Jetons de faction (Alliés, Ennemis, PNJ Neutres) colorés dynamiquement.
                    </li>
                    <li class="flex items-center gap-2">
                        <CircleDot size={12} class="text-burnt-orange" />
                        Dessin de cercles d'auras (diamètre en mètres) pour marquer la portée ou les sorts.
                    </li>
                </ul>
            </div>

            <!-- 2. Combat Tracker -->
            <div class="bg-white p-6 rounded-2xl border border-stone-200/80 shadow-2xs hover:shadow-xs transition-all space-y-4">
                <div class="w-12 h-12 bg-burnt-orange/15 rounded-xl flex items-center justify-center text-burnt-orange border border-burnt-orange/20 shadow-xs">
                    <Sword size={24} />
                </div>
                <h3 class="text-lg font-display font-black text-dark-gray">Suivi de Combat & Initiatives</h3>
                <p class="text-stone-600 text-sm leading-relaxed">
                    Gérez les combats complexes de façon fluide grâce au Combat Tracker temps réel.
                </p>
                <ul class="space-y-2 text-xs text-stone-500 font-medium pl-1">
                    <li class="flex items-center gap-2">
                        <CircleDot size={12} class="text-burnt-orange" />
                        Calcul rapide des PV via formules de dés (ex: <code class="bg-stone-100 px-1 py-0.5 rounded text-stone-600">-5</code> ou <code class="bg-stone-100 px-1 py-0.5 rounded text-stone-600">+2d6</code>).
                    </li>
                    <li class="flex items-center gap-2">
                        <CircleDot size={12} class="text-burnt-orange" />
                        Application d'états standards et d'effets personnalisés de combat.
                    </li>
                    <li class="flex items-center gap-2">
                        <CircleDot size={12} class="text-burnt-orange" />
                        Ordre des initiatives et signalement automatique du tour actif.
                    </li>
                </ul>
            </div>

            <!-- 3. Character & Inventories -->
            <div class="bg-white p-6 rounded-2xl border border-stone-200/80 shadow-2xs hover:shadow-xs transition-all space-y-4">
                <div class="w-12 h-12 bg-burnt-orange/15 rounded-xl flex items-center justify-center text-burnt-orange border border-burnt-orange/20 shadow-xs">
                    <Backpack size={24} />
                </div>
                <h3 class="text-lg font-display font-black text-dark-gray">Fiches de Personnages & Inventaires</h3>
                <p class="text-stone-600 text-sm leading-relaxed">
                    Fini les feuilles de papier volantes et les ratures. Vos statistiques sont sauvegardées de façon sécurisée.
                </p>
                <ul class="space-y-2 text-xs text-stone-500 font-medium pl-1">
                    <li class="flex items-center gap-2">
                        <CircleDot size={12} class="text-burnt-orange" />
                        Suivi des pièces d'or, des caractéristiques, des sorts et de l'expérience.
                    </li>
                    <li class="flex items-center gap-2">
                        <CircleDot size={12} class="text-burnt-orange" />
                        Distribution instantanée de butin dans le sac à dos des joueurs par le MJ.
                    </li>
                    <li class="flex items-center gap-2">
                        <CircleDot size={12} class="text-burnt-orange" />
                        Lancer de dés virtuel (3D physique ou journal de jets) intégré.
                    </li>
                </ul>
            </div>

            <!-- 4. Story & Lore -->
            <div class="bg-white p-6 rounded-2xl border border-stone-200/80 shadow-2xs hover:shadow-xs transition-all space-y-4">
                <div class="w-12 h-12 bg-burnt-orange/15 rounded-xl flex items-center justify-center text-burnt-orange border border-burnt-orange/20 shadow-xs">
                    <BookOpen size={24} />
                </div>
                <h3 class="text-lg font-display font-black text-dark-gray">Histoire & Lore (Codex)</h3>
                <p class="text-stone-600 text-sm leading-relaxed">
                    Rédigez l'histoire de vos parties en Markdown et organisez-la au sein d'un codex interactif.
                </p>
                <ul class="space-y-2 text-xs text-stone-500 font-medium pl-1">
                    <li class="flex items-center gap-2">
                        <CircleDot size={12} class="text-burnt-orange" />
                        Éditeur de texte Markdown double-onglet (Édition / Prévisualisation).
                    </li>
                    <li class="flex items-center gap-2">
                        <CircleDot size={12} class="text-burnt-orange" />
                        Classement en dossiers avec réorganisation par glisser-déposer.
                    </li>
                    <li class="flex items-center gap-2">
                        <CircleDot size={12} class="text-burnt-orange" />
                        Révélation de pages ou dossiers entiers en temps réel pour vos joueurs.
                    </li>
                </ul>
            </div>

            <!-- 5. Community Market -->
            <div class="bg-white p-6 rounded-2xl border border-stone-200/80 shadow-2xs hover:shadow-xs transition-all space-y-4">
                <div class="w-12 h-12 bg-burnt-orange/15 rounded-xl flex items-center justify-center text-burnt-orange border border-burnt-orange/20 shadow-xs">
                    <Users size={24} />
                </div>
                <h3 class="text-lg font-display font-black text-dark-gray">Marché Communautaire</h3>
                <p class="text-stone-600 text-sm leading-relaxed">
                    Profitez de l'imagination débordante de la communauté QuestHub.
                </p>
                <ul class="space-y-2 text-xs text-stone-500 font-medium pl-1">
                    <li class="flex items-center gap-2">
                        <CircleDot size={12} class="text-burnt-orange" />
                        Partagez vos personnages pré-tirés et vos monstres avec les autres MJ.
                    </li>
                    <li class="flex items-center gap-2">
                        <CircleDot size={12} class="text-burnt-orange" />
                        Créez des lots (packs) thématiques regroupant plusieurs créations.
                    </li>
                    <li class="flex items-center gap-2">
                        <CircleDot size={12} class="text-burnt-orange" />
                        Importation directe d'éléments ou de packs entiers dans vos campagnes actives.
                    </li>
                </ul>
            </div>
        </div>
    </div>

    <!-- AI Policy Section -->
    <div class="bg-stone-50 border border-stone-200/80 rounded-3xl p-8 md:p-12 space-y-8">
        <div class="flex items-center gap-3">
            <div class="w-10 h-10 bg-burnt-orange/15 rounded-xl flex items-center justify-center text-burnt-orange border border-burnt-orange/20 shadow-xs">
                <Sparkles size={20} />
            </div>
            <h2 class="text-2xl font-display font-black text-dark-gray">
                Politique & Éthique : L'usage de l'I.A. sur QuestHub
            </h2>
        </div>

        <div class="grid grid-cols-1 lg:grid-cols-2 gap-8 text-stone-600 leading-relaxed text-sm md:text-base">
            <div class="space-y-4">
                <h3 class="text-base font-bold text-dark-gray">Une plateforme ouverte, mais respectueuse des artistes</h3>
                <p>
                    QuestHub est né de la volonté de simplifier l'accès au jeu de rôle sur table. Dans cette optique, l'utilisation de l'intelligence artificielle est tolérée et acceptée sur la plateforme, notamment pour aider à la création et à l'illustration rapide de vos campagnes personnelles.
                </p>
                <p>
                    Cependant, nous croyons fermement que rien ne remplace le talent, l'émotion et la vision des <strong>artistes humains</strong>. Dans un monde idéal, nous encourageons vivement les maîtres de jeu et les joueurs à faire vivre les illustrateurs et illustratrices en utilisant leurs créations ou en leur commandant des œuvres originales pour leurs tables de jeu.
                </p>
            </div>
            <div class="space-y-4 bg-white p-6 rounded-2xl border border-stone-200/50 shadow-3xs">
                <h3 class="text-base font-bold text-dark-gray flex items-center gap-2">
                    <Info size={16} class="text-burnt-orange" />
                    Pourquoi un set par défaut généré par I.A. ?
                </h3>
                <p class="text-stone-600 text-sm">
                    Le créateur de QuestHub est un développeur indépendant travaillant seul sur son temps libre. Proposer dès le lancement un catalogue de monstres et de personnages prêt à l'emploi, riche et visuellement cohérent pour que les joueurs puissent démarrer immédiatement des parties, représentait un défi colossal.
                </p>
                <p class="text-stone-600 text-sm">
                    Faire appel à la génération d'images était la seule alternative réaliste pour y parvenir avec nos ressources actuelles. C'est pourquoi l'intégralité du set par défaut (monstres et personnages) a été générée grâce à l'intelligence artificielle <strong>nano banana pro 2</strong>.
                </p>
                <p class="text-stone-600 text-[13px] italic bg-stone-50 p-3 rounded-lg border border-stone-200/50">
                    Nous tenons à cette totale transparence et espérons que ces illustrations temporaires vous permettront de plonger rapidement dans vos aventures, tout en vous incitant à soutenir la communauté des créateurs.
                </p>
            </div>
        </div>
    </div>

    <!-- Call To Action -->
    <div class="bg-stone-900 text-white rounded-3xl p-8 md:p-12 text-center space-y-6 shadow-xl relative overflow-hidden">
        <!-- Visual background grid effect -->
        <div class="absolute inset-0 bg-[linear-gradient(to_right,#ffffff05_1px,transparent_1px),linear-gradient(to_bottom,#ffffff05_1px,transparent_1px)] bg-[size:24px_24px]"></div>
        
        <div class="relative space-y-4">
            <h2 class="text-2xl md:text-3xl font-display font-black tracking-tight">Prêt à lancer vos dés ?</h2>
            <p class="text-sm text-stone-400 max-w-lg mx-auto leading-relaxed font-medium">
                Rejoignez des milliers de rôlistes et donnez une nouvelle dimension à vos aventures de jeu de rôle sur table.
            </p>
            <div class="pt-4 flex flex-col sm:flex-row gap-3 justify-center items-center">
                <a
                    href="/dashboard"
                    class="px-6 py-3 bg-burnt-orange text-white font-bold rounded-xl text-sm hover:bg-opacity-95 shadow-md hover:shadow-lg transition-all flex items-center gap-2 cursor-pointer"
                >
                    <Play size={16} class="fill-white" />
                    Commencer une partie
                </a>
            </div>
        </div>
    </div>
</main>
