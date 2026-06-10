<script lang="ts">
    import { authClient } from "$lib/auth-client";
    import { Mail, ArrowLeft, Loader2, Check } from "lucide-svelte";

    let email = $state("");
    let loading = $state(false);
    let error = $state("");
    let success = $state(false);

    const handleSubmit = async (event: SubmitEvent) => {
        event.preventDefault();
        error = "";
        loading = true;
        success = false;

        try {
            await authClient.resetPassword(email.trim());
            success = true;
        } catch (e: any) {
            console.error(e);
            error = e.message || "Une erreur est survenue lors de l'envoi de la demande.";
        } finally {
            loading = false;
        }
    };
</script>

<svelte:head>
    <title>Mot de passe oublié - QuestHub</title>
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
                <div class="text-center relative py-4">
                    <div class="relative mx-auto w-20 h-20 mb-6">
                        <div class="absolute inset-0 bg-green-500/10 rounded-full blur-lg animate-pulse"></div>
                        <div class="relative flex items-center justify-center w-full h-full bg-green-50 text-green-600 rounded-full border border-green-200 shadow-inner">
                            <Check size={36} />
                        </div>
                    </div>

                    <h1 class="text-2xl font-display font-black text-dark-gray mb-3">
                        Lien envoyé !
                    </h1>
                    
                    <p class="text-dark-gray/70 text-sm leading-relaxed mb-8">
                        Un e-mail contenant un lien de réinitialisation a été envoyé à :<br />
                        <span class="font-bold text-dark-gray">{email}</span>.
                    </p>

                    <div class="bg-stone-50 rounded-2xl p-4 border border-stone-100 text-left mb-8 text-xs text-dark-gray/70 leading-relaxed space-y-2">
                        <p class="font-semibold text-dark-gray">Que faire ensuite ?</p>
                        <ul class="list-disc pl-4 space-y-1">
                            <li>Consultez votre boîte mail et cliquez sur le lien.</li>
                            <li>Vérifiez le dossier de spams si vous ne recevez rien d'ici quelques minutes.</li>
                        </ul>
                    </div>

                    <a
                        href="/login"
                        class="w-full py-3.5 px-4 rounded-xl bg-burnt-orange hover:bg-opacity-95 text-white font-bold shadow-md transition-all flex items-center justify-center gap-2 hover:-translate-y-0.5 cursor-pointer text-sm"
                    >
                        Retour à la connexion
                    </a>
                </div>
            {:else}
                <div class="text-center mb-8 relative">
                    <h1 class="text-3xl font-display font-bold text-dark-gray mb-2">
                        Mot de passe oublié
                    </h1>
                    <p class="text-dark-gray/60 text-sm leading-relaxed">
                        Entrez votre e-mail pour recevoir un lien de réinitialisation de mot de passe.
                    </p>
                </div>

                <form class="space-y-5 relative" onsubmit={handleSubmit}>
                    <div>
                        <label
                            class="block text-sm font-medium text-dark-gray mb-2 ml-1"
                            for="email">Adresse e-mail</label
                        >
                        <div class="relative">
                            <div class="absolute inset-y-0 left-0 pl-3.5 flex items-center pointer-events-none text-stone-400">
                                <Mail size={18} />
                            </div>
                            <input
                                id="email"
                                type="email"
                                bind:value={email}
                                required
                                class="w-full pl-10 pr-4 py-3 rounded-xl border border-stone-200 focus:outline-none focus:ring-2 focus:ring-burnt-orange/50 focus:border-burnt-orange transition-all bg-stone-50 text-sm"
                                placeholder="ex: vous@exemple.com"
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
                            Envoi en cours...
                        {:else}
                            Envoyer le lien
                        {/if}
                    </button>

                    <div class="text-center mt-6">
                        <a
                            href="/login"
                            class="text-sm text-dark-gray/70 hover:text-dark-gray font-medium transition-all inline-flex items-center gap-2"
                        >
                            <ArrowLeft size={16} />
                            Retour à la connexion
                        </a>
                    </div>
                </form>
            {/if}
        </div>
    </div>
</div>
