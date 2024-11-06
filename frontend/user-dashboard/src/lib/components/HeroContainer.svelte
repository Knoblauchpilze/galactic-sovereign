<script lang="ts">
	import { type Snippet } from 'svelte';

	import heroContainer from '$lib/stores/ui/heroContainer';

	import { FlexContainer } from '$lib/components';

	interface Props {
		children?: Snippet;
	}

	let { children }: Props = $props();

	let style = $derived(
		'absolute left-1/2 top-1/2 -translate-x-1/2 -translate-y-1/2 ' +
			$heroContainer.width +
			' ' +
			$heroContainer.height
	);
</script>

<div class="relative h-full bg-center bg-no-repeat bg-cover bg-primary">
	<FlexContainer bgColor={$heroContainer.color} styling={style}>
		{@render children?.()}
	</FlexContainer>
</div>

<style lang="postcss">
	:global(body),
	:global(html) {
		height: 100%;
	}

	:global(body) {
		margin: 0;
		padding: 0;
	}
</style>
