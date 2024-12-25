<script lang="ts">
	import {
		FlexContainer,
		StyledLink,
		StyledText,
		StyledTitle
	} from '@totocorpsoftwareinc/frontend-toolkit';

	import heroImage, { HOMEPAGE_HERO_IMAGE } from '$lib/stores/ui/heroImage';
	import heroContainer, { HOMEPAGE_HERO_CONTAINER_PROPS } from '$lib/stores/ui/heroContainer';

	heroImage.set(HOMEPAGE_HERO_IMAGE);
	heroContainer.set(HOMEPAGE_HERO_CONTAINER_PROPS);
</script>

<FlexContainer extensible={false} styling="h-3/5">
	<StyledTitle text="Galactic Sovereign" />
	<StyledText text="Battle for the galaxy" />

	<FlexContainer>
		<StyledLink text="Login" link="/login" showAsButton={true} />

		<StyledLink text="Sign up" link="/signup" showAsButton={true} />
	</FlexContainer>
</FlexContainer>
