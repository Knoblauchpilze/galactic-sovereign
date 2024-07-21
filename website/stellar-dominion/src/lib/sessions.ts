import ResponseEnvelope from './responseEnvelope';
import { buildUserUrl, safeFetch } from './api';

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
