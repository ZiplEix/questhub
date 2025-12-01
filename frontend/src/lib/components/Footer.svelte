<script lang="ts">
    const currentYear = new Date().getFullYear();
</script>

<footer class="w-full py-8 text-center text-dark-gray/60 text-sm">
    <div class="flex justify-center gap-6 mb-4">
        <a href="/about" class="hover:text-burnt-orange transition-colors"
            >À propos</a
        >
        <a href="/support" class="hover:text-burnt-orange transition-colors"
            >Support</a
        >
    </div>
    <p>© {currentYear} QuestHub. Tous droits réservés.</p>
</footer>
