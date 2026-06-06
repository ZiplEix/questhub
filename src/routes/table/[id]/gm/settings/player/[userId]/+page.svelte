<script lang="ts">
    import { page } from "$app/state";
    import { goto } from "$app/navigation";
    import Header from "$lib/components/Header.svelte";
    import { onMount } from "svelte";
    import { fetchPlayers, updatePlayerSettings, fetchGame, updateGame } from "$lib/api";
    import { authClient } from "$lib/auth-client";
    import { ArrowLeft, Check } from "lucide-svelte";

    const gameId = page.params.id as string;
    const userId = page.params.userId as string;

    let player = $state<any>(null);
    let loading = $state(true);
    let saving = $state(false);
    let error = $state<string | null>(null);
    let selectedColor = $state("#E07A5F");

    const colorPalette = [
        { name: "Orange", hex: "#E07A5F" },
        { name: "Rouge", hex: "#E63946" },
        { name: "Crimson", hex: "#9B2226" },
        { name: "Bleu", hex: "#457B9D" },
        { name: "Teal", hex: "#2A9D8F" },
        { name: "Vert", hex: "#76C893" },
        { name: "Olive", hex: "#40916C" },
        { name: "Sable", hex: "#F4A261" },
        { name: "Violet", hex: "#8338EC" },
        { name: "Rose", hex: "#FF006E" },
        { name: "Ardoise", hex: "#2B2D42" },
    ];

    onMount(async () => {
        try {
            const { data: tokenData } = await authClient.token();
            if (tokenData?.token) {
                const gameData = await fetchGame(gameId);
                if (userId === gameData.gm_id) {
                    player = {
                        user_id: gameData.gm_id,
                        name: "Maître du Jeu (MJ)",
                        email: "Créateur de la table",
                        ping_color: gameData.gm_ping_color || "#E07A5F"
                    };
                    selectedColor = player.ping_color;
                } else {
                    const players = await fetchPlayers(gameId);
                    const found = players.find(p => p.user_id === userId);
                    if (found) {
                        player = found;
                        selectedColor = found.ping_color || "#E07A5F";
                    } else {
                        error = "Joueur introuvable dans cette partie.";
                    }
                }
            }
        } catch (e) {
            console.error("Failed to fetch player settings:", e);
            error = "Impossible de charger les paramètres du joueur.";
        } finally {
            loading = false;
        }
    });

    async function handleSave() {
        if (!player) return;
        saving = true;
        try {
            const gameData = await fetchGame(gameId);
            if (userId === gameData.gm_id) {
                await updateGame(gameId, { gm_ping_color: selectedColor });
            } else {
                await updatePlayerSettings(gameId, userId, { ping_color: selectedColor });
            }
            goto(`/table/${gameId}/gm/settings?tab=players`);
        } catch (e) {
            console.error("Failed to save player settings:", e);
            alert("Erreur lors de la sauvegarde des paramètres.");
        } finally {
            saving = false;
        }
    }
</script>

<div class="min-h-screen bg-cream">
    <Header />

    <main class="max-w-2xl mx-auto p-6 md:p-12 animate-in fade-in duration-300">
        <!-- Back Button -->
        <a
            href="/table/{gameId}/gm/settings?tab=players"
            class="inline-flex items-center gap-2 text-stone-500 hover:text-dark-gray mb-6 font-medium transition-colors"
        >
            <ArrowLeft size={16} />
            Retour aux paramètres
        </a>

        <!-- Header -->
        <div class="mb-8">
            <h1 class="text-3xl font-display font-bold text-dark-gray mb-2">
                Paramètres du joueur
            </h1>
            <p class="text-stone-500">
                Personnalisez la configuration de ce joueur pour la table en cours.
            </p>
        </div>

        {#if loading}
            <div class="flex justify-center py-12">
                <div
                    class="animate-spin rounded-full h-8 w-8 border-b-2 border-burnt-orange"
                ></div>
            </div>
        {:else if error}
            <div
                class="p-4 bg-red-50 text-red-600 rounded-xl border border-red-100 shadow-sm"
            >
                {error}
            </div>
        {:else if player}
            <div class="bg-white rounded-2xl shadow-sm border border-stone-100 p-6 md:p-8 space-y-6">
                <!-- Info Section -->
                <div class="flex items-center gap-4 pb-6 border-b border-stone-100">
                    <div
                        class="w-12 h-12 rounded-full bg-burnt-orange/20 flex items-center justify-center text-burnt-orange font-bold text-lg"
                    >
                        {player.name.charAt(0)}
                    </div>
                    <div>
                        <h2 class="text-xl font-bold text-dark-gray">{player.name}</h2>
                        <p class="text-sm text-stone-400">{player.email}</p>
                    </div>
                </div>

                <!-- Settings Form -->
                <div class="space-y-4">
                    <label class="block text-sm font-bold text-dark-gray" for="ping-color-picker">
                        Couleur du Ping
                    </label>
                    <p class="text-xs text-stone-500">
                        Cette couleur sera utilisée pour les pings lumineux de ce joueur sur la carte de jeu.
                    </p>

                    <!-- Palettes Swatches -->
                    <div class="grid grid-cols-6 sm:grid-cols-8 gap-3 pt-2">
                        {#each colorPalette as col}
                            <button
                                onclick={() => selectedColor = col.hex}
                                class="w-10 h-10 rounded-xl transition-all relative border border-black/5 hover:scale-105 active:scale-95 shadow-sm flex items-center justify-center"
                                style="background-color: {col.hex}; cursor: pointer;"
                                title={col.name}
                                type="button"
                            >
                                {#if selectedColor === col.hex}
                                    <Check size={16} class="text-white drop-shadow-md" />
                                {/if}
                            </button>
                        {/each}
                        
                        <!-- Custom Color Picker Swatch -->
                        <div class="relative w-10 h-10 rounded-xl border border-stone-200 shadow-sm hover:scale-105 transition-all overflow-hidden flex items-center justify-center bg-stone-50">
                            <input
                                id="ping-color-picker"
                                type="color"
                                bind:value={selectedColor}
                                class="absolute inset-0 w-full h-full opacity-0 cursor-pointer"
                            />
                            <div 
                                class="w-6 h-6 rounded-lg border border-black/5" 
                                style="background-color: {selectedColor};"
                            ></div>
                        </div>
                    </div>
                </div>

                <!-- Form Actions -->
                <div class="pt-6 border-t border-stone-100 flex justify-end gap-3">
                    <a
                        href="/table/{gameId}/gm/settings?tab=players"
                        class="px-5 py-2.5 bg-white border border-stone-200 hover:bg-stone-50 text-stone-600 rounded-xl font-medium transition-all text-sm"
                    >
                        Annuler
                    </a>
                    <button
                        onclick={handleSave}
                        disabled={saving}
                        class="px-5 py-2.5 bg-burnt-orange text-white rounded-xl font-medium hover:bg-opacity-95 disabled:bg-opacity-50 transition-all text-sm shadow-sm flex items-center justify-center min-w-[100px]"
                    >
                        {#if saving}
                            <div class="animate-spin rounded-full h-4 w-4 border-b-2 border-white"></div>
                        {:else}
                            Sauvegarder
                        {/if}
                    </button>
                </div>
            </div>
        {/if}
    </main>
</div>
