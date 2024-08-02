<script lang="ts">
	import '$styles/app.css';
	import { CenteredWrapper, Footer, StyledLink } from '$lib/components';
</script>

<!-- https://stackoverflow.com/questions/70805041/background-image-in-tailwindcss-using-dynamic-url-react-js -->
<CenteredWrapper>
	<slot />
	<Footer>
		<p>Meet the <StyledLink text="author" link="https://github.com/Knoblauchpilze" /> on github!</p>
		<p>
			Found an issue? Report it <StyledLink
				text="here"
				link="https://github.com/Knoblauchpilze/user-service/issues"
			/>
		</p>
	</Footer>
</CenteredWrapper>
