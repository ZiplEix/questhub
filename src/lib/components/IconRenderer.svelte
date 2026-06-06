<script lang="ts">
    import { getIconComponent } from "$lib/utils/iconMap";
    import type { ClassValue } from "clsx";

    interface Props {
        iconName?: string;
        size?: number;
        class?: ClassValue;
    }

    let { iconName, size = 20, class: className }: Props = $props();

    const Icon = getIconComponent(iconName);
</script>

<Icon {size} class={className} />
