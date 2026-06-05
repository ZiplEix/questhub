<script lang="ts">
    import Chat from "../shared/Chat.svelte";

    let { entities = [], currentUserId = "" } = $props<{
        entities: any[];
        currentUserId?: string;
    }>();

    let players = $derived(
        entities
            .filter((e: any) => e.type === "player" && e.userId)
            .map((p: any) => ({ id: p.userId, name: p.name })),
    );
</script>

<Chat isGM={true} {players} {currentUserId} senderName="GM" />
