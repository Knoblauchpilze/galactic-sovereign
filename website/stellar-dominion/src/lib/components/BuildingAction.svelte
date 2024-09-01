<script lang="ts">
	import { type UiBuildingAction } from '$lib/actions';
	import { StyledActionButton, StyledText, Timer } from '$lib/components';

	export let action: UiBuildingAction;

	const title = action.name[0].toUpperCase() + action.name.slice(1);

	// https://stackoverflow.com/questions/14980014/how-can-i-calculate-the-time-between-2-dates-in-typescript
	const serverRemainingMs = action.completedAt.getTime() - Date.now();

	let cancelButtonClass = serverRemainingMs > 0 ? '' : 'hidden';
	let actionCompleted = serverRemainingMs < 0;

	function onActionCompleted() {
		cancelButtonClass = 'hidden';
		actionCompleted = true;
	}
</script>

<div class="p-4 m-2 bg-overlay">
	<StyledText text={title} styling="font-bold" />
	<StyledText text="Upgrade to level {action.nextLevel}" textColor="text-white" />
	<Timer durationMs={serverRemainingMs} onFinished={onActionCompleted} />
	<div class={cancelButtonClass}>
		<form method="POST" action="/planets/{action.planet}/overview?/deleteBuildingAction">
			<input class="hidden" id="action" name="action" value={action.id} />
			<StyledActionButton text="Cancel" enabled={!actionCompleted} negativeConfirmation />
		</form>
	</div>
</div>
