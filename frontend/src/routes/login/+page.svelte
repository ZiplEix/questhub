<script lang="ts">
    import { authClient } from "$lib/auth-client";
    import GoogleButton from "$lib/components/GoogleButton.svelte";

    let email = $state("");
    let password = $state("");
    let error = $state("");
    let loading = $state(false);

    const handleSubmit = async () => {
        error = "";
        loading = true;

        try {
            await authClient.signIn.email(
                {
                    email,
                    password,
                    callbackURL: "/dashboard",
                },
                {
                    onError: (ctx) => {
                        error =
                            ctx.error?.message ??
                            "Une erreur est survenue lors de la connexion.";
                    },
                    onResponse: () => {
                        loading = false;
                    },
                    onSuccess: () => {
                        // redirect in case of no callbackURL
                    },
                },
            );
        } catch (e: any) {
            console.error(e);
            error = "Erreur inatendu";
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
                    Se connecter
                </h1>
                <p class="text-dark-gray/60">
                    Entrez votre email et mot de passe pour continuer.
                </p>
            </div>

            <form class="space-y-5" on:submit|preventDefault={handleSubmit}>
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
                        placeholder="Votre mot de passe"
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
                    class="w-full py-3.5 rounded-xl bg-burnt-orange text-white font-bold shadow-md hover:bg-opacity-90 transition-all hover:-translate-y-0.5 disabled:opacity-60 disabled:hover:translate-y-0 hover:cursor-pointer"
                    disabled={loading}
                >
                    {#if loading}
                        Connexion...
                    {:else}
                        Se connecter
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
                    errorCallbackURL="/login"
                />

                <div class="text-center text-sm text-dark-gray/70 mt-6">
                    Pas encore de compte ? <a
                        href="/signup"
                        class="text-burnt-orange font-medium hover:underline"
                        >S'inscrire</a
                    >
                </div>
            </form>
        </div>
    </div>
</div>
