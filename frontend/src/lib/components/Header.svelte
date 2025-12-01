<script lang="ts">
    import { Dice5, LogOut } from "lucide-svelte";
    import { authClient } from "$lib/auth-client";

    const session = authClient.useSession();

    const handleLogout = async () => {
        await authClient.signOut({
            fetchOptions: {
                onSuccess: () => {
                    window.location.href = "/";
                },
            },
        });
    };
</script>

<header
    class="w-full py-6 px-4 md:px-8 flex justify-between items-center max-w-7xl mx-auto bg-white border-b border-stone-100 sticky top-0 z-50"
>
    <a href="/" class="flex items-center gap-2 text-dark-gray">
        <div class="bg-burnt-orange text-white p-2 rounded-xl shadow-sm">
            <Dice5 size={24} />
        </div>
        <span class="font-display font-bold text-xl tracking-tight"
            >QuestHub</span
        >
    </a>

    <div class="flex items-center gap-4">
        {#if $session.data}
            <div class="flex items-center gap-3 pl-4 border-l border-stone-100">
                <div class="text-right hidden md:block">
                    <p class="text-sm font-bold text-dark-gray">
                        {$session.data.user.name}
                    </p>
                    <p class="text-xs text-dark-gray/60">Joueur</p>
                </div>
                <div
                    class="w-10 h-10 rounded-full bg-mustard-yellow/30 flex items-center justify-center text-mustard-yellow font-bold border-2 border-white shadow-sm"
                >
                    {$session.data.user.name.charAt(0).toUpperCase()}
                </div>
                <button
                    class="text-dark-gray/40 hover:text-red-500 transition-colors ml-2 hover:cursor-pointer"
                    title="Se déconnecter"
                    onclick={handleLogout}
                >
                    <LogOut size={20} />
                </button>
            </div>
        {:else}
            <a
                href="/login"
                class="font-medium text-dark-gray hover:text-burnt-orange transition-colors"
            >
                Connexion
            </a>
            <a
                href="/signup"
                class="bg-burnt-orange text-white px-5 py-2.5 rounded-xl font-medium shadow-md hover:bg-opacity-90 transition-all hover:-translate-y-0.5"
            >
                Créer une table
            </a>
        {/if}
    </div>
</header>
