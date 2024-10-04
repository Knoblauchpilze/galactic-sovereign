<script lang="ts">
	import '$styles/app.css';

	import heroImage, { GAME_HERO_IMAGE } from '$lib/stores/ui/heroImage';
	import heroContainer, { GAME_HERO_CONTAINER_PROPS } from '$lib/stores/ui/heroContainer';
	import pageTitle, { HOMEPAGE_TITLE } from '$lib/stores/ui/pageTitle';

	import { CenteredWrapper, GamePageWrapper } from '$lib/components';

	import { mapPlanetResourcesToUiResources } from '$lib/game/resources';

	export let data;

	heroImage.set(GAME_HERO_IMAGE);
	heroContainer.set(GAME_HERO_CONTAINER_PROPS);
	$: title = HOMEPAGE_TITLE + ' - ' + data.planet.name;
	$: pageTitle.set(title);

	$: planetName = data.planet.name;
	$: universeName = data.universe.name;

	$: resources = mapPlanetResourcesToUiResources(
		data.planet.resources,
		data.planet.productions,
		data.planet.storages,
		data.resources
	);
</script>

<GamePageWrapper {universeName} {planetName} {resources}>
	<CenteredWrapper><p class="text-white">This is being built...</p></CenteredWrapper>
</GamePageWrapper>
