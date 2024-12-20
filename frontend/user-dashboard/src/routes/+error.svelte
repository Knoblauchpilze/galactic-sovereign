<script lang="ts">
	import { page } from '$app/state';
	import { FlexContainer, StyledLink, StyledTitle } from '$lib/components';

	let errorMessage: string = 'Unexpected error happened';
	if (page.error !== null) {
		errorMessage = page.error.message;
	}
</script>

<FlexContainer extensible={false} styling={'h-1/2'}>
	<StyledTitle text="Oh noes, something went wrong..." />

	<table>
		<tbody>
			<tr>
				<td class="text-secondary">Error message:</td>
				<td class="text-error">{errorMessage}</td>
			</tr>
		</tbody>
	</table>

	<StyledLink text="Back to safety" link="/" showAsButton={true} />
</FlexContainer>
