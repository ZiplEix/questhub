<script lang="ts">
    import Header from "$lib/components/Header.svelte";
    import { authClient } from "$lib/auth-client";
    import { submitSupportTicket } from "$lib/api/feedback";
    import { onMount } from "svelte";
    import { 
        Mail, 
        Bug, 
        AlertTriangle, 
        UploadCloud, 
        X, 
        CheckCircle2, 
        Loader, 
        ArrowRight, 
        Info 
    } from "lucide-svelte";

    // Form fields
    let category = $state<'CONTACT' | 'BUG' | 'RECLAMATION'>('CONTACT');
    let email = $state("");
    let subject = $state("");
    let message = $state("");
    let fileInput = $state<HTMLInputElement | null>(null);
    let selectedFile = $state<File | null>(null);
    let previewUrl = $state<string | null>(null);

    // UX states
    let loading = $state(false);
    let success = $state(false);
    let errorMsg = $state<string | null>(null);
    let emailWarning = $state(false);

    // Resolve user session on mount
    onMount(async () => {
        try {
            const { data } = await authClient.getSession();
            if (data?.user) {
                email = data.user.email || "";
            }
        } catch (e) {
            console.error("Failed to load user session for contact form:", e);
        }
    });

    // Handle file picker
    function handleFileChange(event: Event) {
        const target = event.target as HTMLInputElement;
        if (target.files && target.files.length > 0) {
            const file = target.files[0];
            
            // Limit to image attachments (screenshots)
            if (!file.type.startsWith("image/")) {
                alert("Veuillez sélectionner un fichier image (PNG, JPG, WEBP, etc.)");
                return;
            }
            
            // Size limit: 5MB
            if (file.size > 5 * 1024 * 1024) {
                alert("La taille du fichier ne doit pas dépasser 5 Mo.");
                return;
            }

            selectedFile = file;
            previewUrl = URL.createObjectURL(file);
        }
    }

    // Remove file attachment
    function removeAttachment() {
        selectedFile = null;
        if (previewUrl) {
            URL.revokeObjectURL(previewUrl);
            previewUrl = null;
        }
        if (fileInput) {
            fileInput.value = "";
        }
    }

    // Submit form
    async function handleSubmit(event: Event) {
        event.preventDefault();
        errorMsg = null;
        emailWarning = false;

        // Basic validation
        if (!email.trim() || !email.includes("@")) {
            emailWarning = true;
            return;
        }
        if (!subject.trim()) {
            errorMsg = "Veuillez indiquer un sujet.";
            return;
        }
        if (!message.trim()) {
            errorMsg = "Veuillez rédiger votre message.";
            return;
        }

        try {
            loading = true;
            await submitSupportTicket({
                email: email.trim(),
                category,
                subject: subject.trim(),
                message: message.trim(),
                attachmentFile: selectedFile || undefined
            });
            success = true;
            // Clean up resources
            if (previewUrl) URL.revokeObjectURL(previewUrl);
        } catch (err: any) {
            console.error("Failed to submit feedback form:", err);
            errorMsg = err.message || "Une erreur est survenue lors de l'envoi. Veuillez réessayer.";
        } finally {
            loading = false;
        }
    }

    // Reset form for a new ticket
    function resetForm() {
        success = false;
        subject = "";
        message = "";
        removeAttachment();
    }
</script>

<svelte:head>
    <title>Contact & Support | QuestHub</title>
    <meta name="description" content="Formulaire de contact de QuestHub. Signalez des bugs ou soumettez des réclamations." />
</svelte:head>

<Header />

<main class="grow w-full max-w-7xl mx-auto px-4 md:px-8 py-16 flex flex-col justify-center">
    {#if success}
        <!-- Success State View -->
        <div class="bg-white rounded-2xl border border-stone-200/80 shadow-lg p-10 text-center space-y-6 max-w-xl mx-auto animate-in fade-in zoom-in-95 duration-300 w-full">
            <div class="w-16 h-16 bg-green-50 rounded-2xl border border-green-100 flex items-center justify-center text-green-500 mx-auto shadow-xs">
                <CheckCircle2 size={36} />
            </div>
            
            <div class="space-y-2">
                <h1 class="text-2xl font-display font-black text-dark-gray">Message Envoyé !</h1>
                <p class="text-sm text-stone-500 max-w-sm mx-auto leading-relaxed">
                    Merci pour votre retour. Notre équipe a bien reçu votre demande et vous répondra à l'adresse <span class="font-semibold text-stone-700">{email}</span> dans les plus brefs délais.
                </p>
            </div>

            <div class="pt-4 border-t border-stone-100 flex flex-col sm:flex-row gap-3 justify-center">
                <button
                    onclick={resetForm}
                    class="px-5 py-2.5 bg-stone-100 hover:bg-stone-200 text-stone-700 font-bold rounded-xl text-sm transition-colors cursor-pointer"
                >
                    Envoyer un autre message
                </button>
                <a
                    href="/"
                    class="px-5 py-2.5 bg-burnt-orange text-white font-bold rounded-xl text-sm hover:bg-opacity-95 shadow-sm hover:shadow-md transition-all flex items-center justify-center gap-1.5"
                >
                    Retour à l'accueil
                    <ArrowRight size={16} />
                </a>
            </div>
        </div>
    {:else}
        <!-- Form View -->
        <div class="mb-10 text-center">
            <h1 class="text-3xl font-display font-black text-burnt-orange tracking-tight flex items-center justify-center gap-3">
                <Mail class="text-burnt-orange" size={28} />
                Support & Contact
            </h1>
            <p class="text-sm text-stone-500 mt-2 font-medium">
                Une question, un bug à signaler, ou une réclamation ? Dites-nous tout.
            </p>
        </div>

        <div class="bg-white rounded-2xl border border-stone-200/80 shadow-xs p-6 md:p-8 max-w-2xl mx-auto w-full">
            <form onsubmit={handleSubmit} class="space-y-6">
                
                <!-- Category Select Cards -->
                <div class="space-y-2">
                    <label class="block text-xs font-bold text-stone-400 uppercase tracking-wider">
                        Quel est l'objet de votre demande ?
                    </label>
                    <div class="grid grid-cols-1 md:grid-cols-3 gap-3">
                        <!-- General Question Card -->
                        <button
                            type="button"
                            onclick={() => (category = 'CONTACT')}
                            disabled={loading}
                            class="p-4 rounded-xl border-2 text-left transition-all cursor-pointer flex items-start gap-3 flex-row md:flex-col justify-start select-none
                            {category === 'CONTACT' 
                                ? 'border-burnt-orange bg-burnt-orange/5 shadow-xs' 
                                : 'border-stone-200 bg-stone-50/50 hover:bg-stone-50'}"
                        >
                            <div class="p-2 rounded-lg bg-white border border-stone-150 text-burnt-orange shadow-2xs shrink-0
                            {category === 'CONTACT' ? 'border-burnt-orange/30 text-burnt-orange' : 'text-stone-400'}">
                                <Mail size={18} />
                            </div>
                            <div class="min-w-0">
                                <h3 class="font-bold text-sm text-dark-gray leading-tight">Question générale</h3>
                                <p class="text-[11px] text-stone-400 mt-1 leading-snug">Renseignements, suggestions, partenariat...</p>
                            </div>
                        </button>

                        <!-- Bug Report Card -->
                        <button
                            type="button"
                            onclick={() => (category = 'BUG')}
                            disabled={loading}
                            class="p-4 rounded-xl border-2 text-left transition-all cursor-pointer flex items-start gap-3 flex-row md:flex-col justify-start select-none
                            {category === 'BUG' 
                                ? 'border-burnt-orange bg-burnt-orange/5 shadow-xs' 
                                : 'border-stone-200 bg-stone-50/50 hover:bg-stone-50'}"
                        >
                            <div class="p-2 rounded-lg bg-white border border-stone-150 text-burnt-orange shadow-2xs shrink-0
                            {category === 'BUG' ? 'border-burnt-orange/30 text-burnt-orange' : 'text-stone-400'}">
                                <Bug size={18} />
                            </div>
                            <div class="min-w-0">
                                <h3 class="font-bold text-sm text-dark-gray leading-tight">Rapport de bug</h3>
                                <p class="text-[11px] text-stone-400 mt-1 leading-snug">Problème d'affichage, plateau bloqué, crash...</p>
                            </div>
                        </button>

                        <!-- Reclamation Card -->
                        <button
                            type="button"
                            onclick={() => (category = 'RECLAMATION')}
                            disabled={loading}
                            class="p-4 rounded-xl border-2 text-left transition-all cursor-pointer flex items-start gap-3 flex-row md:flex-col justify-start select-none
                            {category === 'RECLAMATION' 
                                ? 'border-burnt-orange bg-burnt-orange/5 shadow-xs' 
                                : 'border-stone-200 bg-stone-50/50 hover:bg-stone-50'}"
                        >
                            <div class="p-2 rounded-lg bg-white border border-stone-150 text-burnt-orange shadow-2xs shrink-0
                            {category === 'RECLAMATION' ? 'border-burnt-orange/30 text-burnt-orange' : 'text-stone-400'}">
                                <AlertTriangle size={18} />
                            </div>
                            <div class="min-w-0">
                                <h3 class="font-bold text-sm text-dark-gray leading-tight">Signalement / Abus</h3>
                                <p class="text-[11px] text-stone-400 mt-1 leading-snug">Réclamation, contenu inapproprié ou abus.</p>
                            </div>
                        </button>
                    </div>
                </div>

                <!-- Input Email -->
                <div class="space-y-1.5">
                    <label for="contact-email" class="block text-xs font-bold text-stone-400 uppercase tracking-wider">
                        Votre adresse e-mail
                    </label>
                    <input
                        type="email"
                        id="contact-email"
                        bind:value={email}
                        disabled={loading}
                        placeholder="nom@exemple.com"
                        class="w-full px-4 py-2.5 bg-stone-50 border rounded-xl text-sm focus:outline-none focus:border-burnt-orange focus:ring-2 focus:ring-burnt-orange/20 transition-all font-medium text-dark-gray
                        {emailWarning ? 'border-red-400 ring-2 ring-red-100' : 'border-stone-200'}"
                    />
                    {#if emailWarning}
                        <p class="text-xs text-red-500 font-semibold flex items-center gap-1">
                            ⚠️ Veuillez indiquer un e-mail valide.
                        </p>
                    {/if}
                </div>

                <!-- Input Subject -->
                <div class="space-y-1.5">
                    <label for="contact-subject" class="block text-xs font-bold text-stone-400 uppercase tracking-wider">
                        Sujet de votre message
                    </label>
                    <input
                        type="text"
                        id="contact-subject"
                        bind:value={subject}
                        disabled={loading}
                        placeholder="Ex: Impossible d'importer mon personnage"
                        class="w-full px-4 py-2.5 bg-stone-50 border border-stone-200 rounded-xl text-sm focus:outline-none focus:border-burnt-orange focus:ring-2 focus:ring-burnt-orange/20 transition-all font-medium text-dark-gray"
                    />
                </div>

                <!-- Input Description / Message -->
                <div class="space-y-1.5">
                    <label for="contact-message" class="block text-xs font-bold text-stone-400 uppercase tracking-wider">
                        Description détaillée
                    </label>
                    <textarea
                        id="contact-message"
                        bind:value={message}
                        disabled={loading}
                        rows="5"
                        placeholder="Racontez-nous ce qui se passe. Plus vous donnez de détails, plus vite nous pourrons vous aider."
                        class="w-full px-4 py-2.5 bg-stone-50 border border-stone-200 rounded-xl text-sm focus:outline-none focus:border-burnt-orange focus:ring-2 focus:ring-burnt-orange/20 transition-all font-medium text-dark-gray resize-none"
                    ></textarea>
                </div>

                <!-- File Attachment Dropzone -->
                <div class="space-y-2">
                    <span class="block text-xs font-bold text-stone-400 uppercase tracking-wider">
                        Capture d'écran / Image (Optionnel)
                    </span>
                    
                    {#if previewUrl}
                        <!-- Preview Thumbnail card -->
                        <div class="relative w-full max-w-sm rounded-xl border border-stone-200 bg-stone-50 p-3 flex items-center justify-between gap-3 animate-in fade-in duration-200">
                            <div class="flex items-center gap-3 min-w-0">
                                <img
                                    src={previewUrl}
                                    alt="Screenshot preview"
                                    class="w-16 h-16 rounded-lg object-cover border border-stone-150 bg-white"
                                />
                                <div class="min-w-0 text-left">
                                    <p class="text-xs font-bold text-stone-700 truncate">{selectedFile?.name}</p>
                                    <p class="text-[10px] text-stone-400 font-semibold mt-0.5">
                                        {(selectedFile!.size / (1024 * 1024)).toFixed(2)} Mo
                                    </p>
                                </div>
                            </div>
                            <button
                                type="button"
                                onclick={removeAttachment}
                                disabled={loading}
                                class="p-2 text-stone-400 hover:text-red-500 rounded-lg hover:bg-red-50 transition-colors cursor-pointer"
                                title="Supprimer l'image"
                            >
                                <X size={16} />
                            </button>
                        </div>
                    {:else}
                        <!-- Dropzone input button -->
                        <div
                            onclick={() => fileInput?.click()}
                            class="border-2 border-dashed border-stone-200 hover:border-burnt-orange/40 rounded-xl p-6 flex flex-col items-center justify-center gap-2 cursor-pointer bg-stone-50/20 hover:bg-stone-50/50 transition-all group"
                        >
                            <UploadCloud class="text-stone-300 group-hover:text-burnt-orange transition-colors" size={28} />
                            <p class="text-xs font-bold text-stone-500 group-hover:text-burnt-orange transition-colors">
                                Choisir ou déposer une capture d'écran
                            </p>
                            <p class="text-[10px] text-stone-400">
                                PNG, JPG ou WEBP jusqu'à 5 Mo
                            </p>
                        </div>
                    {/if}

                    <input
                        type="file"
                        bind:this={fileInput}
                        onchange={handleFileChange}
                        accept="image/*"
                        class="hidden"
                        id="contact-file"
                    />
                </div>

                <!-- Info Box -->
                <div class="bg-stone-50 border border-stone-150 p-4 rounded-xl flex gap-3 text-stone-500 text-xs">
                    <Info size={16} class="shrink-0 text-stone-400 mt-0.5" />
                    <p class="leading-relaxed">
                        Si votre compte est connecté, les détails techniques de votre session et votre identifiant utilisateur 
                        QuestHub seront attachés à votre rapport pour faciliter le diagnostic par notre support technique.
                    </p>
                </div>

                <!-- Error Message -->
                {#if errorMsg}
                    <div class="p-3 bg-red-50 border border-red-100 rounded-xl text-red-500 text-xs font-semibold">
                        ❌ {errorMsg}
                    </div>
                {/if}

                <!-- Actions Submit -->
                <div class="pt-4 border-t border-stone-100 flex justify-end gap-3">
                    <a
                        href="/"
                        class="px-5 py-2.5 bg-stone-100 hover:bg-stone-200 text-stone-700 font-bold rounded-xl text-sm transition-colors cursor-pointer flex items-center justify-center"
                    >
                        Annuler
                    </a>
                    <button
                        type="submit"
                        disabled={loading}
                        class="px-6 py-2.5 bg-burnt-orange text-white font-bold rounded-xl text-sm hover:bg-opacity-95 shadow-sm hover:shadow-md transition-all flex items-center gap-2 cursor-pointer disabled:opacity-50 disabled:cursor-not-allowed"
                    >
                        {#if loading}
                            <Loader size={16} class="animate-spin" />
                            Envoi en cours...
                        {:else}
                            Envoyer la demande
                        {/if}
                    </button>
                </div>
            </form>
        </div>
    {/if}
</main>
