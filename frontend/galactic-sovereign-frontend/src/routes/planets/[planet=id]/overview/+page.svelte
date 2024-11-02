<script lang="ts">
	import heroImage, { GAME_HERO_IMAGE } from '$lib/stores/ui/heroImage';
	import heroContainer, { GAME_HERO_CONTAINER_PROPS } from '$lib/stores/ui/heroContainer';
	import pageTitle from '$lib/stores/ui/pageTitle';
	import activeScreen from '$lib/stores/activeScreen';

	import { BuildingAction, FlexContainer, GamePageWrapper, StyledTitle } from '$lib/components';

	import { invalidate } from '$app/navigation';

	import { formatDate } from '$lib/time';

	// https://svelte.dev/blog/zero-config-type-safety
	// https://svelte.dev/docs/kit/load#Page-data
	let { data } = $props();

	heroImage.set(GAME_HERO_IMAGE);
	heroContainer.set(GAME_HERO_CONTAINER_PROPS);
	pageTitle.set(data.wepageTitle);
	activeScreen.set('overview');

	let colonizationDate = $derived(formatDate(data.planetCreationTime));
	let usedFields = $derived(data.buildings.reduce((used, building) => used + building.level, 0));

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
	<FlexContainer align={'center'}>
		<StyledTitle text="Overview of {data.planetName}" />

		<FlexContainer justify={'center'} bgColor={'bg-overlay'}>
			<table>
				<tbody>
					<tr>
						<td class="text-white">Colonization date:</td>
						<td class="text-white">{colonizationDate}</td>
					</tr>
					<tr>
						<td class="text-white">Used fields:</td>
						<td class="text-white">{usedFields}</td>
					</tr>
				</tbody>
			</table>
		</FlexContainer>
	</FlexContainer>

	<FlexContainer align={'stretch'}>
		<StyledTitle text="Actions running on {data.planetName}" />
		<FlexContainer vertical={false} justify={'start'} align={'start'} styling={'flex-wrap'}>
			{#each data.buildingActions as action}
				<BuildingAction {action} onCompleted={onActionCompleted} />
			{/each}
		</FlexContainer>
	</FlexContainer>
</GamePageWrapper>
