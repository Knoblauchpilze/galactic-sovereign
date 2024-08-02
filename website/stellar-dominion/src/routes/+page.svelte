<script lang="ts">
	import '$styles/app.css';
	import { CenteredWrapper, StyledLink, StyledText, StyledTitle } from '$lib/components';
</script>

<CenteredWrapper width="w-4/5" height="h-3/5">
	<StyledTitle text="Stellar Dominion" />
	<StyledText text="Battle for the galaxy" />

	<CenteredWrapper>
		<!-- https://svelte.dev/docs/component-directives -->
		<StyledLink text="Login" link="/login" showAsButton={true} />

		<StyledLink text="Sign up" link="/signup" showAsButton={true} />
	</CenteredWrapper>
</CenteredWrapper>
