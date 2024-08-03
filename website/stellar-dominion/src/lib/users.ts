import ResponseEnvelope from '$lib/responseEnvelope';
import { buildUserUrl, safeFetch } from '$lib/api';

export default class User {
	readonly id: string = '00000000-0000-0000-0000-000000000000';
	readonly email: string = '';
	readonly password: string = '';
	readonly createdAt: Date = new Date();

	constructor(response: ResponseEnvelope) {
		if (response.error()) {
			return;
		}

		if ('id' in response.details && typeof response.details.id === 'string') {
			this.id = response.details.id;
		}

		if ('email' in response.details && typeof response.details.email === 'string') {
			this.email = response.details.email;
		}

		if ('password' in response.details && typeof response.details.password === 'string') {
			this.password = response.details.password;
		}

		if ('createdAt' in response.details && typeof response.details.createdAt === 'string') {
			this.createdAt = new Date(response.details.createdAt);
		}
	}
}

export async function createUser(email: string, password: string): Promise<ResponseEnvelope> {
	const url = buildUserUrl('');
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
