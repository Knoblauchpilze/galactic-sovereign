<div class="panel">
	<h1>User dashboard</h1>

	<div>
		<form>
			<label for="femail">e-mail:</label><br />
			<input type="text" id="femail" name="femail" /><br />
			<label for="fpassword">password:</label><br />
			<input type="text" id="fpassword" name="fpassword" /><br />
		</form>
	</div>

	<div>
		<button type="button">Login</button>
		<button type="button">Signup</button>
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
</style>
