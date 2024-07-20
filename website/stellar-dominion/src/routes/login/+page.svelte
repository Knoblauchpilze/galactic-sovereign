<script lang="ts">
	/** @type {import('./$types').ActionData} */
	export let form: HTMLFormElement;

	/** @type {import('./$types').PageData} */
	export let data;

	function resetFormError() {
		if (!form) {
			return;
		}
		form.message = '';
	}
</script>

<div class="wrapper">
	<div class="signup-navbar">
		Don't have an account yet? Click <a href="/signup">here</a> to sign-up!
	</div>

	<h1>Stellar Dominion</h1>
	<h2>Login</h2>

	<form method="POST" action="?/login" class="form">
		<div class="field">
			<label for="form-universe">Universe:</label>
			<select name="universe" id="universe">
				{#each data.universes as universe}
					<option value={universe.id}>{universe.name}</option>
				{/each}
			</select>
		</div>
		<div class="field">
			<label for="form-email">e-mail:</label>
			<input
				type="text"
				name="email"
				placeholder="Enter your email address"
				required
				value={form?.email ?? ''}
				on:input={resetFormError}
			/>
		</div>
		<div class="field">
			<label for="form-password">password:</label>
			<input
				type="text"
				name="password"
				placeholder="Enter your password"
				required
				on:input={resetFormError}
			/>
		</div>
		<button class="action-button">Login</button>
	</form>

	{#if form?.message}
		<div class="error-details">
			Failed to login: {form?.message}
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

	.signup-navbar {
		color: #b87333;

		position: fixed;
		left: 1em;
		top: 1em;
	}

	.signup-navbar a:link {
		color: #b87333;
	}

	.signup-navbar a:visited {
		color: #b87333;
	}

	.signup-navbar a:hover {
		color: #fff;
	}

	h1,
	h2 {
		color: #b87333;
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

		color: #b87333;
	}

	input {
		font-size: 1em;
	}

	.action-button {
		padding: 1em 3em;
		border-radius: 8px;
		border-width: 0;

		color: #b87333;
		background-color: #263037;
	}

	.action-button:hover {
		padding: 1em 3em;
		border-radius: 8px;

		color: #b87333;
		background-color: #36454f;
	}

	.error-details {
		color: #fe0075;
		font-size: 1.5em;

		position: fixed;
		bottom: 1em;
	}
</style>
