import { ResponseEnvelope } from '$lib/responseEnvelope';
import { buildUrl, safeFetch } from '$lib/api';
import { Resource, parseResources } from '$lib/game/resources';
import { Building, parseBuildings } from '$lib/game/buildings';

export interface ApiUniverse {
	readonly id: string;
	readonly name: string;
}

export class Universe {
	readonly id: string = '00000000-0000-0000-0000-000000000000';
	readonly name: string = '';

	readonly resources: Resource[] = [];
	readonly buildings: Building[] = [];

	constructor(response: object) {
		if ('id' in response && typeof response.id === 'string') {
			this.id = response.id;
		}

		if ('name' in response && typeof response.name === 'string') {
			this.name = response.name;
		}

		if ('resources' in response && Array.isArray(response.resources)) {
			this.resources = parseResources(response.resources);
		}

		if ('buildings' in response && Array.isArray(response.buildings)) {
			this.buildings = parseBuildings(response.buildings);
		}
	}

	public toJson(): ApiUniverse {
		return {
			id: this.id,
			name: this.name
		};
	}
}

export async function getUniverse(id: string): Promise<ResponseEnvelope> {
	const url = buildUrl('universes/' + id);

	const params = {
		method: 'GET',
	};

	const response = await safeFetch(url, params);
	const jsonContent = await response.json();

	return new ResponseEnvelope(jsonContent);
}

export async function getUniverses(): Promise<ResponseEnvelope> {
	const url = buildUrl('universes');

	const params = {
		method: 'GET'
	};

	const response = await safeFetch(url, params);
	const jsonContent = await response.json();

	return new ResponseEnvelope(jsonContent);
}

export function responseToUniverseArray(response: ResponseEnvelope): Universe[] {
	if (response.error()) {
		return [];
	}

	const details = response.getDetails();
	if (!Array.isArray(details)) {
		return [];
	}

	return details.map((maybeUniverse) => new Universe(maybeUniverse));
}
