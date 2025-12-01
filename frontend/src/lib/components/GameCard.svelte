<script lang="ts">
    import { Calendar, User } from "lucide-svelte";

    let { name, gm, imageUrl, createdAt, isActive } = $props<{
        name: string;
        gm: string;
        imageUrl: string;
        createdAt: string;
        isActive: boolean;
    }>();

    const date = new Date(createdAt).toLocaleDateString("fr-FR", {
        year: "numeric",
        month: "long",
        day: "numeric",
    });
</script>

<div
    class="group bg-white rounded-2xl shadow-sm hover:shadow-md transition-all duration-300 overflow-hidden border border-stone-100 flex flex-col h-full"
>
    <div class="relative h-48 overflow-hidden">
        <img
            src={imageUrl}
            alt={name}
            class="w-full h-full object-cover transition-transform duration-500 group-hover:scale-105"
        />
        {#if !isActive}
            <div
                class="absolute inset-0 bg-dark-gray/60 flex items-center justify-center backdrop-blur-[1px]"
            >
                <span
                    class="bg-white/20 text-white px-3 py-1 rounded-full text-sm font-medium backdrop-blur-md border border-white/30"
                    >Archivée</span
                >
            </div>
        {/if}
    </div>

    <div class="p-5 flex flex-col flex-grow gap-4">
        <div>
            <h3
                class="font-display font-bold text-xl text-dark-gray mb-1 group-hover:text-burnt-orange transition-colors"
            >
                {name}
            </h3>
            <div class="flex items-center gap-2 text-dark-gray/60 text-sm">
                <User size={14} />
                <span>MJ: {gm}</span>
            </div>
        </div>

        <div
            class="mt-auto pt-4 border-t border-stone-100 flex items-center justify-between"
        >
            <div class="flex items-center gap-2 text-dark-gray/50 text-xs">
                <Calendar size={12} />
                <span>Créée le {date}</span>
            </div>

            <button
                class="bg-burnt-orange/10 text-burnt-orange px-4 py-2 rounded-lg text-sm font-bold hover:bg-burnt-orange hover:text-white transition-all"
            >
                Jouer
            </button>
        </div>
    </div>
</div>
