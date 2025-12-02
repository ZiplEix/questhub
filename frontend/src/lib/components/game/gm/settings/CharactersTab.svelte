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
    } from "lucide-svelte";
    import { api } from "$lib/api";
    import { authClient } from "$lib/auth-client";
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

    async function loadCharacters() {
        try {
            const { data: tokenData } = await authClient.token();
            const token = tokenData?.token;
            const res = await api.get(`/table/${gameId}/characters`, {
                headers: {
                    Authorization: `Bearer ${token}`,
                },
            });
            characters = res.data;
        } catch (error) {
            console.error("Failed to load characters:", error);
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
                // We'll use localStorage or similar because state might be lost or complex to handle in the form component if not designed for it.
                // Actually, let's just use the state object in goto.
                // But wait, the CharacterForm needs to read this state.
                // Since I didn't implement state reading in CharacterForm/CreatePage, I should probably do that.
                // Or easier: just pass the data via a store or simply navigate and let the user import on the create page?
                // The user didn't explicitly ask to keep import on the list page, but it's a nice feature.
                // Let's implement import on the create page instead?
                // Or, let's try to pass it via state.

                // For now, let's just move the import button to the create page?
                // No, the user might want to import from the list.

                // Let's use localStorage for simplicity to pass data to the new page
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
        a.download = `${character.name.replace(/\s+/g, "_").toLowerCase()}.json`;
        document.body.appendChild(a);
        a.click();
        document.body.removeChild(a);
        URL.revokeObjectURL(url);

        openMenuId = null;
    }

    async function assignCharacter() {
        if (!characterToAssign || !selectedPlayerId) return;

        try {
            const { data: tokenData } = await authClient.token();
            if (tokenData?.token) {
                await api.post(
                    `/table/${gameId}/characters/${characterToAssign.id}/assign`,
                    { player_id: selectedPlayerId },
                    {
                        headers: {
                            Authorization: `Bearer ${tokenData.token}`,
                        },
                    },
                );
                loadCharacters();
                isAssignModalOpen = false;
                characterToAssign = null;
            }
        } catch (error) {
            console.error("Failed to assign character:", error);
        }
    }

    async function unassignCharacter(character: any) {
        if (
            !confirm(
                `Êtes-vous sûr de vouloir désassigner le personnage ${character.name} ?`,
            )
        ) {
            return;
        }

        try {
            const { data: tokenData } = await authClient.token();
            if (tokenData?.token) {
                await api.post(
                    `/table/${gameId}/characters/${character.id}/assign`,
                    { player_id: "" }, // Empty player_id to unassign
                    {
                        headers: {
                            Authorization: `Bearer ${tokenData.token}`,
                        },
                    },
                );
                loadCharacters();
                openMenuId = null;
            }
        } catch (error) {
            console.error("Failed to unassign character:", error);
        }
    }

    async function deleteCharacter(character: any) {
        if (
            !confirm(
                `Êtes-vous sûr de vouloir supprimer le personnage ${character.name} ?`,
            )
        ) {
            return;
        }

        try {
            const { data: tokenData } = await authClient.token();
            if (tokenData?.token) {
                await api.delete(
                    `/table/${gameId}/characters/${character.id}`,
                    {
                        headers: {
                            Authorization: `Bearer ${tokenData.token}`,
                        },
                    },
                );
                loadCharacters();
            }
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
                class="flex items-center gap-2 px-4 py-2 bg-stone-100 text-stone-700 rounded-xl font-bold hover:bg-stone-200 transition-colors text-sm"
            >
                <Upload size={18} />
                Importer
            </button>
            <button
                onclick={openCreatePage}
                class="flex items-center gap-2 px-4 py-2 bg-burnt-orange text-white rounded-xl font-medium shadow-md hover:bg-opacity-90 transition-all hover:-translate-y-0.5 text-sm"
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
                            alt={character.name}
                            class="w-10 h-10 rounded-full object-cover"
                        />
                    {:else}
                        <div
                            class="w-10 h-10 rounded-full bg-stone-100 flex items-center justify-center text-stone-500 font-bold"
                        >
                            {character.name.charAt(0)}
                        </div>
                    {/if}
                    <div>
                        <div class="flex items-center gap-2">
                            <p class="font-bold text-dark-gray">
                                {character.name}
                            </p>
                            <span
                                class="text-xs px-2 py-0.5 bg-stone-100 text-stone-600 rounded-full border border-stone-200"
                            >
                                {character.race}
                            </span>
                        </div>
                        <p class="text-xs text-stone-500">
                            {#if character.is_npc}
                                PNJ
                            {:else if character.player_name}
                                Joué par {character.player_name}
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
                                class="w-full px-4 py-2 text-left text-sm text-stone-600 hover:bg-stone-50 hover:text-dark-gray flex items-center gap-2"
                            >
                                <Download size={16} />
                                Exporter
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
                    class="px-4 py-2 text-stone-600 font-medium hover:text-dark-gray transition-colors"
                >
                    Annuler
                </button>
                <button
                    onclick={assignCharacter}
                    disabled={!selectedPlayerId}
                    class="px-4 py-2 bg-burnt-orange text-white rounded-xl font-medium shadow-md hover:bg-opacity-90 transition-all disabled:opacity-50 disabled:cursor-not-allowed"
                >
                    Assigner
                </button>
            </div>
        </div>
    </div>
{/if}
