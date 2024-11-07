<script lang="ts">
	interface Props {
		text?: string;
	}

	let { text = 'Unknown error' }: Props = $props();
</script>

<p class="text-error">{text}</p>
