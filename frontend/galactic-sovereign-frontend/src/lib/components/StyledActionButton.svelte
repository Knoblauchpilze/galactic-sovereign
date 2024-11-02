<script lang="ts">
	import { run } from 'svelte/legacy';

	interface Props {
		text?: string;
		enabled?: boolean;
		negativeConfirmation?: boolean;
	}

	let { text = 'Click me', enabled = false, negativeConfirmation = false }: Props = $props();

	let bgColor: string = $state();
	let bgColorHover: string = $state();

	run(() => {
		if (negativeConfirmation) {
			bgColor = enabled ? 'bg-disabled' : 'bg-enabled';
			bgColorHover = enabled ? 'hover:bg-disabled-hover' : 'hover:bg-enabled-hover';
		} else {
			bgColor = enabled ? 'bg-enabled' : 'bg-disabled';
			bgColorHover = enabled ? 'hover:bg-enabled-hover' : 'hover:bg-disabled-hover';
		}
	});

	run(() => {
		if (!enabled) {
			bgColorHover = bgColor;
		}
	});
</script>

{#if enabled}
	<button class="px-8 py-2 rounded-[8px] {bgColor} {bgColorHover} text-white">{text}</button>
{:else}
	<button class="px-8 py-2 rounded-[8px] {bgColor} {bgColorHover} text-white" disabled
		>{text}</button
	>
{/if}
