<script lang="ts">
    import { page } from "$app/state";
    import { goto } from "$app/navigation";
    import { acceptInvitation as acceptInvitationApi, declineInvitation as declineInvitationApi, kickPlayer } from "$lib/api";
    import { UserPlus, Trash2 } from "lucide-svelte";

    let { players, invitations, onRefresh } = $props<{
        players: any[];
        invitations: any[];
        onRefresh: () => void;
    }>();

    async function acceptInvitation(userId: string) {
        const gameId = page.params.id;
        if (!gameId) return;
        try {
            await acceptInvitationApi(gameId, userId);
            onRefresh();
        } catch (error) {
            console.error("Failed to accept invitation:", error);
        }
    }

    async function declineInvitation(userId: string) {
        const gameId = page.params.id;
        if (!gameId) return;
        try {
            await declineInvitationApi(gameId, userId);
            onRefresh();
        } catch (error) {
            console.error("Failed to decline invitation:", error);
        }
    }

    async function removePlayer(userId: string, userName: string) {
        if (
            !confirm(
                `Êtes-vous sûr de vouloir exclure ${userName} de la partie ?`,
            )
        ) {
            return;
        }

        const gameId = page.params.id;
        if (!gameId) return;
        try {
            await kickPlayer(gameId, userId);
            onRefresh();
        } catch (error) {
            console.error("Failed to remove player:", error);
        }
    }
</script>

<div class="space-y-6 animate-in fade-in duration-300">
    <!-- Invitations Section -->
    {#if invitations.length > 0}
        <div class="mb-6">
            <h3
                class="text-lg font-bold text-dark-gray mb-3 flex items-center gap-2"
            >
                <span class="w-2 h-2 rounded-full bg-burnt-orange"></span>
                Demandes en attente ({invitations.length})
            </h3>
            <div class="space-y-3">
                {#each invitations as invitation}
                    <div
                        class="flex items-center justify-between p-4 rounded-xl border border-burnt-orange/20 bg-burnt-orange/5"
                    >
                        <div class="flex items-center gap-4">
                            <div
                                class="w-10 h-10 rounded-full bg-white flex items-center justify-center text-burnt-orange font-bold shadow-sm"
                            >
                                {invitation.user_name.charAt(0)}
                            </div>
                            <div>
                                <p class="font-bold text-dark-gray">
                                    {invitation.user_name}
                                </p>
                                <p class="text-xs text-stone-500">
                                    Demandé le {new Date(
                                        invitation.created_at,
                                    ).toLocaleDateString()}
                                </p>
                            </div>
                        </div>
                        <div class="flex items-center gap-2">
                            <button
                                onclick={() =>
                                    acceptInvitation(invitation.user_id)}
                                class="px-3 py-1.5 bg-burnt-orange text-white text-sm font-medium rounded-lg hover:bg-opacity-90 transition-all"
                            >
                                Accepter
                            </button>
                            <button
                                onclick={() =>
                                    declineInvitation(invitation.user_id)}
                                class="px-3 py-1.5 bg-white border border-stone-200 text-stone-600 text-sm font-medium rounded-lg hover:bg-stone-50 transition-all"
                            >
                                Refuser
                            </button>
                        </div>
                    </div>
                {/each}
            </div>
        </div>
    {/if}

    <div class="flex justify-between items-center mb-4">
        <h3 class="text-lg font-bold text-dark-gray">
            Joueurs ({players.length})
        </h3>
        <button
            class="flex items-center gap-2 px-4 py-2 bg-stone-100 text-dark-gray rounded-xl font-medium hover:bg-stone-200 transition-all text-sm"
        >
            <UserPlus size={18} />
            Inviter
        </button>
    </div>

    <div class="space-y-3">
        {#each players as player}
            {#if player.is_gm}
                <!-- GM Card (Clickable) -->
                <div
                    onclick={() => goto(`/table/${page.params.id}/gm/settings/player/${player.user_id}`)}
                    onkeydown={(e) => e.key === 'Enter' && goto(`/table/${page.params.id}/gm/settings/player/${player.user_id}`)}
                    role="button"
                    tabindex="0"
                    class="flex items-center justify-between p-4 rounded-xl border border-stone-100 hover:border-stone-300 hover:bg-stone-50/50 transition-all cursor-pointer group"
                >
                    <div class="flex items-center gap-4">
                        <div class="relative">
                            {#if player.avatar_url}
                                <img
                                    src={player.avatar_url}
                                    alt={player.name}
                                    class="w-10 h-10 rounded-full object-cover"
                                />
                            {:else}
                                <div
                                    class="w-10 h-10 rounded-full bg-burnt-orange/20 flex items-center justify-center text-burnt-orange font-bold"
                                >
                                    {player.name.charAt(0)}
                                </div>
                            {/if}
                            <span 
                                class="absolute -bottom-0.5 -right-0.5 w-3.5 h-3.5 rounded-full border-2 border-white shadow-sm"
                                style="background-color: {player.ping_color || '#E07A5F'};"
                                title="Couleur du ping"
                            ></span>
                        </div>
                        <div>
                            <div class="flex items-center gap-2">
                                <p class="font-bold text-dark-gray group-hover:text-burnt-orange transition-colors">
                                    {player.name}
                                </p>
                                <span
                                    class="px-2 py-0.5 rounded-full bg-burnt-orange/10 text-burnt-orange text-xs font-bold border border-burnt-orange/20"
                                >
                                    MJ
                                </span>
                            </div>
                            <p class="text-xs text-stone-500">
                                Rejoint le {new Date(
                                    player.joined_at,
                                ).toLocaleDateString()}
                            </p>
                        </div>
                    </div>
                </div>
            {:else}
                <!-- Clickable Player Card -->
                <div
                    onclick={() => goto(`/table/${page.params.id}/gm/settings/player/${player.user_id}`)}
                    onkeydown={(e) => e.key === 'Enter' && goto(`/table/${page.params.id}/gm/settings/player/${player.user_id}`)}
                    role="button"
                    tabindex="0"
                    class="flex items-center justify-between p-4 rounded-xl border border-stone-100 hover:border-stone-300 hover:bg-stone-50/50 transition-all cursor-pointer group"
                >
                    <div class="flex items-center gap-4">
                        <div class="relative">
                            {#if player.avatar_url}
                                <img
                                    src={player.avatar_url}
                                    alt={player.name}
                                    class="w-10 h-10 rounded-full object-cover"
                                />
                            {:else}
                                <div
                                    class="w-10 h-10 rounded-full bg-burnt-orange/20 flex items-center justify-center text-burnt-orange font-bold"
                                >
                                    {player.name.charAt(0)}
                                </div>
                            {/if}
                            <span 
                                class="absolute -bottom-0.5 -right-0.5 w-3.5 h-3.5 rounded-full border-2 border-white shadow-sm"
                                style="background-color: {player.ping_color || '#E07A5F'};"
                                title="Couleur du ping"
                            ></span>
                        </div>
                        <div>
                            <div class="flex items-center gap-2">
                                <p class="font-bold text-dark-gray group-hover:text-burnt-orange transition-colors">
                                    {player.name}
                                </p>
                            </div>
                            {#if player.character_name}
                                <p class="text-sm text-burnt-orange font-medium">
                                    Incarne {player.character_name}
                                </p>
                            {/if}
                            <p class="text-xs text-stone-500">
                                Rejoint le {new Date(
                                    player.joined_at,
                                ).toLocaleDateString()}
                            </p>
                        </div>
                    </div>
                    <div class="flex items-center gap-2">
                        <button
                            onclick={(e) => {
                                e.stopPropagation();
                                removePlayer(player.user_id, player.name);
                            }}
                            class="p-2 text-stone-400 hover:text-red-500 hover:bg-red-50 rounded-lg transition-all"
                            title="Exclure"
                        >
                            <Trash2 size={18} />
                        </button>
                    </div>
                </div>
            {/if}
        {/each}
    </div>
</div>
