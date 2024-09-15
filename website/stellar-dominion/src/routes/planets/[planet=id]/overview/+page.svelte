<script lang="ts">
	import '$styles/app.css';
	import {
		CenteredWrapper,
		Header,
		StyledText,
		StyledTitle,
		Building,
		BuildingAction
	} from '$lib/components';

	import heroImage, { GAME_HERO_IMAGE } from '$lib/stores/ui/heroImage';
	import heroContainer, { GAME_HERO_CONTAINER_PROPS } from '$lib/stores/ui/heroContainer';
	import pageTitle, { HOMEPAGE_TITLE } from '$lib/stores/ui/pageTitle';

	import { mapPlanetResourcesToUiResources } from '$lib/game/resources';
	import { mapPlanetBuildingsToUiBuildings } from '$lib/game/buildings';
	import { mapBuildingActionsToUiActions } from '$lib/game/actions.js';
	import { roundToInteger } from '$lib/displayUtils';
	import { invalidate } from '$app/navigation';

	// https://svelte.dev/blog/zero-config-type-safety
	export let data;

	// https://stackoverflow.com/questions/77047087/why-is-page-svelte-not-reloading-after-page-server-ts-executes-load-function
	const planetName = data.planet.name;
	const universeName = data.universe.name;

	// https://stackoverflow.com/questions/75616911/sveltekit-fetching-on-the-server-and-updating-the-writable-store
	heroImage.set(GAME_HERO_IMAGE);
	heroContainer.set(GAME_HERO_CONTAINER_PROPS);
	const title = HOMEPAGE_TITLE + ' - ' + data.planet.name;
	pageTitle.set(title);

	$: resources = mapPlanetResourcesToUiResources(data.planet.resources, data.resources);
	$: buildings = mapPlanetBuildingsToUiBuildings(
		data.planet.id,
		data.planet.buildings,
		data.planet.buildingActions,
		data.buildings,
		data.resources
	);
	$: actions = mapBuildingActionsToUiActions(data.planet.buildingActions, data.buildings);

	$: anyBuildingActionRunning = data.planet.buildingActions.length !== 0;


	$: console.log("data: ", JSON.stringify(data));

	function onActionCompleted() {
		invalidate('data:planet');
	}
</script>

<CenteredWrapper width="w-4/5" height="h-4/5" bgColor="bg-overlay">
	<Header>
		<StyledText text={universeName} textColor="text-white" />
		<StyledText text={planetName} textColor="text-white" />
		<form method="POST" action="?/logout">
			<button class="hover:underline">Logout</button>
		</form>
	</Header>

	<div class="flex flex-col justify-start flex-grow w-full">
		<div class="flex justify-around bg-black justify-items-stretch">
			{#each resources as resource}
				<div class="flex space-between">
					<StyledText
						text="{resource.name}: {roundToInteger(resource.amount)}"
						textColor="text-white"
						styling="px-1"
					/>
					<StyledText text="(+{roundToInteger(resource.production)}/h)" textColor="text-enabled" />
				</div>
			{/each}
		</div>

		<CenteredWrapper>
			<StyledTitle text="Buildings on {planetName}" />
			<!-- https://tailwindcss.com/docs/align-items -->
			<div class="w-full h-full flex flex-wrap items-start bg-transparent">
				{#each buildings as building}
					<Building
						{building}
						availableResources={resources}
						buildingActionAlreadyRunning={anyBuildingActionRunning}
					/>
				{/each}
			</div>
		</CenteredWrapper>

		<CenteredWrapper>
			<StyledTitle text="Actions running on {planetName}" />
			<div class="w-full h-full flex flex-wrap items-start bg-transparent">
				{#each actions as action}
					<BuildingAction {action} onCompleted={onActionCompleted} />
				{/each}
			</div>
		</CenteredWrapper>
	</div>
</CenteredWrapper>
