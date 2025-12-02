<script lang="ts">
    import { page } from "$app/state";
    import { api } from "$lib/api";
    import { authClient } from "$lib/auth-client";
    import { Mail, Check, X } from "lucide-svelte";

    let { invitations, onRefresh } = $props<{
        invitations: any[];
        onRefresh: () => void;
    }>();

    async function acceptInvitation(userId: string) {
        const gameId = page.params.id;
        try {
            const { data: tokenData } = await authClient.token();
            if (tokenData?.token) {
                await api.post(
                    `/table/${gameId}/invitations/${userId}/accept`,
                    {},
                    {
                        headers: {
                            Authorization: `Bearer ${tokenData.token}`,
                        },
                    },
                );
                onRefresh();
            }
        } catch (error) {
            console.error("Failed to accept invitation:", error);
        }
    }

    async function declineInvitation(userId: string) {
        const gameId = page.params.id;
        try {
            const { data: tokenData } = await authClient.token();
            if (tokenData?.token) {
                await api.post(
                    `/table/${gameId}/invitations/${userId}/decline`,
                    {},
                    {
                        headers: {
                            Authorization: `Bearer ${tokenData.token}`,
                        },
                    },
                );
                onRefresh();
            }
        } catch (error) {
            console.error("Failed to decline invitation:", error);
        }
    }
</script>

<div class="space-y-6 animate-in fade-in duration-300">
    <div class="flex justify-between items-center mb-4">
        <h3 class="text-lg font-bold text-dark-gray">
            Demandes d'adhésion ({invitations.length})
        </h3>
    </div>

    {#if invitations.length === 0}
        <div
            class="flex flex-col items-center justify-center py-12 text-center"
        >
            <div
                class="w-16 h-16 bg-stone-100 rounded-full flex items-center justify-center text-stone-400 mb-4"
            >
                <Mail size={32} />
            </div>
            <h4 class="text-lg font-bold text-dark-gray">
                Aucune demande en attente
            </h4>
            <p class="text-stone-500 max-w-sm mt-2">
                Partagez le lien d'invitation pour permettre aux joueurs de
                rejoindre votre table.
            </p>
        </div>
    {:else}
        <div class="space-y-3">
            {#each invitations as invite}
                <div
                    class="flex flex-col md:flex-row md:items-center justify-between p-4 rounded-xl border border-stone-100 hover:border-stone-200 transition-all gap-4"
                >
                    <div class="flex items-center gap-4">
                        <div
                            class="w-10 h-10 rounded-full bg-burnt-orange/20 flex items-center justify-center text-burnt-orange font-bold"
                        >
                            {invite.user_name.charAt(0)}
                        </div>
                        <div>
                            <p class="font-bold text-dark-gray">
                                {invite.user_name}
                            </p>
                            <p class="text-xs text-stone-500">
                                Demande envoyée le {new Date(
                                    invite.created_at,
                                ).toLocaleDateString()}
                            </p>
                        </div>
                    </div>
                    <div class="flex items-center gap-2">
                        <button
                            onclick={() => acceptInvitation(invite.user_id)}
                            class="flex items-center gap-2 px-3 py-2 bg-green-100 text-green-700 rounded-lg font-medium hover:bg-green-200 transition-colors text-sm"
                        >
                            <Check size={16} />
                            Accepter
                        </button>
                        <button
                            onclick={() => declineInvitation(invite.user_id)}
                            class="flex items-center gap-2 px-3 py-2 bg-red-100 text-red-700 rounded-lg font-medium hover:bg-red-200 transition-colors text-sm"
                        >
                            <X size={16} />
                            Refuser
                        </button>
                    </div>
                </div>
            {/each}
        </div>
    {/if}
</div>
