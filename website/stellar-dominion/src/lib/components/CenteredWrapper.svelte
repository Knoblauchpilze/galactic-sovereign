<script lang="ts">
	import { FlexContainer } from '$lib/components';

	export let width: string = 'w-full';
	export let height: string = 'h-full';

	export let bgColor: string = 'bg-transparent';
</script>

<!-- https://stackoverflow.com/questions/1122381/how-to-force-child-div-to-be-100-of-parent-divs-height-without-specifying-pare -->
<FlexContainer
	horizontal={true}
	justify={'center'}
	align={'stretch'}
	extensible={false}
	styling="{width} {height}"
>
	<FlexContainer {bgColor}>
		<!-- https://svelte.dev/tutorial/slots -->
		<slot />
	</FlexContainer>
</FlexContainer>
