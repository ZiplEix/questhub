<script lang="ts">
    import { onMount } from "svelte";
    import { api } from "$lib/api";
    import { authClient } from "$lib/auth-client";
    import { goto } from "$app/navigation";
    import Header from "$lib/components/Header.svelte";
    import {
        User,
        Activity,
        Swords,
        BookOpen,
        Skull,
        Mail,
    } from "lucide-svelte";
    import type { SessionUser } from "$lib/types/session-user";

    let activeTab = $state("general");
    let user = $state<SessionUser | null>(null);
    let stats = $state<any>(null);
    let campaigns = $state<any[]>([]);
    let loading = $state(true);

    const tabs = [
        { id: "general", label: "Informations", icon: User },
        { id: "stats", label: "Statistiques", icon: Activity },
        { id: "campaigns", label: "Campagnes", icon: Swords },
        { id: "bestiary", label: "Bestiaire", icon: BookOpen },
        { id: "characters", label: "Personnages", icon: Skull },
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

            const tokenResponse = await authClient.token();
            if (tokenResponse.data?.token) {
                const token = tokenResponse.data.token;
                await Promise.all([fetchStats(token), fetchCampaigns(token)]);
            }
        } catch (e) {
            console.error(e);
        } finally {
            loading = false;
        }
    });

    function setTab(tabId: string) {
        activeTab = tabId;
        const url = new URL(window.location.href);
        url.searchParams.set("tab", tabId);
        goto(url.toString(), { replaceState: true, keepFocus: true });
    }

    async function fetchStats(token: string) {
        try {
            const res = await api.get("/user/stats", {
                headers: { Authorization: `Bearer ${token}` },
            });
            stats = res.data;
        } catch (e) {
            console.error("Failed to fetch stats", e);
        }
    }

    async function fetchCampaigns(token: string) {
        try {
            const res = await api.get("/user/campaigns", {
                headers: { Authorization: `Bearer ${token}` },
            });
            campaigns = res.data;
        } catch (e) {
            console.error("Failed to fetch campaigns", e);
        }
    }
</script>

<div class="min-h-screen bg-stone-50">
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
                                            class="px-4 py-2 bg-dark-gray text-white rounded-lg text-sm hover:bg-opacity-90 transition-opacity"
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

                <!-- Bestiary (Mock) -->
                {#if activeTab === "bestiary"}
                    <div
                        class="animate-in fade-in slide-in-from-bottom-4 duration-300"
                    >
                        <div
                            class="grid grid-cols-1 sm:grid-cols-2 gap-4 opacity-50 pointer-events-none"
                        >
                            {#each Array(4) as _}
                                <div
                                    class="p-4 rounded-2xl border border-stone-100 bg-stone-50"
                                >
                                    <div
                                        class="h-4 bg-stone-200 rounded w-1/3 mb-2"
                                    ></div>
                                    <div
                                        class="h-32 bg-stone-200 rounded-xl mb-2"
                                    ></div>
                                    <div
                                        class="h-3 bg-stone-200 rounded w-full"
                                    ></div>
                                </div>
                            {/each}
                        </div>
                        <div class="text-center mt-8 text-dark-gray/60">
                            <p>Fonctionnalité à venir...</p>
                        </div>
                    </div>
                {/if}

                <!-- Characters (Mock) -->
                {#if activeTab === "characters"}
                    <div
                        class="animate-in fade-in slide-in-from-bottom-4 duration-300"
                    >
                        <div
                            class="grid grid-cols-1 sm:grid-cols-2 gap-4 opacity-50 pointer-events-none"
                        >
                            {#each Array(4) as _}
                                <div
                                    class="p-4 rounded-2xl border border-stone-100 bg-stone-50"
                                >
                                    <div class="flex items-center gap-3 mb-4">
                                        <div
                                            class="w-12 h-12 bg-stone-200 rounded-full"
                                        ></div>
                                        <div
                                            class="h-4 bg-stone-200 rounded w-1/3"
                                        ></div>
                                    </div>
                                    <div class="space-y-2">
                                        <div
                                            class="h-3 bg-stone-200 rounded w-full"
                                        ></div>
                                        <div
                                            class="h-3 bg-stone-200 rounded w-5/6"
                                        ></div>
                                    </div>
                                </div>
                            {/each}
                        </div>
                        <div class="text-center mt-8 text-dark-gray/60">
                            <p>Fonctionnalité à venir...</p>
                        </div>
                    </div>
                {/if}
            </div>
        {/if}
    </main>
</div>
