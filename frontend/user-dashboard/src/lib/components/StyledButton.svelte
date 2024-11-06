<script lang="ts">
	interface Props {
		text?: string;
		enabled?: boolean;
	}

	let { text = 'Click me', enabled = true }: Props = $props();
</script>

{#if enabled}
	<button class="px-12 py-4 rounded-[8px] bg-primary hover:bg-primary-hover text-secondary"
		>{text}</button
	>
{:else}
	<button disabled class="px-12 py-4 rounded-[8px] bg-primary-selected text-secondary"
		>{text}</button
	>
{/if}
