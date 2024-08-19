<script lang="ts">
	export let text: string = 'Click me';
	export let enabled: boolean = false;

	const bgColor = enabled ? 'bg-enabled' : 'bg-disabled';
	const bgColorHover = enabled ? 'hover:bg-enabled-hover' : 'hover:bg-disabled';
	const state = enabled ? '' : 'disabled';
</script>

{#if enabled}
	<button class="px-8 py-2 rounded-[8px] {bgColor} {bgColorHover} text-white">{text}</button>
{:else}
	<button class="px-8 py-2 rounded-[8px] {bgColor} {bgColorHover} text-white" disabled
		>{text}</button
	>
{/if}
