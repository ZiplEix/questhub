<script lang="ts">
    import CharacterForm from "$lib/components/game/gm/CharacterForm.svelte";
    import { page } from "$app/state";
    import { api } from "$lib/api";
    import { authClient } from "$lib/auth-client";
    import { onMount } from "svelte";

    import Header from "$lib/components/Header.svelte";

    const gameId = $derived(page.params.id as string);
    const monsterId = $derived(page.params.monsterId as string);

    let monster = $state<any>(null);
    let loading = $state(true);
    let error = $state<string | null>(null);

    async function fetchMonster() {
        loading = true;
        try {
            const { data: tokenData } = await authClient.token();
            if (tokenData?.token) {
                const response = await api.get(
                    `/table/${gameId}/characters/${monsterId}`,
                    {
                        headers: {
                            Authorization: `Bearer ${tokenData.token}`,
                        },
                    },
                );
                monster = response.data;
            }
        } catch (e) {
            console.error("Failed to fetch monster:", e);
            error = "Impossible de charger le monstre.";
        } finally {
            loading = false;
        }
    }

    $effect(() => {
        if (gameId && monsterId) {
            fetchMonster();
        }
    });
</script>

<div class="min-h-screen bg-stone-50">
    <Header />

    <main class="max-w-7xl mx-auto p-6 md:p-12">
        <div class="mb-8">
            <h1 class="text-3xl font-display font-bold text-dark-gray mb-2">
                Modification de créature
            </h1>
            <p class="text-stone-500">
                Modifiez les caractéristiques de votre créature.
            </p>
        </div>

        {#if loading}
            <div class="flex justify-center py-12">
                <div
                    class="w-8 h-8 border-4 border-burnt-orange/30 border-t-burnt-orange rounded-full animate-spin"
                ></div>
            </div>
        {:else if error}
            <div class="text-center py-12 text-red-500 bg-red-50 rounded-xl">
                {error}
            </div>
        {:else if monster}
            <CharacterForm {gameId} character={monster} type="MONSTER" />
        {/if}
    </main>
</div>
