<script lang="ts">
	/** @type {import('./$types').PageData} */
	export let data;

	let id: string = data.planet.id;
	let player: string = data.planet.player;
	let name: string = data.planet.name;
</script>

<div class="wrapper">
	<h1>Welcome to {name}!</h1>
	<h2>Your id is {id}, and you are in the empire of {player}</h2>
	<h2>This page will soon contain more information!</h2>
</div>

<style>
	.wrapper {
		display: flex;
		flex-direction: column;
		justify-content: space-between;
		align-items: center;
	}

	h1, h2 {
		padding: 0.5em 1.5em;

		color: #b87333;
	}
</style>
