<script lang="ts">
	// https://stackoverflow.com/questions/70927735/pass-variable-up-from-page-to-svelte-layout-via-slot
	import heroImage from '$lib/stores/ui/heroImage';
	import heroContainer from '$lib/stores/ui/heroContainer';

	import { FlexContainer } from '$lib/components';

	$: style =
		'absolute left-1/2 top-1/2 -translate-x-1/2 -translate-y-1/2 ' +
		$heroContainer.width +
		' ' +
		$heroContainer.height;
</script>

<!-- https://stackoverflow.com/questions/70805041/background-image-in-tailwindcss-using-dynamic-url-react-js -->
<div class="relative h-full bg-center bg-no-repeat bg-cover {$heroImage}">
	<FlexContainer bgColor={$heroContainer.color} styling={style}>
		<slot />
	</FlexContainer>
</div>

<style>
	:global(body),
	:global(html) {
		height: 100%;
	}

	:global(body) {
		margin: 0;
		padding: 0;
	}
</style>
