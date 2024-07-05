import ResponseEnvelope from './responseEnvelope';
import { buildUrl, safeFetch } from './api';

export async function createUser(email: string, password: string): Promise<ResponseEnvelope> {
	const url = buildUrl('');
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
