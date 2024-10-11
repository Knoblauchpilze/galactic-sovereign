<script lang="ts">
	import '$styles/app.css';

	import pageTitle from '$lib/stores/ui/pageTitle';

	import { HeroContainer } from '$lib/components';
</script>

<!-- https://stackoverflow.com/questions/74243719/how-to-change-page-title-dynamically-in-sveltekit -->
<svelte:head>
	<title>{$pageTitle}</title>
</svelte:head>
<!-- https://stackoverflow.com/questions/70805041/background-image-in-tailwindcss-using-dynamic-url-react-js -->
<HeroContainer>
	<slot />
</HeroContainer>
