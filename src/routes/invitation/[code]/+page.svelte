<script lang="ts">
    import { page } from "$app/stores";
    import { onMount } from "svelte";
    import { joinGame } from "$lib/api";
    import { supabase } from "$lib/supabaseClient";
    import { authClient } from "$lib/auth-client";
    import { goto } from "$app/navigation";
    import Header from "$lib/components/Header.svelte";
    import { Check, Clock } from "lucide-svelte";

    let error = $state("");
    let loading = $state(true);
    let success = $state(false);
    let pending = $state(false);

    onMount(async () => {
        const { data, error: tokenError } = await authClient.token();

        if (tokenError || !data?.token) {
            // Redirect to login with return URL if not authenticated
            goto(`/login?redirectTo=/invitation/${$page.params.code}`);
            return;
        }

        handleJoinGame();
    });

    async function handleJoinGame() {
        try {
            const { id: gameId } = await joinGame($page.params.code || "");

            const { data: { user } } = await supabase.auth.getUser();
            if (!user) throw new Error("Non connecté");

            // Check if GM
            const { data: gameData } = await supabase
                .from('games')
                .select('gm_id')
                .eq('id', gameId)
                .single();

            // Check if already player
            const { data: playerData } = await supabase
                .from('game_players')
                .select('*')
                .eq('game_id', gameId)
                .eq('user_id', user.id)
                .maybeSingle();

            if (gameData?.gm_id === user.id || playerData) {
                // Already in game, redirect to game
                goto(`/table/${gameId}`);
                return;
            }

            // Otherwise, invitation is pending/created
            success = true;
            loading = false;
        } catch (e: any) {
            console.error(e);
            error = "Impossible de rejoindre la partie. Code invalide ou erreur serveur.";
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
                    Traitement de l'invitation...
                </h1>
                <p class="text-dark-gray/60">
                    Veuillez patienter pendant que nous envoyons votre demande.
                </p>
            </div>
        {:else if success}
            <div
                class="bg-white rounded-3xl shadow-sm border border-stone-100 p-8"
            >
                <div
                    class="w-16 h-16 bg-green-100 rounded-full flex items-center justify-center text-green-600 mx-auto mb-6"
                >
                    <Check size={32} />
                </div>
                <h1 class="font-display font-bold text-2xl text-dark-gray mb-4">
                    Demande envoyée !
                </h1>
                <p class="text-dark-gray/60 mb-8">
                    Votre demande pour rejoindre la partie a été envoyée au
                    Maître du Jeu. Vous serez notifié une fois acceptée.
                </p>
                <a
                    href="/dashboard"
                    class="bg-burnt-orange text-white px-6 py-3 rounded-xl font-bold hover:bg-opacity-90 transition-all inline-block"
                >
                    Retour au tableau de bord
                </a>
            </div>
        {:else if pending}
            <div
                class="bg-white rounded-3xl shadow-sm border border-stone-100 p-8"
            >
                <div
                    class="w-16 h-16 bg-orange-100 rounded-full flex items-center justify-center text-orange-600 mx-auto mb-6"
                >
                    <Clock size={32} />
                </div>
                <h1 class="font-display font-bold text-2xl text-dark-gray mb-4">
                    Demande en attente
                </h1>
                <p class="text-dark-gray/60 mb-8">
                    Vous avez déjà envoyé une demande pour rejoindre cette
                    partie. Veuillez attendre la validation du Maître du Jeu.
                </p>
                <a
                    href="/dashboard"
                    class="bg-burnt-orange text-white px-6 py-3 rounded-xl font-bold hover:bg-opacity-90 transition-all inline-block"
                >
                    Retour au tableau de bord
                </a>
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
