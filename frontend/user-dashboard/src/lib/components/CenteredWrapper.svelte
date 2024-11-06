<script lang="ts">
	import { type Snippet } from 'svelte';

	import { FlexContainer } from '$lib/components';

	interface Props {
		width?: string;
		height?: string;
		bgColor?: string;
		children?: Snippet;
	}

	let {
		width = 'w-full',
		height = 'h-full',
		bgColor = 'bg-transparent',
		children
	}: Props = $props();
</script>

<FlexContainer
	vertical={false}
	justify={'center'}
	align={'stretch'}
	extensible={false}
	styling="{width} {height}"
>
	<FlexContainer {bgColor}>
		{@render children?.()}
	</FlexContainer>
</FlexContainer>
