<script lang="ts">
    import { getContext } from "svelte";
    import { api } from "$lib/api";
    import { page } from "$app/state";
    import { authClient } from "$lib/auth-client";
    import { Eye, EyeOff } from "lucide-svelte";

    const { showToast } = getContext<any>("toast");

    let { isGM = false } = $props<{ isGM?: boolean }>();
    let isSecret = $state(false);

    const dice = [
        { label: "d4", val: 4 },
        { label: "d6", val: 6 },
        { label: "d8", val: 8 },
        { label: "d10", val: 10 },
        { label: "d12", val: 12 },
        { label: "d20", val: 20 },
    ];

    async function roll(sides: number) {
        const gameId = page.params.id;
        try {
            const { data: tokenData } = await authClient.token();
            if (!tokenData?.token) return;
            const response = await api.get(
                `/table/${gameId}/roll?sides=${sides}&secret=${isSecret}`,
                {
                    headers: {
                        Authorization: `Bearer ${tokenData.token}`,
                    },
                },
            );
            const result = response.data.result;
            showToast(
                `🎲 d${sides} : ${result} ${isSecret ? "(Secret)" : ""}`,
                "roll",
            );
        } catch (e) {
            console.error(e);
            showToast(`Errore de dé`, "error");
        }
    }
</script>

<div
    class="bg-white border-t border-stone-200 p-3 shadow-[0_-4px_6px_-1px_rgba(0,0,0,0.05)]"
>
    <div class="flex justify-between gap-1 items-center">
        <!-- Secret Toggle for GM -->
        {#if isGM}
            <button
                onclick={() => (isSecret = !isSecret)}
                class="p-2 rounded-lg border transition-all mr-1
                {isSecret
                    ? 'bg-stone-800 text-white border-stone-800'
                    : 'bg-stone-50 text-stone-400 border-stone-200 hover:bg-stone-100'}"
                title="Lancé secret"
            >
                {#if isSecret}
                    <EyeOff size={16} />
                {:else}
                    <Eye size={16} />
                {/if}
            </button>
        {/if}

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
