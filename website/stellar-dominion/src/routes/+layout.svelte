<script lang="ts">
	const width: string = 'w-3/5';
	const height: string = 'h-4/5';
	const bgImage: string = 'bg-homepage';
</script>

<!-- https://stackoverflow.com/questions/70805041/background-image-in-tailwindcss-using-dynamic-url-react-js -->
<div  class="relative h-full bg-center bg-no-repeat bg-cover {bgImage}">
	<div class="flex flex-col items-center justify-around absolute left-1/2 top-1/2 {width} {height} -translate-x-1/2 -translate-y-1/2 bg-[#0005]">
		<slot/>
	</div>
</div>

<style>
	:global(body),
	:global(html) {
		height: 100%;
	}

	:global(body) {
		margin: 0;
		padding: 0;
	}
</style>