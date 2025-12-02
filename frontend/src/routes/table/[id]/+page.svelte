<script lang="ts">
    import PlayerLayout from "$lib/components/game/player/PlayerLayout.svelte";
    import ImmersionZone from "$lib/components/game/player/ImmersionZone.svelte";
    import PlayerDashboard from "$lib/components/game/player/PlayerDashboard.svelte";
    import Header from "$lib/components/Header.svelte";
    import {
        ChevronRight,
        ChevronLeft,
        ShieldAlert,
        User,
    } from "lucide-svelte";
    import { onMount } from "svelte";
    import { page } from "$app/state";
    import { api } from "$lib/api";
    import { authClient } from "$lib/auth-client";

    let isDashboardOpen = $state(true);
    let loading = $state(true);
    let game = $state<any>(null);
    let error = $state<string | null>(null);

    function toggleDashboard() {
        isDashboardOpen = !isDashboardOpen;
    }

    onMount(async () => {
        const gameId = page.params.id;
        try {
            const { data: tokenData } = await authClient.token();
            if (tokenData?.token) {
                const response = await api.get(`/table/${gameId}`, {
                    headers: {
                        Authorization: `Bearer ${tokenData.token}`,
                    },
                });
                game = response.data;
            }
        } catch (e) {
            console.error(e);
            error = "Impossible de charger la partie.";
        } finally {
            loading = false;
        }
    });
</script>

{#if loading}
    <div class="h-screen w-full flex items-center justify-center bg-stone-50">
        <div
            class="animate-spin rounded-full h-12 w-12 border-b-2 border-burnt-orange"
        ></div>
    </div>
{:else if error}
    <div class="min-h-screen bg-stone-50 flex flex-col">
        <Header />
        <div class="flex-1 flex flex-col items-center justify-center gap-4 p-6">
            <ShieldAlert size={48} class="text-red-500" />
            <p class="text-xl font-bold text-dark-gray">{error}</p>
        </div>
    </div>
{:else if game}
    {#if game.is_gm || game.current_character_id}
        <PlayerLayout>
            <div class="relative w-full h-full flex overflow-hidden">
                <!-- Left Zone: Immersion (Flexible) -->
                <div
                    class="h-full relative transition-all duration-500 ease-in-out"
                    style="width: {isDashboardOpen ? '65%' : '100%'}"
                >
                    <ImmersionZone />

                    <!-- Toggle Button (Attached to the right edge of Immersion Zone) -->
                    <button
                        onclick={toggleDashboard}
                        class="absolute top-1/2 -translate-y-1/2 right-0 z-50 bg-white border border-stone-200 shadow-md p-1 rounded-l-lg hover:bg-stone-50 text-stone-500"
                        aria-label={isDashboardOpen
                            ? "Fermer le tableau de bord"
                            : "Ouvrir le tableau de bord"}
                    >
                        {#if isDashboardOpen}
                            <ChevronRight size={20} />
                        {:else}
                            <ChevronLeft size={20} />
                        {/if}
                    </button>
                </div>

                <!-- Right Zone: Dashboard (Collapsible) -->
                <div
                    class="h-full relative transition-all duration-500 ease-in-out border-l border-stone-200 shadow-2xl z-40 bg-stone-50"
                    style="width: {isDashboardOpen
                        ? '35%'
                        : '0%'}; opacity: {isDashboardOpen
                        ? '1'
                        : '0'}; pointer-events: {isDashboardOpen
                        ? 'auto'
                        : 'none'}"
                >
                    <div class="w-full h-full min-w-[350px]">
                        <!-- Prevent content squashing -->
                        <PlayerDashboard />
                    </div>
                </div>
            </div>
        </PlayerLayout>
    {:else}
        <!-- No Character State -->
        <div class="min-h-screen bg-stone-50 flex flex-col">
            <Header />
            <div class="flex-1 flex flex-col items-center justify-center p-6">
                <div class="max-w-lg w-full text-center space-y-8">
                    <div class="relative inline-block">
                        <div
                            class="w-24 h-24 bg-white rounded-full flex items-center justify-center mx-auto text-burnt-orange shadow-sm border border-stone-100"
                        >
                            <User size={48} />
                        </div>
                        <div
                            class="absolute -bottom-2 -right-2 bg-white p-1.5 rounded-full shadow-sm border border-stone-100 text-stone-400"
                        >
                            <ShieldAlert size={24} />
                        </div>
                    </div>

                    <div class="space-y-3">
                        <h1
                            class="text-3xl font-display font-bold text-dark-gray"
                        >
                            En attente d'un personnage
                        </h1>
                        <p class="text-stone-500 text-lg leading-relaxed">
                            Vous avez rejoint la table <span
                                class="font-bold text-dark-gray"
                                >{game.name}</span
                            >, mais aucun personnage ne vous a encore été
                            assigné.
                        </p>
                    </div>

                    <div
                        class="bg-white p-6 rounded-2xl border border-stone-100 shadow-sm text-left space-y-4"
                    >
                        <h3
                            class="font-bold text-dark-gray flex items-center gap-2"
                        >
                            <span
                                class="w-1.5 h-1.5 rounded-full bg-burnt-orange"
                            ></span>
                            Que faire maintenant ?
                        </h3>
                        <p class="text-stone-500 text-sm">
                            Contactez votre Maître du Jeu (MJ) pour qu'il vous
                            crée un personnage ou vous en assigne un existant.
                            Une fois fait, rafraîchissez cette page pour accéder
                            à l'interface de jeu.
                        </p>
                    </div>

                    <div class="pt-4 flex flex-col gap-3">
                        <a
                            href="/dashboard"
                            class="px-6 py-3 bg-white border border-stone-200 hover:border-burnt-orange/50 hover:text-burnt-orange text-stone-600 rounded-xl font-medium transition-all shadow-sm"
                        >
                            Retour au tableau de bord
                        </a>
                    </div>
                </div>
            </div>
        </div>
    {/if}
{/if}
