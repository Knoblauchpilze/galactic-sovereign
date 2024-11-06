<script lang="ts">
	import '$styles/app.css';
	import { type Snippet } from 'svelte';

	import { HeroContainer } from '$lib/components';

	interface Props {
		children?: Snippet;
	}

	let { children }: Props = $props();
</script>

<svelte:head>
	<title>{'Admin Dashboard'}</title>
</svelte:head>
<HeroContainer>
	{@render children?.()}
</HeroContainer>
