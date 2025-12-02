<script lang="ts">
    let { message, type = "info" } = $props();
</script>

<div
    class="animate-bounce-in bg-dark-gray/90 backdrop-blur-md text-white px-8 py-4 rounded-2xl shadow-2xl border border-white/10 flex items-center gap-4"
>
    <span class="text-2xl font-display font-bold tracking-wide">{message}</span>
</div>

<style>
    @keyframes bounce-in {
        0% {
            transform: scale(0.5) translateY(20px);
            opacity: 0;
        }
        60% {
            transform: scale(1.1) translateY(-5px);
            opacity: 1;
        }
        100% {
            transform: scale(1) translateY(0);
            opacity: 1;
        }
    }
    .animate-bounce-in {
        animation: bounce-in 0.3s cubic-bezier(0.175, 0.885, 0.32, 1.275)
            forwards;
    }
</style>
