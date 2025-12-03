<script lang="ts">
    import { Backpack, Coins, Search, Package } from "lucide-svelte";
    import type { Character } from "$lib/types/character";

    let { character } = $props<{ character: Character }>();

    let searchQuery = $state("");

    let filteredItems = $derived(
        character.inventory
            ? character.inventory.filter((item: any) =>
                  item.name.toLowerCase().includes(searchQuery.toLowerCase()),
              )
            : [],
    );
</script>

<div class="h-full flex flex-col">
    <!-- Gold & Search -->
    <div class="p-4 border-b border-stone-100 space-y-3">
        <div
            class="bg-yellow-50 text-yellow-700 p-3 rounded-xl flex items-center justify-between border border-yellow-100"
        >
            <div class="flex items-center gap-2 font-bold">
                <Coins size={18} />
                <span>Or</span>
            </div>
            <span class="font-mono text-lg">{character.money} po</span>
        </div>

        <div class="relative">
            <Search
                size={16}
                class="absolute left-3 top-1/2 -translate-y-1/2 text-stone-400"
            />
            <input
                type="text"
                placeholder="Rechercher..."
                bind:value={searchQuery}
                class="w-full pl-9 pr-4 py-2 bg-stone-50 border-none rounded-xl text-sm focus:ring-2 focus:ring-burnt-orange/20 focus:outline-none"
            />
        </div>
    </div>

    <!-- Items List -->
    <div class="flex-1 overflow-y-auto p-2">
        {#if filteredItems.length === 0}
            <div class="text-center py-8 text-stone-400 text-sm">
                {#if searchQuery}
                    Aucun objet trouvé pour "{searchQuery}"
                {:else}
                    Inventaire vide
                {/if}
            </div>
        {:else}
            {#each filteredItems as item}
                <div
                    class="flex items-center gap-3 p-3 hover:bg-stone-50 rounded-xl transition-colors group cursor-pointer border border-transparent hover:border-stone-100"
                >
                    <div
                        class="w-10 h-10 bg-stone-100 rounded-lg flex items-center justify-center text-xl shadow-sm overflow-hidden"
                    >
                        {#if item.image_url}
                            <img
                                src={item.image_url}
                                alt={item.name}
                                class="w-full h-full object-cover"
                            />
                        {:else if item.icon_name}
                            <!-- Assuming icon_name might be an emoji or text for now -->
                            <span class="text-lg">{item.icon_name}</span>
                        {:else}
                            <Package size={20} class="text-stone-400" />
                        {/if}
                    </div>
                    <div class="flex-1 min-w-0">
                        <div class="font-bold text-dark-gray text-sm">
                            {item.name}
                        </div>
                        {#if item.description}
                            <div class="text-xs text-stone-400 truncate">
                                {item.description}
                            </div>
                        {/if}
                    </div>
                    <div class="flex items-center gap-2">
                        <span
                            class="text-xs font-bold bg-stone-200 text-stone-600 px-2 py-1 rounded-md min-w-[24px] text-center"
                        >
                            x{item.quantity}
                        </span>
                    </div>
                </div>
            {/each}
        {/if}
    </div>

    <!-- Capacity Bar (Optional - removed for now as weight is not in standard item) -->
    <!-- 
    <div class="p-4 border-t border-stone-100 bg-stone-50">
        <div class="flex justify-between text-xs font-bold text-stone-500 mb-1">
            <span>Poids total</span>
            <span>15.5 / 60 kg</span>
        </div>
        <div class="h-2 bg-stone-200 rounded-full overflow-hidden">
            <div class="h-full bg-dark-gray w-1/4 rounded-full"></div>
        </div>
    </div>
    -->
</div>
