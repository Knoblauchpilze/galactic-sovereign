<script lang="ts">
	import { type UiBuildingAction } from '$lib/actions';
	import { StyledText } from '$lib/components';
	import { msToTimeStringOrFinished } from '$lib/time';

	export let action: UiBuildingAction;

	const title = action.name[0].toUpperCase() + action.name.slice(1);

	// https://stackoverflow.com/questions/14980014/how-can-i-calculate-the-time-between-2-dates-in-typescript
	const remainingMs = action.completedAt.getTime() - Date.now();

	let remaining = msToTimeStringOrFinished(remainingMs);
	let textColor = remainingMs <= 0 ? 'text-enabled' : 'text-white';

	console.log('completed at: ', action.completedAt.toString());
	console.log('remaining ms: ', remainingMs);
	console.log('remaining: ', remaining);
</script>

<div class="p-4 m-2 bg-overlay">
	<StyledText text={title} styling="font-bold" />
	<StyledText text="Upgrade to level {action.nextLevel}" textColor="text-white" />
	<StyledText text={remaining} {textColor} />
</div>
