import { PUBLIC_API_BASE_URL } from '$env/static/public';
import { PUBLIC_USER_API_BASE_URL } from '$env/static/public';
import { error, redirect } from '@sveltejs/kit';
import {
	ApiFailure,
	ApiResponse,
	HttpStatus,
	trimTrailingSlash,
	tryGetFailureReason,
	getHttpStatusCodeFromApiFailure
} from '@totocorpsoftwareinc/frontend-toolkit';

export function buildUserUrl(url: string): string {
	const out = trimTrailingSlash(PUBLIC_USER_API_BASE_URL);

	if (url.length === 0) {
		return out;
	}
	return out + '/' + url;
}

export function buildUrl(url: string): string {
	const out = trimTrailingSlash(PUBLIC_API_BASE_URL);

	if (url.length === 0) {
		return out;
	}
	return out + '/' + url;
}

export function redirectToLoginIfNeeded(response: ApiResponse) {
	if (!response.isError()) {
		return;
	}

	const reason = tryGetFailureReason(response);

	switch (reason) {
		case ApiFailure.NOT_AUTHENTICATED:
		case ApiFailure.API_KEY_EXPIRED:
			redirect(HttpStatus.SEE_OTHER, '/login');
	}
}

function getErrorMessageFromApiFailure(failure: ApiFailure): string {
	switch (failure) {
		case ApiFailure.SERVICE_UNAVAILABLE:
			return 'Service is currently unavailable, please try again later';
		case ApiFailure.INVALID_REGISTRATION_DATA:
			return 'The registration data is invalid, please provide a valid email and password';
		case ApiFailure.INVALID_CREDENTIALS:
		case ApiFailure.NO_SUCH_USER:
			return 'The email and/or password are invalid';
		default:
			return 'An unexpected error occurred';
	}
}

export function getErrorMessageFromApiResponse(response: ApiResponse): string {
	if (!response.isError()) {
		return '';
	}

	const reason = tryGetFailureReason(response);
	return getErrorMessageFromApiFailure(reason);
}

export function handleApiError(response: ApiResponse) {
	if (!response.isError()) {
		return '';
	}

	const reason = tryGetFailureReason(response);
	const message = getErrorMessageFromApiFailure(reason);
	const code = getHttpStatusCodeFromApiFailure(reason);

	// https://kit.svelte.dev/docs/errors
	error(code, { message: message });
}
