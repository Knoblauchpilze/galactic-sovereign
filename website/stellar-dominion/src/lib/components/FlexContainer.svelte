<script lang="ts">
	export let horizontal: boolean = false;
	export let center: boolean = false;
	export let extensible: boolean = true;
	export let bgColor: string = 'bg-transparent';

	$: direction = horizontal ? 'flex-row' : 'flex-col';
	$: space = center ? 'justify-center' : 'justify-around';
	$: grow = extensible ? 'flex-1' : 'flex-none';
</script>

<!-- https://stackoverflow.com/questions/29467660/how-to-stretch-children-to-fill-cross-axis -->
<div class="flex {direction} {space} items-center {bgColor} {grow}">
	<slot />
</div>
