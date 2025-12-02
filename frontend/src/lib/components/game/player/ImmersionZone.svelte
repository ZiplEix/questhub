<script lang="ts">
    import Scene from "../Scene.svelte"; // Reusing the Scene component
    import CompactTracker from "./CompactTracker.svelte";
    import OverlayChat from "./OverlayChat.svelte";

    // Mock damage notification
    let damageNotif = $state<{ value: number; visible: boolean }>({
        value: 0,
        visible: false,
    });

    // This could be triggered by a prop or event
    /*
    function showDamage(amount: number) {
        damageNotif = { value: amount, visible: true };
        setTimeout(() => damageNotif.visible = false, 2000);
    }
    */
</script>

<div class="relative w-full h-full">
    <!-- The Scene (Background + Pings + Dice Overlay) -->
    <Scene />

    <!-- Top Left: Compact Tracker -->
    <div class="absolute top-4 left-4 z-10">
        <CompactTracker />
    </div>

    <!-- Bottom Left: Overlay Chat -->
    <div class="absolute bottom-4 left-4 z-10">
        <OverlayChat />
    </div>

    <!-- Floating Damage Text (Example on Avatar/Player) -->
    <!-- In a real app, this would be positioned relative to the player's token or avatar -->
    {#if damageNotif.visible}
        <div
            class="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 text-red-500 font-black text-4xl animate-bounce-in drop-shadow-lg stroke-white"
        >
            -{damageNotif.value} PV
        </div>
    {/if}
</div>
