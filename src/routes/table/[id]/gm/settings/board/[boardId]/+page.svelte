<script lang="ts">
    import { page } from "$app/state";
    import { goto } from "$app/navigation";
    import Header from "$lib/components/Header.svelte";
    import { onMount } from "svelte";
    import { fetchBoard, updateBoard } from "$lib/api/board";
    import { uploadImage } from "$lib/api/storage";
    import { authClient } from "$lib/auth-client";
    import { 
        ArrowLeft, Check, Upload, Link as LinkIcon, 
        Loader2, Image as ImageIcon, Layers 
    } from "lucide-svelte";

    const gameId = page.params.id as string;
    const boardId = page.params.boardId as string;

    let boardName = $state("");
    let boardUrl = $state("");
    let uploadMethod = $state<"file" | "url">("file");
    
    let loading = $state(true);
    let saving = $state(false);
    let isUploadingFile = $state(false);
    let error = $state<string | null>(null);

    let fileInput = $state<HTMLInputElement | null>(null);
    let dragOver = $state(false);

    onMount(async () => {
        try {
            const { data: tokenData } = await authClient.token();
            if (tokenData?.token) {
                const boardData = await fetchBoard(boardId);
                if (boardData) {
                    boardName = boardData.name;
                    boardUrl = boardData.image_url || "";
                    if (boardUrl && !boardUrl.startsWith("http")) {
                        uploadMethod = "file";
                    } else if (boardUrl) {
                        // Default to URL if it looks like an absolute external URL, or fallback to file
                        uploadMethod = boardUrl.startsWith("http") && !boardUrl.includes("supabase") ? "url" : "file";
                    }
                } else {
                    error = "Plateau introuvable.";
                }
            }
        } catch (e) {
            console.error("Failed to fetch board details:", e);
            error = "Impossible de charger les informations du plateau.";
        } finally {
            loading = false;
        }
    });

    async function handleFileUpload(file: File) {
        if (isUploadingFile) return;
        try {
            isUploadingFile = true;
            const publicUrl = await uploadImage(file);
            boardUrl = publicUrl;
        } catch (err) {
            console.error("Erreur lors de l'upload de l'image :", err);
            alert("Erreur lors de l'upload de l'image.");
        } finally {
            isUploadingFile = false;
        }
    }

    function onFileSelected(e: Event) {
        const files = (e.target as HTMLInputElement).files;
        if (files && files.length > 0) {
            handleFileUpload(files[0]);
        }
    }

    function handleDragOver(e: DragEvent) {
        e.preventDefault();
        dragOver = true;
    }

    function handleDragLeave() {
        dragOver = false;
    }

    function handleDrop(e: DragEvent) {
        e.preventDefault();
        dragOver = false;
        const files = e.dataTransfer?.files;
        if (files && files.length > 0) {
            handleFileUpload(files[0]);
        }
    }

    async function handleSave() {
        if (!boardName.trim() || !boardUrl.trim()) return;
        saving = true;
        try {
            await updateBoard(boardId, boardName.trim(), boardUrl.trim());
            goto(`/table/${gameId}/gm/settings?tab=boards`);
        } catch (e) {
            console.error("Failed to save board settings:", e);
            alert("Erreur lors de la sauvegarde du plateau.");
        } finally {
            saving = false;
        }
    }
</script>

<div class="min-h-screen bg-cream">
    <Header />

    <main class="max-w-2xl mx-auto p-6 md:p-12 animate-in fade-in duration-300">
        <!-- Back Button -->
        <a
            href="/table/{gameId}/gm/settings?tab=boards"
            class="inline-flex items-center gap-2 text-stone-500 hover:text-dark-gray mb-6 font-medium transition-colors"
        >
            <ArrowLeft size={16} />
            Retour aux paramètres
        </a>

        <!-- Header -->
        <div class="mb-8">
            <h1 class="text-3xl font-display font-bold text-dark-gray mb-2">
                Modifier le plateau
            </h1>
            <p class="text-stone-500">
                Ajustez le nom et l'image de carte pour ce plateau de jeu.
            </p>
        </div>

        {#if loading}
            <div class="flex justify-center py-12">
                <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-burnt-orange"></div>
            </div>
        {:else if error}
            <div class="p-4 bg-red-50 text-red-600 rounded-xl border border-red-100 shadow-sm">
                {error}
            </div>
        {:else}
            <div class="bg-white rounded-2xl shadow-sm border border-stone-100 p-6 md:p-8 space-y-6">
                <!-- Header Icon & Title -->
                <div class="flex items-center gap-4 pb-6 border-b border-stone-100">
                    <div class="w-12 h-12 rounded-full bg-burnt-orange/10 flex items-center justify-center text-burnt-orange">
                        <Layers size={24} />
                    </div>
                    <div>
                        <h2 class="text-xl font-bold text-dark-gray">Paramètres de la carte</h2>
                        <p class="text-sm text-stone-400">Plateau ID: {boardId}</p>
                    </div>
                </div>

                <!-- Form Section -->
                <div class="space-y-4">
                    <div class="space-y-2">
                        <label for="board-name" class="text-sm font-bold text-dark-gray">Nom du plateau</label>
                        <input
                            id="board-name"
                            type="text"
                            placeholder="Ex: Donjon de l'effroi..."
                            bind:value={boardName}
                            required
                            class="w-full px-4 py-2.5 bg-stone-50 border border-stone-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-burnt-orange/20 focus:border-burnt-orange transition-all font-medium text-dark-gray"
                        />
                    </div>

                    <!-- Image Source Selection Swatches -->
                    <div class="space-y-2">
                        <span class="text-sm font-bold text-dark-gray block">Source de l'image</span>
                        <div class="flex rounded-lg bg-stone-100 p-0.5 gap-0.5 h-[42px] items-center max-w-sm">
                            <button
                                type="button"
                                onclick={() => { uploadMethod = "file"; }}
                                class="flex-1 py-1.5 text-xs font-bold rounded-md transition-all flex items-center justify-center gap-1.5 h-full cursor-pointer
                                {uploadMethod === 'file' ? 'bg-white text-stone-800 shadow-sm' : 'text-stone-500 hover:text-stone-700'}"
                            >
                                <Upload size={14} />
                                Fichier local
                            </button>
                            <button
                                type="button"
                                onclick={() => { uploadMethod = "url"; }}
                                class="flex-1 py-1.5 text-xs font-bold rounded-md transition-all flex items-center justify-center gap-1.5 h-full cursor-pointer
                                {uploadMethod === 'url' ? 'bg-white text-stone-800 shadow-sm' : 'text-stone-500 hover:text-stone-700'}"
                            >
                                <LinkIcon size={14} />
                                Lien externe URL
                            </button>
                        </div>
                    </div>

                    <!-- Map Image Upload/Url Area -->
                    {#if uploadMethod === 'file'}
                        <div
                            role="button"
                            tabindex="0"
                            ondragover={handleDragOver}
                            ondragleave={handleDragLeave}
                            ondrop={handleDrop}
                            onclick={() => fileInput?.click()}
                            onkeydown={(e) => e.key === "Enter" && fileInput?.click()}
                            class="border-2 border-dashed rounded-2xl p-8 text-center cursor-pointer transition-all flex flex-col items-center justify-center gap-3 bg-stone-50
                            {dragOver ? 'border-burnt-orange bg-burnt-orange/5' : 'border-stone-200 hover:border-stone-300'}"
                        >
                            <input
                                type="file"
                                accept="image/*"
                                bind:this={fileInput}
                                onchange={onFileSelected}
                                class="hidden"
                            />
                            {#if isUploadingFile}
                                <Loader2 size={28} class="animate-spin text-burnt-orange" />
                                <span class="text-sm text-stone-500 font-medium">Téléversement de la nouvelle carte...</span>
                            {:else if boardUrl}
                                <Check size={28} class="text-emerald-500" />
                                <span class="text-sm text-emerald-600 font-bold">Image chargée avec succès !</span>
                            {:else}
                                <Upload size={28} class="text-stone-400" />
                                <span class="text-sm text-stone-500">
                                    Glissez la nouvelle image ici ou <strong class="text-burnt-orange">parcourez vos fichiers</strong>
                                </span>
                            {/if}
                        </div>
                    {:else}
                        <div class="space-y-2">
                            <label for="board-url" class="text-sm font-bold text-dark-gray">Adresse URL de la carte</label>
                            <input
                                id="board-url"
                                type="url"
                                placeholder="https://exemple.com/carte.png"
                                bind:value={boardUrl}
                                required
                                class="w-full px-4 py-2.5 bg-stone-50 border border-stone-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-burnt-orange/20 focus:border-burnt-orange transition-all font-medium"
                            />
                        </div>
                    {/if}

                    <!-- Preview Section -->
                    {#if boardUrl}
                        <div class="space-y-2 pt-2">
                            <span class="text-sm font-bold text-dark-gray block">Aperçu de la carte</span>
                            <div class="relative rounded-2xl overflow-hidden border border-stone-200/80 shadow-md bg-stone-100 max-h-64 flex items-center justify-center">
                                <img 
                                    src={boardUrl} 
                                    alt="Aperçu du plateau"
                                    class="w-full h-full object-contain max-h-60"
                                />
                            </div>
                        </div>
                    {/if}
                </div>

                <!-- Form Actions -->
                <div class="pt-6 border-t border-stone-100 flex justify-end gap-3">
                    <a
                        href="/table/{gameId}/gm/settings?tab=boards"
                        class="px-5 py-2.5 bg-white border border-stone-200 hover:bg-stone-50 text-stone-600 rounded-xl font-medium transition-all text-sm"
                    >
                        Annuler
                    </a>
                    <button
                        onclick={handleSave}
                        disabled={saving || isUploadingFile || !boardName.trim() || !boardUrl.trim()}
                        class="px-5 py-2.5 bg-burnt-orange hover:bg-burnt-orange-hover text-white rounded-xl font-medium disabled:opacity-50 transition-all text-sm shadow-sm flex items-center justify-center min-w-[100px]"
                    >
                        {#if saving}
                            <div class="animate-spin rounded-full h-4 w-4 border-b-2 border-white"></div>
                        {:else}
                            Enregistrer
                        {/if}
                    </button>
                </div>
            </div>
        {/if}
    </main>
</div>
