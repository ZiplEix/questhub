<script lang="ts">
    import {
        X,
        Upload,
        Link,
        Plus,
        Trash2,
        Sticker,
        Save,
        ArrowLeft,
    } from "lucide-svelte";
    import { api } from "$lib/api";
    import { authClient } from "$lib/auth-client";
    import { page } from "$app/state";
    import { goto } from "$app/navigation";
    import type { Character } from "$lib/types/character";
    import { untrack } from "svelte";

    let {
        gameId,
        character = null,
        type = "PLAYER",
    } = $props<{
        gameId: string;
        character?: any;
        type?: "PLAYER" | "NPC" | "MONSTER";
    }>();

    // State
    let name = $state("");
    let race = $state("");
    let subRace = $state("");
    let maxHP = $state(10);
    let money = $state(0);
    let isNPC = $state(type === "NPC");
    let characterType = $state(type);
    let avatarType = $state<"upload" | "url">("upload");
    let avatarFile = $state<File | null>(null);
    let avatarURL = $state("");
    let stats = $state<Array<{ key: string; value: string; modifier: string }>>(
        [
            { key: "Force", value: "10", modifier: "0" },
            { key: "Dextérité", value: "10", modifier: "0" },
            { key: "Constitution", value: "10", modifier: "0" },
            { key: "Intelligence", value: "10", modifier: "0" },
            { key: "Sagesse", value: "10", modifier: "0" },
            { key: "Charisme", value: "10", modifier: "0" },
        ],
    );
    let inventory = $state<
        {
            name: string;
            quantity: string;
            imageType: "upload" | "url" | "icon";
            imageFile: File | null;
            imageURL: string;
            iconName: string;
            description: string;
        }[]
    >([]);
    let loading = $state(false);

    // New fields
    let initiative = $state(0);
    let armorClass = $state(10);
    let speed = $state(30);
    let age = $state("");
    let height = $state("");
    let weight = $state("");
    let maxSpells = $state(3);
    let spellsList = $state<
        {
            id: string;
            level: string;
            name: string;
            description: string;
            charges: string;
        }[]
    >([]);
    let abilities = $state("");
    let experience = $state(0);

    $effect(() => {
        const char = character;
        untrack(() => {
            if (char) {
                console.log("Editing character:", char);
                name = char.name || "";
                race = char.race || "";
                subRace = char.sub_race || "";
                maxHP = char.max_hp || 10;
                money = char.money || 0;
                isNPC = char.is_npc || false;
                characterType = char.type || (char.is_npc ? "NPC" : "PLAYER");

                // New fields
                initiative = char.initiative || 0;
                armorClass = char.armor_class || 10;
                speed = char.speed || 30;
                age = char.age || "";
                height = char.height || "";
                weight = char.weight || "";
                maxSpells = char.max_spells || 0;
                abilities = char.abilities || "";
                experience = char.experience || 0;

                // Spells
                spellsList = [];
                if (char.spells) {
                    try {
                        Object.entries(char.spells).forEach(
                            ([level, spells]) => {
                                if (Array.isArray(spells)) {
                                    spells.forEach((spell) => {
                                        if (typeof spell === "string") {
                                            spellsList.push({
                                                id: crypto.randomUUID(),
                                                level,
                                                name: spell,
                                                description: "",
                                                charges: "",
                                            });
                                        } else {
                                            spellsList.push({
                                                id: crypto.randomUUID(),
                                                level,
                                                name: spell.name,
                                                description:
                                                    spell.description || "",
                                                charges: spell.charges || "",
                                            });
                                        }
                                    });
                                }
                            },
                        );
                    } catch (e) {
                        console.error("Error parsing spells:", e);
                    }
                }

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
                        stats = Object.entries(char.stats).map(
                            ([key, data]) => {
                                const statData = data as {
                                    value: number;
                                    modifier: number;
                                };
                                return {
                                    key,
                                    value: statData.value?.toString() || "10",
                                    modifier:
                                        statData.modifier?.toString() || "0",
                                };
                            },
                        );
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
                                description: item.description || "",
                            });
                        });
                    } catch (e) {
                        console.error("Error parsing inventory:", e);
                    }
                }
            } else {
                // Reset form if opening for create
                name = "";
                race = "";
                subRace = "";
                maxHP = 10;
                money = 0;
                isNPC = type === "NPC";
                characterType = type;
                avatarType = "upload";
                avatarFile = null;
                avatarURL = "";
                stats = [];
                inventory = [];
                loading = false;
                initiative = 0;
                armorClass = 10;
                speed = 30;
                age = "";
                height = "";
                weight = "";
                maxSpells = 3;
                spellsList = [];
                abilities = "";
                experience = 0;
            }
        });
    });

    function addStat() {
        stats.push({ key: "", value: "10", modifier: "0" });
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
                description: "",
            },
        ];
    }

    function removeInventoryItem(index: number) {
        inventory = inventory.filter((_, i) => i !== index);
    }

    function addSpell() {
        spellsList = [
            ...spellsList,
            {
                id: crypto.randomUUID(),
                level: "0",
                name: "",
                description: "",
                charges: "",
            },
        ];
    }

    function removeSpell(id: string) {
        spellsList = spellsList.filter((s) => s.id !== id);
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

        try {
            const { data: tokenData } = await authClient.token();
            if (tokenData?.token) {
                const formData = new FormData();
                formData.append("name", name);
                formData.append("race", race);
                if (subRace) formData.append("sub_race", subRace);
                formData.append("max_hp", maxHP.toString());
                formData.append("money", money.toString());
                formData.append("is_npc", isNPC.toString());
                formData.append("type", characterType);

                // New fields
                formData.append("initiative", initiative.toString());
                formData.append("armor_class", armorClass.toString());
                formData.append("speed", speed.toString());
                formData.append("age", age);
                formData.append("height", height);
                formData.append("weight", weight);
                formData.append("max_spells", maxSpells.toString());
                formData.append("abilities", abilities);
                formData.append("experience", experience.toString());

                // Spells
                const spellsObj: Record<
                    string,
                    { name: string; description: string; charges: string }[]
                > = {};
                spellsList.forEach((spell) => {
                    if (spell.name) {
                        if (!spellsObj[spell.level]) {
                            spellsObj[spell.level] = [];
                        }
                        spellsObj[spell.level].push({
                            name: spell.name,
                            description: spell.description,
                            charges: spell.charges,
                        });
                    }
                });
                formData.append("spells", JSON.stringify(spellsObj));

                if (avatarType === "upload" && avatarFile) {
                    formData.append("avatar", avatarFile);
                } else if (avatarURL) {
                    formData.append("avatar_url", avatarURL);
                }

                const statsObj = stats.reduce(
                    (acc, curr) => ({
                        ...acc,
                        [curr.key]: {
                            value: parseInt(curr.value) || 10,
                            modifier: parseInt(curr.modifier) || 0,
                        },
                    }),
                    {},
                );
                formData.append("stats", JSON.stringify(statsObj));

                // Prepare inventory items
                const inventoryItems = inventory.map((item) => ({
                    name: item.name,
                    quantity: item.quantity,
                    image_url:
                        item.imageType === "url" ? item.imageURL : undefined,
                    icon_name:
                        item.imageType === "icon" ? item.iconName : undefined,
                    description: item.description,
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

                // Navigation
                if (characterType === "MONSTER") {
                    goto(`/table/${gameId}/gm/settings?tab=bestiary`);
                } else {
                    goto(`/table/${gameId}/gm/settings?tab=characters`);
                }
            }
        } catch (error) {
            console.error("Failed to save character:", error);
            alert(
                `Erreur lors de la sauvegarde ${characterType === "MONSTER" ? "de la créature" : "du personnage"}.`,
            );
        } finally {
            loading = false;
        }
    }

    function cancel() {
        if (characterType === "MONSTER") {
            goto(`/table/${gameId}/gm/settings?tab=bestiary`);
        } else {
            goto(`/table/${gameId}/gm/settings?tab=characters`);
        }
    }
</script>

<div class="bg-white rounded-2xl shadow-sm border border-stone-100 p-6 md:p-8">
    <div class="flex justify-between items-center mb-8">
        <h2 class="text-2xl font-bold text-dark-gray">
            {character?.id
                ? characterType === "MONSTER"
                    ? "Modifier la créature"
                    : "Modifier le personnage"
                : characterType === "MONSTER"
                  ? "Créer une créature"
                  : "Créer un personnage"}
        </h2>
        <div class="flex gap-3">
            <button
                onclick={cancel}
                class="px-4 py-2 text-stone-600 font-medium hover:text-dark-gray transition-colors flex items-center gap-2"
            >
                <ArrowLeft size={18} />
                Annuler
            </button>
            <button
                onclick={handleSubmit}
                disabled={loading || !name || !race || !maxHP}
                class="px-6 py-2 bg-burnt-orange text-white rounded-xl font-bold shadow-md hover:bg-opacity-90 transition-all hover:-translate-y-0.5 disabled:opacity-50 disabled:cursor-not-allowed flex items-center gap-2"
            >
                {#if loading}
                    <div
                        class="w-4 h-4 border-2 border-white/30 border-t-white rounded-full animate-spin"
                    ></div>
                {:else}
                    <Save size={18} />
                {/if}
                {character ? "Enregistrer" : "Créer"}
            </button>
        </div>
    </div>

    <div class="space-y-8">
        <!-- Basic Info -->
        <section class="space-y-4">
            <h3
                class="text-lg font-bold text-dark-gray border-b border-stone-100 pb-2"
            >
                Informations de base
            </h3>
            <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
                <!-- Avatar -->
                <div class="space-y-2">
                    <label
                        class="text-sm font-bold text-dark-gray block"
                        for="avatar-upload">Avatar</label
                    >
                    <div class="flex gap-4 mb-2">
                        <button
                            class="flex items-center gap-2 text-sm font-medium {avatarType ===
                            'upload'
                                ? 'text-burnt-orange'
                                : 'text-stone-400'}"
                            onclick={() => (avatarType = "upload")}
                            id="avatar-upload"
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

                <div class="space-y-4">
                    <div>
                        <label
                            for="name"
                            class="block text-sm font-bold text-dark-gray mb-1"
                        >
                            Nom
                        </label>
                        <input
                            type="text"
                            id="name"
                            bind:value={name}
                            class="w-full px-4 py-2 rounded-xl border border-stone-200 focus:outline-none focus:ring-2 focus:ring-burnt-orange/20 focus:border-burnt-orange transition-all"
                            placeholder={characterType === "MONSTER"
                                ? "Nom de la créature"
                                : "Nom du personnage"}
                        />
                    </div>

                    <div class="grid grid-cols-2 gap-4">
                        <div>
                            <label
                                for="race"
                                class="block text-sm font-bold text-dark-gray mb-1"
                            >
                                {characterType === "MONSTER"
                                    ? "Espèce"
                                    : "Race"}
                            </label>
                            <input
                                type="text"
                                id="race"
                                bind:value={race}
                                class="w-full px-4 py-2 rounded-xl border border-stone-200 focus:outline-none focus:ring-2 focus:ring-burnt-orange/20 focus:border-burnt-orange transition-all"
                                placeholder={characterType === "MONSTER"
                                    ? "Dragon, Gobelin..."
                                    : "Humain, Elfe..."}
                            />
                        </div>
                        {#if characterType === "MONSTER"}
                            <div>
                                <label
                                    for="subRace"
                                    class="block text-sm font-bold text-dark-gray mb-1"
                                >
                                    Sous-espèce
                                </label>
                                <input
                                    type="text"
                                    id="subRace"
                                    bind:value={subRace}
                                    class="w-full px-4 py-2 rounded-xl border border-stone-200 focus:outline-none focus:ring-2 focus:ring-burnt-orange/20 focus:border-burnt-orange transition-all"
                                    placeholder="Rouge, Archer..."
                                />
                            </div>
                        {/if}
                    </div>

                    <div class="grid grid-cols-2 gap-4">
                        <div>
                            <label
                                for="maxHP"
                                class="block text-sm font-bold text-dark-gray mb-1"
                            >
                                PV Max
                            </label>
                            <input
                                type="number"
                                id="maxHP"
                                bind:value={maxHP}
                                min="1"
                                class="w-full px-4 py-2 rounded-xl border border-stone-200 focus:outline-none focus:ring-2 focus:ring-burnt-orange/20 focus:border-burnt-orange transition-all"
                            />
                        </div>
                        {#if characterType !== "MONSTER"}
                            <div>
                                <label
                                    for="money"
                                    class="block text-sm font-bold text-dark-gray mb-1"
                                >
                                    Argent (PO)
                                </label>
                                <input
                                    type="number"
                                    id="money"
                                    bind:value={money}
                                    min="0"
                                    class="w-full px-4 py-2 rounded-xl border border-stone-200 focus:outline-none focus:ring-2 focus:ring-burnt-orange/20 focus:border-burnt-orange transition-all"
                                />
                            </div>
                        {/if}
                    </div>

                    {#if characterType !== "MONSTER"}
                        <div>
                            <label
                                class="text-sm font-bold text-dark-gray block mb-1"
                                for="character-type">Type</label
                            >
                            <div class="flex bg-stone-100 p-1 rounded-xl w-fit">
                                <button
                                    class="px-6 py-2 rounded-xl text-sm font-medium transition-all {isNPC
                                        ? 'text-stone-500 hover:text-dark-gray'
                                        : 'bg-white text-dark-gray shadow-sm'}"
                                    onclick={() => {
                                        isNPC = false;
                                        characterType = "PLAYER";
                                    }}
                                >
                                    Joueur (PJ)
                                </button>
                                <button
                                    class="px-6 py-2 rounded-xl text-sm font-medium transition-all {!isNPC
                                        ? 'text-stone-500 hover:text-dark-gray'
                                        : 'bg-white text-dark-gray shadow-sm'}"
                                    onclick={() => {
                                        // The statsObj calculation was not used here and caused a type error if modifier was undefined.
                                        // Removing it as it's not needed for the button's primary function.
                                        isNPC = true;
                                        characterType = "NPC";
                                    }}
                                >
                                    PNJ
                                </button>
                            </div>
                        </div>
                    {/if}
                </div>
            </div>
        </section>

        <!-- Physical Attributes -->
        <section class="space-y-4">
            <h3
                class="text-lg font-bold text-dark-gray border-b border-stone-100 pb-2"
            >
                Attributs physiques
            </h3>
            <div class="grid grid-cols-2 md:grid-cols-4 gap-4">
                <div>
                    <label
                        for="initiative"
                        class="block text-sm font-bold text-dark-gray mb-1"
                    >
                        Initiative
                    </label>
                    <input
                        type="number"
                        id="initiative"
                        bind:value={initiative}
                        class="w-full px-4 py-2 rounded-xl border border-stone-200 focus:outline-none focus:ring-2 focus:ring-burnt-orange/20 focus:border-burnt-orange transition-all"
                    />
                </div>
                <div>
                    <label
                        for="armorClass"
                        class="block text-sm font-bold text-dark-gray mb-1"
                    >
                        Classe d'Armure (CA)
                    </label>
                    <input
                        type="number"
                        id="armorClass"
                        bind:value={armorClass}
                        class="w-full px-4 py-2 rounded-xl border border-stone-200 focus:outline-none focus:ring-2 focus:ring-burnt-orange/20 focus:border-burnt-orange transition-all"
                    />
                </div>
                <div>
                    <label
                        for="speed"
                        class="block text-sm font-bold text-dark-gray mb-1"
                    >
                        Vitesse (m)
                    </label>
                    <input
                        type="number"
                        id="speed"
                        bind:value={speed}
                        class="w-full px-4 py-2 rounded-xl border border-stone-200 focus:outline-none focus:ring-2 focus:ring-burnt-orange/20 focus:border-burnt-orange transition-all"
                    />
                </div>
                {#if characterType !== "MONSTER"}
                    <div>
                        <label
                            for="age"
                            class="block text-sm font-bold text-dark-gray mb-1"
                        >
                            Age
                        </label>
                        <input
                            type="text"
                            id="age"
                            bind:value={age}
                            class="w-full px-4 py-2 rounded-xl border border-stone-200 focus:outline-none focus:ring-2 focus:ring-burnt-orange/20 focus:border-burnt-orange transition-all"
                        />
                    </div>
                {/if}
                <div>
                    <label
                        for="height"
                        class="block text-sm font-bold text-dark-gray mb-1"
                    >
                        {characterType === "MONSTER" ? "Taille" : "Taille (cm)"}
                    </label>
                    <input
                        type="text"
                        id="height"
                        bind:value={height}
                        class="w-full px-4 py-2 rounded-xl border border-stone-200 focus:outline-none focus:ring-2 focus:ring-burnt-orange/20 focus:border-burnt-orange transition-all"
                    />
                </div>
                <div>
                    <label
                        for="weight"
                        class="block text-sm font-bold text-dark-gray mb-1"
                    >
                        Poids (kg)
                    </label>
                    <input
                        type="text"
                        id="weight"
                        bind:value={weight}
                        class="w-full px-4 py-2 rounded-xl border border-stone-200 focus:outline-none focus:ring-2 focus:ring-burnt-orange/20 focus:border-burnt-orange transition-all"
                    />
                </div>
            </div>
        </section>

        <!-- Abilities -->
        <section class="space-y-4">
            <h3
                class="text-lg font-bold text-dark-gray border-b border-stone-100 pb-2"
            >
                Capacités & Expérience
            </h3>
            <div class="space-y-2">
                <label for="abilities" class="text-sm font-bold text-dark-gray"
                    >Capacités</label
                >
                <textarea
                    id="abilities"
                    bind:value={abilities}
                    rows="3"
                    class="w-full px-4 py-2 rounded-xl border border-stone-200 focus:outline-none focus:ring-2 focus:ring-burnt-orange/20 focus:border-burnt-orange transition-all resize-y"
                    placeholder="Description des capacités..."
                ></textarea>
            </div>
            {#if characterType !== "MONSTER"}
                <div>
                    <label
                        for="experience"
                        class="block text-sm font-bold text-dark-gray mb-1"
                    >
                        Expérience (XP)
                    </label>
                    <input
                        type="number"
                        id="experience"
                        bind:value={experience}
                        min="0"
                        class="w-full px-4 py-2 rounded-xl border border-stone-200 focus:outline-none focus:ring-2 focus:ring-burnt-orange/20 focus:border-burnt-orange transition-all"
                    />
                </div>
            {/if}
        </section>

        <div class="grid grid-cols-1 lg:grid-cols-2 gap-8">
            <!-- Stats -->
            <div class="space-y-4">
                <div class="flex justify-between items-center">
                    <label class="text-sm font-bold text-dark-gray" for="stats"
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
                    <p
                        class="text-sm text-stone-400 italic bg-stone-50 p-4 rounded-xl border border-stone-100 text-center"
                    >
                        Aucune statistique définie.
                    </p>
                {/if}

                <div class="space-y-4">
                    {#each stats as stat, i}
                        <div
                            class="p-4 bg-stone-50 rounded-xl border border-stone-100 flex gap-4 items-start"
                        >
                            <div class="flex-1 space-y-2">
                                <label
                                    class="text-xs font-bold text-stone-500 uppercase"
                                    for="stat-key">Nom</label
                                >
                                <input
                                    type="text"
                                    bind:value={stat.key}
                                    placeholder="Nom (ex: Force)"
                                    class="w-full px-3 py-2 rounded-lg border border-stone-200 text-sm focus:outline-none focus:border-burnt-orange"
                                />
                            </div>
                            <div class="w-24 space-y-2">
                                <label
                                    class="text-xs font-bold text-stone-500 uppercase"
                                    for="stat-value">Valeur</label
                                >
                                <input
                                    type="text"
                                    bind:value={stat.value}
                                    placeholder="10"
                                    class="w-full px-3 py-2 rounded-lg border border-stone-200 text-sm focus:outline-none focus:border-burnt-orange text-center"
                                />
                            </div>
                            <div class="w-20 space-y-2">
                                <label
                                    class="text-xs font-bold text-stone-500 uppercase"
                                    for="stat-modifier">Mod.</label
                                >
                                <input
                                    type="number"
                                    bind:value={stat.modifier}
                                    placeholder="+0"
                                    class="w-full px-3 py-2 rounded-lg border border-stone-200 text-sm focus:outline-none focus:border-burnt-orange text-center font-mono"
                                />
                            </div>
                            <div class="pt-6">
                                <button
                                    onclick={() => removeStat(i)}
                                    class="p-2 text-stone-400 hover:text-red-500 hover:bg-red-50 rounded-lg transition-all"
                                >
                                    <Trash2 size={16} />
                                </button>
                            </div>
                        </div>
                    {/each}
                </div>
            </div>

            <!-- Inventory -->
            <div class="space-y-4">
                <div class="flex justify-between items-center">
                    {#if characterType === "MONSTER"}
                        <label
                            class="text-sm font-bold text-dark-gray"
                            for="loot">Loot</label
                        >
                    {:else}
                        <label
                            class="text-sm font-bold text-dark-gray"
                            for="inventory">Inventaire</label
                        >
                    {/if}
                    <button
                        onclick={addInventoryItem}
                        class="text-xs flex items-center gap-1 text-burnt-orange font-medium hover:text-burnt-orange/80"
                    >
                        <Plus size={14} /> Ajouter
                    </button>
                </div>

                {#if inventory.length === 0}
                    <p
                        class="text-sm text-stone-400 italic bg-stone-50 p-4 rounded-xl border border-stone-100 text-center"
                    >
                        {#if characterType === "MONSTER"}
                            Loot vide.
                        {:else}
                            Inventaire vide.
                        {/if}
                    </p>
                {/if}

                <div class="space-y-4">
                    {#each inventory as item, i}
                        <div
                            class="p-4 bg-stone-50 rounded-xl border border-stone-100 space-y-3"
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
                            <textarea
                                bind:value={item.description}
                                placeholder="Description de l'objet..."
                                class="w-full px-3 py-2 rounded-lg border border-stone-200 text-sm focus:outline-none focus:border-burnt-orange min-h-[60px]"
                            ></textarea>

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
                                        onclick={() => (item.imageType = "url")}
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
                                                        input.files.length > 0
                                                    ) {
                                                        item.imageFile =
                                                            input.files[0];
                                                    }
                                                }}
                                                class="block w-full text-xs text-stone-500 file:mr-2 file:py-1 file:px-2 file:rounded-lg file:border-0 file:text-xs file:font-medium file:bg-stone-200 file:text-stone-700 hover:file:bg-stone-300"
                                            />
                                        </div>
                                    {:else if item.imageType === "url"}
                                        <input
                                            type="url"
                                            bind:value={item.imageURL}
                                            placeholder="https://..."
                                            class="w-full px-3 py-1 rounded-lg border border-stone-200 text-xs focus:outline-none focus:border-burnt-orange"
                                        />
                                    {:else if item.imageType === "icon"}
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

        <!-- Spells -->
        <div class="space-y-4">
            <div class="flex justify-between items-center">
                <label class="text-sm font-bold text-dark-gray" for="spells"
                    >Sorts</label
                >
                <button
                    onclick={addSpell}
                    class="text-xs flex items-center gap-1 text-burnt-orange font-medium hover:text-burnt-orange/80"
                >
                    <Plus size={14} /> Ajouter
                </button>
            </div>

            {#if spellsList.length === 0}
                <p
                    class="text-sm text-stone-400 italic bg-stone-50 p-4 rounded-xl border border-stone-100 text-center"
                >
                    Aucun sort.
                </p>
            {/if}

            <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                {#each spellsList as spell (spell.id)}
                    <div
                        class="flex flex-col gap-2 w-full p-4 border border-stone-100 rounded-xl bg-stone-50"
                    >
                        <div class="flex gap-2">
                            <select
                                bind:value={spell.level}
                                class="w-20 px-3 py-2 rounded-lg border border-stone-200 text-sm focus:outline-none focus:border-burnt-orange bg-white"
                            >
                                <option value="0">Niv 0</option>
                                <option value="1">Niv 1</option>
                                <option value="2">Niv 2</option>
                                <option value="3">Niv 3</option>
                                <option value="4">Niv 4</option>
                                <option value="5">Niv 5</option>
                                <option value="6">Niv 6</option>
                                <option value="7">Niv 7</option>
                                <option value="8">Niv 8</option>
                                <option value="9">Niv 9</option>
                            </select>
                            <input
                                type="text"
                                bind:value={spell.name}
                                placeholder="Nom du sort"
                                class="flex-1 px-3 py-2 rounded-lg border border-stone-200 text-sm focus:outline-none focus:border-burnt-orange"
                            />
                            <button
                                onclick={() => removeSpell(spell.id)}
                                class="p-2 text-stone-400 hover:text-red-500 hover:bg-red-50 rounded-lg transition-all"
                            >
                                <Trash2 size={16} />
                            </button>
                        </div>
                        <div class="flex gap-2">
                            <input
                                type="text"
                                bind:value={spell.charges}
                                placeholder="Charges (ex: 3/jour)"
                                class="w-1/3 px-3 py-2 rounded-lg border border-stone-200 text-sm focus:outline-none focus:border-burnt-orange"
                            />
                            <textarea
                                bind:value={spell.description}
                                placeholder="Description du sort..."
                                rows="1"
                                class="flex-1 px-3 py-2 rounded-lg border border-stone-200 text-sm focus:outline-none focus:border-burnt-orange resize-none"
                            ></textarea>
                        </div>
                    </div>
                {/each}
            </div>
        </div>
    </div>
</div>
```
