<script lang="ts">
	export let vertical: boolean = true;
	export let justify: string = 'evenly';
	export let align: string = 'center';
	export let extensible: boolean = true;
	export let bgColor: string = 'bg-transparent';
	export let styling: string = '';

	$: direction = vertical ? 'flex-col' : 'flex-row';
	$: justification = 'justify-' + justify;
	$: alignment = 'items-' + align;
	// https://stackoverflow.com/questions/75999354/tailwindcss-content-larger-than-screen-when-adding-components
	$: grow = extensible ? 'flex-1' : 'flex-none';
</script>

<!-- https://stackoverflow.com/questions/29467660/how-to-stretch-children-to-fill-cross-axis -->
<div class="flex {direction} {justification} {alignment} {bgColor} {grow} {styling}">
	<slot />
</div>
