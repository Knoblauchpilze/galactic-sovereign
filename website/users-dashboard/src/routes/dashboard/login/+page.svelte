<script lang="ts">
	// https://learn.svelte.dev/tutorial/lib
	import ResponseEnvelope from '$lib/responseEnvelope';
	import login from '$lib/api';

	// https://stackoverflow.com/questions/73280092/capture-value-from-an-input-with-svelte
	let email: string;
	let password: string;

	let loginError: string = '';
	let loginResponse: ResponseEnvelope;

	async function performLogin() {
		loginError = '';
		loginResponse = await login(email, password);

		if (loginResponse.error()) {
			loginError = String(loginResponse.Details);
		}
	}
</script>

<div class="wrapper">
	<h1>User dashboard</h1>

	<form class="form">
		<div class="field">
			<label for="form-email">e-mail:</label>
			<!-- https://stackoverflow.com/questions/62278480/add-onchange-handler-to-input-in-svelte -->
			<input
				type="text"
				name="email"
				placeholder="Enter your email address"
				bind:value={email}
				on:input={() => (loginError = '')}
			/>
		</div>
		<div class="field">
			<label for="form-password"> password:</label>
			<input
				type="text"
				name="password"
				placeholder="Enter your password"
				bind:value={password}
				on:input={() => (loginError = '')}
			/>
		</div>
	</form>

	<button class="action-button" on:click={performLogin}>Login</button>

	{#if loginError !== ''}
		<div class="error-details">
			Failed to login: {loginError}
		</div>
	{/if}
</div>

<style>
	.wrapper {
		display: flex;
		flex-direction: column;
		justify-content: space-between;
		align-items: center;
	}

	h1 {
		color: #1eb854;
	}

	.form {
		display: flex;
		flex-direction: column;
	}

	.field {
		display: flex;
		flex-direction: column;

		font-size: 1.5em;
		padding: 1em 0em;

		color: #1eb854;
	}

	input {
		font-size: 1em;
	}

	.action-button {
		padding: 1em 3em;
		border-radius: 8px;

		background-color: #1eb854;
	}

	/* https://www.w3schools.com/howto/howto_js_snackbar.asp */
	.error-details {
		color: #fe0075;
		font-size: 1.5em;

		position: fixed;
		bottom: 1em;
	}
</style>
