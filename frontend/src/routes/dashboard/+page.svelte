<script lang="ts">
    import { Dice5, LogOut, Archive, ArchiveRestore } from "lucide-svelte";
    import GameCard from "$lib/components/GameCard.svelte";

    // Mock Data
    const mockGames = [
        {
            id: "1",
            name: "La Malédiction de Strahd",
            gm: "Baptiste",
            imageUrl: "https://placehold.co/600x400/3D405B/F9F7F2?text=Strahd",
            createdAt: "2025-11-15T10:00:00Z",
            isActive: true,
        },
        {
            id: "2",
            name: "Les Mines Oubliées de Phandelver",
            gm: "Sarah",
            imageUrl:
                "https://placehold.co/600x400/E07A5F/white?text=Phandelver",
            createdAt: "2025-10-20T14:30:00Z",
            isActive: true,
        },
        {
            id: "3",
            name: "Campagne Homebrew: Aetheria",
            gm: "Baptiste",
            imageUrl:
                "https://placehold.co/600x400/F2CC8F/3D405B?text=Aetheria",
            createdAt: "2025-12-01T09:00:00Z",
            isActive: true,
        },
        {
            id: "4",
            name: "One-shot: Le Manoir Hanté",
            gm: "Marc",
            imageUrl: "https://placehold.co/600x400/darkgrey/white?text=Manoir",
            createdAt: "2025-01-10T18:00:00Z",
            isActive: false,
        },
    ];

    let showArchived = $state(false);

    let filteredGames = $derived(
        mockGames.filter((g) => (showArchived ? !g.isActive : g.isActive)),
    );
</script>

<div class="min-h-screen bg-cream">
    <!-- Dashboard Header -->
    <header class="bg-white border-b border-stone-100 sticky top-0 z-10">
        <div
            class="max-w-7xl mx-auto px-4 md:px-8 py-4 flex justify-between items-center"
        >
            <div class="flex items-center gap-2 text-dark-gray">
                <div
                    class="bg-burnt-orange text-white p-1.5 rounded-lg shadow-sm"
                >
                    <Dice5 size={20} />
                </div>
                <span class="font-display font-bold text-lg tracking-tight"
                    >QuestHub</span
                >
            </div>

            <div class="flex items-center gap-4">
                <div
                    class="flex items-center gap-3 pl-4 border-l border-stone-100"
                >
                    <div class="text-right hidden md:block">
                        <p class="text-sm font-bold text-dark-gray">Baptiste</p>
                        <p class="text-xs text-dark-gray/60">Maître du Jeu</p>
                    </div>
                    <div
                        class="w-10 h-10 rounded-full bg-mustard-yellow/30 flex items-center justify-center text-mustard-yellow font-bold border-2 border-white shadow-sm"
                    >
                        B
                    </div>
                    <button
                        class="text-dark-gray/40 hover:text-red-500 transition-colors ml-2"
                        title="Se déconnecter"
                    >
                        <LogOut size={20} />
                    </button>
                </div>
            </div>
        </div>
    </header>

    <main class="max-w-7xl mx-auto px-4 md:px-8 py-12">
        <div
            class="flex flex-col md:flex-row justify-between items-start md:items-center gap-4 mb-10"
        >
            <div>
                <h1 class="font-display font-bold text-3xl text-dark-gray mb-2">
                    Mes Parties
                </h1>
                <p class="text-dark-gray/60">
                    Gérez vos aventures et rejoignez vos tables.
                </p>
            </div>

            <div class="flex gap-3">
                <button
                    class="flex items-center gap-2 px-4 py-2 rounded-xl text-sm font-medium transition-all border {showArchived
                        ? 'bg-dark-gray text-white border-dark-gray'
                        : 'bg-white text-dark-gray border-stone-200 hover:bg-stone-50'}"
                    onclick={() => (showArchived = !showArchived)}
                >
                    {#if showArchived}
                        <ArchiveRestore size={16} />
                        Voir les parties actives
                    {:else}
                        <Archive size={16} />
                        Voir les archives
                    {/if}
                </button>

                <button
                    class="bg-burnt-orange text-white px-5 py-2.5 rounded-xl font-bold shadow-md hover:bg-opacity-90 transition-all hover:-translate-y-0.5 flex items-center gap-2"
                >
                    <Dice5 size={18} />
                    Nouvelle Partie
                </button>
            </div>
        </div>

        {#if filteredGames.length > 0}
            <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                {#each filteredGames as game (game.id)}
                    <GameCard
                        name={game.name}
                        gm={game.gm}
                        imageUrl={game.imageUrl}
                        createdAt={game.createdAt}
                        isActive={game.isActive}
                    />
                {/each}
            </div>
        {:else}
            <div
                class="text-center py-20 bg-white rounded-3xl border border-stone-100 shadow-sm"
            >
                <div
                    class="bg-stone-50 w-16 h-16 rounded-full flex items-center justify-center mx-auto mb-4 text-dark-gray/40"
                >
                    <Dice5 size={32} />
                </div>
                <h3 class="font-display font-bold text-xl text-dark-gray mb-2">
                    Aucune partie trouvée
                </h3>
                <p class="text-dark-gray/60 max-w-md mx-auto mb-6">
                    {showArchived
                        ? "Vous n'avez aucune partie archivée."
                        : "Vous n'avez pas encore de partie active. Créez-en une pour commencer l'aventure !"}
                </p>
                {#if !showArchived}
                    <button
                        class="bg-burnt-orange text-white px-6 py-3 rounded-xl font-bold shadow-md hover:bg-opacity-90 transition-all hover:-translate-y-0.5"
                    >
                        Créer une table
                    </button>
                {/if}
            </div>
        {/if}
    </main>
</div>
