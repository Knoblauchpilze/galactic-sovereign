<script lang="ts">
	import { FlexContainer, StyledLink, StyledText, StyledTitle } from '$lib/components';

	import heroImage, { HOMEPAGE_HERO_IMAGE } from '$lib/stores/ui/heroImage';
	import heroContainer, { HOMEPAGE_HERO_CONTAINER_PROPS } from '$lib/stores/ui/heroContainer';

	heroImage.set(HOMEPAGE_HERO_IMAGE);
	heroContainer.set(HOMEPAGE_HERO_CONTAINER_PROPS);
</script>

<!-- <div class="flex flew-col justify-center items-center bg-yellow-200 flex-none h-3/5"> -->
<FlexContainer extensible={false} styling="h-3/5">
	<StyledTitle text="Stellar Dominion" />
	<StyledText text="Battle for the galaxy" />

	<FlexContainer>
		<!-- https://svelte.dev/docs/component-directives -->
		<StyledLink text="Login" link="/login" showAsButton={true} />

		<StyledLink text="Sign up" link="/signup" showAsButton={true} />
	</FlexContainer>
</FlexContainer>
