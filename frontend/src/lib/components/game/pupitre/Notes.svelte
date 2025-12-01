<script lang="ts">
    let note = $state(`- Le garde à l'entrée semble corruptible.
- Ne pas oublier de demander à l'aubergiste pour la clé du sous-sol.
- Loot du dragon : 
  - 1500 po
  - Épée magique (+1)`);
</script>

<div class="h-full flex flex-col bg-[#fffdf0]">
    <textarea
        class="flex-1 w-full h-full p-6 bg-transparent resize-none focus:outline-none text-dark-gray font-serif leading-relaxed"
        bind:value={note}
        placeholder="Vos notes ici..."
    ></textarea>
</div>
