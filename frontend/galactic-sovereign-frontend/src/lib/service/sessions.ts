import { buildUserUrl } from '$lib/rest/api';
import { safeFetchJson, type ApiResponse } from '@totocorpsoftwareinc/frontend-toolkit';

export async function loginUser(email: string, password: string): Promise<ApiResponse> {
	const url = buildUserUrl('sessions');
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

export async function logoutUser(apiKey: string, userId: string): Promise<ApiResponse> {
	const url = buildUserUrl('sessions/' + userId);

	const params = {
		method: 'DELETE',
		headers: {
			'X-Api-Key': apiKey
		}
	};

	return safeFetchJson(url, params);
}
