<script lang="ts">
    import { onMount } from "svelte";
    import { goto } from "$app/navigation";
    import { fade, fly } from "svelte/transition";
    import Header from "$lib/components/Header.svelte";
    import {
        fetchCurrentUserRole,
        fetchAdminStats,
        fetchAdminUsers,
        setAdminUserRole,
        setAdminUserBanned,
        fetchAdminGames,
        archiveAdminGame,
        deleteAdminGame,
        fetchAdminTemplates,
        toggleAdminTemplatePublic,
        deleteAdminTemplate,
        fetchAdminTickets,
        updateAdminTicketStatus,
        fetchAdminUserDetails,
        fetchAdminGameDetails,
        fetchAdminTemplateDetails,
        type AdminStats,
        type AdminUser,
        type AdminGame,
        type AdminTemplate,
        type AdminTicket,
        type AdminUserDetails,
        type AdminGameDetails,
        type AdminTemplateDetails
    } from "$lib/api";
    import {
        Users,
        ShieldAlert,
        ShieldCheck,
        Ban,
        Trash2,
        Archive,
        ArchiveRestore,
        CheckCircle2,
        Lock,
        Globe,
        Ticket,
        LayoutDashboard,
        ChevronRight,
        Loader2,
        Check,
        AlertTriangle,
        X,
        Mail,
        Calendar,
        Award,
        Sparkles,
        Sword,
        Shield,
        Heart,
        Zap,
        Download,
        ExternalLink,
        FileText,
        ChevronDown,
        BookOpen,
        Info
    } from "lucide-svelte";

    // Authentication & Authorization States
    let userRole = $state<'admin' | 'moderator' | null>(null);
    let isBanned = $state(false);
    let loadingRole = $state(true);

    // Active Navigation Tab
    let activeTab = $state<"stats" | "users" | "games" | "templates" | "tickets">("stats");

    // Dashboard Data States
    let stats = $state<AdminStats | null>(null);
    let users = $state<AdminUser[]>([]);
    let games = $state<AdminGame[]>([]);
    let templates = $state<AdminTemplate[]>([]);
    let tickets = $state<AdminTicket[]>([]);

    // Loading & Error States for Data
    let loadingData = $state(false);
    let actionLoading = $state<string | null>(null); // tracks ID of item currently being processed
    let errorMsg = $state<string | null>(null);
    let successMsg = $state<string | null>(null);

    // Filters & Search
    let searchUserQuery = $state("");
    let searchGameQuery = $state("");
    let searchTemplateQuery = $state("");
    let ticketCategoryFilter = $state<string>("ALL");
    let ticketStatusFilter = $state<string>("ALL");
    let templateCategoryFilter = $state<string>("ALL");

    // Detail Drawers States
    let selectedUserId = $state<string | null>(null);
    let selectedUser = $derived(users.find(u => u.id === selectedUserId));
    let userDetails = $state<AdminUserDetails | null>(null);
    let loadingUserDetails = $state(false);

    let selectedGameId = $state<string | null>(null);
    let selectedGame = $derived(games.find(g => g.id === selectedGameId));
    let gameDetails = $state<AdminGameDetails | null>(null);
    let loadingGameDetails = $state(false);

    let selectedTemplate = $state<any>(null); // Stores the template object (raw or virtual item)
    let templateDetails = $state<AdminTemplateDetails | null>(null);
    let loadingTemplateDetails = $state(false);

    let selectedTicketId = $state<string | null>(null);
    let selectedTicket = $derived(tickets.find(t => t.id === selectedTicketId));

    // Drawer controllers
    function closeAllDrawers() {
        selectedUserId = null;
        selectedGameId = null;
        selectedTemplate = null;
        selectedTicketId = null;
        userDetails = null;
        gameDetails = null;
        templateDetails = null;
    }

    async function openUserDetails(userId: string) {
        closeAllDrawers();
        selectedUserId = userId;
        loadingUserDetails = true;
        try {
            userDetails = await fetchAdminUserDetails(userId);
        } catch (err: any) {
            console.error("Failed to load user details:", err);
            triggerError("Impossible de charger les détails de l'utilisateur.");
        } finally {
            loadingUserDetails = false;
        }
    }

    async function openGameDetails(gameId: string) {
        closeAllDrawers();
        selectedGameId = gameId;
        loadingGameDetails = true;
        try {
            gameDetails = await fetchAdminGameDetails(gameId);
        } catch (err: any) {
            console.error("Failed to load game details:", err);
            triggerError("Impossible de charger les détails de la partie.");
        } finally {
            loadingGameDetails = false;
        }
    }

    async function openTemplateDetails(template: any) {
        closeAllDrawers();
        selectedTemplate = template;

        if (template.is_virtual) {
            templateDetails = {
                data: template.data,
                creator_name: template.creator_name
            };
            return;
        }

        loadingTemplateDetails = true;
        try {
            templateDetails = await fetchAdminTemplateDetails(template.id);
        } catch (err: any) {
            console.error("Failed to load template details:", err);
            triggerError("Impossible de charger les détails du modèle.");
        } finally {
            loadingTemplateDetails = false;
        }
    }

    async function openParentTemplateDetails(parentBundleId: string) {
        const parent = templates.find(t => t.id === parentBundleId);
        if (parent) {
            await openTemplateDetails({ ...parent, is_virtual: false });
        } else {
            loadingTemplateDetails = true;
            try {
                const details = await fetchAdminTemplateDetails(parentBundleId);
                const parentTemplate = {
                    id: parentBundleId,
                    name: "Pack Parent",
                    type: 'BUNDLE',
                    creator_name: details.creator_name,
                    is_public: true,
                    uses: 0,
                    created_at: new Date().toISOString(),
                    description: '',
                    data: details.data,
                    is_virtual: false
                };
                selectedTemplate = parentTemplate;
                templateDetails = details;
            } catch (err: any) {
                console.error("Failed to fetch parent bundle template:", err);
                triggerError("Impossible de charger les détails du pack parent.");
            } finally {
                loadingTemplateDetails = false;
            }
        }
    }

    function openTicketDetails(ticketId: string) {
        closeAllDrawers();
        selectedTicketId = ticketId;
    }

    function isImage(url: string | null | undefined): boolean {
        if (!url) return false;
        const lower = url.toLowerCase();
        return lower.endsWith('.png') || lower.endsWith('.jpg') || lower.endsWith('.jpeg') || lower.endsWith('.webp') || lower.endsWith('.gif');
    }

    // Derived filtered lists
    const filteredUsers = $derived(
        users.filter(u => 
            u.name.toLowerCase().includes(searchUserQuery.toLowerCase()) ||
            u.email.toLowerCase().includes(searchUserQuery.toLowerCase())
        )
    );

    const filteredGames = $derived(
        games.filter(g => 
            g.name.toLowerCase().includes(searchGameQuery.toLowerCase()) ||
            g.gm_name.toLowerCase().includes(searchGameQuery.toLowerCase()) ||
            g.gm_email.toLowerCase().includes(searchGameQuery.toLowerCase())
        )
    );

    const processedAdminTemplates = $derived.by(() => {
        let list: any[] = [];
        for (const t of templates) {
            // Keep the original template (the pack/bundle itself)
            list.push({ ...t, is_virtual: false });
            // If it's a BUNDLE, extract its items as virtual templates
            if (t.type === 'BUNDLE' && t.data?.items) {
                t.data.items.forEach((item: any, idx: number) => {
                    list.push({
                        id: `${t.id}-item-${idx}`,
                        parent_bundle_id: t.id,
                        parent_bundle_name: t.name,
                        created_by: t.created_by,
                        creator_name: t.creator_name,
                        name: item.name,
                        description: item.data?.description || `Fait partie du pack "${t.name}"`,
                        type: item.type === 'MONSTER' ? 'MONSTRE' : (item.type === 'PLAYER' ? 'PERSONNAGE' : item.type),
                        data: item.data,
                        is_public: t.is_public,
                        uses: t.uses,
                        created_at: t.created_at,
                        is_virtual: true
                    });
                });
            }
        }

        // Apply filters
        return list.filter(t => {
            // Category Filter
            let matchType = true;
            if (templateCategoryFilter === "BUNDLE") {
                matchType = t.type === "BUNDLE" && !t.is_virtual;
            } else if (templateCategoryFilter !== "ALL") {
                matchType = t.type === templateCategoryFilter;
            }

            // Search Query
            const matchSearch = 
                t.name.toLowerCase().includes(searchTemplateQuery.toLowerCase()) ||
                t.creator_name.toLowerCase().includes(searchTemplateQuery.toLowerCase()) ||
                (t.description && t.description.toLowerCase().includes(searchTemplateQuery.toLowerCase()));

            return matchType && matchSearch;
        });
    });

    const filteredTickets = $derived(
        tickets.filter(t => {
            const matchCat = ticketCategoryFilter === "ALL" || t.category === ticketCategoryFilter;
            const matchStatus = ticketStatusFilter === "ALL" || t.status === ticketStatusFilter;
            return matchCat && matchStatus;
        })
    );

    // Authorization Guard & Data Loader
    onMount(async () => {
        try {
            const roleData = await fetchCurrentUserRole();
            userRole = roleData.role;
            isBanned = roleData.is_banned;

            if (isBanned || (userRole !== "admin" && userRole !== "moderator")) {
                // Not authorized
                loadingRole = false;
                return;
            }

            loadingRole = false;
            // Load initial tab data
            await loadTabValues(activeTab);
        } catch (e) {
            console.error("Authorization check failed", e);
            loadingRole = false;
        }
    });

    async function loadTabValues(tab: typeof activeTab) {
        loadingData = true;
        errorMsg = null;
        try {
            if (tab === "stats") {
                stats = await fetchAdminStats();
            } else if (tab === "users") {
                users = await fetchAdminUsers();
            } else if (tab === "games") {
                games = await fetchAdminGames();
            } else if (tab === "templates") {
                templates = await fetchAdminTemplates();
            } else if (tab === "tickets") {
                tickets = await fetchAdminTickets();
            }
        } catch (err: any) {
            console.error(`Failed to load data for tab ${tab}`, err);
            errorMsg = err.message || "Erreur lors du chargement des données.";
        } finally {
            loadingData = false;
        }
    }

    // Tab Change handler
    async function handleTabChange(newTab: typeof activeTab) {
        activeTab = newTab;
        errorMsg = null;
        successMsg = null;
        await loadTabValues(newTab);
    }

    // Toast helpers
    function triggerSuccess(msg: string) {
        successMsg = msg;
        setTimeout(() => { if (successMsg === msg) successMsg = null; }, 3000);
    }

    function triggerError(msg: string) {
        errorMsg = msg;
        setTimeout(() => { if (errorMsg === msg) errorMsg = null; }, 4000);
    }

    // Role management (Admins only)
    async function handleRoleChange(userId: string, newRole: 'admin' | 'moderator' | null) {
        if (userRole !== 'admin') return;
        actionLoading = `role-${userId}`;
        try {
            await setAdminUserRole(userId, newRole);
            users = users.map(u => u.id === userId ? { ...u, role: newRole } : u);
            triggerSuccess("Rôle de l'utilisateur mis à jour avec succès.");
            if (selectedUserId === userId && userDetails) {
                // Refresh drawer details
                await openUserDetails(userId);
            }
        } catch (err: any) {
            triggerError(err.message || "Impossible de mettre à jour le rôle.");
        } finally {
            actionLoading = null;
        }
    }

    // Ban/Unban management (Admins only)
    async function handleBanToggle(userId: string, currentBanned: boolean) {
        if (userRole !== 'admin') return;
        const confirmMsg = currentBanned 
            ? "Voulez-vous vraiment lever la suspension de cet utilisateur ?" 
            : "Voulez-vous vraiment suspendre cet utilisateur ? Il ne pourra plus se connecter.";
        if (!confirm(confirmMsg)) return;

        actionLoading = `ban-${userId}`;
        try {
            await setAdminUserBanned(userId, !currentBanned);
            users = users.map(u => u.id === userId ? { ...u, is_banned: !currentBanned } : u);
            triggerSuccess(currentBanned ? "Suspension levée." : "Utilisateur suspendu.");
            if (selectedUserId === userId && userDetails) {
                // Refresh drawer details
                await openUserDetails(userId);
            }
        } catch (err: any) {
            triggerError(err.message || "Échec de l'action de bannissement.");
        } finally {
            actionLoading = null;
        }
    }

    // Game Actions (Admins & Moderators)
    async function handleGameArchive(gameId: string, currentActive: boolean) {
        actionLoading = `archive-${gameId}`;
        try {
            await archiveAdminGame(gameId, !currentActive);
            games = games.map(g => g.id === gameId ? { ...g, is_active: !currentActive } : g);
            triggerSuccess(currentActive ? "Partie archivée." : "Partie restaurée.");
            if (selectedGameId === gameId) {
                // Refresh drawer
                await openGameDetails(gameId);
            }
        } catch (err: any) {
            triggerError(err.message || "Échec de l'archivage.");
        } finally {
            actionLoading = null;
        }
    }

    async function handleGameDelete(gameId: string) {
        if (!confirm("ATTENTION : Cette action supprimera définitivement la partie, ses plateaux, ses personnages et ses fichiers de stockage. Continuer ?")) return;
        
        actionLoading = `delete-${gameId}`;
        try {
            await deleteAdminGame(gameId);
            games = games.filter(g => g.id !== gameId);
            triggerSuccess("Partie supprimée définitivement.");
            if (selectedGameId === gameId) {
                closeAllDrawers();
            }
        } catch (err: any) {
            triggerError(err.message || "Échec de la suppression.");
        } finally {
            actionLoading = null;
        }
    }

    // Template Actions (Admins & Moderators)
    async function handleTemplatePublicToggle(templateId: string, currentPublic: boolean) {
        actionLoading = `public-${templateId}`;
        try {
            await toggleAdminTemplatePublic(templateId, !currentPublic);
            templates = templates.map(t => t.id === templateId ? { ...t, is_public: !currentPublic } : t);
            triggerSuccess(currentPublic ? "Modèle rendu privé." : "Modèle publié sur le marché.");
            
            // Update in selectedTemplate and processed list
            if (selectedTemplate && selectedTemplate.id === templateId) {
                selectedTemplate = { ...selectedTemplate, is_public: !currentPublic };
            }
        } catch (err: any) {
            triggerError(err.message || "Échec de la modification.");
        } finally {
            actionLoading = null;
        }
    }

    async function handleTemplateDelete(templateId: string) {
        if (!confirm("Voulez-vous vraiment supprimer ce modèle de la boutique ?")) return;
        
        actionLoading = `delete-${templateId}`;
        try {
            await deleteAdminTemplate(templateId);
            templates = templates.filter(t => t.id !== templateId);
            triggerSuccess("Modèle supprimé avec succès.");
            if (selectedTemplate && (selectedTemplate.id === templateId || selectedTemplate.parent_bundle_id === templateId)) {
                closeAllDrawers();
            }
        } catch (err: any) {
            triggerError(err.message || "Échec de la suppression.");
        } finally {
            actionLoading = null;
        }
    }

    // Support Ticket Actions (Admins & Moderators)
    async function handleTicketStatusChange(ticketId: string, newStatus: string) {
        actionLoading = `ticket-${ticketId}`;
        try {
            await updateAdminTicketStatus(ticketId, newStatus);
            tickets = tickets.map(t => t.id === ticketId ? { ...t, status: newStatus as any } : t);
            triggerSuccess("Statut du ticket mis à jour.");
        } catch (err: any) {
            triggerError(err.message || "Échec de la mise à jour.");
        } finally {
            actionLoading = null;
        }
    }
</script>

<svelte:head>
    <title>Administration — QuestHub</title>
</svelte:head>

<div class="min-h-screen bg-stone-50 flex flex-col font-sans text-stone-850">
    <Header />

    {#if loadingRole}
        <div class="flex-grow flex items-center justify-center">
            <Loader2 class="animate-spin text-burnt-orange" size={40} />
        </div>
    {:else if isBanned || (userRole !== "admin" && userRole !== "moderator")}
        <!-- Unauthorized view -->
        <div class="flex-grow flex flex-col items-center justify-center p-6 text-center max-w-md mx-auto space-y-4">
            <div class="w-16 h-16 bg-red-50 rounded-full flex items-center justify-center text-red-500 shadow-sm border border-red-100">
                <ShieldAlert size={32} />
            </div>
            <h1 class="text-2xl font-bold text-dark-gray font-display">Accès Interdit</h1>
            <p class="text-stone-500 text-sm leading-relaxed">
                Vous n'avez pas les droits nécessaires pour accéder à cette page. Cet espace est réservé aux administrateurs et modérateurs de QuestHub.
            </p>
            <button
                onclick={() => goto("/dashboard")}
                class="px-5 py-2.5 bg-burnt-orange text-white font-medium rounded-xl shadow-md hover:bg-opacity-90 transition-all cursor-pointer"
            >
                Retour au tableau de bord
            </button>
        </div>
    {:else}
        <!-- Toast Alerts -->
        {#if successMsg}
            <div class="fixed bottom-4 right-4 z-100 bg-emerald-500 text-white px-4 py-3 rounded-xl shadow-lg flex items-center gap-2 animate-in slide-in-from-bottom-2 duration-300">
                <Check size={16} />
                <span class="text-sm font-medium">{successMsg}</span>
            </div>
        {/if}
        {#if errorMsg}
            <div class="fixed bottom-4 right-4 z-100 bg-red-500 text-white px-4 py-3 rounded-xl shadow-lg flex items-center gap-2 animate-in slide-in-from-bottom-2 duration-300">
                <AlertTriangle size={16} />
                <span class="text-sm font-medium">{errorMsg}</span>
            </div>
        {/if}

        <div class="flex-grow flex flex-col lg:flex-row">
            <!-- Admin Sidebar Navigation -->
            <aside class="w-full lg:w-64 bg-white border-r border-stone-200 p-4 shrink-0 flex lg:flex-col gap-1 lg:gap-2 overflow-x-auto lg:overflow-x-visible">
                <div class="hidden lg:block px-3 py-2 mb-4">
                    <h2 class="text-xs font-bold text-stone-400 uppercase tracking-widest">Modération</h2>
                    <p class="text-[10px] text-stone-400 mt-0.5">Connecté en tant que <span class="font-bold text-burnt-orange">{userRole === 'admin' ? 'Admin' : 'Modérateur'}</span></p>
                </div>

                <button
                    onclick={() => handleTabChange("stats")}
                    class="flex items-center gap-3 px-4 py-2.5 rounded-xl text-sm font-medium transition-all text-left whitespace-nowrap cursor-pointer
                    {activeTab === 'stats' ? 'bg-burnt-orange/5 text-burnt-orange font-bold shadow-xs' : 'text-stone-500 hover:bg-stone-50'}"
                >
                    <LayoutDashboard size={18} />
                    <span>Vue d'ensemble</span>
                </button>

                <button
                    onclick={() => handleTabChange("users")}
                    class="flex items-center gap-3 px-4 py-2.5 rounded-xl text-sm font-medium transition-all text-left whitespace-nowrap cursor-pointer
                    {activeTab === 'users' ? 'bg-burnt-orange/5 text-burnt-orange font-bold shadow-xs' : 'text-stone-500 hover:bg-stone-50'}"
                >
                    <Users size={18} />
                    <span>Utilisateurs</span>
                </button>

                <button
                    onclick={() => handleTabChange("games")}
                    class="flex items-center gap-3 px-4 py-2.5 rounded-xl text-sm font-medium transition-all text-left whitespace-nowrap cursor-pointer
                    {activeTab === 'games' ? 'bg-burnt-orange/5 text-burnt-orange font-bold shadow-xs' : 'text-stone-500 hover:bg-stone-50'}"
                >
                    <Archive size={18} />
                    <span>Parties (Tables)</span>
                </button>

                <button
                    onclick={() => handleTabChange("templates")}
                    class="flex items-center gap-3 px-4 py-2.5 rounded-xl text-sm font-medium transition-all text-left whitespace-nowrap cursor-pointer
                    {activeTab === 'templates' ? 'bg-burnt-orange/5 text-burnt-orange font-bold shadow-xs' : 'text-stone-500 hover:bg-stone-50'}"
                >
                    <Globe size={18} />
                    <span>Marketplace</span>
                </button>

                <button
                    onclick={() => handleTabChange("tickets")}
                    class="flex items-center gap-3 px-4 py-2.5 rounded-xl text-sm font-medium transition-all text-left whitespace-nowrap cursor-pointer
                    {activeTab === 'tickets' ? 'bg-burnt-orange/5 text-burnt-orange font-bold shadow-xs' : 'text-stone-500 hover:bg-stone-50'}"
                >
                    <Ticket size={18} />
                    <span>Tickets & Support</span>
                    {#if stats?.pending_tickets && stats.pending_tickets > 0}
                        <span class="ml-auto bg-red-500 text-white text-[10px] font-bold px-1.5 py-0.5 rounded-full">{stats.pending_tickets}</span>
                    {/if}
                </button>
            </aside>

            <!-- Admin Tab Content -->
            <main class="flex-1 p-6 md:p-8 overflow-y-auto">
                {#if loadingData}
                    <div class="h-64 flex items-center justify-center">
                        <Loader2 class="animate-spin text-stone-400" size={24} />
                    </div>
                {:else}
                    <!-- 1. STATS TAB -->
                    {#if activeTab === "stats" && stats}
                        <div class="space-y-6">
                            <h2 class="text-xl font-bold text-dark-gray font-display">Vue d'ensemble du réseau</h2>
                            
                            <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
                                <div class="bg-white p-5 rounded-2xl border border-stone-150 shadow-xs flex items-center gap-4">
                                    <div class="w-12 h-12 rounded-xl bg-blue-50 text-blue-500 flex items-center justify-center">
                                        <Users size={24} />
                                    </div>
                                    <div>
                                        <p class="text-xs text-stone-400 font-medium uppercase tracking-wider">Membres total</p>
                                        <h3 class="text-2xl font-bold text-stone-800 mt-1">{stats.total_users}</h3>
                                    </div>
                                </div>

                                <div class="bg-white p-5 rounded-2xl border border-stone-150 shadow-xs flex items-center gap-4">
                                    <div class="w-12 h-12 rounded-xl bg-amber-50 text-amber-500 flex items-center justify-center">
                                        <Archive size={24} />
                                    </div>
                                    <div>
                                        <p class="text-xs text-stone-400 font-medium uppercase tracking-wider">Parties / Tables</p>
                                        <h3 class="text-2xl font-bold text-stone-800 mt-1">{stats.total_games}</h3>
                                    </div>
                                </div>

                                <div class="bg-white p-5 rounded-2xl border border-stone-150 shadow-xs flex items-center gap-4">
                                    <div class="w-12 h-12 rounded-xl bg-emerald-50 text-emerald-500 flex items-center justify-center">
                                        <Globe size={24} />
                                    </div>
                                    <div>
                                        <p class="text-xs text-stone-400 font-medium uppercase tracking-wider">Modèles publics</p>
                                        <h3 class="text-2xl font-bold text-stone-800 mt-1">{stats.total_templates}</h3>
                                    </div>
                                </div>

                                <div class="bg-white p-5 rounded-2xl border border-stone-150 shadow-xs flex items-center gap-4">
                                    <div class="w-12 h-12 rounded-xl bg-red-50 text-red-500 flex items-center justify-center">
                                        <Ticket size={24} />
                                    </div>
                                    <div>
                                        <p class="text-xs text-stone-400 font-medium uppercase tracking-wider">Tickets en attente</p>
                                        <h3 class="text-2xl font-bold text-stone-800 mt-1">{stats.pending_tickets}</h3>
                                    </div>
                                </div>
                            </div>

                            <div class="bg-white p-6 rounded-2xl border border-stone-150 shadow-xs">
                                <h3 class="font-bold text-stone-800 mb-4 font-display">Informations Systèmes</h3>
                                <p class="text-stone-500 text-sm leading-relaxed">
                                    QuestHub est hébergé en mode distribué. L'accès à ce tableau de bord est restreint par les règles de sécurité Row Level Security (RLS) et des procédures stockées exécutées de manière isolée et protégée sur le serveur.
                                </p>
                                <div class="mt-4 flex gap-4 text-xs font-mono text-stone-400">
                                    <div>Version : <span class="text-stone-600 font-bold">1.2.0</span></div>
                                    <div>Status : <span class="text-emerald-500 font-bold">Opérationnel</span></div>
                                    <div>Utilisateurs suspendus : <span class="text-red-500 font-bold">{stats.banned_users}</span></div>
                                </div>
                            </div>
                        </div>

                    <!-- 2. USERS TAB -->
                    {:else if activeTab === "users"}
                        <div class="space-y-6">
                            <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
                                <h2 class="text-xl font-bold text-dark-gray font-display">Gestion des Utilisateurs</h2>
                                <input
                                    type="text"
                                    bind:value={searchUserQuery}
                                    placeholder="Rechercher par nom ou email..."
                                    class="px-4 py-2 bg-white rounded-xl border border-stone-200 text-sm focus:outline-none focus:border-burnt-orange w-full sm:w-64 shadow-xs"
                                />
                            </div>

                            <div class="bg-white rounded-2xl border border-stone-150 shadow-xs overflow-hidden">
                                <div class="overflow-x-auto">
                                    <table class="w-full text-left text-sm border-collapse">
                                        <thead>
                                            <tr class="bg-stone-50 border-b border-stone-150 text-stone-400 font-bold text-xs uppercase tracking-wider">
                                                <th class="px-6 py-3.5">Nom / Email</th>
                                                <th class="px-6 py-3.5">Inscription</th>
                                                <th class="px-6 py-3.5">Statut</th>
                                                <th class="px-6 py-3.5">Rôle</th>
                                                {#if userRole === 'admin'}
                                                    <th class="px-6 py-3.5 text-right">Actions</th>
                                                {/if}
                                            </tr>
                                        </thead>
                                        <tbody class="divide-y divide-stone-100">
                                            {#each filteredUsers as u (u.id)}
                                                <tr 
                                                    onclick={() => openUserDetails(u.id)}
                                                    class="hover:bg-stone-50/50 transition-colors cursor-pointer {u.is_banned ? 'bg-red-50/20' : ''}"
                                                >
                                                    <td class="px-6 py-4">
                                                        <div class="flex flex-col leading-tight">
                                                            <span class="font-bold text-stone-850">{u.name}</span>
                                                            <span class="text-xs text-stone-400">{u.email}</span>
                                                        </div>
                                                    </td>
                                                    <td class="px-6 py-4 text-xs text-stone-500 font-mono">
                                                        {new Date(u.created_at).toLocaleDateString("fr-FR")}
                                                    </td>
                                                    <td class="px-6 py-4">
                                                        {#if u.is_banned}
                                                            <span class="inline-flex items-center gap-1 px-2.5 py-1 rounded-full text-xs font-semibold bg-red-100 text-red-700">
                                                                <Ban size={10} /> Suspendu
                                                            </span>
                                                        {:else}
                                                            <span class="inline-flex items-center gap-1 px-2.5 py-1 rounded-full text-xs font-semibold bg-emerald-100 text-emerald-700">
                                                                Actif
                                                            </span>
                                                        {/if}
                                                    </td>
                                                    <td class="px-6 py-4" onclick={(e) => e.stopPropagation()}>
                                                        {#if userRole === 'admin'}
                                                            <select
                                                                value={u.role || ""}
                                                                onchange={(e) => handleRoleChange(u.id, (e.target as HTMLSelectElement).value === "" ? null : (e.target as HTMLSelectElement).value as any)}
                                                                disabled={actionLoading === `role-${u.id}`}
                                                                class="bg-white border border-stone-200 rounded-lg px-2 py-1 text-xs focus:outline-none focus:border-burnt-orange shadow-xs disabled:opacity-50"
                                                            >
                                                                <option value="">Utilisateur</option>
                                                                <option value="moderator">Modérateur</option>
                                                                <option value="admin">Admin</option>
                                                            </select>
                                                        {:else}
                                                            <span class="text-xs font-medium text-stone-700">
                                                                {u.role === 'admin' ? 'Administrateur' : u.role === 'moderator' ? 'Modérateur' : 'Utilisateur'}
                                                            </span>
                                                        {/if}
                                                    </td>
                                                    {#if userRole === 'admin'}
                                                        <td class="px-6 py-4 text-right" onclick={(e) => e.stopPropagation()}>
                                                            <div class="flex items-center justify-end gap-2">
                                                                <button
                                                                    onclick={() => handleBanToggle(u.id, u.is_banned)}
                                                                    disabled={actionLoading !== null}
                                                                    class="p-1.5 rounded-lg border transition-colors shadow-xs cursor-pointer disabled:opacity-50
                                                                    {u.is_banned 
                                                                        ? 'border-emerald-100 bg-emerald-50 text-emerald-600 hover:bg-emerald-100' 
                                                                        : 'border-red-100 bg-red-50 text-red-600 hover:bg-red-100'}"
                                                                    title={u.is_banned ? "Lever la suspension" : "Suspendre le compte"}
                                                                >
                                                                    {#if actionLoading === `ban-${u.id}`}
                                                                        <Loader2 size={14} class="animate-spin" />
                                                                    {:else}
                                                                        <Ban size={14} />
                                                                    {/if}
                                                                </button>
                                                            </div>
                                                        </td>
                                                    {/if}
                                                </tr>
                                            {:else}
                                                <tr>
                                                    <td colspan="5" class="px-6 py-8 text-center text-stone-400">
                                                        Aucun utilisateur trouvé.
                                                    </td>
                                                </tr>
                                            {/each}
                                        </tbody>
                                    </table>
                                </div>
                            </div>
                        </div>

                    <!-- 3. GAMES TAB -->
                    {:else if activeTab === "games"}
                        <div class="space-y-6">
                            <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
                                <h2 class="text-xl font-bold text-dark-gray font-display">Modération des parties</h2>
                                <input
                                    type="text"
                                    bind:value={searchGameQuery}
                                    placeholder="Rechercher par nom de partie ou MJ..."
                                    class="px-4 py-2 bg-white rounded-xl border border-stone-200 text-sm focus:outline-none focus:border-burnt-orange w-full sm:w-64 shadow-xs"
                                />
                            </div>

                            <div class="bg-white rounded-2xl border border-stone-150 shadow-xs overflow-hidden">
                                <div class="overflow-x-auto">
                                    <table class="w-full text-left text-sm border-collapse">
                                        <thead>
                                            <tr class="bg-stone-50 border-b border-stone-150 text-stone-400 font-bold text-xs uppercase tracking-wider">
                                                <th class="px-6 py-3.5">Nom de la Table</th>
                                                <th class="px-6 py-3.5">Maître du Jeu (MJ)</th>
                                                <th class="px-6 py-3.5">Joueurs</th>
                                                <th class="px-6 py-3.5">Statut</th>
                                                <th class="px-6 py-3.5 text-right">Actions</th>
                                            </tr>
                                        </thead>
                                        <tbody class="divide-y divide-stone-100">
                                            {#each filteredGames as g (g.id)}
                                                <tr 
                                                    onclick={() => openGameDetails(g.id)}
                                                    class="hover:bg-stone-50/50 transition-colors cursor-pointer {!g.is_active ? 'bg-stone-100/50 opacity-70' : ''}"
                                                >
                                                    <td class="px-6 py-4">
                                                        <div class="flex flex-col">
                                                            <span class="font-bold text-stone-850">{g.name}</span>
                                                            <span class="text-[10px] text-stone-450 font-mono mt-0.5">{g.id}</span>
                                                        </div>
                                                    </td>
                                                    <td class="px-6 py-4">
                                                        <div class="flex flex-col leading-tight">
                                                            <span class="font-medium text-stone-850">{g.gm_name}</span>
                                                            <span class="text-xs text-stone-400">{g.gm_email}</span>
                                                        </div>
                                                    </td>
                                                    <td class="px-6 py-4 font-bold text-stone-700">
                                                        {g.player_count}
                                                    </td>
                                                    <td class="px-6 py-4">
                                                        {#if !g.is_active}
                                                            <span class="inline-flex items-center gap-1 px-2.5 py-1 rounded-full text-xs font-semibold bg-stone-100 text-stone-600">
                                                                Archivée
                                                            </span>
                                                        {:else}
                                                            <span class="inline-flex items-center gap-1 px-2.5 py-1 rounded-full text-xs font-semibold bg-emerald-100 text-emerald-700">
                                                                Actif
                                                            </span>
                                                        {/if}
                                                    </td>
                                                    <td class="px-6 py-4 text-right" onclick={(e) => e.stopPropagation()}>
                                                        <div class="flex items-center justify-end gap-2">
                                                            <button
                                                                onclick={() => handleGameArchive(g.id, g.is_active)}
                                                                disabled={actionLoading !== null}
                                                                class="p-1.5 rounded-lg border border-stone-200 bg-white text-stone-600 hover:bg-stone-50 transition-colors shadow-xs cursor-pointer disabled:opacity-50"
                                                                title={g.is_active ? "Archiver la partie" : "Restaurer la partie"}
                                                            >
                                                                {#if actionLoading === `archive-${g.id}`}
                                                                    <Loader2 size={14} class="animate-spin" />
                                                                {:else if g.is_active}
                                                                    <Archive size={14} />
                                                                {:else}
                                                                    <ArchiveRestore size={14} />
                                                                {/if}
                                                            </button>

                                                            <button
                                                                onclick={() => handleGameDelete(g.id)}
                                                                disabled={actionLoading !== null}
                                                                class="p-1.5 rounded-lg border border-red-100 bg-red-50 text-red-600 hover:bg-red-100 transition-colors shadow-xs cursor-pointer disabled:opacity-50"
                                                                title="Supprimer définitivement"
                                                            >
                                                                {#if actionLoading === `delete-${g.id}`}
                                                                    <Loader2 size={14} class="animate-spin" />
                                                                {:else}
                                                                    <Trash2 size={14} />
                                                                {/if}
                                                            </button>
                                                        </div>
                                                    </td>
                                                </tr>
                                            {:else}
                                                <tr>
                                                    <td colspan="5" class="px-6 py-8 text-center text-stone-400">
                                                        Aucune partie trouvée.
                                                    </td>
                                                </tr>
                                            {/each}
                                        </tbody>
                                    </table>
                                </div>
                            </div>
                        </div>

                    <!-- 4. TEMPLATES TAB -->
                    {:else if activeTab === "templates"}
                        <div class="space-y-6">
                            <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
                                <h2 class="text-xl font-bold text-dark-gray font-display">Modération de la Marketplace</h2>
                                <div class="flex flex-wrap items-center gap-3 w-full sm:w-auto">
                                    <select
                                        bind:value={templateCategoryFilter}
                                        class="bg-white border border-stone-200 rounded-xl px-3 py-2 text-sm focus:outline-none focus:border-burnt-orange shadow-xs text-stone-600 font-bold cursor-pointer"
                                    >
                                        <option value="ALL">Tous les types</option>
                                        <option value="BUNDLE">Packs / Bundles</option>
                                        <option value="PERSONNAGE">Personnages</option>
                                        <option value="PNJ">PNJs</option>
                                        <option value="MONSTRE">Monstres</option>
                                    </select>
                                    <input
                                        type="text"
                                        bind:value={searchTemplateQuery}
                                        placeholder="Rechercher par nom de modèle ou créateur..."
                                        class="px-4 py-2 bg-white rounded-xl border border-stone-200 text-sm focus:outline-none focus:border-burnt-orange w-full sm:w-64 shadow-xs"
                                    />
                                </div>
                            </div>

                            <div class="bg-white rounded-2xl border border-stone-150 shadow-xs overflow-hidden">
                                <div class="overflow-x-auto">
                                    <table class="w-full text-left text-sm border-collapse">
                                        <thead>
                                            <tr class="bg-stone-50 border-b border-stone-150 text-stone-400 font-bold text-xs uppercase tracking-wider">
                                                <th class="px-6 py-3.5">Nom du Modèle</th>
                                                <th class="px-6 py-3.5">Type</th>
                                                <th class="px-6 py-3.5">Créateur</th>
                                                <th class="px-6 py-3.5">Visibilité</th>
                                                <th class="px-6 py-3.5 text-center">Importations</th>
                                                <th class="px-6 py-3.5 text-right">Actions</th>
                                            </tr>
                                        </thead>
                                        <tbody class="divide-y divide-stone-100">
                                            {#each processedAdminTemplates as t (t.id)}
                                                <tr 
                                                    onclick={() => openTemplateDetails(t)}
                                                    class="hover:bg-stone-50/50 transition-colors cursor-pointer"
                                                >
                                                    <td class="px-6 py-4">
                                                        <div class="flex flex-col">
                                                            <span class="font-bold text-stone-850">{t.name}</span>
                                                            {#if t.is_virtual}
                                                                <span class="text-[10px] text-burnt-orange font-bold mt-0.5">📦 Pack : {t.parent_bundle_name}</span>
                                                            {/if}
                                                        </div>
                                                    </td>
                                                    <td class="px-6 py-4">
                                                        <span class="text-[10px] bg-stone-100 px-2 py-0.5 rounded text-stone-600 font-bold uppercase tracking-wider">
                                                            {t.type === 'BUNDLE' ? 'Pack' : (t.type === 'PERSONNAGE' ? 'Perso' : t.type)}
                                                        </span>
                                                    </td>
                                                    <td class="px-6 py-4 text-stone-700 font-medium">
                                                        {t.creator_name}
                                                    </td>
                                                    <td class="px-6 py-4">
                                                        {#if t.is_virtual}
                                                            <span class="inline-flex items-center gap-1 px-2.5 py-1 rounded-full text-xs font-semibold bg-stone-100 text-stone-500">
                                                                Hérité du pack
                                                            </span>
                                                        {:else if t.is_public}
                                                            <span class="inline-flex items-center gap-1 px-2.5 py-1 rounded-full text-xs font-semibold bg-emerald-100 text-emerald-700">
                                                                <Globe size={10} /> Public
                                                            </span>
                                                        {:else}
                                                            <span class="inline-flex items-center gap-1 px-2.5 py-1 rounded-full text-xs font-semibold bg-stone-100 text-stone-600">
                                                                <Lock size={10} /> Privé
                                                            </span>
                                                        {/if}
                                                    </td>
                                                    <td class="px-6 py-4 text-center font-bold text-stone-700">
                                                        {t.is_virtual ? '—' : t.uses}
                                                    </td>
                                                    <td class="px-6 py-4 text-right" onclick={(e) => e.stopPropagation()}>
                                                        <div class="flex items-center justify-end gap-2">
                                                            {#if t.is_virtual}
                                                                <button
                                                                    onclick={() => openParentTemplateDetails(t.parent_bundle_id)}
                                                                    class="p-1.5 rounded-lg border border-stone-200 bg-white text-stone-600 hover:bg-stone-50 transition-colors shadow-xs cursor-pointer flex items-center gap-1 text-xs font-medium"
                                                                    title="Gérer le pack parent"
                                                                >
                                                                    <BookOpen size={14} />
                                                                    <span class="hidden md:inline">Voir le pack</span>
                                                                </button>
                                                            {:else}
                                                                <button
                                                                    onclick={() => handleTemplatePublicToggle(t.id, t.is_public)}
                                                                    disabled={actionLoading !== null}
                                                                    class="p-1.5 rounded-lg border border-stone-200 bg-white text-stone-600 hover:bg-stone-50 transition-colors shadow-xs cursor-pointer disabled:opacity-50"
                                                                    title={t.is_public ? "Masquer de la boutique (rendre privé)" : "Publier sur la boutique"}
                                                                >
                                                                    {#if actionLoading === `public-${t.id}`}
                                                                        <Loader2 size={14} class="animate-spin" />
                                                                    {:else if t.is_public}
                                                                        <Lock size={14} />
                                                                    {:else}
                                                                        <Globe size={14} />
                                                                    {/if}
                                                                </button>

                                                                <button
                                                                    onclick={() => handleTemplateDelete(t.id)}
                                                                    disabled={actionLoading !== null}
                                                                    class="p-1.5 rounded-lg border border-red-100 bg-red-50 text-red-600 hover:bg-red-100 transition-colors shadow-xs cursor-pointer disabled:opacity-50"
                                                                    title="Supprimer définitivement"
                                                                >
                                                                    {#if actionLoading === `delete-${t.id}`}
                                                                        <Loader2 size={14} class="animate-spin" />
                                                                    {:else}
                                                                        <Trash2 size={14} />
                                                                    {/if}
                                                                </button>
                                                            {/if}
                                                        </div>
                                                    </td>
                                                </tr>
                                            {:else}
                                                <tr>
                                                    <td colspan="6" class="px-6 py-8 text-center text-stone-400">
                                                        Aucun modèle de la marketplace trouvé.
                                                    </td>
                                                </tr>
                                            {/each}
                                        </tbody>
                                    </table>
                                </div>
                            </div>
                        </div>

                    <!-- 5. TICKETS TAB -->
                    {:else if activeTab === "tickets"}
                        <div class="space-y-6">
                            <div class="flex flex-col md:flex-row md:items-center justify-between gap-4">
                                <h2 class="text-xl font-bold text-dark-gray font-display">Tickets de Support & Contact</h2>
                                
                                <div class="flex flex-wrap items-center gap-3">
                                    <!-- Category filter -->
                                    <select
                                        bind:value={ticketCategoryFilter}
                                        class="bg-white border border-stone-200 rounded-xl px-3 py-2 text-sm focus:outline-none focus:border-burnt-orange shadow-xs text-stone-600 font-bold cursor-pointer"
                                    >
                                        <option value="ALL">Toutes catégories</option>
                                        <option value="CONTACT">Contact / Question</option>
                                        <option value="BUG">Rapport de Bug</option>
                                        <option value="RECLAMATION">Réclamation</option>
                                    </select>

                                    <!-- Status filter -->
                                    <select
                                        bind:value={ticketStatusFilter}
                                        class="bg-white border border-stone-200 rounded-xl px-3 py-2 text-sm focus:outline-none focus:border-burnt-orange shadow-xs text-stone-600 font-bold cursor-pointer"
                                    >
                                        <option value="ALL">Tous les statuts</option>
                                        <option value="PENDING">En attente</option>
                                        <option value="IN_PROGRESS">En cours</option>
                                        <option value="RESOLVED">Résolu</option>
                                        <option value="CLOSED">Fermé</option>
                                    </select>
                                </div>
                            </div>

                            <div class="bg-white rounded-2xl border border-stone-150 shadow-xs overflow-hidden">
                                <div class="overflow-x-auto">
                                    <table class="w-full text-left text-sm border-collapse">
                                        <thead>
                                            <tr class="bg-stone-50 border-b border-stone-150 text-stone-400 font-bold text-xs uppercase tracking-wider">
                                                <th class="px-6 py-3.5">Sujet / Expéditeur</th>
                                                <th class="px-6 py-3.5">Catégorie</th>
                                                <th class="px-6 py-3.5">Soumission</th>
                                                <th class="px-6 py-3.5">Statut</th>
                                                <th class="px-6 py-3.5 text-right">Actions</th>
                                            </tr>
                                        </thead>
                                        <tbody class="divide-y divide-stone-100">
                                            {#each filteredTickets as t (t.id)}
                                                <tr 
                                                    onclick={() => openTicketDetails(t.id)}
                                                    class="hover:bg-stone-50/50 transition-colors cursor-pointer"
                                                >
                                                    <td class="px-6 py-4">
                                                        <div class="flex flex-col leading-tight">
                                                            <span class="font-bold text-stone-850">{t.subject}</span>
                                                            <span class="text-xs text-stone-400">{t.user_name || "Anonyme"} ({t.email})</span>
                                                        </div>
                                                    </td>
                                                    <td class="px-6 py-4">
                                                        <span class="inline-flex items-center gap-1 px-2.5 py-1 rounded-full text-xs font-semibold
                                                            {t.category === 'BUG' ? 'bg-red-100 text-red-700' : ''}
                                                            {t.category === 'CONTACT' ? 'bg-blue-100 text-blue-700' : ''}
                                                            {t.category === 'RECLAMATION' ? 'bg-amber-100 text-amber-700' : ''}"
                                                        >
                                                            {t.category === 'BUG' ? 'Bug' : t.category === 'CONTACT' ? 'Contact' : 'Réclamation'}
                                                        </span>
                                                    </td>
                                                    <td class="px-6 py-4 text-xs text-stone-500 font-mono">
                                                        {new Date(t.created_at).toLocaleDateString("fr-FR")} à {new Date(t.created_at).toLocaleTimeString("fr-FR", { hour: '2-digit', minute: '2-digit' })}
                                                    </td>
                                                    <td class="px-6 py-4">
                                                        <span class="inline-flex items-center gap-1 px-2.5 py-1 rounded-full text-xs font-semibold
                                                            {t.status === 'PENDING' ? 'bg-red-100 text-red-700' : ''}
                                                            {t.status === 'IN_PROGRESS' ? 'bg-blue-100 text-blue-700' : ''}
                                                            {t.status === 'RESOLVED' ? 'bg-emerald-100 text-emerald-700' : ''}
                                                            {t.status === 'CLOSED' ? 'bg-stone-100 text-stone-700' : ''}"
                                                        >
                                                            {t.status === 'PENDING' ? 'En attente' : t.status === 'IN_PROGRESS' ? 'En cours' : t.status === 'RESOLVED' ? 'Résolu' : 'Fermé'}
                                                        </span>
                                                    </td>
                                                    <td class="px-6 py-4 text-right" onclick={(e) => e.stopPropagation()}>
                                                        <div class="flex items-center justify-end gap-2">
                                                            <select
                                                                value={t.status}
                                                                onchange={(e) => handleTicketStatusChange(t.id, (e.target as HTMLSelectElement).value)}
                                                                disabled={actionLoading === `ticket-${t.id}`}
                                                                class="bg-white border border-stone-200 rounded-lg px-2.5 py-1 text-xs font-semibold focus:outline-none focus:border-burnt-orange shadow-xs disabled:opacity-50"
                                                            >
                                                                <option value="PENDING">En attente</option>
                                                                <option value="IN_PROGRESS">En cours</option>
                                                                <option value="RESOLVED">Résolu</option>
                                                                <option value="CLOSED">Fermé</option>
                                                            </select>
                                                        </div>
                                                    </td>
                                                </tr>
                                            {:else}
                                                <tr>
                                                    <td colspan="5" class="px-6 py-8 text-center text-stone-400">
                                                        Aucun ticket trouvé.
                                                    </td>
                                                </tr>
                                            {/each}
                                        </tbody>
                                    </table>
                                </div>
                            </div>
                        </div>
                    {/if}
                {/if}
            </main>
        </div>

        <!-- Details Slide-Over Drawers -->
        <!-- Translucent Backdrop -->
        {#if selectedUserId || selectedGameId || selectedTemplate || selectedTicketId}
            <!-- svelte-ignore a11y_click_events_have_key_events -->
            <!-- svelte-ignore a11y_no_static_element_interactions -->
            <div 
                onclick={closeAllDrawers}
                transition:fade={{ duration: 150 }}
                class="fixed inset-0 bg-stone-900/50 backdrop-blur-xs z-[90] cursor-pointer"
            ></div>
        {/if}

        <!-- 1. USER DETAILS DRAWER -->
        {#if selectedUserId}
            <div 
                transition:fly={{ x: 450, duration: 300 }}
                class="fixed top-0 right-0 bottom-0 w-full sm:max-w-lg bg-stone-50 shadow-2xl border-l border-stone-200 z-[100] flex flex-col h-full overflow-hidden text-stone-800 font-sans"
            >
                <!-- Drawer Header -->
                <div class="bg-white p-5 border-b border-stone-200 flex items-center justify-between shrink-0">
                    <div class="space-y-1">
                        <div class="flex items-center gap-2">
                            <h3 class="font-bold text-lg text-stone-850 font-display">Détails de l'Utilisateur</h3>
                            {#if selectedUser?.is_banned}
                                <span class="px-2 py-0.5 rounded-full text-[10px] font-bold bg-red-100 text-red-700">Suspendu</span>
                            {:else}
                                <span class="px-2 py-0.5 rounded-full text-[10px] font-bold bg-emerald-100 text-emerald-700">Actif</span>
                            {/if}
                        </div>
                        <p class="text-xs text-stone-400 font-mono select-all">{selectedUserId}</p>
                    </div>
                    <button 
                        onclick={closeAllDrawers} 
                        class="p-1.5 rounded-xl border border-stone-200 hover:bg-stone-50 hover:text-burnt-orange transition-colors cursor-pointer text-stone-500"
                    >
                        <X size={18} />
                    </button>
                </div>

                <!-- Drawer Body -->
                <div class="flex-grow p-6 overflow-y-auto space-y-6">
                    {#if loadingUserDetails}
                        <div class="h-64 flex flex-col items-center justify-center gap-3">
                            <Loader2 class="animate-spin text-burnt-orange" size={32} />
                            <p class="text-xs text-stone-400 font-bold font-mono">Chargement des données...</p>
                        </div>
                    {:else if userDetails}
                        <!-- User Card -->
                        <div class="bg-white p-5 rounded-2xl border border-stone-200/80 shadow-xs flex items-center gap-4">
                            <div class="w-14 h-14 rounded-full bg-burnt-orange/10 border border-burnt-orange/20 text-burnt-orange font-bold text-xl flex items-center justify-center shrink-0 shadow-inner">
                                {selectedUser?.name?.slice(0, 2).toUpperCase() || "US"}
                            </div>
                            <div class="space-y-1 truncate">
                                <h4 class="font-bold text-stone-850 truncate">{selectedUser?.name}</h4>
                                <p class="text-xs text-stone-400 truncate flex items-center gap-1.5"><Mail size={12} /> {selectedUser?.email}</p>
                                <p class="text-xs text-stone-400 flex items-center gap-1.5"><Calendar size={12} /> Inscrit le {new Date(selectedUser?.created_at || '').toLocaleDateString("fr-FR")}</p>
                            </div>
                        </div>

                        <!-- Moderation Controls (Admins only) -->
                        {#if userRole === "admin"}
                            <div class="bg-white p-5 rounded-2xl border border-stone-200/80 shadow-xs space-y-4">
                                <h4 class="font-bold text-xs text-stone-400 uppercase tracking-widest flex items-center gap-2"><Shield size={14} /> Contrôle d'Administration</h4>
                                
                                <div class="grid grid-cols-1 sm:grid-cols-2 gap-4 pt-2">
                                    <div class="space-y-1.5">
                                        <label class="text-xs font-bold text-stone-500" for="role-select">Rôle du compte</label>
                                        <select
                                            id="role-select"
                                            value={selectedUser?.role || ""}
                                            onchange={(e) => handleRoleChange(selectedUserId || '', (e.target as HTMLSelectElement).value === "" ? null : (e.target as HTMLSelectElement).value as any)}
                                            disabled={actionLoading !== null}
                                            class="w-full bg-stone-50 border border-stone-200 rounded-xl px-3 py-2 text-sm focus:outline-none focus:border-burnt-orange shadow-xs disabled:opacity-50 font-medium cursor-pointer"
                                        >
                                            <option value="">Utilisateur Standard</option>
                                            <option value="moderator">Modérateur</option>
                                            <option value="admin">Administrateur</option>
                                        </select>
                                    </div>

                                    <div class="flex flex-col justify-end">
                                        <button
                                            onclick={() => handleBanToggle(selectedUserId || '', selectedUser?.is_banned || false)}
                                            disabled={actionLoading !== null}
                                            class="w-full py-2 px-4 rounded-xl border font-bold text-sm shadow-xs transition-all flex items-center justify-center gap-2 cursor-pointer
                                            {selectedUser?.is_banned 
                                                ? 'border-emerald-200 bg-emerald-50 text-emerald-700 hover:bg-emerald-100/70' 
                                                : 'border-red-200 bg-red-50 text-red-700 hover:bg-red-100/70'}"
                                        >
                                            {#if actionLoading === `ban-${selectedUserId}`}
                                                <Loader2 size={16} class="animate-spin" />
                                                <span>En cours...</span>
                                            {:else if selectedUser?.is_banned}
                                                <ShieldCheck size={16} />
                                                <span>Lever la suspension</span>
                                            {:else}
                                                <Ban size={16} />
                                                <span>Suspendre le compte</span>
                                            {/if}
                                        </button>
                                    </div>
                                </div>
                            </div>
                        {/if}

                        <!-- Stats Grid -->
                        <div class="grid grid-cols-3 gap-3">
                            <div class="bg-white p-3.5 rounded-xl border border-stone-200/80 shadow-xs text-center">
                                <p class="text-xs text-stone-400 font-bold uppercase tracking-wider">MJ</p>
                                <p class="text-xl font-bold text-burnt-orange mt-1">{userDetails.games_mastered}</p>
                            </div>
                            <div class="bg-white p-3.5 rounded-xl border border-stone-200/80 shadow-xs text-center">
                                <p class="text-xs text-stone-400 font-bold uppercase tracking-wider">Joueur</p>
                                <p class="text-xl font-bold text-blue-600 mt-1">{userDetails.games_played}</p>
                            </div>
                            <div class="bg-white p-3.5 rounded-xl border border-stone-200/80 shadow-xs text-center">
                                <p class="text-xs text-stone-400 font-bold uppercase tracking-wider">Persos</p>
                                <p class="text-xl font-bold text-stone-800 mt-1">{userDetails.characters_created}</p>
                            </div>
                        </div>

                        <!-- Games List -->
                        <div class="space-y-3">
                            <h4 class="font-bold text-xs text-stone-400 uppercase tracking-widest">Campagnes & Tables ({userDetails.games.length})</h4>
                            {#if userDetails.games.length === 0}
                                <p class="text-sm text-stone-400 italic bg-white p-5 rounded-2xl border border-stone-200/80 text-center">Cet utilisateur ne participe à aucune partie.</p>
                            {:else}
                                <div class="space-y-2">
                                    {#each userDetails.games as g (g.id)}
                                        <div 
                                            onclick={() => openGameDetails(g.id)}
                                            class="bg-white p-3 rounded-xl border border-stone-200/80 hover:border-burnt-orange shadow-xs flex items-center justify-between gap-3 cursor-pointer transition-colors"
                                        >
                                            <div class="flex items-center gap-3 min-w-0">
                                                {#if g.image_url}
                                                    <img src={g.image_url} class="w-10 h-10 rounded-lg object-cover border border-stone-150 shrink-0" alt="Bannière" />
                                                {:else}
                                                    <div class="w-10 h-10 rounded-lg bg-stone-100 border border-stone-200 text-stone-400 flex items-center justify-center shrink-0">
                                                        <BookOpen size={18} />
                                                    </div>
                                                {/if}
                                                <div class="min-w-0">
                                                    <p class="font-bold text-sm text-stone-850 truncate">{g.name}</p>
                                                    <p class="text-[10px] text-stone-400 font-mono mt-0.5 truncate">{g.id}</p>
                                                </div>
                                            </div>
                                            <div class="shrink-0 flex items-center gap-2">
                                                {#if g.is_gm}
                                                    <span class="px-2 py-0.5 rounded-md text-[9px] font-bold bg-burnt-orange/10 text-burnt-orange border border-burnt-orange/20">MJ</span>
                                                {:else}
                                                    <span class="px-2 py-0.5 rounded-md text-[9px] font-bold bg-blue-50 text-blue-600 border border-blue-100">Joueur</span>
                                                {/if}
                                                <ChevronRight size={14} class="text-stone-300" />
                                            </div>
                                        </div>
                                    {/each}
                                </div>
                            {/if}
                        </div>
                    {/if}
                </div>

                <!-- Drawer Footer -->
                <div class="bg-white p-4 border-t border-stone-200 shrink-0 flex justify-end">
                    <button 
                        onclick={closeAllDrawers}
                        class="px-4 py-2 bg-stone-100 text-stone-600 font-bold text-xs uppercase tracking-wider rounded-xl hover:bg-stone-200 transition-colors cursor-pointer"
                    >
                        Fermer
                    </button>
                </div>
            </div>
        {/if}

        <!-- 2. GAME DETAILS DRAWER -->
        {#if selectedGameId}
            <div 
                transition:fly={{ x: 450, duration: 300 }}
                class="fixed top-0 right-0 bottom-0 w-full sm:max-w-lg bg-stone-50 shadow-2xl border-l border-stone-200 z-[100] flex flex-col h-full overflow-hidden text-stone-800 font-sans"
            >
                <!-- Drawer Header -->
                <div class="bg-white p-5 border-b border-stone-200 flex items-center justify-between shrink-0">
                    <div class="space-y-1">
                        <div class="flex items-center gap-2">
                            <h3 class="font-bold text-lg text-stone-850 font-display">Détails de la Partie</h3>
                            {#if selectedGame?.is_active}
                                <span class="px-2 py-0.5 rounded-full text-[10px] font-bold bg-emerald-100 text-emerald-700">Actif</span>
                            {:else}
                                <span class="px-2 py-0.5 rounded-full text-[10px] font-bold bg-stone-100 text-stone-500">Archivée</span>
                            {/if}
                        </div>
                        <p class="text-xs text-stone-400 font-mono select-all">{selectedGameId}</p>
                    </div>
                    <button 
                        onclick={closeAllDrawers} 
                        class="p-1.5 rounded-xl border border-stone-200 hover:bg-stone-50 hover:text-burnt-orange transition-colors cursor-pointer text-stone-500"
                    >
                        <X size={18} />
                    </button>
                </div>

                <!-- Drawer Body -->
                <div class="flex-grow p-6 overflow-y-auto space-y-6">
                    {#if loadingGameDetails}
                        <div class="h-64 flex flex-col items-center justify-center gap-3">
                            <Loader2 class="animate-spin text-burnt-orange" size={32} />
                            <p class="text-xs text-stone-400 font-bold font-mono">Chargement des données...</p>
                        </div>
                    {:else if gameDetails && selectedGame}
                        <!-- Game Header Card -->
                        <div class="bg-white rounded-2xl border border-stone-200/80 shadow-xs overflow-hidden">
                            {#if selectedGame.image_url}
                                <img src={selectedGame.image_url} class="w-full h-32 object-cover border-b border-stone-150" alt="Bannière" />
                            {/if}
                            <div class="p-5 space-y-3">
                                <div>
                                    <h4 class="font-bold text-base text-stone-850 leading-tight">{selectedGame.name}</h4>
                                    <p class="text-xs text-stone-400 mt-1">Créée le {new Date(selectedGame.created_at).toLocaleDateString("fr-FR")}</p>
                                </div>
                                <div class="border-t border-stone-100 pt-3">
                                    <p class="text-xs text-stone-400 font-bold uppercase tracking-wider">Maître du Jeu (MJ)</p>
                                    <div class="flex items-center justify-between gap-3 mt-1.5">
                                        <div class="truncate">
                                            <p class="font-bold text-sm text-stone-800 truncate">{selectedGame.gm_name}</p>
                                            <p class="text-xs text-stone-400 truncate">{selectedGame.gm_email}</p>
                                        </div>
                                        <button 
                                            onclick={() => openUserDetails(selectedGame?.gm_id || '')} 
                                            class="p-1 px-2.5 rounded-lg border border-stone-200 text-stone-500 hover:text-burnt-orange hover:bg-stone-50 text-xs font-bold transition-colors cursor-pointer shrink-0"
                                        >
                                            Gérer MJ
                                        </button>
                                    </div>
                                </div>
                            </div>
                        </div>

                        <!-- Moderation Controls -->
                        <div class="bg-white p-5 rounded-2xl border border-stone-200/80 shadow-xs space-y-4">
                            <h4 class="font-bold text-xs text-stone-400 uppercase tracking-widest flex items-center gap-2"><Shield size={14} /> Outils de Modération</h4>
                            
                            <div class="grid grid-cols-2 gap-3">
                                <button
                                    onclick={() => handleGameArchive(selectedGameId || '', selectedGame?.is_active || false)}
                                    disabled={actionLoading !== null}
                                    class="py-2.5 px-4 rounded-xl border border-stone-200 bg-white text-stone-700 hover:bg-stone-50 font-bold text-sm shadow-xs transition-all flex items-center justify-center gap-2 cursor-pointer"
                                >
                                    {#if actionLoading === `archive-${selectedGameId}`}
                                        <Loader2 size={16} class="animate-spin" />
                                        <span>En cours...</span>
                                    {:else if selectedGame?.is_active}
                                        <Archive size={16} />
                                        <span>Archiver</span>
                                    {:else}
                                        <ArchiveRestore size={16} />
                                        <span>Désarchiver</span>
                                    {/if}
                                </button>

                                <button
                                    onclick={() => handleGameDelete(selectedGameId || '')}
                                    disabled={actionLoading !== null}
                                    class="py-2.5 px-4 rounded-xl border border-red-200 bg-red-50 text-red-700 hover:bg-red-100 font-bold text-sm shadow-xs transition-all flex items-center justify-center gap-2 cursor-pointer"
                                >
                                    {#if actionLoading === `delete-${selectedGameId}`}
                                        <Loader2 size={16} class="animate-spin" />
                                        <span>Suppression...</span>
                                    {:else}
                                        <Trash2 size={16} />
                                        <span>Supprimer</span>
                                    {/if}
                                </button>
                            </div>
                        </div>

                        <!-- Stats Grid -->
                        <div class="grid grid-cols-3 gap-3">
                            <div class="bg-white p-3.5 rounded-xl border border-stone-200/80 shadow-xs text-center">
                                <p class="text-xs text-stone-400 font-bold uppercase tracking-wider">Plateaux</p>
                                <p class="text-xl font-bold text-stone-800 mt-1">{gameDetails.boards_count}</p>
                            </div>
                            <div class="bg-white p-3.5 rounded-xl border border-stone-200/80 shadow-xs text-center">
                                <p class="text-xs text-stone-400 font-bold uppercase tracking-wider">Messages</p>
                                <p class="text-xl font-bold text-stone-800 mt-1">{gameDetails.messages_count}</p>
                            </div>
                            <div class="bg-white p-3.5 rounded-xl border border-stone-200/80 shadow-xs text-center">
                                <p class="text-xs text-stone-400 font-bold uppercase tracking-wider">Personnages</p>
                                <p class="text-xl font-bold text-stone-800 mt-1">{gameDetails.characters_count}</p>
                            </div>
                        </div>

                        <!-- Players List -->
                        <div class="space-y-3">
                            <h4 class="font-bold text-xs text-stone-400 uppercase tracking-widest">Joueurs et Personnages ({gameDetails.players.length})</h4>
                            {#if gameDetails.players.length === 0}
                                <p class="text-sm text-stone-400 italic bg-white p-5 rounded-2xl border border-stone-200/80 text-center">Aucun joueur n'a rejoint cette partie.</p>
                            {:else}
                                <div class="space-y-2.5">
                                    {#each gameDetails.players as p (p.user_id)}
                                        <div class="bg-white p-4 rounded-xl border border-stone-200/80 shadow-xs flex flex-col gap-3">
                                            <div class="flex items-center justify-between gap-3">
                                                <div class="flex items-center gap-2 min-w-0">
                                                    <!-- Ping color indicator -->
                                                    <span class="w-2.5 h-2.5 rounded-full shrink-0 border border-black/10" style="background-color: {p.ping_color || '#d97706'}"></span>
                                                    <div class="min-w-0">
                                                        <p class="font-bold text-sm text-stone-850 truncate">{p.name}</p>
                                                        <p class="text-xs text-stone-400 truncate">{p.email}</p>
                                                    </div>
                                                </div>
                                                <button 
                                                    onclick={() => openUserDetails(p.user_id)}
                                                    class="p-1 px-2.5 rounded-lg border border-stone-150 hover:bg-stone-50 text-[10px] font-bold text-stone-500 transition-colors shrink-0 cursor-pointer"
                                                >
                                                    Profil
                                                </button>
                                            </div>

                                            {#if p.character}
                                                <div class="bg-stone-50 p-2.5 rounded-lg border border-stone-150 flex items-center gap-3">
                                                    {#if p.character.avatar_url}
                                                        <img src={p.character.avatar_url} class="w-8 h-8 rounded-full border border-stone-250 object-cover shrink-0" alt="Avatar" />
                                                    {:else}
                                                        <div class="w-8 h-8 rounded-full bg-stone-200 border border-stone-300 text-stone-500 flex items-center justify-center shrink-0">
                                                            <Users size={14} />
                                                        </div>
                                                    {/if}
                                                    <div>
                                                        <p class="text-xs font-bold text-stone-800">{p.character.name}</p>
                                                        <p class="text-[10px] text-stone-400 uppercase font-bold tracking-wider">Personnage attribué</p>
                                                    </div>
                                                </div>
                                            {:else}
                                                <p class="text-[10px] text-stone-400 italic pl-4">Aucun personnage assigné</p>
                                            {/if}
                                        </div>
                                    {/each}
                                </div>
                            {/if}
                        </div>
                    {/if}
                </div>

                <!-- Drawer Footer -->
                <div class="bg-white p-4 border-t border-stone-200 shrink-0 flex justify-end">
                    <button 
                        onclick={closeAllDrawers}
                        class="px-4 py-2 bg-stone-100 text-stone-600 font-bold text-xs uppercase tracking-wider rounded-xl hover:bg-stone-200 transition-colors cursor-pointer"
                    >
                        Fermer
                    </button>
                </div>
            </div>
        {/if}

        <!-- 3. MARKETPLACE DETAILS DRAWER -->
        {#if selectedTemplate}
            <div 
                transition:fly={{ x: 450, duration: 300 }}
                class="fixed top-0 right-0 bottom-0 w-full sm:max-w-lg bg-stone-50 shadow-2xl border-l border-stone-200 z-[100] flex flex-col h-full overflow-hidden text-stone-800 font-sans"
            >
                <!-- Drawer Header -->
                <div class="bg-white p-5 border-b border-stone-200 flex items-center justify-between shrink-0">
                    <div class="space-y-1">
                        <div class="flex items-center gap-2">
                            <h3 class="font-bold text-lg text-stone-850 font-display truncate max-w-[200px]">{selectedTemplate.name}</h3>
                            <span class="px-2 py-0.5 rounded-full text-[9px] font-bold bg-stone-100 text-stone-600 uppercase tracking-wider shrink-0">
                                {selectedTemplate.type === 'BUNDLE' ? 'Pack' : (selectedTemplate.type === 'PERSONNAGE' ? 'Perso' : selectedTemplate.type)}
                            </span>
                            {#if !selectedTemplate.is_virtual}
                                {#if selectedTemplate.is_public}
                                    <span class="px-2 py-0.5 rounded-full text-[9px] font-bold bg-emerald-100 text-emerald-700 shrink-0">Public</span>
                                {:else}
                                    <span class="px-2 py-0.5 rounded-full text-[9px] font-bold bg-stone-100 text-stone-500 shrink-0">Privé</span>
                                {/if}
                            {/if}
                        </div>
                        <p class="text-xs text-stone-400 font-mono select-all truncate">{selectedTemplate.id}</p>
                    </div>
                    <button 
                        onclick={closeAllDrawers} 
                        class="p-1.5 rounded-xl border border-stone-200 hover:bg-stone-50 hover:text-burnt-orange transition-colors cursor-pointer text-stone-500"
                    >
                        <X size={18} />
                    </button>
                </div>

                <!-- Drawer Body -->
                <div class="flex-grow p-6 overflow-y-auto space-y-6">
                    {#if loadingTemplateDetails}
                        <div class="h-64 flex flex-col items-center justify-center gap-3">
                            <Loader2 class="animate-spin text-burnt-orange" size={32} />
                            <p class="text-xs text-stone-400 font-bold font-mono">Chargement des données...</p>
                        </div>
                    {:else if templateDetails}
                        <!-- Info Card -->
                        <div class="bg-white p-5 rounded-2xl border border-stone-200/80 shadow-xs space-y-3">
                            <div class="flex items-start justify-between gap-3">
                                <div>
                                    <h4 class="font-bold text-stone-850 leading-tight">{selectedTemplate.name}</h4>
                                    <p class="text-xs text-stone-400 mt-1">Créé par <span class="font-bold text-stone-600">{templateDetails.creator_name}</span></p>
                                </div>
                                <div class="text-right shrink-0">
                                    <p class="text-xs text-stone-400 font-bold uppercase tracking-wider">Téléchargements</p>
                                    <p class="text-lg font-bold text-stone-800 mt-0.5">{selectedTemplate.uses}</p>
                                </div>
                            </div>

                            {#if selectedTemplate.description}
                                <p class="text-sm text-stone-600 bg-stone-50 p-3 rounded-xl border border-stone-150 italic whitespace-pre-wrap">{selectedTemplate.description}</p>
                            {/if}
                        </div>

                        <!-- Moderation Actions -->
                        <div class="bg-white p-5 rounded-2xl border border-stone-200/80 shadow-xs space-y-3">
                            <h4 class="font-bold text-xs text-stone-400 uppercase tracking-widest flex items-center gap-2"><Shield size={14} /> Actions de modération</h4>
                            
                            {#if selectedTemplate.is_virtual}
                                <div class="bg-amber-50 p-4 rounded-xl border border-amber-200 text-amber-900 text-xs space-y-3">
                                    <p class="leading-relaxed flex gap-2 font-medium">
                                        <Info size={16} class="shrink-0 text-amber-600" />
                                        <span>Cet élément est virtuel et provient du pack <strong>{selectedTemplate.parent_bundle_name}</strong>. Vous devez modérer le pack parent pour modifier sa visibilité ou le retirer de la boutique.</span>
                                    </p>
                                    <button 
                                        onclick={() => openParentTemplateDetails(selectedTemplate.parent_bundle_id)}
                                        class="px-3.5 py-1.5 bg-amber-600 hover:bg-amber-700 text-white font-bold rounded-lg text-[10px] tracking-wide uppercase transition-colors flex items-center gap-1.5 cursor-pointer shadow-sm"
                                    >
                                        <BookOpen size={12} />
                                        <span>Gérer le pack parent</span>
                                    </button>
                                </div>
                            {:else}
                                <div class="grid grid-cols-2 gap-3">
                                    <button
                                        onclick={() => handleTemplatePublicToggle(selectedTemplate.id, selectedTemplate.is_public)}
                                        disabled={actionLoading !== null}
                                        class="py-2.5 px-4 rounded-xl border border-stone-200 bg-white text-stone-700 hover:bg-stone-50 font-bold text-sm shadow-xs transition-all flex items-center justify-center gap-2 cursor-pointer"
                                    >
                                        {#if actionLoading === `public-${selectedTemplate.id}`}
                                            <Loader2 size={16} class="animate-spin" />
                                            <span>En cours...</span>
                                        {:else if selectedTemplate.is_public}
                                            <Lock size={16} />
                                            <span>Rendre Privé</span>
                                        {:else}
                                            <Globe size={16} />
                                            <span>Publier</span>
                                        {/if}
                                    </button>

                                    <button
                                        onclick={() => handleTemplateDelete(selectedTemplate.id)}
                                        disabled={actionLoading !== null}
                                        class="py-2.5 px-4 rounded-xl border border-red-200 bg-red-50 text-red-700 hover:bg-red-100 font-bold text-sm shadow-xs transition-all flex items-center justify-center gap-2 cursor-pointer"
                                    >
                                        {#if actionLoading === `delete-${selectedTemplate.id}`}
                                            <Loader2 size={16} class="animate-spin" />
                                            <span>Suppression...</span>
                                        {:else}
                                            <Trash2 size={16} />
                                            <span>Supprimer</span>
                                        {/if}
                                    </button>
                                </div>
                            {/if}
                        </div>

                        <!-- Template Content -->
                        {#if selectedTemplate.type === "BUNDLE"}
                            <!-- Bundle list of items -->
                            <div class="space-y-3">
                                <h4 class="font-bold text-xs text-stone-400 uppercase tracking-widest">Contenu du Pack ({templateDetails.data?.items?.length || 0} éléments)</h4>
                                
                                <div class="space-y-2">
                                    {#each (templateDetails.data?.items || []) as item, idx}
                                        <div class="bg-white p-3.5 rounded-xl border border-stone-200/80 shadow-xs flex gap-3 items-center">
                                            {#if item.data?.avatar_url}
                                                <img src={item.data.avatar_url} class="w-10 h-10 rounded-lg object-cover border border-stone-250 shrink-0" alt="Avatar" />
                                            {:else}
                                                <div class="w-10 h-10 rounded-lg bg-stone-100 border border-stone-200 text-stone-400 flex items-center justify-center shrink-0">
                                                    <Users size={16} />
                                                </div>
                                            {/if}
                                            <div class="min-w-0 flex-1">
                                                <div class="flex items-center gap-2">
                                                    <p class="font-bold text-sm text-stone-850 truncate">{item.name}</p>
                                                    <span class="px-1.5 py-0.5 rounded text-[8px] font-bold bg-stone-100 text-stone-500 uppercase tracking-wider shrink-0">
                                                        {item.type === 'MONSTER' ? 'MONSTRE' : (item.type === 'PLAYER' ? 'PERSO' : item.type)}
                                                    </span>
                                                </div>
                                                <p class="text-xs text-stone-400 truncate mt-0.5">{item.data?.description || "Aucune description."}</p>
                                            </div>
                                            <button 
                                                onclick={() => {
                                                    // Load this virtual item in the drawer!
                                                    const virtualItem = {
                                                        id: `${selectedTemplate.id}-item-${idx}`,
                                                        parent_bundle_id: selectedTemplate.id,
                                                        parent_bundle_name: selectedTemplate.name,
                                                        created_by: selectedTemplate.created_by,
                                                        creator_name: templateDetails?.creator_name || "Créateur",
                                                        name: item.name,
                                                        description: item.data?.description || `Fait partie du pack "${selectedTemplate.name}"`,
                                                        type: item.type === 'MONSTER' ? 'MONSTRE' : (item.type === 'PLAYER' ? 'PERSONNAGE' : item.type),
                                                        data: item.data,
                                                        is_public: selectedTemplate.is_public,
                                                        uses: selectedTemplate.uses,
                                                        created_at: selectedTemplate.created_at,
                                                        is_virtual: true
                                                    };
                                                    openTemplateDetails(virtualItem);
                                                }}
                                                class="p-1 rounded-lg border border-stone-150 text-stone-400 hover:text-burnt-orange hover:bg-stone-50 transition-colors cursor-pointer shrink-0"
                                                title="Inspecter les stats"
                                            >
                                                <ChevronRight size={16} />
                                            </button>
                                        </div>
                                    {/each}
                                </div>
                            </div>
                        {:else}
                            <!-- Character / Monster Stats rendering -->
                            <div class="space-y-6">
                                <!-- Potential Avatar image -->
                                {#if templateDetails.data?.avatar_url}
                                    <div class="flex justify-center py-2 bg-white rounded-2xl border border-stone-200/80 p-4 shadow-xs">
                                        <img src={templateDetails.data.avatar_url} class="w-32 h-32 rounded-full object-cover border-4 border-stone-100 shadow-md" alt="Avatar" />
                                    </div>
                                {/if}

                                <!-- Main Attributes block -->
                                <div class="bg-white p-4 rounded-2xl border border-stone-200/80 shadow-xs grid grid-cols-4 gap-2 text-center">
                                    <div>
                                        <p class="text-[10px] text-stone-400 font-bold uppercase tracking-wider">PV Max</p>
                                        <p class="text-base font-extrabold text-red-600 mt-0.5">{templateDetails.data?.max_hp || 10}</p>
                                    </div>
                                    <div>
                                        <p class="text-[10px] text-stone-400 font-bold uppercase tracking-wider">CA</p>
                                        <p class="text-base font-extrabold text-blue-600 mt-0.5">{templateDetails.data?.armor_class || 10}</p>
                                    </div>
                                    <div>
                                        <p class="text-[10px] text-stone-400 font-bold uppercase tracking-wider">Vitesse</p>
                                        <p class="text-base font-extrabold text-emerald-600 mt-0.5">{templateDetails.data?.speed || 9}m</p>
                                    </div>
                                    <div>
                                        <p class="text-[10px] text-stone-400 font-bold uppercase tracking-wider">Initiative</p>
                                        <p class="text-base font-extrabold text-stone-800 mt-0.5">+{templateDetails.data?.initiative || 0}</p>
                                    </div>
                                </div>

                                <!-- Core Stats Grid -->
                                <div class="bg-white p-5 rounded-2xl border border-stone-200/80 shadow-xs space-y-3">
                                    <h4 class="font-bold text-xs text-stone-400 uppercase tracking-widest flex items-center gap-1.5"><Award size={14} /> Caractéristiques</h4>
                                    
                                    <div class="grid grid-cols-6 gap-2">
                                        <!-- STR -->
                                        <div class="bg-stone-50 p-2 rounded-lg text-center border border-stone-150">
                                            <p class="text-[9px] font-bold text-stone-400 uppercase">FOR</p>
                                            <p class="text-sm font-bold text-stone-800 mt-0.5">{templateDetails.data?.strength || 10}</p>
                                            <p class="text-[9px] text-stone-400 mt-0.5 font-bold">
                                                {templateDetails.data?.strength_mod >= 0 ? '+' : ''}{templateDetails.data?.strength_mod || 0}
                                            </p>
                                        </div>
                                        <!-- DEX -->
                                        <div class="bg-stone-50 p-2 rounded-lg text-center border border-stone-150">
                                            <p class="text-[9px] font-bold text-stone-400 uppercase">DEX</p>
                                            <p class="text-sm font-bold text-stone-800 mt-0.5">{templateDetails.data?.dexterity || 10}</p>
                                            <p class="text-[9px] text-stone-400 mt-0.5 font-bold">
                                                {templateDetails.data?.dexterity_mod >= 0 ? '+' : ''}{templateDetails.data?.dexterity_mod || 0}
                                            </p>
                                        </div>
                                        <!-- CON -->
                                        <div class="bg-stone-50 p-2 rounded-lg text-center border border-stone-150">
                                            <p class="text-[9px] font-bold text-stone-400 uppercase">CON</p>
                                            <p class="text-sm font-bold text-stone-800 mt-0.5">{templateDetails.data?.constitution || 10}</p>
                                            <p class="text-[9px] text-stone-400 mt-0.5 font-bold">
                                                {templateDetails.data?.constitution_mod >= 0 ? '+' : ''}{templateDetails.data?.constitution_mod || 0}
                                            </p>
                                        </div>
                                        <!-- INT -->
                                        <div class="bg-stone-50 p-2 rounded-lg text-center border border-stone-150">
                                            <p class="text-[9px] font-bold text-stone-400 uppercase">INT</p>
                                            <p class="text-sm font-bold text-stone-800 mt-0.5">{templateDetails.data?.intelligence || 10}</p>
                                            <p class="text-[9px] text-stone-400 mt-0.5 font-bold">
                                                {templateDetails.data?.intelligence_mod >= 0 ? '+' : ''}{templateDetails.data?.intelligence_mod || 0}
                                            </p>
                                        </div>
                                        <!-- WIS -->
                                        <div class="bg-stone-50 p-2 rounded-lg text-center border border-stone-150">
                                            <p class="text-[9px] font-bold text-stone-400 uppercase">SAG</p>
                                            <p class="text-sm font-bold text-stone-800 mt-0.5">{templateDetails.data?.wisdom || 10}</p>
                                            <p class="text-[9px] text-stone-400 mt-0.5 font-bold">
                                                {templateDetails.data?.wisdom_mod >= 0 ? '+' : ''}{templateDetails.data?.wisdom_mod || 0}
                                            </p>
                                        </div>
                                        <!-- CHA -->
                                        <div class="bg-stone-50 p-2 rounded-lg text-center border border-stone-150">
                                            <p class="text-[9px] font-bold text-stone-400 uppercase">CHA</p>
                                            <p class="text-sm font-bold text-stone-800 mt-0.5">{templateDetails.data?.charisma || 10}</p>
                                            <p class="text-[9px] text-stone-400 mt-0.5 font-bold">
                                                {templateDetails.data?.charisma_mod >= 0 ? '+' : ''}{templateDetails.data?.charisma_mod || 0}
                                            </p>
                                        </div>
                                    </div>
                                </div>

                                <!-- Abilities block -->
                                {#if templateDetails.data?.abilities}
                                    <div class="bg-white p-5 rounded-2xl border border-stone-200/80 shadow-xs space-y-2.5">
                                        <h4 class="font-bold text-xs text-stone-400 uppercase tracking-widest flex items-center gap-1.5"><Sparkles size={14} /> Capacités & Traits</h4>
                                        <p class="text-xs text-stone-600 leading-relaxed whitespace-pre-wrap">{templateDetails.data.abilities}</p>
                                    </div>
                                {/if}

                                <!-- Attacks / Spells block -->
                                {#if templateDetails.data?.spells && Object.keys(templateDetails.data.spells).length > 0}
                                    <div class="bg-white p-5 rounded-2xl border border-stone-200/80 shadow-xs space-y-4">
                                        <h4 class="font-bold text-xs text-stone-400 uppercase tracking-widest flex items-center gap-1.5"><Sword size={14} /> Actions & Sorts</h4>
                                        
                                        <div class="space-y-3">
                                            {#each Object.keys(templateDetails.data.spells) as category}
                                                <div class="space-y-2">
                                                    <h5 class="text-xs font-extrabold text-burnt-orange border-b border-stone-100 pb-1">{category}</h5>
                                                    <div class="space-y-2">
                                                        {#each (templateDetails.data.spells[category] || []) as spell}
                                                            <div class="bg-stone-50 p-2.5 rounded-lg border border-stone-150 text-xs">
                                                                <div class="flex items-center justify-between gap-2">
                                                                    <p class="font-bold text-stone-850">{spell.name}</p>
                                                                    {#if spell.charges}
                                                                        <span class="text-[9px] font-bold px-1.5 py-0.5 rounded bg-amber-100 text-amber-800">{spell.charges}</span>
                                                                    {/if}
                                                                </div>
                                                                <p class="text-stone-500 mt-1 leading-relaxed">{spell.description}</p>
                                                            </div>
                                                        {/each}
                                                    </div>
                                                </div>
                                            {/each}
                                        </div>
                                    </div>
                                {/if}
                            </div>
                        {/if}
                    {/if}
                </div>

                <!-- Drawer Footer -->
                <div class="bg-white p-4 border-t border-stone-200 shrink-0 flex justify-between items-center">
                    {#if selectedTemplate.is_virtual}
                        <button 
                            onclick={() => openParentTemplateDetails(selectedTemplate.parent_bundle_id)}
                            class="px-3 py-2 border border-stone-200 hover:bg-stone-50 text-stone-600 font-bold text-xs uppercase tracking-wider rounded-xl transition-colors cursor-pointer flex items-center gap-1.5"
                        >
                            <BookOpen size={14} />
                            <span>Voir le pack</span>
                        </button>
                    {:else}
                        <div></div>
                    {/if}
                    <button 
                        onclick={closeAllDrawers}
                        class="px-4 py-2 bg-stone-100 text-stone-600 font-bold text-xs uppercase tracking-wider rounded-xl hover:bg-stone-200 transition-colors cursor-pointer"
                    >
                        Fermer
                    </button>
                </div>
            </div>
        {/if}

        <!-- 4. SUPPORT TICKET DETAILS DRAWER -->
        {#if selectedTicketId && selectedTicket}
            <div 
                transition:fly={{ x: 450, duration: 300 }}
                class="fixed top-0 right-0 bottom-0 w-full sm:max-w-lg bg-stone-50 shadow-2xl border-l border-stone-200 z-[100] flex flex-col h-full overflow-hidden text-stone-800 font-sans"
            >
                <!-- Drawer Header -->
                <div class="bg-white p-5 border-b border-stone-200 flex items-center justify-between shrink-0">
                    <div class="space-y-1">
                        <div class="flex items-center gap-2">
                            <h3 class="font-bold text-lg text-stone-850 font-display">Détails du Ticket</h3>
                            <span class="px-2 py-0.5 rounded-full text-[9px] font-bold uppercase tracking-wider
                                {selectedTicket.category === 'BUG' ? 'bg-red-150 text-red-700' : ''}
                                {selectedTicket.category === 'CONTACT' ? 'bg-blue-150 text-blue-700' : ''}
                                {selectedTicket.category === 'RECLAMATION' ? 'bg-amber-150 text-amber-700' : ''}"
                            >
                                {selectedTicket.category === 'BUG' ? 'Bug' : selectedTicket.category === 'CONTACT' ? 'Contact' : 'Réclamation'}
                            </span>
                        </div>
                        <p class="text-xs text-stone-400 font-mono select-all truncate">{selectedTicketId}</p>
                    </div>
                    <button 
                        onclick={closeAllDrawers} 
                        class="p-1.5 rounded-xl border border-stone-200 hover:bg-stone-50 hover:text-burnt-orange transition-colors cursor-pointer text-stone-500"
                    >
                        <X size={18} />
                    </button>
                </div>

                <!-- Drawer Body -->
                <div class="flex-grow p-6 overflow-y-auto space-y-6">
                    <!-- Ticket Info Card -->
                    <div class="bg-white p-5 rounded-2xl border border-stone-200/80 shadow-xs space-y-3">
                        <div>
                            <p class="text-[10px] text-stone-400 font-bold uppercase tracking-widest">Sujet du message</p>
                            <h4 class="font-bold text-base text-stone-850 mt-1 leading-tight">{selectedTicket.subject}</h4>
                        </div>

                        <div class="border-t border-stone-100 pt-3 grid grid-cols-2 gap-3 text-xs">
                            <div>
                                <p class="text-[10px] text-stone-400 font-bold uppercase tracking-wider">Soumis par</p>
                                <p class="font-bold text-stone-800 mt-1">{selectedTicket.user_name || "Utilisateur anonyme"}</p>
                                <p class="text-stone-500 font-medium">{selectedTicket.email}</p>
                            </div>
                            <div>
                                <p class="text-[10px] text-stone-400 font-bold uppercase tracking-wider">Date d'envoi</p>
                                <p class="font-semibold text-stone-700 mt-1">{new Date(selectedTicket.created_at).toLocaleDateString("fr-FR")}</p>
                                <p class="text-stone-400 font-medium">à {new Date(selectedTicket.created_at).toLocaleTimeString("fr-FR", { hour: '2-digit', minute: '2-digit' })}</p>
                            </div>
                        </div>
                    </div>

                    <!-- Status Change Dropdown -->
                    <div class="bg-white p-5 rounded-2xl border border-stone-200/80 shadow-xs space-y-2">
                        <label class="text-xs font-bold text-stone-400 uppercase tracking-widest" for="ticket-status-select">Mettre à jour le statut</label>
                        <select
                            id="ticket-status-select"
                            value={selectedTicket.status}
                            onchange={(e) => handleTicketStatusChange(selectedTicket?.id || '', (e.target as HTMLSelectElement).value)}
                            disabled={actionLoading === `ticket-${selectedTicket.id}`}
                            class="w-full bg-stone-50 border border-stone-200 rounded-xl px-3 py-2 text-sm font-semibold focus:outline-none focus:border-burnt-orange shadow-xs disabled:opacity-50 cursor-pointer"
                        >
                            <option value="PENDING">En attente (Non traité)</option>
                            <option value="IN_PROGRESS">En cours de traitement</option>
                            <option value="RESOLVED">Résolu (Fermé avec succès)</option>
                            <option value="CLOSED">Fermé sans résolution</option>
                        </select>
                    </div>

                    <!-- Ticket Message -->
                    <div class="space-y-2">
                        <h4 class="font-bold text-xs text-stone-400 uppercase tracking-widest flex items-center gap-1.5"><FileText size={14} /> Message reçu</h4>
                        <div class="bg-white p-5 rounded-2xl border border-stone-200/80 shadow-xs text-sm text-stone-700 whitespace-pre-wrap leading-relaxed min-h-[120px]">
                            {selectedTicket.message}
                        </div>
                    </div>

                    <!-- Attachment Rendering -->
                    {#if selectedTicket.attachment_url}
                        <div class="space-y-2">
                            <h4 class="font-bold text-xs text-stone-400 uppercase tracking-widest">Pièce jointe</h4>
                            
                            {#if isImage(selectedTicket.attachment_url)}
                                <!-- Render Image Directly -->
                                <div class="bg-white p-4 rounded-2xl border border-stone-200/80 shadow-xs overflow-hidden flex flex-col gap-2">
                                    <div class="group relative rounded-xl border border-stone-200 overflow-hidden bg-stone-100 flex items-center justify-center">
                                        <img 
                                            src={selectedTicket.attachment_url} 
                                            class="max-w-full max-h-[320px] object-contain transition-transform duration-300 group-hover:scale-[1.02]" 
                                            alt="Pièce jointe du ticket" 
                                        />
                                        <a 
                                            href={selectedTicket.attachment_url} 
                                            target="_blank" 
                                            class="absolute inset-0 bg-black/40 opacity-0 group-hover:opacity-100 flex items-center justify-center text-white text-xs font-bold gap-1.5 transition-opacity duration-200 cursor-pointer"
                                        >
                                            <ExternalLink size={14} />
                                            <span>Ouvrir en taille réelle</span>
                                        </a>
                                    </div>
                                    <a 
                                        href={selectedTicket.attachment_url} 
                                        download 
                                        class="text-xs text-burnt-orange hover:underline font-bold flex items-center gap-1.5 mt-1 self-start"
                                    >
                                        <Download size={14} />
                                        <span>Télécharger l'image</span>
                                    </a>
                                </div>
                            {:else}
                                <!-- Render Download Link Card -->
                                <a 
                                    href={selectedTicket.attachment_url} 
                                    target="_blank"
                                    download
                                    class="bg-white p-4 rounded-xl border border-stone-200/80 hover:border-burnt-orange shadow-xs flex items-center justify-between gap-3 transition-colors cursor-pointer"
                                >
                                    <div class="flex items-center gap-3 min-w-0">
                                        <div class="w-10 h-10 rounded-lg bg-stone-100 border border-stone-200 text-stone-500 flex items-center justify-center shrink-0">
                                            <FileText size={20} />
                                        </div>
                                        <div class="min-w-0">
                                            <p class="font-bold text-sm text-stone-850 truncate">Fichier joint</p>
                                            <p class="text-xs text-stone-400 truncate">Cliquez pour télécharger le document</p>
                                        </div>
                                    </div>
                                    <Download size={18} class="text-stone-400 shrink-0" />
                                </a>
                            {/if}
                        </div>
                    {/if}
                </div>

                <!-- Drawer Footer -->
                <div class="bg-white p-4 border-t border-stone-200 shrink-0 flex justify-end">
                    <button 
                        onclick={closeAllDrawers}
                        class="px-4 py-2 bg-stone-100 text-stone-600 font-bold text-xs uppercase tracking-wider rounded-xl hover:bg-stone-200 transition-colors cursor-pointer"
                    >
                        Fermer
                    </button>
                </div>
            </div>
        {/if}
    {/if}
</div>
