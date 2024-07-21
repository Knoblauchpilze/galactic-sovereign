import ResponseEnvelope, { createEmptySuccessResponseEnvelope } from './responseEnvelope';
import { buildUserUrl, safeFetch } from './api';
import HttpStatus from './httpStatuses';

export default class ApiKey {
	readonly user: string = '00000000-0000-0000-0000-000000000000';
	readonly key: string = '00000000-0000-0000-0000-000000000000';
	readonly validUntil: Date = new Date();

	constructor(response: ResponseEnvelope) {
		if (response.error()) {
			return;
		}

		// https://stackoverflow.com/questions/455338/how-do-i-check-if-an-object-has-a-key-in-javascript
		if ('user' in response.details && typeof response.details.user === 'string') {
			this.user = response.details.user;
		}

		if ('key' in response.details && typeof response.details.key === 'string') {
			this.key = response.details.key;
		}

		if ('validUntil' in response.details && typeof response.details.validUntil === 'string') {
			this.validUntil = new Date(response.details.validUntil);
		}
	}
}

export async function loginUser(email: string, password: string): Promise<ResponseEnvelope> {
	const url = buildUserUrl('sessions');
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
	const url = buildUserUrl('sessions/' + userId);

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
