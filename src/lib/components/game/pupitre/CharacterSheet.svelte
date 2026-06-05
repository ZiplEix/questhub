<script lang="ts">
    import { Shield, Heart, Zap, Sword, Book, Swords } from "lucide-svelte";
    import type { Character } from "$lib/types/character";

    let { character } = $props<{ character: Character }>();

    let stats = $derived(
        character.stats
            ? Object.entries(character.stats).map(([key, data]) => {
                  let value = "10";
                  let mod = 0;

                  if (
                      typeof data === "object" &&
                      data !== null &&
                      "value" in data
                  ) {
                      const statData = data as {
                          value: number;
                          modifier: number;
                      };
                      value = statData.value?.toString() || "10";
                      mod = statData.modifier || 0;
                  } else {
                      value = String(data);
                  }

                  return {
                      key,
                      value,
                      mod,
                  };
              })
            : [],
    );

    let spellsByLevel = $derived(
        character.spells
            ? (Object.entries(character.spells).sort(([a], [b]) =>
                  a.localeCompare(b),
              ) as [
                  string,
                  { name: string; description: string; charges: string }[],
              ][])
            : [],
    );

    let hpPercent = $derived((character.current_hp / character.max_hp) * 100);
</script>

<div class="h-full overflow-y-auto p-4 space-y-6">
    <!-- Vital Header Section -->
    <div class="bg-white rounded-2xl p-4 border border-stone-100 shadow-sm">
        <div class="flex items-center gap-4 mb-4">
            <img
                src={character.avatar_url ||
                    `https://api.dicebear.com/7.x/avataaars/svg?seed=${character.name}`}
                alt={character.name}
                class="w-16 h-16 rounded-2xl bg-stone-100 shadow-sm border-2 border-white object-cover"
            />
            <div>
                <h2 class="font-display font-bold text-2xl text-dark-gray">
                    {character.name}
                </h2>
                <div class="text-sm text-stone-500 font-medium">
                    {character.race}
                    {character.sub_race ? `(${character.sub_race})` : ""}
                    {#if character.experience}
                        • {character.experience} XP
                    {/if}
                </div>
            </div>
        </div>

        <!-- HP Bar -->
        <div class="mb-4">
            <div class="flex justify-between text-sm font-bold mb-1">
                <span class="text-stone-400">Points de Vie</span>
                <span
                    class={character.current_hp < character.max_hp / 3
                        ? "text-red-500"
                        : "text-green-600"}
                >
                    {character.current_hp} / {character.max_hp}
                </span>
            </div>
            <div
                class="h-4 bg-stone-100 rounded-full overflow-hidden border border-stone-200"
            >
                <div
                    class="h-full rounded-full transition-all duration-500 {character.current_hp <
                    character.max_hp / 3
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
                <span class="font-bold text-lg text-dark-gray"
                    >{character.armor_class}</span
                >
                <span class="text-[10px] uppercase font-bold text-stone-400"
                    >CA</span
                >
            </div>
            <div
                class="flex flex-col items-center p-2 bg-stone-50 rounded-xl border border-stone-100"
            >
                <Zap size={16} class="text-stone-400 mb-1" />
                <span class="font-bold text-lg text-dark-gray"
                    >{character.speed}</span
                >
                <span class="text-[10px] uppercase font-bold text-stone-400"
                    >Vitesse</span
                >
            </div>
            <div
                class="flex flex-col items-center p-2 bg-stone-50 rounded-xl border border-stone-100"
            >
                <Swords size={16} class="text-stone-400 mb-1" />
                <span class="font-bold text-lg text-dark-gray"
                    >{character.initiative}</span
                >
                <span class="text-[10px] uppercase font-bold text-stone-400"
                    >Init</span
                >
            </div>
        </div>
    </div>

    <!-- Ability Scores -->
    {#if stats.length > 0}
        <div>
            <h3
                class="text-xs font-bold text-stone-400 uppercase tracking-wider mb-3"
            >
                Caractéristiques
            </h3>
            <div class="grid grid-cols-3 gap-2">
                {#each stats as stat}
                    <div class="bg-stone-50 p-2 rounded-lg text-center">
                        <div
                            class="text-[10px] uppercase font-bold text-stone-400 truncate"
                            title={stat.key}
                        >
                            {stat.key}
                        </div>
                        <div class="font-bold text-dark-gray text-lg">
                            {stat.value}
                        </div>
                        <div
                            class="text-xs font-bold font-mono text-burnt-orange"
                        >
                            {#if !isNaN(stat.mod)}
                                {stat.mod >= 0 ? "+" : ""}{stat.mod}
                            {/if}
                        </div>
                    </div>
                {/each}
            </div>
        </div>
    {/if}

    <!-- Abilities / Skills (Text) -->
    {#if character.abilities}
        <div>
            <h3
                class="text-xs font-bold text-stone-400 uppercase tracking-wider mb-3"
            >
                Capacités
            </h3>
            <div
                class="bg-white p-4 rounded-xl border border-stone-100 text-sm text-stone-600 whitespace-pre-wrap"
            >
                {character.abilities}
            </div>
        </div>
    {/if}

    <!-- Spells -->
    {#if spellsByLevel.length > 0}
        <div>
            <h3
                class="text-xs font-bold text-stone-400 uppercase tracking-wider mb-3"
            >
                Sorts
            </h3>
            <div class="space-y-4">
                {#each spellsByLevel as [level, spells]}
                    <div>
                        <h4
                            class="text-xs font-bold text-burnt-orange mb-2 flex items-center gap-2"
                        >
                            <Book size={12} />
                            Niveau {level}
                        </h4>
                        <div class="space-y-2">
                            {#each spells as spell}
                                <div
                                    class="bg-white p-3 rounded-xl border border-stone-100 shadow-sm"
                                >
                                    <div
                                        class="flex justify-between items-start mb-1"
                                    >
                                        <span
                                            class="font-bold text-dark-gray text-sm"
                                            >{spell.name}</span
                                        >
                                        {#if spell.charges}
                                            <span
                                                class="text-[10px] font-mono bg-stone-100 px-1.5 py-0.5 rounded text-stone-500"
                                            >
                                                {spell.charges} charges
                                            </span>
                                        {/if}
                                    </div>
                                    {#if spell.description}
                                        <p class="text-xs text-stone-500">
                                            {spell.description}
                                        </p>
                                    {/if}
                                </div>
                            {/each}
                        </div>
                    </div>
                {/each}
            </div>
        </div>
    {/if}
</div>
