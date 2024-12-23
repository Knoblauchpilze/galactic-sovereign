<script lang="ts">
	import { StyledActionButton } from '@totocorpsoftwareinc/frontend-toolkit';

	interface Props {
		text?: string;
		enabled?: boolean;
	}

	let { text = 'Click me', enabled = true }: Props = $props();

	let bgColorHover = $derived(enabled ? 'hover:bg-disabled-hover' : 'hover:bg-disabled');
</script>

<StyledActionButton {text} {enabled} bgColor={'bg-disabled'} {bgColorHover} />
