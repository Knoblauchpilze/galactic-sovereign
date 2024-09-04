<script lang="ts">
	export let text: string = 'Click me';
	export let enabled: boolean = false;
	export let negativeConfirmation: boolean = false;

	let bgColor: string;
	let bgColorHover: string;

	$: if (negativeConfirmation) {
		bgColor = enabled ? 'bg-disabled' : 'bg-enabled';
		bgColorHover = enabled ? 'hover:bg-disabled-hover' : 'hover:bg-enabled-hover';
	} else {
		bgColor = enabled ? 'bg-enabled' : 'bg-disabled';
		bgColorHover = enabled ? 'hover:bg-enabled-hover' : 'hover:bg-disabled-hover';
	}

	$: if (!enabled) {
		bgColorHover = bgColor;
	}
</script>

{#if enabled}
	<button class="px-8 py-2 rounded-[8px] {bgColor} {bgColorHover} text-white">{text}</button>
{:else}
	<button class="px-8 py-2 rounded-[8px] {bgColor} {bgColorHover} text-white" disabled
		>{text}</button
	>
{/if}
