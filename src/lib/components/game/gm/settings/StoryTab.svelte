<script lang="ts">
    import { onMount } from "svelte";
    import { 
        Folder, 
        FolderPlus, 
        FileText, 
        FilePlus, 
        Eye, 
        EyeOff, 
        Trash2, 
        Edit, 
        ChevronRight, 
        ChevronDown, 
        Loader2, 
        CheckCircle, 
        AlertCircle, 
        Plus, 
        Check, 
        X,
        BookOpen,
        Image as ImageIcon,
        UploadCloud,
        Copy,
        FileImage,
        FolderHeart
    } from "lucide-svelte";
    import MediaSelectorModal from "$lib/components/MediaSelectorModal.svelte";
    import { marked } from "marked";
    import { 
        fetchStoryFolders, 
        fetchStoryPages, 
        createStoryFolder, 
        updateStoryFolder, 
        deleteStoryFolder, 
        createStoryPage, 
        updateStoryPage, 
        deleteStoryPage,
        type StoryFolder,
        type StoryPage
    } from "$lib/api/story";
    import {
        uploadStoryAsset,
        listStoryAssets,
        deleteStoryAsset,
        validateImage
    } from "$lib/api/storage";

    let { gameId } = $props<{ gameId: string }>();

    let folders = $state<StoryFolder[]>([]);
    let pages = $state<StoryPage[]>([]);
    let loading = $state(true);
    let selectedPage = $state<StoryPage | null>(null);
    let activeEditorTab = $state<"edit" | "preview">("edit");

    // Collapsed state for folders in the tree
    let collapsedFolders = $state<Record<string, boolean>>({});

    // Inline Creation forms state
    let isCreatingFolder = $state(false);
    let newFolderName = $state("");
    let isCreatingPage = $state(false);
    let targetFolderId = $state<string | null>(null);

    // Folder Editing / Renaming
    let editingFolderId = $state<string | null>(null);
    let editingFolderName = $state("");

    // Auto-save status
    let saveStatus = $state<"saved" | "saving" | "error">("saved");
    let saveTimeout: ReturnType<typeof setTimeout>;

    // Drag and Drop state
    let draggedPageId = $state<string | null>(null);
    let activeDropTargetId = $state<string | "root" | null>(null);

    // Media Library (Drive) state
    let assets = $state<Array<{ name: string; url: string; created_at: string; size: number }>>([]);
    let assetsLoading = $state(true);
    let assetsCollapsed = $state(true);
    let isUploadingAsset = $state(false);
    let isMediaModalOpen = $state(false);
    let isImportingFromMedia = $state(false);
    let assetsError = $state<string | null>(null);
    let fileInputRef = $state<HTMLInputElement | null>(null);
    let textareaElement = $state<HTMLTextAreaElement | null>(null);

    function handleDragStart(e: DragEvent, page: StoryPage) {
        draggedPageId = page.id;
        if (e.dataTransfer) {
            e.dataTransfer.effectAllowed = "move";
            e.dataTransfer.setData("text/plain", page.id);
        }
    }

    function handleAssetDragStart(e: DragEvent, asset: typeof assets[0]) {
        if (e.dataTransfer) {
            const altText = asset.name.split('_').slice(1).join('_').split('.').shift() || "image";
            const markdown = `![${altText}](${asset.url})`;
            e.dataTransfer.setData("text/plain", markdown);
            e.dataTransfer.effectAllowed = "copy";
        }
    }

    function handleDragOver(e: DragEvent, targetId: string | null) {
        e.preventDefault();
        activeDropTargetId = targetId === null ? "root" : targetId;
    }

    function handleDragLeave() {
        activeDropTargetId = null;
    }

    async function handleDrop(e: DragEvent, targetFolderId: string | null) {
        e.preventDefault();
        const pageId = e.dataTransfer?.getData("text/plain") || draggedPageId;
        draggedPageId = null;
        activeDropTargetId = null;

        if (!pageId) return;

        const page = pages.find(p => p.id === pageId);
        if (!page) return;
        if (page.folder_id === targetFolderId) return;

        try {
            const updated = await updateStoryPage(pageId, { folder_id: targetFolderId });
            pages = pages.map(p => p.id === pageId ? updated : p);
            if (selectedPage && selectedPage.id === pageId) {
                selectedPage.folder_id = targetFolderId;
            }
        } catch (error) {
            console.error("Failed to drag and drop page:", error);
        }
    }

    async function loadAssets() {
        try {
            assets = await listStoryAssets(gameId);
        } catch (error) {
            console.error("Failed to load story assets:", error);
            assetsError = "Impossible de charger les images.";
        }
    }

    async function handleUploadFiles(files: FileList | File[]) {
        if (files.length === 0) return;
        
        const imageFiles: File[] = [];
        for (const file of Array.from(files)) {
            const validation = validateImage(file, 2);
            if (!validation.valid) {
                alert(`Fichier "${file.name}" rejeté : ${validation.error}`);
            } else {
                imageFiles.push(file);
            }
        }
        if (imageFiles.length === 0) return;

        isUploadingAsset = true;
        saveStatus = "saving";
        assetsError = null;

        try {
            for (const file of imageFiles) {
                const url = await uploadStoryAsset(gameId, file);
                if (selectedPage) {
                    const altText = file.name.split('.').shift() || "image";
                    insertMarkdown(`\n![${altText}](${url})\n`);
                }
            }
            await loadAssets();
            saveStatus = "saved";
        } catch (error) {
            console.error("Failed to upload files:", error);
            saveStatus = "error";
            assetsError = "Échec du téléversement.";
        } finally {
            isUploadingAsset = false;
        }
    }

    async function handleDeleteAsset(fileName: string) {
        if (!confirm("Voulez-vous vraiment supprimer cette image du Drive de la campagne ? Elle ne s'affichera plus dans le lore.")) return;
        try {
            await deleteStoryAsset(gameId, fileName);
            assets = assets.filter(a => a.name !== fileName);
        } catch (error) {
            console.error("Failed to delete asset:", error);
            alert("Erreur lors de la suppression de l'image.");
        }
    }

    async function handleSelectFromMediaLibrary(url: string) {
        isMediaModalOpen = false;
        isImportingFromMedia = true;
        saveStatus = "saving";
        assetsError = null;

        try {
            const response = await fetch(url);
            if (!response.ok) throw new Error("Impossible de télécharger l'image depuis la médiathèque.");
            
            const blob = await response.blob();
            const urlPath = new URL(url).pathname;
            const originalName = urlPath.split('/').pop() || 'imported_image.png';
            const nameWithoutTimestamp = originalName.replace(/^\d+_(.+)/, '$1');
            const file = new File([blob], nameWithoutTimestamp, { type: blob.type });

            const storyAssetUrl = await uploadStoryAsset(gameId, file);
            
            if (selectedPage) {
                const altText = nameWithoutTimestamp.split('.').shift() || "image";
                insertMarkdown(`\n![${altText}](${storyAssetUrl})\n`);
            }
            await loadAssets();
            saveStatus = "saved";
        } catch (error: any) {
            console.error("Failed to import image from media library:", error);
            saveStatus = "error";
            assetsError = error.message || "Échec de l'importation de l'image.";
            alert(`Erreur lors de l'importation : ${assetsError}`);
        } finally {
            isImportingFromMedia = false;
        }
    }

    function insertMarkdown(text: string) {
        if (!selectedPage || !textareaElement) return;
        
        const start = textareaElement.selectionStart;
        const end = textareaElement.selectionEnd;
        const currentContent = selectedPage.content || "";
        
        selectedPage.content = currentContent.substring(0, start) + text + currentContent.substring(end);
        
        setTimeout(() => {
            if (textareaElement) {
                textareaElement.focus();
                const newCursorPos = start + text.length;
                textareaElement.setSelectionRange(newCursorPos, newCursorPos);
            }
        }, 50);

        triggerAutoSave();
    }

    function handleTextareaDrop(e: DragEvent) {
        if (e.dataTransfer && e.dataTransfer.files.length > 0) {
            const files = Array.from(e.dataTransfer.files);
            const containsImages = files.some(file => file.type.startsWith("image/"));
            if (containsImages) {
                e.preventDefault();
                handleUploadFiles(files);
            }
        }
    }

    function handleTextareaPaste(e: ClipboardEvent) {
        if (e.clipboardData && e.clipboardData.files.length > 0) {
            const files = Array.from(e.clipboardData.files);
            const containsImages = files.some(file => file.type.startsWith("image/"));
            if (containsImages) {
                e.preventDefault();
                handleUploadFiles(files);
            }
        }
    }

    onMount(async () => {
        try {
            await Promise.all([
                refreshData(),
                loadAssets()
            ]);
        } catch (error) {
            console.error("Failed to load story elements:", error);
        } finally {
            loading = false;
            assetsLoading = false;
        }
    });

    async function refreshData() {
        const [foldersData, pagesData] = await Promise.all([
            fetchStoryFolders(gameId),
            fetchStoryPages(gameId)
        ]);
        folders = foldersData;
        pages = pagesData;
    }

    // FOLDERS CRUD
    async function handleCreateFolder() {
        if (!newFolderName.trim()) return;
        try {
            const folder = await createStoryFolder(gameId, newFolderName.trim());
            folders = [...folders, folder];
            newFolderName = "";
            isCreatingFolder = false;
        } catch (error) {
            console.error("Failed to create folder:", error);
        }
    }

    async function handleRenameFolder(folderId: string) {
        if (!editingFolderName.trim()) return;
        try {
            const updated = await updateStoryFolder(folderId, { name: editingFolderName.trim() });
            folders = folders.map(f => f.id === folderId ? updated : f);
            editingFolderId = null;
            editingFolderName = "";
        } catch (error) {
            console.error("Failed to rename folder:", error);
        }
    }

    async function handleToggleFolderVisibility(folder: StoryFolder) {
        try {
            const nextVisibility = !folder.is_visible;
            const updated = await updateStoryFolder(folder.id, { is_visible: nextVisibility });
            folders = folders.map(f => f.id === folder.id ? updated : f);
            
            // Re-fetch pages to update their inherited states locally, or just let realtime subscription handle it if any.
            // But immediate local sync is better:
            await refreshData();
        } catch (error) {
            console.error("Failed to toggle folder visibility:", error);
        }
    }

    async function handleDeleteFolder(folderId: string) {
        if (!confirm("Voulez-vous vraiment supprimer ce dossier ? Tous ses articles seront supprimés ou détachés selon la base de données (cascade activée).")) return;
        try {
            await deleteStoryFolder(folderId);
            folders = folders.filter(f => f.id !== folderId);
            pages = pages.filter(p => p.folder_id !== folderId);
            if (selectedPage && selectedPage.folder_id === folderId) {
                selectedPage = null;
            }
        } catch (error) {
            console.error("Failed to delete folder:", error);
        }
    }

    // PAGES CRUD
    async function handleCreatePage(folderId: string | null = null) {
        const defaultTitle = "Nouvel Article";
        try {
            const page = await createStoryPage(gameId, folderId, defaultTitle, "");
            pages = [...pages, page];
            selectedPage = page;
            activeEditorTab = "edit";
        } catch (error) {
            console.error("Failed to create page:", error);
        }
    }

    async function handleTogglePageVisibility(page: StoryPage) {
        try {
            const nextVisibility = !page.is_visible;
            const updated = await updateStoryPage(page.id, { is_visible: nextVisibility });
            pages = pages.map(p => p.id === page.id ? updated : p);
            if (selectedPage && selectedPage.id === page.id) {
                selectedPage.is_visible = nextVisibility;
            }
        } catch (error) {
            console.error("Failed to toggle page visibility:", error);
        }
    }

    async function handlePageFolderChange(folderId: string | "" | null) {
        if (!selectedPage) return;
        const targetFolder = folderId === "" ? null : folderId;
        try {
            const updated = await updateStoryPage(selectedPage.id, { folder_id: targetFolder });
            pages = pages.map(p => p.id === selectedPage!.id ? updated : p);
            selectedPage.folder_id = targetFolder;
        } catch (error) {
            console.error("Failed to move page folder:", error);
        }
    }

    async function handleDeletePage(pageId: string) {
        if (!confirm("Voulez-vous vraiment supprimer cet article ?")) return;
        try {
            await deleteStoryPage(pageId);
            pages = pages.filter(p => p.id !== pageId);
            if (selectedPage && selectedPage.id === pageId) {
                selectedPage = null;
            }
        } catch (error) {
            console.error("Failed to delete page:", error);
        }
    }

    // EDITOR AUTO-SAVE
    function triggerAutoSave() {
        if (!selectedPage) return;
        saveStatus = "saving";
        clearTimeout(saveTimeout);
        saveTimeout = setTimeout(async () => {
            if (!selectedPage) return;
            try {
                const updated = await updateStoryPage(selectedPage.id, {
                    title: selectedPage.title,
                    content: selectedPage.content
                });
                pages = pages.map(p => p.id === selectedPage!.id ? updated : p);
                saveStatus = "saved";
            } catch (error) {
                console.error("Auto-save failed:", error);
                saveStatus = "error";
            }
        }, 1000);
    }

    function selectPage(page: StoryPage) {
        // Clear pending auto-saves of previous page if any
        if (saveStatus === "saving" && selectedPage) {
            clearTimeout(saveTimeout);
            // Save immediately
            updateStoryPage(selectedPage.id, {
                title: selectedPage.title,
                content: selectedPage.content
            }).then(updated => {
                pages = pages.map(p => p.id === updated.id ? updated : p);
            });
        }
        selectedPage = { ...page };
        saveStatus = "saved";
    }

    // Markdown derived HTML
    let previewHtml = $derived.by(() => {
        if (!selectedPage) return "";
        try {
            return marked.parse(selectedPage.content || "") as string;
        } catch (e) {
            console.error("Markdown parse error:", e);
            return "<p class='text-red-500'>Erreur de rendu Markdown.</p>";
        }
    });
</script>

<div class="h-[650px] flex gap-6 text-stone-700 bg-white">
    <!-- LEFT SIDEBAR: Story Tree -->
    <div class="w-80 flex flex-col border-r border-stone-100 pr-4 h-full overflow-y-auto">
        <!-- Sidebar controls -->
        <div class="flex gap-2 mb-4">
            <button
                onclick={() => {
                    isCreatingFolder = !isCreatingFolder;
                    isCreatingPage = false;
                }}
                class="flex-1 flex items-center justify-center gap-1.5 py-2 px-3 bg-stone-50 hover:bg-stone-100 border border-stone-200 text-stone-600 rounded-xl text-xs font-bold transition-all"
            >
                <FolderPlus size={14} />
                <span>Nouveau Dossier</span>
            </button>
            <button
                onclick={() => handleCreatePage(null)}
                class="flex-1 flex items-center justify-center gap-1.5 py-2 px-3 bg-burnt-orange hover:bg-burnt-orange-dark text-white rounded-xl text-xs font-bold transition-all"
            >
                <FilePlus size={14} />
                <span>Nouvel Article</span>
            </button>
        </div>

        <!-- Inline Folder Creation Input -->
        {#if isCreatingFolder}
            <div class="mb-4 p-3 bg-stone-50 rounded-xl border border-stone-100 flex gap-2">
                <input
                    type="text"
                    bind:value={newFolderName}
                    placeholder="Nom du dossier..."
                    class="flex-1 px-3 py-1.5 text-xs border border-stone-200 rounded-lg focus:outline-none focus:border-burnt-orange bg-white"
                    onkeydown={(e) => e.key === "Enter" && handleCreateFolder()}
                />
                <button
                    onclick={handleCreateFolder}
                    class="p-1.5 bg-green-500 hover:bg-green-600 text-white rounded-lg transition-colors"
                >
                    <Check size={14} />
                </button>
                <button
                    onclick={() => {
                        isCreatingFolder = false;
                        newFolderName = "";
                    }}
                    class="p-1.5 bg-stone-200 hover:bg-stone-300 text-stone-500 rounded-lg transition-colors"
                >
                    <X size={14} />
                </button>
            </div>
        {/if}

        <!-- List Tree -->
        {#if loading}
            <div class="flex-1 flex items-center justify-center">
                <Loader2 class="animate-spin text-stone-400" />
            </div>
        {:else}
            <div class="flex-1 space-y-4">
                <!-- FOLDERS -->
                <div class="space-y-1">
                    <h4 class="text-xs font-bold uppercase tracking-wider text-stone-400 mb-2 px-2">Dossiers</h4>
                    {#if folders.length === 0}
                        <p class="text-xs text-stone-400 italic px-2">Aucun dossier créé.</p>
                    {/if}
                    {#each folders as folder}
                        <div 
                            ondragover={(e) => handleDragOver(e, folder.id)}
                            ondragleave={handleDragLeave}
                            ondrop={(e) => handleDrop(e, folder.id)}
                            class="group rounded-xl transition-all border border-transparent 
                            {editingFolderId === folder.id ? 'bg-stone-50 border-stone-200 p-2' : ''}
                            {activeDropTargetId === folder.id ? 'bg-burnt-orange/10 border-burnt-orange border-dashed' : ''}"
                        >
                            <!-- Folder Header Row -->
                            <div class="flex items-center justify-between p-2 rounded-lg hover:bg-stone-50 transition-colors">
                                <button
                                    onclick={() => collapsedFolders[folder.id] = !collapsedFolders[folder.id]}
                                    class="flex items-center gap-2 text-stone-700 font-semibold text-sm flex-1 text-left"
                                >
                                    {#if collapsedFolders[folder.id]}
                                        <ChevronRight size={16} class="text-stone-400" />
                                    {:else}
                                        <ChevronDown size={16} class="text-stone-400" />
                                    {/if}
                                    
                                    <Folder size={16} class="text-burnt-orange fill-burnt-orange/10" />
                                    
                                    {#if editingFolderId === folder.id}
                                        <input
                                            type="text"
                                            bind:value={editingFolderName}
                                            class="px-2 py-0.5 text-sm border border-stone-200 rounded bg-white font-normal"
                                            onclick={(e) => e.stopPropagation()}
                                            onkeydown={(e) => e.key === "Enter" && handleRenameFolder(folder.id)}
                                        />
                                    {:else}
                                        <span class="truncate">{folder.name}</span>
                                    {/if}
                                </button>

                                <div class="flex items-center gap-1.5">
                                    <!-- Always visible folder actions -->
                                    {#if editingFolderId !== folder.id}
                                        <!-- Visibility Toggle -->
                                        <button 
                                            onclick={(e) => {
                                                e.stopPropagation();
                                                handleToggleFolderVisibility(folder);
                                            }} 
                                            class="p-1 text-stone-400 hover:text-dark-gray transition-colors" 
                                            title={folder.is_visible ? "Visible pour les joueurs" : "Masqué pour les joueurs"}
                                        >
                                            {#if folder.is_visible}
                                                <Eye size={14} class="text-green-500" />
                                            {:else}
                                                <EyeOff size={14} />
                                            {/if}
                                        </button>

                                        <!-- Rename Folder Toggle -->
                                        <button 
                                            onclick={(e) => {
                                                e.stopPropagation();
                                                editingFolderId = folder.id;
                                                editingFolderName = folder.name;
                                            }} 
                                            class="p-1 text-stone-400 hover:text-dark-gray transition-colors" 
                                            title="Renommer"
                                        >
                                            <Edit size={14} />
                                        </button>
                                    {/if}

                                    <!-- Edit/Manage Actions (Hidden by default, shown on hover/group-hover) -->
                                    <div class="flex items-center gap-1 opacity-0 group-hover:opacity-100 transition-opacity">
                                        {#if editingFolderId === folder.id}
                                            <button onclick={() => handleRenameFolder(folder.id)} class="p-1 hover:text-green-500" title="Valider">
                                                <Check size={14} />
                                            </button>
                                            <button onclick={() => editingFolderId = null} class="p-1 hover:text-red-500" title="Annuler">
                                                <X size={14} />
                                            </button>
                                        {:else}
                                            <!-- Add page to folder button -->
                                            <button onclick={() => handleCreatePage(folder.id)} class="p-1 text-stone-400 hover:text-burnt-orange" title="Ajouter un article">
                                                <Plus size={14} />
                                            </button>
                                            <!-- Delete Folder -->
                                            <button onclick={() => handleDeleteFolder(folder.id)} class="p-1 text-stone-400 hover:text-red-500" title="Supprimer">
                                                <Trash2 size={14} />
                                            </button>
                                        {/if}
                                    </div>
                                </div>
                            </div>

                            <!-- Folder Children Pages -->
                            {#if !collapsedFolders[folder.id]}
                                {@const folderPages = pages.filter(p => p.folder_id === folder.id)}
                               <div class="pl-6 space-y-0.5 border-l border-stone-100 ml-4 mt-1">
                                    {#if folderPages.length === 0}
                                        <span class="text-xs text-stone-400 italic block py-1">Vide</span>
                                    {/if}
                                    {#each folderPages as page}
                                        <div
                                            draggable="true"
                                            ondragstart={(e) => handleDragStart(e, page)}
                                            onclick={() => selectPage(page)}
                                            onkeydown={(e) => e.key === "Enter" && selectPage(page)}
                                            role="button"
                                            tabindex="0"
                                            class="w-full flex items-center justify-between p-1.5 rounded-lg text-left text-xs transition-colors hover:bg-stone-50 cursor-grab active:cursor-grabbing {selectedPage?.id === page.id ? 'bg-stone-100 font-bold text-burnt-orange' : 'text-stone-600'}"
                                        >
                                            <div class="flex items-center gap-2 truncate">
                                                <FileText size={12} class="shrink-0" />
                                                <span class="truncate">{page.title}</span>
                                            </div>
                                            <!-- Page specific visibility -->
                                            <button
                                                onclick={(e) => {
                                                    e.stopPropagation();
                                                    handleTogglePageVisibility(page);
                                                }}
                                                class="flex items-center gap-1 shrink-0 hover:bg-stone-200/50 p-0.5 rounded transition-colors"
                                                title={page.is_visible ? (folder.is_visible ? "Visible pour les joueurs" : "Masqué par héritage du dossier") : "Masqué pour les joueurs"}
                                            >
                                                {#if page.is_visible}
                                                    <Eye size={12} class={folder.is_visible ? "text-green-500" : "text-green-300"} />
                                                {:else}
                                                    <EyeOff size={12} class="text-stone-400" />
                                                {/if}
                                            </button>
                                        </div>
                                    {/each}
                                </div>
                            {/if}
                        </div>
                    {/each}
                </div>

                <!-- ROOT PAGES -->
                {#if pages.filter(p => p.folder_id === null).length >= 0}
                    {@const rootPages = pages.filter(p => p.folder_id === null)}
                     <div 
                        ondragover={(e) => handleDragOver(e, null)}
                        ondragleave={handleDragLeave}
                        ondrop={(e) => handleDrop(e, null)}
                        class="space-y-1 pt-2 border border-transparent rounded-xl transition-all 
                        {activeDropTargetId === 'root' ? 'bg-burnt-orange/10 border-burnt-orange border-dashed p-2' : ''}"
                    >
                        <h4 class="text-xs font-bold uppercase tracking-wider text-stone-400 mb-2 px-2">Articles à la racine</h4>
                        {#if rootPages.length === 0}
                            <p class="text-xs text-stone-400 italic px-2">Aucun article à la racine.</p>
                        {/if}
                        {#each rootPages as page}
                            <div
                                draggable="true"
                                ondragstart={(e) => handleDragStart(e, page)}
                                onclick={() => selectPage(page)}
                                onkeydown={(e) => e.key === "Enter" && selectPage(page)}
                                role="button"
                                tabindex="0"
                                class="w-full flex items-center justify-between p-2 rounded-xl text-left text-sm transition-colors hover:bg-stone-50 group cursor-grab active:cursor-grabbing {selectedPage?.id === page.id ? 'bg-stone-50 border border-stone-200 font-bold text-burnt-orange' : 'text-stone-700'}"
                            >
                                <div class="flex items-center gap-2 truncate">
                                    <FileText size={14} class="text-stone-400" />
                                    <span class="truncate">{page.title}</span>
                                </div>
                                <div class="flex items-center gap-2">
                                    <!-- Visibility Toggle (Always visible) -->
                                    <button 
                                        onclick={(e) => {
                                            e.stopPropagation();
                                            handleTogglePageVisibility(page);
                                        }}
                                        class="p-0.5 text-stone-400 hover:text-dark-gray transition-colors"
                                        title={page.is_visible ? "Visible pour les joueurs" : "Masqué pour les joueurs"}
                                    >
                                        {#if page.is_visible}
                                            <Eye size={14} class="text-green-500" />
                                        {:else}
                                            <EyeOff size={14} />
                                        {/if}
                                    </button>

                                    <!-- Delete button (Hidden by default, shown on hover/group-hover) -->
                                    <button 
                                        onclick={(e) => {
                                            e.stopPropagation();
                                            handleDeletePage(page.id);
                                        }}
                                        class="p-0.5 text-stone-400 hover:text-red-500 opacity-0 group-hover:opacity-100 transition-opacity"
                                        title="Supprimer"
                                    >
                                        <Trash2 size={14} />
                                    </button>
                                </div>
                            </div>
                        {/each}
                    </div>
                {/if}

                <!-- MEDIATHEQUE (DRIVE) -->
                <div class="mt-6 pt-6 border-t border-stone-100 space-y-4">
                    <div class="flex items-center justify-between px-2">
                        <button
                            onclick={() => assetsCollapsed = !assetsCollapsed}
                            class="flex items-center gap-2 text-xs font-bold uppercase tracking-wider text-stone-400 hover:text-stone-600 transition-colors"
                        >
                            {#if assetsCollapsed}
                                <ChevronRight size={14} />
                            {:else}
                                <ChevronDown size={14} />
                            {/if}
                            <ImageIcon size={14} class="text-burnt-orange" />
                            <span>Médiathèque (Drive)</span>
                            <span class="text-[10px] font-normal text-stone-400 font-mono">({assets.length})</span>
                        </button>

                        <div class="flex items-center gap-1">
                            <!-- Import from media library button -->
                            <button
                                onclick={() => (isMediaModalOpen = true)}
                                class="text-stone-400 hover:text-burnt-orange p-1 rounded transition-colors"
                                title="Choisir depuis la médiathèque"
                            >
                                <FolderHeart size={16} />
                            </button>

                            <!-- Upload file button -->
                            <button
                                onclick={() => fileInputRef?.click()}
                                class="text-stone-400 hover:text-burnt-orange p-1 rounded transition-colors"
                                title="Importer une image depuis votre appareil"
                            >
                                <UploadCloud size={16} />
                            </button>
                        </div>
                        <input
                            type="file"
                            accept="image/*"
                            multiple
                            class="hidden"
                            bind:this={fileInputRef}
                            onchange={(e) => e.currentTarget.files && handleUploadFiles(e.currentTarget.files)}
                        />
                    </div>

                    {#if !assetsCollapsed}
                        {#if assetsLoading || isImportingFromMedia}
                            <div class="flex items-center justify-center py-4">
                                <Loader2 class="animate-spin text-stone-400" size={16} />
                            </div>
                        {:else if assets.length === 0}
                            <div class="p-4 border border-dashed border-stone-200 rounded-xl text-center space-y-2">
                                <FileImage size={24} class="text-stone-300 mx-auto" />
                                <p class="text-[10px] text-stone-400 leading-normal">
                                    Aucune image importée. Déposez des images dans l'éditeur ou importez-les ici.
                                </p>
                            </div>
                        {:else}
                            {#if assetsError}
                                <p class="text-[10px] text-red-500 px-2">{assetsError}</p>
                            {/if}
                            <div class="grid grid-cols-3 gap-2 px-1">
                                {#each assets as asset}
                                    <div class="group relative aspect-square rounded-lg border border-stone-150 overflow-hidden bg-stone-50 hover:border-burnt-orange transition-all">
                                        <!-- Thumbnail click inserts markdown image -->
                                        <button
                                            onclick={() => {
                                                const altText = asset.name.split('_').slice(1).join('_').split('.').shift() || "image";
                                                insertMarkdown(`![${altText}](${asset.url})`);
                                            }}
                                            class="w-full h-full p-0 border-0 bg-transparent block cursor-pointer"
                                            title="Insérer dans l'article"
                                        >
                                            <img
                                                src={asset.url}
                                                alt={asset.name}
                                                draggable="true"
                                                ondragstart={(e) => handleAssetDragStart(e, asset)}
                                                class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-300"
                                            />
                                        </button>

                                        <!-- Asset Actions -->
                                        <div class="absolute top-1 right-1 flex gap-1 opacity-0 group-hover:opacity-100 transition-opacity">
                                            <!-- Delete asset -->
                                            <button
                                                onclick={() => handleDeleteAsset(asset.name)}
                                                class="p-1 bg-white/90 hover:bg-red-50 text-stone-400 hover:text-red-500 rounded shadow-xs backdrop-blur-xs transition-colors cursor-pointer"
                                                title="Supprimer définitivement"
                                            >
                                                <Trash2 size={10} />
                                            </button>
                                        </div>
                                    </div>
                                {/each}
                            </div>
                        {/if}
                    {/if}
                </div>
            </div>
        {/if}
    </div>

    <!-- RIGHT EDITOR / PREVIEW -->
    <div class="flex-1 flex flex-col h-full bg-stone-50/50 rounded-2xl border border-stone-100 p-6 overflow-hidden">
        {#if selectedPage}
            <!-- 1. Header with Metadata controls -->
            <div class="flex flex-col gap-4 pb-4 border-b border-stone-100 mb-4">
                <div class="flex items-center justify-between gap-4">
                    <!-- Title Input -->
                    <input
                        type="text"
                        bind:value={selectedPage.title}
                        oninput={triggerAutoSave}
                        class="text-2xl font-display font-bold text-dark-gray bg-transparent border-b border-transparent hover:border-stone-200 focus:border-burnt-orange focus:outline-none w-full pb-1 transition-colors"
                        placeholder="Titre de l'article..."
                    />

                    <!-- Save Status Indicator -->
                    <div class="shrink-0 flex items-center">
                        {#if saveStatus === "saving"}
                            <span class="text-stone-400 flex items-center gap-1 text-xs">
                                <Loader2 size={14} class="animate-spin" /> Sauvegarde...
                            </span>
                        {:else if saveStatus === "saved"}
                            <span class="text-green-500 flex items-center gap-1 text-xs">
                                <CheckCircle size={14} /> Enregistré
                            </span>
                        {:else if saveStatus === "error"}
                            <span class="text-red-500 flex items-center gap-1 text-xs">
                                <AlertCircle size={14} /> Erreur de sauvegarde
                            </span>
                        {/if}
                    </div>
                </div>

                <!-- Parent Folder and Visibility Row -->
                <div class="flex flex-wrap items-center gap-4 text-xs text-stone-500">
                    <!-- Folder Selection -->
                    <div class="flex items-center gap-1.5">
                        <span>Dossier :</span>
                        <select
                            value={selectedPage.folder_id || ""}
                            onchange={(e) => handlePageFolderChange(e.currentTarget.value)}
                            class="px-2 py-1 bg-white border border-stone-200 rounded-lg focus:outline-none focus:border-burnt-orange"
                        >
                            <option value="">(Aucun dossier)</option>
                            {#each folders as folder}
                                <option value={folder.id}>{folder.name}</option>
                            {/each}
                        </select>
                    </div>

                    <!-- Visibility Toggle switch -->
                    <button
                        onclick={() => handleTogglePageVisibility(selectedPage!)}
                        class="flex items-center gap-1.5 px-3 py-1 bg-white hover:bg-stone-50 border border-stone-200 rounded-lg transition-all"
                    >
                        {#if selectedPage.is_visible}
                            <Eye size={14} class="text-green-500" />
                            <span class="text-green-600 font-semibold">Visible pour les joueurs</span>
                        {:else}
                            <EyeOff size={14} class="text-stone-400" />
                            <span class="text-stone-500">Masqué pour les joueurs</span>
                        {/if}
                    </button>
                </div>
            </div>

            <!-- 2. Editor Tabs (Edit / Preview) -->
            <div class="flex border-b border-stone-200 mb-4 bg-white rounded-lg p-1 self-start shadow-sm">
                <button
                    onclick={() => activeEditorTab = "edit"}
                    class="px-4 py-1.5 rounded-md text-xs font-bold transition-all {activeEditorTab === 'edit' ? 'bg-dark-gray text-white' : 'text-stone-500 hover:text-dark-gray'}"
                >
                    Éditer (Markdown)
                </button>
                <button
                    onclick={() => activeEditorTab = "preview"}
                    class="px-4 py-1.5 rounded-md text-xs font-bold transition-all {activeEditorTab === 'preview' ? 'bg-dark-gray text-white' : 'text-stone-500 hover:text-dark-gray'}"
                >
                    Prévisualiser
                </button>
            </div>

            <!-- 3. Tab Contents -->
            <div class="flex-1 overflow-y-auto bg-white rounded-xl border border-stone-100 p-4 relative flex flex-col">
                {#if activeEditorTab === "edit"}
                    <div class="relative flex-1 flex flex-col">
                        <textarea
                            bind:this={textareaElement}
                            bind:value={selectedPage.content}
                            oninput={triggerAutoSave}
                            ondrop={handleTextareaDrop}
                            onpaste={handleTextareaPaste}
                            class="w-full flex-1 resize-none focus:outline-none font-mono text-sm leading-relaxed text-stone-700 bg-white"
                            placeholder="# Titre du lore&#10;&#10;Écrivez votre lore en Markdown ici. Vous pouvez glisser-déposer ou coller une image directement ici pour l'importer dans l'article.&#10;&#10;Exemple de Markdown :&#10;- **gras**, *italique*&#10;- [liens](https://...)&#10;- listes et tableaux..."
                        ></textarea>
                        {#if isUploadingAsset}
                            <div class="absolute inset-0 bg-white/75 backdrop-blur-xs flex flex-col items-center justify-center gap-2 text-stone-600 z-10">
                                <Loader2 class="animate-spin text-burnt-orange" size={24} />
                                <span class="text-xs font-bold uppercase tracking-wider text-stone-500">Téléversement de l'image...</span>
                            </div>
                        {/if}
                    </div>
                {:else}
                    <div class="markdown-content flex-1 max-w-none">
                        {#if selectedPage.content}
                            {@html previewHtml}
                        {:else}
                            <p class="text-stone-400 italic text-sm">Le contenu est vide. Écrivez quelque chose dans l'onglet Éditer pour le voir s'afficher ici.</p>
                        {/if}
                    </div>
                {/if}
            </div>

            <!-- 4. Danger zone (footer) -->
            <div class="pt-4 mt-4 border-t border-stone-100 flex justify-end">
                <button
                    onclick={() => handleDeletePage(selectedPage!.id)}
                    class="flex items-center gap-1.5 px-3 py-2 text-xs font-bold text-red-500 hover:bg-red-50 rounded-xl transition-all"
                >
                    <Trash2 size={14} />
                    <span>Supprimer cet article</span>
                </button>
            </div>
        {:else}
            <!-- Blank State -->
            <div class="flex-1 flex flex-col items-center justify-center text-center p-6 space-y-4">
                <div class="w-16 h-16 bg-stone-100 rounded-full flex items-center justify-center text-stone-400">
                    <BookOpen size={32} />
                </div>
                <div>
                    <h3 class="font-bold text-dark-gray text-lg mb-1">Aucun article sélectionné</h3>
                    <p class="text-stone-400 text-sm max-w-sm">
                        Sélectionnez un article existant dans le panneau de gauche ou créez-en un nouveau pour commencer à rédiger l'histoire de votre partie.
                    </p>
                </div>
            </div>
        {/if}
    </div>
</div>

<MediaSelectorModal
    isOpen={isMediaModalOpen}
    onSelect={handleSelectFromMediaLibrary}
    onClose={() => {
        isMediaModalOpen = false;
    }}
/>

<style>
    /* Premium Markdown styling for Story preview and presentation */
    .markdown-content :global(h1) {
        font-size: 1.75rem;
        font-weight: 800;
        color: #2D2A26;
        margin-top: 1.5rem;
        margin-bottom: 0.75rem;
        border-bottom: 2px solid #F3F4F6;
        padding-bottom: 0.25rem;
        font-family: inherit;
    }
    .markdown-content :global(h2) {
        font-size: 1.35rem;
        font-weight: 700;
        color: #2D2A26;
        margin-top: 1.25rem;
        margin-bottom: 0.5rem;
    }
    .markdown-content :global(h3) {
        font-size: 1.15rem;
        font-weight: 600;
        color: #2D2A26;
        margin-top: 1rem;
        margin-bottom: 0.5rem;
    }
    .markdown-content :global(p) {
        margin-bottom: 1rem;
        line-height: 1.625;
        color: #4A4A4A;
    }
    .markdown-content :global(ul) {
        list-style-type: disc;
        padding-left: 1.5rem;
        margin-bottom: 1rem;
    }
    .markdown-content :global(ol) {
        list-style-type: decimal;
        padding-left: 1.5rem;
        margin-bottom: 1rem;
    }
    .markdown-content :global(li) {
        margin-bottom: 0.25rem;
    }
    .markdown-content :global(strong) {
        font-weight: 700;
        color: #1A1A1A;
    }
    .markdown-content :global(em) {
        font-style: italic;
    }
    .markdown-content :global(code) {
        font-family: monospace;
        background-color: #F3F4F6;
        padding: 0.125rem 0.25rem;
        border-radius: 0.25rem;
        font-size: 0.85em;
        color: #E07A5F;
    }
    .markdown-content :global(pre) {
        background-color: #1F2937;
        color: #F9FAFB;
        padding: 1rem;
        border-radius: 0.75rem;
        overflow-x: auto;
        margin-bottom: 1rem;
        font-size: 0.875rem;
    }
    .markdown-content :global(pre code) {
        background-color: transparent;
        padding: 0;
        color: inherit;
        font-size: inherit;
    }
    .markdown-content :global(blockquote) {
        border-left: 4px solid #E07A5F;
        padding-left: 1rem;
        font-style: italic;
        color: #6B7280;
        margin-bottom: 1rem;
        background-color: #FAFAF9;
        padding-top: 0.5rem;
        padding-bottom: 0.5rem;
        border-radius: 0 0.5rem 0.5rem 0;
    }
    .markdown-content :global(a) {
        color: #E07A5F;
        text-decoration: underline;
        font-weight: 500;
    }
    .markdown-content :global(table) {
        width: 100%;
        border-collapse: collapse;
        margin-bottom: 1rem;
        font-size: 0.875rem;
    }
    .markdown-content :global(th), .markdown-content :global(td) {
        border: 1px solid #E5E7EB;
        padding: 0.5rem 0.75rem;
    }
    .markdown-content :global(th) {
        background-color: #F9FAFB;
        font-weight: 600;
        text-align: left;
    }
</style>
