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

export interface PlanetResourceProduction {
	readonly resource: string;
	readonly building: string;
	readonly production: number;
}

export function parsePlanetResourceProductions(data: object[]): PlanetResourceProduction[] {
	const out: PlanetResourceProduction[] = [];

	for (const maybeProduction of data) {
		const hasResource =
			'resource' in maybeProduction && typeof maybeProduction.resource === 'string';
		const hasBuilding =
			'building' in maybeProduction && typeof maybeProduction.building === 'string';
		const hasProduction =
			'production' in maybeProduction && typeof maybeProduction.production === 'number';

		if (hasResource && hasProduction) {
			const res: PlanetResourceProduction = {
				resource: maybeProduction.resource as string,
				building: hasBuilding ? (maybeProduction.building as string) : '',
				production: maybeProduction.production as number
			};

			out.push(res);
		}
	}

	return out;
}

export interface PlanetResourceStorage {
	readonly resource: string;
	readonly storage: number;
}

export function parsePlanetResourceStorages(data: object[]): PlanetResourceStorage[] {
	const out: PlanetResourceStorage[] = [];

	for (const maybeStorage of data) {
		const hasResource = 'resource' in maybeStorage && typeof maybeStorage.resource === 'string';
		const hasStorage = 'storage' in maybeStorage && typeof maybeStorage.storage === 'number';

		if (hasResource && hasStorage) {
			const res: PlanetResourceStorage = {
				resource: maybeStorage.resource as string,
				storage: maybeStorage.storage as number
			};

			out.push(res);
		}
	}

	return out;
}

export interface UiResource {
	readonly name: string;
	readonly amount: number;
	readonly production: number;
	readonly storage: number;
}

export function mapPlanetResourcesToUiResources(
	planetResources: PlanetResource[],
	planetProductions: PlanetResourceProduction[],
	planetStorages: PlanetResourceStorage[],
	apiResources: ApiResource[]
): UiResource[] {
	return apiResources.map((apiResource) => {
		const maybeResource = planetResources.find((r) => r.id === apiResource.id);

		const production = planetProductions.reduce((currentProduction, resource) => {
			if (resource.resource === apiResource.id) {
				return currentProduction + resource.production;
			}
			return currentProduction;
		}, 0);

		const planetStorage = planetStorages.find((s) => s.resource === apiResource.id);
		const storage = planetStorage === undefined ? 0 : planetStorage.storage;

		if (maybeResource === undefined) {
			const isProducing = storage > 0;
			return {
				name: apiResource.name,
				amount: 0,
				production: isProducing ? production : 0,
				storage: storage
			};
		} else {
			const isProducing = storage > maybeResource.amount;
			return {
				name: apiResource.name,
				amount: maybeResource.amount,
				production: isProducing ? production : 0,
				storage: storage
			};
		}
	});
}
