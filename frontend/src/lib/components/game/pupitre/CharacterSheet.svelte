<script lang="ts">
    import { Shield, Heart, Zap, Sword } from "lucide-svelte";

    let stats = {
        str: 16,
        dex: 14,
        con: 15,
        int: 10,
        wis: 12,
        cha: 8,
    };

    let skills = [
        { name: "Athlétisme", mod: "+5" },
        { name: "Perception", mod: "+3" },
        { name: "Discrétion", mod: "+4" },
    ];

    let equipment = [
        { name: "Épée longue", damage: "1d8 + 3", type: "slashing" },
        { name: "Arc court", damage: "1d6 + 2", type: "piercing" },
    ];
</script>

<div class="h-full overflow-y-auto p-4 space-y-6">
    <!-- Header Stats -->
    <div class="grid grid-cols-3 gap-3">
        <div
            class="bg-white p-3 rounded-xl border border-stone-100 shadow-sm flex flex-col items-center"
        >
            <Shield size={20} class="text-stone-400 mb-1" />
            <span class="text-2xl font-display font-bold text-dark-gray"
                >16</span
            >
            <span class="text-[10px] uppercase font-bold text-stone-400"
                >CA</span
            >
        </div>
        <div
            class="bg-white p-3 rounded-xl border border-stone-100 shadow-sm flex flex-col items-center"
        >
            <Heart size={20} class="text-red-400 mb-1" />
            <span class="text-2xl font-display font-bold text-dark-gray"
                >45</span
            >
            <span class="text-[10px] uppercase font-bold text-stone-400"
                >PV Max</span
            >
        </div>
        <div
            class="bg-white p-3 rounded-xl border border-stone-100 shadow-sm flex flex-col items-center"
        >
            <Zap size={20} class="text-yellow-400 mb-1" />
            <span class="text-2xl font-display font-bold text-dark-gray"
                >30</span
            >
            <span class="text-[10px] uppercase font-bold text-stone-400"
                >Vitesse</span
            >
        </div>
    </div>

    <!-- Ability Scores -->
    <div>
        <h3
            class="text-xs font-bold text-stone-400 uppercase tracking-wider mb-3"
        >
            Caractéristiques
        </h3>
        <div class="grid grid-cols-3 gap-2">
            {#each Object.entries(stats) as [stat, value]}
                <div class="bg-stone-50 p-2 rounded-lg text-center">
                    <div class="text-[10px] uppercase font-bold text-stone-400">
                        {stat}
                    </div>
                    <div class="font-bold text-dark-gray">{value}</div>
                    <div class="text-xs font-mono text-burnt-orange">
                        +{Math.floor((value - 10) / 2)}
                    </div>
                </div>
            {/each}
        </div>
    </div>

    <!-- Quick Actions / Weapons -->
    <div>
        <h3
            class="text-xs font-bold text-stone-400 uppercase tracking-wider mb-3"
        >
            Attaques
        </h3>
        <div class="space-y-2">
            {#each equipment as item}
                <button
                    class="w-full bg-white p-3 rounded-xl border border-stone-100 shadow-sm hover:border-burnt-orange/30 hover:shadow-md transition-all flex items-center justify-between group"
                >
                    <div class="flex items-center gap-3">
                        <div
                            class="w-8 h-8 rounded-lg bg-stone-100 flex items-center justify-center text-stone-400 group-hover:bg-burnt-orange/10 group-hover:text-burnt-orange transition-colors"
                        >
                            <Sword size={16} />
                        </div>
                        <div class="text-left">
                            <div class="font-bold text-dark-gray text-sm">
                                {item.name}
                            </div>
                            <div class="text-[10px] text-stone-400">
                                {item.type}
                            </div>
                        </div>
                    </div>
                    <div
                        class="font-mono font-bold text-dark-gray text-sm bg-stone-50 px-2 py-1 rounded"
                    >
                        {item.damage}
                    </div>
                </button>
            {/each}
        </div>
    </div>

    <!-- Skills -->
    <div>
        <h3
            class="text-xs font-bold text-stone-400 uppercase tracking-wider mb-3"
        >
            Compétences
        </h3>
        <div class="space-y-1">
            {#each skills as skill}
                <div
                    class="flex items-center justify-between p-2 hover:bg-stone-50 rounded-lg cursor-pointer transition-colors"
                >
                    <span class="text-sm font-medium text-dark-gray"
                        >{skill.name}</span
                    >
                    <span class="font-mono text-sm font-bold text-stone-500"
                        >{skill.mod}</span
                    >
                </div>
            {/each}
        </div>
    </div>
</div>
