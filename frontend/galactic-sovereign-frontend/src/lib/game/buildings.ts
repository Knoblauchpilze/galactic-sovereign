import { type ApiResource } from '$lib/game/resources';
import { type ApiPlanet } from '$lib/game/planets';

export interface ApiBuilding {
	readonly id: string;
	readonly name: string;
	readonly costs: ApiBuildingCost[];
	readonly resourceProductions: ApiBuildingResourceProduction[];
}

export interface ApiBuildingCost {
	readonly resource: string;
	readonly cost: number;
	readonly progress: number;
}

export interface ApiBuildingResourceProduction {
	readonly resource: string;
	readonly base: number;
	readonly progress: number;
}

export class Building {
	readonly id: string = '00000000-0000-0000-0000-000000000000';
	readonly name: string = '';
	readonly costs: BuildingCost[] = [];
	readonly resourceProductions: BuildingResourceProduction[] = [];

	constructor(response: object) {
		if ('id' in response && typeof response.id === 'string') {
			this.id = response.id;
		}

		if ('name' in response && typeof response.name === 'string') {
			this.name = response.name;
		}

		if ('costs' in response && Array.isArray(response.costs)) {
			this.costs = parseBuildingCosts(response.costs);
		}

		if ('productions' in response && Array.isArray(response.productions)) {
			this.resourceProductions = parseBuildingResourceProductions(response.productions);
		}
	}

	public toJson(): ApiBuilding {
		return {
			id: this.id,
			name: this.name,
			costs: this.costs.map((c) => ({
				resource: c.resource,
				cost: c.cost,
				progress: c.progress
			})),
			resourceProductions: this.resourceProductions.map((p) => ({
				resource: p.resource,
				base: p.base,
				progress: p.progress
			}))
		};
	}
}

export function parseBuildings(data: object[]): Building[] {
	const out: Building[] = [];

	for (const maybeBuilding of data) {
		const hasBuilding = 'id' in maybeBuilding && typeof maybeBuilding.id === 'string';
		const hasName = 'name' in maybeBuilding && typeof maybeBuilding.name === 'string';

		if (hasBuilding && hasName) {
			out.push(new Building(maybeBuilding));
		}
	}

	return out;
}

export interface BuildingCost {
	readonly resource: string;
	readonly cost: number;
	readonly progress: number;
}

export function parseBuildingCosts(data: object[]): BuildingCost[] {
	const out: BuildingCost[] = [];

	for (const maybeCost of data) {
		const hasResource = 'resource' in maybeCost && typeof maybeCost.resource === 'string';
		const hasCost = 'cost' in maybeCost && typeof maybeCost.cost === 'number';
		const hasProgress = 'progress' in maybeCost && typeof maybeCost.progress === 'number';

		if (hasResource && hasCost && hasProgress) {
			const cost: BuildingCost = {
				resource: maybeCost.resource as string,
				cost: maybeCost.cost as number,
				progress: maybeCost.progress as number
			};

			out.push(cost);
		}
	}

	return out;
}

export interface BuildingResourceProduction {
	readonly resource: string;
	readonly base: number;
	readonly progress: number;
}

export function parseBuildingResourceProductions(data: object[]): BuildingResourceProduction[] {
	const out: BuildingResourceProduction[] = [];

	for (const maybeCost of data) {
		const hasResource = 'resource' in maybeCost && typeof maybeCost.resource === 'string';
		const hasBase = 'base' in maybeCost && typeof maybeCost.base === 'number';
		const hasProgress = 'progress' in maybeCost && typeof maybeCost.progress === 'number';

		if (hasResource && hasBase && hasProgress) {
			const resourceProduction: BuildingResourceProduction = {
				resource: maybeCost.resource as string,
				base: maybeCost.base as number,
				progress: maybeCost.progress as number
			};

			out.push(resourceProduction);
		}
	}

	return out;
}

export interface PlanetBuilding {
	readonly id: string;
	readonly level: number;
}

export function parsePlanetBuildings(data: object[]): PlanetBuilding[] {
	const out: PlanetBuilding[] = [];

	for (const maybeBuilding of data) {
		const hasBuilding = 'building' in maybeBuilding && typeof maybeBuilding.building === 'string';
		const hasLevel = 'level' in maybeBuilding && typeof maybeBuilding.level === 'number';

		if (hasBuilding && hasLevel) {
			const res: PlanetBuilding = {
				id: maybeBuilding.building as string,
				level: maybeBuilding.level as number
			};

			out.push(res);
		}
	}

	return out;
}

export interface UiBuildingCost {
	readonly resource: string;
	readonly cost: number;
}

export interface UiBuildingGain {
	readonly resource: string;
	readonly nextProduction: number;
	readonly gain: number;
}

export class UiBuilding {
	readonly id: string = '';
	readonly name: string = '';
	readonly level: number = 0;

	readonly planet: string = '';

	readonly hasAction: boolean = false;
	readonly action?: string = '';
	readonly nextLevel?: number = 1;
	readonly completedAt?: Date = new Date();

	readonly costs: UiBuildingCost[] = [];
	readonly resourcesProduction: UiBuildingGain[] = [];
}

function computeCostForLevel(cost: number, progress: number, level: number): number {
	return Math.floor(cost * Math.pow(progress, level - 1));
}

function mapBuildingCostsToUiBuildingCosts(
	buildingCosts: BuildingCost[],
	apiResources: ApiResource[],
	level: number
): UiBuildingCost[] {
	return buildingCosts.map((cost) => {
		const maybeResource = apiResources.find((r) => r.id === cost.resource);
		if (maybeResource === undefined) {
			return {
				resource: 'Unknown resource',
				cost: computeCostForLevel(cost.cost, cost.progress, level + 1)
			};
		} else {
			return {
				resource: maybeResource.name,
				cost: computeCostForLevel(cost.cost, cost.progress, level + 1)
			};
		}
	});
}

function computeProductionForLevel(production: number, progress: number, level: number): number {
	return Math.floor(production * Math.pow(progress, level));
}

function mapBuildingResourceProductionsToUiBuildingGains(
	resourceProductions: ApiBuildingResourceProduction[],
	apiResources: ApiResource[],
	level: number
): UiBuildingGain[] {
	return resourceProductions.map((production) => {
		const currentProduction = computeProductionForLevel(
			production.base,
			production.progress,
			level
		);
		const nextProduction = computeProductionForLevel(
			production.base,
			production.progress,
			level + 1
		);
		const gain = nextProduction - currentProduction;

		const maybeResource = apiResources.find((r) => r.id === production.resource);
		if (maybeResource === undefined) {
			return {
				resource: 'Unknown resource',
				nextProduction: nextProduction,
				gain: gain
			};
		} else {
			return {
				resource: maybeResource.name,
				nextProduction: nextProduction,
				gain: gain
			};
		}
	});
}

export function mapPlanetBuildingsToUiBuildings(
	planet: ApiPlanet,
	apiBuildings: ApiBuilding[],
	apiResources: ApiResource[]
): UiBuilding[] {
	return apiBuildings.map((apiBuilding) => {
		const maybeBuilding = planet.buildings.find((r) => r.id === apiBuilding.id);
		if (maybeBuilding === undefined) {
			return {
				id: apiBuilding.id,
				name: apiBuilding.name,
				level: 0,
				planet: planet.name,
				hasAction: false,
				costs: mapBuildingCostsToUiBuildingCosts(apiBuilding.costs, apiResources, 0),
				resourcesProduction: mapBuildingResourceProductionsToUiBuildingGains(
					apiBuilding.resourceProductions,
					apiResources,
					0
				)
			};
		} else {
			const maybeAction = planet.buildingActions.find((a) => a.building === maybeBuilding.id);
			if (maybeAction === undefined) {
				return {
					id: apiBuilding.id,
					name: apiBuilding.name,
					level: maybeBuilding.level,
					planet: planet.name,
					hasAction: false,
					costs: mapBuildingCostsToUiBuildingCosts(
						apiBuilding.costs,
						apiResources,
						maybeBuilding.level
					),
					resourcesProduction: mapBuildingResourceProductionsToUiBuildingGains(
						apiBuilding.resourceProductions,
						apiResources,
						maybeBuilding.level
					)
				};
			}

			return {
				id: apiBuilding.id,
				name: apiBuilding.name,
				level: maybeBuilding.level,
				planet: planet.name,
				hasAction: true,
				action: maybeAction.id,
				nextLevel: maybeAction.desiredLevel,
				completedAt: maybeAction.completedAt,
				costs: mapBuildingCostsToUiBuildingCosts(
					apiBuilding.costs,
					apiResources,
					maybeBuilding.level
				),
				resourcesProduction: mapBuildingResourceProductionsToUiBuildingGains(
					apiBuilding.resourceProductions,
					apiResources,
					maybeBuilding.level
				)
			};
		}
	});
}
