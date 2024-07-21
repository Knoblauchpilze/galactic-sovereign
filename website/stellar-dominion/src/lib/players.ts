import ResponseEnvelope from './responseEnvelope';
import { buildUrl, safeFetch } from './api';
import User, { createUser } from '$lib/users';

export async function createPlayer(
	apiUserId: string,
	universeId: string,
	playerName: string
): Promise<ResponseEnvelope> {
	const url = buildUrl('players');
	const body = JSON.stringify({ api_user: apiUserId, universe: universeId, name: playerName });

	const params = {
		method: 'POST',
		body: body,
		headers: {
			'content-type': 'application/json'
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

	console.log('signup response: ', JSON.stringify(apiUser));

	const playerResponse = await createPlayer(apiUser.id, universeId, playerName);

	console.log('player response: ', JSON.stringify(playerResponse));
	return playerResponse;
}
