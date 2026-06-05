<script lang="ts">
    import {
        Dice5,
        LogOut,
        User,
        LayoutDashboard,
        ChevronDown,
        Settings,
    } from "lucide-svelte";
    import { authClient } from "$lib/auth-client";
    import { clickOutside } from "$lib/actions/clickOutside";
    import { page } from "$app/state";

    const session = authClient.useSession();
    let isMenuOpen = $state(false);

    let path = $state(page.url.pathname);

    const handleLogout = async () => {
        await authClient.signOut({
            fetchOptions: {
                onSuccess: () => {
                    window.location.href = "/";
                },
            },
        });
    };

    function toggleMenu() {
        isMenuOpen = !isMenuOpen;
    }

    function closeMenu() {
        isMenuOpen = false;
    }
</script>

<header
    class="w-full py-2 px-4 md:px-8 flex justify-between items-center mx-auto bg-white border-b border-stone-100 sticky top-0 z-50"
>
    <a
        href={$session.data ? `/dashboard` : `/`}
        class="flex items-center gap-2 text-dark-gray"
    >
        <div class="bg-burnt-orange text-white p-2 rounded-xl shadow-sm">
            <Dice5 size={24} />
        </div>
        <span class="font-display font-bold text-xl tracking-tight"
            >QuestHub</span
        >
    </a>

    <!-- Centered on the header -->
    <!-- Centered on the header -->
    {#if path.endsWith("/gm")}
        <div
            class="absolute left-1/2 top-1/2 -translate-x-1/2 -translate-y-1/2"
        >
            <a
                href="{path}/settings"
                class="flex items-center gap-2 px-4 py-2 rounded-xl bg-white border border-stone-200 text-dark-gray hover:bg-stone-50 hover:border-burnt-orange/30 hover:text-burnt-orange transition-all shadow-sm"
                title="Paramètres de la partie"
            >
                <Settings size={18} />
                <span class="font-medium text-sm hidden md:inline"
                    >Paramètres de la partie</span
                >
            </a>
        </div>
    {:else if path.endsWith("/gm/settings")}
        <div
            class="absolute left-1/2 top-1/2 -translate-x-1/2 -translate-y-1/2"
        >
            <a
                href={path.replace("/settings", "")}
                class="flex items-center gap-2 px-4 py-2 rounded-xl bg-white border border-stone-200 text-dark-gray hover:bg-stone-50 hover:border-burnt-orange/30 hover:text-burnt-orange transition-all shadow-sm"
                title="Retour au jeu"
            >
                <LayoutDashboard size={18} />
                <span class="font-medium text-sm hidden md:inline"
                    >Retour au jeu</span
                >
            </a>
        </div>
    {/if}
    <div class="flex items-center gap-4">
        {#if $session.data}
            <div class="relative" use:clickOutside={closeMenu}>
                <button
                    class="flex items-center gap-3 pl-4 border-l border-stone-100 hover:bg-stone-100 transition-colors rounded-xl py-1 pr-2 hover:cursor-pointer"
                    onclick={toggleMenu}
                >
                    <div class="text-right hidden md:block">
                        <p class="text-sm font-bold text-dark-gray">
                            {$session.data.user.name}
                        </p>
                        <p class="text-xs text-dark-gray/60">Maître du Jeu</p>
                    </div>
                    {#if $session.data.user.image}
                        <img
                            src={$session.data.user.image}
                            alt={$session.data.user.name}
                            class="w-10 h-10 rounded-full border-2 border-white shadow-sm object-cover"
                        />
                    {:else}
                        <div
                            class="w-10 h-10 rounded-full bg-mustard-yellow/30 flex items-center justify-center text-mustard-yellow font-bold border-2 border-white shadow-sm"
                        >
                            {$session.data.user.name.charAt(0).toUpperCase()}
                        </div>
                    {/if}
                    <ChevronDown
                        size={16}
                        class="text-dark-gray/40 transition-transform duration-200 {isMenuOpen
                            ? 'rotate-180'
                            : ''}"
                    />
                </button>

                {#if isMenuOpen}
                    <div
                        class="absolute right-0 top-full mt-2 w-56 bg-white rounded-xl shadow-lg border border-stone-100 py-2 overflow-hidden z-50 animate-in fade-in slide-in-from-top-2 duration-200"
                    >
                        <a
                            href="/dashboard"
                            class="flex items-center gap-3 px-4 py-2.5 text-sm text-dark-gray hover:bg-stone-50 transition-colors"
                            onclick={closeMenu}
                        >
                            <LayoutDashboard
                                size={18}
                                class="text-dark-gray/60"
                            />
                            Dashboard
                        </a>
                        <a
                            href="/me"
                            class="flex items-center gap-3 px-4 py-2.5 text-sm text-dark-gray hover:bg-stone-50 transition-colors"
                            onclick={closeMenu}
                        >
                            <User size={18} class="text-dark-gray/60" />
                            Profil
                        </a>
                        <div class="h-px bg-stone-100 my-1"></div>
                        <button
                            class="w-full flex items-center gap-3 px-4 py-2.5 text-sm text-red-500 hover:bg-red-50 transition-colors text-left hover:cursor-pointer"
                            onclick={handleLogout}
                        >
                            <LogOut size={18} />
                            Se déconnecter
                        </button>
                    </div>
                {/if}
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
