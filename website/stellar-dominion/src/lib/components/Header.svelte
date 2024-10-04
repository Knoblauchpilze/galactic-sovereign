<script lang="ts">
	export let height: string = 'h-1/10';
</script>

<div
	class="fixed top-0 left-0 w-full {height} flex justify-around items-center bg-black text-white"
>
	<slot />
</div>
