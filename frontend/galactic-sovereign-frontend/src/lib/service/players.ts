import { buildUrl } from '$lib/rest/api';
import { safeFetchJson, type ApiResponse } from '@totocorpsoftwareinc/frontend-toolkit';

export async function createPlayer(
	apiUserId: string,
	universeId: string,
	playerName: string,
	apiKey: string
): Promise<ApiResponse> {
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

	return safeFetchJson(url, params);
}

export async function fetchPlayerFromApiUser(
	apiUserId: string,
	apiKey: string
): Promise<ApiResponse> {
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

	return safeFetchJson(url, params);
}
