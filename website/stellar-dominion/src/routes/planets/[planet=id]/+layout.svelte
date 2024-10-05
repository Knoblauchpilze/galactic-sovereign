<script lang="ts">
	import '$styles/app.css';
	import { CenteredWrapper, Footer, StyledLink } from '$lib/components';

	import activeScreen from '$lib/stores/activeScreen';
</script>

<div class="w-[90%] h-[90%] flex justify-center items-stretch bg-transparent">
	<div class="flex flex-col justify-center items-stretch">
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
	</div>

	<CenteredWrapper>
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
	</CenteredWrapper>
</div>
