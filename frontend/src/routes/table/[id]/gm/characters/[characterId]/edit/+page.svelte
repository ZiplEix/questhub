<script lang="ts">
    import { page } from "$app/state";
    import CharacterForm from "$lib/components/game/gm/CharacterForm.svelte";
    import Header from "$lib/components/Header.svelte";
    import { onMount } from "svelte";
    import { api } from "$lib/api";
    import { authClient } from "$lib/auth-client";

    const gameId = page.params.id as string;
    const characterId = page.params.characterId as string;

    let character = $state<any>(null);
    let loading = $state(true);
    let error = $state<string | null>(null);

    onMount(async () => {
        try {
            const { data: tokenData } = await authClient.token();
            if (tokenData?.token) {
                const res = await api.get(
                    `/table/${gameId}/characters/${characterId}`,
                    {
                        headers: {
                            Authorization: `Bearer ${tokenData.token}`,
                        },
                    },
                );
                character = res.data;
            }
        } catch (e) {
            console.error("Failed to fetch character:", e);
            error = "Impossible de charger le personnage.";
        } finally {
            loading = false;
        }
    });
</script>

<div class="min-h-screen bg-stone-50">
    <Header />

    <main class="max-w-7xl mx-auto p-6 md:p-12">
        <div class="mb-8">
            <h1 class="text-3xl font-display font-bold text-dark-gray mb-2">
                Modification de personnage
            </h1>
            <p class="text-stone-500">
                Modifiez les informations de votre personnage.
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
                class="p-4 bg-red-50 text-red-600 rounded-xl border border-red-100"
            >
                {error}
            </div>
        {:else if character}
            <CharacterForm {gameId} {character} />
        {/if}
    </main>
</div>
