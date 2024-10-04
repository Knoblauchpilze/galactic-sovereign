<script lang="ts">
	import '$styles/app.css';
	import { CenteredWrapper, Footer, StyledLink } from '$lib/components';
</script>

<div class="w-[90%] h-[90%] flex justify-center items-stretch bg-transparent">
	<div class=" flex items-center">
		<ul class="list-none p-0 m-0">
			<li><a href="overview">Overview</a></li>
			<li><a href="buildings">Buildings</a></li>
		</ul>
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
