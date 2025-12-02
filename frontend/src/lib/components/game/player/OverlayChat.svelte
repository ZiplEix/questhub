<script lang="ts">
    import { MessageSquare, ChevronDown, ChevronUp } from "lucide-svelte";

    let collapsed = $state(false);

    let messages = $state([
        { id: 1, sender: "GM", text: "Le dragon rugit !", type: "text" },
        {
            id: 2,
            sender: "Lyra",
            text: "Je me cache derrière le rocher.",
            type: "text",
        },
        { id: 3, sender: "Moi", result: 18, label: "Discrétion", type: "roll" },
    ]);

    function toggle() {
        collapsed = !collapsed;
    }
</script>

<div
    class="bg-black/60 backdrop-blur-sm rounded-xl border border-white/10 text-white transition-all duration-300 flex flex-col
    {collapsed ? 'h-10 w-40' : 'h-64 w-80'}"
>
    <!-- Header / Toggle -->
    <button
        onclick={toggle}
        class="flex items-center justify-between p-2 hover:bg-white/5 rounded-t-xl transition-colors w-full"
    >
        <div
            class="flex items-center gap-2 text-xs font-bold text-stone-300 uppercase tracking-wider"
        >
            <MessageSquare size={14} />
            Journal
        </div>
        {#if collapsed}
            <ChevronUp size={14} />
        {:else}
            <ChevronDown size={14} />
        {/if}
    </button>

    <!-- Content -->
    {#if !collapsed}
        <div
            class="flex-1 overflow-y-auto p-2 space-y-2 text-sm scrollbar-thin scrollbar-thumb-white/20 scrollbar-track-transparent"
        >
            {#each messages as msg}
                <div class="flex flex-col gap-0.5">
                    <span class="text-[10px] font-bold text-stone-400"
                        >{msg.sender}</span
                    >
                    {#if msg.type === "roll"}
                        <div class="text-burnt-orange font-bold">
                            🎲 {msg.label}: {msg.result}
                        </div>
                    {:else}
                        <div class="text-stone-200 leading-snug">
                            {msg.text}
                        </div>
                    {/if}
                </div>
            {/each}
        </div>
        <!-- Input placeholder -->
        <div class="p-2 border-t border-white/10">
            <input
                type="text"
                placeholder="Message..."
                class="w-full bg-white/10 border-none rounded-lg px-2 py-1 text-xs text-white placeholder-stone-500 focus:ring-1 focus:ring-burnt-orange"
            />
        </div>
    {/if}
</div>
