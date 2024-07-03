<div class="hero-image">
	<div class="hero-panel">
		<slot></slot>
	</div>
</div>

<style>
	:global(body),
	:global(html) {
		height: 100%;
	}

	:global(body) {
		margin: 0;
		padding: 0;
	}

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
		top: 50%;
		left: 50%;
		transform: translate(-50%, -50%);

		width: 50%;
		height: 50%;

		display: flex;
		flex-direction: column;
		justify-content: space-around;

		background-color: #0005;
	}
</style>
