<script lang="ts">
	import heroImage, { GAME_HERO_IMAGE } from '$lib/stores/ui/heroImage';
	import heroContainer, { GAME_HERO_CONTAINER_PROPS } from '$lib/stores/ui/heroContainer';
	import pageTitle, { HOMEPAGE_TITLE } from '$lib/stores/ui/pageTitle';
	import activeScreen from '$lib/stores/activeScreen';

	import { BuildingAction, FlexContainer, GamePageWrapper, StyledTitle } from '$lib/components';

	import { invalidate } from '$app/navigation';

	import { formatDate } from '$lib/time';
	import { mapPlanetResourcesToUiResources } from '$lib/game/resources';
	import { mapBuildingActionsToUiActions } from '$lib/game/actions.js';

	export let data;

	heroImage.set(GAME_HERO_IMAGE);
	heroContainer.set(GAME_HERO_CONTAINER_PROPS);
	$: title = HOMEPAGE_TITLE + ' - ' + data.planet.name;
	$: pageTitle.set(title);
	$: activeScreen.set('overview');

	$: playerName = data.playerName;
	$: planetName = data.planet.name;
	$: universeName = data.universe.name;

	$: colonizationDate = formatDate(data.planet.createdAt);
	$: usedFields = data.planet.buildings.reduce((used, building) => used + building.level, 0);
	$: actions = mapBuildingActionsToUiActions(data.planet.buildingActions, data.buildings);

	$: resources = mapPlanetResourcesToUiResources(
		data.planet.resources,
		data.planet.productions,
		data.planet.storages,
		data.resources
	);

	function onActionCompleted() {
		invalidate('data:planet');
	}
</script>

<GamePageWrapper {universeName} {playerName} {planetName} {resources}>
	<FlexContainer align={'center'}>
		<StyledTitle text="Overview of {planetName}" />

		<FlexContainer justify={'center'} bgColor={'bg-overlay'}>
			<table>
				<tr>
					<td class="text-white">Colonization date:</td>
					<td class="text-white">{colonizationDate}</td>
				</tr>
				<tr>
					<td class="text-white">Used fields:</td>
					<td class="text-white">{usedFields}</td>
				</tr>
			</table>
		</FlexContainer>
	</FlexContainer>

	<FlexContainer align={'stretch'}>
		<StyledTitle text="Actions running on {planetName}" />
		<FlexContainer vertical={false} justify={'start'} align={'start'} styling={'flex-wrap'}>
			{#each actions as action}
				<BuildingAction {action} onCompleted={onActionCompleted} />
			{/each}
		</FlexContainer>
	</FlexContainer>
</GamePageWrapper>
