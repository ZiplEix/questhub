<script lang="ts">
    import { getContext } from "svelte";

    const { showToast } = getContext<any>("toast");

    const dice = [
        { label: "d4", val: 4 },
        { label: "d6", val: 6 },
        { label: "d8", val: 8 },
        { label: "d10", val: 10 },
        { label: "d12", val: 12 },
        { label: "d20", val: 20 },
    ];

    function roll(sides: number) {
        const result = Math.floor(Math.random() * sides) + 1;
        showToast(`🎲 d${sides} : ${result}`, "roll");
    }
</script>

<div
    class="bg-white border-t border-stone-200 p-3 shadow-[0_-4px_6px_-1px_rgba(0,0,0,0.05)]"
>
    <div class="flex justify-between gap-1">
        {#each dice as die}
            <button
                onclick={() => roll(die.val)}
                class="flex-1 bg-stone-50 border border-stone-200 text-stone-600 font-display font-bold py-2 rounded-lg hover:bg-burnt-orange hover:text-white hover:border-burnt-orange transition-all text-xs active:scale-95"
            >
                {die.label}
            </button>
        {/each}
    </div>
</div>
