<script lang="ts">
    import { Plus, MoreVertical, Pencil, Trash2 } from "lucide-svelte";
    import CharacterCreationModal from "$lib/components/game/gm/CharacterCreationModal.svelte";
    import { api } from "$lib/api";
    import { authClient } from "$lib/auth-client";
    import { page } from "$app/state";

    let { characters, onRefresh } = $props<{
        characters: any[];
        onRefresh: () => void;
    }>();

    let isCharacterModalOpen = $state(false);
    let selectedCharacter = $state<any>(null);
    let openMenuId = $state<string | null>(null);

    function handleCharacterCreated() {
        onRefresh();
        selectedCharacter = null;
    }

    function openEditModal(character: any) {
        selectedCharacter = character;
        isCharacterModalOpen = true;
        openMenuId = null;
    }

    function openCreateModal() {
        selectedCharacter = null;
        isCharacterModalOpen = true;
    }

    async function deleteCharacter(character: any) {
        if (
            !confirm(
                `Êtes-vous sûr de vouloir supprimer le personnage ${character.name} ?`,
            )
        ) {
            return;
        }

        const gameId = page.params.id;
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
                onRefresh();
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
</script>

<svelte:window onclick={handleClickOutside} />

<div class="space-y-6 animate-in fade-in duration-300">
    <div class="flex justify-between items-center mb-4">
        <h3 class="text-lg font-bold text-dark-gray">
            Personnages ({characters.length})
        </h3>
        <button
            onclick={openCreateModal}
            class="flex items-center gap-2 px-4 py-2 bg-burnt-orange text-white rounded-xl font-medium shadow-md hover:bg-opacity-90 transition-all hover:-translate-y-0.5 text-sm"
        >
            <Plus size={18} />
            Créer un personnage
        </button>
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
                            <button
                                onclick={() => openEditModal(character)}
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

<CharacterCreationModal
    isOpen={isCharacterModalOpen}
    onClose={() => (isCharacterModalOpen = false)}
    onCharacterCreated={handleCharacterCreated}
    character={selectedCharacter}
/>
