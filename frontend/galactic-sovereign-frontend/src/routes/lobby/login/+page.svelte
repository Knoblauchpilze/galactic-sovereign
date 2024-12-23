<script lang="ts">
	import { FlexContainer, StyledLink, StyledTitle } from '@totocorpsoftwareinc/frontend-toolkit';
	import { SessionCard, StyledError, StyledText } from '$lib/components';

	import heroImage, { HOMEPAGE_HERO_IMAGE } from '$lib/stores/ui/heroImage';
	import heroContainer, { HOMEPAGE_HERO_CONTAINER_PROPS } from '$lib/stores/ui/heroContainer';

	let { data } = $props();

	heroImage.set(HOMEPAGE_HERO_IMAGE);
	heroContainer.set(HOMEPAGE_HERO_CONTAINER_PROPS);
</script>

<FlexContainer>
	<div class="fixed right-4 top-4">
		<p class="text-secondary">
			Back to the <StyledLink text="lobby" link="/lobby" />
		</p>
	</div>

	<FlexContainer extensible={false} styling="h-1/5">
		<StyledTitle text="Galactic Sovereign" />
		<StyledText text="Resume existing session" />
	</FlexContainer>

	<FlexContainer extensible={false} styling="h-3/5">
		{#if data.players.length > 0}
			<FlexContainer justify={'center'} align={'stretch'}>
				{#each data.players as player}
					<SessionCard {player} />
				{/each}
			</FlexContainer>
		{:else}
			<StyledError text="You don't have an account yet, please register!" />
			<StyledLink text="Register" link="/lobby/register" showAsButton={true} />
		{/if}
	</FlexContainer>
</FlexContainer>
