<script lang="ts">
	import '$styles/app.css';
	import { CenteredWrapper, Header, StyledText, StyledTitle } from '$lib/components';

	import heroImage from '$lib/stores/heroImage';
	import heroContainer from '$lib/stores/heroContainer';

	/** @type {import('./$types').PageData} */
	export let data;

	let id: string = data.planet.id;
	let player: string = data.planet.player;
	let name: string = data.planet.name;

	// https://stackoverflow.com/questions/75616911/sveltekit-fetching-on-the-server-and-updating-the-writable-store
	heroImage.set('bg-overview');
	heroContainer.set({
		width: 'w-full',
		height: 'h-full',
		color: 'bg-transparent'
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
			{#each data.planet.resources as resource}
				<StyledText text="{resource.id}: {resource.amount}" textColor="text-white" />
			{/each}
		</div>

		<div class="flex-grow">
			<StyledTitle text="Welcome to {name}!" />
			<StyledText text="Your id is {id}, and you are in the empire of {player}" />
			<StyledText text="This page will soon contain more information!" />
		</div>
	</div>
</CenteredWrapper>
