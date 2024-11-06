<script lang="ts">
	import { type Snippet } from 'svelte';

	import { FlexContainer } from '$lib/components';

	interface Props {
		label?: string;
		labelId: string;
		children?: Snippet;
	}

	let { label = 'My label', labelId, children }: Props = $props();
</script>

<FlexContainer extensible={false} align={'stretch'} styling={'text-primary'}>
	<label for={labelId}>{label}</label>
	{@render children?.()}
</FlexContainer>
