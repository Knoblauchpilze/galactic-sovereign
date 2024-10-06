<script lang="ts">
	import '$styles/app.css';
	import {
		CenteredWrapper,
		FlexContainer,
		StyledLink,
		StyledText,
		StyledTitle
	} from '$lib/components';

	import heroImage, { HOMEPAGE_HERO_IMAGE } from '$lib/stores/ui/heroImage';
	import heroContainer, { HOMEPAGE_HERO_CONTAINER_PROPS } from '$lib/stores/ui/heroContainer';

	heroImage.set(HOMEPAGE_HERO_IMAGE);
	heroContainer.set(HOMEPAGE_HERO_CONTAINER_PROPS);
</script>

<CenteredWrapper width="w-4/5" height="h-3/5">
	<StyledTitle text="Stellar Dominion" />
	<StyledText text="Battle for the galaxy" />

	<FlexContainer>
		<!-- https://svelte.dev/docs/component-directives -->
		<StyledLink text="Login" link="/login" showAsButton={true} />

		<StyledLink text="Sign up" link="/signup" showAsButton={true} />
	</FlexContainer>
</CenteredWrapper>
