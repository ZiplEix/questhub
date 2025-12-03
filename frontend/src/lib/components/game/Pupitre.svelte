<script lang="ts">
    import { MessageSquare, User, Backpack, NotebookPen } from "lucide-svelte";
    import Chat from "./pupitre/Chat.svelte";
    import CharacterSheet from "./pupitre/CharacterSheet.svelte";
    import Inventory from "./pupitre/Inventory.svelte";
    import Notes from "./pupitre/Notes.svelte";
    import DiceTray from "./DiceTray.svelte";

    import type { Character } from "$lib/types/character";

    let activeTab = $state("chat");
    let { onRoll, character } = $props<{ onRoll: any; character: Character }>();

    const tabs = [
        { id: "chat", icon: MessageSquare, label: "Chat" },
        { id: "sheet", icon: User, label: "Fiche" },
        { id: "inventory", icon: Backpack, label: "Sac" },
        { id: "notes", icon: NotebookPen, label: "Notes" },
    ];
</script>

<div
    class="h-full flex flex-col bg-white border-l border-stone-200 shadow-xl z-20"
>
    <!-- Tabs Header -->
    <div class="flex border-b border-stone-200 bg-stone-50">
        {#each tabs as tab}
            <button
                onclick={() => (activeTab = tab.id)}
                class="flex-1 py-3 flex flex-col items-center justify-center gap-1 text-[10px] font-bold uppercase tracking-wider transition-all border-b-2
                {activeTab === tab.id
                    ? 'border-burnt-orange text-burnt-orange bg-white'
                    : 'border-transparent text-stone-400 hover:text-stone-600 hover:bg-stone-100'}"
            >
                <tab.icon size={18} />
                {tab.label}
            </button>
        {/each}
    </div>

    <!-- Tab Content -->
    <div class="flex-1 overflow-hidden relative">
        {#if activeTab === "chat"}
            <Chat />
        {:else if activeTab === "sheet"}
            <CharacterSheet {character} />
        {:else if activeTab === "inventory"}
            <Inventory {character} />
        {:else if activeTab === "notes"}
            <Notes />
        {/if}
    </div>

    <!-- Fixed Dice Tray -->
    <DiceTray {onRoll} />
</div>
