<!-- https://svelte.dev/docs/svelte-components -->
<script lang="ts">
	export let link: string = "/";
	export let text: string = "Home";
	export let showAsButton: boolean = false;
</script>

<!-- https://tailwindcss.com/docs/padding -->
<!-- https://tailwindcss.com/docs/hover-focus-and-other-states -->
{#if showAsButton}
	<a class="px-12 py-4 rounded-[8px] bg-primary hover:bg-primary-hover text-secondary" href={link}>{text}</a>
{:else}
	<a class="text-secondary hover:text-secondary-hover" href={link}>{text}</a>
{/if}
