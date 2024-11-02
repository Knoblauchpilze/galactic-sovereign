<script lang="ts">
	import { StyledActionButton } from '$lib/components';

	interface Props {
		text?: string;
		enabled?: boolean;
	}

	let { text = 'Click me', enabled = true }: Props = $props();

	let bgColor = $derived(enabled ? 'bg-enabled' : 'bg-disabled');
	let bgColorHover = $derived(enabled ? 'hover:bg-enabled-hover' : 'hover:bg-disabled');
</script>

<StyledActionButton {text} {enabled} {bgColor} {bgColorHover} />
