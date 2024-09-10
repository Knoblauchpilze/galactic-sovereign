import {
	type ApiBuilding,
	type ApiBuildingResourceProduction,
	type PlanetBuilding
} from '$lib/buildings';

export interface ApiResource {
	readonly id: string;
	readonly name: string;
}

export class Resource {
	readonly id: string = '00000000-0000-0000-0000-000000000000';
	readonly name: string = '';

	constructor(response: object) {
		if ('id' in response && typeof response.id === 'string') {
			this.id = response.id;
		}

		if ('name' in response && typeof response.name === 'string') {
			this.name = response.name;
		}
	}

	// https://stackoverflow.com/questions/65512526/cannot-stringify-arbitrary-non-pojos-and-invalid-prop-type-check-failed-for
	public toJson(): ApiResource {
		return {
			id: this.id,
			name: this.name
		};
	}
}

export function parseResources(data: object[]): Resource[] {
	const out: Resource[] = [];

	for (const maybeResource of data) {
		const hasResource = 'id' in maybeResource && typeof maybeResource.id === 'string';
		const hasName = 'name' in maybeResource && typeof maybeResource.name === 'string';

		if (hasResource && hasName) {
			out.push(new Resource(maybeResource));
		}
	}

	return out;
}

export interface PlanetResource {
	readonly id: string;
	readonly amount: number;
}

export function parsePlanetResources(data: object[]): PlanetResource[] {
	const out: PlanetResource[] = [];

	for (const maybeResource of data) {
		const hasResource = 'resource' in maybeResource && typeof maybeResource.resource === 'string';
		const hasAmount = 'amount' in maybeResource && typeof maybeResource.amount === 'number';

		if (hasResource && hasAmount) {
			const res: PlanetResource = {
				id: maybeResource.resource as string,
				amount: maybeResource.amount as number
			};

			out.push(res);
		}
	}

	return out;
}

function computeProductionForLevel(
	production: ApiBuildingResourceProduction,
	level: number
): number {
	return Math.floor(production.base * Math.pow(production.progress, level));
}

export interface UiResource {
	readonly name: string;
	readonly amount: number;
	readonly production: number;
}

interface ResourceProduction {
	readonly production: ApiBuildingResourceProduction;
	readonly level: number;
}

function mapApiBuildingToBuildingProduction(
	apiBuilding: ApiBuilding,
	planetBuildings: PlanetBuilding[]
): ResourceProduction[] {
	const maybePlanetBuilding = planetBuildings.find((pb) => pb.id === apiBuilding.id);
	const level = maybePlanetBuilding === undefined ? 0 : maybePlanetBuilding.level;

	return apiBuilding.resourceProductions.map((rp) => {
		return {
			production: rp,
			level: level
		};
	});
}

function mapApiBuildingsToBuildingProductions(
	apiBuildings: ApiBuilding[],
	planetBuildings: PlanetBuilding[]
): ResourceProduction[] {
	return apiBuildings
		.map((b) => mapApiBuildingToBuildingProduction(b, planetBuildings))
		.reduce((productions, currentProduction) => productions.concat(currentProduction), []);
}

export function mapPlanetResourcesAndBuildingsToUiResources(
	planetResources: PlanetResource[],
	apiResources: ApiResource[],
	apiBuildings: ApiBuilding[],
	planetBuildings: PlanetBuilding[]
): UiResource[] {
	const resourceProductions = mapApiBuildingsToBuildingProductions(apiBuildings, planetBuildings);

	return apiResources.map((apiResource) => {
		const maybeResource = planetResources.find((r) => r.id === apiResource.id);

		const production = resourceProductions
			.filter((p) => p.production.resource === apiResource.id)
			.reduce(
				(production, currentProduction) =>
					production +
					computeProductionForLevel(currentProduction.production, currentProduction.level),
				0
			);

		if (maybeResource === undefined) {
			return {
				name: apiResource.name,
				amount: 0,
				production: production
			};
		} else {
			return {
				name: apiResource.name,
				amount: maybeResource.amount,
				production: production
			};
		}
	});
}
