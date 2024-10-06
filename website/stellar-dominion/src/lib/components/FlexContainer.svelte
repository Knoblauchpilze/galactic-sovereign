<script lang="ts">
	export let horizontal: boolean = false;
	export let bgColor: string = 'bg-transparent';

	$: direction = horizontal ? 'flex-row' : 'flex-col';
</script>

<!-- https://stackoverflow.com/questions/29467660/how-to-stretch-children-to-fill-cross-axis -->
<div class="flex {direction} justify-around items-center {bgColor} flex-1">
	<slot />
</div>
