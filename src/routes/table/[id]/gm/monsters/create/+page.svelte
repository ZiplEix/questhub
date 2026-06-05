<script lang="ts">
    import CharacterForm from "$lib/components/game/gm/CharacterForm.svelte";
    import { page } from "$app/state";
    import Header from "$lib/components/Header.svelte";
    import { onMount } from "svelte";

    const gameId = $derived(page.params.id as string);
    let importedMonster = $state<any>(null);

    onMount(() => {
        const urlParams = new URLSearchParams(window.location.search);
        if (urlParams.get("import") === "true") {
            const data = localStorage.getItem("importedMonster");
            if (data) {
                importedMonster = JSON.parse(data);
                // Remove ID to ensure it's treated as a new creation
                delete importedMonster.id;
                localStorage.removeItem("importedMonster");
            }
        }
    });
</script>

<div class="min-h-screen bg-stone-50">
    <Header />

    <main class="max-w-7xl mx-auto p-6 md:p-12">
        <div class="mb-8">
            <h1 class="text-3xl font-display font-bold text-dark-gray mb-2">
                Création de créature
            </h1>
            <p class="text-stone-500">
                Ajoutez une nouvelle créature à votre bestiaire.
            </p>
        </div>

        <CharacterForm {gameId} character={importedMonster} type="MONSTER" />
    </main>
</div>
