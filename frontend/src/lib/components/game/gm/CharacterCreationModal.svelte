<script lang="ts">
    import { X, Upload, Link, Plus, Trash2, Sticker } from "lucide-svelte";
    import { api } from "$lib/api";
    import { authClient } from "$lib/auth-client";
    import { page } from "$app/state";

    let {
        isOpen,
        onClose,
        onCharacterCreated,
        character = null,
    } = $props<{
        isOpen: boolean;
        onClose: () => void;
        onCharacterCreated: () => void;
        character?: any;
    }>();

    let name = $state("");
    let race = $state("");
    let maxHP = $state(10);
    let money = $state(0);
    let isNPC = $state(false);
    let avatarType = $state<"upload" | "url">("upload");
    let avatarFile = $state<File | null>(null);
    let avatarURL = $state("");
    let stats = $state<{ key: string; value: string }[]>([]);
    let inventory = $state<
        {
            name: string;
            quantity: string;
            imageType: "upload" | "url" | "icon";
            imageFile: File | null;
            imageURL: string;
            iconName: string;
        }[]
    >([]);
    let loading = $state(false);

    import { untrack } from "svelte";

    $effect(() => {
        // Track dependencies
        const open = isOpen;
        const char = character;

        untrack(() => {
            if (open && char) {
                console.log("Editing character:", char);
                name = char.name || "";
                race = char.race || "";
                maxHP = char.max_hp || 10;
                money = char.money || 0;
                isNPC = char.is_npc || false;

                if (char.avatar_url) {
                    avatarType = "url";
                    avatarURL = char.avatar_url;
                } else {
                    avatarType = "upload";
                    avatarURL = "";
                }

                // Stats
                stats = [];
                if (char.stats && typeof char.stats === "object") {
                    try {
                        Object.entries(char.stats).forEach(([key, value]) => {
                            stats.push({ key, value: String(value) });
                        });
                    } catch (e) {
                        console.error("Error parsing stats:", e);
                    }
                }

                // Inventory
                inventory = [];
                if (char.inventory && Array.isArray(char.inventory)) {
                    try {
                        char.inventory.forEach((item: any) => {
                            inventory.push({
                                name: item.name,
                                quantity: item.quantity,
                                imageType: item.icon_name
                                    ? "icon"
                                    : item.image_url
                                      ? "url"
                                      : "upload",
                                imageFile: null,
                                imageURL: item.image_url || "",
                                iconName: item.icon_name || "",
                            });
                        });
                    } catch (e) {
                        console.error("Error parsing inventory:", e);
                    }
                } else if (char.inventory) {
                    console.warn("Inventory is not an array:", char.inventory);
                }
            } else if (open && !char) {
                // Reset form if opening for create
                name = "";
                race = "";
                maxHP = 10;
                money = 0;
                isNPC = false;
                avatarFile = null;
                avatarURL = "";
                stats = [];
                inventory = [];
                avatarType = "upload";
            }
        });
    });

    function addStat() {
        stats = [...stats, { key: "", value: "" }];
    }

    function removeStat(index: number) {
        stats = stats.filter((_, i) => i !== index);
    }

    function addInventoryItem() {
        inventory = [
            ...inventory,
            {
                name: "",
                quantity: "",
                imageType: "upload",
                imageFile: null,
                imageURL: "",
                iconName: "",
            },
        ];
    }

    function removeInventoryItem(index: number) {
        inventory = inventory.filter((_, i) => i !== index);
    }

    function handleFileChange(event: Event) {
        const input = event.target as HTMLInputElement;
        if (input.files && input.files.length > 0) {
            avatarFile = input.files[0];
        }
    }

    async function handleSubmit() {
        if (!name || !race || !maxHP) return;

        loading = true;
        const gameId = page.params.id;

        try {
            const { data: tokenData } = await authClient.token();
            if (tokenData?.token) {
                const formData = new FormData();
                formData.append("name", name);
                formData.append("race", race);
                formData.append("max_hp", maxHP.toString());
                formData.append("money", money.toString());
                formData.append("is_npc", isNPC.toString());

                if (avatarType === "upload" && avatarFile) {
                    formData.append("avatar", avatarFile);
                } else if (avatarType === "url" && avatarURL) {
                    formData.append("avatar_url", avatarURL);
                }

                // Convert stats array to object
                const statsObj: Record<string, string> = {};
                stats.forEach((stat) => {
                    if (stat.key) statsObj[stat.key] = stat.value;
                });
                formData.append("stats", JSON.stringify(statsObj));

                // Prepare inventory items
                const inventoryItems = inventory.map((item) => ({
                    name: item.name,
                    quantity: item.quantity,
                    image_url:
                        item.imageType === "url" ? item.imageURL : undefined,
                    icon_name:
                        item.imageType === "icon" ? item.iconName : undefined,
                }));
                formData.append("inventory", JSON.stringify(inventoryItems));

                // Append inventory images
                inventory.forEach((item, index) => {
                    if (item.imageType === "upload" && item.imageFile) {
                        formData.append(
                            `inventory_image_${index}`,
                            item.imageFile,
                        );
                    }
                });

                if (character && character.id) {
                    await api.put(
                        `/table/${gameId}/characters/${character.id}`,
                        formData,
                        {
                            headers: {
                                Authorization: `Bearer ${tokenData.token}`,
                                "Content-Type": "multipart/form-data",
                            },
                        },
                    );
                } else {
                    await api.post(`/table/${gameId}/characters`, formData, {
                        headers: {
                            Authorization: `Bearer ${tokenData.token}`,
                            "Content-Type": "multipart/form-data",
                        },
                    });
                }

                onCharacterCreated();
                onClose();
            }
        } catch (error) {
            console.error("Failed to save character:", error);
            alert("Erreur lors de la sauvegarde du personnage.");
        } finally {
            loading = false;
        }
    }
</script>

{#if isOpen}
    <div
        class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/50 backdrop-blur-sm"
    >
        <div
            class="bg-white rounded-2xl shadow-xl w-full max-w-lg max-h-[90vh] overflow-y-auto animate-in fade-in zoom-in duration-200"
        >
            <div
                class="p-6 border-b border-stone-100 flex justify-between items-center sticky top-0 bg-white z-10"
            >
                <h2 class="text-xl font-bold text-dark-gray">
                    {character?.id
                        ? "Modifier le personnage"
                        : "Créer un personnage"}
                </h2>
                <button
                    onclick={onClose}
                    class="p-2 hover:bg-stone-100 rounded-full transition-colors"
                >
                    <X size={20} />
                </button>
            </div>

            <div class="p-6 space-y-6">
                <!-- Name & Race & HP -->
                <div class="grid grid-cols-3 gap-4">
                    <div class="space-y-2">
                        <label
                            for="name"
                            class="text-sm font-bold text-dark-gray">Nom</label
                        >
                        <input
                            type="text"
                            id="name"
                            bind:value={name}
                            class="w-full px-4 py-2 rounded-xl border border-stone-200 focus:outline-none focus:ring-2 focus:ring-burnt-orange/20 focus:border-burnt-orange transition-all"
                            placeholder="Ex: Gandalf"
                        />
                    </div>
                    <div class="space-y-2">
                        <label
                            for="race"
                            class="text-sm font-bold text-dark-gray">Race</label
                        >
                        <input
                            type="text"
                            id="race"
                            bind:value={race}
                            class="w-full px-4 py-2 rounded-xl border border-stone-200 focus:outline-none focus:ring-2 focus:ring-burnt-orange/20 focus:border-burnt-orange transition-all"
                            placeholder="Ex: Humain"
                        />
                    </div>
                    <div class="space-y-2">
                        <label for="hp" class="text-sm font-bold text-dark-gray"
                            >PV Max</label
                        >
                        <input
                            type="number"
                            id="hp"
                            bind:value={maxHP}
                            min="1"
                            class="w-full px-4 py-2 rounded-xl border border-stone-200 focus:outline-none focus:ring-2 focus:ring-burnt-orange/20 focus:border-burnt-orange transition-all"
                        />
                    </div>
                    <div class="space-y-2">
                        <label
                            for="money"
                            class="text-sm font-bold text-dark-gray"
                            >Argent</label
                        >
                        <input
                            type="number"
                            id="money"
                            bind:value={money}
                            min="0"
                            class="w-full px-4 py-2 rounded-xl border border-stone-200 focus:outline-none focus:ring-2 focus:ring-burnt-orange/20 focus:border-burnt-orange transition-all"
                        />
                    </div>
                </div>

                <!-- Type (PC/NPC) -->
                <div class="space-y-2">
                    <label class="text-sm font-bold text-dark-gray block"
                        >Type</label
                    >
                    <div class="flex bg-stone-100 p-1 rounded-xl">
                        <button
                            class="flex-1 py-2 rounded-lg text-sm font-medium transition-all {isNPC
                                ? 'text-stone-500 hover:text-dark-gray'
                                : 'bg-white text-dark-gray shadow-sm'}"
                            onclick={() => (isNPC = false)}
                        >
                            Joueur (PJ)
                        </button>
                        <button
                            class="flex-1 py-2 rounded-lg text-sm font-medium transition-all {!isNPC
                                ? 'text-stone-500 hover:text-dark-gray'
                                : 'bg-white text-dark-gray shadow-sm'}"
                            onclick={() => (isNPC = true)}
                        >
                            Non-Joueur (PNJ)
                        </button>
                    </div>
                </div>

                <!-- Avatar -->
                <div class="space-y-2">
                    <label class="text-sm font-bold text-dark-gray block"
                        >Avatar</label
                    >
                    <div class="flex gap-4 mb-2">
                        <button
                            class="flex items-center gap-2 text-sm font-medium {avatarType ===
                            'upload'
                                ? 'text-burnt-orange'
                                : 'text-stone-400'}"
                            onclick={() => (avatarType = "upload")}
                        >
                            <Upload size={16} /> Upload
                        </button>
                        <button
                            class="flex items-center gap-2 text-sm font-medium {avatarType ===
                            'url'
                                ? 'text-burnt-orange'
                                : 'text-stone-400'}"
                            onclick={() => (avatarType = "url")}
                        >
                            <Link size={16} /> URL
                        </button>
                    </div>

                    {#if avatarType === "upload"}
                        <div
                            class="border-2 border-dashed border-stone-200 rounded-xl p-8 text-center hover:border-burnt-orange/50 transition-colors cursor-pointer relative"
                        >
                            <input
                                type="file"
                                accept="image/*"
                                onchange={handleFileChange}
                                class="absolute inset-0 w-full h-full opacity-0 cursor-pointer"
                            />
                            {#if avatarFile}
                                <p class="text-sm text-dark-gray font-medium">
                                    {avatarFile.name}
                                </p>
                            {:else}
                                <p class="text-sm text-stone-500">
                                    Cliquez ou glissez une image ici
                                </p>
                            {/if}
                        </div>
                    {:else}
                        <input
                            type="url"
                            bind:value={avatarURL}
                            class="w-full px-4 py-2 rounded-xl border border-stone-200 focus:outline-none focus:ring-2 focus:ring-burnt-orange/20 focus:border-burnt-orange transition-all"
                            placeholder="https://example.com/avatar.png"
                        />
                    {/if}
                </div>

                <!-- Stats -->
                <div class="space-y-3">
                    <div class="flex justify-between items-center">
                        <label class="text-sm font-bold text-dark-gray"
                            >Statistiques</label
                        >
                        <button
                            onclick={addStat}
                            class="text-xs flex items-center gap-1 text-burnt-orange font-medium hover:text-burnt-orange/80"
                        >
                            <Plus size={14} /> Ajouter
                        </button>
                    </div>

                    {#if stats.length === 0}
                        <p class="text-xs text-stone-400 italic">
                            Aucune statistique définie.
                        </p>
                    {/if}

                    <div class="space-y-2">
                        {#each stats as stat, i}
                            <div class="flex gap-2">
                                <input
                                    type="text"
                                    bind:value={stat.key}
                                    placeholder="Nom (ex: Force)"
                                    class="flex-1 px-3 py-2 rounded-lg border border-stone-200 text-sm focus:outline-none focus:border-burnt-orange"
                                />
                                <input
                                    type="text"
                                    bind:value={stat.value}
                                    placeholder="Valeur (ex: 15)"
                                    class="flex-1 px-3 py-2 rounded-lg border border-stone-200 text-sm focus:outline-none focus:border-burnt-orange"
                                />
                                <button
                                    onclick={() => removeStat(i)}
                                    class="p-2 text-stone-400 hover:text-red-500 hover:bg-red-50 rounded-lg transition-all"
                                >
                                    <Trash2 size={16} />
                                </button>
                            </div>
                        {/each}
                    </div>
                </div>

                <!-- Inventory -->
                <div class="space-y-3">
                    <div class="flex justify-between items-center">
                        <label class="text-sm font-bold text-dark-gray"
                            >Inventaire</label
                        >
                        <button
                            onclick={addInventoryItem}
                            class="text-xs flex items-center gap-1 text-burnt-orange font-medium hover:text-burnt-orange/80"
                        >
                            <Plus size={14} /> Ajouter
                        </button>
                    </div>

                    {#if inventory.length === 0}
                        <p class="text-xs text-stone-400 italic">
                            Inventaire vide.
                        </p>
                    {/if}

                    <div class="space-y-4">
                        {#each inventory as item, i}
                            <div
                                class="p-3 bg-stone-50 rounded-xl border border-stone-100 space-y-3"
                            >
                                <div class="flex gap-2">
                                    <input
                                        type="text"
                                        bind:value={item.name}
                                        placeholder="Objet (ex: Épée)"
                                        class="flex-1 px-3 py-2 rounded-lg border border-stone-200 text-sm focus:outline-none focus:border-burnt-orange"
                                    />
                                    <input
                                        type="text"
                                        bind:value={item.quantity}
                                        placeholder="Qté (ex: 1)"
                                        class="w-24 px-3 py-2 rounded-lg border border-stone-200 text-sm focus:outline-none focus:border-burnt-orange"
                                    />
                                    <button
                                        onclick={() => removeInventoryItem(i)}
                                        class="p-2 text-stone-400 hover:text-red-500 hover:bg-red-50 rounded-lg transition-all"
                                    >
                                        <Trash2 size={16} />
                                    </button>
                                </div>

                                <!-- Inventory Item Image -->
                                <div class="flex items-center gap-3">
                                    <div class="flex gap-2">
                                        <button
                                            class="text-xs font-medium px-2 py-1 rounded-lg transition-colors {item.imageType ===
                                            'upload'
                                                ? 'bg-burnt-orange text-white'
                                                : 'bg-stone-200 text-stone-600'}"
                                            onclick={() =>
                                                (item.imageType = "upload")}
                                        >
                                            Upload
                                        </button>
                                        <button
                                            class="text-xs font-medium px-2 py-1 rounded-lg transition-colors {item.imageType ===
                                            'url'
                                                ? 'bg-burnt-orange text-white'
                                                : 'bg-stone-200 text-stone-600'}"
                                            onclick={() =>
                                                (item.imageType = "url")}
                                        >
                                            URL
                                        </button>
                                        <button
                                            class="text-xs font-medium px-2 py-1 rounded-lg transition-colors {item.imageType ===
                                            'icon'
                                                ? 'bg-burnt-orange text-white'
                                                : 'bg-stone-200 text-stone-600'}"
                                            onclick={() =>
                                                (item.imageType = "icon")}
                                        >
                                            Icone
                                        </button>
                                    </div>

                                    <div class="flex-1">
                                        {#if item.imageType === "upload"}
                                            <div class="relative">
                                                <input
                                                    type="file"
                                                    accept="image/*"
                                                    onchange={(e) => {
                                                        const input =
                                                            e.target as HTMLInputElement;
                                                        if (
                                                            input.files &&
                                                            input.files.length >
                                                                0
                                                        ) {
                                                            item.imageFile =
                                                                input.files[0];
                                                        }
                                                    }}
                                                    class="block w-full text-xs text-stone-500 file:mr-2 file:py-1 file:px-2 file:rounded-lg file:border-0 file:text-xs file:font-medium file:bg-stone-200 file:text-stone-700 hover:file:bg-stone-300"
                                                />
                                            </div>
                                        {:else}
                                            <input
                                                type="url"
                                                bind:value={item.imageURL}
                                                placeholder="https://..."
                                                class="w-full px-3 py-1 rounded-lg border border-stone-200 text-xs focus:outline-none focus:border-burnt-orange"
                                            />
                                        {/if}
                                        {#if item.imageType === "icon"}
                                            <input
                                                type="text"
                                                bind:value={item.iconName}
                                                placeholder="Nom de l'icône (ex: Sword)"
                                                class="w-full px-3 py-1 rounded-lg border border-stone-200 text-xs focus:outline-none focus:border-burnt-orange"
                                            />
                                        {/if}
                                    </div>
                                </div>
                            </div>
                        {/each}
                    </div>
                </div>
            </div>

            <div
                class="p-6 border-t border-stone-100 flex justify-end gap-3 sticky bottom-0 bg-white z-10"
            >
                <button
                    onclick={onClose}
                    class="px-4 py-2 text-stone-600 font-medium hover:bg-stone-100 rounded-xl transition-all"
                >
                    Annuler
                </button>
                <button
                    onclick={handleSubmit}
                    disabled={loading || !name || !race}
                    class="px-6 py-2 bg-burnt-orange text-white font-medium rounded-xl shadow-md hover:bg-opacity-90 transition-all disabled:opacity-50 disabled:cursor-not-allowed flex items-center gap-2"
                >
                    {#if loading}
                        <div
                            class="w-4 h-4 border-2 border-white/30 border-t-white rounded-full animate-spin"
                        ></div>
                    {/if}
                    {character?.id ? "Enregistrer" : "Créer"}
                </button>
            </div>
        </div>
    </div>
{/if}
