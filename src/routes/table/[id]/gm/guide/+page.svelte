<script lang="ts">
    import Header from "$lib/components/Header.svelte";
    import { 
        BookOpen, 
        Heart, 
        Activity, 
        Skull, 
        Shield, 
        RotateCcw, 
        Play, 
        Users, 
        Flame, 
        Eye, 
        EyeOff,
        Trash2,
        X,
        Plus,
        CheckCircle2,
        HelpCircle,
        Code,
        Settings
    } from "lucide-svelte";
    import { page } from "$app/state";

    let activeSection = $state("combat");

    const sections = [
        { id: "combat", label: "Tracker de Combat", icon: Play },
        { id: "hp", label: "Gestion des PV (Dés)", icon: Heart },
        { id: "conditions", label: "États & Effets", icon: Activity },
        { id: "factions", label: "Factions & Camps", icon: Users },
        { id: "visibility", label: "Visibilité & Joueurs", icon: Eye }
    ];
</script>

<div class="h-screen w-screen flex flex-col bg-stone-50 overflow-hidden font-sans text-stone-850">
    <Header />
    
    <div class="flex-1 flex overflow-hidden">
        <!-- Sidebar Navigation -->
        <aside class="w-64 bg-white border-r border-stone-200 flex flex-col p-4 shrink-0 justify-between">
            <div class="space-y-6">
                <div class="flex items-center gap-2 px-2">
                    <BookOpen class="text-burnt-orange" size={20} />
                    <span class="font-display font-bold text-base text-dark-gray">Guide du MJ</span>
                </div>

                <nav class="space-y-1">
                    {#each sections as sec}
                        <button
                            onclick={() => (activeSection = sec.id)}
                            class="w-full flex items-center gap-3 px-3 py-2.5 rounded-xl text-sm font-semibold transition-all text-left cursor-pointer
                            {activeSection === sec.id
                                ? 'bg-stone-900 text-white shadow-md'
                                : 'text-stone-500 hover:bg-stone-100 hover:text-stone-800'}"
                        >
                            <sec.icon size={16} />
                            <span>{sec.label}</span>
                        </button>
                    {/each}
                </nav>
            </div>

            <!-- Back link -->
            <div class="pt-4 border-t border-stone-100">
                <a
                    href={page.url.pathname.replace("/guide", "")}
                    class="flex items-center justify-center gap-2 py-2 px-4 rounded-xl border border-stone-200 text-stone-600 hover:bg-stone-50 hover:text-stone-900 transition-colors text-xs font-bold font-display"
                >
                    Retour à la Table
                </a>
            </div>
        </aside>

        <!-- Main Content Panel -->
        <main class="flex-1 overflow-y-auto bg-stone-50 p-6 md:p-10">
            <div class="max-w-4xl mx-auto space-y-8">
                
                {#if activeSection === "combat"}
                    <section class="space-y-6 animate-in fade-in duration-200">
                        <div class="border-b border-stone-200 pb-4">
                            <h1 class="text-3xl font-display font-extrabold text-dark-gray flex items-center gap-3">
                                <Play class="text-burnt-orange fill-burnt-orange/10" size={32} />
                                Tracker de Combat & Initiatives
                            </h1>
                            <p class="text-stone-500 mt-2">
                                Apprenez à gérer les tours, les initiatives et l'enchaînement des rounds de combat.
                            </p>
                        </div>

                        <!-- Card Step 1 -->
                        <div class="bg-white rounded-2xl border border-stone-200/80 p-6 shadow-sm space-y-4">
                            <h2 class="text-lg font-bold text-stone-900 flex items-center gap-2">
                                <span class="w-6 h-6 rounded-full bg-burnt-orange/10 text-burnt-orange flex items-center justify-center text-xs font-black">1</span>
                                Préparation du Combat
                            </h2>
                            <p class="text-sm text-stone-600 leading-relaxed">
                                Avant de lancer le combat, le tracker est en mode **Préparation**. Une liste de tous les jetons (tokens) présents sur la carte s'affiche.
                            </p>
                            <ul class="list-disc pl-5 space-y-2 text-sm text-stone-600">
                                <li>Cochez ou décochez les cases à côté des jetons pour les inclure ou non dans l'initiative.</li>
                                <li>Ajustez les initiatives initiales directement dans les champs numériques si nécessaire.</li>
                                <li>Cliquez sur le bouton orange <strong>"Commencer le Combat"</strong>. Les initiatives des monstres et PNJ seront tirées automatiquement si elles ne sont pas spécifiées.</li>
                            </ul>
                        </div>

                        <!-- Card Step 2 -->
                        <div class="bg-white rounded-2xl border border-stone-200/80 p-6 shadow-sm space-y-4">
                            <h2 class="text-lg font-bold text-stone-900 flex items-center gap-2">
                                <span class="w-6 h-6 rounded-full bg-burnt-orange/10 text-burnt-orange flex items-center justify-center text-xs font-black">2</span>
                                Déroulement des Tours
                            </h2>
                            <p class="text-sm text-stone-600 leading-relaxed">
                                Une fois le combat lancé, le tracker trie automatiquement les participants par ordre d'initiative décroissante.
                            </p>
                            <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                                <div class="p-4 rounded-xl bg-stone-50 border border-stone-200/60 space-y-2">
                                    <h3 class="text-xs font-bold text-stone-500 uppercase tracking-wider">Tour Actif</h3>
                                    <p class="text-xs text-stone-600">
                                        Le personnage dont c'est le tour est mis en avant avec une bordure orange et un halo lumineux. Un point clignotant s'affiche sur son avatar.
                                    </p>
                                </div>
                                <div class="p-4 rounded-xl bg-stone-50 border border-stone-200/60 space-y-2">
                                    <h3 class="text-xs font-bold text-stone-500 uppercase tracking-wider">Actions Rapides</h3>
                                    <p class="text-xs text-stone-600">
                                        Utilisez le bouton <strong>"Tour Suivant"</strong> pour passer la main, ou la flèche de gauche pour revenir en arrière en cas d'erreur. Le round s'incrémente automatiquement.
                                    </p>
                                </div>
                            </div>
                        </div>
                    </section>

                {:else if activeSection === "hp"}
                    <section class="space-y-6 animate-in fade-in duration-200">
                        <div class="border-b border-stone-200 pb-4">
                            <h1 class="text-3xl font-display font-extrabold text-dark-gray flex items-center gap-3">
                                <Heart class="text-rose-500 fill-rose-500/10" size={32} />
                                Gestion des Points de Vie & Formules de Dés
                            </h1>
                            <p class="text-stone-500 mt-2">
                                Modifiez rapidement la santé des monstres et des joueurs grâce à des raccourcis mathématiques et de dés.
                            </p>
                        </div>

                        <!-- HP Input guide -->
                        <div class="bg-white rounded-2xl border border-stone-200/80 p-6 shadow-sm space-y-4">
                            <h2 class="text-lg font-bold text-stone-900 flex items-center gap-2">
                                <Code class="text-stone-500" size={20} />
                                Syntaxes supportées dans le champ PV
                            </h2>
                            <p class="text-sm text-stone-600">
                                Cliquez sur la valeur de points de vie d'un jeton (par exemple <code>45/60</code>) et saisissez l'une des syntaxes suivantes, puis appuyez sur <strong>Entrée</strong> :
                            </p>
                            
                            <div class="overflow-x-auto">
                                <table class="w-full text-left border-collapse text-sm">
                                    <thead>
                                        <tr class="border-b border-stone-200 text-stone-400 font-semibold">
                                            <th class="py-2.5">Saisie</th>
                                            <th class="py-2.5">Action</th>
                                            <th class="py-2.5">Exemple (pour PV actuel = 20)</th>
                                        </tr>
                                    </thead>
                                    <tbody class="divide-y divide-stone-100 text-stone-600 font-medium">
                                        <tr>
                                            <td class="py-3 font-mono text-stone-900 bg-stone-50/50 px-2 rounded">35</td>
                                            <td class="py-3">Définit directement les PV à cette valeur absolue</td>
                                            <td class="py-3">PV devient <strong>35</strong></td>
                                        </tr>
                                        <tr>
                                            <td class="py-3 font-mono text-stone-900 bg-stone-50/50 px-2 rounded">-5</td>
                                            <td class="py-3">Retire des dégâts directs</td>
                                            <td class="py-3">PV devient 20 - 5 = <strong>15</strong></td>
                                        </tr>
                                        <tr>
                                            <td class="py-3 font-mono text-stone-900 bg-stone-50/50 px-2 rounded">+12</td>
                                            <td class="py-3">Ajoute du soin direct (bridé au maximum du jeton)</td>
                                            <td class="py-3">PV devient 20 + 12 = <strong>32</strong></td>
                                        </tr>
                                        <tr>
                                            <td class="py-3 font-mono text-stone-900 bg-stone-50/50 px-2 rounded">-2d6</td>
                                            <td class="py-3">Lance les dés et applique le total en dégâts</td>
                                            <td class="py-3">Lance 2d6 (ex: 7), PV devient 20 - 7 = <strong>13</strong></td>
                                        </tr>
                                        <tr>
                                            <td class="py-3 font-mono text-stone-900 bg-stone-50/50 px-2 rounded">+1d8+3</td>
                                            <td class="py-3">Lance les dés avec modificateur fixe pour soigner</td>
                                            <td class="py-3">Lance 1d8+3 (ex: 8), PV devient 20 + 8 = <strong>28</strong></td>
                                        </tr>
                                    </tbody>
                                </table>
                            </div>
                        </div>

                        <!-- Graphical Health Bar visual -->
                        <div class="bg-white rounded-2xl border border-stone-200/80 p-6 shadow-sm space-y-4">
                            <h2 class="text-lg font-bold text-stone-900 flex items-center gap-2">
                                <Shield class="text-stone-500" size={20} />
                                Indicateur Visuel & Barre de Vie
                            </h2>
                            <p class="text-sm text-stone-600 leading-relaxed">
                                Chaque participant dispose d'une barre de vie horizontale intégrée directement dans sa carte. Elle change de couleur dynamiquement :
                            </p>
                            <div class="space-y-2.5">
                                <div class="flex items-center gap-3">
                                    <div class="w-24 bg-stone-150 h-2 rounded-full overflow-hidden shrink-0"><div class="bg-emerald-500 h-full w-full"></div></div>
                                    <span class="text-xs text-stone-600">Vert (Plus de 50% des PV) : Personnage en bonne santé.</span>
                                </div>
                                <div class="flex items-center gap-3">
                                    <div class="w-24 bg-stone-150 h-2 rounded-full overflow-hidden shrink-0"><div class="bg-amber-500 h-full w-[40%]"></div></div>
                                    <span class="text-xs text-stone-600">Jaune (Entre 20% et 50% des PV) : Blessures modérées.</span>
                                </div>
                                <div class="flex items-center gap-3">
                                    <div class="w-24 bg-stone-150 h-2 rounded-full overflow-hidden shrink-0"><div class="bg-rose-500 h-full w-[10%]"></div></div>
                                    <span class="text-xs text-stone-600">Rouge (Moins de 20% des PV) : État critique, danger de mort !</span>
                                </div>
                            </div>
                        </div>
                    </section>

                {:else if activeSection === "conditions"}
                    <section class="space-y-6 animate-in fade-in duration-200">
                        <div class="border-b border-stone-200 pb-4">
                            <h1 class="text-3xl font-display font-extrabold text-dark-gray flex items-center gap-3">
                                <Activity class="text-emerald-500 fill-emerald-500/10" size={32} />
                                Gestion des États & Altérations (Conditions)
                            </h1>
                            <p class="text-stone-500 mt-2">
                                Appliquez ou retirez des statuts comme Empoisonné, Invisible, ou créez vos propres états personnalisés.
                            </p>
                        </div>

                        <!-- How to Add/Remove -->
                        <div class="bg-white rounded-2xl border border-stone-200/80 p-6 shadow-sm space-y-4">
                            <h2 class="text-lg font-bold text-stone-900 flex items-center gap-2">
                                <Plus class="text-stone-500" size={20} />
                                Assigner un État
                            </h2>
                            <p class="text-sm text-stone-600 leading-relaxed">
                                Pour attribuer un statut à un personnage :
                            </p>
                            <ol class="list-decimal pl-5 space-y-2 text-sm text-stone-600">
                                <li>Cliquez sur le bouton <strong>"Gérer"</strong> à droite de la section "États & Effets" sur la carte du personnage.</li>
                                <li>Le menu de sélection opaque s'ouvre. Cliquez sur l'un des états prédéfinis (comme <em>Invisible</em>, <em>Étourdi</em>, etc.) pour l'activer.</li>
                                <li>Un indicateur orange s'affiche à côté de l'état actif, et un badge apparaît sur sa fiche.</li>
                            </ol>
                        </div>

                        <!-- Custom Conditions -->
                        <div class="bg-white rounded-2xl border border-stone-200/80 p-6 shadow-sm space-y-4">
                            <h2 class="text-lg font-bold text-stone-900 flex items-center gap-2">
                                <Code class="text-stone-500" size={20} />
                                Conditions Libres (Personnalisées)
                            </h2>
                            <p class="text-sm text-stone-600 leading-relaxed">
                                Si une situation réclame un effet unique (ex: <em>"Bénéficie d'un abri"</em> ou <em>"Effrayé par le Dragon"</em>) :
                            </p>
                            <ul class="list-disc pl-5 space-y-2 text-sm text-stone-600">
                                <li>Dans le menu déroulant, sous <strong>"Condition Libre"</strong>, saisissez le nom de votre effet.</li>
                                <li>Cliquez sur <strong>OK</strong>. L'état personnalisé s'ajoute immédiatement. Le menu se ferme pour fluidifier la partie.</li>
                                <li><em>Note de raccourci :</em> Si vous n'avez rien saisi et que vous cliquez sur <strong>OK</strong>, le menu se fermera simplement.</li>
                            </ul>
                        </div>

                        <!-- Removing conditions -->
                        <div class="bg-white rounded-2xl border border-stone-200/80 p-6 shadow-sm space-y-4">
                            <h2 class="text-lg font-bold text-stone-900 flex items-center gap-2">
                                <Trash2 class="text-stone-500" size={20} />
                                Retirer un État
                            </h2>
                            <p class="text-sm text-stone-600 leading-relaxed">
                                Pour supprimer une condition active, vous pouvez :
                            </p>
                            <ul class="list-disc pl-5 space-y-2 text-sm text-stone-600">
                                <li>Cliquer à nouveau sur l'état dans le menu <strong>"Gérer"</strong> pour le désactiver.</li>
                                <li>Ou cliquer directement sur la petite croix <code>×</code> située à droite du badge de l'état sur la fiche du personnage pour le supprimer instantanément sans ouvrir le menu.</li>
                            </ul>
                        </div>
                    </section>

                {:else if activeSection === "factions"}
                    <section class="space-y-6 animate-in fade-in duration-200">
                        <div class="border-b border-stone-200 pb-4">
                            <h1 class="text-3xl font-display font-extrabold text-dark-gray flex items-center gap-3">
                                <Users class="text-burnt-orange" size={32} />
                                Factions & Camps (Alliés, Ennemis, Neutres)
                            </h1>
                            <p class="text-stone-500 mt-2">
                                Distribuez les jetons dans les bonnes factions pour colorer les anneaux de repérage et structurer le tracker.
                            </p>
                        </div>

                        <!-- Faction controls -->
                        <div class="bg-white rounded-2xl border border-stone-200/80 p-6 shadow-sm space-y-4">
                            <h2 class="text-lg font-bold text-stone-900 flex items-center gap-2">
                                <Shield class="text-stone-500" size={20} />
                                Les trois camps de combat
                            </h2>
                            <p class="text-sm text-stone-600 leading-relaxed">
                                Chaque jeton sur la carte appartient à l'une des 3 factions. La faction influe sur sa couleur d'anneau sur le plateau :
                            </p>
                            <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
                                <div class="p-4 rounded-xl bg-emerald-50/50 border border-emerald-200/60 space-y-2">
                                    <span class="px-2 py-0.5 rounded text-[10px] font-black bg-emerald-100 text-emerald-700 uppercase tracking-wide">Allié</span>
                                    <p class="text-xs text-stone-600">Anneau bleu sur le plateau de jeu. Représente généralement les joueurs et leurs compagnons.</p>
                                </div>
                                <div class="p-4 rounded-xl bg-rose-50/50 border border-rose-200/60 space-y-2">
                                    <span class="px-2 py-0.5 rounded text-[10px] font-black bg-rose-100 text-rose-700 uppercase tracking-wide">Ennemi</span>
                                    <p class="text-xs text-stone-600">Anneau rouge sur le plateau de jeu. Représente les monstres ou les PNJ hostiles.</p>
                                </div>
                                <div class="p-4 rounded-xl bg-amber-50/50 border border-amber-200/60 space-y-2">
                                    <span class="px-2 py-0.5 rounded text-[10px] font-black bg-amber-100 text-amber-700 uppercase tracking-wide">Neutre</span>
                                    <p class="text-xs text-stone-600">Anneau jaune/doré sur le plateau. Représente les civils, monstres neutres ou forces tierces.</p>
                                </div>
                            </div>
                        </div>

                        <!-- Cycling factions -->
                        <div class="bg-white rounded-2xl border border-stone-200/80 p-6 shadow-sm space-y-4">
                            <h2 class="text-lg font-bold text-stone-900 flex items-center gap-2">
                                <RotateCcw class="text-stone-500" size={20} />
                                Comment changer la faction d'un jeton ?
                            </h2>
                            <p class="text-sm text-stone-600 leading-relaxed">
                                Le MJ a deux manières rapides de modifier l'affiliation d'une entité :
                            </p>
                            <ul class="list-disc pl-5 space-y-2 text-sm text-stone-600">
                                <li><strong>Depuis le Tracker de Combat</strong> : Cliquez directement sur le badge de faction (ex: <em>"Ennemi"</em>) sous le nom du personnage. Cela fera défiler la faction dans l'ordre : Allié &rarr; Ennemi &rarr; Neutre.</li>
                                <li><strong>Depuis la carte de jeu</strong> : Faites un clic droit sur le jeton concerné sur le plateau. Sélectionnez l'option <strong>Factions</strong> dans le menu contextuel et choisissez le camp voulu.</li>
                            </ul>
                        </div>
                    </section>

                {:else if activeSection === "visibility"}
                    <section class="space-y-6 animate-in fade-in duration-200">
                        <div class="border-b border-stone-200 pb-4">
                            <h1 class="text-3xl font-display font-extrabold text-dark-gray flex items-center gap-3">
                                <Eye class="text-indigo-500 fill-indigo-500/10" size={32} />
                                Masquage & Visibilité des Joueurs
                            </h1>
                            <p class="text-stone-500 mt-2">
                                Configurez ce que les joueurs peuvent voir ou non sur la carte et dans le combat.
                            </p>
                        </div>

                        <!-- Hide/Show on Board -->
                        <div class="bg-white rounded-2xl border border-stone-200/80 p-6 shadow-sm space-y-4">
                            <h2 class="text-lg font-bold text-stone-900 flex items-center gap-2">
                                <Eye class="text-stone-500" size={20} />
                                Masquage de Jetons (Brouillard de guerre)
                            </h2>
                            <p class="text-sm text-stone-600 leading-relaxed">
                                Cliquez sur l'icône d'œil (<code><Eye size={12} class="inline text-emerald-500" /> / <EyeOff size={12} class="inline text-rose-500" /></code>) située sur la fiche du combatant dans le tracker :
                            </p>
                            <ul class="list-disc pl-5 space-y-2 text-sm text-stone-600">
                                <li><strong>Masqué (Œil barré rouge)</strong> : Le monstre ou PNJ est invisible pour les joueurs. Son jeton n'apparaît pas sur leur carte et son entrée n'apparaît pas dans leur liste d'initiative.</li>
                                <li><strong>Visible (Œil vert)</strong> : Les joueurs peuvent voir le jeton sur le plateau et dans le tracker de combat.</li>
                            </ul>
                        </div>

                        <!-- Hide Health checkbox -->
                        <div class="bg-white rounded-2xl border border-stone-200/80 p-6 shadow-sm space-y-4">
                            <h2 class="text-lg font-bold text-stone-900 flex items-center gap-2">
                                <Settings class="text-stone-500" size={20} />
                                Masquer la santé des monstres aux joueurs
                            </h2>
                            <p class="text-sm text-stone-600 leading-relaxed">
                                Dans le panneau de <strong>Préparation de Combat</strong>, vous disposez de l'option :
                            </p>
                            <div class="p-4 rounded-xl bg-stone-900 text-stone-200 font-medium text-xs border border-stone-800 flex items-center gap-3">
                                <input type="checkbox" checked={true} class="rounded text-burnt-orange w-4 h-4 pointer-events-none" />
                                <span>Masquer complètement la santé des monstres</span>
                            </div>
                            <p class="text-sm text-stone-600 leading-relaxed">
                                Lorsque cette case est **cochée**, les joueurs ne verront aucune valeur numérique (ex: <code>45/45</code>) ni barre de vie sur les jetons identifiés comme **PNJ / Ennemis**. Leurs informations de santé restent confidentielles pour le MJ.
                            </p>
                        </div>
                    </section>
                {/if}

            </div>
        </main>
    </div>
</div>
