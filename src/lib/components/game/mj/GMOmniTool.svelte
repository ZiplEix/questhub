<script lang="ts">
    import { MessageSquare, Ghost, Map, NotebookPen } from "lucide-svelte";
    import QuickDice from "../player/QuickDice.svelte";

    // Components
    import GMChat from "./GMChat.svelte";
    import Bestiary from "./Bestiary.svelte";
    import SceneSelector from "./SceneSelector.svelte";
    import NotesTab from "../shared/NotesTab.svelte";
    import { page } from "$app/state";

    let {
        onSpawn,
        entities = [],
        currentUserId = "",
        gmCharacterId = "",
    } = $props();

    let activeTab = $state("chat");

    const tabs = [
        { id: "chat", icon: MessageSquare, label: "Chat" },
        { id: "bestiary", icon: Ghost, label: "Bestiaire" },
        { id: "scenes", icon: Map, label: "Scènes" },
        { id: "notes", icon: NotebookPen, label: "Notes" },
    ];
</script>

<div class="h-full flex flex-col bg-white border-l border-stone-200 shadow-xl">
    <!-- Tabs Header -->
    <div class="flex p-2 gap-2 bg-white border-b border-stone-100">
        {#each tabs as tab}
            <button
                onclick={() => (activeTab = tab.id)}
                class="flex-1 py-2 rounded-lg text-sm font-bold transition-all flex items-center justify-center gap-2
                {activeTab === tab.id
                    ? 'bg-dark-gray text-white shadow-md'
                    : 'bg-stone-50 text-stone-400 hover:bg-stone-100'}"
            >
                <tab.icon size={16} />
                <span>{tab.label}</span>
            </button>
        {/each}
    </div>

    <!-- Content -->
    <div class="flex-1 overflow-hidden relative">
        {#if activeTab === "chat"}
            <GMChat {entities} {currentUserId} />
        {:else if activeTab === "bestiary"}
            <Bestiary {onSpawn} />
        {:else if activeTab === "scenes"}
            <SceneSelector />
        {:else if activeTab === "notes"}
            <NotesTab
                characterId={gmCharacterId}
                gameId={page.params.id || ""}
            />
        {/if}
    </div>

    <!-- Quick Dice (Sticky Bottom) -->
    <QuickDice isGM={true} />
</div>
