<script lang="ts">
	import '$styles/app.css';
	import { page } from '$app/stores';
	import { CenteredWrapper, StyledLink, StyledTitle } from '$lib/components';

	import heroImage, { HOMEPAGE_HERO_IMAGE } from '$lib/stores/ui/heroImage';
	import heroContainer, { HOMEPAGE_HERO_CONTAINER_PROPS } from '$lib/stores/ui/heroContainer';

	let errorMessage: string = 'Unexpected error happened';
	if ($page.error !== null) {
		errorMessage = $page.error.message;
	}

	heroImage.set(HOMEPAGE_HERO_IMAGE);
	heroContainer.set(HOMEPAGE_HERO_CONTAINER_PROPS);
</script>

<CenteredWrapper height="h-1/2">
	<StyledTitle text="Oh noes, something went wrong..." />

	<table>
		<tr>
			<td class="text-secondary">Error message:</td>
			<td class="text-error">{errorMessage}</td>
		</tr>
	</table>

	<StyledLink text="Back to safety" link="/" showAsButton={true} />
</CenteredWrapper>
