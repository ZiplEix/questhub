<script lang="ts">
    import { onMount } from "svelte";
    import Header from "$lib/components/Header.svelte";
    import { 
        Search, 
        PackageOpen, 
        Plus, 
        Globe, 
        FolderPlus, 
        Trash2, 
        Edit, 
        Check, 
        Package, 
        ExternalLink, 
        Lock, 
        Users, 
        Skull, 
        BookOpen,
        X,
        Sparkles
    } from "lucide-svelte";
    import { 
        fetchMarketplaceTemplates, 
        createTemplate, 
        updateTemplate, 
        deleteTemplate, 
        importTemplateToGame, 
        fetchGames,
        fetchUserCharacters,
        fetchUserMonsters,
        type MarketplaceTemplate
    } from "$lib/api";
    import { authClient } from "$lib/auth-client";
    import { goto } from "$app/navigation";

    let user = $state<any>(null);
    let rawTemplates = $state<MarketplaceTemplate[]>([]);
    let myGames = $state<any[]>([]);
    
    // User character & monster lists for bundle creation
    let userCharacters = $state<any[]>([]);
    let userMonsters = $state<any[]>([]);

    let loading = $state(true);
    let searchQuery = $state("");
    let selectedType = $state("TOUT");
    let showOnlyMyPublications = $state(false);

    // Derived templates representing both parents and virtual children
    let processedTemplates = $derived.by(() => {
        let list: any[] = [];
        for (const t of rawTemplates) {
            // Keep the original template (monsters, characters, PNJ, or bundle itself)
            list.push(t);
            // If it's a bundle, also extract all its items as individual virtual templates
            if (t.type === 'BUNDLE' && t.data?.items) {
                t.data.items.forEach((item: any, idx: number) => {
                    list.push({
                        id: `${t.id}-item-${idx}`,
                        parent_bundle_id: t.id,
                        parent_bundle_name: t.name,
                        created_by: t.created_by,
                        author_name: t.author_name,
                        name: item.name,
                        description: item.data?.description || `Fait partie du pack "${t.name}"`,
                        type: item.type, // 'PERSONNAGE' | 'PNJ' | 'MONSTRE'
                        data: item.data, // stats
                        is_public: t.is_public,
                        uses: t.uses,
                        created_at: t.created_at,
                        is_virtual: true
                    });
                });
            }
        }

        // Apply type filters on the client
        if (selectedType === "PACK") {
            return list.filter(t => t.type === 'BUNDLE');
        } else if (selectedType !== "TOUT") {
            return list.filter(t => t.type === selectedType);
        }
        
        return list;
    });

    // Modal States
    let isImportModalOpen = $state(false);
    let selectedTemplateForImport = $state<MarketplaceTemplate | null>(null);
    let targetGameId = $state("");
    let isImporting = $state(false);

    let isBundleDetailsModalOpen = $state(false);
    let selectedBundleForDetails = $state<MarketplaceTemplate | null>(null);

    let isCreateBundleModalOpen = $state(false);
    let bundleName = $state("");
    let bundleDescription = $state("");
    let bundleIsPublic = $state(true);
    let bundleCoverUrl = $state("");
    let selectedItemIds = $state<string[]>([]);
    let isCreatingBundle = $state(false);

    let isEditModalOpen = $state(false);
    let selectedTemplateForEdit = $state<MarketplaceTemplate | null>(null);
    let editName = $state("");
    let editDescription = $state("");
    let editIsPublic = $state(true);
    let isSavingEdit = $state(false);

    // Fetch initial data
    async function loadData() {
        try {
            loading = true;
            
            // Check auth session
            const { data } = await authClient.getSession();
            if (data?.user) {
                user = data.user;
                
                // Fetch games where user is GM
                const gamesList = await fetchGames();
                myGames = gamesList.filter(g => g.gm_id === user.id);

                // Fetch characters and monsters owned by user for bundle creation
                const [chars, monsts] = await Promise.all([
                    fetchUserCharacters(),
                    fetchUserMonsters()
                ]);
                userCharacters = chars;
                userMonsters = monsts;
            }

            await loadTemplates();
        } catch (error) {
            console.error("Failed to load marketplace data:", error);
        } finally {
            loading = false;
        }
    }

    async function loadTemplates() {
        try {
            // Fetch all templates matching search/mine filters, and we'll filter types client-side
            const list = await fetchMarketplaceTemplates({
                search: searchQuery,
                onlyUser: showOnlyMyPublications
            });
            rawTemplates = list;
        } catch (error) {
            console.error("Failed to load templates:", error);
        }
    }

    // React to filter changes
    $effect(() => {
        // Dependencies to trigger reload on changes
        const _mine = showOnlyMyPublications;
        loadTemplates();
    });

    onMount(() => {
        loadData();
    });

    // Handle single template/bundle delete
    async function handleDelete(templateId: string, event: Event) {
        event.stopPropagation();
        if (!confirm("Êtes-vous sûr de vouloir supprimer cette publication du marché ?")) return;

        try {
            await deleteTemplate(templateId);
            rawTemplates = rawTemplates.filter(t => t.id !== templateId);
            alert("Publication supprimée avec succès.");
        } catch (error) {
            console.error("Failed to delete template:", error);
            alert("Erreur lors de la suppression.");
        }
    }

    // Open Edit Modal
    function openEditModal(template: MarketplaceTemplate, event: Event) {
        event.stopPropagation();
        selectedTemplateForEdit = template;
        editName = template.name;
        editDescription = template.description || "";
        editIsPublic = template.is_public;
        isEditModalOpen = true;
    }

    // Save edited template metadata
    async function handleSaveEdit() {
        if (!selectedTemplateForEdit) return;
        try {
            isSavingEdit = true;
            await updateTemplate(selectedTemplateForEdit.id, {
                name: editName,
                description: editDescription || null,
                is_public: editIsPublic
            });
            
            // Update local state
            rawTemplates = rawTemplates.map(t => 
                t.id === selectedTemplateForEdit!.id 
                    ? { ...t, name: editName, description: editDescription, is_public: editIsPublic }
                    : t
            );
            
            isEditModalOpen = false;
            selectedTemplateForEdit = null;
            alert("Publication mise à jour avec succès.");
        } catch (error) {
            console.error("Failed to update template:", error);
            alert("Erreur lors de la mise à jour.");
        } finally {
            isSavingEdit = false;
        }
    }

    // Open Import Dialog
    function openImportModal(template: MarketplaceTemplate) {
        if (!user) {
            goto("/login");
            return;
        }
        selectedTemplateForImport = template;
        targetGameId = myGames.length > 0 ? myGames[0].id : "";
        isImportModalOpen = true;
    }

    // Open Bundle Details Modal
    function openBundleDetails(bundle: MarketplaceTemplate) {
        selectedBundleForDetails = bundle;
        isBundleDetailsModalOpen = true;
    }

    // Import template or bundle into game
    async function handleImport() {
        if (!selectedTemplateForImport || !targetGameId) return;
        try {
            isImporting = true;
            
            await importTemplateToGame(selectedTemplateForImport, targetGameId);
            
            // Increment uses count locally
            const targetId = (selectedTemplateForImport as any).is_virtual
                ? (selectedTemplateForImport as any).parent_bundle_id
                : selectedTemplateForImport.id;

            rawTemplates = rawTemplates.map(t => 
                t.id === targetId 
                    ? { ...t, uses: (t.uses || 0) + 1 } 
                    : t
            );

            alert("Importation réussie dans votre partie ! Retrouvez la création dans la liste de configuration correspondante.");
            isImportModalOpen = false;
            selectedTemplateForImport = null;
        } catch (error) {
            console.error("Failed to import template:", error);
            alert("Erreur lors de l'importation.");
        } finally {
            isImporting = false;
        }
    }

    // Create Bundle Form Submit
    async function handleCreateBundle() {
        if (!bundleName) {
            alert("Veuillez donner un nom à votre lot/bundle.");
            return;
        }
        if (selectedItemIds.length === 0) {
            alert("Veuillez sélectionner au moins un personnage ou un monstre à inclure dans le lot.");
            return;
        }

        try {
            isCreatingBundle = true;

            // Gather full character data payloads
            const items = selectedItemIds.map(id => {
                const char = userCharacters.find(c => c.id === id) || userMonsters.find(m => m.id === id);
                if (!char) return null;
                
                // Strip database IDs
                const { id: _, game_id, user_id, created_at, ...cleanData } = char;
                
                let itemType: 'PERSONNAGE' | 'PNJ' | 'MONSTRE' = 'PERSONNAGE';
                if (char.type === 'MONSTER') {
                    itemType = 'MONSTRE';
                } else if (char.is_npc) {
                    itemType = 'PNJ';
                }

                return {
                    name: char.name,
                    type: itemType,
                    data: cleanData
                };
            }).filter(Boolean);

            // Use first character avatar as default cover url if empty
            let coverUrl = bundleCoverUrl;
            if (!coverUrl && items.length > 0) {
                coverUrl = items[0]?.data?.avatar_url || "";
            }

            await createTemplate({
                name: bundleName,
                description: bundleDescription || null,
                type: 'BUNDLE',
                data: {
                    cover_url: coverUrl,
                    items: items
                },
                is_public: bundleIsPublic
            });

            alert("Lot/Bundle créé et publié avec succès !");
            isCreateBundleModalOpen = false;
            
            // Clear inputs
            bundleName = "";
            bundleDescription = "";
            bundleCoverUrl = "";
            selectedItemIds = [];
            
            // Refresh list
            loadTemplates();
        } catch (error) {
            console.error("Failed to create bundle:", error);
            alert("Erreur lors de la création du bundle.");
        } finally {
            isCreatingBundle = false;
        }
    }

    // Toggle character selection in bundle form
    function toggleItemSelection(id: string) {
        if (selectedItemIds.includes(id)) {
            selectedItemIds = selectedItemIds.filter(itemId => itemId !== id);
        } else {
            selectedItemIds = [...selectedItemIds, id];
        }
    }

    // Search input handler
    let searchTimeout: any;
    function handleSearchInput() {
        clearTimeout(searchTimeout);
        searchTimeout = setTimeout(() => {
            loadTemplates();
        }, 300);
    }

    const types = [
        { id: "TOUT", label: "Tout" },
        { id: "PERSONNAGE", label: "Personnages" },
        { id: "PNJ", label: "PNJs" },
        { id: "MONSTRE", label: "Monstres" },
        { id: "PACK", label: "Packs/Bundles" }
    ];
</script>

<div class="min-h-screen bg-cream font-sans text-dark-gray">
    <Header />

    <main class="max-w-7xl mx-auto px-4 md:px-8 py-16">
        <div class="mb-12 text-center">
            <h1 class="text-4xl font-display font-black mb-4 text-burnt-orange tracking-tight flex items-center justify-center gap-3">
                <Sparkles class="text-burnt-orange fill-burnt-orange/20 animate-pulse" size={36} />
                Marché Communautaire
            </h1>
            <p class="text-lg text-stone-500 max-w-2xl mx-auto font-medium">
                Partagez vos monstres, vos personnages pré-tirés, ou importez des bundles complets créés par les rôlistes de la communauté.
            </p>
        </div>

        <!-- Search & Actions Bar -->
        <div class="mb-12 flex flex-col md:flex-row items-center justify-between gap-6 bg-white p-6 rounded-2xl border border-stone-200/60 shadow-xs">
            <div class="relative w-full max-w-md">
                <Search class="absolute left-3.5 top-1/2 -translate-y-1/2 text-stone-400" size={18} />
                <input
                    type="text"
                    placeholder="Rechercher par nom, auteur, description..."
                    bind:value={searchQuery}
                    oninput={handleSearchInput}
                    class="w-full pl-10 pr-4 py-2.5 bg-stone-50 border border-stone-200 rounded-xl text-sm focus:outline-none focus:border-burnt-orange focus:ring-2 focus:ring-burnt-orange/20 transition-all font-medium text-dark-gray"
                />
            </div>

            <!-- Tabs Filters -->
            <div class="flex flex-wrap gap-2">
                {#each types as type}
                    <button
                        onclick={() => (selectedType = type.id)}
                        class="px-4 py-2 rounded-xl text-xs font-bold uppercase tracking-wider border transition-all cursor-pointer
                        {selectedType === type.id
                            ? 'bg-dark-gray text-white border-dark-gray shadow-md'
                            : 'bg-stone-50 text-stone-500 border-stone-200 hover:border-burnt-orange hover:text-burnt-orange'}"
                    >
                        {type.label}
                    </button>
                {/each}
            </div>

            <!-- Action buttons -->
            {#if user}
                <div class="flex items-center gap-3 shrink-0">
                    <button
                        onclick={() => (showOnlyMyPublications = !showOnlyMyPublications)}
                        class="px-4 py-2.5 rounded-xl text-xs font-bold uppercase tracking-wider border transition-all cursor-pointer flex items-center gap-2
                        {showOnlyMyPublications
                            ? 'bg-burnt-orange/10 text-burnt-orange border-burnt-orange/30'
                            : 'bg-stone-50 text-stone-500 border-stone-200 hover:border-stone-400'}"
                    >
                        {showOnlyMyPublications ? "Toutes les publications" : "Mes publications"}
                    </button>
                    
                    <button
                        onclick={() => (isCreateBundleModalOpen = true)}
                        class="px-4 py-2.5 bg-burnt-orange text-white rounded-xl text-xs font-bold uppercase tracking-wider hover:bg-opacity-95 shadow-sm hover:shadow-md transition-all flex items-center gap-2 cursor-pointer"
                    >
                        <FolderPlus size={16} />
                        Créer un pack
                    </button>
                </div>
            {/if}
        </div>

        <!-- Grid of Templates -->
        {#if loading}
            <div class="flex flex-col justify-center items-center py-24 gap-4">
                <div class="animate-spin rounded-full h-10 w-10 border-b-2 border-burnt-orange"></div>
                <p class="text-sm text-stone-500 font-bold font-mono">Chargement du marché...</p>
            </div>
        {:else if processedTemplates.length === 0}
            <div class="text-center py-24 bg-white rounded-2xl border border-stone-100 shadow-xs flex flex-col items-center gap-4">
                <PackageOpen size={48} class="text-stone-300" strokeWidth={1.5} />
                <p class="text-xl font-bold text-dark-gray">Aucune création trouvée</p>
                <p class="text-stone-400 text-sm max-w-sm">
                    Revenez plus tard ou publiez vos propres monstres et personnages depuis les paramètres de vos campagnes.
                </p>
            </div>
        {:else}
            <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
                {#each processedTemplates as template}
                    <div 
                        onclick={() => template.type === 'BUNDLE' && !template.is_virtual && openBundleDetails(template)}
                        class="group bg-white rounded-2xl border border-stone-200/80 shadow-xs hover:shadow-md hover:border-stone-300 transition-all duration-300 flex flex-col h-full overflow-hidden relative {template.type === 'BUNDLE' && !template.is_virtual ? 'cursor-pointer' : ''}"
                    >
                        <!-- Parent bundle indicator for virtual templates -->
                        {#if template.is_virtual}
                            <div class="absolute top-0 right-0 left-0 bg-stone-900/90 text-white text-[9px] font-bold py-1.5 px-3 flex items-center justify-between shadow-xs backdrop-blur-md z-10">
                                <span class="truncate">📦 Pack : {template.parent_bundle_name}</span>
                                <button 
                                    onclick={(e) => {
                                        e.stopPropagation();
                                        const parent = rawTemplates.find(p => p.id === template.parent_bundle_id);
                                        if (parent) openBundleDetails(parent);
                                    }}
                                    class="hover:text-burnt-orange underline cursor-pointer text-[8px] tracking-wide shrink-0"
                                >
                                    Voir le pack
                                </button>
                            </div>
                        {/if}

                        <!-- Cover Image / Icon -->
                        <div class="relative h-44 overflow-hidden bg-stone-50 border-b border-stone-100 flex items-center justify-center shrink-0">
                            {#if template.type === 'BUNDLE'}
                                {#if template.data?.cover_url}
                                    <img
                                        src={template.data.cover_url}
                                        alt={template.name}
                                        class="w-full h-full object-cover group-hover:scale-103 transition-transform duration-500"
                                    />
                                {:else}
                                    <div class="w-16 h-16 bg-burnt-orange/15 rounded-2xl flex items-center justify-center text-burnt-orange border border-burnt-orange/20 shadow-xs">
                                        <Package size={36} />
                                    </div>
                                {/if}
                            {:else}
                                <img
                                    src={template.data?.avatar_url || `https://api.dicebear.com/9.x/avataaars/svg?seed=${template.name}`}
                                    alt={template.name}
                                    class="w-full h-full object-cover group-hover:scale-103 transition-transform duration-500"
                                />
                            {/if}

                            <!-- Type badge -->
                            <div class="absolute top-3 left-3 flex gap-1.5 {template.is_virtual ? 'mt-6' : ''}">
                                <span class="bg-white/90 text-dark-gray px-2.5 py-1 rounded-lg text-[10px] font-extrabold shadow-xs backdrop-blur-md border border-stone-100 uppercase tracking-wider">
                                    {template.type === 'BUNDLE' ? 'PACK / BUNDLE' : template.type}
                                </span>
                                {#if !template.is_public}
                                    <span class="bg-stone-900/90 text-white p-1 rounded-lg shadow-xs backdrop-blur-md" title="Bibliothèque privée">
                                        <Lock size={12} />
                                    </span>
                                {/if}
                            </div>

                            <!-- Uses badge -->
                            <div class="absolute bottom-3 right-3">
                                <span class="bg-black/60 text-white px-2 py-0.5 rounded-md text-[9px] font-bold tracking-wider uppercase">
                                    {template.uses || 0} utilisations
                                </span>
                            </div>
                        </div>

                        <!-- Card Body -->
                        <div class="p-5 flex-1 flex flex-col justify-between">
                            <div class="space-y-2">
                                <div>
                                    <h3 class="font-display font-bold text-lg text-dark-gray leading-tight group-hover:text-burnt-orange transition-colors truncate">
                                        {template.name}
                                    </h3>
                                    <p class="text-xs text-stone-400 font-semibold">
                                        par @{template.author_name}
                                    </p>
                                </div>
                                
                                {#if template.description}
                                    <p class="text-xs text-stone-500 line-clamp-3 leading-relaxed">
                                        {template.description}
                                    </p>
                                {:else}
                                    <p class="text-xs text-stone-300 italic">
                                        Aucune description fournie.
                                    </p>
                                {/if}

                                <!-- Bundle specific summary -->
                                {#if template.type === 'BUNDLE' && template.data?.items}
                                    <div class="pt-2">
                                        <span class="text-[10px] uppercase font-bold text-stone-400 tracking-wider">Contenu ({template.data.items.length}) :</span>
                                        <div class="flex flex-wrap gap-1 mt-1">
                                            {#each template.data.items.slice(0, 3) as item}
                                                <span class="text-[9px] bg-stone-50 border border-stone-100 rounded px-1.5 py-0.5 text-stone-600 font-semibold">
                                                    {item.name} ({item.type})
                                                </span>
                                            {/each}
                                            {#if template.data.items.length > 3}
                                                <span class="text-[9px] text-stone-400 font-semibold align-middle pl-1">
                                                    +{template.data.items.length - 3} autres
                                                </span>
                                            {/if}
                                        </div>
                                    </div>
                                {/if}
                            </div>

                            <!-- Card Footer -->
                            <div class="pt-4 mt-4 border-t border-stone-100 flex items-center justify-between">
                                <div class="flex items-center gap-1.5">
                                    <!-- Only show edit/delete for non-virtual templates owned by the user -->
                                    {#if user && template.created_by === user.id && !template.is_virtual}
                                        <button
                                            onclick={(e) => openEditModal(template, e)}
                                            class="p-2 text-stone-400 hover:text-burnt-orange hover:bg-stone-50 rounded-lg transition-colors cursor-pointer"
                                            title="Modifier les détails"
                                        >
                                            <Edit size={14} />
                                        </button>
                                        <button
                                            onclick={(e) => handleDelete(template.id, e)}
                                            class="p-2 text-stone-400 hover:text-red-500 hover:bg-red-50 rounded-lg transition-colors cursor-pointer"
                                            title="Supprimer la publication"
                                        >
                                            <Trash2 size={14} />
                                        </button>
                                    {/if}
                                </div>

                                <div class="flex gap-1.5">
                                    {#if template.type === 'BUNDLE' && !template.is_virtual}
                                        <button
                                            onclick={(e) => { e.stopPropagation(); openBundleDetails(template); }}
                                            class="px-3 py-1.5 bg-stone-100 hover:bg-stone-200 text-stone-700 text-xs font-bold uppercase tracking-wider rounded-xl transition-all cursor-pointer"
                                        >
                                            Contenu
                                        </button>
                                    {/if}
                                    <button
                                        onclick={(e) => { e.stopPropagation(); openImportModal(template); }}
                                        class="px-4 py-1.5 bg-burnt-orange text-white text-xs font-bold uppercase tracking-wider rounded-xl shadow-xs hover:bg-opacity-95 hover:shadow-md transition-all cursor-pointer flex items-center gap-1.5"
                                    >
                                        Importer
                                        <ExternalLink size={12} />
                                    </button>
                                </div>
                            </div>
                        </div>
                    </div>
                {/each}
            </div>
        {/if}
    </main>
</div>

<!-- Import Target Selector Modal -->
{#if isImportModalOpen}
    <div
        class="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4 animate-in fade-in duration-200"
    >
        <div
            class="bg-white rounded-2xl shadow-xl w-full max-w-md overflow-hidden animate-in zoom-in-95 duration-200"
        >
            <div class="p-6 border-b border-stone-100">
                <h3 class="text-xl font-bold text-dark-gray flex items-center gap-2">
                    <FolderPlus class="text-burnt-orange" size={22} />
                    Importer dans une campagne
                </h3>
                <p class="text-sm text-stone-500 mt-1">
                    Choisissez la table de jeu de destination pour {selectedTemplateForImport?.name}.
                </p>
            </div>
            
            <div class="p-6">
                {#if myGames.length === 0}
                    <div class="bg-red-50 text-red-600 text-xs p-4 rounded-xl border border-red-100 font-semibold">
                        Vous n'êtes le Maître du Jeu d'aucune table active. Vous devez d'abord créer une table en tant que GM pour pouvoir y importer des modèles.
                    </div>
                {:else}
                    <div class="space-y-4">
                        <div>
                            <label for="import-game-select" class="block text-sm font-medium text-stone-700 mb-1">
                                Campagne (Table GM)
                            </label>
                            <select
                                id="import-game-select"
                                bind:value={targetGameId}
                                class="w-full px-4 py-2.5 rounded-xl border border-stone-200 focus:border-burnt-orange focus:ring-2 focus:ring-burnt-orange/20 outline-none transition-all text-sm font-medium"
                            >
                                {#each myGames as game}
                                    <option value={game.id}>{game.name}</option>
                                {/each}
                            </select>
                        </div>
                        
                        {#if selectedTemplateForImport?.type === 'BUNDLE'}
                            <div class="p-3 bg-stone-50 rounded-xl border border-stone-100">
                                <span class="text-[11px] font-bold text-stone-500 block">Détails du bundle :</span>
                                <span class="text-xs text-stone-400 mt-1 block">
                                    Cet import créera {selectedTemplateForImport?.data?.items?.length || 0} fiches de personnages/monstres distinctes dans votre partie en une seule opération.
                                </span>
                            </div>
                        {/if}
                    </div>
                {/if}
            </div>
            
            <div class="p-6 bg-stone-50 flex justify-end gap-3">
                <button
                    onclick={() => (isImportModalOpen = false)}
                    class="px-4 py-2 text-stone-600 font-medium hover:text-dark-gray transition-colors cursor-pointer"
                    disabled={isImporting}
                >
                    Annuler
                </button>
                {#if myGames.length > 0}
                    <button
                        onclick={handleImport}
                        disabled={isImporting}
                        class="px-4 py-2 bg-burnt-orange text-white rounded-xl font-bold shadow-md hover:bg-opacity-90 flex items-center gap-2 cursor-pointer"
                    >
                        {#if isImporting}
                            <div class="w-4 h-4 border-2 border-white/30 border-t-white rounded-full animate-spin"></div>
                        {/if}
                        Confirmer l'import
                    </button>
                {/if}
            </div>
        </div>
    </div>
{/if}

<!-- Edit Template Modal -->
{#if isEditModalOpen}
    <div
        class="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4 animate-in fade-in duration-200"
    >
        <div
            class="bg-white rounded-2xl shadow-xl w-full max-w-md overflow-hidden animate-in zoom-in-95 duration-200"
        >
            <div class="p-6 border-b border-stone-100">
                <h3 class="text-xl font-bold text-dark-gray flex items-center gap-2">
                    <Edit class="text-burnt-orange" size={20} />
                    Modifier la publication
                </h3>
                <p class="text-sm text-stone-500 mt-1">
                    Modifiez les informations de votre modèle sur le marché.
                </p>
            </div>
            
            <div class="p-6 space-y-4">
                <div>
                    <label for="edit-name" class="block text-sm font-medium text-stone-700 mb-1">
                        Nom
                    </label>
                    <input
                        id="edit-name"
                        type="text"
                        bind:value={editName}
                        class="w-full px-4 py-2 rounded-xl border border-stone-200 focus:border-burnt-orange focus:ring-2 focus:ring-burnt-orange/20 outline-none transition-all text-sm font-medium"
                    />
                </div>
                <div>
                    <label for="edit-desc" class="block text-sm font-medium text-stone-700 mb-1">
                        Description
                    </label>
                    <textarea
                        id="edit-desc"
                        bind:value={editDescription}
                        class="w-full px-4 py-2 rounded-xl border border-stone-200 focus:border-burnt-orange focus:ring-2 focus:ring-burnt-orange/20 outline-none transition-all text-sm h-24 resize-none"
                    ></textarea>
                </div>
                
                <div class="flex items-center justify-between bg-stone-50 p-3 rounded-xl border border-stone-100">
                    <div>
                        <span class="text-sm font-bold text-stone-800 block">Rendre public</span>
                        <span class="text-xs text-stone-400">Si décoché, cette création deviendra privée</span>
                    </div>
                    <label class="relative inline-flex items-center cursor-pointer">
                        <input type="checkbox" bind:checked={editIsPublic} class="sr-only peer">
                        <div class="w-11 h-6 bg-stone-200 peer-focus:outline-none rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-stone-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-burnt-orange"></div>
                    </label>
                </div>
            </div>
            
            <div class="p-6 bg-stone-50 flex justify-end gap-3">
                <button
                    onclick={() => (isEditModalOpen = false)}
                    class="px-4 py-2 text-stone-600 font-medium hover:text-dark-gray transition-colors cursor-pointer"
                    disabled={isSavingEdit}
                >
                    Annuler
                </button>
                <button
                    onclick={handleSaveEdit}
                    disabled={isSavingEdit}
                    class="px-4 py-2 bg-burnt-orange text-white rounded-xl font-bold shadow-md hover:bg-opacity-90 flex items-center gap-2 cursor-pointer"
                >
                    {#if isSavingEdit}
                        <div class="w-4 h-4 border-2 border-white/30 border-t-white rounded-full animate-spin"></div>
                    {/if}
                    Enregistrer
                </button>
            </div>
        </div>
    </div>
{/if}

<!-- Create Bundle/Pack Modal -->
{#if isCreateBundleModalOpen}
    <div
        class="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4 animate-in fade-in duration-200"
    >
        <div
            class="bg-white rounded-2xl shadow-xl w-full max-w-2xl h-[640px] flex flex-col overflow-hidden animate-in zoom-in-95 duration-200"
        >
            <div class="p-6 border-b border-stone-100 shrink-0">
                <div class="flex justify-between items-start">
                    <div>
                        <h3 class="text-xl font-bold text-dark-gray flex items-center gap-2">
                            <FolderPlus class="text-burnt-orange" size={22} />
                            Créer un nouveau Pack (Bundle)
                        </h3>
                        <p class="text-sm text-stone-500 mt-1">
                            Regroupez plusieurs créations de vos campagnes existantes dans un lot thématique.
                        </p>
                    </div>
                    <button
                        onclick={() => (isCreateBundleModalOpen = false)}
                        class="p-1 text-stone-400 hover:text-stone-600 rounded-lg transition-colors cursor-pointer"
                    >
                        ✕
                    </button>
                </div>
            </div>
            
            <div class="flex-1 overflow-y-auto p-6 space-y-6">
                <!-- Metadata section -->
                <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <div>
                        <label for="bundle-name" class="block text-xs font-bold text-stone-500 uppercase tracking-wider mb-1">
                            Nom du Lot *
                        </label>
                        <input
                            id="bundle-name"
                            type="text"
                            bind:value={bundleName}
                            placeholder="ex: Clan Gobelin de la Mine Morte"
                            class="w-full px-4 py-2 rounded-xl border border-stone-200 focus:border-burnt-orange focus:ring-2 focus:ring-burnt-orange/20 outline-none transition-all text-sm font-medium"
                        />
                    </div>
                    <div>
                        <label for="bundle-cover" class="block text-xs font-bold text-stone-500 uppercase tracking-wider mb-1">
                            Image de couverture (URL)
                        </label>
                        <input
                            id="bundle-cover"
                            type="text"
                            bind:value={bundleCoverUrl}
                            placeholder="URL de l'image (optionnel - utilisera l'avatar du 1er jeton par défaut)"
                            class="w-full px-4 py-2 rounded-xl border border-stone-200 focus:border-burnt-orange focus:ring-2 focus:ring-burnt-orange/20 outline-none transition-all text-sm"
                        />
                    </div>
                </div>

                <div>
                    <label for="bundle-desc" class="block text-xs font-bold text-stone-500 uppercase tracking-wider mb-1">
                        Description du Pack
                    </label>
                    <textarea
                        id="bundle-desc"
                        bind:value={bundleDescription}
                        placeholder="Qu'est-ce qui unit ces créatures ? Une histoire commune, une illustration d'un même style, une faction..."
                        class="w-full px-4 py-2 rounded-xl border border-stone-200 focus:border-burnt-orange focus:ring-2 focus:ring-burnt-orange/20 outline-none transition-all text-sm h-20 resize-none"
                    ></textarea>
                </div>

                <div class="flex items-center justify-between bg-stone-50 p-3 rounded-xl border border-stone-100">
                    <div>
                        <span class="text-sm font-bold text-stone-800 block">Rendre public</span>
                        <span class="text-xs text-stone-400">Si décoché, ce bundle sera stocké dans votre bibliothèque personnelle</span>
                    </div>
                    <label class="relative inline-flex items-center cursor-pointer">
                        <input type="checkbox" bind:checked={bundleIsPublic} class="sr-only peer">
                        <div class="w-11 h-6 bg-stone-200 peer-focus:outline-none rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-stone-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-burnt-orange"></div>
                    </label>
                </div>

                <!-- Item selectors -->
                <div class="space-y-3">
                    <label class="block text-xs font-bold text-stone-500 uppercase tracking-wider">
                        Sélectionnez les personnages / monstres à inclure :
                    </label>
                    
                    {#if userCharacters.length === 0 && userMonsters.length === 0}
                        <p class="text-xs text-stone-400 italic">Vous n'avez aucun personnage ou monstre dans vos parties à ajouter.</p>
                    {:else}
                        <div class="space-y-4 max-h-64 overflow-y-auto border border-stone-200 rounded-xl p-3 bg-stone-50/50">
                            <!-- User Characters -->
                            {#if userCharacters.length > 0}
                                <div>
                                    <h5 class="text-xs font-bold text-stone-400 uppercase tracking-wide mb-2">Personnages</h5>
                                    <div class="grid grid-cols-1 sm:grid-cols-2 gap-2">
                                        {#each userCharacters as char}
                                            <button
                                                type="button"
                                                onclick={() => toggleItemSelection(char.id)}
                                                class="flex items-center gap-3 p-2.5 rounded-lg border text-left transition-all bg-white hover:border-burnt-orange/50 cursor-pointer
                                                {selectedItemIds.includes(char.id) ? 'border-burnt-orange ring-1 ring-burnt-orange' : 'border-stone-200'}"
                                            >
                                                <div class="w-4 h-4 rounded border border-stone-300 flex items-center justify-center shrink-0">
                                                    {#if selectedItemIds.includes(char.id)}
                                                        <Check size={12} class="text-burnt-orange" />
                                                    {/if}
                                                </div>
                                                <img
                                                    src={char.avatar_url || `https://api.dicebear.com/9.x/avataaars/svg?seed=${char.name}`}
                                                    alt={char.name}
                                                    class="w-8 h-8 rounded-full object-cover shrink-0 bg-stone-50 border border-stone-100"
                                                />
                                                <div class="min-w-0">
                                                    <p class="text-xs font-bold text-stone-800 truncate">{char.name}</p>
                                                    <p class="text-[10px] text-stone-400 font-semibold truncate">
                                                        {char.race || 'PJ'} • Table : {char.game?.name || 'Inconnue'}
                                                    </p>
                                                </div>
                                            </button>
                                        {/each}
                                    </div>
                                </div>
                            {/if}

                            <!-- User Monsters -->
                            {#if userMonsters.length > 0}
                                <div class="mt-4">
                                    <h5 class="text-xs font-bold text-stone-400 uppercase tracking-wide mb-2">Bestiaire / Monstres</h5>
                                    <div class="grid grid-cols-1 sm:grid-cols-2 gap-2">
                                        {#each userMonsters as monst}
                                            <button
                                                type="button"
                                                onclick={() => toggleItemSelection(monst.id)}
                                                class="flex items-center gap-3 p-2.5 rounded-lg border text-left transition-all bg-white hover:border-burnt-orange/50 cursor-pointer
                                                {selectedItemIds.includes(monst.id) ? 'border-burnt-orange ring-1 ring-burnt-orange' : 'border-stone-200'}"
                                            >
                                                <div class="w-4 h-4 rounded border border-stone-300 flex items-center justify-center shrink-0">
                                                    {#if selectedItemIds.includes(monst.id)}
                                                        <Check size={12} class="text-burnt-orange" />
                                                    {/if}
                                                </div>
                                                <img
                                                    src={monst.avatar_url || `https://api.dicebear.com/9.x/bottts/svg?seed=${monst.name}`}
                                                    alt={monst.name}
                                                    class="w-8 h-8 rounded-lg object-cover shrink-0 bg-stone-50 border border-stone-100"
                                                />
                                                <div class="min-w-0">
                                                    <p class="text-xs font-bold text-stone-800 truncate">{monst.name}</p>
                                                    <p class="text-[10px] text-stone-400 font-semibold truncate font-mono">
                                                        {monst.race || 'Monstre'} • Table : {monst.game?.name || 'Inconnue'}
                                                    </p>
                                                </div>
                                            </button>
                                        {/each}
                                    </div>
                                </div>
                            {/if}
                        </div>
                    {/if}
                </div>
            </div>
            
            <div class="p-6 border-t border-stone-100 shrink-0 bg-stone-50 flex justify-between items-center">
                <span class="text-xs text-stone-500 font-semibold">
                    {selectedItemIds.length} élément(s) sélectionné(s)
                </span>
                
                <div class="flex gap-3">
                    <button
                        onclick={() => (isCreateBundleModalOpen = false)}
                        class="px-4 py-2 bg-white border border-stone-200 hover:bg-stone-50 text-stone-700 font-bold rounded-xl text-sm transition-colors cursor-pointer"
                        disabled={isCreatingBundle}
                    >
                        Annuler
                    </button>
                    <button
                        onclick={handleCreateBundle}
                        disabled={isCreatingBundle || selectedItemIds.length === 0}
                        class="px-4 py-2 bg-burnt-orange text-white rounded-xl font-bold shadow-md hover:bg-opacity-90 flex items-center gap-2 cursor-pointer disabled:opacity-50 disabled:cursor-not-allowed"
                    >
                        {#if isCreatingBundle}
                            <div class="w-4 h-4 border-2 border-white/30 border-t-white rounded-full animate-spin"></div>
                        {/if}
                        Créer le Pack
                    </button>
                </div>
            </div>
        </div>
    </div>
{/if}

<!-- Bundle Details Modal -->
{#if isBundleDetailsModalOpen && selectedBundleForDetails}
    <div
        class="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4 animate-in fade-in duration-200"
    >
        <div
            class="bg-white rounded-2xl shadow-xl w-full max-w-2xl h-[560px] flex flex-col overflow-hidden animate-in zoom-in-95 duration-200"
        >
            <div class="p-6 border-b border-stone-100 shrink-0">
                <div class="flex justify-between items-start">
                    <div class="flex items-center gap-3">
                        <div class="w-12 h-12 bg-burnt-orange/15 rounded-xl flex items-center justify-center text-burnt-orange border border-burnt-orange/20 shadow-xs shrink-0">
                            <Package size={26} />
                        </div>
                        <div class="min-w-0">
                            <h3 class="text-xl font-bold text-dark-gray truncate leading-tight">
                                {selectedBundleForDetails.name}
                            </h3>
                            <p class="text-xs text-stone-400 font-semibold mt-0.5">
                                Créé par @{selectedBundleForDetails.author_name} • {selectedBundleForDetails.uses || 0} utilisations
                            </p>
                        </div>
                    </div>
                    <button
                        onclick={() => (isBundleDetailsModalOpen = false)}
                        class="p-1 text-stone-400 hover:text-stone-600 rounded-lg transition-colors cursor-pointer"
                    >
                        ✕
                    </button>
                </div>
                
                {#if selectedBundleForDetails.description}
                    <p class="text-xs text-stone-500 mt-4 bg-stone-50 p-3 rounded-xl border border-stone-150 leading-relaxed">
                        {selectedBundleForDetails.description}
                    </p>
                {/if}
            </div>
            
            <div class="flex-1 overflow-y-auto p-6 bg-stone-50/50 space-y-3">
                <span class="text-xs font-bold text-stone-400 uppercase tracking-wider block">
                    Créatures et personnages inclus ({selectedBundleForDetails.data?.items?.length || 0}) :
                </span>
                
                <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
                    {#each selectedBundleForDetails.data?.items || [] as item, idx}
                        <div class="bg-white p-4 rounded-xl border border-stone-200 shadow-xs hover:border-burnt-orange/20 transition-all flex flex-col justify-between h-full group">
                            <div class="flex items-center gap-3">
                                <img
                                    src={item.data?.avatar_url || `https://api.dicebear.com/9.x/avataaars/svg?seed=${item.name}`}
                                    alt={item.name}
                                    class="w-12 h-12 rounded-xl object-cover bg-stone-50 border border-stone-150 shrink-0"
                                />
                                <div class="min-w-0">
                                    <h4 class="font-bold text-sm text-stone-800 truncate group-hover:text-burnt-orange transition-colors">
                                        {item.name}
                                    </h4>
                                    <div class="flex gap-1.5 mt-1">
                                        <span class="inline-block px-1.5 py-0.5 rounded text-[8px] font-bold uppercase tracking-wider bg-stone-100 text-stone-500 border border-stone-200">
                                            {item.type}
                                        </span>
                                        <span class="text-[9px] text-stone-400 font-semibold self-center">
                                            {item.data?.race || 'Race inconnue'}
                                        </span>
                                    </div>
                                </div>
                            </div>
                            
                            <div class="flex justify-between items-center mt-4 pt-3 border-t border-stone-100">
                                <div class="flex gap-2 text-[10px] text-stone-400 font-mono">
                                    <span>HP {item.data?.max_hp || 10}</span>
                                    <span>AC {item.data?.armor_class || 10}</span>
                                </div>
                                <button
                                    onclick={() => {
                                        // Construct virtual single item payload to import
                                        const virtualItem = {
                                            name: item.name,
                                            type: item.type,
                                            data: item.data,
                                            is_virtual: true,
                                            parent_bundle_id: selectedBundleForDetails!.id
                                        };
                                        openImportModal(virtualItem as any);
                                    }}
                                    class="px-2.5 py-1 bg-burnt-orange/10 hover:bg-burnt-orange hover:text-white text-burnt-orange text-[10px] font-bold rounded-lg transition-colors cursor-pointer"
                                >
                                    Importer seul
                                </button>
                            </div>
                        </div>
                    {/each}
                </div>
            </div>
            
            <div class="p-6 border-t border-stone-150 shrink-0 bg-white flex justify-end gap-3">
                <button
                    onclick={() => (isBundleDetailsModalOpen = false)}
                    class="px-4 py-2 bg-stone-100 hover:bg-stone-200 text-stone-700 font-bold rounded-xl text-sm transition-colors cursor-pointer"
                >
                    Annuler
                </button>
                <button
                    onclick={() => {
                        isBundleDetailsModalOpen = false;
                        openImportModal(selectedBundleForDetails!);
                    }}
                    class="px-4 py-2 bg-burnt-orange text-white rounded-xl font-bold shadow-md hover:bg-opacity-90 transition-all cursor-pointer flex items-center gap-1.5"
                >
                    Importer tout le pack
                    <ExternalLink size={14} />
                </button>
            </div>
        </div>
    </div>
{/if}

