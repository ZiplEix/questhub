<script lang="ts">
    import Header from "$lib/components/Header.svelte";
    import { setContext } from "svelte";
    import ToastNotification from "../player/ToastNotification.svelte";

    let { children } = $props();

    // Simple Toast System
    let toast = $state<{ message: string; visible: boolean; type: string }>({
        message: "",
        visible: false,
        type: "info",
    });

    function showToast(message: string, type = "info") {
        toast = { message, visible: true, type };
        setTimeout(() => {
            toast.visible = false;
        }, 2000);
    }

    setContext("toast", { showToast });
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
        <!-- Toast Layer (Center Screen) -->
        <div
            class="absolute inset-0 pointer-events-none z-100 flex items-center justify-center"
        >
            {#if toast.visible}
                <ToastNotification message={toast.message} type={toast.type} />
            {/if}
        </div>

        {@render children()}
    </div>
</div>
