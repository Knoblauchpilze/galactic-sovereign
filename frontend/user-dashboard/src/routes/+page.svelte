<script lang="ts">
	import { FlexContainer, StyledLink, StyledTitle } from '$lib/components';
</script>

<FlexContainer>
	<StyledTitle text="User dashboard" />
	<StyledLink text="Login" link="/login" showAsButton={true} />
	<StyledLink text="Sign up" link="/signup" showAsButton={true} />
</FlexContainer>
