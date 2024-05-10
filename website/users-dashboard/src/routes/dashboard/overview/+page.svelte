<script lang="ts">
	let id: string = '00-00-00-00';
	let email: string = 'some-user';
	let password: string = 'some-password';

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
	const createdAt = new Date().toLocaleDateString('en-US', options);

	async function performLogout() {}
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
			<td class="label">Member since:</td>
			<td class="field">{createdAt}</td>
		</tr>
	</table>

	<button class="action-button" on:click={performLogout}>Logout</button>
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
</style>
