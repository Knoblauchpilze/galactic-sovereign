<script lang="ts">
	/** @type {import('./$types').ActionData} */
	export let form: HTMLFormElement;

	function resetFormError() {
		if (!form) {
			return;
		}
		form.message = '';
	}
</script>

<div class="wrapper">
	<h1>User dashboard</h1>
	<h2>Sign up</h2>

	<form method="POST" action="?/signup" class="form">
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
		<button class="action-button">Sign up</button>
	</form>

	{#if form?.message}
		<div class="error-details">
			Failed to sign up: {form?.message}
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

	h1,
	h2 {
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

	.error-details {
		color: #fe0075;
		font-size: 1.5em;

		position: fixed;
		bottom: 1em;
	}
</style>
