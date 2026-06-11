<script lang="ts">
    import { onMount } from "svelte";
    import { fetchUserStats, fetchUserCampaigns, fetchUserCharacters, fetchUserMonsters, deleteUserAccount, fetchMediaLibrary, uploadToMediaLibrary, deleteFromMediaLibrary, getUserStorageUsage, STORAGE_LIMIT, type MediaAsset } from "$lib/api";
    import { uploadImage, validateImage } from "$lib/api/storage";
    import { authClient } from "$lib/auth-client";
    import { goto } from "$app/navigation";
    import Header from "$lib/components/Header.svelte";
    import {
        User,
        Activity,
        Swords,
        Skull,
        BookOpen,
        Mail,
        Camera,
        Loader2,
        Check,
        Lock,
        ShieldAlert,
        Trash2,
        Link,
        RefreshCw,
        X
    } from "lucide-svelte";
    import type { SessionUser } from "$lib/types/session-user";

    let activeTab = $state("general");
    let user = $state<SessionUser | null>(null);
    let stats = $state<any>(null);
    let campaigns = $state<any[]>([]);
    let characters = $state<any[]>([]);
    let monsters = $state<any[]>([]);
    let loading = $state(true);
    let avatarUploading = $state(false);
    let avatarSuccess = $state(false);
    let avatarError = $state("");
    let fileInput = $state<HTMLInputElement | null>(null);
 
    let showDeleteConfirm = $state(false);
    let deleteMonsters = $state(true);
    let deleteGames = $state(true);
    let deleteCharacters = $state(true);
    let deleteImages = $state(true);
    let deletingAccount = $state(false);
    let deleteError = $state("");

    let editName = $state("");
    let nameSaving = $state(false);
    let nameError = $state("");
    let nameSuccess = $state(false);

    let editEmail = $state("");
    let emailSaving = $state(false);
    let emailError = $state("");
    let emailSuccess = $state(false);

    let newPassword = $state("");
    let confirmPassword = $state("");
    let passwordSaving = $state(false);
    let passwordError = $state("");
    let passwordSuccess = $state(false);

    let mediaAssets = $state<MediaAsset[]>([]);
    let mediaLoading = $state(false);
    let mediaUploading = $state(false);
    let mediaError = $state("");
    let mediaFileInput = $state<HTMLInputElement | null>(null);
    let storageUsage = $state(0);
    let previewImage = $state<string | null>(null);
    let previewImageName = $state("");

    async function loadStorageUsage() {
        try {
            storageUsage = await getUserStorageUsage();
        } catch (e) {
            console.error("Failed to load storage usage", e);
        }
    }

    const tabs = [
        { id: "general", label: "Informations", icon: User },
        { id: "stats", label: "Statistiques", icon: Activity },
        { id: "campaigns", label: "Campagnes", icon: Swords },
        { id: "characters", label: "Personnages", icon: Skull },
        { id: "bestiary", label: "Bestiaire", icon: BookOpen },
        { id: "media", label: "Médiathèque", icon: Camera },
    ];

    onMount(async () => {
        const urlParams = new URLSearchParams(window.location.search);
        const tab = urlParams.get("tab");
        if (tab && tabs.find((t) => t.id === tab)) {
            activeTab = tab;
        }

        try {
            const { data, error } = await authClient.getSession();
            if (error || !data?.user) {
                goto("/login");
                return;
            }
            user = data.user;
            editName = user.name || "";
            editEmail = user.email || "";

            await Promise.all([
                fetchStats(), 
                fetchCampaigns(), 
                fetchCharacters(), 
                fetchMonsters(),
                fetchMedia(),
                loadStorageUsage()
            ]);
        } catch (e) {
            console.error(e);
        } finally {
            loading = false;
        }
    });

    // Tab switcher
    function setTab(tabId: string) {
        activeTab = tabId;
        const url = new URL(window.location.href);
        url.searchParams.set("tab", tabId);
        goto(url.toString(), { replaceState: true, keepFocus: true });
    }

    async function fetchStats() {
        try {
            stats = await fetchUserStats();
        } catch (e) {
            console.error("Failed to fetch stats", e);
        }
    }

    async function fetchCampaigns() {
        try {
            campaigns = await fetchUserCampaigns();
        } catch (e) {
            console.error("Failed to fetch campaigns", e);
        }
    }

    async function fetchCharacters() {
        try {
            characters = await fetchUserCharacters();
        } catch (e) {
            console.error("Failed to fetch characters", e);
        }
    }

    async function handleAvatarUpload(event: Event) {
        const input = event.target as HTMLInputElement;
        const file = input.files?.[0];
        if (!file) return;

        const validation = validateImage(file, 2); // 2MB limit for avatar
        if (!validation.valid) {
            avatarError = validation.error || "Image invalide.";
            return;
        }

        avatarUploading = true;
        avatarError = "";
        avatarSuccess = false;

        try {
            const publicUrl = await uploadImage(file);
            await authClient.updateUser({ avatar_url: publicUrl });

            // Update the local user state
            if (user) {
                user = { ...user, image: publicUrl };
            }
            avatarSuccess = true;
            await loadStorageUsage();
            setTimeout(() => { avatarSuccess = false; }, 2500);
        } catch (e: any) {
            console.error("Failed to upload avatar", e);
            avatarError = e.message || "Erreur lors de l'upload de la photo.";
        } finally {
            avatarUploading = false;
            // Reset the input so the same file can be re-selected
            if (input) input.value = "";
        }
    }

    async function fetchMedia() {
        try {
            mediaLoading = true;
            mediaAssets = await fetchMediaLibrary();
        } catch (e) {
            console.error("Failed to fetch media library", e);
        } finally {
            mediaLoading = false;
        }
    }

    async function handleMediaUpload(event: Event) {
        const input = event.target as HTMLInputElement;
        const files = input.files;
        if (!files || files.length === 0) return;

        mediaUploading = true;
        mediaError = "";

        try {
            for (let i = 0; i < files.length; i++) {
                const file = files[i];
                const validation = validateImage(file, 10);
                if (!validation.valid) {
                    throw new Error(validation.error || `Le fichier "${file.name}" est invalide.`);
                }
                await uploadToMediaLibrary(file);
            }
            await fetchMedia();
            await loadStorageUsage();
        } catch (e: any) {
            console.error("Failed to upload to media library", e);
            mediaError = e.message || "Erreur lors du téléversement.";
        } finally {
            mediaUploading = false;
            if (input) input.value = "";
        }
    }

    async function handleDeleteMedia(id: string, url: string) {
        if (!confirm("Voulez-vous vraiment supprimer cette image de votre médiathèque ?")) return;

        try {
            await deleteFromMediaLibrary(id, url);
            await fetchMedia();
            await loadStorageUsage();
        } catch (e: any) {
            console.error("Failed to delete media asset", e);
            alert("Erreur lors de la suppression de l'image : " + (e.message || "Erreur inconnue"));
        }
    }

    async function fetchMonsters() {
        try {
            monsters = await fetchUserMonsters();
        } catch (e) {
            console.error("Failed to fetch monsters", e);
        }
    }

    async function handleDeleteAccount() {
        deletingAccount = true;
        deleteError = "";
        try {
            await deleteUserAccount({
                deleteMonsters,
                deleteGames,
                deleteCharacters,
                deleteImages
            });
            await authClient.signOut();
            goto("/");
        } catch (e: any) {
            console.error("Failed to delete account", e);
            deleteError = e.message || "Une erreur est survenue lors de la suppression du compte.";
            deletingAccount = false;
        }
    }

    async function handleUpdateName(event: SubmitEvent) {
        event.preventDefault();
        if (!editName.trim()) {
            nameError = "Le pseudo ne peut pas être vide.";
            return;
        }
        nameSaving = true;
        nameError = "";
        nameSuccess = false;
        try {
            await authClient.updateUser({ name: editName.trim() });
            if (user) {
                user = { ...user, name: editName.trim() };
            }
            nameSuccess = true;
            setTimeout(() => { nameSuccess = false; }, 3000);
        } catch (e: any) {
            console.error(e);
            nameError = e.message || "Erreur lors de la mise à jour du pseudo.";
        } finally {
            nameSaving = false;
        }
    }

    async function handleUpdateEmail(event: SubmitEvent) {
        event.preventDefault();
        if (!editEmail.trim()) {
            emailError = "L'adresse email ne peut pas être vide.";
            return;
        }
        emailSaving = true;
        emailError = "";
        emailSuccess = false;
        try {
            await authClient.updateEmail(editEmail.trim());
            emailSuccess = true;
        } catch (e: any) {
            console.error(e);
            emailError = e.message || "Erreur lors de la mise à jour de l'email.";
        } finally {
            emailSaving = false;
        }
    }

    async function handleUpdatePassword(event: SubmitEvent) {
        event.preventDefault();
        if (newPassword.length < 6) {
            passwordError = "Le mot de passe doit faire au moins 6 caractères.";
            return;
        }
        if (newPassword !== confirmPassword) {
            passwordError = "Les mots de passe ne correspondent pas.";
            return;
        }
        passwordSaving = true;
        passwordError = "";
        passwordSuccess = false;
        try {
            await authClient.updatePassword(newPassword);
            passwordSuccess = true;
            newPassword = "";
            confirmPassword = "";
            setTimeout(() => { passwordSuccess = false; }, 3000);
        } catch (e: any) {
            console.error(e);
            passwordError = e.message || "Erreur lors de la mise à jour du mot de passe.";
        } finally {
            passwordSaving = false;
        }
    }
</script>

<div class="min-h-screen bg-cream">
    <Header />

    <main class="max-w-7xl mx-auto px-4 md:px-8 py-12">
        <!-- Header -->
        <div class="mb-8">
            <h1 class="font-display font-bold text-4xl text-dark-gray mb-2">
                Mon Profil
            </h1>
            <p class="text-stone-500">
                Gérez vos informations, vos statistiques et vos campagnes.
            </p>
        </div>

        {#if loading}
            <div class="flex justify-center items-center h-64">
                <div
                    class="animate-spin rounded-full h-12 w-12 border-b-2 border-burnt-orange"
                ></div>
            </div>
        {:else}
            <!-- Tabs Navigation -->
            <div
                class="flex gap-2 mb-8 border-b border-stone-200 overflow-x-auto"
            >
                {#each tabs as tab}
                    <button
                        onclick={() => setTab(tab.id)}
                        class="flex items-center gap-2 px-4 py-3 font-medium text-sm transition-all relative {activeTab ===
                        tab.id
                            ? 'text-burnt-orange'
                            : 'text-stone-500 hover:text-dark-gray'}"
                    >
                        <tab.icon size={18} />
                        {tab.label}
                        {#if activeTab === tab.id}
                            <div
                                class="absolute bottom-0 left-0 w-full h-0.5 bg-burnt-orange rounded-t-full"
                            ></div>
                        {/if}
                    </button>
                {/each}
            </div>

            <!-- Content Area -->
            <div
                class="bg-white rounded-2xl shadow-sm border border-stone-100 p-6 md:p-8 min-h-[500px]"
            >
                <!-- General Info -->
                {#if activeTab === "general"}
                    <div
                        class="animate-in fade-in slide-in-from-bottom-4 duration-300 space-y-8"
                    >
                        <div class="grid grid-cols-1 lg:grid-cols-12 gap-8">
                            <!-- Left Column: Avatar & Profile Info (5 cols) -->
                            <div class="lg:col-span-5 bg-stone-50/50 rounded-2xl border border-stone-200/60 p-6 flex flex-col items-center text-center space-y-6">
                                <div class="w-full flex justify-between items-center pb-2 border-b border-stone-100">
                                    <h3 class="text-xs font-bold uppercase tracking-wider text-stone-400">Photo & Identité</h3>
                                    <span class="inline-flex items-center px-2 py-0.5 rounded-full text-[10px] font-semibold bg-stone-200 text-stone-600 border border-stone-300/30">
                                        Actif
                                    </span>
                                </div>
                                
                                <!-- Avatar Upload -->
                                <div class="relative group">
                                    <div class="absolute inset-0 bg-burnt-orange/10 rounded-full blur-md opacity-0 group-hover:opacity-100 transition-opacity"></div>
                                    <img
                                        src={user?.image ||
                                            `https://ui-avatars.com/api/?name=${user?.name}&background=random`}
                                        alt={user?.name}
                                        class="w-28 h-28 rounded-full object-cover border-4 border-white shadow-md transition-all duration-200 group-hover:scale-[1.02] {avatarUploading ? 'opacity-50' : ''}"
                                    />
                                    <!-- Upload overlay -->
                                    <button
                                        type="button"
                                        onclick={() => fileInput?.click()}
                                        disabled={avatarUploading}
                                        class="absolute inset-0 flex flex-col items-center justify-center rounded-full cursor-pointer bg-black/40 opacity-0 group-hover:opacity-100 transition-all duration-200 disabled:opacity-50"
                                        title="Changer la photo de profil"
                                    >
                                        {#if avatarUploading}
                                            <Loader2 size={24} class="text-white animate-spin" />
                                        {:else if avatarSuccess}
                                            <Check size={24} class="text-green-400 animate-bounce" />
                                        {:else}
                                            <Camera size={24} class="text-white drop-shadow-md mb-1" />
                                            <span class="text-white text-xs font-semibold drop-shadow-md">Modifier</span>
                                        {/if}
                                    </button>
                                    <input
                                        type="file"
                                        accept="image/*"
                                        class="hidden"
                                        bind:this={fileInput}
                                        onchange={handleAvatarUpload}
                                    />
                                    {#if avatarError}
                                        <div class="absolute -bottom-12 left-1/2 -translate-x-1/2 whitespace-nowrap bg-red-50 text-red-600 text-xs px-3 py-1.5 rounded-xl border border-red-100 shadow-sm z-10 animate-fade-in">
                                            {avatarError}
                                        </div>
                                    {/if}
                                    {#if avatarSuccess}
                                        <div class="absolute -bottom-12 left-1/2 -translate-x-1/2 whitespace-nowrap bg-green-50 text-green-600 text-xs px-3 py-1.5 rounded-xl border border-green-100 shadow-sm z-10 animate-fade-in">
                                            Photo mise à jour !
                                        </div>
                                    {/if}
                                </div>

                                <div class="w-full space-y-1">
                                    <h4 class="font-display font-black text-xl text-dark-gray leading-tight">{user?.name}</h4>
                                    <p class="text-xs text-stone-400 font-mono truncate max-w-full px-2" title={user?.email}>{user?.email}</p>
                                </div>

                                <!-- Pseudo Edit Form -->
                                <form onsubmit={handleUpdateName} class="w-full space-y-3 pt-5 border-t border-stone-200/60">
                                    <div class="text-left">
                                        <label class="block text-[10px] font-bold text-dark-gray/60 uppercase tracking-wider mb-2 ml-1" for="username">
                                            Modifier le pseudo
                                        </label>
                                        <div class="flex gap-2">
                                            <input
                                                id="username"
                                                type="text"
                                                bind:value={editName}
                                                required
                                                class="flex-1 px-4 py-2.5 rounded-xl border border-stone-200 focus:outline-none focus:ring-2 focus:ring-burnt-orange/20 focus:border-burnt-orange transition-all bg-white text-sm"
                                                placeholder="Pseudo"
                                            />
                                            <button
                                                type="submit"
                                                disabled={nameSaving || editName.trim() === user?.name || !editName.trim()}
                                                class="px-4 py-2.5 bg-burnt-orange hover:bg-opacity-95 text-white font-bold rounded-xl text-sm transition-all hover:shadow-md cursor-pointer disabled:bg-stone-200 disabled:text-stone-400 disabled:shadow-none disabled:cursor-not-allowed flex items-center gap-1.5 shrink-0"
                                            >
                                                {#if nameSaving}
                                                    <Loader2 size={16} class="animate-spin" />
                                                {/if}
                                                Enregistrer
                                            </button>
                                        </div>
                                    </div>
                                    
                                    {#if nameError}
                                        <p class="text-xs text-red-600 text-left">{nameError}</p>
                                    {/if}
                                    {#if nameSuccess}
                                        <p class="text-xs text-green-600 font-medium flex items-center gap-1 text-left">
                                            <Check size={14} /> Pseudo mis à jour !
                                        </p>
                                    {/if}
                                </form>
                            </div>

                            <!-- Right Column: Settings & Security (7 cols) -->
                            <div class="lg:col-span-7 space-y-6">
                                <!-- Section 2: Adresse Email -->
                                <div class="bg-stone-50/20 rounded-2xl border border-stone-200/60 p-6 space-y-4 shadow-3xs">
                                    <div class="flex items-center gap-2 pb-2 border-b border-stone-100">
                                        <Mail size={18} class="text-burnt-orange" />
                                        <h3 class="font-display font-bold text-base text-dark-gray">Adresse E-mail</h3>
                                    </div>
                                    
                                    <div class="bg-burnt-orange/5 border border-burnt-orange/10 rounded-xl p-3.5 text-xs text-dark-gray/70 leading-relaxed space-y-1">
                                        <p class="font-bold text-dark-gray">ℹ️ Processus de changement d'e-mail :</p>
                                        <p>
                                            Par sécurité, nous exigeons la confirmation via <strong>deux e-mails distincts</strong>.
                                            Vous devrez cliquer sur le lien de confirmation envoyé à votre <strong>ancienne adresse e-mail</strong>,
                                            puis sur celui envoyé à votre <strong>nouvelle adresse e-mail</strong> pour finaliser la mise à jour.
                                        </p>
                                    </div>

                                    <form onsubmit={handleUpdateEmail} class="space-y-3">
                                        <div>
                                            <label class="block text-[10px] font-bold text-dark-gray/60 uppercase tracking-wider mb-2 ml-1" for="email">
                                                Nouvelle adresse e-mail
                                            </label>
                                            <div class="flex gap-2">
                                                <div class="relative flex-1">
                                                    <div class="absolute inset-y-0 left-0 pl-3.5 flex items-center pointer-events-none text-stone-400">
                                                        <Mail size={16} />
                                                    </div>
                                                    <input
                                                        id="email"
                                                        type="email"
                                                        bind:value={editEmail}
                                                        required
                                                        class="w-full pl-10 pr-4 py-2.5 rounded-xl border border-stone-200 focus:outline-none focus:ring-2 focus:ring-burnt-orange/20 focus:border-burnt-orange transition-all bg-white text-sm"
                                                        placeholder="nouveau@mail.com"
                                                    />
                                                </div>
                                                <button
                                                    type="submit"
                                                    disabled={emailSaving || editEmail.trim() === user?.email || !editEmail.trim()}
                                                    class="px-4 py-2.5 bg-burnt-orange hover:bg-opacity-95 text-white font-bold rounded-xl text-sm transition-all hover:shadow-md cursor-pointer disabled:bg-stone-200 disabled:text-stone-400 disabled:shadow-none disabled:cursor-not-allowed flex items-center gap-1.5 shrink-0"
                                                >
                                                    {#if emailSaving}
                                                        <Loader2 size={16} class="animate-spin" />
                                                    {/if}
                                                    Enregistrer
                                                </button>
                                            </div>
                                        </div>
                                        
                                        {#if emailError}
                                            <p class="text-xs text-red-600">{emailError}</p>
                                        {/if}
                                        {#if emailSuccess}
                                            <div class="bg-blue-50 text-blue-700 p-3.5 rounded-xl text-xs border border-blue-100 leading-relaxed animate-fade-in space-y-2">
                                                <p class="font-bold">Demande de changement envoyée !</p>
                                                <ul class="list-disc pl-4 space-y-1">
                                                    <li>Vérifiez votre <strong>ancienne boîte mail</strong> et cliquez sur le lien de validation.</li>
                                                    <li>Vérifiez la boîte de <strong>{editEmail}</strong> et cliquez sur le lien de confirmation reçu pour valider la nouvelle adresse.</li>
                                                </ul>
                                            </div>
                                        {/if}
                                    </form>
                                </div>

                                <!-- Section 3: Sécurité & Mot de passe -->
                                <div class="bg-stone-50/20 rounded-2xl border border-stone-200/60 p-6 space-y-4 shadow-3xs">
                                    <div class="flex items-center gap-2 pb-2 border-b border-stone-100">
                                        <Lock size={18} class="text-burnt-orange" />
                                        <h3 class="font-display font-bold text-base text-dark-gray">Sécurité & Mot de passe</h3>
                                    </div>
                                    
                                    <form onsubmit={handleUpdatePassword} class="space-y-4">
                                        <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
                                            <div>
                                                <label class="block text-[10px] font-bold text-dark-gray/60 uppercase tracking-wider mb-2 ml-1" for="new-password">
                                                    Nouveau mot de passe
                                                </label>
                                                <div class="relative">
                                                    <div class="absolute inset-y-0 left-0 pl-3.5 flex items-center pointer-events-none text-stone-400">
                                                        <Lock size={16} />
                                                    </div>
                                                    <input
                                                        id="new-password"
                                                        type="password"
                                                        bind:value={newPassword}
                                                        required
                                                        class="w-full pl-10 pr-4 py-2.5 rounded-xl border border-stone-200 focus:outline-none focus:ring-2 focus:ring-burnt-orange/20 focus:border-burnt-orange transition-all bg-white text-sm"
                                                        placeholder="Min. 6 caractères"
                                                    />
                                                </div>
                                            </div>
                                            <div>
                                                <label class="block text-[10px] font-bold text-dark-gray/60 uppercase tracking-wider mb-2 ml-1" for="confirm-password">
                                                    Confirmer le mot de passe
                                                </label>
                                                <div class="relative">
                                                    <div class="absolute inset-y-0 left-0 pl-3.5 flex items-center pointer-events-none text-stone-400">
                                                        <Lock size={16} />
                                                    </div>
                                                    <input
                                                        id="confirm-password"
                                                        type="password"
                                                        bind:value={confirmPassword}
                                                        required
                                                        class="w-full pl-10 pr-4 py-2.5 rounded-xl border border-stone-200 focus:outline-none focus:ring-2 focus:ring-burnt-orange/20 focus:border-burnt-orange transition-all bg-white text-sm"
                                                        placeholder="Confirmez"
                                                    />
                                                </div>
                                            </div>
                                        </div>

                                        <div class="flex justify-end pt-2">
                                            <button
                                                type="submit"
                                                disabled={passwordSaving || !newPassword || !confirmPassword}
                                                class="px-5 py-2.5 bg-burnt-orange hover:bg-opacity-95 text-white font-bold rounded-xl text-sm transition-all hover:shadow-md cursor-pointer disabled:bg-stone-200 disabled:text-stone-400 disabled:shadow-none disabled:cursor-not-allowed flex items-center gap-1.5"
                                            >
                                                {#if passwordSaving}
                                                    <Loader2 size={16} class="animate-spin" />
                                                {/if}
                                                Mettre à jour le mot de passe
                                            </button>
                                        </div>
                                        
                                        {#if passwordError}
                                            <p class="text-xs text-red-600">{passwordError}</p>
                                        {/if}
                                        {#if passwordSuccess}
                                            <p class="text-xs text-green-600 font-medium flex items-center gap-1">
                                                <Check size={14} /> Mot de passe modifié !
                                            </p>
                                        {/if}
                                    </form>
                                </div>
                            </div>
                        </div>

                        <!-- Danger Zone -->
                        <div class="mt-12 pt-8 border-t border-stone-200 space-y-6">
                            <div>
                                <h3 class="text-lg font-display font-black text-red-600">Zone de danger</h3>
                                <p class="text-sm text-stone-500">
                                    Gérez la suppression définitive de vos données et de votre compte utilisateur.
                                </p>
                            </div>

                            <div class="bg-red-50/20 border border-red-100 rounded-2xl p-6 flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4">
                                <div class="space-y-1">
                                    <h4 class="font-bold text-dark-gray text-sm">Supprimer mon compte QuestHub</h4>
                                    <p class="text-xs text-stone-500">
                                        Cette action supprimera définitivement votre profil d'utilisateur et vous permettra de choisir les données associées à effacer.
                                    </p>
                                </div>
                                <button
                                    type="button"
                                    onclick={() => showDeleteConfirm = true}
                                    class="px-5 py-2.5 bg-red-600 hover:bg-red-700 text-white font-bold rounded-xl text-sm transition-all hover:shadow-md cursor-pointer flex items-center gap-2 shrink-0"
                                >
                                    Supprimer mon compte
                                </button>
                            </div>
                        </div>
                    </div>
                {/if}

                <!-- Stats -->
                {#if activeTab === "stats"}
                    <div
                        class="animate-in fade-in slide-in-from-bottom-4 duration-300"
                    >
                        <div
                            class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6"
                        >
                            <div
                                class="bg-stone-50 rounded-2xl p-6 border border-stone-100 hover:shadow-md transition-shadow"
                            >
                                <p class="text-dark-gray/60 font-medium mb-2">
                                    Parties Jouées
                                </p>
                                <p
                                    class="text-4xl font-display font-bold text-burnt-orange"
                                >
                                    {stats?.games_played || 0}
                                </p>
                            </div>
                            <div
                                class="bg-stone-50 rounded-2xl p-6 border border-stone-100 hover:shadow-md transition-shadow"
                            >
                                <p class="text-dark-gray/60 font-medium mb-2">
                                    Parties en tant que GM
                                </p>
                                <p
                                    class="text-4xl font-display font-bold text-burnt-orange"
                                >
                                    {stats?.games_mastered || 0}
                                </p>
                            </div>
                            <div
                                class="bg-stone-50 rounded-2xl p-6 border border-stone-100 hover:shadow-md transition-shadow"
                            >
                                <p class="text-dark-gray/60 font-medium mb-2">
                                    Personnages Créés
                                </p>
                                <p
                                    class="text-4xl font-display font-bold text-burnt-orange"
                                >
                                    {stats?.characters_created || 0}
                                </p>
                            </div>
                        </div>
                    </div>
                {/if}

                <!-- Campaigns -->
                {#if activeTab === "campaigns"}
                    <div
                        class="animate-in fade-in slide-in-from-bottom-4 duration-300"
                    >
                        {#if campaigns.length > 0}
                            <div class="grid gap-4">
                                {#each campaigns as camp}
                                    <div
                                        class="flex items-center gap-4 p-4 rounded-2xl border border-stone-100 hover:bg-stone-50 transition-colors"
                                    >
                                        <img
                                            src={camp.game_image_url ||
                                                `https://placehold.co/100x100/3D405B/F9F7F2?text=${camp.game_name}`}
                                            alt={camp.game_name}
                                            class="w-16 h-16 rounded-xl object-cover"
                                        />
                                        <div class="flex-1">
                                            <h3
                                                class="font-bold text-lg text-dark-gray"
                                            >
                                                {camp.game_name}
                                            </h3>
                                            <div
                                                class="flex items-center gap-2 text-sm text-dark-gray/60"
                                            >
                                                <span
                                                    >Personnage: <span
                                                        class="font-medium text-burnt-orange"
                                                        >{camp.character_name}</span
                                                    ></span
                                                >
                                            </div>
                                        </div>
                                        <a
                                            href="/table/{camp.game_id}"
                                            class="px-4 py-2 bg-dark-gray text-white rounded-xl text-sm hover:bg-opacity-90 transition-opacity"
                                        >
                                            Rejoindre
                                        </a>
                                    </div>
                                {/each}
                            </div>
                        {:else}
                            <div class="text-center py-12 text-dark-gray/40">
                                <p>Aucune campagne trouvée.</p>
                            </div>
                        {/if}
                    </div>
                {/if}

                <!-- Characters Tab -->
                {#if activeTab === "characters"}
                    <div
                        class="animate-in fade-in slide-in-from-bottom-4 duration-300"
                    >
                        {#if characters.length > 0}
                            <div class="grid grid-cols-1 sm:grid-cols-2 gap-6">
                                {#each characters as char}
                                    <div
                                        class="bg-white rounded-2xl border border-stone-100 p-5 shadow-xs hover:shadow-md hover:border-burnt-orange/20 transition-all flex flex-col justify-between h-full group"
                                    >
                                        <div class="space-y-4">
                                            <!-- Top section -->
                                            <div class="flex items-start justify-between gap-4">
                                                <div class="flex items-center gap-4">
                                                    <img
                                                        src={char.avatar_url || `https://api.dicebear.com/9.x/avataaars/svg?seed=${char.name}`}
                                                        alt={char.name}
                                                        class="w-16 h-16 rounded-xl object-cover bg-stone-50 border-2 border-stone-100"
                                                    />
                                                    <div>
                                                        <h3
                                                            class="font-display font-bold text-xl text-dark-gray group-hover:text-burnt-orange transition-colors"
                                                        >
                                                            {char.name}
                                                        </h3>
                                                        <p class="text-stone-400 text-sm font-medium">
                                                            {char.race} {char.sub_race ? `(${char.sub_race})` : ''}
                                                        </p>
                                                    </div>
                                                </div>
                                                <div>
                                                    {#if char.association_type === 'assigned'}
                                                        <span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-semibold bg-blue-50 text-blue-700 border border-blue-100">
                                                            Attribué
                                                        </span>
                                                    {:else}
                                                        <span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-semibold bg-purple-50 text-purple-700 border border-purple-100">
                                                            Créé (MJ)
                                                        </span>
                                                    {/if}
                                                </div>
                                            </div>
 
                                            <!-- Vitals (HP & XP/Level) -->
                                            <div class="space-y-2 bg-stone-50 p-3 rounded-xl border border-stone-100">
                                                <div class="flex justify-between text-xs font-bold text-stone-500">
                                                    <span>Points de Vie</span>
                                                    <span>{char.current_hp} / {char.max_hp} PV</span>
                                                </div>
                                                <div class="h-2 bg-stone-200 rounded-full overflow-hidden">
                                                    <div
                                                        class="h-full rounded-full transition-all duration-300 {char.current_hp < char.max_hp / 3 ? 'bg-red-500' : 'bg-green-500'}"
                                                        style="width: {(char.current_hp / char.max_hp) * 100}%"
                                                    ></div>
                                                </div>
                                                {#if char.experience !== undefined}
                                                    <div class="flex justify-between text-[11px] text-stone-400 font-semibold pt-1 border-t border-stone-200/50 mt-1">
                                                        <span>Expérience</span>
                                                        <span>{char.experience} XP</span>
                                                    </div>
                                                {/if}
                                            </div>
 
                                            <!-- Associated Game info -->
                                            {#if char.game}
                                                <div class="flex items-center gap-2 pt-2">
                                                    <img
                                                        src={char.game.image_url || `https://placehold.co/100x100/3D405B/F9F7F2?text=${char.game.name}`}
                                                        alt={char.game.name}
                                                        class="w-6 h-6 rounded-md object-cover"
                                                    />
                                                    <span class="text-xs font-medium text-stone-500 truncate">
                                                        Partie : <span class="font-semibold text-dark-gray">{char.game.name}</span>
                                                    </span>
                                                </div>
                                            {/if}
                                        </div>
 
                                        <!-- Join button -->
                                        <div class="pt-5 mt-4 border-t border-stone-100/60 flex justify-end">
                                            <a
                                                href="/table/{char.game_id || char.game?.id}"
                                                class="px-4 py-2 bg-burnt-orange text-white rounded-xl text-sm font-bold shadow-xs hover:bg-opacity-90 hover:shadow-md hover:-translate-y-0.5 transition-all flex items-center gap-2"
                                            >
                                                {#if char.association_type === 'assigned'}
                                                    Jouer ce personnage
                                                {:else}
                                                    Rejoindre la partie (MJ)
                                                {/if}
                                            </a>
                                        </div>
                                    </div>
                                {/each}
                            </div>
                        {:else}
                            <div class="text-center py-20 bg-stone-50/50 rounded-2xl border border-stone-100 flex flex-col items-center justify-center gap-4">
                                <div class="w-16 h-16 bg-white rounded-full flex items-center justify-center text-stone-300 shadow-sm border border-stone-100">
                                    <Skull size={32} />
                                </div>
                                <div class="space-y-1">
                                    <h4 class="font-bold text-dark-gray text-lg">Aucun personnage</h4>
                                    <p class="text-stone-400 text-sm max-w-sm">
                                        Vous n'avez aucun personnage pour le moment. Rejoignez une partie ou demandez à votre Maître du Jeu de vous en assigner un !
                                    </p>
                                </div>
                            </div>
                        {/if}
                    </div>
                {/if}

                <!-- Bestiary Tab -->
                {#if activeTab === "bestiary"}
                    <div
                        class="animate-in fade-in slide-in-from-bottom-4 duration-300"
                    >
                        {#if monsters.length > 0}
                            <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6">
                                {#each monsters as monster}
                                    <div
                                        class="bg-white rounded-2xl border border-stone-100 p-5 shadow-xs hover:shadow-md hover:border-burnt-orange/20 transition-all flex flex-col justify-between h-full group"
                                    >
                                        <div class="space-y-4">
                                            <!-- Top section -->
                                            <div class="flex items-start justify-between gap-4">
                                                <div class="flex items-center gap-4">
                                                    <img
                                                        src={monster.avatar_url || `https://api.dicebear.com/9.x/bottts/svg?seed=${monster.name}`}
                                                        alt={monster.name}
                                                        class="w-16 h-16 rounded-xl object-cover bg-stone-50 border-2 border-stone-100"
                                                    />
                                                    <div>
                                                        <h3
                                                            class="font-display font-bold text-xl text-dark-gray group-hover:text-burnt-orange transition-colors"
                                                        >
                                                            {monster.name}
                                                        </h3>
                                                        <p class="text-stone-400 text-sm font-medium">
                                                            {monster.race || 'Monstre'}
                                                        </p>
                                                    </div>
                                                </div>
                                            </div>

                                            <!-- Vitals (HP & AC) -->
                                            <div class="grid grid-cols-2 gap-2 bg-stone-50 p-3 rounded-xl border border-stone-100 text-center">
                                                <div>
                                                    <span class="block text-[10px] uppercase tracking-wider text-stone-400 font-bold">Classe d'armure</span>
                                                    <span class="text-lg font-bold text-dark-gray">{monster.armor_class ?? 10}</span>
                                                </div>
                                                <div>
                                                    <span class="block text-[10px] uppercase tracking-wider text-stone-400 font-bold">Points de Vie Max</span>
                                                    <span class="text-lg font-bold text-dark-gray">{monster.max_hp} PV</span>
                                                </div>
                                            </div>

                                            <!-- Associated Game info -->
                                            {#if monster.game}
                                                <div class="flex items-center gap-2 pt-2">
                                                    <img
                                                        src={monster.game.image_url || `https://placehold.co/100x100/3D405B/F9F7F2?text=${monster.game.name}`}
                                                        alt={monster.game.name}
                                                        class="w-6 h-6 rounded-md object-cover"
                                                    />
                                                    <span class="text-xs font-medium text-stone-500 truncate font-semibold">
                                                        Partie : <span class="font-bold text-dark-gray">{monster.game.name}</span>
                                                    </span>
                                                </div>
                                            {/if}
                                        </div>

                                        <!-- Actions -->
                                        <div class="pt-5 mt-4 border-t border-stone-100/60 flex justify-between items-center gap-2">
                                            <a
                                                href="/table/{monster.game_id || monster.game?.id}/gm/settings?tab=bestiary"
                                                class="text-xs font-bold text-stone-500 hover:text-burnt-orange transition-colors"
                                            >
                                                Modifier le monstre
                                            </a>
                                            <a
                                                href="/table/{monster.game_id || monster.game?.id}"
                                                class="px-3 py-1.5 bg-dark-gray text-white rounded-lg text-xs font-bold shadow-xs hover:bg-opacity-90 hover:shadow-md transition-all"
                                            >
                                                Aller à la table
                                            </a>
                                        </div>
                                    </div>
                                {/each}
                            </div>
                        {:else}
                            <div class="text-center py-20 bg-stone-50/50 rounded-2xl border border-stone-100 flex flex-col items-center justify-center gap-4">
                                <div class="w-16 h-16 bg-white rounded-full flex items-center justify-center text-stone-300 shadow-sm border border-stone-100">
                                    <BookOpen size={32} />
                                </div>
                                <div class="space-y-1">
                                    <h4 class="font-bold text-dark-gray text-lg">Aucun monstre créé</h4>
                                    <p class="text-stone-400 text-sm max-w-sm">
                                        Vous n'avez pas encore créé de monstre. Allez dans les paramètres d'une partie dont vous êtes le GM pour en ajouter au bestiaire !
                                    </p>
                                </div>
                            </div>
                        {/if}
                    </div>
                {/if}

                {#if activeTab === "media"}
                    <div class="animate-in fade-in slide-in-from-bottom-4 duration-300 space-y-6">
                        <div class="flex flex-col sm:flex-row justify-between items-center gap-4 bg-white p-6 rounded-2xl border border-stone-200/60 shadow-2xs">
                            <div class="space-y-1 text-center sm:text-left">
                                <h3 class="font-display font-black text-xl text-dark-gray">Médiathèque Personnelle</h3>
                                <p class="text-xs text-stone-400">Stockez et réutilisez vos images pour vos cartes, monstres, jetons et illustrations.</p>
                            </div>
                            <div>
                                <input
                                    type="file"
                                    multiple
                                    accept="image/*"
                                    bind:this={mediaFileInput}
                                    class="hidden"
                                    onchange={handleMediaUpload}
                                />
                                <button
                                    onclick={() => mediaFileInput?.click()}
                                    disabled={mediaUploading}
                                    class="px-5 py-2.5 bg-burnt-orange hover:bg-opacity-95 text-white font-bold rounded-xl text-sm transition-all flex items-center gap-2 cursor-pointer shadow-sm disabled:opacity-60 disabled:cursor-not-allowed"
                                >
                                    {#if mediaUploading}
                                        <Loader2 size={16} class="animate-spin" />
                                        Téléversement...
                                    {:else}
                                        <RefreshCw size={16} />
                                        Ajouter des images
                                    {/if}
                                </button>
                            </div>
                        </div>

                        <!-- Stockage Quota Box -->
                        <div class="bg-white p-6 rounded-2xl border border-stone-200/60 shadow-2xs space-y-4">
                            <div class="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-2">
                                <div>
                                    <h4 class="font-bold text-stone-850 text-sm">Utilisation du stockage</h4>
                                    <p class="text-xs text-stone-400">Quota total alloué par compte : 50 Mo</p>
                                </div>
                                <span class="text-xs font-bold font-mono text-stone-600">
                                    {(storageUsage / (1024 * 1024)).toFixed(2)} Mo / 50 Mo ({((storageUsage / (50 * 1024 * 1024)) * 100).toFixed(1)}%)
                                </span>
                            </div>
                            
                            <!-- Progress Bar -->
                            <div class="w-full bg-stone-100 rounded-full h-2.5 overflow-hidden">
                                <div 
                                    class="h-full rounded-full transition-all duration-500 {((storageUsage / (50 * 1024 * 1024)) * 100) > 85 ? 'bg-red-500' : (((storageUsage / (50 * 1024 * 1024)) * 100) > 50 ? 'bg-amber-500' : 'bg-burnt-orange')}" 
                                    style="width: {Math.min(100, (storageUsage / (50 * 1024 * 1024)) * 100)}%"
                                ></div>
                            </div>

                            <div class="flex items-start gap-2 bg-amber-50/50 border border-amber-200/40 rounded-xl p-3 text-amber-800 text-xs">
                                <ShieldAlert size={14} class="shrink-0 mt-0.5" />
                                <div class="space-y-1">
                                    <p class="font-bold">Conseil de stockage</p>
                                    <p>Le stockage direct étant limité, nous vous conseillons de <strong>privilégier l'utilisation de liens URL d'images externes</strong> (par ex. Imgur, Discord, etc.) dans vos configurations de personnages, cartes ou jetons plutôt que de téléverser des fichiers.</p>
                                </div>
                            </div>
                        </div>

                        {#if mediaError}
                            <div class="bg-red-50 text-red-600 px-4 py-3 rounded-xl text-sm border border-red-100">
                                {mediaError}
                            </div>
                        {/if}

                        {#if mediaLoading}
                            <div class="flex items-center justify-center py-20">
                                <Loader2 class="animate-spin text-burnt-orange" size={32} />
                            </div>
                        {:else if mediaAssets.length > 0}
                            <div class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 gap-4">
                                {#each mediaAssets as asset}
                                    <div class="group relative bg-white/70 backdrop-blur-md rounded-2xl border border-stone-200/50 overflow-hidden shadow-xs hover:-translate-y-1 hover:shadow-lg hover:border-burnt-orange/30 transition-all duration-300 flex flex-col justify-between">
                                        <!-- Image Container -->
                                        <div 
                                            onclick={() => { previewImage = asset.url; previewImageName = asset.name; }}
                                            class="relative aspect-square bg-stone-50 overflow-hidden flex items-center justify-center cursor-zoom-in"
                                            role="button"
                                            tabindex="0"
                                            onkeydown={(e) => e.key === 'Enter' && (previewImage = asset.url, previewImageName = asset.name)}
                                        >
                                            <img
                                                src={asset.url}
                                                alt={asset.name}
                                                class="w-full h-full object-cover transition-transform duration-500 ease-out group-hover:scale-105"
                                            />
                                            <!-- Blur Overlay on Hover -->
                                            <div class="absolute inset-0 bg-black/30 backdrop-blur-[2px] opacity-0 group-hover:opacity-100 transition-all duration-300 flex items-center justify-center gap-2">
                                                <button
                                                    onclick={(e) => {
                                                        e.stopPropagation();
                                                        navigator.clipboard.writeText(asset.url);
                                                        alert("URL copiée dans le presse-papiers !");
                                                    }}
                                                    class="p-2.5 bg-white text-stone-750 hover:bg-burnt-orange hover:text-white rounded-xl shadow-md scale-90 group-hover:scale-100 transition-all duration-350 ease-out cursor-pointer hover:rotate-6"
                                                    title="Copier le lien de l'image"
                                                >
                                                    <Link size={16} />
                                                </button>
                                                <button
                                                    onclick={(e) => {
                                                        e.stopPropagation();
                                                        handleDeleteMedia(asset.id, asset.url);
                                                    }}
                                                    class="p-2.5 bg-white text-stone-750 hover:bg-red-500 hover:text-white rounded-xl shadow-md scale-90 group-hover:scale-100 transition-all duration-350 ease-out cursor-pointer hover:-rotate-6"
                                                    title="Supprimer l'image"
                                                >
                                                    <Trash2 size={16} />
                                                </button>
                                            </div>
                                        </div>
                                        
                                        <!-- Info Details -->
                                        <div class="p-3.5 space-y-2.5 bg-gradient-to-b from-transparent to-stone-50/30">
                                            <p class="text-xs font-bold text-dark-gray truncate leading-tight group-hover:text-burnt-orange transition-colors" title={asset.name}>
                                                {asset.name}
                                            </p>
                                            <div class="flex items-center justify-between gap-1.5 pt-0.5">
                                                <!-- Size Pill -->
                                                <span class="inline-flex items-center text-[9px] font-bold font-sans text-stone-550 bg-stone-100 px-2 py-0.5 rounded-md border border-stone-200/30">
                                                    {Math.round(asset.size / 1024)} Ko
                                                </span>
                                                <!-- Format Pill -->
                                                <span class="inline-flex items-center text-[9px] font-black font-sans uppercase tracking-wide text-burnt-orange bg-burnt-orange/10 px-2 py-0.5 rounded-md border border-burnt-orange/15">
                                                    {asset.mime_type.split('/')[1] || 'IMAGE'}
                                                </span>
                                            </div>
                                        </div>
                                    </div>
                                {/each}
                            </div>
                        {:else}
                            <div class="text-center py-20 bg-stone-50/50 rounded-2xl border border-stone-100 flex flex-col items-center justify-center gap-4">
                                <div class="w-16 h-16 bg-white rounded-full flex items-center justify-center text-stone-300 shadow-sm border border-stone-100">
                                    <Camera size={32} />
                                </div>
                                <div class="space-y-1">
                                    <h4 class="font-bold text-dark-gray text-lg">Votre médiathèque est vide</h4>
                                    <p class="text-stone-400 text-sm max-w-sm mx-auto">
                                        Importez vos images de cartes, de jetons de combat ou d'illustrations pour les stocker et y accéder facilement à tout moment.
                                    </p>
                                </div>
                            </div>
                        {/if}
                    </div>
                {/if}

            </div>
        {/if}
    </main>

    {#if showDeleteConfirm}
        <!-- Confirmation Dialog -->
        <div class="fixed inset-0 bg-black/50 backdrop-blur-xs flex items-center justify-center p-4 z-50 animate-in fade-in duration-200">
            <div class="bg-white rounded-3xl p-6 md:p-8 max-w-xl w-full border border-stone-200 shadow-xl space-y-6">
                <div>
                    <h3 class="font-display font-black text-2xl text-red-600">Supprimer définitivement le compte ?</h3>
                    <p class="text-sm text-stone-500 mt-2 leading-relaxed">
                        Cette action est irréversible. Cochez ci-dessous les données que vous souhaitez supprimer en même temps que votre compte :
                    </p>
                </div>

                <!-- Checkboxes inside the popup -->
                <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
                    <label class="flex items-start gap-2.5 p-3 bg-stone-50 rounded-xl border border-stone-200/60 shadow-3xs cursor-pointer hover:border-red-200 hover:bg-red-50/10 transition-all select-none">
                        <input
                            type="checkbox"
                            bind:checked={deleteMonsters}
                            class="mt-1 rounded text-red-600 focus:ring-red-500 border-stone-300 w-4 h-4 cursor-pointer"
                        />
                        <div>
                            <span class="block text-xs font-bold text-dark-gray">Monstres du bestiaire</span>
                            <span class="block text-[10px] text-stone-400 mt-0.5">Supprime tous les monstres créés dans vos parties.</span>
                        </div>
                    </label>

                    <label class="flex items-start gap-2.5 p-3 bg-stone-50 rounded-xl border border-stone-200/60 shadow-3xs cursor-pointer hover:border-red-200 hover:bg-red-50/10 transition-all select-none">
                        <input
                            type="checkbox"
                            bind:checked={deleteGames}
                            class="mt-1 rounded text-red-600 focus:ring-red-500 border-stone-300 w-4 h-4 cursor-pointer"
                        />
                        <div>
                            <span class="block text-xs font-bold text-dark-gray">Tables & Campagnes (si MJ)</span>
                            <span class="block text-[10px] text-stone-400 mt-0.5">Supprime définitivement toutes les tables dont vous êtes le MJ.</span>
                        </div>
                    </label>

                    <label class="flex items-start gap-2.5 p-3 bg-stone-50 rounded-xl border border-stone-200/60 shadow-3xs cursor-pointer hover:border-red-200 hover:bg-red-50/10 transition-all select-none">
                        <input
                            type="checkbox"
                            bind:checked={deleteCharacters}
                            class="mt-1 rounded text-red-600 focus:ring-red-500 border-stone-300 w-4 h-4 cursor-pointer"
                        />
                        <div>
                            <span class="block text-xs font-bold text-dark-gray">Personnages de jeu</span>
                            <span class="block text-[10px] text-stone-400 mt-0.5">Supprime vos fiches de personnages ainsi que celles créées pour vos parties.</span>
                        </div>
                    </label>

                    <label class="flex items-start gap-2.5 p-3 bg-stone-50 rounded-xl border border-stone-200/60 shadow-3xs cursor-pointer hover:border-red-200 hover:bg-red-50/10 transition-all select-none">
                        <input
                            type="checkbox"
                            bind:checked={deleteImages}
                            class="mt-1 rounded text-red-600 focus:ring-red-500 border-stone-300 w-4 h-4 cursor-pointer"
                        />
                        <div>
                            <span class="block text-xs font-bold text-dark-gray">Photos, cartes & médias</span>
                            <span class="block text-[10px] text-stone-400 mt-0.5">Supprime vos images et photos importées de notre stockage.</span>
                        </div>
                    </label>
                </div>

                <p class="text-xs text-stone-400 leading-relaxed italic bg-stone-50 p-3 rounded-lg border border-stone-200/50">
                    Remarque : L'historique et les données des autres joueurs sur vos tables supprimées seront également perdus.
                </p>

                <div class="flex gap-3 justify-end pt-2 border-t border-stone-100">
                    <button
                        type="button"
                        onclick={() => showDeleteConfirm = false}
                        class="px-4 py-2 border border-stone-200 hover:bg-stone-50 text-stone-600 font-bold rounded-xl text-sm transition-colors cursor-pointer"
                    >
                        Annuler
                    </button>
                    <button
                        type="button"
                        onclick={handleDeleteAccount}
                        disabled={deletingAccount}
                        class="px-5 py-2 bg-red-600 hover:bg-red-700 text-white font-bold rounded-xl text-sm transition-all flex items-center gap-2 cursor-pointer disabled:opacity-50"
                    >
                        {#if deletingAccount}
                            <Loader2 size={16} class="animate-spin" />
                            Suppression...
                        {:else}
                            Confirmer la suppression
                        {/if}
                    </button>
                </div>
                {#if deleteError}
                    <p class="text-xs text-red-600 bg-red-50 p-3 rounded-xl border border-red-100">{deleteError}</p>
                {/if}
            </div>
        </div>
    {/if}

    <!-- Fullscreen Image Preview Lightbox -->
    {#if previewImage}
        <div 
            class="fixed inset-0 bg-black/90 backdrop-blur-md flex flex-col items-center justify-center z-[9999] p-4 animate-in fade-in duration-200"
            onclick={() => { previewImage = null; }}
            role="button"
            tabindex="0"
            onkeydown={(e) => e.key === 'Escape' && (previewImage = null)}
        >
            <!-- Close Button -->
            <button 
                onclick={() => { previewImage = null; }}
                class="absolute top-6 right-6 p-2.5 bg-white/10 hover:bg-white/20 text-white rounded-full transition-all cursor-pointer border border-white/10"
                aria-label="Fermer l'aperçu"
            >
                <X size={24} />
            </button>

            <!-- Image Wrapper -->
            <div 
                class="relative max-w-5xl max-h-[80vh] flex items-center justify-center overflow-hidden rounded-2xl border border-white/10 shadow-2xl animate-in zoom-in-95 duration-250"
                onclick={(e) => e.stopPropagation()}
            >
                <img 
                    src={previewImage} 
                    alt={previewImageName} 
                    class="max-w-full max-h-[80vh] object-contain select-none"
                />
            </div>

            <!-- Image Info Footer -->
            <div 
                class="mt-6 text-center space-y-2.5 animate-in slide-in-from-bottom-4 duration-300"
                onclick={(e) => e.stopPropagation()}
            >
                <h4 class="text-white font-bold text-lg">{previewImageName}</h4>
                <div class="flex items-center gap-3 justify-center">
                    <button
                        onclick={() => {
                            navigator.clipboard.writeText(previewImage || '');
                            alert("URL copiée dans le presse-papiers !");
                        }}
                        class="px-4 py-2 bg-white hover:bg-stone-100 text-stone-900 rounded-xl text-xs font-bold transition-all flex items-center gap-1.5 cursor-pointer shadow-md"
                    >
                        <Link size={14} /> Copier le lien
                    </button>
                    <a
                        href={previewImage}
                        target="_blank"
                        rel="noopener noreferrer"
                        class="px-4 py-2 bg-stone-900 hover:bg-stone-800 text-white rounded-xl text-xs font-bold transition-all flex items-center gap-1.5 border border-white/15 cursor-pointer"
                    >
                        Ouvrir dans un nouvel onglet
                    </a>
                </div>
            </div>
        </div>
    {/if}
</div>
