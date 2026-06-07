<script lang="ts">
    import QuickDice from "./QuickDice.svelte";
    import NotesTab from "../shared/NotesTab.svelte";
    import StorySidebarTab from "../shared/StorySidebarTab.svelte";

    // Shared components from GM view
    import CharacterSheet from "../pupitre/CharacterSheet.svelte";
    import Inventory from "../pupitre/Inventory.svelte";
    import type { Character } from "$lib/types/character";
    import { authClient } from "$lib/auth-client";
    // Removed unused imports

    import { onMount } from "svelte";
    import Chat from "../shared/Chat.svelte";
    import PlayerTracker from "./PlayerTracker.svelte";
    import { activeBoardStore } from "$lib/websocket";

    let {
        character,
        players = [],
        currentUserId = "",
    } = $props<{
        character: Character;
        players?: any[];
        currentUserId?: string;
    }>();

    let activeTab = $state("fiche");

    const tabs = [
        { id: "fiche", label: "Fiche" },
        { id: "sac", label: "Sac" },
        { id: "notes", label: "Notes" },
        { id: "story", label: "Histoire" },
        { id: "chat", label: "Chat" },
        { id: "combat", label: "Combat" },
    ];

    $effect(() => {
        if ($activeBoardStore?.combat_active) {
            activeTab = "combat";
        }
    });
</script>

<div
    class="h-full flex flex-col bg-stone-50 border-l border-stone-200 shadow-xl z-20"
>
    <!-- 1. Tabs Navigation (Top) -->
    <div class="flex p-2 gap-2 bg-white border-b border-stone-100">
        {#each tabs as tab}
            <button
                onclick={() => (activeTab = tab.id)}
                class="flex-1 py-2 rounded-lg text-[11px] md:text-xs lg:text-sm font-bold transition-all
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
        {:else if activeTab === "sac"}
            <Inventory {character} />
        {:else if activeTab === "notes"}
            <NotesTab characterId={character.id} gameId={character.game_id} />
        {:else if activeTab === "story"}
            <StorySidebarTab gameId={character.game_id} isGM={false} />
        {:else if activeTab === "chat"}
            <Chat {players} {currentUserId} senderName={character.name} />
        {:else if activeTab === "combat"}
            <PlayerTracker />
        {/if}
    </div>

    <!-- 4. Quick Dice (Sticky Bottom) -->
    <QuickDice />
</div>
