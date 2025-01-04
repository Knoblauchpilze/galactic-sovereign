import { buildUrl } from '$lib/rest/api';
import { safeFetchJson, type ApiResponse } from '@totocorpsoftwareinc/frontend-toolkit';

export async function getUniverse(id: string, apiKey: string): Promise<ApiResponse> {
	const url = buildUrl('universes/' + id);

	const params = {
		method: 'GET',
		headers: {
			'X-Api-Key': apiKey
		}
	};

	return safeFetchJson(url, params);
}

export async function getUniverses(apiKey: string): Promise<ApiResponse> {
	const url = buildUrl('universes');

	const params = {
		method: 'GET',
		headers: {
			'X-Api-Key': apiKey
		}
	};

	return safeFetchJson(url, params);
}
