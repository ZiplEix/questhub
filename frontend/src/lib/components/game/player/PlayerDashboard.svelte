<script lang="ts">
    import QuickDice from "./QuickDice.svelte";
    import ActionTab from "./ActionTab.svelte";

    // Shared components from GM view
    import CharacterSheet from "../pupitre/CharacterSheet.svelte";
    import Inventory from "../pupitre/Inventory.svelte";
    import type { Character } from "$lib/types/character";
    import { api } from "$lib/api";
    import { authClient } from "$lib/auth-client";
    import { Loader2, CheckCircle, AlertCircle } from "lucide-svelte";

    import { onMount } from "svelte";

    let { character } = $props<{ character: Character }>();

    let activeTab = $state("fiche");
    let notes = $state("");
    let saveStatus = $state<"saved" | "saving" | "error">("saved");
    let timeout: ReturnType<typeof setTimeout>;

    const tabs = [
        { id: "fiche", label: "Fiche" },
        { id: "action", label: "Action" },
        { id: "sac", label: "Sac" },
        { id: "notes", label: "Notes" },
    ];

    onMount(async () => {
        try {
            const { data: tokenData } = await authClient.token();
            const token = tokenData?.token;
            if (!token) return;

            const response = await api.get(
                `/table/${character.game_id}/characters/${character.id}/notes`,
                {
                    headers: {
                        Authorization: `Bearer ${token}`,
                    },
                },
            );
            notes = response.data.content;
        } catch (e) {
            console.error("Failed to fetch notes:", e);
        }
    });

    function handleInput() {
        saveStatus = "saving";
        clearTimeout(timeout);
        timeout = setTimeout(saveNotes, 1000);
    }

    async function saveNotes() {
        try {
            const { data: tokenData } = await authClient.token();
            const token = tokenData?.token;
            if (!token) return;

            await api.put(
                `/table/${character.game_id}/characters/${character.id}/notes`,
                { content: notes },
                {
                    headers: {
                        Authorization: `Bearer ${token}`,
                    },
                },
            );
            saveStatus = "saved";
        } catch (e) {
            console.error(e);
            saveStatus = "error";
        }
    }
</script>

<div
    class="h-full flex flex-col bg-stone-50 border-l border-stone-200 shadow-xl z-20"
>
    <!-- 1. Tabs Navigation (Top) -->
    <div class="flex p-2 gap-2 bg-white border-b border-stone-100">
        {#each tabs as tab}
            <button
                onclick={() => (activeTab = tab.id)}
                class="flex-1 py-2 rounded-lg text-sm font-bold transition-all
                {activeTab === tab.id
                    ? 'bg-dark-gray text-white shadow-md'
                    : 'bg-stone-50 text-stone-400 hover:bg-stone-100'}"
            >
                {tab.label}
            </button>
        {/each}
    </div>

    <!-- 2. Tab Content (Scrollable) -->
    <div class="flex-1 overflow-y-auto bg-stone-50/50 relative">
        {#if activeTab === "fiche"}
            <CharacterSheet {character} />
        {:else if activeTab === "action"}
            <ActionTab />
        {:else if activeTab === "sac"}
            <Inventory {character} />
        {:else if activeTab === "notes"}
            <div class="p-4 h-full flex flex-col gap-2">
                <div class="flex justify-between items-center">
                    <h3 class="font-bold text-dark-gray">Notes de voyage</h3>
                    {#if saveStatus === "saving"}
                        <span
                            class="text-stone-400 flex items-center gap-1 text-xs"
                            ><Loader2 size={14} class="animate-spin" /> Sauvegarde...</span
                        >
                    {:else if saveStatus === "saved"}
                        <span
                            class="text-green-500 flex items-center gap-1 text-xs"
                            ><CheckCircle size={14} /> Enregistré</span
                        >
                    {:else if saveStatus === "error"}
                        <span
                            class="text-red-500 flex items-center gap-1 text-xs"
                            ><AlertCircle size={14} /> Erreur</span
                        >
                    {/if}
                </div>
                <textarea
                    bind:value={notes}
                    oninput={handleInput}
                    class="w-full flex-1 p-3 rounded-lg border border-stone-200 resize-none focus:outline-none focus:border-burnt-orange bg-white text-dark-gray"
                    placeholder="Prenez des notes ici..."
                ></textarea>
            </div>
        {/if}
    </div>

    <!-- 4. Quick Dice (Sticky Bottom) -->
    <QuickDice />
</div>
