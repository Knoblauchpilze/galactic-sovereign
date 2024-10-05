<!-- https://svelte.dev/docs/svelte-components -->
<script lang="ts">
	import { StyledButton } from '$lib/components';

	export let text: string = 'Home';
	export let link: string = '/';
	export let showAsButton: boolean = false;
	export let disabled: boolean = false;
</script>

<!-- https://tailwindcss.com/docs/padding -->
<!-- https://tailwindcss.com/docs/hover-focus-and-other-states -->
{#if showAsButton}
	{#if disabled}
		<a href={link}>
			<StyledButton {text} enabled={false} />
		</a>
	{:else}
		<a href={link}>
			<StyledButton {text} />
		</a>
	{/if}
{:else}
	<a class="text-secondary hover:text-secondary-hover" href={link}>{text}</a>
{/if}
