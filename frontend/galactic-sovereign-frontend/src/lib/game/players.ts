import { ResponseEnvelope } from '$lib/responseEnvelope';
import { buildUrl, safeFetch } from '$lib/api';
import { type ApiUniverse } from '$lib/game/universes';

export interface ApiPlayer {
	readonly id: string;
	readonly name: string;
	readonly universe: string;
}

export class Player {
	readonly id: string = '00000000-0000-0000-0000-000000000000';
	readonly apiUser: string = '';
	readonly universe: string = '';
	readonly name: string = '';
	readonly createdAt: Date = new Date();

	constructor(response: object) {
		if ('id' in response && typeof response.id === 'string') {
			this.id = response.id;
		}

		if ('apiUser' in response && typeof response.apiUser === 'string') {
			this.apiUser = response.apiUser;
		}

		if ('universe' in response && typeof response.universe === 'string') {
			this.universe = response.universe;
		}

		if ('name' in response && typeof response.name === 'string') {
			this.name = response.name;
		}

		if ('createdAt' in response && typeof response.createdAt === 'string') {
			this.createdAt = new Date(response.createdAt);
		}
	}

	public toJson(): ApiPlayer {
		return {
			id: this.id,
			name: this.name,
			universe: this.universe
		};
	}
}

export async function createPlayer(
	apiUserId: string,
	universeId: string,
	playerName: string,
	apiKey: string
): Promise<ResponseEnvelope> {
	const url = buildUrl('players');
	const body = JSON.stringify({ api_user: apiUserId, universe: universeId, name: playerName });

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

export async function fetchPlayerFromApiUser(
	apiUserId: string,
	apiKey: string
): Promise<ResponseEnvelope> {
	let url = buildUrl('players');

	// https://medium.com/meta-box/how-to-send-get-and-post-requests-with-javascript-fetch-api-d0685b7ee6ed
	const queryParams = {
		api_user: apiUserId
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

export function responseToPlayerArray(response: ResponseEnvelope): Player[] {
	if (response.error()) {
		return [];
	}

	const details = response.getDetails();
	if (!Array.isArray(details)) {
		return [];
	}

	// https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Array/map
	return details.map((maybePlayer) => new Player(maybePlayer));
}

export class UiPlayer {
	readonly id: string = '';
	readonly name: string = '';
	readonly universeId: string = '';
	readonly universeName: string = '';
}

export function mapPlayersToUiPlayers(
	apiPlayers: Player[],
	apiUniverses: ApiUniverse[]
): UiPlayer[] {
	return apiPlayers.map((apiPlayer) => {
		const universe = apiUniverses.find((u) => u.id === apiPlayer.universe);

		if (universe === undefined) {
			return {
				id: apiPlayer.id,
				name: apiPlayer.name,
				universeId: '',
				universeName: 'Unknown universe'
			};
		} else {
			return {
				id: apiPlayer.id,
				name: apiPlayer.name,
				universeId: universe.id,
				universeName: universe.name
			};
		}
	});
}
