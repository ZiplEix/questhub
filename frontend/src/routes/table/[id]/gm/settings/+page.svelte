<script lang="ts">
    import { page } from "$app/state";
    import {
        Settings,
        Users,
        Mail,
        Save,
        Trash2,
        Plus,
        Copy,
        RefreshCw,
        Check,
        X,
        UserPlus,
    } from "lucide-svelte";
    import Header from "$lib/components/Header.svelte";
    import { onMount } from "svelte";
    import { api } from "$lib/api";
    import { authClient } from "$lib/auth-client";

    let activeTab = $state("general");
    let game = $state<any>(null);
    let loading = $state(true);
    let copied = $state(false);
    let pendingInvitations = $state<any[]>([]);

    const tabs = [
        { id: "general", label: "Général", icon: Settings },
        { id: "players", label: "Joueurs", icon: Users },
        { id: "invitations", label: "Invitations", icon: Mail },
    ];

    // Mocked Data for other tabs
    let players = $state<any[]>([]);

    async function fetchPlayers(gameId: string, token: string) {
        try {
            const response = await api.get(`/table/${gameId}/players`, {
                headers: {
                    Authorization: `Bearer ${token}`,
                },
            });
            players = response.data;
        } catch (error) {
            console.error("Failed to fetch players:", error);
        }
    }

    onMount(async () => {
        const gameId = page.params.id!;
        try {
            const { data: tokenData } = await authClient.token();
            if (tokenData?.token) {
                const response = await api.get(`/table/${gameId}`, {
                    headers: {
                        Authorization: `Bearer ${tokenData.token}`,
                    },
                });
                game = response.data;

                // Fetch pending invitations
                fetchInvitations(gameId, tokenData.token);
                // Fetch players
                fetchPlayers(gameId, tokenData.token);
            }
        } catch (error) {
            console.error("Failed to fetch game data:", error);
        } finally {
            loading = false;
        }
    });

    async function fetchInvitations(gameId: string, token: string) {
        try {
            const response = await api.get(`/table/${gameId}/invitations`, {
                headers: {
                    Authorization: `Bearer ${token}`,
                },
            });
            pendingInvitations = response.data;
        } catch (error) {
            console.error("Failed to fetch invitations:", error);
        }
    }

    function setActiveTab(tabId: string) {
        activeTab = tabId;
    }

    function copyInviteLink() {
        if (!game?.invite_code) return;
        const link = `${window.location.origin}/invitation/${game.invite_code}`;
        navigator.clipboard.writeText(link);
        copied = true;
        setTimeout(() => (copied = false), 2000);
    }

    async function acceptInvitation(userId: string) {
        const gameId = page.params.id;
        try {
            const { data: tokenData } = await authClient.token();
            if (tokenData?.token) {
                await api.post(
                    `/table/${gameId}/invitations/${userId}/accept`,
                    {},
                    {
                        headers: {
                            Authorization: `Bearer ${tokenData.token}`,
                        },
                    },
                );
                // Refresh invitations
                fetchInvitations(gameId!, tokenData.token);
            }
        } catch (error) {
            console.error("Failed to accept invitation:", error);
        }
    }

    async function declineInvitation(userId: string) {
        const gameId = page.params.id;
        try {
            const { data: tokenData } = await authClient.token();
            if (tokenData?.token) {
                await api.post(
                    `/table/${gameId}/invitations/${userId}/decline`,
                    {},
                    {
                        headers: {
                            Authorization: `Bearer ${tokenData.token}`,
                        },
                    },
                );
                // Refresh invitations
                fetchInvitations(gameId!, tokenData.token);
            }
        } catch (error) {
            console.error("Failed to decline invitation:", error);
        }
    }
    async function regenerateInviteCode() {
        if (
            !confirm(
                "Êtes-vous sûr de vouloir régénérer le code d'invitation ? L'ancien lien ne fonctionnera plus.",
            )
        ) {
            return;
        }

        const gameId = page.params.id;
        try {
            const { data: tokenData } = await authClient.token();
            if (tokenData?.token) {
                const response = await api.post(
                    `/table/${gameId}/invite-code`,
                    {},
                    {
                        headers: {
                            Authorization: `Bearer ${tokenData.token}`,
                        },
                    },
                );
                game.invite_code = response.data.invite_code;
            }
        } catch (error) {
            console.error("Failed to regenerate invite code:", error);
        }
    }

    async function removePlayer(userId: string, userName: string) {
        if (
            !confirm(
                `Êtes-vous sûr de vouloir exclure ${userName} de la partie ?`,
            )
        ) {
            return;
        }

        const gameId = page.params.id;
        try {
            const { data: tokenData } = await authClient.token();
            if (tokenData?.token) {
                await api.delete(`/table/${gameId}/players/${userId}`, {
                    headers: {
                        Authorization: `Bearer ${tokenData.token}`,
                    },
                });
                // Refresh players
                fetchPlayers(gameId!, tokenData.token);
            }
        } catch (error) {
            console.error("Failed to remove player:", error);
        }
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
                        onclick={() => setActiveTab(tab.id)}
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
                    <div class="space-y-6 animate-in fade-in duration-300">
                        <div class="space-y-2">
                            <label
                                for="tableName"
                                class="text-sm font-bold text-dark-gray"
                                >Nom de la table</label
                            >
                            <input
                                type="text"
                                id="tableName"
                                bind:value={game.name}
                                class="w-full px-4 py-2 rounded-xl border border-stone-200 focus:outline-none focus:ring-2 focus:ring-burnt-orange/20 focus:border-burnt-orange transition-all"
                            />
                        </div>

                        <div class="space-y-2">
                            <label
                                class="text-sm font-bold text-dark-gray block"
                                for="inviteLink"
                            >
                                Lien d'invitation
                            </label>
                            <div class="flex gap-2" id="inviteLink">
                                <div
                                    class="flex-1 bg-stone-50 px-4 py-2 rounded-xl border border-stone-200 text-stone-600 font-mono text-sm truncate"
                                >
                                    {window.location
                                        .origin}/invitation/{game.invite_code}
                                </div>
                                <button
                                    onclick={copyInviteLink}
                                    class="flex items-center gap-2 px-4 py-2 bg-white border border-stone-200 text-dark-gray rounded-xl font-medium hover:bg-stone-50 hover:border-burnt-orange/30 hover:text-burnt-orange transition-all shadow-sm min-w-[100px] justify-center"
                                >
                                    {#if copied}
                                        <Check size={18} />
                                        <span>Copié</span>
                                    {:else}
                                        <Copy size={18} />
                                        <span>Copier</span>
                                    {/if}
                                </button>
                                <button
                                    onclick={regenerateInviteCode}
                                    class="flex items-center gap-2 px-4 py-2 bg-white border border-stone-200 text-dark-gray rounded-xl font-medium hover:bg-stone-50 hover:border-burnt-orange/30 hover:text-burnt-orange transition-all shadow-sm"
                                    title="Régénérer le lien"
                                >
                                    <RefreshCw size={18} />
                                </button>
                            </div>
                            <p class="text-xs text-stone-500">
                                Partagez ce lien avec vos joueurs pour qu'ils
                                puissent rejoindre la partie.
                            </p>
                        </div>

                        <div
                            class="pt-4 border-t border-stone-100 flex justify-end"
                        >
                            <button
                                class="flex items-center gap-2 px-6 py-2.5 bg-burnt-orange text-white rounded-xl font-medium shadow-md hover:bg-opacity-90 transition-all hover:-translate-y-0.5"
                            >
                                <Save size={18} />
                                Enregistrer
                            </button>
                        </div>
                    </div>
                {/if}

                <!-- PLAYERS TAB -->
                {#if activeTab === "players"}
                    <div class="space-y-6 animate-in fade-in duration-300">
                        <div class="flex justify-between items-center mb-4">
                            <h3 class="text-lg font-bold text-dark-gray">
                                Joueurs ({players.length})
                            </h3>
                            <button
                                class="flex items-center gap-2 px-4 py-2 bg-stone-100 text-dark-gray rounded-xl font-medium hover:bg-stone-200 transition-all text-sm"
                            >
                                <UserPlus size={18} />
                                Inviter
                            </button>
                        </div>

                        <div class="space-y-3">
                            {#each players as player}
                                <div
                                    class="flex items-center justify-between p-4 rounded-xl border border-stone-100 hover:border-stone-200 transition-all"
                                >
                                    <div class="flex items-center gap-4">
                                        {#if player.avatar_url}
                                            <img
                                                src={player.avatar_url}
                                                alt={player.name}
                                                class="w-10 h-10 rounded-full object-cover"
                                            />
                                        {:else}
                                            <div
                                                class="w-10 h-10 rounded-full bg-burnt-orange/20 flex items-center justify-center text-burnt-orange font-bold"
                                            >
                                                {player.name.charAt(0)}
                                            </div>
                                        {/if}
                                        <div>
                                            <div
                                                class="flex items-center gap-2"
                                            >
                                                <p
                                                    class="font-bold text-dark-gray"
                                                >
                                                    {player.name}
                                                </p>
                                                {#if player.is_gm}
                                                    <span
                                                        class="px-2 py-0.5 rounded-full bg-burnt-orange/10 text-burnt-orange text-xs font-bold border border-burnt-orange/20"
                                                    >
                                                        MJ
                                                    </span>
                                                {/if}
                                            </div>
                                            <p class="text-xs text-stone-500">
                                                Rejoint le {new Date(
                                                    player.joined_at,
                                                ).toLocaleDateString()}
                                            </p>
                                        </div>
                                    </div>
                                    <div class="flex items-center gap-2">
                                        {#if !player.is_gm}
                                            <button
                                                onclick={() =>
                                                    removePlayer(
                                                        player.user_id,
                                                        player.name,
                                                    )}
                                                class="p-2 text-stone-400 hover:text-red-500 hover:bg-red-50 rounded-lg transition-all"
                                                title="Exclure"
                                            >
                                                <Trash2 size={18} />
                                            </button>
                                        {/if}
                                    </div>
                                </div>
                            {/each}
                        </div>
                    </div>
                {/if}

                <!-- INVITATIONS TAB -->
                {#if activeTab === "invitations"}
                    <div class="space-y-6 animate-in fade-in duration-300">
                        <div class="flex justify-between items-center mb-4">
                            <h3 class="text-lg font-bold text-dark-gray">
                                Demandes d'adhésion ({pendingInvitations.length})
                            </h3>
                        </div>

                        {#if pendingInvitations.length === 0}
                            <div
                                class="flex flex-col items-center justify-center py-12 text-center"
                            >
                                <div
                                    class="w-16 h-16 bg-stone-100 rounded-full flex items-center justify-center text-stone-400 mb-4"
                                >
                                    <Mail size={32} />
                                </div>
                                <h4 class="text-lg font-bold text-dark-gray">
                                    Aucune demande en attente
                                </h4>
                                <p class="text-stone-500 max-w-sm mt-2">
                                    Partagez le lien d'invitation pour permettre
                                    aux joueurs de rejoindre votre table.
                                </p>
                            </div>
                        {:else}
                            <div class="space-y-3">
                                {#each pendingInvitations as invite}
                                    <div
                                        class="flex flex-col md:flex-row md:items-center justify-between p-4 rounded-xl border border-stone-100 hover:border-stone-200 transition-all gap-4"
                                    >
                                        <div class="flex items-center gap-4">
                                            <div
                                                class="w-10 h-10 rounded-full bg-burnt-orange/20 flex items-center justify-center text-burnt-orange font-bold"
                                            >
                                                {invite.user_name.charAt(0)}
                                            </div>
                                            <div>
                                                <p
                                                    class="font-bold text-dark-gray"
                                                >
                                                    {invite.user_name}
                                                </p>
                                                <p
                                                    class="text-xs text-stone-500"
                                                >
                                                    Demande envoyée le {new Date(
                                                        invite.created_at,
                                                    ).toLocaleDateString()}
                                                </p>
                                            </div>
                                        </div>
                                        <div class="flex items-center gap-2">
                                            <button
                                                onclick={() =>
                                                    acceptInvitation(
                                                        invite.user_id,
                                                    )}
                                                class="flex items-center gap-2 px-3 py-2 bg-green-100 text-green-700 rounded-lg font-medium hover:bg-green-200 transition-colors text-sm"
                                            >
                                                <Check size={16} />
                                                Accepter
                                            </button>
                                            <button
                                                onclick={() =>
                                                    declineInvitation(
                                                        invite.user_id,
                                                    )}
                                                class="flex items-center gap-2 px-3 py-2 bg-red-100 text-red-700 rounded-lg font-medium hover:bg-red-200 transition-colors text-sm"
                                            >
                                                <X size={16} />
                                                Refuser
                                            </button>
                                        </div>
                                    </div>
                                {/each}
                            </div>
                        {/if}
                    </div>
                {/if}
            </div>
        {/if}
    </main>
</div>
