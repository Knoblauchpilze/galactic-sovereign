import { type RequestEvent, redirect } from '@sveltejs/kit';
import { logoutUser } from '$lib/service/sessions';
import { loadSessionCookiesOrRedirectToLogin } from '$lib/cookies';
import { HttpStatus, tryGetFailureReason } from '@totocorpsoftwareinc/frontend-toolkit';

export const logout = async ({ cookies }: RequestEvent) => {
	const sessionCookies = loadSessionCookiesOrRedirectToLogin(cookies);

	const apiResponse = await logoutUser(sessionCookies.apiKey, sessionCookies.apiUser);

	if (apiResponse.isError()) {
		return {
			success: false,
			message: tryGetFailureReason(apiResponse)
		};
	}

	redirect(HttpStatus.SEE_OTHER, '/login');
};
