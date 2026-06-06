<script lang="ts">
    const currentYear = new Date().getFullYear();
</script>

<footer class="w-full py-8 text-center text-stone-400 text-xs mt-auto border-t border-stone-200/50 bg-stone-50/20 backdrop-blur-xs shrink-0">
    <div class="flex justify-center flex-wrap gap-x-6 gap-y-2 mb-4 font-semibold">
        <a href="/about" class="hover:text-burnt-orange transition-colors">À propos</a>
        <a href="/mentions-legales" class="hover:text-burnt-orange transition-colors">Mentions légales</a>
        <a href="/contact" class="hover:text-burnt-orange transition-colors">Contact & Signalement</a>
    </div>
    <p class="font-medium">© {currentYear} QuestHub. Tous droits réservés.</p>
</footer>
