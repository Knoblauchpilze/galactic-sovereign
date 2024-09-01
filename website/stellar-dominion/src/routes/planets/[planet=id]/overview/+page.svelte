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

	import { mapPlanetResourcesToUiResources } from '$lib/resources';
	import { mapPlanetBuildingsToUiBuildings } from '$lib/buildings';
	import { mapBuildingActionsToUiActions } from '$lib/actions.js';
	import { invalidate } from '$app/navigation';

	// https://svelte.dev/blog/zero-config-type-safety
	export let data;

	const planetName: string = data.planet.name;
	const universeName: string = data.universe.name;

	// https://stackoverflow.com/questions/75616911/sveltekit-fetching-on-the-server-and-updating-the-writable-store
	heroImage.set(GAME_HERO_IMAGE);
	heroContainer.set(GAME_HERO_CONTAINER_PROPS);
	const title = HOMEPAGE_TITLE + ' - ' + data.planet.name;
	pageTitle.set(title);

	const resources = mapPlanetResourcesToUiResources(data.planet.resources, data.resources);
	const buildings = mapPlanetBuildingsToUiBuildings(
		data.planet.id,
		data.planet.buildings,
		data.planet.buildingActions,
		data.buildings,
		data.resources
	);
	const actions = mapBuildingActionsToUiActions(data.planet.buildingActions, data.buildings);

	const anyBuildingActionRunning = data.planet.buildingActions.length !== 0;

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
		<div class="flex justify-around bg-black">
			{#each resources as resource}
				<StyledText text="{resource.name}: {resource.amount}" textColor="text-white" />
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
