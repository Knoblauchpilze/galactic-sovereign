import { type ApiResource } from '$lib/resources';
import { type ApiBuildingAction } from './actions';

export interface ApiBuilding {
	readonly id: string;
	readonly name: string;
	readonly costs: ApiBuildingCost[];
}

export interface ApiBuildingCost {
	readonly resource: string;
	readonly cost: number;
	readonly progress: number;
}

export class Building {
	readonly id: string = '00000000-0000-0000-0000-000000000000';
	readonly name: string = '';
	readonly costs: BuildingCost[] = [];

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
	}

	public toJson(): ApiBuilding {
		return {
			id: this.id,
			name: this.name,
			costs: this.costs.map((c) => ({
				resource: c.resource,
				cost: c.cost,
				progress: c.progress
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

export function mapPlanetBuildingsToUiBuildings(
	planet: string,
	planetBuildings: PlanetBuilding[],
	planetActions: ApiBuildingAction[],
	apiBuildings: ApiBuilding[],
	apiResources: ApiResource[]
): UiBuilding[] {
	return apiBuildings.map((apiBuilding) => {
		const maybeBuilding = planetBuildings.find((r) => r.id === apiBuilding.id);
		if (maybeBuilding === undefined) {
			return {
				id: apiBuilding.id,
				name: apiBuilding.name,
				level: 0,
				planet: planet,
				hasAction: false,
				costs: mapBuildingCostsToUiBuildingCosts(apiBuilding.costs, apiResources, 0)
			};
		} else {
			const maybeAction = planetActions.find((a) => a.building === maybeBuilding.id);
			if (maybeAction === undefined) {
				return {
					id: apiBuilding.id,
					name: apiBuilding.name,
					level: maybeBuilding.level,
					planet: planet,
					hasAction: false,
					costs: mapBuildingCostsToUiBuildingCosts(apiBuilding.costs, apiResources, maybeBuilding.level)
				};
			}

			return {
				id: apiBuilding.id,
				name: apiBuilding.name,
				level: maybeBuilding.level,
				planet: planet,
				hasAction: true,
				action: maybeAction.id,
				nextLevel: maybeAction.desiredLevel,
				completedAt: maybeAction.completedAt,
				costs: mapBuildingCostsToUiBuildingCosts(apiBuilding.costs, apiResources, maybeBuilding.level)
			};
		}
	});
}
