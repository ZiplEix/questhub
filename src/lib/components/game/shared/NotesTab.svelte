<script lang="ts">
    import { fetchNotes as fetchNotesApi, updateNotes } from "$lib/api";
    import { Loader2, CheckCircle, AlertCircle } from "lucide-svelte";
    import { getContext } from "svelte";

    let { characterId, gameId } = $props<{
        characterId: string;
        gameId: string;
    }>();

    const tableCtx = getContext<{ isReadOnly: boolean; isPaused?: boolean }>("tableContext");
    const isReadOnly = $derived(tableCtx?.isReadOnly || false);
    const isPaused = $derived(tableCtx?.isPaused || false);

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
        if (isReadOnly) return;
        saveStatus = "saving";
        clearTimeout(timeout);
        timeout = setTimeout(saveNotes, 1000);
    }

    async function saveNotes() {
        if (!characterId || isReadOnly) return;
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
        {#if isReadOnly}
            <span class="text-stone-400 flex items-center gap-1 text-xs">Lecture seule</span>
        {:else if saveStatus === "saving"}
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
        disabled={isReadOnly}
        class="w-full flex-1 p-3 rounded-lg border border-stone-200 resize-none focus:outline-none focus:border-burnt-orange bg-white text-dark-gray disabled:bg-stone-50 disabled:text-stone-400"
        placeholder={isReadOnly ? (isPaused ? "Cette partie est en pause. Les notes sont en lecture seule." : "Cette partie est archivée. Les notes sont en lecture seule.") : "Prenez des notes ici..."}
    ></textarea>
</div>
