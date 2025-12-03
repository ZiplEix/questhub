<script lang="ts">
    import { page } from "$app/state";
    import Header from "$lib/components/Header.svelte";
    import { onMount } from "svelte";
    import { api } from "$lib/api";
    import { authClient } from "$lib/auth-client";
    import { goto } from "$app/navigation";
    import GeneralTab from "$lib/components/game/gm/settings/GeneralTab.svelte";
    import PlayersTab from "$lib/components/game/gm/settings/PlayersTab.svelte";
    import CharactersTab from "$lib/components/game/gm/settings/CharactersTab.svelte";
    import InvitationsTab from "$lib/components/game/gm/settings/InvitationsTab.svelte";
    import BestiaryTab from "$lib/components/game/gm/settings/BestiaryTab.svelte";
    import { Settings, Users, Mail, User, Skull } from "lucide-svelte";

    let activeTab = $state("general");
    let game = $state<any>(null);
    let loading = $state(true);
    let players = $state<any[]>([]);
    let invitations = $state<any[]>([]);

    const tabs = [
        { id: "general", label: "Général", icon: Settings },
        { id: "players", label: "Joueurs", icon: Users },
        { id: "characters", label: "Personnages", icon: User },
        { id: "bestiary", label: "Bestiaire", icon: Skull },
        { id: "invitations", label: "Invitations", icon: Mail },
    ];

    async function fetchGame(id: string, token: string) {
        try {
            const response = await api.get(`/table/${id}`, {
                headers: {
                    Authorization: `Bearer ${token}`,
                },
            });
            game = response.data;
        } catch (error) {
            console.error("Failed to fetch game:", error);
        }
    }

    async function fetchPlayers(id: string, token: string) {
        try {
            const response = await api.get(`/table/${id}/players`, {
                headers: {
                    Authorization: `Bearer ${token}`,
                },
            });
            players = response.data;
        } catch (error) {
            console.error("Failed to fetch players:", error);
        }
    }

    async function fetchInvitations(id: string, token: string) {
        try {
            const response = await api.get(`/table/${id}/invitations`, {
                headers: {
                    Authorization: `Bearer ${token}`,
                },
            });
            invitations = response.data;
        } catch (error) {
            console.error("Failed to fetch invitations:", error);
        }
    }

    async function refreshData() {
        const gameId = page.params.id;
        if (gameId) {
            const { data: tokenData } = await authClient.token();
            if (tokenData?.token) {
                await Promise.all([
                    fetchPlayers(gameId, tokenData.token),
                    fetchInvitations(gameId, tokenData.token),
                ]);
            }
        }
    }

    onMount(async () => {
        const gameId = page.params.id;
        if (gameId) {
            try {
                const { data: tokenData } = await authClient.token();
                if (tokenData?.token) {
                    await Promise.all([
                        fetchGame(gameId, tokenData.token),
                        fetchPlayers(gameId, tokenData.token),
                        fetchInvitations(gameId, tokenData.token),
                    ]);
                }
            } catch (error) {
                console.error("Failed to fetch game data:", error);
            } finally {
                loading = false;
            }
        }
    });

    // Replaced setActiveTab with setTab from the provided edit
    function setTab(tabId: string) {
        activeTab = tabId;
        goto(`?tab=${tabId}`, { replaceState: true, keepFocus: true });
    }
</script>

<div class="min-h-screen bg-stone-50">
    <Header />

    <main class="max-w-4xl mx-auto p-6 md:p-12">
        <!-- Header -->
        <div class="mb-8">
            <h1 class="text-3xl font-display font-bold text-dark-gray mb-2">
                Paramètres de la partie
            </h1>
            <p class="text-stone-500">
                Gérez les configurations, les joueurs et les accès à votre
                table.
            </p>
        </div>

        {#if loading}
            <div class="flex justify-center py-12">
                <div
                    class="animate-spin rounded-full h-8 w-8 border-b-2 border-burnt-orange"
                ></div>
            </div>
        {:else if game}
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
                class="bg-white rounded-2xl shadow-sm border border-stone-100 p-6 md:p-8"
            >
                <!-- GENERAL TAB -->
                {#if activeTab === "general"}
                    <GeneralTab bind:game />
                {/if}

                <!-- PLAYERS TAB -->
                {#if activeTab === "players"}
                    <PlayersTab
                        {players}
                        {invitations}
                        onRefresh={refreshData}
                    />
                {/if}

                <!-- CHARACTERS TAB -->
                {#if activeTab === "characters"}
                    <CharactersTab {players} gameId={page.params.id || ""} />
                {/if}

                <!-- BESTIARY TAB -->
                {#if activeTab === "bestiary"}
                    <BestiaryTab gameId={page.params.id || ""} />
                {/if}

                <!-- INVITATIONS TAB -->
                {#if activeTab === "invitations"}
                    <InvitationsTab {invitations} onRefresh={refreshData} />
                {/if}
            </div>
        {/if}
    </main>
</div>
