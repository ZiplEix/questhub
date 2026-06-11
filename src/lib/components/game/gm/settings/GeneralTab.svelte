<script lang="ts">
    import { page } from "$app/state";
    import { updateGame, regenerateInviteCode as regenerateInviteCodeApi } from "$lib/api";
    import { Save, Copy, RefreshCw, Check, Play, Pause, Archive, ArchiveRestore } from "lucide-svelte";
    import { goto } from "$app/navigation";

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
        if (!gameId) return;
        try {
            const newCode = await regenerateInviteCodeApi(gameId);
            game.invite_code = newCode;
        } catch (error) {
            console.error("Failed to regenerate invite code:", error);
        }
    }

    async function updateGameState(newState: string) {
        if (game.state === newState) return;

        const gameId = page.params.id;
        if (!gameId) return;
        // Optimistic update
        const oldState = game.state;
        game.state = newState;

        try {
            await updateGame(gameId, { state: newState });
        } catch (error) {
            console.error("Failed to update game state:", error);
            // Revert on error
            game.state = oldState;
            alert("Erreur lors de la mise à jour du statut");
        }
    }

    async function saveGameName() {
        const gameId = page.params.id;
        if (!gameId) return;
        try {
            await updateGame(gameId, { name: game.name });
            alert("Nom de la table mis à jour");
        } catch (error) {
            console.error("Failed to update game name:", error);
            alert("Erreur lors de la mise à jour du nom");
        }
    }

    async function toggleArchiveGame() {
        const gameId = page.params.id;
        if (!gameId) return;

        const newActiveState = !game.is_active;
        const actionText = newActiveState ? "désarchiver" : "archiver";

        if (
            !confirm(
                `Êtes-vous sûr de vouloir ${actionText} cette partie ? ${
                    newActiveState
                        ? "Elle redeviendra active et modifiable par les joueurs et vous."
                        : "Elle sera verrouillée en lecture seule pour tous."
                }`,
            )
        ) {
            return;
        }

        try {
            await updateGame(gameId, { is_active: newActiveState });
            game.is_active = newActiveState;
            alert(
                `La partie a bien été ${
                    newActiveState ? "désarchivée" : "archivée"
                }.`,
            );
            if (!newActiveState) {
                // Redirect to dashboard after archiving
                goto("/dashboard");
            }
        } catch (error) {
            console.error(`Failed to ${actionText} game:`, error);
            alert(`Erreur lors de la mise à jour du statut d'archivage.`);
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

    <div class="space-y-2 pt-4 border-t border-stone-100">
        <label class="text-sm font-bold text-dark-gray block" for="archiveGame"
            >Archivage de la partie</label
        >
        <div class="flex flex-col md:flex-row gap-4 items-start md:items-center justify-between">
            <p class="text-xs text-stone-500 max-w-xl">
                Archiver la partie la rendra complètement en lecture seule pour vous et tous vos joueurs. Vous pourrez toujours consulter les données (chat, notes, plateaux) depuis le tableau de bord, mais aucune nouvelle action ne pourra être effectuée.
            </p>
            <button
                id="archiveGame"
                onclick={toggleArchiveGame}
                class="flex items-center gap-2 px-5 py-3 rounded-xl border font-semibold transition-all hover:cursor-pointer shadow-sm shrink-0
                {game.is_active
                    ? 'bg-red-50 border-red-200 text-red-700 hover:bg-red-100 hover:border-red-300'
                    : 'bg-emerald-50 border-emerald-200 text-emerald-700 hover:bg-emerald-100 hover:border-emerald-300'}"
            >
                {#if game.is_active}
                    <Archive size={18} />
                    <span>Archiver la partie</span>
                {:else}
                    <ArchiveRestore size={18} />
                    <span>Désarchiver la partie</span>
                {/if}
            </button>
        </div>
    </div>

    <div class="pt-4 border-t border-stone-100 flex justify-end">
        <button
            onclick={saveGameName}
            class="flex items-center gap-2 px-6 py-2.5 bg-burnt-orange text-white rounded-xl font-medium shadow-md hover:bg-opacity-90 transition-all hover:-translate-y-0.5"
        >
            <Save size={18} />
            Enregistrer
        </button>
    </div>
</div>
