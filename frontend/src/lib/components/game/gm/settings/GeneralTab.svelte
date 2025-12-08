<script lang="ts">
    import { page } from "$app/state";
    import { api } from "$lib/api";
    import { authClient } from "$lib/auth-client";
    import { Save, Copy, RefreshCw, Check, Play, Pause } from "lucide-svelte";

    let { game = $bindable() } = $props();

    let copied = $state(false);

    async function copyInviteLink() {
        if (!game?.invite_code) return;
        const link = `${window.location.origin}/invitation/${game.invite_code}`;
        await navigator.clipboard.writeText(link);
        copied = true;
        setTimeout(() => (copied = false), 2000);
    }

    async function regenerateInviteCode() {
        if (
            !confirm(
                "Êtes-vous sûr de vouloir régénérer le code d'invitation ? L'ancien lien ne fonctionnera plus.",
            )
        ) {
            return;
        }

        const gameId = page.params.id;
        try {
            const { data: tokenData } = await authClient.token();
            if (tokenData?.token) {
                const response = await api.post(
                    `/table/${gameId}/invite-code`,
                    {},
                    {
                        headers: {
                            Authorization: `Bearer ${tokenData.token}`,
                        },
                    },
                );
                game.invite_code = response.data.invite_code;
            }
        } catch (error) {
            console.error("Failed to regenerate invite code:", error);
        }
    }

    async function updateGameState(newState: string) {
        if (game.state === newState) return;

        const gameId = page.params.id;
        // Optimistic update
        const oldState = game.state;
        game.state = newState;

        try {
            const { data: tokenData } = await authClient.token();
            if (tokenData?.token) {
                await api.put(
                    `/table/${gameId}/state`,
                    { state: newState },
                    {
                        headers: {
                            Authorization: `Bearer ${tokenData.token}`,
                        },
                    },
                );
            }
        } catch (error) {
            console.error("Failed to update game state:", error);
            // Revert on error
            game.state = oldState;
            alert("Erreur lors de la mise à jour du statut");
        }
    }
</script>

<div class="space-y-6 animate-in fade-in duration-300">
    <div class="space-y-2">
        <label for="tableName" class="text-sm font-bold text-dark-gray"
            >Nom de la table</label
        >
        <input
            type="text"
            id="tableName"
            bind:value={game.name}
            class="w-full px-4 py-2 rounded-xl border border-stone-200 focus:outline-none focus:ring-2 focus:ring-burnt-orange/20 focus:border-burnt-orange transition-all"
        />
    </div>

    <div class="space-y-2">
        <label class="text-sm font-bold text-dark-gray block" for="inviteLink">
            Lien d'invitation
        </label>
        <div class="flex gap-2" id="inviteLink">
            <div
                class="flex-1 bg-stone-50 px-4 py-2 rounded-xl border border-stone-200 text-stone-600 font-mono text-sm truncate"
            >
                {window.location.origin}/invitation/{game.invite_code}
            </div>
            <button
                onclick={copyInviteLink}
                class="flex items-center gap-2 px-4 py-2 bg-white border border-stone-200 text-dark-gray rounded-xl font-medium hover:bg-stone-50 hover:border-burnt-orange/30 hover:text-burnt-orange transition-all shadow-sm min-w-[100px] justify-center"
            >
                {#if copied}
                    <Check size={18} />
                    <span>Copié</span>
                {:else}
                    <Copy size={18} />
                    <span>Copier</span>
                {/if}
            </button>
            <button
                onclick={regenerateInviteCode}
                class="flex items-center gap-2 px-4 py-2 bg-white border border-stone-200 text-dark-gray rounded-xl font-medium hover:bg-stone-50 hover:border-burnt-orange/30 hover:text-burnt-orange transition-all shadow-sm"
                title="Régénérer le lien"
            >
                <RefreshCw size={18} />
            </button>
        </div>
        <p class="text-xs text-stone-500">
            Partagez ce lien avec vos joueurs pour qu'ils puissent rejoindre la
            partie.
        </p>
    </div>

    <div class="space-y-2">
        <label class="text-sm font-bold text-dark-gray block" for="gameState"
            >Status de la partie</label
        >
        <div class="flex gap-4">
            <button
                onclick={() => updateGameState("ongoing")}
                class="flex-1 flex items-center justify-center gap-2 px-4 py-3 rounded-xl border transition-all {game.state ===
                'ongoing'
                    ? 'bg-emerald-50 border-emerald-500 text-emerald-700'
                    : 'bg-white border-stone-200 text-stone-500 hover:border-emerald-200 hover:text-emerald-600'}"
            >
                <Play size={20} />
                <span class="font-medium">En cours</span>
            </button>
            <button
                onclick={() => updateGameState("paused")}
                class="flex-1 flex items-center justify-center gap-2 px-4 py-3 rounded-xl border transition-all {game.state ===
                'paused'
                    ? 'bg-amber-50 border-amber-500 text-amber-700'
                    : 'bg-white border-stone-200 text-stone-500 hover:border-amber-200 hover:text-amber-600'}"
            >
                <Pause size={20} />
                <span class="font-medium">En pause</span>
            </button>
        </div>
        <p class="text-xs text-stone-500">
            Mettre la partie en pause empêche les joueurs d'envoyer des messages
            ou d'interagir avec la table.
        </p>
    </div>

    <div class="pt-4 border-t border-stone-100 flex justify-end">
        <button
            class="flex items-center gap-2 px-6 py-2.5 bg-burnt-orange text-white rounded-xl font-medium shadow-md hover:bg-opacity-90 transition-all hover:-translate-y-0.5"
        >
            <Save size={18} />
            Enregistrer
        </button>
    </div>
</div>
