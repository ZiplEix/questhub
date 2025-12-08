<script lang="ts">
    import { api } from "$lib/api";
    import { authClient } from "$lib/auth-client";
    import { Loader2, CheckCircle, AlertCircle } from "lucide-svelte";
    import { onMount } from "svelte";
    import { page } from "$app/state";

    let { characterId, gameId } = $props<{
        characterId: string;
        gameId: string;
    }>();

    let notes = $state("");
    let saveStatus = $state<"saved" | "saving" | "error">("saved");
    let timeout: ReturnType<typeof setTimeout>;

    async function fetchNotes() {
        if (!characterId) return;
        try {
            const { data: tokenData } = await authClient.token();
            const token = tokenData?.token;
            if (!token) return;

            console.log("NotesTab fetching notes for:", characterId);
            const response = await api.get(
                `/table/${gameId}/characters/${characterId}/notes`,
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
    }

    $effect(() => {
        if (characterId) {
            fetchNotes();
        }
    });

    function handleInput() {
        saveStatus = "saving";
        clearTimeout(timeout);
        timeout = setTimeout(saveNotes, 1000);
    }

    async function saveNotes() {
        if (!characterId) return;
        try {
            const { data: tokenData } = await authClient.token();
            const token = tokenData?.token;
            if (!token) return;

            await api.put(
                `/table/${gameId}/characters/${characterId}/notes`,
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

<div class="p-4 h-full flex flex-col gap-2">
    <div class="flex justify-between items-center">
        <h3 class="font-bold text-dark-gray">Notes</h3>
        {#if saveStatus === "saving"}
            <span class="text-stone-400 flex items-center gap-1 text-xs"
                ><Loader2 size={14} class="animate-spin" /> Sauvegarde...</span
            >
        {:else if saveStatus === "saved"}
            <span class="text-green-500 flex items-center gap-1 text-xs"
                ><CheckCircle size={14} /> Enregistré</span
            >
        {:else if saveStatus === "error"}
            <span class="text-red-500 flex items-center gap-1 text-xs"
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
