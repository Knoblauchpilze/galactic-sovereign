<script lang="ts">
	export let width: string = 'w-full';
	export let height: string = 'h-full';
</script>

<!-- https://stackoverflow.com/questions/67852559/pass-svelte-component-as-props -->
<!-- https://svelte.dev/tutorial/slots -->
<div class="{width} {height} flex flex-col justify-around items-center">
	<slot />
</div>
