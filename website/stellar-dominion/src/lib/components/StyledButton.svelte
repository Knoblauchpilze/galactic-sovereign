<script lang="ts">
	export let text: string = 'Click me';
	export let enabled: boolean = true;
</script>

{#if enabled}
	<button class="px-12 py-4 rounded-[8px] bg-primary hover:bg-primary-hover text-secondary"
		>{text}</button
	>
{:else}
	<!-- https://stackoverflow.com/questions/34684467/dont-redirect-on-button-click-with-css -->
	<button disabled class="px-12 py-4 rounded-[8px] bg-primary-selected text-secondary"
		>{text}</button
	>
{/if}
