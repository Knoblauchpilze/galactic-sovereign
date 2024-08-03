<script lang="ts">
	import '$styles/app.css';
	import { CenteredWrapper, Header, StyledText, StyledTitle } from '$lib/components';

	import heroImage, { GAME_HERO_IMAGE } from '$lib/stores/ui/heroImage';
	import heroContainer, { GAME_HERO_CONTAINER_PROPS } from '$lib/stores/ui/heroContainer';

	// https://svelte.dev/blog/zero-config-type-safety
	export let data;

	let id: string = data.planet.id;
	let player: string = data.planet.player;
	let name: string = data.planet.name;

	// https://stackoverflow.com/questions/75616911/sveltekit-fetching-on-the-server-and-updating-the-writable-store
	heroImage.set(GAME_HERO_IMAGE);
	heroContainer.set(GAME_HERO_CONTAINER_PROPS);

	let resources = data.resources.map((apiResource) => {
		const maybeResource = data.planet.resources.find((r) => r.id === apiResource.id);
		if (maybeResource === undefined) {
			return {
				name: apiResource.name,
				amount: 0
			};
		} else {
			return {
				name: apiResource.name,
				amount: maybeResource.amount
			};
		}
	});
</script>

<CenteredWrapper width="w-4/5" height="h-4/5" bgColor="bg-overlay">
	<Header>
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

		<div class="flex-grow">
			<StyledTitle text="Welcome to {name}!" />
			<StyledText text="Your id is {id}, and you are in the empire of {player}" />
			<StyledText text="This page will soon contain more information!" />
		</div>
	</div>
</CenteredWrapper>
