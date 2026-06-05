<script lang="ts">
    let messages = $state([
        {
            id: 1,
            sender: "GM",
            text: "Bienvenue dans les Ruines de l'Aube !",
            type: "text",
            time: "20:01",
        },
        {
            id: 2,
            sender: "Eldric",
            text: "Je lance une perception.",
            type: "text",
            time: "20:02",
        },
        {
            id: 3,
            sender: "Eldric",
            result: 18,
            formula: "1d20 + 4",
            type: "roll",
            time: "20:02",
        },
        {
            id: 4,
            sender: "GM",
            text: "Tu remarques une légère brise venant du couloir nord.",
            type: "text",
            time: "20:03",
        },
    ]);

    let newMessage = $state("");

    function sendMessage() {
        if (!newMessage.trim()) return;
        messages = [
            ...messages,
            {
                id: Date.now(),
                sender: "Moi",
                text: newMessage,
                type: "text",
                time: new Date().toLocaleTimeString([], {
                    hour: "2-digit",
                    minute: "2-digit",
                }),
            },
        ];
        newMessage = "";
    }

    function handleKeydown(e: KeyboardEvent) {
        if (e.key === "Enter" && !e.shiftKey) {
            e.preventDefault();
            sendMessage();
        }
    }
</script>

<div class="flex flex-col h-full">
    <div class="flex-1 overflow-y-auto p-4 space-y-4">
        {#each messages as msg}
            <div
                class="flex flex-col gap-1 {msg.sender === 'Moi'
                    ? 'items-end'
                    : 'items-start'}"
            >
                <div class="flex items-baseline gap-2">
                    <span class="text-xs font-bold text-stone-500"
                        >{msg.sender}</span
                    >
                    <span class="text-[10px] text-stone-400">{msg.time}</span>
                </div>

                {#if msg.type === "roll"}
                    <div
                        class="bg-stone-100 border border-stone-200 rounded-xl p-3 min-w-[150px]"
                    >
                        <div class="text-xs text-stone-500 mb-1">
                            Jet de dés
                        </div>
                        <div class="flex items-center gap-3">
                            <div
                                class="bg-white w-8 h-8 rounded-lg shadow-sm flex items-center justify-center font-bold text-dark-gray border border-stone-100"
                            >
                                {msg.result}
                            </div>
                            <div class="text-xs font-mono text-stone-400">
                                {msg.formula}
                            </div>
                        </div>
                    </div>
                {:else}
                    <div
                        class="px-4 py-2 rounded-2xl text-sm max-w-[90%]
                        {msg.sender === 'Moi'
                            ? 'bg-burnt-orange text-white rounded-tr-none'
                            : 'bg-white border border-stone-100 text-dark-gray rounded-tl-none shadow-sm'}"
                    >
                        {msg.text}
                    </div>
                {/if}
            </div>
        {/each}
    </div>

    <div class="p-3 border-t border-stone-100 bg-white">
        <div class="relative">
            <input
                type="text"
                bind:value={newMessage}
                onkeydown={handleKeydown}
                placeholder="Message..."
                class="w-full pl-4 pr-10 py-2 bg-stone-50 border-none rounded-xl focus:ring-2 focus:ring-burnt-orange/20 focus:outline-none text-sm"
            />
            <button
                onclick={sendMessage}
                class="absolute right-2 top-1/2 -translate-y-1/2 p-1 text-burnt-orange hover:bg-burnt-orange/10 rounded-lg transition-colors"
                aria-label="Envoyer le message"
            >
                <svg
                    xmlns="http://www.w3.org/2000/svg"
                    width="16"
                    height="16"
                    viewBox="0 0 24 24"
                    fill="none"
                    stroke="currentColor"
                    stroke-width="2"
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    ><path d="m22 2-7 20-4-9-9-4Z" /><path
                        d="M22 2 11 13"
                    /></svg
                >
            </button>
        </div>
    </div>
</div>
