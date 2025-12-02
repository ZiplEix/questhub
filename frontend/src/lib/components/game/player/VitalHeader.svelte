<script lang="ts">
    import { Shield, Zap, Swords } from "lucide-svelte";

    let player = $state({
        name: "Eldric",
        avatar: "https://api.dicebear.com/7.x/avataaars/svg?seed=Eldric",
        hp: 15,
        maxHp: 24,
        ac: 16,
        speed: 30,
        initiative: "+2",
    });

    let hpPercent = $derived((player.hp / player.maxHp) * 100);
</script>

<div class="p-6 bg-white border-b border-stone-100">
    <div class="flex items-center gap-4 mb-4">
        <img
            src={player.avatar}
            alt={player.name}
            class="w-16 h-16 rounded-2xl bg-stone-100 shadow-sm border-2 border-white"
        />
        <div>
            <h2 class="font-display font-bold text-2xl text-dark-gray">
                {player.name}
            </h2>
            <div class="text-sm text-stone-500 font-medium">
                Guerrier Niv. 3
            </div>
        </div>
    </div>

    <!-- HP Bar -->
    <div class="mb-4">
        <div class="flex justify-between text-sm font-bold mb-1">
            <span class="text-stone-400">Points de Vie</span>
            <span
                class={player.hp < player.maxHp / 3
                    ? "text-red-500"
                    : "text-green-600"}
            >
                {player.hp} / {player.maxHp}
            </span>
        </div>
        <div
            class="h-4 bg-stone-100 rounded-full overflow-hidden border border-stone-200"
        >
            <div
                class="h-full rounded-full transition-all duration-500 {player.hp <
                player.maxHp / 3
                    ? 'bg-red-500'
                    : 'bg-green-500'}"
                style="width: {hpPercent}%"
            ></div>
        </div>
    </div>

    <!-- Defenses -->
    <div class="grid grid-cols-3 gap-3">
        <div
            class="flex flex-col items-center p-2 bg-stone-50 rounded-xl border border-stone-100"
        >
            <Shield size={16} class="text-stone-400 mb-1" />
            <span class="font-bold text-lg text-dark-gray">{player.ac}</span>
            <span class="text-[10px] uppercase font-bold text-stone-400"
                >CA</span
            >
        </div>
        <div
            class="flex flex-col items-center p-2 bg-stone-50 rounded-xl border border-stone-100"
        >
            <Zap size={16} class="text-stone-400 mb-1" />
            <span class="font-bold text-lg text-dark-gray">{player.speed}</span>
            <span class="text-[10px] uppercase font-bold text-stone-400"
                >Vitesse</span
            >
        </div>
        <div
            class="flex flex-col items-center p-2 bg-stone-50 rounded-xl border border-stone-100"
        >
            <Swords size={16} class="text-stone-400 mb-1" />
            <span class="font-bold text-lg text-dark-gray"
                >{player.initiative}</span
            >
            <span class="text-[10px] uppercase font-bold text-stone-400"
                >Init</span
            >
        </div>
    </div>
</div>
