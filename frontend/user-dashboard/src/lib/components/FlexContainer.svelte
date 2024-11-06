<script lang="ts">
	import { type Snippet } from 'svelte';

	interface Props {
		vertical?: boolean;
		justify?: string;
		align?: string;
		extensible?: boolean;
		bgColor?: string;
		styling?: string;
		children?: Snippet;
	}

	let {
		vertical = true,
		justify = 'evenly',
		align = 'center',
		extensible = true,
		bgColor = 'bg-transparent',
		styling = '',
		children
	}: Props = $props();

	function generateJustifyClass(justify: string): string {
		switch (justify) {
			case 'center':
				return 'justify-center';
			case 'start':
				return 'justify-start';
			case 'between':
				return 'justify-between';
			case 'around':
				return 'justify-around';
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
	let direction = $derived(vertical ? 'flex-col' : 'flex-row');
	let justification = $derived(generateJustifyClass(justify));
	let alignment = $derived(generateAlignClass(align));
	let grow = $derived(extensible ? 'flex-1' : 'flex-none');
</script>

<div class="flex {direction} {justification} {alignment} {bgColor} {grow} {styling}">
	{@render children?.()}
</div>
