<script lang="ts">
    import type { ComponentType } from "svelte";
    import type { Icon as LucideIcon } from "lucide-svelte";

    let {
        icon: Icon,
        title,
        description,
    } = $props<{
        icon: ComponentType<LucideIcon>;
        title: string;
        description: string;
    }>();
</script>

<div
    class="bg-white p-8 rounded-2xl shadow-sm hover:shadow-md transition-shadow duration-300 flex flex-col items-start gap-4 border border-stone-100"
>
    <div class="bg-mustard-yellow/20 p-3 rounded-xl text-dark-gray">
        <Icon size={28} strokeWidth={1.5} />
    </div>
    <h3 class="font-display font-bold text-xl text-dark-gray">{title}</h3>
    <p class="text-dark-gray/80 leading-relaxed">{description}</p>
</div>
