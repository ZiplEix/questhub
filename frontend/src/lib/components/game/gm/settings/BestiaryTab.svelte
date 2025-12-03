<script lang="ts">
    import {
        Plus,
        MoreVertical,
        Pencil,
        Trash2,
        Search,
        Skull,
        Copy,
    } from "lucide-svelte";
    import { api } from "$lib/api";
    import { authClient } from "$lib/auth-client";
    import { onMount } from "svelte";
    import { goto } from "$app/navigation";
    import type { Character } from "$lib/types/character";

    let { gameId } = $props<{ gameId: string }>();

    let monsters = $state<Character[]>([]);
    let loading = $state(true);
    let error = $state<string | null>(null);
    let searchQuery = $state("");

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
            const { data: tokenData } = await authClient.token();
            if (tokenData?.token) {
                const response = await api.get(`/table/${gameId}/monsters`, {
                    headers: {
                        Authorization: `Bearer ${tokenData.token}`,
                    },
                });
                monsters = response.data;
            }
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
            const { data: tokenData } = await authClient.token();
            if (tokenData?.token) {
                await api.delete(`/table/${gameId}/characters/${monsterId}`, {
                    headers: {
                        Authorization: `Bearer ${tokenData.token}`,
                    },
                });
                monsters = monsters.filter((m) => m.id !== monsterId);
            }
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

    onMount(() => {
        fetchMonsters();
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
        <button
            onclick={openCreatePage}
            class="flex items-center gap-2 px-4 py-2 bg-burnt-orange text-white rounded-xl font-bold shadow-sm hover:bg-opacity-90 transition-all hover:-translate-y-0.5"
        >
            <Plus size={20} />
            Créer un monstre
        </button>
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
                            class="w-16 h-16 rounded-xl bg-stone-100 flex-shrink-0 overflow-hidden border border-stone-100"
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
                                                class="w-full px-4 py-2 text-left text-sm text-stone-600 hover:bg-stone-50 hover:text-burnt-orange flex items-center gap-2"
                                            >
                                                <Pencil size={16} />
                                                Modifier
                                            </button>
                                            <button
                                                onclick={() =>
                                                    duplicateMonster(monster)}
                                                class="w-full px-4 py-2 text-left text-sm text-stone-600 hover:bg-stone-50 hover:text-burnt-orange flex items-center gap-2"
                                            >
                                                <Copy size={16} />
                                                Dupliquer
                                            </button>
                                            <button
                                                onclick={() =>
                                                    deleteMonster(monster.id)}
                                                class="w-full px-4 py-2 text-left text-sm text-red-500 hover:bg-red-50 flex items-center gap-2"
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
