<script lang="ts">
	import { page } from '$app/stores';

	let errorMessage: string = 'Unexpected error happened';
	if ($page.error !== null) {
		errorMessage = $page.error.message;
	}
</script>

<div class="wrapper">
	<h1>Oh noes, something went wrong...</h1>

	<table>
		<tr>
			<td class="label">Error message:</td>
			<td class="field">{errorMessage}</td>
		</tr>
	</table>

	<a href="/">Back to safety</a>
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

		color: #b87333;
	}

	.field {
		color: #d92d0f;
	}

	a {
		padding: 1em 3em;

		color: #b87333;
	}
</style>
