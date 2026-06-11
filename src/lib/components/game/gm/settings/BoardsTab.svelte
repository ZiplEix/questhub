<script lang="ts">
    import { page } from "$app/state";
    import { boardsStore } from "$lib/websocket";
    import { addBoard, deleteBoard, activateBoard } from "$lib/api/board";
    import { uploadImage, validateImage } from "$lib/api/storage";
    import { 
        Plus, Trash2, Check, Image as ImageIcon, 
        Link as LinkIcon, Upload, Loader2, Layers,
        ChevronRight, Edit
    } from "lucide-svelte";
    import MediaSelectorModal from "$lib/components/MediaSelectorModal.svelte";

    // Props
    let {
        mode = "sidebar"
    } = $props<{
        mode?: "sidebar" | "settings";
    }>();

    let newBoardName = $state("");
    let newBoardUrl = $state("");
    let uploadMethod = $state<"file" | "url" | "media">("file");
    
    let isUploadingFile = $state(false);
    let isCreatingBoard = $state(false);
    let isMediaModalOpen = $state(false);
    
    let fileInput = $state<HTMLInputElement | null>(null);
    let dragOver = $state(false);

    const gameId = $derived(page.params.id || "");

    // --- Board Actions ---

    async function handleFileUpload(file: File) {
        if (isUploadingFile) return;

        const validation = validateImage(file, 10); // 10MB limit for maps/boards
        if (!validation.valid) {
            alert(validation.error);
            return;
        }

        try {
            isUploadingFile = true;
            const publicUrl = await uploadImage(file, 10);
            newBoardUrl = publicUrl;
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

    async function handleCreateBoard(e: Event) {
        e.preventDefault();
        if (!newBoardName.trim() || !newBoardUrl.trim() || isCreatingBoard) return;

        try {
            isCreatingBoard = true;
            await addBoard(gameId, newBoardName.trim(), newBoardUrl.trim());
            
            // Reset fields
            newBoardName = "";
            newBoardUrl = "";
        } catch (err) {
            console.error("Erreur lors de la création du plateau :", err);
        } finally {
            isCreatingBoard = false;
        }
    }

    async function handleActivateBoard(boardId: string) {
        try {
            await activateBoard(gameId, boardId);
        } catch (err) {
            console.error("Erreur lors de l'activation du plateau :", err);
        }
    }

    async function handleDeleteBoard(boardId: string) {
        if (!confirm("Voulez-vous vraiment supprimer ce plateau et sa carte associée ?")) return;
        try {
            await deleteBoard(boardId);
        } catch (err) {
            console.error("Erreur lors de la suppression du plateau :", err);
        }
    }
</script>

<div class="{mode === 'settings' ? 'space-y-8 animate-in fade-in duration-300' : 'h-full flex flex-col bg-stone-50 overflow-y-auto p-4 space-y-6'}">
    <!-- --- FORMULAIRE DE CRÉATION DE PLATEAU --- -->
    {#if mode === 'settings'}
        <div class="p-6 bg-stone-50 rounded-2xl border border-stone-200/80 shrink-0">
            <h3 class="font-bold text-dark-gray flex items-center gap-2 border-b border-stone-200 pb-3 mb-4">
                <Plus size={18} class="text-burnt-orange" />
                <span>Créer un nouveau plateau de jeu</span>
            </h3>

            <form onsubmit={handleCreateBoard} class="space-y-4">
                <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <div class="space-y-2">
                        <label for="board-name" class="text-sm font-bold text-dark-gray">Nom du plateau</label>
                        <input
                            id="board-name"
                            type="text"
                            placeholder="Ex: Donjon de l'effroi, Taverne..."
                            bind:value={newBoardName}
                            required
                            class="w-full px-4 py-2 bg-white border border-stone-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-burnt-orange/20 focus:border-burnt-orange transition-all"
                        />
                    </div>

                    <!-- Selector for Upload Method -->
                    <div class="space-y-2">
                        <span class="text-sm font-bold text-dark-gray block">Source de l'image</span>
                        <div class="flex rounded-lg bg-stone-200/50 p-0.5 gap-0.5 h-[42px] items-center">
                            <button
                                type="button"
                                onclick={() => { uploadMethod = "file"; newBoardUrl = ""; }}
                                class="flex-1 py-1.5 text-xs font-bold rounded-md transition-all flex items-center justify-center gap-1.5 h-full
                                {uploadMethod === 'file' ? 'bg-white text-stone-800 shadow-sm' : 'text-stone-500 hover:text-stone-700'}"
                            >
                                <Upload size={14} />
                                Fichier local
                            </button>
                            <button
                                type="button"
                                onclick={() => { uploadMethod = "url"; newBoardUrl = ""; }}
                                class="flex-1 py-1.5 text-xs font-bold rounded-md transition-all flex items-center justify-center gap-1.5 h-full
                                {uploadMethod === 'url' ? 'bg-white text-stone-800 shadow-sm' : 'text-stone-500 hover:text-stone-700'}"
                            >
                                <LinkIcon size={14} />
                                Lien externe URL
                            </button>
                            <button
                                type="button"
                                onclick={() => { uploadMethod = "media"; isMediaModalOpen = true; }}
                                class="flex-1 py-1.5 text-xs font-bold rounded-md transition-all flex items-center justify-center gap-1.5 h-full
                                {uploadMethod === 'media' ? 'bg-white text-stone-800 shadow-sm' : 'text-stone-500 hover:text-stone-700'}"
                            >
                                <ImageIcon size={14} />
                                Médiathèque
                            </button>
                        </div>
                    </div>
                </div>

                <!-- Image input -->
                {#if uploadMethod === 'file'}
                    <div
                        role="button"
                        tabindex="0"
                        ondragover={handleDragOver}
                        ondragleave={handleDragLeave}
                        ondrop={handleDrop}
                        onclick={() => fileInput?.click()}
                        onkeydown={(e) => e.key === "Enter" && fileInput?.click()}
                        class="border-2 border-dashed rounded-2xl p-6 text-center cursor-pointer transition-all flex flex-col items-center justify-center gap-3 bg-white
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
                            <span class="text-sm text-stone-500 font-medium">Téléversement de l'image de la carte...</span>
                        {:else if newBoardUrl}
                            <Check size={28} class="text-emerald-500" />
                            <span class="text-sm text-emerald-600 font-bold">L'image a été importée avec succès !</span>
                        {:else}
                            <Upload size={28} class="text-stone-400 animate-bounce" />
                            <span class="text-sm text-stone-500">
                                Glissez l'image de la carte ici ou <strong class="text-burnt-orange">parcourez vos fichiers</strong>
                            </span>
                        {/if}
                    </div>
                {:else if uploadMethod === 'url'}
                    <div class="space-y-2">
                        <label for="board-url" class="text-sm font-bold text-dark-gray">Adresse URL de l'image</label>
                        <input
                            id="board-url"
                            type="url"
                            placeholder="https://exemple.com/ma-carte.png"
                            bind:value={newBoardUrl}
                            required
                            class="w-full px-4 py-2 bg-white border border-stone-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-burnt-orange/20 focus:border-burnt-orange transition-all"
                        />
                    </div>
                {:else if uploadMethod === 'media'}
                    <div 
                        class="border-2 border-dashed rounded-2xl p-6 text-center cursor-pointer transition-all flex flex-col items-center justify-center gap-3 bg-white border-stone-200 hover:border-stone-300"
                        onclick={() => (isMediaModalOpen = true)}
                        role="button"
                        tabindex="0"
                        onkeydown={(e) => e.key === "Enter" && (isMediaModalOpen = true)}
                    >
                        {#if newBoardUrl}
                            <Check size={28} class="text-emerald-500" />
                            <span class="text-sm text-emerald-600 font-bold">Image sélectionnée depuis la médiathèque !</span>
                            <span class="text-xs text-stone-400 truncate max-w-full font-mono">{newBoardUrl}</span>
                        {:else}
                            <ImageIcon size={28} class="text-stone-400" />
                            <span class="text-sm text-stone-500">
                                Cliquer pour choisir une image dans la <strong class="text-burnt-orange">médiathèque</strong>
                            </span>
                        {/if}
                    </div>
                {/if}

                <div class="flex justify-end pt-2">
                    <button
                        type="submit"
                        disabled={isCreatingBoard || isUploadingFile || !newBoardName.trim() || !newBoardUrl.trim()}
                        class="px-6 py-2.5 bg-burnt-orange hover:bg-burnt-orange-hover text-white rounded-xl text-sm font-bold flex items-center justify-center gap-2 transition-all disabled:opacity-50 hover:-translate-y-0.5 shadow-md"
                    >
                        {#if isCreatingBoard}
                            <Loader2 size={16} class="animate-spin" />
                            Création en cours...
                        {:else}
                            <Plus size={16} />
                            Créer le plateau
                        {/if}
                    </button>
                </div>
            </form>
        </div>
    {/if}

    <!-- --- LISTE DES PLATEAUX --- -->
    <div class="space-y-4">
        <h3 class="font-bold text-dark-gray flex items-center gap-2 border-b border-stone-100 pb-2">
            <Layers size={18} class="text-burnt-orange" />
            <span>Plateaux existants</span>
        </h3>

        {#if $boardsStore.length === 0}
            <div class="text-sm text-stone-400 py-12 text-center bg-stone-50 rounded-2xl border border-dashed border-stone-200">
                Aucun plateau n'a encore été créé. Utilisez le formulaire ci-dessus pour ajouter votre premier plateau de jeu.
            </div>
        {:else}
            <div class={mode === 'settings' ? "grid grid-cols-1 md:grid-cols-2 gap-4" : "flex flex-col gap-2.5"}>
                {#each $boardsStore as board (board.id)}
                    {#if mode === 'sidebar'}
                        <!-- --- SIDEBAR VIEW (COMPACT & PREMIUM) --- -->
                        <button 
                            onclick={() => !board.is_active && handleActivateBoard(board.id)}
                            disabled={board.is_active}
                            class="group relative rounded-xl border flex items-center gap-3 p-2 bg-white transition-all duration-200 w-full text-left focus:outline-none focus:ring-2 focus:ring-emerald-500/20
                            {board.is_active 
                                ? 'border-emerald-500/40 bg-emerald-50/15 shadow-sm cursor-default' 
                                : 'border-stone-200/80 hover:border-stone-300 hover:bg-stone-50/50 hover:shadow-sm cursor-pointer'}"
                        >
                            <!-- Thumbnail -->
                            <div class="relative w-16 h-12 bg-stone-100 rounded-lg overflow-hidden shrink-0 border border-stone-200/60 shadow-inner">
                                {#if board.image_url}
                                    <img 
                                        src={board.image_url} 
                                        alt={board.name}
                                        class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-300"
                                    />
                                {:else}
                                    <div class="w-full h-full flex items-center justify-center text-stone-400">
                                        <ImageIcon size={16} />
                                    </div>
                                {/if}
                                {#if board.is_active}
                                    <div class="absolute inset-0 bg-emerald-500/10 pointer-events-none"></div>
                                    <div class="absolute top-1 left-1 bg-emerald-500 text-white p-0.5 rounded-full shadow-sm">
                                        <Check size={8} />
                                    </div>
                                {/if}
                            </div>

                            <!-- Details -->
                            <div class="flex-1 min-w-0 flex flex-col justify-center">
                                <span class="text-sm font-bold text-stone-800 truncate" title={board.name}>
                                    {board.name}
                                </span>
                                <span class="text-[11px] font-medium flex items-center gap-1 mt-0.5">
                                    {#if board.is_active}
                                        <span class="inline-block w-1.5 h-1.5 rounded-full bg-emerald-500 animate-pulse"></span>
                                        <span class="text-emerald-700 font-semibold">Actif en jeu</span>
                                    {:else}
                                        <span class="inline-block w-1.5 h-1.5 rounded-full bg-stone-300"></span>
                                        <span class="text-stone-400 group-hover:text-stone-500 transition-colors">Cliquez pour activer</span>
                                    {/if}
                                </span>
                            </div>

                            <!-- Actions -->
                            <div class="shrink-0 flex items-center pr-1">
                                {#if board.is_active}
                                    <Check size={16} class="text-emerald-600 font-bold" />
                                {:else}
                                    <ChevronRight size={16} class="text-stone-300 opacity-0 group-hover:opacity-100 transition-opacity duration-200" />
                                {/if}
                            </div>
                        </button>
                    {:else}
                        <!-- --- SETTINGS VIEW (DETAILED GRID) --- -->
                        <div class="group relative rounded-2xl border overflow-hidden flex bg-white border-stone-200 hover:shadow-md hover:border-stone-300 transition-all">
                            <!-- Preview Thumbnail -->
                            <div class="relative w-36 h-24 bg-stone-50 overflow-hidden shrink-0 border-r border-stone-100">
                                {#if board.image_url}
                                    <img 
                                        src={board.image_url} 
                                        alt={board.name}
                                        class="w-full h-full object-cover"
                                    />
                                {:else}
                                    <div class="w-full h-full flex items-center justify-center text-stone-300">
                                        <ImageIcon size={28} />
                                    </div>
                                {/if}
                                {#if board.is_active}
                                    <div class="absolute inset-0 bg-emerald-500/10 border-2 border-emerald-500 pointer-events-none"></div>
                                    <div class="absolute top-2 left-2 bg-emerald-500 text-white p-0.5 rounded-full shadow-sm">
                                        <Check size={12} />
                                    </div>
                                {/if}
                            </div>
                            
                            <!-- Card Body -->
                            <div class="p-4 flex flex-col justify-between flex-1 min-w-0">
                                <div class="flex items-start justify-between gap-2">
                                    <span class="text-sm font-bold text-stone-800 truncate" title={board.name}>{board.name}</span>
                                    {#if board.is_active}
                                        <span class="text-[9px] bg-emerald-100 text-emerald-800 font-extrabold px-2.5 py-0.5 rounded-full shrink-0">Actif</span>
                                    {/if}
                                </div>
                                <div class="flex items-center gap-2 mt-2">
                                    {#if !board.is_active}
                                        <button
                                            onclick={() => handleActivateBoard(board.id)}
                                            class="flex-1 py-1.5 px-3 bg-emerald-50 hover:bg-emerald-100 text-emerald-700 hover:text-emerald-800 rounded-xl text-xs font-bold transition-all flex items-center justify-center gap-1"
                                        >
                                            <Check size={12} /> Activer
                                        </button>
                                    {:else}
                                        <span class="flex-1 text-center text-xs font-bold text-emerald-600 bg-emerald-50/50 py-1.5 rounded-xl">Affiché en jeu</span>
                                    {/if}
                                    <a
                                        href="/table/{gameId}/gm/settings/board/{board.id}"
                                        class="p-1.5 text-stone-400 hover:text-burnt-orange hover:bg-stone-50 rounded-xl transition-colors border border-stone-100 hover:border-stone-200 shrink-0 flex items-center justify-center cursor-pointer"
                                        title="Modifier le plateau"
                                    >
                                        <Edit size={16} />
                                    </a>
                                    <button
                                        onclick={() => handleDeleteBoard(board.id)}
                                        class="p-1.5 text-stone-400 hover:text-red-500 hover:bg-red-50 rounded-xl transition-colors border border-stone-100 hover:border-red-100 shrink-0"
                                        title="Supprimer le plateau"
                                    >
                                        <Trash2 size={16} />
                                    </button>
                                </div>
                            </div>
                        </div>
                    {/if}
                {/each}
            </div>
        {/if}
    </div>

    <MediaSelectorModal
        isOpen={isMediaModalOpen}
        onSelect={(url) => {
            newBoardUrl = url;
            isMediaModalOpen = false;
        }}
        onClose={() => {
            isMediaModalOpen = false;
        }}
    />
</div>
