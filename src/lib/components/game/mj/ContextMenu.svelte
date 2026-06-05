<script lang="ts">
    import { Skull, EyeOff, Copy, Edit, Trash2 } from "lucide-svelte";

    let { x, y, onClose } = $props();

    function handleAction(action: string) {
        console.log("Action:", action);
        onClose();
    }
</script>

<!-- Backdrop to close menu -->
<button
    class="fixed inset-0 z-50 cursor-default"
    onclick={onClose}
    oncontextmenu={(e) => {
        e.preventDefault();
        onClose();
    }}
    aria-label="Close context menu"
></button>

<div
    class="fixed z-50 bg-white rounded-lg shadow-xl border border-stone-200 py-1 w-48 animate-in fade-in zoom-in-95 duration-100"
    style="top: {y}px; left: {x}px;"
>
    <button
        onclick={() => handleAction("edit")}
        class="w-full text-left px-3 py-2 text-sm hover:bg-stone-100 flex items-center gap-2 text-stone-700"
    >
        <Edit size={14} /> Modifier PV
    </button>
    <button
        onclick={() => handleAction("hide")}
        class="w-full text-left px-3 py-2 text-sm hover:bg-stone-100 flex items-center gap-2 text-stone-700"
    >
        <EyeOff size={14} /> Cacher / Révéler
    </button>
    <button
        onclick={() => handleAction("duplicate")}
        class="w-full text-left px-3 py-2 text-sm hover:bg-stone-100 flex items-center gap-2 text-stone-700"
    >
        <Copy size={14} /> Dupliquer
    </button>
    <div class="h-px bg-stone-100 my-1"></div>
    <button
        onclick={() => handleAction("kill")}
        class="w-full text-left px-3 py-2 text-sm hover:bg-red-50 flex items-center gap-2 text-red-600"
    >
        <Skull size={14} /> Tuer
    </button>
    <button
        onclick={() => handleAction("delete")}
        class="w-full text-left px-3 py-2 text-sm hover:bg-red-50 flex items-center gap-2 text-red-600"
    >
        <Trash2 size={14} /> Supprimer
    </button>
</div>
