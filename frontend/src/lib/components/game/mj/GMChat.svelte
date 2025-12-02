<script lang="ts">
    import { Send, Dices, EyeOff, MessageSquare } from "lucide-svelte";

    let messages = $state([
        {
            id: 1,
            sender: "Eldric",
            text: "Je lance une attaque !",
            type: "chat",
        },
        {
            id: 2,
            sender: "System",
            text: "🎲 Eldric : 18 (Touché)",
            type: "roll",
        },
        {
            id: 3,
            sender: "GM",
            text: "Le gobelin esquive de justesse...",
            type: "chat",
        },
        {
            id: 4,
            sender: "GM (Secret)",
            text: "🎲 Perception (Passive) : 12",
            type: "secret",
        },
    ]);

    let newMessage = $state("");
    let isSecretRoll = $state(false);
    let whisperTarget = $state("");

    function sendMessage() {
        if (!newMessage.trim()) return;

        messages.push({
            id: Date.now(),
            sender: "GM",
            text: newMessage,
            type: isSecretRoll ? "secret" : "chat",
        });
        newMessage = "";
    }
</script>

<div class="h-full flex flex-col bg-stone-50">
    <!-- Messages Area -->
    <div class="flex-1 overflow-y-auto p-4 space-y-3">
        {#each messages as msg}
            <div
                class="flex flex-col {msg.sender === 'GM' ||
                msg.sender === 'GM (Secret)'
                    ? 'items-end'
                    : 'items-start'}"
            >
                <div class="flex items-baseline gap-2 mb-1">
                    <span class="text-xs font-bold text-stone-500"
                        >{msg.sender}</span
                    >
                </div>
                <div
                    class="max-w-[85%] p-2.5 rounded-xl text-sm shadow-sm
                    {msg.type === 'secret'
                        ? 'bg-stone-200 text-stone-600 border border-stone-300 italic'
                        : msg.sender.startsWith('GM')
                          ? 'bg-burnt-orange text-white rounded-tr-none'
                          : 'bg-white text-stone-800 border border-stone-200 rounded-tl-none'}"
                >
                    {#if msg.type === "secret"}
                        <div
                            class="flex items-center gap-1 mb-1 text-xs font-bold uppercase opacity-70"
                        >
                            <EyeOff size={10} /> Secret Roll
                        </div>
                    {/if}

                    {msg.text}
                </div>
            </div>
        {/each}
    </div>

    <!-- Input Area -->
    <div class="p-3 bg-white border-t border-stone-200">
        <!-- Tools -->
        <div class="flex gap-2 mb-2">
            <button
                class="flex items-center gap-1 px-2 py-1 rounded text-xs font-bold transition-colors border
                {isSecretRoll
                    ? 'bg-stone-800 text-white border-stone-800'
                    : 'bg-stone-100 text-stone-500 border-stone-200 hover:bg-stone-200'}"
                onclick={() => (isSecretRoll = !isSecretRoll)}
                title="Lancer en secret"
            >
                <EyeOff size={12} />
                Secret
            </button>
            <select
                bind:value={whisperTarget}
                class="px-2 py-1 rounded text-xs font-bold bg-stone-100 text-stone-500 border border-stone-200 outline-none focus:border-burnt-orange"
            >
                <option value="">À tous</option>
                <option value="Eldric">Eldric</option>
                <option value="Lyra">Lyra</option>
            </select>
        </div>

        <div class="relative">
            <input
                type="text"
                bind:value={newMessage}
                onkeydown={(e) => e.key === "Enter" && sendMessage()}
                placeholder={isSecretRoll ? "Message secret..." : "Message..."}
                class="w-full pl-4 pr-10 py-2.5 bg-stone-50 border border-stone-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-burnt-orange/20 focus:border-burnt-orange transition-all"
            />
            <button
                onclick={sendMessage}
                class="absolute right-2 top-1/2 -translate-y-1/2 p-1.5 text-burnt-orange hover:bg-burnt-orange/10 rounded-lg transition-colors"
            >
                <Send size={16} />
            </button>
        </div>
    </div>
</div>
