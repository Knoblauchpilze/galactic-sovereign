
import ResponseEnvelope from './responseEnvelope';
import {createFailedResponseEnvelope} from './responseEnvelope';

const port : number = 60001;
const baseUrl: string = "http://localhost:{port}/v1/users";

function buildUrl(url: string): string {
	return baseUrl.replace("{port}", port.toString()) + "/" + url;
}

const genericFailureReason : string = "Unknown failure";

function analyzeFetchFailureReasone(reason: object): Response {
	// https://developer.mozilla.org/en-US/docs/Web/API/Response/Response
	let failureReason = genericFailureReason;
	if (reason instanceof TypeError) {
		failureReason = (reason as TypeError).message;
	}

	const responseEnvelope = createFailedResponseEnvelope((failureReason as unknown) as object);

	const body = new Blob([JSON.stringify(responseEnvelope)]);
	return new Response(body);
}

async function safeFetch(url: URL | RequestInfo, init?: RequestInit | undefined): Promise<Response> {
	return await fetch(url, init).catch(reason => analyzeFetchFailureReasone(reason));
}

export async function login(email: string, password: string): Promise<ResponseEnvelope> {
	const url = buildUrl("sessions");
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

export async function logout(userId: string): Promise<ResponseEnvelope> {
	const url = buildUrl("sessions/" + userId);

	const params = {
		method: 'DELETE'
	};

	const response = await safeFetch(url, params);
	const jsonContent = await response.json();

	return new ResponseEnvelope(jsonContent);
}