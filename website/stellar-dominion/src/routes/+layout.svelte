<div class="hero-image">
	<div class="hero-panel">
		<div class="main-panel">
			<h1>Stellar dominion</h1>
			<p>Battle for the galaxy</p>
			<slot></slot>
		</div>
	</div>
</div>

<style>
	:global(html),
	:global(body) {
		height: 100%;
	}

	:global(body) {
		margin: 0;
		padding: 0;
	}

	/* https://kit.svelte.dev/docs/images */
	.hero-image {
		background-image: url('$lib/assets/background.webp');

		height: 100%;

		background-position: center;
		background-repeat: no-repeat;
		background-size: cover;
		position: relative;
	}

	.hero-panel {
		text-align: center;
		position: absolute;
		padding: 2em;
		top: 50%;
		left: 50%;

		transform: translate(-50%, -50%);
		background-color: #0005;
	}

	h1,
	p {
		font-weight: bold;
		font-style: italic;

		color: #b87333;
	}
</style>
