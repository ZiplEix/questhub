<script lang="ts">
    import { page } from "$app/stores";
    import { onMount } from "svelte";
    import { api } from "$lib/api";
    import { authClient } from "$lib/auth-client";
    import { goto } from "$app/navigation";
    import Header from "$lib/components/Header.svelte";
    import { Crown } from "lucide-svelte";

    let game = $state<any>(null);
    let loading = $state(true);
    let error = $state("");
    let userId = $state<string | null>(null);

    const session = authClient.useSession();

    $effect(() => {
        if ($session.data?.user) {
            userId = $session.data.user.id;
        }
    });

    onMount(async () => {
        const { data, error: tokenError } = await authClient.token();

        if (tokenError || !data?.token) {
            goto("/login");
            return;
        }

        fetchGame(data.token);
    });

    async function fetchGame(token: string) {
        try {
            const response = await api.get(`/table/${$page.params.id}`, {
                headers: {
                    Authorization: `Bearer ${token}`,
                },
            });
            game = response.data;
        } catch (e) {
            console.error(e);
            error = "Impossible de charger la partie.";
        } finally {
            loading = false;
        }
    }
</script>

<div class="min-h-screen bg-cream">
    <Header />

    <main class="max-w-4xl mx-auto px-4 py-16">
        {#if loading}
            <div class="text-center text-dark-gray">Chargement...</div>
        {:else if error}
            <div class="text-center text-red-500">{error}</div>
        {:else if game}
            <div
                class="bg-white rounded-3xl shadow-sm border border-stone-100 p-8"
            >
                <div class="flex items-center gap-4 mb-6">
                    <h1 class="font-display font-bold text-3xl text-dark-gray">
                        {game.name}
                    </h1>
                    {#if userId === game.gm_id}
                        <div
                            class="bg-burnt-orange/10 text-burnt-orange px-3 py-1 rounded-full flex items-center gap-1.5 text-sm font-bold border border-burnt-orange/20"
                        >
                            <Crown size={14} />
                            MJ
                        </div>
                    {/if}
                </div>

                <div class="space-y-4">
                    <div>
                        <span class="font-bold text-dark-gray">ID:</span>
                        <span class="text-dark-gray/80 ml-2">{game.id}</span>
                    </div>
                    <div>
                        <span class="font-bold text-dark-gray"
                            >Code d'invitation:</span
                        >
                        <span
                            class="bg-stone-100 px-3 py-1 rounded-lg font-mono text-burnt-orange ml-2"
                            >{game.invite_code}</span
                        >
                    </div>
                    <div>
                        <span class="font-bold text-dark-gray">GM ID:</span>
                        <span class="text-dark-gray/80 ml-2">{game.gm_id}</span>
                    </div>
                    <div>
                        <span class="font-bold text-dark-gray">Statut:</span>
                        <span
                            class="ml-2 {game.is_active
                                ? 'text-green-600'
                                : 'text-red-500'}"
                        >
                            {game.is_active ? "Active" : "Archivée"}
                        </span>
                    </div>
                </div>
            </div>
        {/if}
    </main>
</div>
