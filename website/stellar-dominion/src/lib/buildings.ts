import { ResponseEnvelope } from '$lib/responseEnvelope';

export interface ApiBuilding {
	readonly id: string;
	readonly name: string;
}

export class Building {
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

	public toJson(): ApiBuilding {
		return {
			id: this.id,
			name: this.name
		};
	}
}

export function responseToBuildingsArray(response: ResponseEnvelope): Building[] {
	if (response.error()) {
		return [];
	}

	const details = response.getDetails();
	if (!Array.isArray(details)) {
		return [];
	}

	return details.map((maybeBuilding) => new Building(maybeBuilding));
}

export interface PlanetBuilding {
	readonly id: string;
	readonly level: number;
}

export function parseBuildings(data: object[]): PlanetBuilding[] {
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

export interface UiBuilding {
	readonly name: string;
	readonly level: number;
}

export function mapPlanetBuildingsToApiBuildings(
	planetBuildings: PlanetBuilding[],
	apiBuildings: ApiBuilding[]
): UiBuilding[] {
	return apiBuildings.map((apiBuilding) => {
		const maybeBuilding = planetBuildings.find((r) => r.id === apiBuilding.id);
		if (maybeBuilding === undefined) {
			return {
				name: apiBuilding.name,
				level: 0
			};
		} else {
			return {
				name: apiBuilding.name,
				level: maybeBuilding.level
			};
		}
	});
}
