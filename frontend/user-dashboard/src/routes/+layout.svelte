<div class="root">
	<div class="panel">
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

	.root {
		height: 100%;
		display: flex;
		flex-direction: column;
		justify-content: center;
		align-items: center;

		/* https://svelte-ux.techniq.dev/theme use dark theme: forest */
		background-color: #19362d;
	}

	.panel {
		width: 65%;
		height: 65%;
	}
</style>
