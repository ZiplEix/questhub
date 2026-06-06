<script lang="ts">
    import { MessageSquare, Ghost, NotebookPen, Layers, Users } from "lucide-svelte";
    import QuickDice from "../player/QuickDice.svelte";

    // Components
    import GMChat from "./GMChat.svelte";
    import Bestiary from "./Bestiary.svelte";
    import NotesTab from "../shared/NotesTab.svelte";
    import BoardsTab from "../gm/settings/BoardsTab.svelte";
    import CharactersSidebarTab from "./CharactersSidebarTab.svelte";
    import { page } from "$app/state";

    let {
        currentUserId = "",
        gmCharacterId = "",
    } = $props();

    let activeTab = $state("chat");

    const tabs = [
        { id: "chat", icon: MessageSquare, label: "Chat" },
        { id: "characters", icon: Users, label: "Personnages" },
        { id: "bestiary", icon: Ghost, label: "Bestiaire" },
        { id: "notes", icon: NotebookPen, label: "Notes" },
        { id: "boards", icon: Layers, label: "Plateaux" },
    ];
</script>

<div class="h-full flex flex-col bg-white border-l border-stone-200 shadow-xl">
    <!-- Tabs Header -->
    <div class="flex p-2 gap-2 bg-white border-b border-stone-100 overflow-x-auto">
        {#each tabs as tab}
            <button
                onclick={() => (activeTab = tab.id)}
                class="flex-1 py-2 px-3 rounded-lg text-xs font-bold transition-all flex items-center justify-center gap-1.5 shrink-0
                {activeTab === tab.id
                    ? 'bg-dark-gray text-white shadow-md'
                    : 'bg-stone-50 text-stone-400 hover:bg-stone-100'}"
            >
                <tab.icon size={14} />
                <span>{tab.label}</span>
            </button>
        {/each}
    </div>

    <!-- Content -->
    <div class="flex-1 overflow-hidden relative">
        {#if activeTab === "chat"}
            <GMChat {currentUserId} />
        {:else if activeTab === "characters"}
            <CharactersSidebarTab />
        {:else if activeTab === "bestiary"}
            <Bestiary />
        {:else if activeTab === "notes"}
            <NotesTab
                characterId={gmCharacterId}
                gameId={page.params.id || ""}
            />
        {:else if activeTab === "boards"}
            <BoardsTab mode="sidebar" />
        {/if}
    </div>

    <!-- Quick Dice (Sticky Bottom) -->
    <QuickDice isGM={true} />
</div>
