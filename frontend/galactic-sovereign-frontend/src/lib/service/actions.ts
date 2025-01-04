import { buildUrl } from '$lib/rest/api';
import { safeFetchJson, type ApiResponse } from '@totocorpsoftwareinc/frontend-toolkit';

export async function createBuildingAction(
	apiKey: string,
	planet: string,
	building: string
): Promise<ApiResponse> {
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

	return safeFetchJson(url, params);
}

export async function deleteBuildingAction(apiKey: string, action: string): Promise<ApiResponse> {
	const url = buildUrl('actions/' + action);

	const params = {
		method: 'DELETE',
		headers: {
			'X-Api-Key': apiKey
		}
	};

	return safeFetchJson(url, params);
}
