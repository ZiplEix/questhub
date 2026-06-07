<script lang="ts">
    import { onMount } from "svelte";
    import { 
        Folder, 
        FileText, 
        ChevronRight, 
        ChevronDown, 
        ArrowLeft, 
        BookOpen, 
        EyeOff, 
        Loader2,
        ShieldAlert
    } from "lucide-svelte";
    import { marked } from "marked";
    import { supabase } from "$lib/supabaseClient";
    import { 
        fetchStoryFolders, 
        fetchStoryPages,
        type StoryFolder,
        type StoryPage
    } from "$lib/api/story";

    let { gameId, isGM = false } = $props<{ gameId: string; isGM?: boolean }>();

    let folders = $state<StoryFolder[]>([]);
    let pages = $state<StoryPage[]>([]);
    let loading = $state(true);
    let selectedPage = $state<StoryPage | null>(null);

    // Collapsed state for folders
    let collapsedFolders = $state<Record<string, boolean>>({});

    onMount(async () => {
        try {
            await refreshData();
        } catch (error) {
            console.error("Failed to load story sidebar:", error);
        } finally {
            loading = false;
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

    // Realtime Postgres sync
    $effect(() => {
        if (!gameId) return;

        const channel = supabase.channel(`story_sidebar_${gameId}`)
            .on(
                'postgres_changes',
                { event: '*', schema: 'public', table: 'story_folders', filter: `game_id=eq.${gameId}` },
                (payload) => {
                    if (payload.eventType === 'INSERT') {
                        folders = [...folders, payload.new as StoryFolder].sort((a, b) => a.name.localeCompare(b.name));
                    } else if (payload.eventType === 'UPDATE') {
                        folders = folders.map(f => f.id === payload.new.id ? payload.new as StoryFolder : f).sort((a, b) => a.name.localeCompare(b.name));
                    } else if (payload.eventType === 'DELETE') {
                        folders = folders.filter(f => f.id !== payload.old.id);
                    }
                }
            )
            .on(
                'postgres_changes',
                { event: '*', schema: 'public', table: 'story_pages', filter: `game_id=eq.${gameId}` },
                (payload) => {
                    if (payload.eventType === 'INSERT') {
                        pages = [...pages, payload.new as StoryPage].sort((a, b) => a.title.localeCompare(b.title));
                    } else if (payload.eventType === 'UPDATE') {
                        pages = pages.map(p => p.id === payload.new.id ? payload.new as StoryPage : p).sort((a, b) => a.title.localeCompare(b.title));
                        if (selectedPage && selectedPage.id === payload.new.id) {
                            selectedPage = payload.new as StoryPage;
                        }
                    } else if (payload.eventType === 'DELETE') {
                        pages = pages.filter(p => p.id !== payload.old.id);
                        if (selectedPage && selectedPage.id === payload.old.id) {
                            selectedPage = null;
                        }
                    }
                }
            )
            .subscribe();

        return () => {
            supabase.removeChannel(channel);
        };
    });

    // Reactive visibility calculations
    let visibleFolders = $derived.by(() => {
        if (isGM) return folders;
        return folders.filter(f => f.is_visible);
    });

    let visiblePages = $derived.by(() => {
        if (isGM) return pages;
        return pages.filter(p => {
            if (!p.is_visible) return false;
            if (p.folder_id === null) return true;
            const parentFolder = folders.find(f => f.id === p.folder_id);
            return parentFolder ? parentFolder.is_visible : false;
        });
    });

    // Derived markdown compiled HTML
    let selectedPageHtml = $derived.by(() => {
        if (!selectedPage) return "";
        try {
            return marked.parse(selectedPage.content || "") as string;
        } catch (e) {
            console.error(e);
            return "<p class='text-red-500'>Erreur de rendu.</p>";
        }
    });

    // Format Date helper
    function formatDate(dateStr: string) {
        if (!dateStr) return "";
        const date = new Date(dateStr);
        return date.toLocaleDateString("fr-FR", {
            day: "numeric",
            month: "long",
            year: "numeric"
        });
    }

    // Helper to check if a page is actually hidden from players
    function isPageHiddenForPlayers(page: StoryPage) {
        if (!page.is_visible) return true;
        if (page.folder_id) {
            const folder = folders.find(f => f.id === page.folder_id);
            return folder ? !folder.is_visible : true;
        }
        return false;
    }
</script>

<div class="h-full flex flex-col bg-stone-50 overflow-hidden text-stone-700">
    {#if selectedPage}
        <!-- PAGE DETAIL VIEW -->
        <div class="flex-1 flex flex-col overflow-hidden bg-white">
            <!-- Header and Back button -->
            <div class="p-4 border-b border-stone-100 flex items-center gap-3">
                <button
                    onclick={() => selectedPage = null}
                    class="p-1.5 hover:bg-stone-50 text-stone-500 rounded-lg transition-colors"
                    aria-label="Retour à la liste"
                >
                    <ArrowLeft size={16} />
                </button>
                <div class="truncate">
                    <h3 class="font-bold text-dark-gray truncate text-sm">{selectedPage.title}</h3>
                    <p class="text-[10px] text-stone-400">Mis à jour le {formatDate(selectedPage.updated_at)}</p>
                </div>
            </div>

            <!-- Scrollable Content -->
            <div class="flex-1 overflow-y-auto p-4 space-y-4">
                <!-- GM Warning Banner -->
                {#if isGM && isPageHiddenForPlayers(selectedPage)}
                    <div class="p-3 bg-amber-50 border border-amber-200 text-amber-800 rounded-xl text-xs flex items-start gap-2 shadow-sm">
                        <EyeOff size={16} class="shrink-0 text-amber-600 mt-0.5" />
                        <div>
                            <span class="font-bold">Masqué pour les joueurs</span>
                            <p class="text-[11px] text-amber-700/80 mt-0.5">
                                {#if !selectedPage.is_visible}
                                    Cet article est marqué comme masqué.
                                {:else}
                                    Le dossier contenant cet article est masqué.
                                {/if}
                            </p>
                        </div>
                    </div>
                {/if}

                <!-- HTML Rendered Article -->
                <div class="markdown-content max-w-none">
                    {#if selectedPage.content}
                        {@html selectedPageHtml}
                    {:else}
                        <p class="text-stone-400 italic text-sm text-center py-6">Cet article est vide.</p>
                    {/if}
                </div>
            </div>
        </div>
    {:else}
        <!-- TREE LIST VIEW -->
        <div class="p-4 border-b border-stone-100 bg-white">
            <h3 class="font-bold text-dark-gray flex items-center gap-2 text-sm">
                <BookOpen size={16} class="text-burnt-orange" />
                <span>Histoire & Lore</span>
            </h3>
        </div>

        {#if loading}
            <div class="flex-1 flex items-center justify-center">
                <Loader2 class="animate-spin text-stone-400" />
            </div>
        {:else if visibleFolders.length === 0 && visiblePages.filter(p => p.folder_id === null).length === 0}
            <div class="flex-1 flex flex-col items-center justify-center text-center p-6 space-y-4">
                <div class="w-12 h-12 bg-stone-100 rounded-full flex items-center justify-center text-stone-400">
                    <BookOpen size={24} />
                </div>
                <div>
                    <h4 class="font-bold text-stone-600 text-sm">Rien à afficher</h4>
                    <p class="text-stone-400 text-xs max-w-xs mt-1">
                        {#if isGM}
                            Aucun élément d'histoire créé. Rendez-vous dans les paramètres de la table pour écrire du lore.
                        {:else}
                            Aucune histoire ou élément de lore n'a encore été révélé par le Maître du Jeu.
                        {/if}
                    </p>
                </div>
            </div>
        {:else}
            <div class="flex-1 overflow-y-auto p-4 space-y-4">
                <!-- FOLDERS -->
                {#if visibleFolders.length > 0}
                    <div class="space-y-1">
                        <span class="text-[10px] font-bold text-stone-400 uppercase tracking-wider block mb-1">Dossiers</span>
                        {#each visibleFolders as folder}
                            <div class="space-y-0.5">
                                <!-- Folder header -->
                                <button
                                    onclick={() => collapsedFolders[folder.id] = !collapsedFolders[folder.id]}
                                    class="w-full flex items-center justify-between p-2 rounded-xl text-left hover:bg-stone-100 transition-colors"
                                >
                                    <div class="flex items-center gap-2 text-xs font-semibold text-stone-700 truncate">
                                        {#if collapsedFolders[folder.id]}
                                            <ChevronRight size={14} class="text-stone-400" />
                                        {:else}
                                            <ChevronDown size={14} class="text-stone-400" />
                                        {/if}
                                        <Folder size={14} class="text-burnt-orange fill-burnt-orange/10" />
                                        <span class="truncate">{folder.name}</span>
                                    </div>
                                    {#if isGM && !folder.is_visible}
                                        <span title="Masqué pour les joueurs">
                                            <EyeOff size={12} class="text-stone-400" />
                                        </span>
                                    {/if}
                                </button>

                                <!-- Folder Pages -->
                                {#if !collapsedFolders[folder.id]}
                                    {@const folderPages = visiblePages.filter(p => p.folder_id === folder.id)}
                                    <div class="pl-4 ml-3 border-l border-stone-200/60 space-y-0.5">
                                        {#if folderPages.length === 0}
                                            <span class="text-[10px] text-stone-400 italic block p-2">Aucun article</span>
                                        {:else}
                                            {#each folderPages as page}
                                                <button
                                                    onclick={() => selectedPage = page}
                                                    class="w-full flex items-center justify-between p-1.5 rounded-lg text-left text-xs text-stone-600 hover:bg-stone-100 transition-colors"
                                                >
                                                    <div class="flex items-center gap-1.5 truncate">
                                                        <FileText size={12} class="text-stone-400" />
                                                        <span class="truncate">{page.title}</span>
                                                    </div>
                                                    {#if isGM && !page.is_visible}
                                                        <span title="Masqué">
                                                            <EyeOff size={10} class="text-stone-400" />
                                                        </span>
                                                    {/if}
                                                </button>
                                            {/each}
                                        {/if}
                                    </div>
                                {/if}
                            </div>
                        {/each}
                    </div>
                {/if}

                <!-- ROOT PAGES -->
                {#if visiblePages.filter(p => p.folder_id === null).length > 0}
                    {@const rootPages = visiblePages.filter(p => p.folder_id === null)}
                    <div class="space-y-1">
                        <span class="text-[10px] font-bold text-stone-400 uppercase tracking-wider block mb-1">Articles</span>
                        {#each rootPages as page}
                            <button
                                onclick={() => selectedPage = page}
                                class="w-full flex items-center justify-between p-2 rounded-xl text-left hover:bg-stone-100 transition-colors text-xs text-stone-700"
                            >
                                <div class="flex items-center gap-2 truncate">
                                    <FileText size={14} class="text-stone-400" />
                                    <span class="truncate font-medium">{page.title}</span>
                                </div>
                                {#if isGM && !page.is_visible}
                                    <span title="Masqué pour les joueurs">
                                        <EyeOff size={12} class="text-stone-400" />
                                    </span>
                                {/if}
                            </button>
                        {/each}
                    </div>
                {/if}
            </div>
        {/if}
    {/if}
</div>

<style>
    /* Premium Markdown styling for presentation in sidebar */
    .markdown-content :global(h1) {
        font-size: 1.35rem;
        font-weight: 800;
        color: #2D2A26;
        margin-top: 1.25rem;
        margin-bottom: 0.5rem;
        border-bottom: 1px solid #F3F4F6;
        padding-bottom: 0.125rem;
    }
    .markdown-content :global(h2) {
        font-size: 1.15rem;
        font-weight: 700;
        color: #2D2A26;
        margin-top: 1rem;
        margin-bottom: 0.35rem;
    }
    .markdown-content :global(h3) {
        font-size: 1rem;
        font-weight: 600;
        color: #2D2A26;
        margin-top: 0.75rem;
        margin-bottom: 0.25rem;
    }
    .markdown-content :global(p) {
        margin-bottom: 0.75rem;
        line-height: 1.5;
        font-size: 0.825rem;
        color: #4A4A4A;
    }
    .markdown-content :global(ul) {
        list-style-type: disc;
        padding-left: 1.25rem;
        margin-bottom: 0.75rem;
        font-size: 0.825rem;
    }
    .markdown-content :global(ol) {
        list-style-type: decimal;
        padding-left: 1.25rem;
        margin-bottom: 0.75rem;
        font-size: 0.825rem;
    }
    .markdown-content :global(li) {
        margin-bottom: 0.2rem;
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
        padding: 0.1rem 0.2rem;
        border-radius: 0.25rem;
        font-size: 0.85em;
        color: #E07A5F;
    }
    .markdown-content :global(pre) {
        background-color: #1F2937;
        color: #F9FAFB;
        padding: 0.75rem;
        border-radius: 0.5rem;
        overflow-x: auto;
        margin-bottom: 0.75rem;
        font-size: 0.775rem;
    }
    .markdown-content :global(pre code) {
        background-color: transparent;
        padding: 0;
        color: inherit;
        font-size: inherit;
    }
    .markdown-content :global(blockquote) {
        border-left: 3px solid #E07A5F;
        padding-left: 0.75rem;
        font-style: italic;
        color: #6B7280;
        margin-bottom: 0.75rem;
        background-color: #FAFAF9;
        padding-top: 0.35rem;
        padding-bottom: 0.35rem;
        border-radius: 0 0.35rem 0.35rem 0;
        font-size: 0.8rem;
    }
    .markdown-content :global(a) {
        color: #E07A5F;
        text-decoration: underline;
    }
    .markdown-content :global(table) {
        width: 100%;
        border-collapse: collapse;
        margin-bottom: 0.75rem;
        font-size: 0.775rem;
    }
    .markdown-content :global(th), .markdown-content :global(td) {
        border: 1px solid #E5E7EB;
        padding: 0.4rem 0.6rem;
    }
    .markdown-content :global(th) {
        background-color: #F9FAFB;
        font-weight: 600;
        text-align: left;
    }
</style>
