<script lang="ts">
    import { getContext } from "svelte";
    import { rollDice } from "$lib/api";
    import { page } from "$app/state";
    import { Eye, EyeOff } from "lucide-svelte";
    import { parseDiceAndMath } from "$lib/utils/diceParser";
    import { sendMessage } from "$lib/chat";
    import { supabase } from "$lib/supabaseClient";
    import { slide } from "svelte/transition";

    const { showToast } = getContext<any>("toast");

    const tableCtx = getContext<{ isReadOnly: boolean }>("tableContext");
    const isReadOnly = $derived(tableCtx?.isReadOnly || false);

    let { isGM = false } = $props<{ isGM?: boolean }>();
    let isSecret = $state(false);
    let showCustomInput = $state(false);
    let customFormula = $state("");

    const dice = [
        { label: "d4", val: 4 },
        { label: "d6", val: 6 },
        { label: "d8", val: 8 },
        { label: "d10", val: 10 },
        { label: "d20", val: 20 },
    ];

    async function roll(sides: number) {
        if (isReadOnly) return;
        const gameId = page.params.id;
        if (!gameId) return;
        try {
            const response = await rollDice(gameId, sides, isSecret);
            const result = response.result;
            showToast(
                `🎲 d${sides} : ${result} ${isSecret ? "(Secret)" : ""}`,
                "roll",
            );
        } catch (e) {
            console.error(e);
            showToast(`Erreur de dé`, "error");
        }
    }

    async function handleCustomRoll() {
        if (isReadOnly) return;
        const gameId = page.params.id;
        if (!gameId) return;

        const formula = customFormula.trim();
        if (!formula) return;

        const parseResult = parseDiceAndMath(formula);
        if (!parseResult) {
            showToast("Formule de dés invalide", "error");
            return;
        }

        try {
            const { data: { user } } = await supabase.auth.getUser();
            if (!user) return;

            let msgType = "EVENT";
            let targetId: string | undefined = undefined;

            if (isSecret && isGM) {
                msgType = "CHAT_PRIVATE";
                targetId = user.id;
            }

            const formattedMsg = `🎲 Jet : \`${formula}\` ➔ **${parseResult.logText}**`;

            await sendMessage({
                game_id: gameId,
                content: formattedMsg,
                type: msgType,
                sender_name: isGM ? "GM" : "Joueur",
                target_id: targetId
            });

            showToast(
                `🎲 Custom (${formula}) : ${parseResult.evaluatedTotal} ${isSecret && isGM ? "(Secret)" : ""}`,
                "roll"
            );

            customFormula = "";
            showCustomInput = false;
        } catch (e) {
            console.error(e);
            showToast("Erreur de dé", "error");
        }
    }
</script>

<div
    class="bg-white border-t border-stone-200 p-3 shadow-[0_-4px_6px_-1px_rgba(0,0,0,0.05)] space-y-2"
>
    <div class="flex justify-between gap-1 items-center">
        <!-- Secret Toggle for GM -->
        {#if isGM}
            <button
                onclick={() => (isSecret = !isSecret)}
                disabled={isReadOnly}
                class="p-2 rounded-lg border transition-all mr-1
                {isSecret
                    ? 'bg-stone-800 text-white border-stone-800'
                    : 'bg-stone-50 text-stone-400 border-stone-200 hover:bg-stone-100'}
                disabled:opacity-50 disabled:cursor-not-allowed"
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
                disabled={isReadOnly}
                class="flex-1 bg-stone-50 border border-stone-200 text-stone-600 font-display font-bold py-2 rounded-lg hover:bg-burnt-orange hover:text-white hover:border-burnt-orange transition-all text-xs active:scale-95 disabled:opacity-50 disabled:cursor-not-allowed"
            >
                {die.label}
            </button>
        {/each}

        <!-- Custom button -->
        <button
            onclick={() => showCustomInput = !showCustomInput}
            disabled={isReadOnly}
            class="flex-1 border text-stone-600 font-display font-bold py-2 rounded-lg hover:bg-burnt-orange hover:text-white hover:border-burnt-orange transition-all text-[10px] uppercase active:scale-95 disabled:opacity-50 disabled:cursor-not-allowed
            {showCustomInput
                ? 'bg-burnt-orange text-white border-burnt-orange'
                : 'bg-stone-50 border-stone-200'}"
            title="Jet de dés personnalisé"
        >
            Custom
        </button>
    </div>

    {#if showCustomInput && !isReadOnly}
        <div class="flex gap-1.5 items-center pt-1" transition:slide|local>
            <input
                type="text"
                bind:value={customFormula}
                disabled={isReadOnly}
                placeholder="Ex: 2d6+4, d100, 3d10-2..."
                class="flex-1 px-3 py-1.5 text-xs border border-stone-200 rounded-lg focus:outline-none focus:border-burnt-orange bg-white font-mono disabled:opacity-50 disabled:cursor-not-allowed"
                onkeydown={(e) => {
                    if (e.key === "Enter") handleCustomRoll();
                    if (e.key === "Escape") {
                        showCustomInput = false;
                        customFormula = "";
                    }
                }}
                autofocus
            />
            <button
                onclick={handleCustomRoll}
                disabled={isReadOnly}
                class="py-1.5 px-3 bg-burnt-orange hover:bg-burnt-orange-dark text-white font-bold text-xs rounded-lg active:scale-95 transition-all disabled:opacity-50 disabled:cursor-not-allowed"
            >
                Lancer
            </button>
        </div>
    {/if}
</div>
