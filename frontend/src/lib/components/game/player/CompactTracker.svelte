<script lang="ts">
    import { Swords, ChevronUp, ChevronDown } from "lucide-svelte";

    // Mock data
    let combatActive = $state(true);
    let currentTurnId = $state(1); // Player's ID
    let collapsed = $state(false);

    let initiative = $state([
        {
            id: 2,
            name: "Gobelin Chef",
            avatar: "https://api.dicebear.com/7.x/avataaars/svg?seed=Gob",
            isEnemy: true,
        },
        {
            id: 1,
            name: "Moi",
            avatar: "https://api.dicebear.com/7.x/avataaars/svg?seed=Eldric",
            isEnemy: false,
        },
        {
            id: 3,
            name: "Lyra",
            avatar: "https://api.dicebear.com/7.x/avataaars/svg?seed=Lyra",
            isEnemy: false,
        },
    ]);

    function toggle() {
        collapsed = !collapsed;
    }
</script>

{#if combatActive}
    <div
        class="bg-black/60 backdrop-blur-sm rounded-xl border border-white/10 text-white transition-all duration-300
        {collapsed ? 'w-40' : 'w-48 p-2'}"
    >
        <!-- Header / Toggle -->
        <button
            onclick={toggle}
            class="flex items-center justify-between w-full {collapsed
                ? 'p-2'
                : 'mb-2 px-1'}"
        >
            <div
                class="flex items-center gap-2 text-xs font-bold text-stone-300 uppercase tracking-wider"
            >
                <Swords size={14} />
                Initiative
            </div>
            {#if collapsed}
                <ChevronDown size={14} class="text-stone-400" />
            {:else}
                <ChevronUp size={14} class="text-stone-400" />
            {/if}
        </button>

        <!-- List Content -->
        {#if !collapsed}
            <div class="space-y-1">
                {#each initiative as entity}
                    <div
                        class="flex items-center gap-2 p-1.5 rounded-lg transition-all
                        {currentTurnId === entity.id
                            ? 'bg-burnt-orange/80 shadow-lg scale-105 border border-white/20'
                            : 'opacity-70'}"
                    >
                        <img
                            src={entity.avatar}
                            alt={entity.name}
                            class="w-6 h-6 rounded-full bg-stone-700"
                        />
                        <span class="text-sm font-bold truncate"
                            >{entity.name}</span
                        >
                        {#if currentTurnId === entity.id && entity.id === 1}
                            <span class="absolute -right-2 -top-2 flex h-3 w-3">
                                <span
                                    class="animate-ping absolute inline-flex h-full w-full rounded-full bg-burnt-orange opacity-75"
                                ></span>
                                <span
                                    class="relative inline-flex rounded-full h-3 w-3 bg-burnt-orange"
                                ></span>
                            </span>
                        {/if}
                    </div>
                {/each}
            </div>
            {#if currentTurnId === 1}
                <div
                    class="mt-2 text-center text-xs font-bold text-burnt-orange animate-pulse"
                >
                    À toi de jouer !
                </div>
            {/if}
        {/if}
    </div>
{/if}
