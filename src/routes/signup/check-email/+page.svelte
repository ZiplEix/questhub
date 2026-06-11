<script lang="ts">
    import { page } from "$app/stores";
    import { Mail, ArrowLeft, ShieldCheck, HelpCircle, RefreshCw, CheckCircle2 } from "lucide-svelte";
    import { supabase } from "$lib/supabaseClient";

    // Get email from URL parameters or fallback
    let emailParam = $derived($page.url.searchParams.get("email") || "");
    let emailDisplay = $derived(emailParam || "votre adresse e-mail");

    let resending = $state(false);
    let resendSuccess = $state(false);
    let resendError = $state("");
    let cooldown = $state(0);
    let timerId = $state<any>(null);

    function startCooldown() {
        cooldown = 60;
        if (timerId) clearInterval(timerId);
        timerId = setInterval(() => {
            if (cooldown > 0) {
                cooldown -= 1;
            } else {
                clearInterval(timerId);
            }
        }, 1000);
    }

    async function handleResend() {
        if (!emailParam) {
            resendError = "Adresse e-mail manquante dans l'URL. Veuillez retourner à l'inscription.";
            return;
        }
        if (cooldown > 0 || resending) return;

        resending = true;
        resendError = "";
        resendSuccess = false;

        try {
            const { error } = await supabase.auth.resend({
                type: 'signup',
                email: emailParam,
                options: {
                    emailRedirectTo: `${window.location.origin}/auth/callback`
                }
            });

            if (error) {
                resendError = error.message || "Une erreur est survenue lors du renvoi du mail.";
            } else {
                resendSuccess = true;
                startCooldown();
            }
        } catch (e: any) {
            resendError = e.message || "Une erreur inattendue est survenue.";
        } finally {
            resending = false;
        }
    }
</script>

<svelte:head>
    <title>Vérifiez votre e-mail - QuestHub</title>
</svelte:head>

<div class="min-h-screen flex items-center justify-center bg-cream px-4 py-12">
    <div class="w-full max-w-md">
        <div
            class="bg-white shadow-xl rounded-3xl p-8 md:p-10 border border-stone-100 text-center relative overflow-hidden"
        >
            <!-- Background Accent Graphic -->
            <div class="absolute -top-24 -right-24 w-48 h-48 rounded-full bg-mustard-yellow/10 blur-2xl"></div>
            <div class="absolute -bottom-24 -left-24 w-48 h-48 rounded-full bg-burnt-orange/10 blur-2xl"></div>

            <!-- Icon Header -->
            <div class="relative mx-auto w-24 h-24 mb-6">
                <!-- Glowing effect -->
                <div class="absolute inset-0 bg-burnt-orange/15 rounded-full blur-xl animate-pulse"></div>
                <!-- Main Circle -->
                <div class="relative flex items-center justify-center w-full h-full bg-burnt-orange/10 text-burnt-orange rounded-full border border-burnt-orange/20 shadow-inner">
                    <Mail size={42} class="animate-bounce duration-1000" />
                </div>
                <!-- Mini overlay badge -->
                <div class="absolute bottom-1 right-1 bg-green-500 text-white rounded-full p-1 border-2 border-white shadow-sm">
                    <ShieldCheck size={14} />
                </div>
            </div>

            <!-- Text Content -->
            <h1 class="text-3xl font-display font-black text-dark-gray mb-3 leading-tight">
                Vérifiez votre boîte mail
            </h1>
            
            <p class="text-dark-gray/70 text-sm leading-relaxed mb-6">
                Nous avons envoyé un e-mail de confirmation à :<br />
                <span class="font-bold text-dark-gray bg-stone-50 px-2 py-1 rounded-md border border-stone-200/50 mt-1 inline-block select-all">
                    {emailDisplay}
                </span>
            </p>

            <div class="bg-stone-50/80 rounded-2xl p-5 border border-stone-100 text-left space-y-3 mb-8">
                <div class="flex gap-2.5 items-start">
                    <div class="bg-burnt-orange/10 text-burnt-orange rounded-lg p-1.5 shrink-0 mt-0.5">
                        <ShieldCheck size={16} />
                    </div>
                    <p class="text-xs text-dark-gray/80 leading-relaxed">
                        Cliquez sur le lien contenu dans le message pour activer votre compte.
                    </p>
                </div>
                <div class="flex gap-2.5 items-start">
                    <div class="bg-mustard-yellow/20 text-dark-gray/80 rounded-lg p-1.5 shrink-0 mt-0.5">
                        <HelpCircle size={16} />
                    </div>
                    <div class="space-y-1">
                        <p class="text-xs font-bold text-dark-gray">Vous ne trouvez pas l'e-mail ?</p>
                        <p class="text-[11px] text-dark-gray/60 leading-relaxed">
                            Pensez à vérifier votre dossier de courriers indésirables (spams).
                        </p>
                    </div>
                </div>
            </div>

            {#if emailParam}
                <div class="mb-8 text-center">
                    {#if resendSuccess}
                        <div class="mb-4 bg-green-50 text-green-700 border border-green-100 text-xs px-4 py-3 rounded-xl flex items-center justify-center gap-2">
                            <CheckCircle2 size={16} class="shrink-0" />
                            <span>Un nouvel e-mail de confirmation a été envoyé !</span>
                        </div>
                    {/if}

                    {#if resendError}
                        <div class="mb-4 bg-red-50 text-red-600 border border-red-100 text-xs px-4 py-3 rounded-xl">
                            {resendError}
                        </div>
                    {/if}

                    <button
                        onclick={handleResend}
                        disabled={resending || cooldown > 0}
                        class="text-xs text-burnt-orange font-bold hover:underline inline-flex items-center gap-1.5 cursor-pointer disabled:opacity-50 disabled:cursor-not-allowed transition-all"
                    >
                        {#if resending}
                            <RefreshCw size={14} class="animate-spin" />
                            Renvoi en cours...
                        {:else if cooldown > 0}
                            Renvoyer l'e-mail ({cooldown}s)
                        {:else}
                            <RefreshCw size={14} />
                            Renvoyer l'e-mail de confirmation
                        {/if}
                    </button>
                </div>
            {/if}

            <!-- Action buttons -->
            <div class="space-y-3">
                <a
                    href="/login"
                    class="w-full py-3.5 px-4 rounded-xl bg-burnt-orange hover:bg-opacity-95 text-white font-bold shadow-md hover:shadow-lg transition-all flex items-center justify-center gap-2 hover:-translate-y-0.5 cursor-pointer"
                >
                    Aller à la page de connexion
                </a>

                <a
                    href="/signup"
                    class="w-full py-3 px-4 rounded-xl border border-stone-200 hover:bg-stone-50 text-dark-gray/70 hover:text-dark-gray font-semibold transition-all flex items-center justify-center gap-2 cursor-pointer text-sm"
                >
                    <ArrowLeft size={16} />
                    Retour à l'inscription
                </a>
            </div>
        </div>
    </div>
</div>
