
import {createFailedResponseEnvelope} from './responseEnvelope';
import { PUBLIC_API_BASE_URL } from '$env/static/public';

// TODO: Check how to trim it of the last '/' character
export function buildUrl(url: string): string {
	if (url.length === 0) {
		return PUBLIC_API_BASE_URL;
	}
	return PUBLIC_API_BASE_URL + "/" + url;
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