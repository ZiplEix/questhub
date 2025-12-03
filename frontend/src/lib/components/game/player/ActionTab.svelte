<script lang="ts">
    import { Sword, Wind, Footprints, Shield, Hand } from "lucide-svelte";
    import { getContext } from "svelte";

    const { showToast } = getContext<any>("toast");

    let attacks = [
        { name: "Épée longue", damage: "1d8 + 3", type: "Tranchant" },
        { name: "Dague", damage: "1d4 + 3", type: "Perçant" },
    ];

    let actions = [
        { name: "Se déplacer", icon: Footprints },
        { name: "Esquiver", icon: Shield },
        { name: "Foncer", icon: Wind },
        { name: "Se désengager", icon: Footprints },
    ];

    function handleAttack(attack: any) {
        showToast(`⚔️ Attaque : ${attack.name}`, "combat");
        // Logic to roll dice would go here
    }

    function handleAction(action: any) {
        showToast(`🏃 ${action.name}`, "info");
    }

    function requestSpeech() {
        showToast("✋ Demande de parole envoyée au MJ", "info");
    }
</script>

<div class="p-4 space-y-6">
    <!-- Speech Request -->
    <button
        onclick={requestSpeech}
        class="w-full bg-linear-to-r from-burnt-orange to-orange-600 text-white p-4 rounded-xl shadow-md hover:shadow-lg hover:scale-[1.02] transition-all flex items-center justify-center gap-3 font-bold text-lg"
    >
        <Hand size={24} />
        Demander la parole
    </button>
    <!-- Attacks -->
    <div>
        <h3
            class="text-xs font-bold text-stone-400 uppercase tracking-wider mb-3"
        >
            Attaques
        </h3>
        <div class="space-y-2">
            {#each attacks as attack}
                <button
                    onclick={() => handleAttack(attack)}
                    class="w-full bg-white p-3 rounded-xl border border-stone-200 shadow-sm hover:border-burnt-orange hover:shadow-md transition-all flex items-center justify-between group text-left"
                >
                    <div class="flex items-center gap-3">
                        <div
                            class="w-10 h-10 rounded-lg bg-stone-100 flex items-center justify-center text-stone-400 group-hover:bg-burnt-orange/10 group-hover:text-burnt-orange transition-colors"
                        >
                            <Sword size={20} />
                        </div>
                        <div>
                            <div
                                class="font-bold text-dark-gray group-hover:text-burnt-orange transition-colors"
                            >
                                {attack.name}
                            </div>
                            <div class="text-xs text-stone-400">
                                {attack.type}
                            </div>
                        </div>
                    </div>
                    <div
                        class="font-mono font-bold text-dark-gray bg-stone-50 px-2 py-1 rounded border border-stone-100"
                    >
                        {attack.damage}
                    </div>
                </button>
            {/each}
        </div>
    </div>

    <!-- Quick Actions -->
    <div>
        <h3
            class="text-xs font-bold text-stone-400 uppercase tracking-wider mb-3"
        >
            Actions Rapides
        </h3>
        <div class="grid grid-cols-3 gap-2">
            {#each actions as action}
                <button
                    onclick={() => handleAction(action)}
                    class="flex flex-col items-center justify-center p-3 bg-stone-50 rounded-xl border border-stone-100 hover:bg-white hover:border-stone-300 hover:shadow-sm transition-all gap-2"
                >
                    <action.icon size={20} class="text-stone-500" />
                    <span class="text-xs font-bold text-dark-gray text-center"
                        >{action.name}</span
                    >
                </button>
            {/each}
        </div>
    </div>

    <!-- Spells (Placeholder) -->
    <div>
        <h3
            class="text-xs font-bold text-stone-400 uppercase tracking-wider mb-3"
        >
            Sorts Préparés
        </h3>
        <div
            class="text-sm text-stone-500 italic text-center py-4 bg-stone-50 rounded-xl border border-stone-100 border-dashed"
        >
            Aucun sort préparé
        </div>
    </div>
</div>
