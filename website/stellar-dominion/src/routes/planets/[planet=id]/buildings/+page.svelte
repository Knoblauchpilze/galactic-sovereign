<script lang="ts">
	import '$styles/app.css';

	import heroImage, { GAME_HERO_IMAGE } from '$lib/stores/ui/heroImage';
	import heroContainer, { GAME_HERO_CONTAINER_PROPS } from '$lib/stores/ui/heroContainer';
	import pageTitle, { HOMEPAGE_TITLE } from '$lib/stores/ui/pageTitle';
	import activeScreen from '$lib/stores/activeScreen';

	import {
		FlexContainer,
		GamePageWrapper,
		StyledTitle,
		Building,
		BuildingAction
	} from '$lib/components';

	import { invalidate } from '$app/navigation';

	import { mapPlanetResourcesToUiResources } from '$lib/game/resources';
	import { mapPlanetBuildingsToUiBuildings } from '$lib/game/buildings';
	import { mapBuildingActionsToUiActions } from '$lib/game/actions.js';

	// https://svelte.dev/blog/zero-config-type-safety
	export let data;

	// https://stackoverflow.com/questions/77047087/why-is-page-svelte-not-reloading-after-page-server-ts-executes-load-function
	$: playerName = data.playerName;
	$: planetName = data.planet.name;
	$: universeName = data.universe.name;

	// https://stackoverflow.com/questions/75616911/sveltekit-fetching-on-the-server-and-updating-the-writable-store
	heroImage.set(GAME_HERO_IMAGE);
	heroContainer.set(GAME_HERO_CONTAINER_PROPS);
	$: title = HOMEPAGE_TITLE + ' - ' + planetName;
	$: pageTitle.set(title);
	$: activeScreen.set('buildings');

	$: resources = mapPlanetResourcesToUiResources(
		data.planet.resources,
		data.planet.productions,
		data.planet.storages,
		data.resources
	);
	$: buildings = mapPlanetBuildingsToUiBuildings(
		data.planet.id,
		data.planet.buildings,
		data.planet.buildingActions,
		data.buildings,
		data.resources
	);
	$: actions = mapBuildingActionsToUiActions(data.planet.buildingActions, data.buildings);

	$: anyBuildingActionRunning = data.planet.buildingActions.length !== 0;

	function onActionCompleted() {
		invalidate('data:planet');
	}
</script>

<GamePageWrapper {universeName} {playerName} {planetName} {resources}>
	<FlexContainer align={'stretch'}>
		<StyledTitle text="Buildings on {planetName}" />
		<!-- https://tailwindcss.com/docs/align-items -->
		<FlexContainer horizontal={true} justify={'start'} align={'start'} styling={'flex-wrap'}>
			{#each buildings as building}
				<Building
					{building}
					availableResources={resources}
					buildingActionAlreadyRunning={anyBuildingActionRunning}
				/>
			{/each}
		</FlexContainer>
	</FlexContainer>

	<FlexContainer align={'stretch'}>
		<StyledTitle text="Actions running on {planetName}" />
		<FlexContainer horizontal={true} justify={'start'} align={'start'} styling={'flex-wrap'}>
			{#each actions as action}
				<BuildingAction {action} onCompleted={onActionCompleted} />
			{/each}
		</FlexContainer>
	</FlexContainer>
</GamePageWrapper>
