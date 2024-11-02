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

<!-- https://stackoverflow.com/questions/1122381/how-to-force-child-div-to-be-100-of-parent-divs-height-without-specifying-pare -->
<FlexContainer
	vertical={false}
	justify={'center'}
	align={'stretch'}
	extensible={false}
	styling="{width} {height}"
>
	<FlexContainer {bgColor}>
		<!-- https://svelte.dev/tutorial/slots -->
		{@render children?.()}
	</FlexContainer>
</FlexContainer>
