<script lang="ts">
	import activeScreen from '$lib/stores/activeScreen';

	import { NegativeConfirmationButton, StyledText, Timer } from '$lib/components';

	import { type UiBuildingAction } from '$lib/game/actions';

	interface Props {
		action: UiBuildingAction;
		onCompleted: () => void;
	}

	let { action, onCompleted }: Props = $props();

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

	let cancelButtonClass = $state(serverRemainingMs > 0 ? '' : 'hidden');
	let actionCompleted = $state(serverRemainingMs < 0);

	// The image path looks like this:
	// /src/lib/assets/buildings/crystal_mine.webp
	// The name of the action looks like this:
	// crystal mine
	let images = $derived(
		Object.keys(modules).map((imagePath) => {
			return {
				building: imagePath
					// this removes the extension
					.replace(/^.*[\\/]/, '')
					// this removes anything before the last '/' character (keeping only the building name)
					.replace(/\..*$/, '')
					// this replaces the '_' by a space to match the action name
					.replace(/\_/, ' '),
				data: modules[imagePath].default
			};
		})
	);
	let actionImage = $derived(images.find((image) => image.building === action.name));

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
		<form method="POST" action="/planets/{action.planet}/{$activeScreen}?/deleteBuildingAction">
			<input class="hidden" id="action" name="action" value={action.id} />
			<NegativeConfirmationButton text="Cancel" enabled={!actionCompleted} />
		</form>
	</div>
</div>
