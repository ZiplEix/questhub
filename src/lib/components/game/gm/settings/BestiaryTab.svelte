<script lang="ts">
    import {
        Plus,
        MoreVertical,
        Pencil,
        Trash2,
        Search,
        Skull,
        Copy,
        Share2,
        Globe,
        Package,
        Download,
        Upload
    } from "lucide-svelte";
    import { 
        fetchMonsters as fetchMonstersApi, 
        deleteCharacter,
        createTemplate,
        fetchMarketplaceTemplates,
        importTemplateToGame,
        deleteTemplate,
        updateTemplate
    } from "$lib/api";
    import { authClient } from "$lib/auth-client";
    import { onMount } from "svelte";
    import { goto } from "$app/navigation";
    import type { Character } from "$lib/types/character";

    let { gameId } = $props<{ gameId: string }>();

    let monsters = $state<Character[]>([]);
    let loading = $state(true);
    let error = $state<string | null>(null);
    let searchQuery = $state("");
    let fileInput: HTMLInputElement;

    // Sharing Modal State
    let isShareModalOpen = $state(false);
    let monsterToShare = $state<any>(null);
    let shareDescription = $state("");
    let shareIsPublic = $state(true);
    let isSharing = $state(false);

    // Import from Market Modal State
    let isImportMarketModalOpen = $state(false);
    let marketSearchQuery = $state("");
    let marketTemplates = $state<any[]>([]);
    let marketLoading = $state(false);
    let importingTemplateId = $state<string | null>(null);
    let currentUser = $state<any>(null);

    async function handleDeleteTemplate(template: any, event: Event) {
        event.stopPropagation();
        
        if (template.is_virtual) {
            if (!confirm(`Êtes-vous sûr de vouloir retirer "${template.name}" du pack "${template.parent_bundle_name}" ?`)) return;
            
            try {
                // Find parent template in marketTemplates
                const parentTemplate = marketTemplates.find(t => t.id === template.parent_bundle_id);
                if (!parentTemplate) return;
                
                const items = parentTemplate.data?.items || [];
                const match = template.id.match(/-item-(\d+)$/);
                const itemIdx = match ? parseInt(match[1]) : -1;
                
                let updatedItems;
                if (itemIdx >= 0 && itemIdx < items.length) {
                    updatedItems = items.filter((_: any, idx: number) => idx !== itemIdx);
                } else {
                    updatedItems = items.filter((item: any) => item.name !== template.name);
                }
                const updatedData = { ...parentTemplate.data, items: updatedItems };
                
                // Update parent bundle in database
                await updateTemplate(parentTemplate.id, { data: updatedData });
                
                // Reload market templates to refresh the list
                await loadMarketTemplates();
                alert("Élément retiré du pack avec succès.");
            } catch (e) {
                console.error("Failed to remove item from bundle:", e);
                alert("Erreur lors du retrait de l'élément du pack.");
            }
        } else {
            if (!confirm("Êtes-vous sûr de vouloir retirer ce modèle du marché ?")) return;
            try {
                await deleteTemplate(template.id);
                marketTemplates = marketTemplates.filter(t => t.id !== template.id);
                alert("Modèle retiré du marché avec succès.");
            } catch (e) {
                console.error("Failed to delete template:", e);
                alert("Erreur lors de la suppression du modèle.");
            }
        }
    }

    async function loadMarketTemplates() {
        try {
            marketLoading = true;
            const list = await fetchMarketplaceTemplates({
                search: marketSearchQuery
            });
            
            let flattened: any[] = [];
            for (const t of list) {
                if (t.type === 'BUNDLE') {
                    // Keep the bundle itself so GM can import the whole pack
                    flattened.push(t);
                    // Also extract and include individual monsters from the bundle
                    if (t.data?.items) {
                        t.data.items.forEach((item: any, idx: number) => {
                            if (item.type === 'MONSTRE') {
                                flattened.push({
                                    id: `${t.id}-item-${idx}`,
                                    parent_bundle_id: t.id,
                                    parent_bundle_name: t.name,
                                    created_by: t.created_by,
                                    author_name: t.author_name,
                                    name: item.name,
                                    description: item.data?.description || `Fait partie du pack "${t.name}"`,
                                    type: item.type,
                                    data: item.data,
                                    is_public: t.is_public,
                                    uses: t.uses,
                                    created_at: t.created_at,
                                    is_virtual: true
                                });
                            }
                        });
                    }
                } else if (t.type === 'MONSTRE') {
                    flattened.push(t);
                }
            }
            marketTemplates = flattened;
        } catch (e) {
            console.error("Failed to load market templates:", e);
        } finally {
            marketLoading = false;
        }
    }

    function openShareModal(monster: any) {
        monsterToShare = monster;
        shareDescription = "";
        shareIsPublic = true;
        isShareModalOpen = true;
        openMenuId = null;
    }

    async function handleShare() {
        if (!monsterToShare) return;
        try {
            isSharing = true;
            // Clean game-specific IDs
            const { id, game_id, user_id, created_at, ...cleanData } = monsterToShare;
            
            await createTemplate({
                name: monsterToShare.name,
                description: shareDescription || null,
                type: 'MONSTRE',
                data: cleanData,
                is_public: shareIsPublic
            });
            
            alert("Monstre partagé avec succès sur le Marché !");
            isShareModalOpen = false;
            monsterToShare = null;
        } catch (e) {
            console.error("Failed to share monster:", e);
            alert("Erreur lors du partage du monstre.");
        } finally {
            isSharing = false;
        }
    }

    function openImportMarketModal() {
        marketSearchQuery = "";
        marketTemplates = [];
        isImportMarketModalOpen = true;
        loadMarketTemplates();
    }

    async function handleImportTemplate(template: any) {
        try {
            importingTemplateId = template.id;
            await importTemplateToGame(template, gameId);
            alert("Création importée avec succès dans le bestiaire !");
            fetchMonsters();
            isImportMarketModalOpen = false;
        } catch (e) {
            console.error("Failed to import template:", e);
            alert("Erreur lors de l'importation.");
        } finally {
            importingTemplateId = null;
        }
    }

    // Filtered monsters
    let filteredMonsters = $derived(
        monsters.filter(
            (m) =>
                m.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
                m.race.toLowerCase().includes(searchQuery.toLowerCase()),
        ),
    );

    async function fetchMonsters() {
        loading = true;
        try {
            monsters = await fetchMonstersApi(gameId);
        } catch (e) {
            console.error("Failed to fetch monsters:", e);
            error = "Impossible de charger le bestiaire.";
        } finally {
            loading = false;
        }
    }

    async function deleteMonster(monsterId: string) {
        if (!confirm("Êtes-vous sûr de vouloir supprimer ce monstre ?")) return;

        try {
            await deleteCharacter(gameId, monsterId);
            monsters = monsters.filter((m) => m.id !== monsterId);
        } catch (e) {
            console.error("Failed to delete monster:", e);
            alert("Erreur lors de la suppression.");
        }
    }

    function openCreatePage() {
        goto(`/table/${gameId}/gm/monsters/create`);
    }

    function openEditPage(monster: Character) {
        goto(`/table/${gameId}/gm/monsters/${monster.id}/edit`);
    }

    function duplicateMonster(monster: Character) {
        localStorage.setItem("importedMonster", JSON.stringify(monster));
        goto(`/table/${gameId}/gm/monsters/create?import=true`);
    }

    async function handleImport(event: Event) {
        const target = event.target as HTMLInputElement;
        if (!target.files || target.files.length === 0) return;

        const file = target.files[0];
        const reader = new FileReader();

        reader.onload = async (e) => {
            try {
                const json = e.target?.result as string;
                const monsterData = JSON.parse(json);

                // Remove system fields to treat as new monster
                const { id, game_id, user_id, created_at, ...cleanData } = monsterData;

                // Ensure type is MONSTER
                cleanData.type = "MONSTER";

                localStorage.setItem("importedMonster", JSON.stringify(cleanData));
                goto(`/table/${gameId}/gm/monsters/create?import=true`);

                target.value = "";
            } catch (error) {
                console.error("Failed to import monster:", error);
                alert("Erreur lors de l'import du monstre. Vérifiez le format du fichier JSON.");
            }
        };

        reader.readAsText(file);
    }

    function exportMonster(monster: Character) {
        const {
            id,
            game_id,
            user_id,
            created_at,
            player_name,
            ...monsterData
        } = monster as any;

        const exportData = JSON.stringify(monsterData, null, 2);
        const blob = new Blob([exportData], { type: "application/json" });
        const url = URL.createObjectURL(blob);

        const a = document.createElement("a");
        a.href = url;
        a.download = `${(monster.name || "monstre").replace(/\s+/g, "_").toLowerCase()}.json`;
        document.body.appendChild(a);
        a.click();
        document.body.removeChild(a);
        URL.revokeObjectURL(url);

        openMenuId = null;
    }

    onMount(async () => {
        fetchMonsters();
        const { data } = await authClient.getSession();
        currentUser = data?.user || null;
    });
    let openMenuId = $state<string | null>(null);

    function toggleMenu(e: MouseEvent, id: string) {
        e.stopPropagation();
        if (openMenuId === id) {
            openMenuId = null;
        } else {
            openMenuId = id;
        }
    }

    function closeMenu() {
        openMenuId = null;
    }
</script>

<svelte:window onclick={closeMenu} />

<div class="space-y-6">
    <div class="flex items-center justify-between">
        <div class="relative flex-1 max-w-md">
            <Search
                class="absolute left-3 top-1/2 -translate-y-1/2 text-stone-400"
                size={20}
            />
            <input
                type="text"
                bind:value={searchQuery}
                placeholder="Rechercher un monstre..."
                class="w-full pl-10 pr-4 py-2 rounded-xl border border-stone-200 focus:outline-none focus:ring-2 focus:ring-burnt-orange/20 focus:border-burnt-orange transition-all"
            />
        </div>
        <div class="flex gap-2">
            <input
                bind:this={fileInput}
                type="file"
                accept=".json"
                class="hidden"
                onchange={handleImport}
            />
            <button
                onclick={() => fileInput.click()}
                class="flex items-center gap-2 px-4 py-2 bg-stone-100 text-stone-700 rounded-xl font-bold hover:bg-stone-200 transition-colors text-sm cursor-pointer"
            >
                <Upload size={18} />
                Importer
            </button>
            <button
                onclick={openImportMarketModal}
                class="flex items-center gap-2 px-4 py-2 bg-stone-900 text-white rounded-xl font-bold hover:bg-stone-800 transition-colors shadow-sm cursor-pointer"
            >
                <Globe size={18} />
                Importer du Marché
            </button>
            <button
                onclick={openCreatePage}
                class="flex items-center gap-2 px-4 py-2 bg-burnt-orange text-white rounded-xl font-bold shadow-sm hover:bg-opacity-90 transition-all hover:-translate-y-0.5 cursor-pointer"
            >
                <Plus size={20} />
                Créer un monstre
            </button>
        </div>
    </div>

    {#if loading}
        <div class="flex justify-center py-12">
            <div
                class="w-8 h-8 border-4 border-burnt-orange/30 border-t-burnt-orange rounded-full animate-spin"
            ></div>
        </div>
    {:else if error}
        <div class="text-center py-12 text-red-500 bg-red-50 rounded-xl">
            {error}
        </div>
    {:else if filteredMonsters.length === 0}
        <div
            class="text-center py-12 bg-stone-50 rounded-2xl border-2 border-dashed border-stone-200"
        >
            <div
                class="w-16 h-16 bg-stone-100 rounded-full flex items-center justify-center mx-auto mb-4"
            >
                <Skull class="text-stone-400" size={32} />
            </div>
            <h3 class="text-lg font-bold text-dark-gray mb-2">
                Le bestiaire est vide
            </h3>
            <p class="text-stone-500 mb-6 max-w-md mx-auto">
                Commencez par ajouter des monstres pour peupler votre monde.
            </p>
            <button
                onclick={openCreatePage}
                class="px-6 py-2 bg-white border border-stone-200 text-dark-gray font-bold rounded-xl hover:bg-stone-50 transition-colors shadow-sm"
            >
                Créer le premier monstre
            </button>
        </div>
    {:else}
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
            {#each filteredMonsters as monster}
                <div
                    class="bg-white p-4 rounded-xl border border-stone-200 shadow-sm hover:shadow-md transition-all group"
                >
                    <div class="flex items-start gap-4">
                        <div
                            class="w-16 h-16 rounded-xl bg-stone-100 shrink-0 overflow-hidden border border-stone-100"
                        >
                            {#if monster.avatar_url}
                                <img
                                    src={monster.avatar_url}
                                    alt={monster.name}
                                    class="w-full h-full object-cover"
                                />
                            {:else}
                                <div
                                    class="w-full h-full flex items-center justify-center"
                                >
                                    <Skull class="text-stone-300" size={32} />
                                </div>
                            {/if}
                        </div>
                        <div class="flex-1 min-w-0">
                            <div class="flex items-start justify-between">
                                <div>
                                    <h3
                                        class="font-bold text-dark-gray truncate"
                                    >
                                        {monster.name}
                                    </h3>
                                    <p class="text-sm text-stone-500 truncate">
                                        {monster.race}
                                        {#if monster.sub_race}
                                            <span class="text-stone-400"
                                                >• {monster.sub_race}</span
                                            >
                                        {/if}
                                    </p>
                                </div>
                                <div class="relative">
                                    <button
                                        onclick={(e) =>
                                            toggleMenu(e, monster.id)}
                                        class="p-1 text-stone-400 hover:text-dark-gray rounded-lg hover:bg-stone-100 transition-colors {openMenuId ===
                                        monster.id
                                            ? 'bg-stone-100 text-dark-gray'
                                            : ''}"
                                    >
                                        <MoreVertical size={20} />
                                    </button>
                                    {#if openMenuId === monster.id}
                                        <div
                                            class="absolute right-0 top-full mt-1 w-48 bg-white rounded-xl shadow-lg border border-stone-100 py-1 z-10"
                                        >
                                            <button
                                                onclick={() =>
                                                    openEditPage(monster)}
                                                class="w-full px-4 py-2 text-left text-sm text-stone-600 hover:bg-stone-50 hover:text-burnt-orange flex items-center gap-2 cursor-pointer"
                                            >
                                                <Pencil size={16} />
                                                Modifier
                                            </button>
                                            <button
                                                onclick={() =>
                                                    duplicateMonster(monster)}
                                                class="w-full px-4 py-2 text-left text-sm text-stone-600 hover:bg-stone-50 hover:text-burnt-orange flex items-center gap-2 cursor-pointer"
                                            >
                                                <Copy size={16} />
                                                Dupliquer
                                            </button>
                                            <button
                                                onclick={() => openShareModal(monster)}
                                                class="w-full px-4 py-2 text-left text-sm text-stone-600 hover:bg-stone-50 hover:text-burnt-orange flex items-center gap-2 cursor-pointer"
                                            >
                                                <Share2 size={16} class="text-burnt-orange" />
                                                Partager sur le marché
                                            </button>
                                            <button
                                                onclick={() =>
                                                    exportMonster(monster)}
                                                class="w-full px-4 py-2 text-left text-sm text-stone-600 hover:bg-stone-50 hover:text-burnt-orange flex items-center gap-2"
                                            >
                                                <Download size={16} />
                                                Exporter
                                            </button>
                                            <button
                                                onclick={() =>
                                                    deleteMonster(monster.id)}
                                                class="w-full px-4 py-2 text-left text-sm text-red-500 hover:bg-red-50 flex items-center gap-2 cursor-pointer"
                                            >
                                                <Trash2 size={16} />
                                                Supprimer
                                            </button>
                                        </div>
                                    {/if}
                                </div>
                            </div>
                            <div class="mt-3 flex items-center gap-4 text-sm">
                                <div class="flex items-center gap-1.5">
                                    <div
                                        class="w-2 h-2 rounded-full bg-red-400"
                                    ></div>
                                    <span class="font-medium text-stone-600"
                                        >{monster.max_hp} PV</span
                                    >
                                </div>
                                <div class="flex items-center gap-1.5">
                                    <div
                                        class="w-2 h-2 rounded-full bg-blue-400"
                                    ></div>
                                    <span class="font-medium text-stone-600"
                                        >Init {monster.initiative}</span
                                    >
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
            {/each}
        </div>
    {/if}
</div>

<!-- Share Monster to Market Modal -->
{#if isShareModalOpen}
    <div
        class="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4 animate-in fade-in duration-200"
    >
        <div
            class="bg-white rounded-2xl shadow-xl w-full max-w-md overflow-hidden animate-in zoom-in-95 duration-200"
        >
            <div class="p-6 border-b border-stone-100">
                <h3 class="text-xl font-bold text-dark-gray flex items-center gap-2">
                    <Share2 class="text-burnt-orange" size={20} />
                    Partager {monsterToShare?.name}
                </h3>
                <p class="text-sm text-stone-500 mt-1">
                    Publiez ce monstre sur le Marché Communautaire.
                </p>
            </div>
            <div class="p-6 space-y-4">
                <div>
                    <label
                        for="share-desc-monster"
                        class="block text-sm font-medium text-stone-700 mb-1"
                    >
                        Description du monstre / de la créature
                    </label>
                    <textarea
                        id="share-desc-monster"
                        bind:value={shareDescription}
                        placeholder="Particularités, comportement ou histoire de ce monstre..."
                        class="w-full px-4 py-2 rounded-xl border border-stone-200 focus:border-burnt-orange focus:ring-2 focus:ring-burnt-orange/20 outline-none transition-all text-sm h-24 resize-none"
                    ></textarea>
                </div>
                
                <div class="flex items-center justify-between bg-stone-50 p-3 rounded-xl border border-stone-100">
                    <div>
                        <span class="text-sm font-bold text-stone-800 block">Rendre public</span>
                        <span class="text-xs text-stone-400">Tout le monde pourra l'importer</span>
                    </div>
                    <label class="relative inline-flex items-center cursor-pointer">
                        <input type="checkbox" bind:checked={shareIsPublic} class="sr-only peer">
                        <div class="w-11 h-6 bg-stone-200 peer-focus:outline-none rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-stone-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-burnt-orange"></div>
                    </label>
                </div>
            </div>
            <div class="p-6 bg-stone-50 flex justify-end gap-3">
                <button
                    onclick={() => (isShareModalOpen = false)}
                    class="px-4 py-2 text-stone-600 font-medium hover:text-dark-gray transition-colors cursor-pointer"
                    disabled={isSharing}
                >
                    Annuler
                </button>
                <button
                    onclick={handleShare}
                    disabled={isSharing}
                    class="px-4 py-2 bg-burnt-orange text-white rounded-xl font-medium shadow-md hover:bg-opacity-90 transition-all flex items-center gap-2 cursor-pointer"
                >
                    {#if isSharing}
                        <div class="w-4 h-4 border-2 border-white/30 border-t-white rounded-full animate-spin"></div>
                    {/if}
                    Partager
                </button>
            </div>
        </div>
    </div>
{/if}

<!-- Import Monster from Market Modal -->
{#if isImportMarketModalOpen}
    <div
        class="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4 animate-in fade-in duration-200"
    >
        <div
            class="bg-white rounded-2xl shadow-xl w-full max-w-2xl h-[600px] flex flex-col overflow-hidden animate-in zoom-in-95 duration-200"
        >
            <div class="p-6 border-b border-stone-100 shrink-0">
                <div class="flex justify-between items-start">
                    <div>
                        <h3 class="text-xl font-bold text-dark-gray flex items-center gap-2">
                            <Globe class="text-stone-700" size={20} />
                            Importer des Monstres du Marché
                        </h3>
                        <p class="text-sm text-stone-500 mt-1">
                            Recherchez et importez des monstres ou des bundles communautaires.
                        </p>
                    </div>
                    <button
                        onclick={() => (isImportMarketModalOpen = false)}
                        class="p-1 text-stone-400 hover:text-stone-600 rounded-lg transition-colors cursor-pointer"
                    >
                        ✕
                    </button>
                </div>
                
                <!-- Search bar -->
                <div class="mt-4 relative">
                    <Search class="absolute left-3 top-1/2 -translate-y-1/2 text-stone-400" size={16} />
                    <input
                        type="text"
                        bind:value={marketSearchQuery}
                        oninput={loadMarketTemplates}
                        placeholder="Rechercher par nom, auteur, description..."
                        class="w-full pl-9 pr-4 py-2.5 bg-stone-50 border border-stone-200 rounded-xl text-sm focus:outline-none focus:border-burnt-orange transition-colors"
                    />
                </div>
            </div>
            
            <!-- Templates List -->
            <div class="flex-1 overflow-y-auto p-6 bg-stone-50/50 space-y-4">
                {#if marketLoading}
                    <div class="flex flex-col items-center justify-center py-20 gap-3">
                        <div class="w-8 h-8 border-4 border-burnt-orange/30 border-t-burnt-orange rounded-full animate-spin"></div>
                        <p class="text-sm text-stone-500">Recherche de monstres...</p>
                    </div>
                {:else if marketTemplates.length === 0}
                    <div class="text-center py-20 text-stone-400 flex flex-col items-center gap-2">
                        <Package size={32} class="text-stone-300" />
                        <p class="text-sm font-semibold">Aucun modèle ou bundle trouvé.</p>
                        <p class="text-xs text-stone-400">Essayez d'autres mots-clés ou partagez d'abord vos monstres.</p>
                    </div>
                {:else}
                    <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
                        {#each marketTemplates as template}
                            <div class="bg-white p-4 rounded-xl border border-stone-200 shadow-sm hover:shadow-md transition-all flex flex-col justify-between h-full group relative {template.is_virtual ? 'pt-8' : ''}">
                                {#if currentUser && template.created_by === currentUser.id}
                                    <button
                                        onclick={(e) => handleDeleteTemplate(template, e)}
                                        class="absolute top-2 right-2 p-1.5 bg-white/90 hover:bg-red-55 text-stone-400 hover:text-red-500 rounded-lg shadow-xs border border-stone-100 transition-all cursor-pointer z-20"
                                        title={template.is_virtual ? "Retirer le pack parent du marché" : "Retirer du marché"}
                                    >
                                        <Trash2 size={12} />
                                    </button>
                                {/if}
                                {#if template.is_virtual}
                                    <div class="absolute top-0 right-0 left-0 bg-stone-900/90 text-white text-[9px] font-bold py-1 px-3 flex items-center justify-between rounded-t-xl shadow-xs backdrop-blur-md z-10 font-sans">
                                        <span class="truncate">📦 Pack : {template.parent_bundle_name}</span>
                                    </div>
                                {/if}
                                <div class="space-y-3">
                                    <div class="flex items-start justify-between gap-3">
                                        <div class="flex items-center gap-3">
                                            {#if template.type === 'BUNDLE'}
                                                <div class="w-12 h-12 rounded-xl bg-burnt-orange/10 flex items-center justify-center text-burnt-orange border border-burnt-orange/20 shrink-0">
                                                    <Package size={24} />
                                                </div>
                                            {:else}
                                                <img
                                                    src={template.data?.avatar_url || `https://api.dicebear.com/9.x/bottts/svg?seed=${template.name}`}
                                                    alt={template.name}
                                                    class="w-12 h-12 rounded-xl object-cover bg-stone-50 border-2 border-stone-100 shrink-0"
                                                />
                                            {/if}
                                            <div class="min-w-0">
                                                <h4 class="font-bold text-stone-800 truncate group-hover:text-burnt-orange transition-colors">
                                                    {template.name}
                                                </h4>
                                                <span class="inline-block px-2 py-0.5 rounded-full text-[10px] font-bold uppercase tracking-wider bg-stone-100 text-stone-500 border border-stone-200">
                                                    {template.type === 'BUNDLE' ? 'PACK' : template.type}
                                                </span>
                                            </div>
                                        </div>
                                    </div>
                                    
                                    {#if template.description}
                                        <p class="text-xs text-stone-500 line-clamp-2 leading-relaxed">
                                            {template.description}
                                        </p>
                                    {/if}
                                    
                                    <div class="flex items-center justify-between text-[10px] text-stone-400 pt-2 border-t border-stone-100">
                                        <span>par @{template.author_name}</span>
                                        <span>{template.uses || 0} util.</span>
                                    </div>
                                </div>
                                
                                <button
                                    onclick={() => handleImportTemplate(template)}
                                    disabled={importingTemplateId === template.id}
                                    class="mt-4 w-full py-1.5 bg-burnt-orange/10 hover:bg-burnt-orange hover:text-white text-burnt-orange text-xs font-bold rounded-lg transition-all flex items-center justify-center gap-2 cursor-pointer disabled:opacity-50"
                                >
                                    {#if importingTemplateId === template.id}
                                        <div class="w-3.5 h-3.5 border-2 border-burnt-orange/30 border-t-burnt-orange rounded-full animate-spin"></div>
                                    {/if}
                                    Importer dans le Bestiaire
                                </button>
                            </div>
                        {/each}
                    </div>
                {/if}
            </div>
            
            <div class="p-6 border-t border-stone-100 shrink-0 bg-white flex justify-end">
                <button
                    onclick={() => (isImportMarketModalOpen = false)}
                    class="px-4 py-2 bg-stone-100 hover:bg-stone-200 text-stone-700 font-bold rounded-xl text-sm transition-colors cursor-pointer"
                >
                    Fermer
                </button>
            </div>
        </div>
    </div>
{/if}
