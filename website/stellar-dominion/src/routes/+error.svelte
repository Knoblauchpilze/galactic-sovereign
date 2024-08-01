<script lang="ts">
	import '../app.css';
	import { page } from '$app/stores';
	import { CenteredWrapper, StyledLink, StyledTitle } from '$lib/components';

	let errorMessage: string = 'Unexpected error happened';
	if ($page.error !== null) {
		errorMessage = $page.error.message;
	}
</script>

<CenteredWrapper height="h-1/2">
	<StyledTitle text="Oh noes, something went wrong..." />

	<table>
		<tr>
			<td class="text-secondary">Error message:</td>
			<td class="text-error">{errorMessage}</td>
		</tr>
	</table>

	<StyledLink link="/" text="Back to safety" />
</CenteredWrapper>
