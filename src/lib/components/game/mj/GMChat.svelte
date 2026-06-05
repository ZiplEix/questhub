<script lang="ts">
    import Chat from "../shared/Chat.svelte";
    import { fetchPlayers } from "$lib/api";
    import { page } from "$app/state";
    import { onMount } from "svelte";

    let { currentUserId = "" } = $props<{
        currentUserId?: string;
    }>();

    let players = $state<{ id: string; name: string }[]>([]);

    onMount(async () => {
        try {
            const gameId = page.params.id;
            if (gameId) {
                const data = await fetchPlayers(gameId);
                players = data.map((p: any) => ({ id: p.user_id, name: p.name }));
            }
        } catch (e) {
            console.error("Failed to fetch players for chat:", e);
        }
    });
</script>

<Chat isGM={true} {players} {currentUserId} senderName="GM" />
