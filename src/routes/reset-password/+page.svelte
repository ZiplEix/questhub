<script lang="ts">
    import { authClient } from "$lib/auth-client";
    import { goto } from "$app/navigation";
    import { Lock, Loader2, Check } from "lucide-svelte";

    let password = $state("");
    let confirmPassword = $state("");
    let loading = $state(false);
    let error = $state("");
    let success = $state(false);

    const handleSubmit = async (event: SubmitEvent) => {
        event.preventDefault();
        error = "";

        if (password.length < 6) {
            error = "Le mot de passe doit contenir au moins 6 caractères.";
            return;
        }

        if (password !== confirmPassword) {
            error = "Les mots de passe ne correspondent pas.";
            return;
        }

        loading = true;

        try {
            await authClient.updatePassword(password);
            success = true;
            setTimeout(() => {
                goto("/dashboard");
            }, 1800);
        } catch (e: any) {
            console.error(e);
            error = e.message || "Une erreur est survenue lors de la réinitialisation du mot de passe.";
            loading = false;
        }
    };
</script>

<svelte:head>
    <title>Réinitialiser le mot de passe - QuestHub</title>
</svelte:head>

<div class="min-h-screen flex items-center justify-center bg-cream px-4 py-12">
    <div class="w-full max-w-md">
        <div
            class="bg-white shadow-xl rounded-3xl p-8 md:p-10 border border-stone-100 relative overflow-hidden"
        >
            <!-- Background Accent Graphic -->
            <div class="absolute -top-24 -right-24 w-48 h-48 rounded-full bg-mustard-yellow/10 blur-2xl"></div>
            <div class="absolute -bottom-24 -left-24 w-48 h-48 rounded-full bg-burnt-orange/10 blur-2xl"></div>

            {#if success}
                <div class="text-center relative py-6">
                    <div class="relative mx-auto w-20 h-20 mb-6">
                        <div class="absolute inset-0 bg-green-500/10 rounded-full blur-lg animate-pulse"></div>
                        <div class="relative flex items-center justify-center w-full h-full bg-green-50 text-green-600 rounded-full border border-green-200 shadow-inner">
                            <Check size={36} />
                        </div>
                    </div>

                    <h1 class="text-2xl font-display font-black text-dark-gray mb-3">
                        Mot de passe enregistré !
                    </h1>
                    
                    <p class="text-dark-gray/60 text-sm leading-relaxed">
                        Votre mot de passe a été réinitialisé avec succès. Vous allez être redirigé vers votre tableau de bord...
                    </p>
                </div>
            {:else}
                <div class="text-center mb-8 relative">
                    <h1 class="text-3xl font-display font-bold text-dark-gray mb-2">
                        Nouveau mot de passe
                    </h1>
                    <p class="text-dark-gray/60 text-sm leading-relaxed">
                        Saisissez votre nouveau mot de passe ci-dessous.
                    </p>
                </div>

                <form class="space-y-5 relative" onsubmit={handleSubmit}>
                    <div>
                        <label
                            class="block text-sm font-medium text-dark-gray mb-2 ml-1"
                            for="password">Nouveau mot de passe</label
                        >
                        <div class="relative">
                            <div class="absolute inset-y-0 left-0 pl-3.5 flex items-center pointer-events-none text-stone-400">
                                <Lock size={18} />
                            </div>
                            <input
                                id="password"
                                type="password"
                                bind:value={password}
                                required
                                class="w-full pl-10 pr-4 py-3 rounded-xl border border-stone-200 focus:outline-none focus:ring-2 focus:ring-burnt-orange/50 focus:border-burnt-orange transition-all bg-stone-50 text-sm"
                                placeholder="Min. 6 caractères"
                            />
                        </div>
                    </div>

                    <div>
                        <label
                            class="block text-sm font-medium text-dark-gray mb-2 ml-1"
                            for="confirm-password">Confirmer le mot de passe</label
                        >
                        <div class="relative">
                            <div class="absolute inset-y-0 left-0 pl-3.5 flex items-center pointer-events-none text-stone-400">
                                <Lock size={18} />
                            </div>
                            <input
                                id="confirm-password"
                                type="password"
                                bind:value={confirmPassword}
                                required
                                class="w-full pl-10 pr-4 py-3 rounded-xl border border-stone-200 focus:outline-none focus:ring-2 focus:ring-burnt-orange/50 focus:border-burnt-orange transition-all bg-stone-50 text-sm"
                                placeholder="Confirmez le mot de passe"
                            />
                        </div>
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
                        class="w-full py-3.5 rounded-xl bg-burnt-orange text-white font-bold shadow-md hover:bg-opacity-95 transition-all hover:-translate-y-0.5 disabled:opacity-60 disabled:hover:translate-y-0 hover:cursor-pointer flex items-center justify-center gap-2"
                        disabled={loading}
                    >
                        {#if loading}
                            <Loader2 size={18} class="animate-spin" />
                            Enregistrement...
                        {:else}
                            Confirmer
                        {/if}
                    </button>
                </form>
            {/if}
        </div>
    </div>
</div>
