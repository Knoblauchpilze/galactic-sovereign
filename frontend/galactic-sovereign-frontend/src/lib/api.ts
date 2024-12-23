import { createFailedResponseEnvelope } from '$lib/responseEnvelope';
import { PUBLIC_API_BASE_URL } from '$env/static/public';
import { PUBLIC_USER_API_BASE_URL } from '$env/static/public';

function trimTrailingSlash(url: string): string {
	if (!url.endsWith('/')) {
		return url;
	}

	return url.substring(0, url.length - 1);
}

export function buildUrl(url: string): string {
	const out = trimTrailingSlash(PUBLIC_API_BASE_URL);

	if (url.length === 0) {
		return out;
	}
	return out + '/' + url;
}

export function buildUserUrl(url: string): string {
	const out = trimTrailingSlash(PUBLIC_USER_API_BASE_URL);

	if (url.length === 0) {
		return out;
	}
	return out + '/' + url;
}

const genericFailureReason: string = 'Unknown failure';

function analyzeFetchFailureReason(reason: object): Response {
	let failureReason = genericFailureReason;
	if (reason instanceof TypeError) {
		failureReason = (reason as TypeError).message;
	}

	const responseEnvelope = createFailedResponseEnvelope(failureReason as unknown as object);

	const body = new Blob([JSON.stringify(responseEnvelope)]);
	return new Response(body);
}

export async function safeFetch(
	url: URL | RequestInfo,
	init?: RequestInit | undefined
): Promise<Response> {
	return await fetch(url, init).catch((reason) => analyzeFetchFailureReason(reason));
}
