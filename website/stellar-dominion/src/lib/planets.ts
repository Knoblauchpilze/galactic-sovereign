import ResponseEnvelope from './responseEnvelope';
import { buildUrl, safeFetch } from './api';

export default class Planet {
	readonly id: string = '00000000-0000-0000-0000-000000000000';
	readonly player: string = '00000000-0000-0000-0000-000000000000';
	readonly name: string = '';

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
	}
}

export async function getPlanet(apiKey: string, id: string): Promise<ResponseEnvelope> {
	const url = buildUrl('planet/' + id);

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
