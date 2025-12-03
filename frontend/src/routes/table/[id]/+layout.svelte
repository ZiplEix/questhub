<script lang="ts">
    import { onMount, onDestroy } from "svelte";
    import { authClient } from "$lib/auth-client";
    import { connectWebSocket, closeWebSocket } from "$lib/websocket";

    let { children } = $props();

    onMount(async () => {
        try {
            const { data: tokenData } = await authClient.token();
            if (tokenData?.token) {
                connectWebSocket(tokenData.token);
            }
        } catch (e) {
            console.error("Failed to connect to WebSocket:", e);
        }
    });

    onDestroy(() => {
        closeWebSocket();
    });
</script>

{@render children()}
