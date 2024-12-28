import { PUBLIC_API_BASE_URL } from '$env/static/public';
import { error, redirect } from '@sveltejs/kit';
import {
	ApiFailure,
	ApiResponse,
	HttpStatus,
	trimTrailingSlash,
	tryGetFailureReason
} from '@totocorpsoftwareinc/frontend-toolkit';

export function buildUserUrl(url: string): string {
	const out = trimTrailingSlash(PUBLIC_API_BASE_URL);

	if (url.length === 0) {
		return out;
	}
	return out + '/' + url;
}

export function analyzeApiResponseAndRedirectIfNeeded(response: ApiResponse) {
	if (!response.isError()) {
		return;
	}

	const reason = tryGetFailureReason(response);

	switch (reason) {
		case ApiFailure.NOT_AUTHENTICATED:
		case ApiFailure.API_KEY_EXPIRED:
			redirect(HttpStatus.SEE_OTHER, '/login');
	}

	// https://kit.svelte.dev/docs/errors
	error(HttpStatus.NOT_FOUND, { message: 'Request failed with code: ' + reason });
}
