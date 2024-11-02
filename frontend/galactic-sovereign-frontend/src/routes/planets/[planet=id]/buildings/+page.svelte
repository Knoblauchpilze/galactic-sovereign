<script lang="ts">
	import heroImage, { GAME_HERO_IMAGE } from '$lib/stores/ui/heroImage';
	import heroContainer, { GAME_HERO_CONTAINER_PROPS } from '$lib/stores/ui/heroContainer';
	import pageTitle from '$lib/stores/ui/pageTitle';
	import activeScreen from '$lib/stores/activeScreen';

	import {
		FlexContainer,
		GamePageWrapper,
		StyledTitle,
		Building,
		BuildingAction
	} from '$lib/components';

	import { invalidate } from '$app/navigation';

	let { data } = $props();

	// https://stackoverflow.com/questions/75616911/sveltekit-fetching-on-the-server-and-updating-the-writable-store
	heroImage.set(GAME_HERO_IMAGE);
	heroContainer.set(GAME_HERO_CONTAINER_PROPS);
	pageTitle.set(data.wepageTitle);
	activeScreen.set('buildings');

	let anyBuildingActionRunning = $derived(data.buildingActions.length !== 0);

	function onActionCompleted() {
		invalidate('data:planet');
	}
</script>

<GamePageWrapper
	universeName={data.universeName}
	playerName={data.playerName}
	planetName={data.planetName}
	resources={data.resources}
>
	<FlexContainer align={'stretch'}>
		<StyledTitle text="Buildings on {data.planetName}" />
		<!-- https://tailwindcss.com/docs/align-items -->
		<FlexContainer vertical={false} justify={'start'} align={'start'} styling={'flex-wrap'}>
			{#each data.buildings as building}
				<Building
					{building}
					availableResources={data.resources}
					buildingActionAlreadyRunning={anyBuildingActionRunning}
				/>
			{/each}
		</FlexContainer>
	</FlexContainer>

	<FlexContainer align={'stretch'}>
		<StyledTitle text="Actions running on {data.planetName}" />
		<FlexContainer vertical={false} justify={'start'} align={'start'} styling={'flex-wrap'}>
			{#each data.buildingActions as buildingAction}
				<BuildingAction action={buildingAction} onCompleted={onActionCompleted} />
			{/each}
		</FlexContainer>
	</FlexContainer>
</GamePageWrapper>
