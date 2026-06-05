<script lang="ts">
    import { page } from "$app/state";
    import { connectWebSocket, closeWebSocket } from "$lib/websocket";

    let { children } = $props();

    $effect(() => {
        const gameId = page.params.id;
        if (gameId) {
            connectWebSocket(gameId);
        }
        return () => {
            closeWebSocket();
        };
    });
</script>

{@render children()}
