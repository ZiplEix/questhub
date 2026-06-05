<script lang="ts">
    import { Swords, EyeOff, Droplets, Skull } from "lucide-svelte";

    let { entities = [], onSelect } = $props();
</script>

<div
    class="h-full flex flex-col bg-stone-800 text-stone-200 border-r border-stone-700"
>
    <!-- Entity List Header -->
    <div class="p-4 border-b border-stone-700 font-display font-bold text-center text-stone-400 tracking-wider text-sm">
        INITIATIVE ET MONSTRES
    </div>

    <!-- Entity List -->
    <div class="flex-1 overflow-y-auto p-2 space-y-2">
        {#each entities as entity}
            <button
                class="w-full text-left p-2 rounded-lg border-l-4 bg-stone-700/50 hover:bg-stone-700 transition-colors
                {entity.type === 'player'
                    ? 'border-blue-500'
                    : entity.type === 'monster'
                      ? 'border-red-500'
                      : 'border-stone-500'}"
                onclick={() => onSelect?.(entity)}
            >
                <div class="flex justify-between items-start mb-1">
                    <span class="font-bold text-sm truncate">{entity.name}</span>
                    <span class="text-xs font-mono text-stone-400"
                        >Init: {entity.init}</span
                    >
                </div>

                <!-- HP & Status -->
                <div class="flex items-center justify-between gap-2">
                    <!-- svelte-ignore a11y_click_events_have_key_events -->
                    <!-- svelte-ignore a11y_no_static_element_interactions -->
                    <div
                        class="flex items-center bg-black/30 rounded px-1.5 py-0.5 text-xs font-mono"
                        onclick={(e) => e.stopPropagation()}
                        aria-label="HP"
                    >
                        <input
                            type="number"
                            value={entity.hp}
                            class="w-6 bg-transparent text-right outline-none {entity.hp <
                            entity.maxHp / 3
                                ? 'text-red-400'
                                : 'text-green-400'}"
                        />
                        <span class="text-stone-500 mx-1">/</span>
                        <span class="text-stone-400">{entity.maxHp}</span>
                    </div>

                    <div class="flex gap-1">
                        {#if entity.status && entity.status.includes("invisible")}
                            <EyeOff size={12} class="text-blue-400" />
                        {/if}
                        {#if entity.status && entity.status.includes("dead")}
                            <Skull size={12} class="text-stone-500" />
                        {/if}
                    </div>
                </div>
            </button>
        {/each}
    </div>
</div>
