<script lang="ts">
    import { page } from "$app/state";
    import Header from "$lib/components/Header.svelte";
    import { onMount } from "svelte";
    import { fetchGame, fetchPlayers, fetchInvitations } from "$lib/api";
    import { authClient } from "$lib/auth-client";
    import { goto } from "$app/navigation";
    import GeneralTab from "$lib/components/game/gm/settings/GeneralTab.svelte";
    import PlayersTab from "$lib/components/game/gm/settings/PlayersTab.svelte";
    import CharactersTab from "$lib/components/game/gm/settings/CharactersTab.svelte";
    import InvitationsTab from "$lib/components/game/gm/settings/InvitationsTab.svelte";
    import BestiaryTab from "$lib/components/game/gm/settings/BestiaryTab.svelte";
    import BoardsTab from "$lib/components/game/gm/settings/BoardsTab.svelte";
    import { Settings, Users, Mail, User, Skull, Layers } from "lucide-svelte";

    let activeTab = $state("general");
    let game = $state<any>(null);
    let loading = $state(true);
    let players = $state<any[]>([]);
    let invitations = $state<any[]>([]);

    const tabs = [
        { id: "general", label: "Général", icon: Settings },
        { id: "boards", label: "Plateaux", icon: Layers },
        { id: "players", label: "Joueurs", icon: Users },
        { id: "characters", label: "Personnages", icon: User },
        { id: "bestiary", label: "Bestiaire", icon: Skull },
        { id: "invitations", label: "Invitations", icon: Mail },
    ];

    async function loadGame(id: string) {
        try {
            game = await fetchGame(id);
        } catch (error) {
            console.error("Failed to fetch game:", error);
        }
    }

    async function loadPlayers(id: string) {
        try {
            players = await fetchPlayers(id);
        } catch (error) {
            console.error("Failed to fetch players:", error);
        }
    }

    async function loadInvitations(id: string) {
        try {
            invitations = await fetchInvitations(id);
        } catch (error) {
            console.error("Failed to fetch invitations:", error);
        }
    }

    async function refreshData() {
        const gameId = page.params.id;
        if (gameId) {
            await Promise.all([
                loadPlayers(gameId),
                loadInvitations(gameId),
            ]);
        }
    }

    onMount(async () => {
        const urlParams = new URLSearchParams(window.location.search);
        const tab = urlParams.get("tab");
        if (tab && tabs.find((t) => t.id === tab)) {
            activeTab = tab;
        }

        const gameId = page.params.id;
        if (gameId) {
            try {
                await Promise.all([
                    loadGame(gameId),
                    loadPlayers(gameId),
                    loadInvitations(gameId),
                ]);
            } catch (error) {
                console.error("Failed to fetch game data:", error);
            } finally {
                loading = false;
            }
        }
    });

    function setTab(tabId: string) {
        activeTab = tabId;
        const url = new URL(window.location.href);
        url.searchParams.set("tab", tabId);
        goto(url.toString(), { replaceState: true, keepFocus: true });
    }
</script>

<div class="min-h-screen bg-cream">
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

                <!-- BOARDS TAB -->
                {#if activeTab === "boards"}
                    <BoardsTab mode="settings" />
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
