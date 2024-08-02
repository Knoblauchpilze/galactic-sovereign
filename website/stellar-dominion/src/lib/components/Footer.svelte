<script lang="ts">
	export let height: string = 'h-1/10';
</script>

<div class="fixed bottom-0 w-full {height} flex justify-around items-center bg-black text-white">
	<slot />
</div>
