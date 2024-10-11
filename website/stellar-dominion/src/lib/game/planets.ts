import { ResponseEnvelope, createEmptySuccessResponseEnvelope } from '$lib/responseEnvelope';
import { buildUrl, safeFetch } from '$lib/api';
import {
	type PlanetResource,
	type PlanetResourceProduction,
	type PlanetResourceStorage,
	parsePlanetResources,
	parsePlanetResourceProductions,
	parsePlanetResourceStorages
} from '$lib/game/resources';
import { type PlanetBuilding, parsePlanetBuildings } from '$lib/game/buildings';
import {
	type BuildingAction,
	type ApiBuildingAction,
	parseBuildingActions
} from '$lib/game/actions';
import HttpStatus from '$lib/httpStatuses';

export interface ApiPlanet {
	readonly id: string;
	readonly player: string;
	readonly name: string;

	readonly createdAt: Date;

	readonly resources: PlanetResource[];
	readonly productions: PlanetResourceProduction[];
	readonly storages: PlanetResourceStorage[];
	readonly buildings: PlanetBuilding[];
	readonly buildingActions: ApiBuildingAction[];
}

export class Planet {
	readonly id: string = '00000000-0000-0000-0000-000000000000';
	readonly player: string = '00000000-0000-0000-0000-000000000000';
	readonly name: string = '';

	readonly createdAt: Date = new Date();

	readonly resources: PlanetResource[] = [];
	readonly productions: PlanetResourceProduction[] = [];
	readonly storages: PlanetResourceStorage[] = [];
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

		if ('createdAt' in response && typeof response.createdAt === 'string') {
			this.createdAt = new Date(response.createdAt);
		}

		if ('resources' in response && Array.isArray(response.resources)) {
			this.resources = parsePlanetResources(response.resources);
		}

		if ('productions' in response && Array.isArray(response.productions)) {
			this.productions = parsePlanetResourceProductions(response.productions);
		}

		if ('storages' in response && Array.isArray(response.storages)) {
			this.storages = parsePlanetResourceStorages(response.storages);
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

			createdAt: this.createdAt,

			resources: this.resources,
			productions: this.productions,
			storages: this.storages,
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
	const url = buildUrl('planets/' + planet + '/actions');
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

export async function deleteBuildingAction(
	apiKey: string,
	action: string
): Promise<ResponseEnvelope> {
	const url = buildUrl('actions/' + action);

	const params = {
		method: 'DELETE',
		headers: {
			'X-Api-Key': apiKey
		}
	};

	const response = await safeFetch(url, params);

	if (response.status == HttpStatus.NO_CONTENT) {
		return createEmptySuccessResponseEnvelope();
	}

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
