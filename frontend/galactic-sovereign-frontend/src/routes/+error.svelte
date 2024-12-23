<script lang="ts">
	import { page } from '$app/state';
	import { StyledLink, StyledTitle, FlexContainer } from '@totocorpsoftwareinc/frontend-toolkit';

	import heroImage, { HOMEPAGE_HERO_IMAGE } from '$lib/stores/ui/heroImage';
	import heroContainer, { HOMEPAGE_HERO_CONTAINER_PROPS } from '$lib/stores/ui/heroContainer';

	let errorMessage: string = 'Unexpected error happened';
	if (page.error !== null) {
		errorMessage = page.error.message;
	}

	heroImage.set(HOMEPAGE_HERO_IMAGE);
	heroContainer.set(HOMEPAGE_HERO_CONTAINER_PROPS);
</script>

<FlexContainer extensible={false} styling={'h-1/2'}>
	<StyledTitle text="Oh noes, something went wrong..." />

	<table>
		<tbody>
			<tr>
				<td class="text-secondary">Error message:</td>
				<td class="text-error">{errorMessage}</td>
			</tr>
		</tbody>
	</table>

	<StyledLink text="Back to safety" link="/" showAsButton={true} />
</FlexContainer>
