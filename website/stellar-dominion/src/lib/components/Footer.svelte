<script lang="ts">
	import { FlexContainer } from '$lib/components';

	export let height: string = 'h-1/10';
</script>

<FlexContainer
	horizontal={true}
	bgColor={'bg-black'}
	styling="fixed bottom-0 w-full {height} text-white"
>
	<slot />
</FlexContainer>
