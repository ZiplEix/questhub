<script lang="ts">
	import favicon from "$lib/assets/favicon.svg";
	import "../app.css";
	import { page } from "$app/state";
	import Footer from "$lib/components/Footer.svelte";

	let { children } = $props();

	const showFooter = $derived(!page.url.pathname.startsWith("/table/"));
</script>

<svelte:head>
	<link rel="icon" href={favicon} />
</svelte:head>

{#if showFooter}
	<div class="min-h-screen flex flex-col">
		<main class="grow flex flex-col">
			{@render children()}
		</main>
		<Footer />
	</div>
{:else}
	{@render children()}
{/if}
