
import ResponseEnvelope from './responseEnvelope';

const port : number = 60001;
const baseUrl: string = "http://localhost:{port}/v1/users";

function buildUrl(url: string): string {
	return baseUrl.replace("{port}", port.toString()) + "/" + url;
}

export default async function login(email: string, password: string): Promise<ResponseEnvelope> {
	const url = buildUrl("sessions");
	const body = JSON.stringify({ email: email, password: password });

	const response = await fetch(url, {
		method: 'POST',
		body: body,
		headers: {
			'content-type': 'application/json'
		}
	});

	const jsonContent = await response.json();

	return new ResponseEnvelope(jsonContent);
}