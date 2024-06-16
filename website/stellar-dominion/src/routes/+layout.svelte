<div class="hero-image">
	<div class="hero-panel">
		<slot></slot>
	</div>
</div>

<style>
	/* https://stackoverflow.com/questions/19026884/flexbox-center-horizontally-and-vertically */
	/* https://stackoverflow.com/questions/71760177/styling-the-body-element-in-svelte */
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

		display: flex;
		flex-direction: column;
		justify-content: center;
		align-items: center;

		font-weight: bold;
		font-style: italic;

		transform: translate(-50%, -50%);
		background-color: #0005;
	}
</style>
