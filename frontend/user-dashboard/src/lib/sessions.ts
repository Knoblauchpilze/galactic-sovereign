import { ResponseEnvelope, createEmptySuccessResponseEnvelope } from '$lib/responseEnvelope';
import { buildUrl, safeFetch } from '$lib/api';
import HttpStatus from '$lib/httpStatuses';

export async function loginUser(email: string, password: string): Promise<ResponseEnvelope> {
	const url = buildUrl('sessions');
	const body = JSON.stringify({ email: email, password: password });

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

export async function logoutUser(apiKey: string, userId: string): Promise<ResponseEnvelope> {
	const url = buildUrl('sessions/' + userId);

	const params = {
		method: 'DELETE',
		headers: {
			'X-Api-Key': apiKey
		}
	};

	const response = await safeFetch(url, params);

	if (response.status == HttpStatus.NO_CONTENT) {
		return createEmptySuccessResponseEnvelope();
	}

	const jsonContent = await response.json();

	return new ResponseEnvelope(jsonContent);
}
