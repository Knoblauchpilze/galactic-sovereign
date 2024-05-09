<script lang="ts">
	// https://learn.svelte.dev/tutorial/lib
	import ResponseEnvelope from '$lib/responseEnvelope';
	import login from '$lib/api';

	// https://stackoverflow.com/questions/73280092/capture-value-from-an-input-with-svelte
	let email: string;
	let password: string;

	let loginResponse: ResponseEnvelope;

	async function performLogin() {
		loginResponse = await login(email, password);
	}
</script>

<div class="panel">
	<h1>User dashboard</h1>

	<form class="form">
		<label for="form-email">e-mail:</label>
		<input type="text" name="email" bind:value={email} />
		<label for="form-password"> password:</label>
		<input type="text" name="password" bind:value={password} />

		<div class="form-submit">
			<button on:click={performLogin}>Login</button>
		</div>
	</form>

	<div>
		{JSON.stringify(loginResponse)}
	</div>
</div>

<!-- https://stackoverflow.com/questions/71760177/styling-the-body-element-in-svelte -->
<!-- https://stackoverflow.com/questions/19026884/flexbox-center-horizontally-and-vertically -->
<style>
	:global(html),
	:global(body) {
		height: 100%;
	}

	:global(body) {
		margin: 0;
		padding: 0;
	}

	.panel {
		height: 100%;
		display: flex;
		flex-direction: column;
		justify-content: center;
		align-items: center;
	}

	.form {
		display: flex;
		flex-direction: column;
	}

	.form-submit {
		display: flex;
		flex-direction: column;
		align-items: center;
	}
</style>
