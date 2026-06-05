<script lang="ts">
    import { page } from "$app/state";
    import CharacterForm from "$lib/components/game/gm/CharacterForm.svelte";
    import Header from "$lib/components/Header.svelte";
    import { onMount } from "svelte";

    const gameId = page.params.id as string;
    let importedCharacter = $state<any>(null);

    onMount(() => {
        const urlParams = new URLSearchParams(window.location.search);
        if (urlParams.get("import") === "true") {
            const data = localStorage.getItem("importedCharacter");
            if (data) {
                importedCharacter = JSON.parse(data);
                localStorage.removeItem("importedCharacter"); // Clean up
            }
        }
    });
</script>

<div class="min-h-screen bg-stone-50">
    <Header />

    <main class="max-w-7xl mx-auto p-6 md:p-12">
        <div class="mb-8">
            <h1 class="text-3xl font-display font-bold text-dark-gray mb-2">
                Création de personnage
            </h1>
            <p class="text-stone-500">
                Créez un nouveau personnage pour votre table.
            </p>
        </div>

        <CharacterForm {gameId} character={importedCharacter} />
    </main>
</div>
