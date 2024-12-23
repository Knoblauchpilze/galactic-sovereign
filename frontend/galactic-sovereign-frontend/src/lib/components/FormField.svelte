<script lang="ts">
	import { type Snippet } from 'svelte';

	import { FlexContainer } from '@totocorpsoftwareinc/frontend-toolkit';

	interface Props {
		label?: string;
		labelId: string;
		children?: Snippet;
	}

	let { label = 'My label', labelId, children }: Props = $props();
</script>

<FlexContainer extensible={false} align={'stretch'} styling={'text-secondary'}>
	<label for={labelId}>{label}</label>
	{@render children?.()}
</FlexContainer>
