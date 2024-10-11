<script lang="ts">
	import heroImage, { GAME_HERO_IMAGE } from '$lib/stores/ui/heroImage';
	import heroContainer, { GAME_HERO_CONTAINER_PROPS } from '$lib/stores/ui/heroContainer';
	import pageTitle, { HOMEPAGE_TITLE } from '$lib/stores/ui/pageTitle';
	import activeScreen from '$lib/stores/activeScreen';

	import { FlexContainer, GamePageWrapper } from '$lib/components';

	import { mapPlanetResourcesToUiResources } from '$lib/game/resources';

	export let data;

	heroImage.set(GAME_HERO_IMAGE);
	heroContainer.set(GAME_HERO_CONTAINER_PROPS);
	$: title = HOMEPAGE_TITLE + ' - ' + data.planet.name;
	$: pageTitle.set(title);
	$: activeScreen.set('overview');

	$: playerName = data.playerName;
	$: planetName = data.planet.name;
	$: universeName = data.universe.name;

	$: resources = mapPlanetResourcesToUiResources(
		data.planet.resources,
		data.planet.productions,
		data.planet.storages,
		data.resources
	);
</script>

<GamePageWrapper {universeName} {playerName} {planetName} {resources}>
	<FlexContainer>
		<p class="text-white">This is being built...</p>
	</FlexContainer>
</GamePageWrapper>
