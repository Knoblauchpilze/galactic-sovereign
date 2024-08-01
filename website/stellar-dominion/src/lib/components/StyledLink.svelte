<!-- https://svelte.dev/docs/svelte-components -->
<script>
	export let link = "/";
	export let text = "Home";
</script>

<!-- https://tailwindcss.com/docs/padding -->
<!-- https://tailwindcss.com/docs/hover-focus-and-other-states -->
<a class="px-12 py-4 rounded-[8px] bg-primary hover:bg-primary-hover text-secondary" href={link}>{text}</a>
