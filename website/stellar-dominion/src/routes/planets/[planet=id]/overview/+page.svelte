<script lang="ts">
	import '$styles/app.css';
	import { CenteredWrapper, Header, StyledText, StyledTitle } from '$lib/components';

	import heroImage, { GAME_HERO_IMAGE } from '$lib/stores/ui/heroImage';
	import heroContainer, { GAME_HERO_CONTAINER_PROPS } from '$lib/stores/ui/heroContainer';

	import { mapPlanetResourcesToApiResources } from '$lib/resources';
	import { mapPlanetBuildingsToApiBuildings } from '$lib/buildings';

	// https://svelte.dev/blog/zero-config-type-safety
	export let data;

	const id: string = data.planet.id;
	const player: string = data.planet.player;
	const planetName: string = data.planet.name;
	const universeName: string = data.universe.name;

	// https://stackoverflow.com/questions/75616911/sveltekit-fetching-on-the-server-and-updating-the-writable-store
	heroImage.set(GAME_HERO_IMAGE);
	heroContainer.set(GAME_HERO_CONTAINER_PROPS);

	const resources = mapPlanetResourcesToApiResources(data.planet.resources, data.resources);
	const buildings = mapPlanetBuildingsToApiBuildings(data.planet.buildings, data.buildings);
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
		<!-- TODO: Interpret this with a new component -->
		<div class="flex justify-around bg-black">
			{#each buildings as building}
				<StyledText text="{building.name}: {building.level}" textColor="text-white" />
			{/each}
		</div>

		<div class="flex-grow">
			<StyledTitle text="Welcome to {planetName}!" />
			<StyledText text="Your id is {id}, and you are in the empire of {player}" />
			<StyledText text="This page will soon contain more information!" />
		</div>
	</div>
</CenteredWrapper>
