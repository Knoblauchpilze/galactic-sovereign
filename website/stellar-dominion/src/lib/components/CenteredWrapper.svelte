<script>
  export let width = 'w-full';
  export let height = 'h-full';
</script>

<!-- https://stackoverflow.com/questions/67852559/pass-svelte-component-as-props -->
<!-- https://svelte.dev/tutorial/slots -->
<div class="{width} {height} flex flex-col justify-around items-center">
	<slot />
</div>
