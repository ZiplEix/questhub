<script lang="ts">
    import { X, Shield, Heart, Zap, Sword } from "lucide-svelte";

    let { monster, onClose } = $props();

    // Mock full stats if not provided
    let stats = {
        str: 14,
        dex: 12,
        con: 12,
        int: 10,
        wis: 8,
        cha: 6,
    };

    let actions = [
        {
            name: "Cimetière",
            type: "Melee",
            hit: "+4",
            damage: "1d6 + 2",
            damageType: "slashing",
        },
        {
            name: "Javelot",
            type: "Ranged",
            hit: "+4",
            damage: "1d6 + 2",
            damageType: "piercing",
        },
    ];
</script>

<div
    class="fixed inset-y-0 right-[400px] w-[350px] bg-white shadow-2xl z-40 border-l border-stone-200 transform transition-transform duration-300 ease-out flex flex-col"
>
    <!-- Header -->
    <div class="p-4 bg-stone-900 text-white flex justify-between items-start">
        <div>
            <h2 class="font-display font-bold text-xl">
                {monster?.name || "Monstre"}
            </h2>
            <div class="text-xs text-stone-400 font-mono mt-1">
                Gobelin • Humanoïde (Gobelinoïde)
            </div>
        </div>
        <button
            onclick={onClose}
            class="text-stone-400 hover:text-white transition-colors"
        >
            <X size={20} />
        </button>
    </div>

    <!-- Content -->
    <div class="flex-1 overflow-y-auto p-4 space-y-6">
        <!-- Vitals -->
        <div class="grid grid-cols-3 gap-3">
            <div
                class="bg-stone-50 p-2 rounded-lg text-center border border-stone-100"
            >
                <Shield size={16} class="mx-auto text-stone-400 mb-1" />
                <div class="font-bold text-lg">{monster?.ac || 15}</div>
                <div class="text-[10px] uppercase font-bold text-stone-400">
                    CA
                </div>
            </div>
            <div
                class="bg-stone-50 p-2 rounded-lg text-center border border-stone-100"
            >
                <Heart size={16} class="mx-auto text-red-400 mb-1" />
                <div class="font-bold text-lg">{monster?.hp || 7}</div>
                <div class="text-[10px] uppercase font-bold text-stone-400">
                    PV
                </div>
            </div>
            <div
                class="bg-stone-50 p-2 rounded-lg text-center border border-stone-100"
            >
                <Zap size={16} class="mx-auto text-yellow-400 mb-1" />
                <div class="font-bold text-lg">30</div>
                <div class="text-[10px] uppercase font-bold text-stone-400">
                    VIT
                </div>
            </div>
        </div>

        <!-- Stats -->
        <div class="grid grid-cols-6 gap-1 text-center">
            {#each Object.entries(stats) as [stat, val]}
                <div class="flex flex-col">
                    <span class="text-[10px] font-bold uppercase text-stone-400"
                        >{stat}</span
                    >
                    <span class="font-bold text-stone-800">{val}</span>
                    <span class="text-[10px] text-stone-500"
                        >+{Math.floor((val - 10) / 2)}</span
                    >
                </div>
            {/each}
        </div>

        <!-- Actions -->
        <div>
            <h3
                class="text-xs font-bold text-stone-400 uppercase tracking-wider mb-3 border-b border-stone-100 pb-1"
            >
                Actions
            </h3>
            <div class="space-y-2">
                {#each actions as action}
                    <button
                        class="w-full text-left bg-stone-50 hover:bg-stone-100 p-3 rounded-lg border border-stone-200 transition-colors group"
                    >
                        <div class="flex justify-between items-center mb-1">
                            <span
                                class="font-bold text-stone-800 flex items-center gap-2"
                            >
                                <Sword size={14} class="text-burnt-orange" />
                                {action.name}
                            </span>
                            <span
                                class="text-xs font-mono bg-white px-1.5 py-0.5 rounded border border-stone-200 text-stone-500 group-hover:border-burnt-orange group-hover:text-burnt-orange transition-colors"
                            >
                                {action.hit}
                            </span>
                        </div>
                        <div class="text-xs text-stone-500 pl-6">
                            {action.type} •
                            <span class="font-mono text-stone-700"
                                >{action.damage}</span
                            >
                            {action.damageType}
                        </div>
                    </button>
                {/each}
            </div>
        </div>
    </div>
</div>
