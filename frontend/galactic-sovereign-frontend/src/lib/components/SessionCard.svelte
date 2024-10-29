<script lang="ts">
	import { FlexContainer, StyledActionButton } from '$lib/components';

	import { type UiPlayer } from '$lib/game/players';

	export let player: UiPlayer;
</script>

<FlexContainer
	vertical={false}
	bgColor={'bg-overlay'}
	justify={'between'}
	styling={'px-8 py-2 my-2 rounded-[8px]'}
>
	<table class="mx-2">
		<tr>
			<td class="text-white capitalize">Universe:</td>
			<td class="text-secondary">{player.universeName}</td>
		</tr>
		<tr>
			<td class="text-white capitalize">Player:</td>
			<td class="text-secondary">{player.name}</td>
		</tr>
	</table>

	<form method="POST" action="?/login">
		<input class="hidden" id="universe" name="universe" value={player.universeId} />
		<input class="hidden" id="player" name="player" value={player.name} />
		<StyledActionButton text="Join" enabled={true} />
	</form>
</FlexContainer>
