<script lang="ts">
	export let horizontal: boolean = false;
	export let justify: string = 'around';
	export let align: string = 'center';
	export let extensible: boolean = true;
	export let bgColor: string = 'bg-transparent';
	export let styling: string = '';

	$: direction = horizontal ? 'flex-row' : 'flex-col';
	$: justification = 'justify-' + justify;
	$: alignment = 'items-' + align;
	$: grow = extensible ? 'flex-1' : 'flex-none';
</script>

<!-- https://stackoverflow.com/questions/29467660/how-to-stretch-children-to-fill-cross-axis -->
<div class="flex {direction} {justification} {alignment} {bgColor} {grow} {styling}">
	<slot />
</div>
