<script lang="ts">
    import { Dice5 } from "lucide-svelte";
</script>

<header
    class="w-full py-6 px-4 md:px-8 flex justify-between items-center max-w-7xl mx-auto"
>
    <div class="flex items-center gap-2 text-dark-gray">
        <div class="bg-burnt-orange text-white p-2 rounded-xl shadow-sm">
            <Dice5 size={24} />
        </div>
        <span class="font-display font-bold text-xl tracking-tight"
            >QuestHub</span
        >
    </div>

    <div class="flex items-center gap-4">
        <button
            class="font-medium text-dark-gray hover:text-burnt-orange transition-colors"
        >
            Connexion
        </button>
        <button
            class="bg-burnt-orange text-white px-5 py-2.5 rounded-xl font-medium shadow-md hover:bg-opacity-90 transition-all hover:-translate-y-0.5"
        >
            Créer une table
        </button>
    </div>
</header>
