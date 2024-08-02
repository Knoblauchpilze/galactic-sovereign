<script lang="ts">
	import '$styles/app.css';
	import { CenteredWrapper, StyledText, StyledTitle } from '$lib/components';

	import heroImage from '$lib/heroImage';

	/** @type {import('./$types').PageData} */
	export let data;

	let id: string = data.planet.id;
	let player: string = data.planet.player;
	let name: string = data.planet.name;

	// https://stackoverflow.com/questions/75616911/sveltekit-fetching-on-the-server-and-updating-the-writable-store
	heroImage.set('bg-overview');
</script>

<CenteredWrapper>
	<StyledTitle text="Welcome to {name}!" />
	<StyledText text="Your id is {id}, and you are in the empire of {player}" />
	<StyledText text="This page will soon contain more information!" />
</CenteredWrapper>
