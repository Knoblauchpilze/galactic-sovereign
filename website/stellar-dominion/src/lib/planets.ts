import { ResponseEnvelope } from '$lib/responseEnvelope';
import { buildUrl, safeFetch } from '$lib/api';
import { type PlanetResource, parsePlanetResources } from '$lib/resources';
import { type PlanetBuilding, parsePlanetBuildings } from '$lib/buildings';
import { type BuildingAction, type ApiBuildingAction, parseBuildingActions } from '$lib/actions';

export interface ApiPlanet {
	readonly id: string;
	readonly player: string;
	readonly name: string;

	readonly resources: PlanetResource[];
	readonly buildings: PlanetBuilding[];
	readonly buildingActions: ApiBuildingAction[];
}

export class Planet {
	readonly id: string = '00000000-0000-0000-0000-000000000000';
	readonly player: string = '00000000-0000-0000-0000-000000000000';
	readonly name: string = '';

	readonly resources: PlanetResource[] = [];
	readonly buildings: PlanetBuilding[] = [];
	readonly buildingActions: BuildingAction[] = [];

	constructor(response: object) {
		if ('id' in response && typeof response.id === 'string') {
			this.id = response.id;
		}

		if ('player' in response && typeof response.player === 'string') {
			this.player = response.player;
		}

		if ('name' in response && typeof response.name === 'string') {
			this.name = response.name;
		}

		if ('resources' in response && Array.isArray(response.resources)) {
			this.resources = parsePlanetResources(response.resources);
		}

		if ('buildings' in response && Array.isArray(response.buildings)) {
			this.buildings = parsePlanetBuildings(response.buildings);
		}

		if ('buildingActions' in response && Array.isArray(response.buildingActions)) {
			this.buildingActions = parseBuildingActions(response.buildingActions);
		}
	}

	public toJson(): ApiPlanet {
		return {
			id: this.id,
			player: this.player,
			name: this.name,

			resources: this.resources,
			buildings: this.buildings,
			buildingActions: this.buildingActions.map((a) => a.toJson())
		};
	}
}

export async function getPlanet(apiKey: string, id: string): Promise<ResponseEnvelope> {
	const url = buildUrl('planets/' + id);

	const params = {
		method: 'GET',
		headers: {
			'X-Api-Key': apiKey
		}
	};

	const response = await safeFetch(url, params);
	const jsonContent = await response.json();

	return new ResponseEnvelope(jsonContent);
}

export async function fetchPlanetsFromPlayer(
	playerId: string,
	apiKey: string
): Promise<ResponseEnvelope> {
	let url = buildUrl('planets');

	const queryParams = {
		player: playerId
	};
	url += '?' + new URLSearchParams(queryParams).toString();

	const params = {
		method: 'GET',
		headers: {
			'X-Api-Key': apiKey
		}
	};

	const response = await safeFetch(url, params);
	const jsonContent = await response.json();

	return new ResponseEnvelope(jsonContent);
}

export async function createBuildingAction(
	apiKey: string,
	planet: string,
	building: string
): Promise<ResponseEnvelope> {
	const url = buildUrl('planets/' + planet + '/buildings');
	const body = JSON.stringify({ planet: planet, building: building });

	const params = {
		method: 'POST',
		body: body,
		headers: {
			'content-type': 'application/json',
			'X-Api-Key': apiKey
		}
	};

	const response = await safeFetch(url, params);
	const jsonContent = await response.json();

	return new ResponseEnvelope(jsonContent);
}

export function responseToPlanetArray(response: ResponseEnvelope): Planet[] {
	if (response.error()) {
		return [];
	}

	const details = response.getDetails();
	if (!Array.isArray(details)) {
		return [];
	}

	return details.map((maybePlanet) => new Planet(maybePlanet));
}
