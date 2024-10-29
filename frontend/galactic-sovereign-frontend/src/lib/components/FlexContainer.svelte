<script lang="ts">
	export let vertical: boolean = true;
	export let justify: string = 'evenly';
	export let align: string = 'center';
	export let extensible: boolean = true;
	export let bgColor: string = 'bg-transparent';
	export let styling: string = '';

	$: direction = vertical ? 'flex-col' : 'flex-row';
	$: justification = generateJustifyClass(justify);
	$: alignment = generateAlignClass(align);
	// https://stackoverflow.com/questions/75999354/tailwindcss-content-larger-than-screen-when-adding-components
	$: grow = extensible ? 'flex-1' : 'flex-none';

	// https://tailwindcss.com/docs/content-configuration#dynamic-class-names
	function generateJustifyClass(justify: string): string {
		switch (justify) {
			case 'center':
				return 'justify-center';
			case 'start':
				return 'justify-start';
			case 'between':
				return 'justify-between';
			default:
				return 'justify-evenly';
		}
	}

	function generateAlignClass(align: string): string {
		switch (align) {
			case 'start':
				return 'items-start';
			case 'stretch':
				return 'items-stretch';
			default:
				return 'items-center';
		}
	}
</script>

<!-- https://stackoverflow.com/questions/29467660/how-to-stretch-children-to-fill-cross-axis -->
<div class="flex {direction} {justification} {alignment} {bgColor} {grow} {styling}">
	<slot />
</div>
