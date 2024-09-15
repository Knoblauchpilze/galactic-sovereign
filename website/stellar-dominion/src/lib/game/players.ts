import { ResponseEnvelope } from '$lib/responseEnvelope';
import { buildUrl, safeFetch } from '$lib/api';
import { ApiKey, loginUser, logoutUser } from '$lib/sessions';
import { User, createUser } from '$lib/users';

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

export async function registerPlayer(
	email: string,
	password: string,
	universeId: string,
	playerName: string
): Promise<ResponseEnvelope> {
	const signupResponse = await createUser(email, password);
	if (signupResponse.error()) {
		return signupResponse;
	}

	const apiUser = new User(signupResponse);

	const loginResponse = await loginUser(email, password);
	if (loginResponse.error()) {
		return loginResponse;
	}

	const apiKey = new ApiKey(loginResponse);

	const playerResponse = await createPlayer(apiUser.id, universeId, playerName, apiKey.key);
	if (playerResponse.error()) {
		return playerResponse;
	}

	const logoutResponse = await logoutUser(apiKey.key, apiKey.user);
	if (logoutResponse.error()) {
		return logoutResponse;
	}

	return playerResponse;
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
