<script lang="ts">
	interface Props {
		text?: string;
		textColor?: string;
		styling?: string;
	}

	let { text = 'Missing text', textColor = 'text-secondary', styling = '' }: Props = $props();
</script>

<p class="{textColor} {styling}">{text}</p>
