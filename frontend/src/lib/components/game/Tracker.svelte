<script lang="ts">
    import { Users, Swords, Skull, Heart, ShieldAlert } from "lucide-svelte";

    // Mock data for development
    let mode = $state<"exploration" | "combat">("exploration");

    let players = $state([
        {
            id: 1,
            name: "Eldric",
            avatar: "https://api.dicebear.com/7.x/avataaars/svg?seed=Eldric",
            hp: 45,
            maxHp: 50,
            initiative: 12,
            conditions: [],
        },
        {
            id: 2,
            name: "Lyra",
            avatar: "https://api.dicebear.com/7.x/avataaars/svg?seed=Lyra",
            hp: 28,
            maxHp: 35,
            initiative: 18,
            conditions: ["poisoned"],
        },
        {
            id: 3,
            name: "Thorn",
            avatar: "https://api.dicebear.com/7.x/avataaars/svg?seed=Thorn",
            hp: 60,
            maxHp: 60,
            initiative: 8,
            conditions: [],
        },
    ]);

    let activeEntityId = $state(2); // Mock active entity in combat

    function toggleMode() {
        mode = mode === "exploration" ? "combat" : "exploration";
    }
</script>

<div class="h-full bg-white border-r border-stone-200 flex flex-col">
    <!-- Header / Mode Switcher -->
    <div class="p-4 border-b border-stone-100">
        <button
            onclick={toggleMode}
            class="w-full flex items-center justify-center gap-2 py-2 px-4 rounded-xl font-display font-bold transition-all duration-200
            {mode === 'exploration'
                ? 'bg-stone-100 text-dark-gray hover:bg-stone-200'
                : 'bg-red-50 text-red-600 hover:bg-red-100'}"
        >
            {#if mode === "exploration"}
                <Users size={18} />
                <span>Exploration</span>
            {:else}
                <Swords size={18} />
                <span>Combat</span>
            {/if}
        </button>
    </div>

    <!-- List -->
    <div class="flex-1 overflow-y-auto p-2 space-y-2">
        {#if mode === "exploration"}
            <div
                class="px-2 py-1 text-xs font-bold text-stone-400 uppercase tracking-wider"
            >
                Joueurs connectés
            </div>
            {#each players as player}
                <div
                    class="flex items-center gap-3 p-2 rounded-xl hover:bg-stone-50 transition-colors group cursor-pointer"
                >
                    <img
                        src={player.avatar}
                        alt={player.name}
                        class="w-10 h-10 rounded-full bg-stone-200 border-2 border-white shadow-sm"
                    />
                    <div class="flex-1 min-w-0">
                        <div class="font-bold text-dark-gray truncate">
                            {player.name}
                        </div>
                        <div
                            class="text-xs text-stone-500 flex items-center gap-1"
                        >
                            <div
                                class="w-2 h-2 rounded-full bg-green-500"
                            ></div>
                            En ligne
                        </div>
                    </div>
                </div>
            {/each}
        {:else}
            <div
                class="px-2 py-1 text-xs font-bold text-stone-400 uppercase tracking-wider"
            >
                Initiative
            </div>
            {#each players.sort((a, b) => b.initiative - a.initiative) as player}
                <div
                    class="relative p-2 rounded-xl transition-all duration-200 border border-transparent
                    {activeEntityId === player.id
                        ? 'bg-white shadow-md border-burnt-orange/30 scale-[1.02] z-10'
                        : 'hover:bg-stone-50'}"
                >
                    {#if activeEntityId === player.id}
                        <div
                            class="absolute left-0 top-1/2 -translate-y-1/2 -translate-x-1 w-1.5 h-8 bg-burnt-orange rounded-r-full"
                        ></div>
                    {/if}

                    <div class="flex items-center gap-3 mb-2">
                        <img
                            src={player.avatar}
                            alt={player.name}
                            class="w-10 h-10 rounded-full bg-stone-200 border-2 border-white shadow-sm"
                        />
                        <div class="flex-1 min-w-0">
                            <div class="flex items-center justify-between">
                                <span class="font-bold text-dark-gray truncate"
                                    >{player.name}</span
                                >
                                <span
                                    class="text-xs font-mono font-bold text-stone-400 bg-stone-100 px-1.5 py-0.5 rounded"
                                >
                                    {player.initiative}
                                </span>
                            </div>
                            <!-- Conditions -->
                            {#if player.conditions.length > 0}
                                <div class="flex gap-1 mt-0.5">
                                    {#each player.conditions as condition}
                                        <span
                                            class="text-[10px] bg-purple-100 text-purple-700 px-1.5 rounded-full flex items-center gap-0.5"
                                        >
                                            <ShieldAlert size={8} />
                                            {condition}
                                        </span>
                                    {/each}
                                </div>
                            {/if}
                        </div>
                    </div>

                    <!-- HP Bar -->
                    <div
                        class="relative h-2 bg-stone-100 rounded-full overflow-hidden"
                    >
                        <div
                            class="absolute top-0 left-0 h-full rounded-full transition-all duration-500
                            {player.hp / player.maxHp < 0.3
                                ? 'bg-red-500'
                                : 'bg-green-500'}"
                            style="width: {(player.hp / player.maxHp) * 100}%"
                        ></div>
                    </div>
                    <div
                        class="flex justify-between text-[10px] font-bold text-stone-400 mt-1 px-0.5"
                    >
                        <span>{player.hp} PV</span>
                        <span>{player.maxHp} Max</span>
                    </div>
                </div>
            {/each}
        {/if}
    </div>

    <!-- Footer Actions (GM only usually) -->
    <div class="p-4 border-t border-stone-100">
        <button
            class="w-full py-2 text-sm font-bold text-stone-400 hover:text-dark-gray transition-colors flex items-center justify-center gap-2"
        >
            <Skull size={16} />
            Ajouter un monstre
        </button>
    </div>
</div>
