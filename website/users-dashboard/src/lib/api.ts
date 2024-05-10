
import {createFailedResponseEnvelope} from './responseEnvelope';

const port : number = 60001;
const baseUrl: string = "http://localhost:{port}/v1/users";

export function buildUrl(url: string): string {
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

export async function safeFetch(url: URL | RequestInfo, init?: RequestInit | undefined): Promise<Response> {
	return await fetch(url, init).catch(reason => analyzeFetchFailureReasone(reason));
}