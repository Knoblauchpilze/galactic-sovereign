<script lang="ts">
	interface Props {
		text?: string;
	}

	let { text = 'Home' }: Props = $props();
</script>

<h1 class="px-12 py-4 text-3xl font-bold italic text-secondary">{text}</h1>
