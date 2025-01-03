<script lang="ts">
	import '$styles/app.css';
	import { type Snippet } from 'svelte';

	import { HeroContainer } from '@totocorpsoftwareinc/frontend-toolkit';

	interface Props {
		children?: Snippet;
	}

	let { children }: Props = $props();
</script>

<svelte:head>
	<title>{'Admin Dashboard'}</title>
</svelte:head>
<HeroContainer width={'w-4/12'} height={'h-3/6'}>
	{@render children?.()}
</HeroContainer>
