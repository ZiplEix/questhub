<script lang="ts">
    import { page } from "$app/stores";
    import { onMount } from "svelte";
    import { api } from "$lib/api";
    import { authClient } from "$lib/auth-client";
    import { goto } from "$app/navigation";
    import Header from "$lib/components/Header.svelte";

    let error = $state("");
    let loading = $state(true);

    onMount(async () => {
        const { data, error: tokenError } = await authClient.token();

        if (tokenError || !data?.token) {
            // Redirect to login with return URL if not authenticated
            goto(`/login?redirectTo=/invitation/${$page.params.code}`);
            return;
        }

        joinGame(data.token);
    });

    async function joinGame(token: string) {
        try {
            const response = await api.post(
                "/table/join",
                { invite_code: $page.params.code },
                {
                    headers: {
                        Authorization: `Bearer ${token}`,
                    },
                },
            );
            goto(`/table/${response.data.id}`);
        } catch (e: any) {
            console.error(e);
            if (e.response && e.response.status === 409 && e.response.data.id) {
                // Already joined, redirect to game
                goto(`/table/${e.response.data.id}`);
                return;
            }
            error =
                "Impossible de rejoindre la partie. Code invalide ou erreur serveur.";
            loading = false;
        }
    }
</script>

<div class="min-h-screen bg-cream">
    <Header />

    <main class="max-w-md mx-auto px-4 py-24 text-center">
        {#if loading}
            <div
                class="bg-white rounded-3xl shadow-sm border border-stone-100 p-8"
            >
                <h1 class="font-display font-bold text-2xl text-dark-gray mb-4">
                    Rejoindre la partie...
                </h1>
                <p class="text-dark-gray/60">
                    Veuillez patienter pendant que nous validons votre
                    invitation.
                </p>
            </div>
        {:else if error}
            <div
                class="bg-white rounded-3xl shadow-sm border border-stone-100 p-8"
            >
                <h1 class="font-display font-bold text-2xl text-red-500 mb-4">
                    Erreur
                </h1>
                <p class="text-dark-gray/60 mb-6">{error}</p>
                <a
                    href="/dashboard"
                    class="bg-burnt-orange text-white px-6 py-3 rounded-xl font-bold hover:bg-opacity-90 transition-all inline-block"
                >
                    Retour au tableau de bord
                </a>
            </div>
        {/if}
    </main>
</div>
