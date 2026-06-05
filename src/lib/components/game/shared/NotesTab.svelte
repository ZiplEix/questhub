<script lang="ts">
    import { fetchNotes as fetchNotesApi, updateNotes } from "$lib/api";
    import { Loader2, CheckCircle, AlertCircle } from "lucide-svelte";

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
            console.log("NotesTab fetching notes for:", characterId);
            const res = await fetchNotesApi(gameId, characterId);
            notes = res.content;
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
            await updateNotes(gameId, characterId, notes);
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
