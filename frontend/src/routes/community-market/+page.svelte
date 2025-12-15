<script lang="ts">
    import { onMount } from "svelte";
    import Header from "$lib/components/Header.svelte";

    // Mock Data
    interface Template {
        id: string;
        name: string;
        type: "PERSONNAGE" | "PNJ" | "MONSTRE" | "CRÉATURE";
        author: string;
        uses: number;
        image: string;
        description: string;
    }

    let templates: Template[] = [
        {
            id: "1",
            name: "Eldric Shadoweaver",
            type: "PERSONNAGE",
            author: "MysticMage",
            uses: 124,
            image: "https://api.dicebear.com/9.x/avataaars/svg?seed=Eldric",
            description:
                "Un magicien de niveau 5 spécialisé dans la magie d'illusion. Inclut une liste complète de sorts et une histoire.",
        },
        {
            id: "2",
            name: "Groupe d'embuscade Gobelin",
            type: "MONSTRE",
            author: "DungeonMaster99",
            uses: 89,
            image: "https://api.dicebear.com/9.x/avataaars/svg?seed=Goblin",
            description:
                "Un groupe pré-généré de 5 gobelins avec des armes et des tactiques variées. Parfait pour les rencontres de bas niveau.",
        },
        {
            id: "3",
            name: "Seraphina Lightbringer",
            type: "PNJ",
            author: "HolyCleric",
            uses: 256,
            image: "https://api.dicebear.com/9.x/avataaars/svg?seed=Seraphina",
            description:
                "Une marchande sournoise qui vend des objets maudits. Livrée avec une liste de marchandises et des motivations cachées.",
        },
        {
            id: "4",
            name: "Troll des Cavernes",
            type: "CRÉATURE",
            author: "NatureLover",
            uses: 45,
            image: "https://api.dicebear.com/9.x/avataaars/svg?seed=Troll",
            description:
                "Troll des cavernes standard avec des capacités de régénération implémentées comme un trait personnalisé.",
        },
        {
            id: "5",
            name: "Lyra Silvertongue",
            type: "PERSONNAGE",
            author: "BardExtraordinaire",
            uses: 312,
            image: "https://api.dicebear.com/9.x/avataaars/svg?seed=Lyra",
            description:
                "Barde du Collège de l'Éloquence. Optimisée pour les rencontres sociales et le soutien.",
        },
    ];

    let searchQuery = "";
    let selectedType = "TOUT";

    $: filteredTemplates = templates.filter((t) => {
        const matchesSearch =
            t.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
            t.author.toLowerCase().includes(searchQuery.toLowerCase());
        const matchesType = selectedType === "TOUT" || t.type === selectedType;
        return matchesSearch && matchesType;
    });

    const types = ["TOUT", "PERSONNAGE", "PNJ", "MONSTRE", "CRÉATURE"];
</script>

<div class="min-h-screen bg-cream font-sans text-dark-gray">
    <Header />

    <main class="max-w-7xl mx-auto px-4 md:px-8 py-16">
        <div class="mb-12 text-center">
            <h1 class="text-4xl font-bold mb-4 font-display text-burnt-orange">
                Marché Communautaire
            </h1>
            <p class="text-xl opacity-80 max-w-2xl mx-auto">
                Découvrez et partagez des personnages, monstres et plus encore
                avec la communauté Questhub.
            </p>
        </div>

        <!-- Search & Filters -->
        <div class="mb-12 flex flex-col items-center gap-8">
            <div class="relative w-full max-w-xl group">
                <div
                    class="flex items-center gap-4 w-full px-6 py-4 bg-white rounded-full shadow-sm border border-stone-200 group-focus-within:border-burnt-orange group-focus-within:ring-4 group-focus-within:ring-burnt-orange/10 transition-all duration-300"
                >
                    <svg
                        xmlns="http://www.w3.org/2000/svg"
                        class="h-6 w-6 text-gray-400 group-focus-within:text-burnt-orange transition-colors"
                        fill="none"
                        viewBox="0 0 24 24"
                        stroke="currentColor"
                    >
                        <path
                            stroke-linecap="round"
                            stroke-linejoin="round"
                            stroke-width="2"
                            d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"
                        />
                    </svg>
                    <input
                        type="text"
                        placeholder="Rechercher des modèles..."
                        bind:value={searchQuery}
                        class="grow bg-transparent border-none outline-none text-lg text-dark-gray placeholder:text-gray-400 font-medium"
                    />
                </div>
            </div>

            <div class="flex flex-wrap justify-center gap-3">
                {#each types as type}
                    <button
                        class="px-5 py-2 rounded-full text-sm font-bold tracking-wide transition-all border
                      {selectedType === type
                            ? 'bg-dark-gray text-white border-dark-gray shadow-lg'
                            : 'bg-white text-gray-500 border-gray-100 hover:border-burnt-orange hover:text-burnt-orange'}"
                        on:click={() => (selectedType = type)}
                    >
                        {type === "TOUT" ? "Tout" : type}
                    </button>
                {/each}
            </div>
        </div>

        <!-- Grid -->
        <div
            class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-8"
        >
            {#each filteredTemplates as template}
                <div
                    class="group bg-white rounded-2xl shadow-sm hover:shadow-md transition-all duration-300 overflow-hidden border border-stone-100 flex flex-col h-full relative"
                >
                    <!-- Image Section -->
                    <div class="relative h-48 overflow-hidden">
                        <img
                            src={template.image}
                            alt={template.name}
                            class="w-full h-full object-cover transition-transform duration-500 group-hover:scale-105 bg-cream"
                        />

                        <!-- Type Badge -->
                        <div class="absolute top-3 left-3">
                            <span
                                class="bg-white/90 text-dark-gray px-3 py-1 rounded-full text-xs font-bold shadow-sm backdrop-blur-md flex items-center gap-1 border border-white/30 uppercase tracking-wider"
                            >
                                {template.type}
                            </span>
                        </div>
                    </div>

                    <!-- Content Section -->
                    <div class="p-5 flex flex-col grow gap-3">
                        <div>
                            <h3
                                class="font-display font-bold text-xl text-dark-gray mb-1 group-hover:text-burnt-orange transition-colors leading-tight"
                            >
                                {template.name}
                            </h3>
                            <div
                                class="flex items-center gap-2 text-dark-gray/60 text-sm font-medium"
                            >
                                <span>par @{template.author}</span>
                            </div>
                        </div>

                        <p
                            class="text-gray-500 text-sm line-clamp-2 leading-relaxed"
                        >
                            {template.description}
                        </p>

                        <!-- Footer -->
                        <div
                            class="mt-auto pt-4 border-t border-stone-100 flex items-center justify-between"
                        >
                            <div
                                class="flex items-center gap-2 text-dark-gray/50 text-xs font-semibold"
                            >
                                <svg
                                    xmlns="http://www.w3.org/2000/svg"
                                    class="h-3.5 w-3.5"
                                    viewBox="0 0 24 24"
                                    fill="none"
                                    stroke="currentColor"
                                    stroke-width="2"
                                    stroke-linecap="round"
                                    stroke-linejoin="round"
                                >
                                    <circle cx="12" cy="12" r="10"></circle>
                                    <polyline points="12 6 12 12 16 14"
                                    ></polyline>
                                </svg>
                                <span>{template.uses} utilisations</span>
                            </div>

                            <button
                                class="bg-burnt-orange/10 text-burnt-orange px-4 py-2 rounded-lg text-sm font-bold hover:bg-burnt-orange hover:text-white transition-all hover:cursor-pointer"
                            >
                                Voir
                            </button>
                        </div>
                    </div>
                </div>
            {/each}
        </div>

        {#if filteredTemplates.length === 0}
            <div class="text-center py-20 opacity-50">
                <p class="text-xl">
                    Aucun modèle trouvé correspondant à vos critères.
                </p>
            </div>
        {/if}
    </main>
</div>
