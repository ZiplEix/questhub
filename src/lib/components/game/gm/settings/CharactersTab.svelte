<script lang="ts">
    import {
        Plus,
        MoreVertical,
        Pencil,
        Trash2,
        UserPlus,
        UserMinus,
        Download,
        Upload,
        Share2,
        Globe,
        Search,
        Package
    } from "lucide-svelte";
    import { 
        fetchCharacters, 
        assignCharacter as assignCharacterApi, 
        deleteCharacter as deleteCharacterApi,
        createTemplate,
        fetchMarketplaceTemplates,
        importTemplateToGame
    } from "$lib/api";
    import { page } from "$app/state";
    import { goto } from "$app/navigation";

    let { players, gameId } = $props<{
        players: any[];
        gameId: string;
    }>();

    let characters = $state<any[]>([]);
    let fileInput: HTMLInputElement;

    let openMenuId = $state<string | null>(null);

    // Assignment Modal State
    let isAssignModalOpen = $state(false);
    let characterToAssign = $state<any>(null);
    let selectedPlayerId = $state<string>("");

    // Sharing Modal State
    let isShareModalOpen = $state(false);
    let characterToShare = $state<any>(null);
    let shareDescription = $state("");
    let shareIsPublic = $state(true);
    let isSharing = $state(false);

    // Import from Market Modal State
    let isImportMarketModalOpen = $state(false);
    let marketSearchQuery = $state("");
    let marketTemplates = $state<any[]>([]);
    let marketLoading = $state(false);
    let importingTemplateId = $state<string | null>(null);

    async function loadCharacters() {
        try {
            characters = await fetchCharacters(gameId);
        } catch (error) {
            console.error("Failed to load characters:", error);
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
                    // Also extract and include individual characters/PNJs from the bundle
                    if (t.data?.items) {
                        t.data.items.forEach((item: any, idx: number) => {
                            if (item.type !== 'MONSTRE') {
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
                } else if (t.type !== 'MONSTRE') {
                    flattened.push(t);
                }
            }
            marketTemplates = flattened;
        } catch (error) {
            console.error("Failed to load market templates:", error);
        } finally {
            marketLoading = false;
        }
    }

    function openShareModal(character: any) {
        characterToShare = character;
        shareDescription = "";
        shareIsPublic = true;
        isShareModalOpen = true;
        openMenuId = null;
    }

    async function handleShare() {
        if (!characterToShare) return;
        try {
            isSharing = true;
            // Clean game-specific IDs
            const { id, game_id, user_id, created_at, ...cleanData } = characterToShare;
            
            await createTemplate({
                name: characterToShare.name,
                description: shareDescription || null,
                type: characterToShare.is_npc ? 'PNJ' : 'PERSONNAGE',
                data: cleanData,
                is_public: shareIsPublic
            });
            
            alert("Personnage partagé avec succès sur le Marché !");
            isShareModalOpen = false;
            characterToShare = null;
        } catch (error) {
            console.error("Failed to share character:", error);
            alert("Erreur lors du partage du personnage.");
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
            alert("Personnage importé avec succès !");
            loadCharacters();
            isImportMarketModalOpen = false;
        } catch (error) {
            console.error("Failed to import template:", error);
            alert("Erreur lors de l'importation.");
        } finally {
            importingTemplateId = null;
        }
    }

    function openEditPage(character: any) {
        goto(`/table/${gameId}/gm/characters/${character.id}/edit`);
    }

    function openCreatePage() {
        goto(`/table/${gameId}/gm/characters/create`);
    }

    function openAssignModal(character: any) {
        characterToAssign = character;
        selectedPlayerId = ""; // Reset selection
        isAssignModalOpen = true;
        openMenuId = null;
    }

    async function handleImport(event: Event) {
        const target = event.target as HTMLInputElement;
        if (!target.files || target.files.length === 0) return;

        const file = target.files[0];
        const reader = new FileReader();

        reader.onload = async (e) => {
            try {
                const json = e.target?.result as string;
                const charData = JSON.parse(json);

                // Remove system fields to treat as new character
                const { id, game_id, user_id, created_at, ...cleanData } =
                    charData;

                // Navigate to create page with imported data
                // We'll use localStorage for simplicity to pass data to the new page
                localStorage.setItem(
                    "importedCharacter",
                    JSON.stringify(cleanData),
                );
                goto(`/table/${gameId}/gm/characters/create?import=true`);

                target.value = ""; // Reset input
            } catch (error) {
                console.error("Failed to import character:", error);
                alert("Erreur lors de l'import du personnage.");
            }
        };

        reader.readAsText(file);
    }

    function exportCharacter(character: any) {
        // Create a clean copy of the character data for export
        const {
            id,
            game_id,
            user_id,
            created_at,
            player_name,
            ...characterData
        } = character;

        const exportData = JSON.stringify(characterData, null, 2);
        const blob = new Blob([exportData], { type: "application/json" });
        const url = URL.createObjectURL(blob);

        const a = document.createElement("a");
        a.href = url;
        a.download = `${(character.name || "personnage").replace(/\s+/g, "_").toLowerCase()}.json`;
        document.body.appendChild(a);
        a.click();
        document.body.removeChild(a);
        URL.revokeObjectURL(url);

        openMenuId = null;
    }

    async function assignCharacter() {
        if (!characterToAssign || !selectedPlayerId) return;

        try {
            await assignCharacterApi(gameId, characterToAssign.id, selectedPlayerId);
            loadCharacters();
            isAssignModalOpen = false;
            characterToAssign = null;
        } catch (error) {
            console.error("Failed to assign character:", error);
        }
    }

    async function unassignCharacter(character: any) {
        if (
            !confirm(
                `Êtes-vous sûr de vouloir désassigner le personnage ${character.name || ""} ?`,
            )
        ) {
            return;
        }

        try {
            await assignCharacterApi(gameId, character.id, null);
            loadCharacters();
            openMenuId = null;
        } catch (error) {
            console.error("Failed to unassign character:", error);
        }
    }

    async function deleteCharacter(character: any) {
        if (
            !confirm(
                `Êtes-vous sûr de vouloir supprimer le personnage ${character.name || ""} ?`,
            )
        ) {
            return;
        }

        try {
            await deleteCharacterApi(gameId, character.id);
            loadCharacters();
        } catch (error) {
            console.error("Failed to delete character:", error);
        }
    }

    function toggleMenu(id: string, event: Event) {
        event.stopPropagation();
        if (openMenuId === id) {
            openMenuId = null;
        } else {
            openMenuId = id;
        }
    }

    // Close menu when clicking outside
    function handleClickOutside(event: MouseEvent) {
        if (openMenuId) {
            openMenuId = null;
        }
    }

    $effect(() => {
        loadCharacters();
    });
</script>

<svelte:window onclick={handleClickOutside} />

<div class="space-y-6 animate-in fade-in duration-300">
    <div class="flex justify-between items-center mb-4">
        <h3 class="text-lg font-bold text-dark-gray">
            Personnages ({characters.length})
        </h3>
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
                Importer (Fichier)
            </button>
            <button
                onclick={openImportMarketModal}
                class="flex items-center gap-2 px-4 py-2 bg-stone-900 text-white rounded-xl font-bold hover:bg-stone-800 transition-colors text-sm cursor-pointer"
            >
                <Globe size={18} />
                Importer du Marché
            </button>
            <button
                onclick={openCreatePage}
                class="flex items-center gap-2 px-4 py-2 bg-burnt-orange text-white rounded-xl font-medium shadow-md hover:bg-opacity-90 transition-all hover:-translate-y-0.5 text-sm cursor-pointer"
            >
                <Plus size={18} />
                Créer un personnage
            </button>
        </div>
    </div>

    <div class="space-y-3">
        {#each characters as character}
            <div
                class="flex items-center justify-between p-4 rounded-xl border border-stone-100 hover:border-stone-200 transition-all relative"
            >
                <div class="flex items-center gap-4">
                    {#if character.avatar_url}
                        <img
                            src={character.avatar_url}
                            alt={character.name || "avatar"}
                            class="w-10 h-10 rounded-full object-cover"
                        />
                    {:else}
                        <div
                            class="w-10 h-10 rounded-full bg-stone-100 flex items-center justify-center text-stone-500 font-bold"
                        >
                            {(character.name || "?").charAt(0)}
                        </div>
                    {/if}
                    <div>
                        <div class="flex items-center gap-2">
                            <p class="font-bold text-dark-gray">
                                {character.name || "Sans nom"}
                            </p>
                            <span
                                class="text-xs px-2 py-0.5 bg-stone-100 text-stone-600 rounded-full border border-stone-200"
                            >
                                {character.race || "Inconnue"}
                            </span>
                        </div>
                        <p class="text-xs text-stone-500">
                            {#if character.is_npc}
                                PNJ
                            {:else if character.user_id && players.find((p: any) => p.user_id === character.user_id)}
                                Joué par {players.find((p: any) => p.user_id === character.user_id).name}
                            {:else}
                                Non assigné
                            {/if}
                        </p>
                    </div>
                </div>

                <!-- Menu Button -->
                <div class="relative">
                    <button
                        onclick={(e) => toggleMenu(character.id, e)}
                        class="p-2 text-stone-400 hover:text-dark-gray hover:bg-stone-100 rounded-lg transition-all"
                    >
                        <MoreVertical size={18} />
                    </button>

                    {#if openMenuId === character.id}
                        <div
                            class="absolute right-0 mt-2 w-48 bg-white rounded-xl shadow-lg border border-stone-100 py-1 z-20 animate-in fade-in zoom-in duration-100 origin-top-right"
                        >
                            {#if !character.is_npc}
                                <button
                                    onclick={() => openAssignModal(character)}
                                    class="w-full px-4 py-2 text-left text-sm text-stone-600 hover:bg-stone-50 hover:text-dark-gray flex items-center gap-2"
                                >
                                    <UserPlus size={16} />
                                    Assigner à un joueur
                                </button>
                                {#if character.user_id}
                                    <button
                                        onclick={() =>
                                            unassignCharacter(character)}
                                        class="w-full px-4 py-2 text-left text-sm text-stone-600 hover:bg-stone-50 hover:text-dark-gray flex items-center gap-2"
                                    >
                                        <UserMinus size={16} />
                                        Désassigner
                                    </button>
                                {/if}
                            {/if}
                            <button
                                onclick={() => exportCharacter(character)}
                                class="w-full px-4 py-2 text-left text-sm text-stone-600 hover:bg-stone-50 hover:text-dark-gray flex items-center gap-2 cursor-pointer"
                            >
                                <Download size={16} />
                                Exporter (JSON)
                            </button>
                            <button
                                onclick={() => openShareModal(character)}
                                class="w-full px-4 py-2 text-left text-sm text-stone-600 hover:bg-stone-50 hover:text-dark-gray flex items-center gap-2 cursor-pointer"
                            >
                                <Share2 size={16} class="text-burnt-orange" />
                                Partager sur le marché
                            </button>
                            <button
                                onclick={() => openEditPage(character)}
                                class="w-full px-4 py-2 text-left text-sm text-stone-600 hover:bg-stone-50 hover:text-dark-gray flex items-center gap-2"
                            >
                                <Pencil size={16} />
                                Modifier
                            </button>
                            <button
                                onclick={() => deleteCharacter(character)}
                                class="w-full px-4 py-2 text-left text-sm text-red-500 hover:bg-red-50 flex items-center gap-2"
                            >
                                <Trash2 size={16} />
                                Supprimer
                            </button>
                        </div>
                    {/if}
                </div>
            </div>
        {/each}
        {#if characters.length === 0}
            <div
                class="text-center py-12 text-stone-500 bg-stone-50 rounded-xl border border-stone-100 border-dashed"
            >
                <p>Aucun personnage créé pour le moment.</p>
            </div>
        {/if}
    </div>
</div>

<!-- Assignment Modal -->
{#if isAssignModalOpen}
    <div
        class="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4 animate-in fade-in duration-200"
    >
        <div
            class="bg-white rounded-2xl shadow-xl w-full max-w-md overflow-hidden animate-in zoom-in-95 duration-200"
        >
            <div class="p-6 border-b border-stone-100">
                <h3 class="text-xl font-bold text-dark-gray">
                    Assigner {characterToAssign?.name}
                </h3>
                <p class="text-sm text-stone-500 mt-1">
                    Sélectionnez un joueur pour ce personnage.
                </p>
            </div>
            <div class="p-6">
                <div class="space-y-4">
                    <div>
                        <label
                            for="player-select"
                            class="block text-sm font-medium text-stone-700 mb-1"
                        >
                            Joueur
                        </label>
                        <select
                            id="player-select"
                            bind:value={selectedPlayerId}
                            class="w-full px-4 py-2 rounded-xl border border-stone-200 focus:border-burnt-orange focus:ring-2 focus:ring-burnt-orange/20 outline-none transition-all"
                        >
                            <option value="" disabled selected>
                                Choisir un joueur...
                            </option>
                            {#each players as player}
                                <option value={player.user_id}
                                    >{player.name}</option
                                >
                            {/each}
                        </select>
                    </div>
                </div>
            </div>
            <div class="p-6 bg-stone-50 flex justify-end gap-3">
                <button
                    onclick={() => (isAssignModalOpen = false)}
                    class="px-4 py-2 text-stone-600 font-medium hover:text-dark-gray transition-colors cursor-pointer"
                >
                    Annuler
                </button>
                <button
                    onclick={assignCharacter}
                    disabled={!selectedPlayerId}
                    class="px-4 py-2 bg-burnt-orange text-white rounded-xl font-medium shadow-md hover:bg-opacity-90 transition-all disabled:opacity-50 disabled:cursor-not-allowed cursor-pointer"
                >
                    Assigner
                </button>
            </div>
        </div>
    </div>
{/if}

<!-- Share to Market Modal -->
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
                    Partager {characterToShare?.name}
                </h3>
                <p class="text-sm text-stone-500 mt-1">
                    Publiez ce personnage sur le Marché Communautaire.
                </p>
            </div>
            <div class="p-6 space-y-4">
                <div>
                    <label
                        for="share-desc"
                        class="block text-sm font-medium text-stone-700 mb-1"
                    >
                        Description de la création
                    </label>
                    <textarea
                        id="share-desc"
                        bind:value={shareDescription}
                        placeholder="Expliquez ce qui rend ce personnage unique (classes, particularités, histoire...)"
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

<!-- Import from Market Modal -->
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
                            Importer depuis le Marché Communautaire
                        </h3>
                        <p class="text-sm text-stone-500 mt-1">
                            Recherchez et importez des personnages créés par la communauté.
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
                        <p class="text-sm text-stone-500">Recherche de créations...</p>
                    </div>
                {:else if marketTemplates.length === 0}
                    <div class="text-center py-20 text-stone-400 flex flex-col items-center gap-2">
                        <Package size={32} class="text-stone-300" />
                        <p class="text-sm font-semibold">Aucun modèle trouvé.</p>
                        <p class="text-xs text-stone-400">Essayez d'autres mots-clés ou partagez d'abord vos personnages.</p>
                    </div>
                {:else}
                    <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
                        {#each marketTemplates as template}
                            <div class="bg-white p-4 rounded-xl border border-stone-200 shadow-sm hover:shadow-md transition-all flex flex-col justify-between h-full group relative {template.is_virtual ? 'pt-8' : ''}">
                                {#if template.is_virtual}
                                    <div class="absolute top-0 right-0 left-0 bg-stone-900/90 text-white text-[9px] font-bold py-1 px-3 flex items-center justify-between rounded-t-xl shadow-xs backdrop-blur-md z-10">
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
                                                    src={template.data?.avatar_url || `https://api.dicebear.com/9.x/avataaars/svg?seed=${template.name}`}
                                                    alt={template.name}
                                                    class="w-12 h-12 rounded-xl object-cover bg-stone-50 border-2 border-stone-100 shrink-0"
                                                />
                                            {/if}
                                            <div class="min-w-0">
                                                <h4 class="font-bold text-stone-800 truncate group-hover:text-burnt-orange transition-colors">
                                                    {template.name}
                                                </h4>
                                                <span class="inline-block px-2 py-0.5 rounded-full text-[10px] font-bold uppercase tracking-wider bg-stone-100 text-stone-500 border border-stone-200">
                                                    {template.type === 'BUNDLE' ? 'PACK / BUNDLE' : template.type}
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
                                    Importer dans cette table
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

