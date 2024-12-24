<script lang="ts">
	import '$styles/app.css';
	import { type Snippet } from 'svelte';

	import pageTitle from '$lib/stores/ui/pageTitle';
	import heroImage from '$lib/stores/ui/heroImage';
	import heroContainer from '$lib/stores/ui/heroContainer';

	import { HeroContainer } from '@totocorpsoftwareinc/frontend-toolkit';

	interface Props {
		children?: Snippet;
	}

	let { children }: Props = $props();
</script>

<!-- https://stackoverflow.com/questions/74243719/how-to-change-page-title-dynamically-in-sveltekit -->
<svelte:head>
	<title>{$pageTitle}</title>
</svelte:head>
<HeroContainer
	width={$heroContainer.width}
	height={$heroContainer.height}
	bgColor={$heroImage}
	bgOverlay={$heroContainer.color}
>
	{@render children?.()}
</HeroContainer>
