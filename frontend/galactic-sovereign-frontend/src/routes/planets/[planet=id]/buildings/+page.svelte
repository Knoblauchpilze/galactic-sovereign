<script lang="ts">
	import heroImage, { GAME_HERO_IMAGE } from '$lib/stores/ui/heroImage';
	import heroContainer, { GAME_HERO_CONTAINER_PROPS } from '$lib/stores/ui/heroContainer';
	import pageTitle from '$lib/stores/ui/pageTitle';
	import activeScreen from '$lib/stores/activeScreen';

	import { type UiResource } from '$lib/game/resources';
	import { type UiBuilding } from '$lib/game/buildings';
	import { type UiBuildingAction } from '$lib/game/actions';

	import {
		FlexContainer,
		GamePageWrapper,
		StyledTitle,
		Building,
		BuildingAction
	} from '$lib/components';

	import { invalidate } from '$app/navigation';

	interface Props {
		wepageTitle: string;
		universeName: string;
		playerName: string;
		planetName: string;
		resources: UiResource[];
		buildings: UiBuilding[];
		buildingActions: UiBuildingAction[];
	}

	// https://svelte.dev/blog/zero-config-type-safety
	let {
		wepageTitle,
		universeName,
		playerName,
		planetName,
		resources,
		buildings,
		buildingActions
	}: Props = $props();

	// https://stackoverflow.com/questions/75616911/sveltekit-fetching-on-the-server-and-updating-the-writable-store
	heroImage.set(GAME_HERO_IMAGE);
	heroContainer.set(GAME_HERO_CONTAINER_PROPS);
	pageTitle.set(wepageTitle);
	activeScreen.set('buildings');

	let anyBuildingActionRunning = $derived(buildingActions.length !== 0);

	function onActionCompleted() {
		invalidate('data:planet');
	}
</script>

<GamePageWrapper {universeName} {playerName} {planetName} {resources}>
	<FlexContainer align={'stretch'}>
		<StyledTitle text="Buildings on {planetName}" />
		<!-- https://tailwindcss.com/docs/align-items -->
		<FlexContainer vertical={false} justify={'start'} align={'start'} styling={'flex-wrap'}>
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
		<FlexContainer vertical={false} justify={'start'} align={'start'} styling={'flex-wrap'}>
			{#each buildingActions as buildingAction}
				<BuildingAction action={buildingAction} onCompleted={onActionCompleted} />
			{/each}
		</FlexContainer>
	</FlexContainer>
</GamePageWrapper>
