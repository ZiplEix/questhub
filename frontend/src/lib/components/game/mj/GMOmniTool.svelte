<script lang="ts">
    import { MessageSquare, Ghost, Map, NotebookPen } from "lucide-svelte";

    // Components
    import GMChat from "./GMChat.svelte";
    import Bestiary from "./Bestiary.svelte";
    import SceneSelector from "./SceneSelector.svelte";
    import Notes from "../pupitre/Notes.svelte";

    let { onSpawn, entities = [], currentUserId = "" } = $props();

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
    <div class="flex border-b border-stone-200 bg-stone-50">
        {#each tabs as tab}
            <button
                onclick={() => (activeTab = tab.id)}
                class="flex-1 py-3 flex flex-col items-center justify-center gap-1 text-xs font-bold uppercase tracking-wider transition-colors border-b-2
                {activeTab === tab.id
                    ? 'border-burnt-orange text-burnt-orange bg-white'
                    : 'border-transparent text-stone-400 hover:text-stone-600 hover:bg-stone-100'}"
            >
                <tab.icon size={18} />
                {tab.label}
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
            <Notes />
        {/if}
    </div>
</div>
