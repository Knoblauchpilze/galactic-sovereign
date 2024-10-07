<script lang="ts">
	export let horizontal: boolean = false;
	export let center: boolean = false;
	export let stretch: boolean = false;
	export let extensible: boolean = true;
	export let bgColor: string = 'bg-transparent';
	export let styling: string = '';

	$: direction = horizontal ? 'flex-row' : 'flex-col';
	$: justify = center ? 'justify-center' : 'justify-around';
	$: items = stretch ? 'items-stretch' : 'items-center';
	$: grow = extensible ? 'flex-1' : 'flex-none';
</script>

<!-- https://stackoverflow.com/questions/29467660/how-to-stretch-children-to-fill-cross-axis -->
<div class="flex {direction} {justify} {items} {bgColor} {grow} {styling}">
	<slot />
</div>
