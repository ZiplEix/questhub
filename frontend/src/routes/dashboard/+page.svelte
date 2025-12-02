<script lang="ts">
    import { Dice5, LogOut, Archive, ArchiveRestore, X } from "lucide-svelte";
    import GameCard from "$lib/components/GameCard.svelte";
    import Header from "$lib/components/Header.svelte";
    import { api } from "$lib/api";
    import { goto } from "$app/navigation";
    import { onMount } from "svelte";
    import { authClient } from "$lib/auth-client";

    let jwtToken = $state<string | null>(null);
    let games = $state<any[]>([]);
    let loading = $state(true);

    onMount(async () => {
        try {
            const { data, error: tokenError } = await authClient.token();

            if (tokenError || !data?.token) {
                console.error(tokenError);
                goto("/login");
                return;
            }

            jwtToken = data.token;
            await fetchGames(jwtToken);
        } catch (error) {
            console.error(error);
            goto("/login");
        }
    });

    async function fetchGames(token: string | null) {
        if (!token) return;
        try {
            const response = await api.get("/table", {
                headers: {
                    Authorization: `Bearer ${token}`,
                },
            });
            games = response.data;
        } catch (e) {
            console.error(e);
        } finally {
            loading = false;
        }
    }

    let showArchived = $state(false);
    let dialog: HTMLDialogElement;
    let newGameName = $state("");
    let creating = $state(false);
    let imageType = $state<"url" | "upload">("url");
    let imageUrl = $state("");
    let isUploading = $state(false);

    async function handleFileUpload(e: Event) {
        const input = e.target as HTMLInputElement;
        if (!input.files || input.files.length === 0) return;

        const file = input.files[0];
        const formData = new FormData();
        formData.append("image", file);

        isUploading = true;
        try {
            const response = await api.post("/upload/image", formData, {
                headers: {
                    Authorization: `Bearer ${jwtToken}`,
                    "Content-Type": "multipart/form-data",
                },
            });
            imageUrl = response.data.url;
        } catch (error) {
            console.error("Upload failed:", error);
            alert("Erreur lors de l'upload de l'image");
        } finally {
            isUploading = false;
        }
    }

    let filteredGames = $derived(
        games.filter((g) => (showArchived ? !g.is_active : g.is_active)),
    );

    function openDialog() {
        dialog.showModal();
    }

    function closeDialog() {
        dialog.close();
        newGameName = "";
        imageUrl = "";
        imageType = "url";
    }

    async function createGame(e: SubmitEvent) {
        e.preventDefault();
        creating = true;

        let response: any;

        try {
            response = await api.post(
                "/table",
                { name: newGameName, image_url: imageUrl },
                {
                    headers: {
                        Authorization: `Bearer ${jwtToken}`,
                    },
                },
            );
            console.log(response.data);
            await fetchGames(jwtToken); // Refresh list
            goto(`/table/${response.data.id}`);
        } catch (e: any) {
            console.error(e);
        } finally {
            creating = false;
        }

        closeDialog();
    }

    async function handleDeleteGame(id: string, name: string) {
        if (
            !confirm(
                `Êtes-vous sûr de vouloir supprimer la partie "${name}" ? Cette action est irréversible.`,
            )
        ) {
            return;
        }

        try {
            await api.delete(`/table/${id}`, {
                headers: {
                    Authorization: `Bearer ${jwtToken}`,
                },
            });
            await fetchGames(jwtToken);
        } catch (e) {
            console.error("Failed to delete game:", e);
            alert("Erreur lors de la suppression de la partie");
        }
    }
</script>

<div class="min-h-screen bg-cream">
    <Header />

    <main class="max-w-7xl mx-auto px-4 md:px-8 py-16">
        <div
            class="flex flex-col md:flex-row justify-between items-start md:items-center gap-6 mb-12"
        >
            <div>
                <h1 class="font-display font-bold text-4xl text-dark-gray mb-3">
                    Mes Parties
                </h1>
                <p class="text-dark-gray/60 text-lg">
                    Gérez vos aventures et rejoignez vos tables.
                </p>
            </div>

            <div class="flex gap-4">
                <button
                    class="flex items-center gap-2 px-5 py-3 rounded-2xl text-sm font-medium transition-all border {showArchived
                        ? 'bg-dark-gray text-white border-dark-gray'
                        : 'bg-white text-dark-gray border-stone-200 hover:bg-stone-50'}"
                    onclick={() => (showArchived = !showArchived)}
                >
                    {#if showArchived}
                        <ArchiveRestore size={18} />
                        Voir les parties actives
                    {:else}
                        <Archive size={18} />
                        Voir les archives
                    {/if}
                </button>

                <button
                    class="bg-burnt-orange text-white px-6 py-3 rounded-2xl font-bold shadow-md hover:bg-opacity-90 transition-all hover:-translate-y-0.5 flex items-center gap-2 hover:cursor-pointer"
                    onclick={openDialog}
                >
                    <Dice5 size={20} />
                    Nouvelle Partie
                </button>
            </div>
        </div>

        {#if loading}
            <div class="text-center py-24">
                <p class="text-dark-gray/60 text-lg">
                    Chargement des parties...
                </p>
            </div>
        {:else if filteredGames.length > 0}
            <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
                {#each filteredGames as game (game.id)}
                    <GameCard
                        id={game.id}
                        name={game.name}
                        gm={game.gm_name}
                        imageUrl={game.image_url ||
                            `https://placehold.co/600x400/3D405B/F9F7F2?text=${game.name}`}
                        createdAt={game.created_at}
                        isActive={game.is_active}
                        isGm={jwtToken ? true : false}
                        onDelete={() => handleDeleteGame(game.id, game.name)}
                    />
                {/each}
            </div>
        {:else}
            <div
                class="text-center py-24 bg-white rounded-3xl border border-stone-100 shadow-sm"
            >
                <div
                    class="bg-stone-50 w-20 h-20 rounded-full flex items-center justify-center mx-auto mb-6 text-dark-gray/40"
                >
                    <Dice5 size={40} />
                </div>
                <h3 class="font-display font-bold text-2xl text-dark-gray mb-3">
                    Aucune partie trouvée
                </h3>
                <p class="text-dark-gray/60 max-w-md mx-auto mb-8 text-lg">
                    {showArchived
                        ? "Vous n'avez aucune partie archivée."
                        : "Vous n'avez pas encore de partie active. Créez-en une pour commencer l'aventure !"}
                </p>
                {#if !showArchived}
                    <button
                        class="bg-burnt-orange text-white px-8 py-4 rounded-2xl font-bold shadow-md hover:bg-opacity-90 transition-all hover:-translate-y-0.5 text-lg hover:cursor-pointer"
                        onclick={openDialog}
                    >
                        Créer une table
                    </button>
                {/if}
            </div>
        {/if}
    </main>

    <dialog
        bind:this={dialog}
        class="m-auto rounded-3xl bg-white shadow-xl backdrop:bg-black/50 p-0 w-full max-w-lg open:animate-in open:fade-in open:zoom-in-95 open:duration-200 backdrop:animate-in backdrop:fade-in backdrop:duration-200"
    >
        <div class="p-8 md:p-10">
            <div class="flex justify-between items-center mb-6">
                <h2 class="font-display font-bold text-2xl text-dark-gray">
                    Nouvelle Partie
                </h2>
                <button
                    onclick={closeDialog}
                    class="text-dark-gray/40 hover:text-dark-gray transition-colors hover:cursor-pointer"
                >
                    <X size={24} />
                </button>
            </div>

            <form onsubmit={createGame} class="space-y-6">
                <div>
                    <label
                        class="block text-sm font-medium text-dark-gray mb-2 ml-1"
                        for="gameName"
                    >
                        Nom de la partie
                    </label>
                    <input
                        type="text"
                        id="gameName"
                        bind:value={newGameName}
                        required
                        class="w-full px-4 py-3 rounded-xl border border-stone-200 focus:outline-none focus:ring-2 focus:ring-burnt-orange/50 focus:border-burnt-orange transition-all bg-stone-50"
                        placeholder="Ex: La Légende de Vox Machina"
                    />
                </div>

                <div>
                    <label
                        class="block text-sm font-medium text-dark-gray mb-2 ml-1"
                        for="imageType"
                    >
                        Image de couverture (optionnel)
                    </label>

                    <div class="flex gap-4 mb-4">
                        <button
                            type="button"
                            class="flex-1 py-2 rounded-lg text-sm font-medium transition-colors {imageType ===
                            'url'
                                ? 'bg-dark-gray text-white'
                                : 'bg-stone-100 text-dark-gray hover:bg-stone-200'}"
                            onclick={() => (imageType = "url")}
                        >
                            Lien URL
                        </button>
                        <button
                            type="button"
                            class="flex-1 py-2 rounded-lg text-sm font-medium transition-colors {imageType ===
                            'upload'
                                ? 'bg-dark-gray text-white'
                                : 'bg-stone-100 text-dark-gray hover:bg-stone-200'}"
                            onclick={() => (imageType = "upload")}
                        >
                            Upload
                        </button>
                    </div>

                    {#if imageType === "url"}
                        <input
                            type="url"
                            bind:value={imageUrl}
                            class="w-full px-4 py-3 rounded-xl border border-stone-200 focus:outline-none focus:ring-2 focus:ring-burnt-orange/50 focus:border-burnt-orange transition-all bg-stone-50"
                            placeholder="https://example.com/image.jpg"
                        />
                    {:else}
                        <div class="relative">
                            <input
                                type="file"
                                accept="image/*"
                                onchange={handleFileUpload}
                                class="hidden"
                                id="fileInput"
                            />
                            <label
                                for="fileInput"
                                class="flex flex-col items-center justify-center w-full h-32 border-2 border-dashed border-stone-300 rounded-xl cursor-pointer hover:bg-stone-50 transition-colors"
                            >
                                {#if isUploading}
                                    <div
                                        class="animate-spin rounded-full h-8 w-8 border-b-2 border-burnt-orange"
                                    ></div>
                                {:else if imageUrl && imageType === "upload"}
                                    <img
                                        src={imageUrl}
                                        alt="Preview"
                                        class="h-full w-full object-cover rounded-xl opacity-50"
                                    />
                                    <div
                                        class="absolute inset-0 flex items-center justify-center"
                                    >
                                        <span
                                            class="bg-white/80 px-3 py-1 rounded-full text-sm font-medium text-dark-gray"
                                            >Changer l'image</span
                                        >
                                    </div>
                                {:else}
                                    <div class="text-center p-4">
                                        <p
                                            class="text-sm font-medium text-dark-gray"
                                        >
                                            Cliquez pour uploader
                                        </p>
                                        <p
                                            class="text-xs text-dark-gray/60 mt-1"
                                        >
                                            PNG, JPG jusqu'à 5MB
                                        </p>
                                    </div>
                                {/if}
                            </label>
                        </div>
                        <p
                            class="text-xs text-amber-600 mt-2 flex items-center gap-1"
                        >
                            ⚠️ Attention : toute image uploadée est publique.
                        </p>
                    {/if}
                </div>

                <div class="flex justify-end gap-3 mt-8">
                    <button
                        type="button"
                        onclick={closeDialog}
                        class="px-6 py-3 rounded-xl font-medium text-dark-gray hover:bg-stone-50 transition-colors hover:cursor-pointer"
                    >
                        Annuler
                    </button>
                    <button
                        type="submit"
                        disabled={creating}
                        class="bg-burnt-orange text-white px-6 py-3 rounded-xl font-bold shadow-md hover:bg-opacity-90 transition-all hover:-translate-y-0.5 disabled:opacity-60 disabled:hover:translate-y-0 hover:cursor-pointer"
                    >
                        {creating ? "Création..." : "Créer la partie"}
                    </button>
                </div>
            </form>
        </div>
    </dialog>
</div>
