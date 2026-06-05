<script lang="ts">
    import { AVAILABLE_ICONS } from "$lib/utils/iconMap";
    import { Search, X } from "lucide-svelte";

    interface Props {
        value?: string;
        onSelect?: (iconName: string) => void;
    }

    const { value = "", onSelect } = $props<Props>();

    let searchQuery = $state("");
    let isOpen = $state(false);

    // Grouper les icônes par catégorie
    const groupedIcons = $derived(() => {
        const filtered = AVAILABLE_ICONS.filter((icon) =>
            [icon.name, icon.label].some((text) =>
                text.toLowerCase().includes(searchQuery.toLowerCase())
            )
        );

        const grouped = filtered.reduce(
            (acc, icon) => {
                if (!acc[icon.category]) {
                    acc[icon.category] = [];
                }
                acc[icon.category].push(icon);
                return acc;
            },
            {} as Record<string, typeof AVAILABLE_ICONS>
        );

        return Object.entries(grouped).sort(([a], [b]) => a.localeCompare(b));
    });

    const selectedIcon = $derived.by(() => {
        return AVAILABLE_ICONS.find((icon) => icon.name === value);
    });

    function handleSelect(iconName: string) {
        onSelect?.(iconName);
        isOpen = false;
        searchQuery = "";
    }

    function handleClear() {
        onSelect?.("");
        isOpen = false;
    }
</script>

<div class="relative">
    <!-- Input/Button -->
    <button
        onclick={() => (isOpen = !isOpen)}
        class="w-full px-3 py-2 rounded-lg border border-stone-200 text-sm focus:outline-none focus:border-burnt-orange text-left flex items-center justify-between bg-white hover:bg-stone-50 transition-colors group"
    >
        <div class="flex items-center gap-2 flex-1 min-w-0">
            {#if selectedIcon}
                <div class="w-5 h-5 flex items-center justify-center flex-shrink-0">
                    <svelte:component this={selectedIcon.component} size={18} class="text-stone-600" />
                </div>
                <span class="text-stone-700 truncate">{selectedIcon.label}</span>
            {:else}
                <span class="text-stone-400">Sélectionner une icône...</span>
            {/if}
        </div>
        {#if selectedIcon}
            <button
                onclick={(e) => {
                    e.stopPropagation();
                    handleClear();
                }}
                class="p-1 text-stone-400 hover:text-red-500 rounded hover:bg-red-50 transition-colors flex-shrink-0"
            >
                <X size={16} />
            </button>
        {/if}
    </button>

    <!-- Dropdown -->
    {#if isOpen}
        <div class="absolute top-full left-0 right-0 mt-2 z-50 bg-white border border-stone-200 rounded-lg shadow-lg overflow-hidden">
            <!-- Search -->
            <div class="p-3 border-b border-stone-100 sticky top-0 bg-white">
                <div class="relative">
                    <Search
                        size={16}
                        class="absolute left-3 top-1/2 -translate-y-1/2 text-stone-400"
                    />
                    <input
                        type="text"
                        bind:value={searchQuery}
                        placeholder="Rechercher une icône..."
                        class="w-full pl-9 pr-4 py-2 bg-stone-50 border border-stone-200 rounded-lg text-sm focus:outline-none focus:border-burnt-orange transition-colors"
                        onkeydown={(e) => {
                            if (e.key === "Escape") {
                                isOpen = false;
                            }
                        }}
                    />
                </div>
            </div>

            <!-- Icons Grid -->
            <div class="max-h-80 overflow-y-auto p-3 space-y-3">
                {#if groupedIcons().length === 0}
                    <div class="text-center py-8 text-stone-400 text-sm">
                        Aucune icône trouvée
                    </div>
                {:else}
                    {#each groupedIcons() as [category, icons]}
                        <div>
                            <h3 class="text-xs font-bold text-stone-500 uppercase tracking-wide mb-2 px-1">
                                {category}
                            </h3>
                            <div class="grid grid-cols-6 gap-2">
                                {#each icons as icon}
                                    <button
                                        onclick={() => handleSelect(icon.name)}
                                        class="aspect-square flex items-center justify-center rounded-lg border border-stone-200 hover:bg-burnt-orange/10 hover:border-burnt-orange transition-all active:scale-95 group"
                                        title={icon.label}
                                    >
                                        <div class="flex items-center justify-center">
                                            <svelte:component
                                                this={icon.component}
                                                size={20}
                                                class="text-stone-600 group-hover:text-burnt-orange transition-colors"
                                            />
                                        </div>
                                    </button>
                                {/each}
                            </div>
                        </div>
                    {/each}
                {/if}
            </div>
        </div>
    {/if}
</div>
