import { buildUrl } from '$lib/rest/api';
import { safeFetchJson, type ApiResponse } from '@totocorpsoftwareinc/frontend-toolkit';

export async function getPlanet(apiKey: string, id: string): Promise<ApiResponse> {
	const url = buildUrl('planets/' + id);

	const params = {
		method: 'GET',
		headers: {
			'X-Api-Key': apiKey
		}
	};

	return safeFetchJson(url, params);
}

export async function fetchPlanetsFromPlayer(
	playerId: string,
	apiKey: string
): Promise<ApiResponse> {
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

	return safeFetchJson(url, params);
}
