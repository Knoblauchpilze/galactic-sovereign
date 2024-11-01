<script lang="ts">
	import '$styles/app.css';
	import { type Snippet } from 'svelte';

	import pageTitle from '$lib/stores/ui/pageTitle';

	import { HeroContainer } from '$lib/components';

	interface Props {
		children?: Snippet;
	}

	let { children }: Props = $props();
</script>

<!-- https://stackoverflow.com/questions/74243719/how-to-change-page-title-dynamically-in-sveltekit -->
<svelte:head>
	<title>{$pageTitle}</title>
</svelte:head>
<!-- https://stackoverflow.com/questions/70805041/background-image-in-tailwindcss-using-dynamic-url-react-js -->
<HeroContainer>
	{@render children?.()}
</HeroContainer>
