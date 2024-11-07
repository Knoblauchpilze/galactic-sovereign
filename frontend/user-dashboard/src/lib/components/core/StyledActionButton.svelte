<script lang="ts">
	interface Props {
		text?: string;
		enabled?: boolean;
		bgColor?: string;
		bgColorHover?: string;
	}

	let {
		text = 'Click me',
		enabled = true,
		bgColor = 'bg-enabled',
		bgColorHover = 'hover:bg-enabled-hover'
	}: Props = $props();
</script>

{#if enabled}
	<button class="px-8 py-2 rounded-[8px] {bgColor} {bgColorHover} text-white">{text}</button>
{:else}
	<button class="px-8 py-2 rounded-[8px] {bgColor} {bgColorHover} text-white" disabled
		>{text}</button
	>
{/if}
