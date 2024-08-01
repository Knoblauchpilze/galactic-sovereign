<script>
	import '../app.css';
	import { CenteredWrapper, StyledLink, StyledText, StyledTitle } from '$lib/components';
</script>

<CenteredWrapper width="w-4/5" height="h-3/5">
	<StyledTitle text="Stellar Dominion" />
	<StyledText text="Battle for the galaxy" />

	<CenteredWrapper>
		<!-- https://svelte.dev/docs/component-directives -->
		<StyledLink link="/login" text="Login" />

		<StyledLink link="/signup" text="Sign up" />
	</CenteredWrapper>
</CenteredWrapper>
