<script lang="ts">
	import favicon from "$lib/assets/favicon.svg";
	import "../app.css";
	import { page } from "$app/state";
	import Footer from "$lib/components/Footer.svelte";
	import { env } from "$env/dynamic/public";

	let { children } = $props();

	const showFooter = $derived(!page.url.pathname.startsWith("/table/"));

	const umamiUrl = $derived(env.PUBLIC_UMAMI_URL);
	const umamiId = $derived(env.PUBLIC_UMAMI_WEBSITE_ID);
</script>

<svelte:head>
	<link rel="icon" href={favicon} />
	{#if umamiUrl && umamiId}
		<script defer src={umamiUrl} data-website-id={umamiId}></script>
	{/if}
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
