<script lang="ts">
    import { authClient } from "$lib/auth-client";
    import GoogleButton from "$lib/components/GoogleButton.svelte";
    import { goto } from "$app/navigation";

    let name = $state("");
    let email = $state("");
    let password = $state("");
    let confirmPassword = $state("");
    let acceptedCgu = $state(false);
    let error = $state("");
    let loading = $state(false);

    const handleSubmit = async (event: SubmitEvent) => {
        event.preventDefault();
        error = "";

        if (password !== confirmPassword) {
            error = "Les mots de passe ne correspondent pas.";
            return;
        }

        if (!acceptedCgu) {
            error = "Vous devez avoir lu et accepté les CGU.";
            return;
        }

        loading = true;

        try {
            await authClient.signUp.email(
                {
                    email,
                    password,
                    name,
                    callbackURL: "/dashboard", // ou '/login' si tu préfères
                },
                {
                    onError: (ctx) => {
                        error =
                            ctx.error?.message ??
                            "Une erreur est survenue lors de l'inscription.";
                    },
                    onResponse: () => {
                        loading = false;
                    },
                    onSuccess: (ctx) => {
                        // If no session is returned, email confirmation is active, redirect to notification page
                        if (!ctx.data?.session) {
                            goto(`/signup/check-email?email=${encodeURIComponent(email)}`);
                        }
                    },
                },
            );
        } catch (e: any) {
            console.error(e);
            error = "Erreur inattendue";
            loading = false;
        }
    };
</script>

<svelte:head>
    <title>Inscription — QuestHub</title>
    <meta name="description" content="Créez votre compte QuestHub pour commencer à organiser vos campagnes de jeu de rôle et jouer en ligne." />
</svelte:head>

<div class="min-h-screen flex items-center justify-center bg-cream px-4 py-12">
    <div class="w-full max-w-md">
        <div
            class="bg-white shadow-lg rounded-3xl p-8 md:p-10 border border-stone-100"
        >
            <div class="text-center mb-8">
                <h1 class="text-3xl font-display font-bold text-dark-gray mb-2">
                    S'inscrire
                </h1>
                <p class="text-dark-gray/60">Créez un compte pour commencer.</p>
            </div>

            <form class="space-y-5" onsubmit={handleSubmit}>
                <div>
                    <label
                        class="block text-sm font-medium text-dark-gray mb-2 ml-1"
                        for="name">Nom</label
                    >
                    <input
                        type="text"
                        bind:value={name}
                        required
                        class="w-full px-4 py-3 rounded-xl border border-stone-200 focus:outline-none focus:ring-2 focus:ring-burnt-orange/50 focus:border-burnt-orange transition-all bg-stone-50"
                        placeholder="Votre nom complet"
                    />
                </div>

                <div>
                    <label
                        class="block text-sm font-medium text-dark-gray mb-2 ml-1"
                        for="email">Email</label
                    >
                    <input
                        type="email"
                        bind:value={email}
                        required
                        class="w-full px-4 py-3 rounded-xl border border-stone-200 focus:outline-none focus:ring-2 focus:ring-burnt-orange/50 focus:border-burnt-orange transition-all bg-stone-50"
                        placeholder="ex: vous@exemple.com"
                    />
                </div>

                <div>
                    <label
                        class="block text-sm font-medium text-dark-gray mb-2 ml-1"
                        for="password">Mot de passe</label
                    >
                    <input
                        type="password"
                        bind:value={password}
                        required
                        class="w-full px-4 py-3 rounded-xl border border-stone-200 focus:outline-none focus:ring-2 focus:ring-burnt-orange/50 focus:border-burnt-orange transition-all bg-stone-50"
                        placeholder="Créer un mot de passe"
                    />
                </div>

                <div>
                    <label
                        class="block text-sm font-medium text-dark-gray mb-2 ml-1"
                        for="confirmPassword">Confirmer le mot de passe</label
                    >
                    <input
                        type="password"
                        bind:value={confirmPassword}
                        required
                        class="w-full px-4 py-3 rounded-xl border border-stone-200 focus:outline-none focus:ring-2 focus:ring-burnt-orange/50 focus:border-burnt-orange transition-all bg-stone-50"
                        placeholder="Confirmez le mot de passe"
                    />
                </div>

                <div class="flex items-start gap-3 rounded-xl border border-stone-200 bg-stone-50 px-3 py-2.5">
                    <input
                        id="accept-cgu"
                        type="checkbox"
                        bind:checked={acceptedCgu}
                        required
                        class="mt-0.5 h-4 w-4 rounded border-stone-300 text-burnt-orange focus:ring-burnt-orange"
                    />
                    <label for="accept-cgu" class="text-sm text-dark-gray/80 leading-relaxed">
                        J'ai lu et j'accepte les
                        <a href="/conditions-utilisation" class="text-burnt-orange font-medium hover:underline">conditions générales d'utilisation (CGU)</a>.
                    </label>
                </div>

                {#if error}
                    <div
                        class="bg-red-50 text-red-600 px-4 py-3 rounded-xl text-sm border border-red-100"
                    >
                        {error}
                    </div>
                {/if}

                <button
                    type="submit"
                    class="w-full py-3.5 rounded-xl bg-burnt-orange text-white font-bold shadow-md hover:bg-opacity-90 transition-all hover:-translate-y-0.5 disabled:opacity-60 disabled:hover:translate-y-0 hover:cursor-pointer"
                    disabled={loading}
                >
                    {#if loading}
                        Inscription...
                    {:else}
                        S'inscrire
                    {/if}
                </button>

                <div class="relative py-2">
                    <div class="absolute inset-0 flex items-center">
                        <div class="w-full border-t border-stone-200"></div>
                    </div>
                    <div class="relative flex justify-center text-sm">
                        <span class="px-4 bg-white text-dark-gray/50">ou</span>
                    </div>
                </div>

                <GoogleButton
                    callbackURL="/dashboard"
                    errorCallbackURL="/signup"
                />

                <div class="text-center text-sm text-dark-gray/70 mt-6">
                    Déjà un compte ?
                    <a
                        href="/login"
                        class="text-burnt-orange font-medium hover:underline"
                        >Se connecter</a
                    >
                </div>
            </form>
        </div>
    </div>
</div>
