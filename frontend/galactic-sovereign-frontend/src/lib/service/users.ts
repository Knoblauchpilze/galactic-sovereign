import { buildUserUrl } from '$lib/rest/api';
import { safeFetchJson, type ApiResponse } from '@totocorpsoftwareinc/frontend-toolkit';

export async function createUser(email: string, password: string): Promise<ApiResponse> {
	const url = buildUserUrl('');
	const body = JSON.stringify({ email: email, password: password });

	const params = {
		method: 'POST',
		body: body,
		headers: {
			'content-type': 'application/json'
		}
	};

	return safeFetchJson(url, params);
}

export async function getUser(apiKey: string, id: string): Promise<ApiResponse> {
	const url = buildUserUrl(id);

	const params = {
		method: 'GET',
		headers: {
			'X-Api-Key': apiKey
		}
	};

	return safeFetchJson(url, params);
}
