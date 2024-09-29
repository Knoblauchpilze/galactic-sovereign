<script lang="ts">
	import { type UiBuildingAction } from '$lib/game/actions';
	import { StyledActionButton, StyledText, Timer } from '$lib/components';

	export let action: UiBuildingAction;
	export let onCompleted: () => void;

	// https://kit.svelte.dev/docs/images#sveltejs-enhanced-img-dynamically-choosing-an-image
	// https://github.com/vitejs/vite/issues/9599#issuecomment-1209333753
	const modules = import.meta.glob<Record<string, string>>('$lib/assets/buildings/*.webp', {
		eager: true,
		query: {
			enhanced: true
		}
	});

	const title = action.name[0].toUpperCase() + action.name.slice(1);

	// https://stackoverflow.com/questions/14980014/how-can-i-calculate-the-time-between-2-dates-in-typescript
	const serverRemainingMs = action.completedAt.getTime() - Date.now();

	let cancelButtonClass = serverRemainingMs > 0 ? '' : 'hidden';
	let actionCompleted = serverRemainingMs < 0;

	$: images = Object.keys(modules).map((imagePath) => {
		return {
			building: imagePath
				.replace(/^.*[\\/]/, '')
				.replace(/\..*$/, '')
				.replace(/\_/, ' '),
			data: modules[imagePath].default
		};
	});
	$: actionImage = images.find((image) => image.building === action.name);

	function onActionCompleted() {
		cancelButtonClass = 'hidden';
		actionCompleted = true;

		onCompleted();
	}
</script>

<div class="p-4 m-2 bg-overlay">
	<StyledText text={title} styling="font-bold" />
	{#if actionImage !== undefined}
		<enhanced:img src={actionImage.data} alt="Building visual" width="150" height="150" />
	{/if}
	<StyledText text="Upgrade to level {action.nextLevel}" textColor="text-white" />
	<Timer durationMs={serverRemainingMs} onFinished={onActionCompleted} />
	<div class={cancelButtonClass}>
		<form method="POST" action="/planets/{action.planet}/overview?/deleteBuildingAction">
			<input class="hidden" id="action" name="action" value={action.id} />
			<StyledActionButton text="Cancel" enabled={!actionCompleted} negativeConfirmation />
		</form>
	</div>
</div>
