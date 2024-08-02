<script lang="ts">
	import '$styles/app.css';
	import { HeroContainer } from '$lib/components';
</script>

<!-- https://stackoverflow.com/questions/70805041/background-image-in-tailwindcss-using-dynamic-url-react-js -->
<HeroContainer>
	<slot/>
</HeroContainer>