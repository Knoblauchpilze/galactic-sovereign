<script lang="ts">
	export let text: string = 'Missing text';
	export let textColor: string = 'text-secondary';
	export let styling: string = '';
</script>

<p class="{textColor} {styling}">{text}</p>
