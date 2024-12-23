<script lang="ts">
	import { type Snippet } from 'svelte';

	import { FlexContainer } from '@totocorpsoftwareinc/frontend-toolkit';

	interface Props {
		height?: string;
		children?: Snippet;
	}

	let { height = 'h-1/10', children }: Props = $props();
</script>

<FlexContainer
	vertical={false}
	bgColor={'bg-black'}
	styling="fixed bottom-0 w-full {height} text-white"
>
	{@render children?.()}
</FlexContainer>
