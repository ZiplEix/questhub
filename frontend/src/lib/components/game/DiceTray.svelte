<script lang="ts">
    import { Dices } from "lucide-svelte";

    let customFormula = $state("");

    const dice = [
        {
            label: "d4",
            value: 4,
            color: "bg-red-100 text-red-700 border-red-200",
        },
        {
            label: "d6",
            value: 6,
            color: "bg-blue-100 text-blue-700 border-blue-200",
        },
        {
            label: "d8",
            value: 8,
            color: "bg-green-100 text-green-700 border-green-200",
        },
        {
            label: "d10",
            value: 10,
            color: "bg-purple-100 text-purple-700 border-purple-200",
        },
        {
            label: "d12",
            value: 12,
            color: "bg-orange-100 text-orange-700 border-orange-200",
        },
        {
            label: "d20",
            value: 20,
            color: "bg-yellow-100 text-yellow-700 border-yellow-200",
        },
    ];

    let { onRoll } = $props();

    function roll(sides: number) {
        const result = Math.floor(Math.random() * sides) + 1;
        onRoll?.(result);
    }
</script>

<div
    class="bg-white border-t border-stone-200 p-4 shadow-[0_-4px_6px_-1px_rgba(0,0,0,0.05)]"
>
    <!-- Quick Dice Grid -->
    <div class="grid grid-cols-3 gap-2 mb-3">
        {#each dice as die}
            <button
                onclick={() => roll(die.value)}
                class="{die.color} border font-display font-bold py-2 rounded-lg hover:brightness-95 active:scale-95 transition-all text-sm"
            >
                {die.label}
            </button>
        {/each}
    </div>

    <!-- Custom Formula -->
    <div class="relative">
        <input
            type="text"
            bind:value={customFormula}
            placeholder="2d6 + 4"
            class="w-full pl-9 pr-4 py-2 bg-stone-50 border border-stone-200 rounded-xl text-sm focus:ring-2 focus:ring-burnt-orange/20 focus:outline-none focus:border-burnt-orange"
        />
        <Dices
            size={16}
            class="absolute left-3 top-1/2 -translate-y-1/2 text-stone-400"
        />
    </div>
</div>
