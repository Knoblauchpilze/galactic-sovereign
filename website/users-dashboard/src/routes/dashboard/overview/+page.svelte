<script lang="ts">
	/** @type {import('./$types').PageData} */
	export let data;

	let id: string = data.id;
	let email: string = data.email;
	let password: string = data.password;
	let apiKey: string = data.apiKey;

	// https://stackoverflow.com/questions/3552461/how-do-i-format-a-date-in-javascript
	const options: Intl.DateTimeFormatOptions = {
		year: 'numeric',
		month: 'long',
		day: 'numeric',
		hour: 'numeric',
		minute: 'numeric',
		second: 'numeric',
		hour12: false
	};
	const createdAt = data.createdAt.toLocaleDateString('en-US', options);

	/** @type {import('./$types').ActionData} */
	export let form: HTMLFormElement;
</script>

<div class="wrapper">
	<h1>User details</h1>

	<table>
		<tr>
			<td class="label">ID:</td>
			<td class="field">{id}</td>
		</tr>
		<tr>
			<td class="label">E-mail:</td>
			<td class="field">{email}</td>
		</tr>
		<tr>
			<td class="label">Password:</td>
			<td class="field">{password}</td>
		</tr>
		<tr>
			<td class="label">API key:</td>
			<td class="field">{apiKey}</td>
		</tr>
		<tr>
			<td class="label">Member since:</td>
			<td class="field">{createdAt}</td>
		</tr>
	</table>

	<form method="POST" action="?/logout">
		<button class="action-button">Logout</button>
	</form>

	{#if form?.message}
		<div class="error-details">
			Failed to logout: {form?.message}
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
	.label {
		padding: 0.5em 1.5em;

		color: #1eb854;
	}

	.field {
		color: #1db8ab;
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
