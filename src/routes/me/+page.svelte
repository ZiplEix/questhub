<script lang="ts">
    import { onMount } from "svelte";
    import { fetchUserStats, fetchUserCampaigns, fetchUserCharacters, fetchUserMonsters } from "$lib/api";
    import { authClient } from "$lib/auth-client";
    import { goto } from "$app/navigation";
    import Header from "$lib/components/Header.svelte";
    import {
        User,
        Activity,
        Swords,
        Skull,
        BookOpen,
        Mail,
    } from "lucide-svelte";
    import type { SessionUser } from "$lib/types/session-user";

    let activeTab = $state("general");
    let user = $state<SessionUser | null>(null);
    let stats = $state<any>(null);
    let campaigns = $state<any[]>([]);
    let characters = $state<any[]>([]);
    let monsters = $state<any[]>([]);
    let loading = $state(true);

    const tabs = [
        { id: "general", label: "Informations", icon: User },
        { id: "stats", label: "Statistiques", icon: Activity },
        { id: "campaigns", label: "Campagnes", icon: Swords },
        { id: "characters", label: "Personnages", icon: Skull },
        { id: "bestiary", label: "Bestiaire", icon: BookOpen },
    ];

    onMount(async () => {
        const urlParams = new URLSearchParams(window.location.search);
        const tab = urlParams.get("tab");
        if (tab && tabs.find((t) => t.id === tab)) {
            activeTab = tab;
        }

        try {
            const { data, error } = await authClient.getSession();
            if (error || !data?.user) {
                goto("/login");
                return;
            }
            user = data.user;

            await Promise.all([fetchStats(), fetchCampaigns(), fetchCharacters(), fetchMonsters()]);
        } catch (e) {
            console.error(e);
        } finally {
            loading = false;
        }
    });

    // Tab switcher
    function setTab(tabId: string) {
        activeTab = tabId;
        const url = new URL(window.location.href);
        url.searchParams.set("tab", tabId);
        goto(url.toString(), { replaceState: true, keepFocus: true });
    }

    async function fetchStats() {
        try {
            stats = await fetchUserStats();
        } catch (e) {
            console.error("Failed to fetch stats", e);
        }
    }

    async function fetchCampaigns() {
        try {
            campaigns = await fetchUserCampaigns();
        } catch (e) {
            console.error("Failed to fetch campaigns", e);
        }
    }

    async function fetchCharacters() {
        try {
            characters = await fetchUserCharacters();
        } catch (e) {
            console.error("Failed to fetch characters", e);
        }
    }

    async function fetchMonsters() {
        try {
            monsters = await fetchUserMonsters();
        } catch (e) {
            console.error("Failed to fetch monsters", e);
        }
    }
</script>

<div class="min-h-screen bg-cream">
    <Header />

    <main class="max-w-7xl mx-auto px-4 md:px-8 py-12">
        <!-- Header -->
        <div class="mb-8">
            <h1 class="font-display font-bold text-4xl text-dark-gray mb-2">
                Mon Profil
            </h1>
            <p class="text-stone-500">
                Gérez vos informations, vos statistiques et vos campagnes.
            </p>
        </div>

        {#if loading}
            <div class="flex justify-center items-center h-64">
                <div
                    class="animate-spin rounded-full h-12 w-12 border-b-2 border-burnt-orange"
                ></div>
            </div>
        {:else}
            <!-- Tabs Navigation -->
            <div
                class="flex gap-2 mb-8 border-b border-stone-200 overflow-x-auto"
            >
                {#each tabs as tab}
                    <button
                        onclick={() => setTab(tab.id)}
                        class="flex items-center gap-2 px-4 py-3 font-medium text-sm transition-all relative {activeTab ===
                        tab.id
                            ? 'text-burnt-orange'
                            : 'text-stone-500 hover:text-dark-gray'}"
                    >
                        <tab.icon size={18} />
                        {tab.label}
                        {#if activeTab === tab.id}
                            <div
                                class="absolute bottom-0 left-0 w-full h-0.5 bg-burnt-orange rounded-t-full"
                            ></div>
                        {/if}
                    </button>
                {/each}
            </div>

            <!-- Content Area -->
            <div
                class="bg-white rounded-2xl shadow-sm border border-stone-100 p-6 md:p-8 min-h-[500px]"
            >
                <!-- General Info -->
                {#if activeTab === "general"}
                    <div
                        class="animate-in fade-in slide-in-from-bottom-4 duration-300"
                    >
                        <div class="flex items-start gap-8">
                            <div class="relative">
                                <img
                                    src={user?.image ||
                                        `https://ui-avatars.com/api/?name=${user?.name}&background=random`}
                                    alt={user?.name}
                                    class="w-32 h-32 rounded-full object-cover border-4 border-stone-100 shadow-md"
                                />
                            </div>
                            <div class="space-y-4 flex-1">
                                <div>
                                    <label
                                        class="block text-sm font-medium text-dark-gray/40 mb-1"
                                        >Nom d'utilisateur</label
                                    >
                                    <p class="text-xl font-bold text-dark-gray">
                                        {user?.name}
                                    </p>
                                </div>
                                <div>
                                    <label
                                        class="block text-sm font-medium text-dark-gray/40 mb-1"
                                        >Adresse Email</label
                                    >
                                    <div
                                        class="flex items-center gap-2 text-dark-gray"
                                    >
                                        <Mail size={18} />
                                        <span>{user?.email}</span>
                                    </div>
                                </div>
                            </div>
                        </div>
                    </div>
                {/if}

                <!-- Stats -->
                {#if activeTab === "stats"}
                    <div
                        class="animate-in fade-in slide-in-from-bottom-4 duration-300"
                    >
                        <div
                            class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6"
                        >
                            <div
                                class="bg-stone-50 rounded-2xl p-6 border border-stone-100 hover:shadow-md transition-shadow"
                            >
                                <p class="text-dark-gray/60 font-medium mb-2">
                                    Parties Jouées
                                </p>
                                <p
                                    class="text-4xl font-display font-bold text-burnt-orange"
                                >
                                    {stats?.games_played || 0}
                                </p>
                            </div>
                            <div
                                class="bg-stone-50 rounded-2xl p-6 border border-stone-100 hover:shadow-md transition-shadow"
                            >
                                <p class="text-dark-gray/60 font-medium mb-2">
                                    Parties en tant que GM
                                </p>
                                <p
                                    class="text-4xl font-display font-bold text-burnt-orange"
                                >
                                    {stats?.games_mastered || 0}
                                </p>
                            </div>
                            <div
                                class="bg-stone-50 rounded-2xl p-6 border border-stone-100 hover:shadow-md transition-shadow"
                            >
                                <p class="text-dark-gray/60 font-medium mb-2">
                                    Personnages Créés
                                </p>
                                <p
                                    class="text-4xl font-display font-bold text-burnt-orange"
                                >
                                    {stats?.characters_created || 0}
                                </p>
                            </div>
                        </div>
                    </div>
                {/if}

                <!-- Campaigns -->
                {#if activeTab === "campaigns"}
                    <div
                        class="animate-in fade-in slide-in-from-bottom-4 duration-300"
                    >
                        {#if campaigns.length > 0}
                            <div class="grid gap-4">
                                {#each campaigns as camp}
                                    <div
                                        class="flex items-center gap-4 p-4 rounded-2xl border border-stone-100 hover:bg-stone-50 transition-colors"
                                    >
                                        <img
                                            src={camp.game_image_url ||
                                                `https://placehold.co/100x100/3D405B/F9F7F2?text=${camp.game_name}`}
                                            alt={camp.game_name}
                                            class="w-16 h-16 rounded-xl object-cover"
                                        />
                                        <div class="flex-1">
                                            <h3
                                                class="font-bold text-lg text-dark-gray"
                                            >
                                                {camp.game_name}
                                            </h3>
                                            <div
                                                class="flex items-center gap-2 text-sm text-dark-gray/60"
                                            >
                                                <span
                                                    >Personnage: <span
                                                        class="font-medium text-burnt-orange"
                                                        >{camp.character_name}</span
                                                    ></span
                                                >
                                            </div>
                                        </div>
                                        <a
                                            href="/table/{camp.game_id}"
                                            class="px-4 py-2 bg-dark-gray text-white rounded-xl text-sm hover:bg-opacity-90 transition-opacity"
                                        >
                                            Rejoindre
                                        </a>
                                    </div>
                                {/each}
                            </div>
                        {:else}
                            <div class="text-center py-12 text-dark-gray/40">
                                <p>Aucune campagne trouvée.</p>
                            </div>
                        {/if}
                    </div>
                {/if}

                <!-- Characters Tab -->
                {#if activeTab === "characters"}
                    <div
                        class="animate-in fade-in slide-in-from-bottom-4 duration-300"
                    >
                        {#if characters.length > 0}
                            <div class="grid grid-cols-1 sm:grid-cols-2 gap-6">
                                {#each characters as char}
                                    <div
                                        class="bg-white rounded-2xl border border-stone-100 p-5 shadow-xs hover:shadow-md hover:border-burnt-orange/20 transition-all flex flex-col justify-between h-full group"
                                    >
                                        <div class="space-y-4">
                                            <!-- Top section -->
                                            <div class="flex items-start justify-between gap-4">
                                                <div class="flex items-center gap-4">
                                                    <img
                                                        src={char.avatar_url || `https://api.dicebear.com/9.x/avataaars/svg?seed=${char.name}`}
                                                        alt={char.name}
                                                        class="w-16 h-16 rounded-xl object-cover bg-stone-50 border-2 border-stone-100"
                                                    />
                                                    <div>
                                                        <h3
                                                            class="font-display font-bold text-xl text-dark-gray group-hover:text-burnt-orange transition-colors"
                                                        >
                                                            {char.name}
                                                        </h3>
                                                        <p class="text-stone-400 text-sm font-medium">
                                                            {char.race} {char.sub_race ? `(${char.sub_race})` : ''}
                                                        </p>
                                                    </div>
                                                </div>
                                                <div>
                                                    {#if char.association_type === 'assigned'}
                                                        <span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-semibold bg-blue-50 text-blue-700 border border-blue-100">
                                                            Attribué
                                                        </span>
                                                    {:else}
                                                        <span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-semibold bg-purple-50 text-purple-700 border border-purple-100">
                                                            Créé (MJ)
                                                        </span>
                                                    {/if}
                                                </div>
                                            </div>
 
                                            <!-- Vitals (HP & XP/Level) -->
                                            <div class="space-y-2 bg-stone-50 p-3 rounded-xl border border-stone-100">
                                                <div class="flex justify-between text-xs font-bold text-stone-500">
                                                    <span>Points de Vie</span>
                                                    <span>{char.current_hp} / {char.max_hp} PV</span>
                                                </div>
                                                <div class="h-2 bg-stone-200 rounded-full overflow-hidden">
                                                    <div
                                                        class="h-full rounded-full transition-all duration-300 {char.current_hp < char.max_hp / 3 ? 'bg-red-500' : 'bg-green-500'}"
                                                        style="width: {(char.current_hp / char.max_hp) * 100}%"
                                                    ></div>
                                                </div>
                                                {#if char.experience !== undefined}
                                                    <div class="flex justify-between text-[11px] text-stone-400 font-semibold pt-1 border-t border-stone-200/50 mt-1">
                                                        <span>Expérience</span>
                                                        <span>{char.experience} XP</span>
                                                    </div>
                                                {/if}
                                            </div>
 
                                            <!-- Associated Game info -->
                                            {#if char.game}
                                                <div class="flex items-center gap-2 pt-2">
                                                    <img
                                                        src={char.game.image_url || `https://placehold.co/100x100/3D405B/F9F7F2?text=${char.game.name}`}
                                                        alt={char.game.name}
                                                        class="w-6 h-6 rounded-md object-cover"
                                                    />
                                                    <span class="text-xs font-medium text-stone-500 truncate">
                                                        Partie : <span class="font-semibold text-dark-gray">{char.game.name}</span>
                                                    </span>
                                                </div>
                                            {/if}
                                        </div>
 
                                        <!-- Join button -->
                                        <div class="pt-5 mt-4 border-t border-stone-100/60 flex justify-end">
                                            <a
                                                href="/table/{char.game_id || char.game?.id}"
                                                class="px-4 py-2 bg-burnt-orange text-white rounded-xl text-sm font-bold shadow-xs hover:bg-opacity-90 hover:shadow-md hover:-translate-y-0.5 transition-all flex items-center gap-2"
                                            >
                                                {#if char.association_type === 'assigned'}
                                                    Jouer ce personnage
                                                {:else}
                                                    Rejoindre la partie (MJ)
                                                {/if}
                                            </a>
                                        </div>
                                    </div>
                                {/each}
                            </div>
                        {:else}
                            <div class="text-center py-20 bg-stone-50/50 rounded-2xl border border-stone-100 flex flex-col items-center justify-center gap-4">
                                <div class="w-16 h-16 bg-white rounded-full flex items-center justify-center text-stone-300 shadow-sm border border-stone-100">
                                    <Skull size={32} />
                                </div>
                                <div class="space-y-1">
                                    <h4 class="font-bold text-dark-gray text-lg">Aucun personnage</h4>
                                    <p class="text-stone-400 text-sm max-w-sm">
                                        Vous n'avez aucun personnage pour le moment. Rejoignez une partie ou demandez à votre Maître du Jeu de vous en assigner un !
                                    </p>
                                </div>
                            </div>
                        {/if}
                    </div>
                {/if}

                <!-- Bestiary Tab -->
                {#if activeTab === "bestiary"}
                    <div
                        class="animate-in fade-in slide-in-from-bottom-4 duration-300"
                    >
                        {#if monsters.length > 0}
                            <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6">
                                {#each monsters as monster}
                                    <div
                                        class="bg-white rounded-2xl border border-stone-100 p-5 shadow-xs hover:shadow-md hover:border-burnt-orange/20 transition-all flex flex-col justify-between h-full group"
                                    >
                                        <div class="space-y-4">
                                            <!-- Top section -->
                                            <div class="flex items-start justify-between gap-4">
                                                <div class="flex items-center gap-4">
                                                    <img
                                                        src={monster.avatar_url || `https://api.dicebear.com/9.x/bottts/svg?seed=${monster.name}`}
                                                        alt={monster.name}
                                                        class="w-16 h-16 rounded-xl object-cover bg-stone-50 border-2 border-stone-100"
                                                    />
                                                    <div>
                                                        <h3
                                                            class="font-display font-bold text-xl text-dark-gray group-hover:text-burnt-orange transition-colors"
                                                        >
                                                            {monster.name}
                                                        </h3>
                                                        <p class="text-stone-400 text-sm font-medium">
                                                            {monster.race || 'Monstre'}
                                                        </p>
                                                    </div>
                                                </div>
                                            </div>

                                            <!-- Vitals (HP & AC) -->
                                            <div class="grid grid-cols-2 gap-2 bg-stone-50 p-3 rounded-xl border border-stone-100 text-center">
                                                <div>
                                                    <span class="block text-[10px] uppercase tracking-wider text-stone-400 font-bold">Classe d'armure</span>
                                                    <span class="text-lg font-bold text-dark-gray">{monster.armor_class ?? 10}</span>
                                                </div>
                                                <div>
                                                    <span class="block text-[10px] uppercase tracking-wider text-stone-400 font-bold">Points de Vie Max</span>
                                                    <span class="text-lg font-bold text-dark-gray">{monster.max_hp} PV</span>
                                                </div>
                                            </div>

                                            <!-- Associated Game info -->
                                            {#if monster.game}
                                                <div class="flex items-center gap-2 pt-2">
                                                    <img
                                                        src={monster.game.image_url || `https://placehold.co/100x100/3D405B/F9F7F2?text=${monster.game.name}`}
                                                        alt={monster.game.name}
                                                        class="w-6 h-6 rounded-md object-cover"
                                                    />
                                                    <span class="text-xs font-medium text-stone-500 truncate font-semibold">
                                                        Partie : <span class="font-bold text-dark-gray">{monster.game.name}</span>
                                                    </span>
                                                </div>
                                            {/if}
                                        </div>

                                        <!-- Actions -->
                                        <div class="pt-5 mt-4 border-t border-stone-100/60 flex justify-between items-center gap-2">
                                            <a
                                                href="/table/{monster.game_id || monster.game?.id}/gm/settings?tab=bestiary"
                                                class="text-xs font-bold text-stone-500 hover:text-burnt-orange transition-colors"
                                            >
                                                Modifier le monstre
                                            </a>
                                            <a
                                                href="/table/{monster.game_id || monster.game?.id}"
                                                class="px-3 py-1.5 bg-dark-gray text-white rounded-lg text-xs font-bold shadow-xs hover:bg-opacity-90 hover:shadow-md transition-all"
                                            >
                                                Aller à la table
                                            </a>
                                        </div>
                                    </div>
                                {/each}
                            </div>
                        {:else}
                            <div class="text-center py-20 bg-stone-50/50 rounded-2xl border border-stone-100 flex flex-col items-center justify-center gap-4">
                                <div class="w-16 h-16 bg-white rounded-full flex items-center justify-center text-stone-300 shadow-sm border border-stone-100">
                                    <BookOpen size={32} />
                                </div>
                                <div class="space-y-1">
                                    <h4 class="font-bold text-dark-gray text-lg">Aucun monstre créé</h4>
                                    <p class="text-stone-400 text-sm max-w-sm">
                                        Vous n'avez pas encore créé de monstre. Allez dans les paramètres d'une partie dont vous êtes le GM pour en ajouter au bestiaire !
                                    </p>
                                </div>
                            </div>
                        {/if}
                    </div>
                {/if}

            </div>
        {/if}
    </main>
</div>
