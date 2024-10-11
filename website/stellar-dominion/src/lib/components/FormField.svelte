<script lang="ts">
	import { FlexContainer } from '$lib/components';

	export let label: string = 'My label';
	export let labelId: string;
</script>

<FlexContainer extensible={false} align={'stretch'} styling={'text-secondary'}>
	<label for={labelId}>{label}</label>
	<slot />
</FlexContainer>
