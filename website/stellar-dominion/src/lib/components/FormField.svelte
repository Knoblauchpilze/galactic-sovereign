<script lang="ts">
	export let label: string = 'My label';
	export let labelId: string;
</script>

<div class="flex flex-col text-secondary">
	<label for={labelId}>{label}</label>
	<slot />
</div>
