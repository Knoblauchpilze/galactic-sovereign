<script lang="ts">
	import '$styles/app.css';
	import { CenteredWrapper, Footer, StyledLink } from '$lib/components';

	import activeScreen from '$lib/stores/activeScreen';
	import FlexContainer from '$lib/components/FlexContainer.svelte';
</script>

<FlexContainer
	horizontal={true}
	extensible={false}
	center={true}
	stretch={true}
	styling={'w-[90%] h-[90%]'}
>
	<FlexContainer center={true} stretch={true} extensible={false}>
		<StyledLink
			text="Overview"
			link="overview"
			showAsButton={true}
			disabled={$activeScreen === 'overview'}
		/>
		<StyledLink
			text="Buildings"
			link="buildings"
			showAsButton={true}
			disabled={$activeScreen === 'buildings'}
		/>
	</FlexContainer>

	<slot />

	<Footer>
		<p>
			Meet the <StyledLink text="author" link="https://github.com/Knoblauchpilze" /> on github!
		</p>
		<p>
			Found an issue? Report it <StyledLink
				text="here"
				link="https://github.com/Knoblauchpilze/user-service/issues"
			/>
		</p>
	</Footer>
</FlexContainer>
