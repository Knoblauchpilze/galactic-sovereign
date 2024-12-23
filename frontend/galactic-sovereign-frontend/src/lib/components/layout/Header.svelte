<script lang="ts">
	import { type Snippet } from 'svelte';

	import { FlexContainer } from '@totocorpsoftwareinc/frontend-toolkit';

	interface Props {
		children?: Snippet;
	}

	let { children }: Props = $props();
</script>

<FlexContainer
	vertical={false}
	bgColor={'bg-black'}
	styling={'w-full h-1/10 fixed top-0 left-0 text-white'}
>
	{@render children?.()}
</FlexContainer>
