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

	<a href="/dashboard/login">Login</a>
	<a href="/dashboard/signup">Sign up</a>
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

	a {
		padding: 1em 3em;

		color: #1eb854;
	}
</style>
