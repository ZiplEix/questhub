<script lang="ts">
    import { authClient } from "$lib/auth-client";

    let name = $state("");
    let email = $state("");
    let password = $state("");
    let confirmPassword = $state("");
    let error = $state("");
    let loading = $state(false);
    let socialLoading = $state(false);

    const handleGoogleSignIn = async () => {
        error = "";
        socialLoading = true;

        try {
            await authClient.signIn.social({
                provider: "google",
                callbackURL: "/",
                errorCallbackURL: "/signup",
            });
        } catch (e: any) {
            console.error(e);
            error = "Impossible de se connecter avec Google.";
            socialLoading = false;
        }
    };

    const handleSubmit = async () => {
        error = "";

        if (password !== confirmPassword) {
            error = "Les mots de passe ne correspondent pas.";
            return;
        }

        loading = true;

        try {
            await authClient.signUp.email(
                {
                    email,
                    password,
                    name,
                    callbackURL: "/", // ou '/login' si tu préfères
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
                    onSuccess: () => {
                        // si pas de callbackURL, tu peux rediriger ici
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

            <form class="space-y-5" on:submit|preventDefault={handleSubmit}>
                <div>
                    <label
                        class="block text-sm font-medium text-dark-gray mb-2 ml-1"
                        >Nom</label
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
                        >Email</label
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
                        >Mot de passe</label
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
                        >Confirmer le mot de passe</label
                    >
                    <input
                        type="password"
                        bind:value={confirmPassword}
                        required
                        class="w-full px-4 py-3 rounded-xl border border-stone-200 focus:outline-none focus:ring-2 focus:ring-burnt-orange/50 focus:border-burnt-orange transition-all bg-stone-50"
                        placeholder="Confirmez le mot de passe"
                    />
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
                    class="w-full py-3.5 rounded-xl bg-burnt-orange text-white font-bold shadow-md hover:bg-opacity-90 transition-all hover:-translate-y-0.5 disabled:opacity-60 disabled:hover:translate-y-0"
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

                <button
                    type="button"
                    class="w-full py-3.5 rounded-xl border border-stone-200 bg-white text-dark-gray font-medium hover:bg-stone-50 transition-all flex items-center justify-center gap-3"
                    on:click={handleGoogleSignIn}
                    disabled={socialLoading}
                >
                    {#if socialLoading}
                        Connexion avec Google...
                    {:else}
                        <svg class="w-5 h-5" viewBox="0 0 24 24">
                            <path
                                d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z"
                                fill="#4285F4"
                            />
                            <path
                                d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z"
                                fill="#34A853"
                            />
                            <path
                                d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z"
                                fill="#FBBC05"
                            />
                            <path
                                d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z"
                                fill="#EA4335"
                            />
                        </svg>
                        <span>Continuer avec Google</span>
                    {/if}
                </button>

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
