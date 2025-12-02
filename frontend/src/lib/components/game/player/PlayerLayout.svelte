<script lang="ts">
    import Header from "$lib/components/Header.svelte";
    import { setContext } from "svelte";
    import ToastNotification from "./ToastNotification.svelte";

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

<div class="h-screen flex flex-col bg-stone-900 overflow-hidden">
    <!-- Global Header (Fixed at top, maybe overlay or solid?) -->
    <!-- User asked to KEEP the header. Usually headers are at the top. -->
    <div class="z-50 relative bg-cream shadow-sm">
        <Header />
    </div>

    <div class="flex-1 flex overflow-hidden relative">
        <!-- Toast Layer (Center Screen) -->
        <div
            class="absolute inset-0 pointer-events-none z-50 flex items-center justify-center"
        >
            {#if toast.visible}
                <ToastNotification message={toast.message} type={toast.type} />
            {/if}
        </div>

        {@render children()}
    </div>
</div>
