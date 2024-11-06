<script lang="ts">
	interface Props {
		text?: string;
		link?: string;
		showAsButton?: boolean;
		disabled?: boolean;
	}

	let { text = 'Home', link = '/', showAsButton = false, disabled = false }: Props = $props();
</script>

{#if showAsButton}
	{#if disabled}
		<a href={link} class="px-12 py-4 rounded-[8px] bg-primary-selected text-secondary">
			{text}
		</a>
	{:else}
		<a
			href={link}
			class="px-12 py-4 rounded-[8px] bg-primary hover:bg-primary-hover text-secondary"
		>
			{text}
		</a>
	{/if}
{:else}
	<a class="text-secondary hover:text-secondary-hover" href={link}>{text}</a>
{/if}
