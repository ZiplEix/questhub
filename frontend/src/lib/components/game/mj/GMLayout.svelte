<script lang="ts">
    import Header from "$lib/components/Header.svelte";

    let { children } = $props();
</script>

<div
    class="h-screen w-screen flex flex-col bg-stone-900 overflow-hidden font-sans text-stone-800"
>
    <!-- Global Header -->
    <div class="flex-shrink-0 z-50 relative">
        <Header />
    </div>

    <!-- Main Content Area -->
    <div class="flex-1 flex overflow-hidden relative">
        {@render children()}
    </div>
</div>
