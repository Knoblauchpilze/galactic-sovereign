<script lang="ts">
	import { type UiBuilding, type UiBuildingCost } from '$lib/game/buildings';
	import { type UiResource } from '$lib/game/resources';
	import { StyledActionButton, StyledText } from '$lib/components';

	// https://kit.svelte.dev/docs/images#sveltejs-enhanced-img-dynamically-choosing-an-image
	// https://github.com/vitejs/vite/issues/9599#issuecomment-1209333753
	const modules = import.meta.glob<Record<string, string>>('$lib/assets/buildings/*.webp', {
		eager: true,
		query: {
			enhanced: true
		}
	});

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
	$: title = building.name[0].toUpperCase() + building.name.slice(1);

	// https://stackoverflow.com/questions/68060723/glob-import-of-image-urls-in-sveltekit
	// https://stackoverflow.com/questions/423376/how-to-get-the-file-name-from-a-full-path-using-javascript
	$: images = Object.keys(modules).map((imagePath) => {
		return {
			building: imagePath
				.replace(/^.*[\\/]/, '')
				.replace(/\..*$/, '')
				.replace(/\_/, ' '),
			data: modules[imagePath].default
		};
	});

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

	$: buildingImage = images.find((image) => image.building === building.name);

	const isAffordable = building.costs.reduce(
		(currentlyAffordable, cost) => currentlyAffordable && canAfford(cost, availableResources),
		true
	);
</script>

<div class="p-4 m-2 bg-overlay">
	<StyledText text="{title} (level {building.level})" styling="font-bold" />
	{#if buildingImage !== undefined}
		<enhanced:img src={buildingImage.data} alt="Building visual" width="100" height="100" />
	{/if}
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
