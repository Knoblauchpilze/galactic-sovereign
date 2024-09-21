<script lang="ts">
	import { type UiBuilding, type UiBuildingCost } from '$lib/game/buildings';
	import { type UiResource } from '$lib/game/resources';
	import { StyledActionButton, StyledText } from '$lib/components';

	export let building: UiBuilding;
	export let availableResources: UiResource[];
	export let buildingActionAlreadyRunning: boolean;

	function canAfford(cost: UiBuildingCost, availableResources: UiResource[]): boolean {
		const maybeResource = availableResources.find((r) => r.name === cost.resource);
		return maybeResource === undefined || maybeResource.amount >= cost.cost;
	}

	function textColor(cost: UiBuildingCost, availableResources: UiResource[]): string {
		const affordable = canAfford(cost, availableResources);
		if (affordable) {
			return 'text-enabled';
		}
		return 'text-disabled';
	}

	// https://stackoverflow.com/questions/49296458/capitalize-first-letter-of-a-string-using-angular-or-typescript
	const title = building.name[0].toUpperCase() + building.name.slice(1);

	$: costs = building.costs.map((c) => ({
		resource: c.resource,
		cost: c.cost,
		color: textColor(c, availableResources)
	}));

	$: gains = building.resourcesProduction.map((rp) => ({
		resource: rp.resource,
		nextProduction: rp.nextProduction,
		gain: rp.gain
	}));
	$: needsGainsSection = gains.length > 0;

	const isAffordable = building.costs.reduce(
		(currentlyAffordable, cost) => currentlyAffordable && canAfford(cost, availableResources),
		true
	);
</script>

<div class="p-4 m-2 bg-overlay">
	<StyledText text="{title} (level {building.level})" styling="font-bold" />
	<StyledText text="Required for level {building.level + 1}:" textColor="text-white" />
	<table>
		{#each costs as cost}
			<tr>
				<td class="text-white capitalize">{cost.resource}:</td>
				<td class={cost.color}>{cost.cost}</td>
			</tr>
		{/each}
	</table>
	{#if needsGainsSection}
		<StyledText text="Production:" textColor="text-white" />
		<table>
			{#each gains as gain}
				<tr>
					<td class="text-white capitalize">{gain.resource}:</td>
					<td class="text-enabled">{gain.nextProduction}(+{gain.gain})</td>
				</tr>
			{/each}
		</table>
	{/if}
	<!-- https://kit.svelte.dev/docs/form-actions#default-actions -->
	<form method="POST" action="/planets/{building.planet}/overview?/createBuildingAction">
		<input class="hidden" id="building" name="building" value={building.id} />
		<StyledActionButton text="Upgrade" enabled={isAffordable && !buildingActionAlreadyRunning} />
	</form>
</div>
