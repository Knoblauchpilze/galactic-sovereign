import { type RequestEvent, redirect } from '@sveltejs/kit';
import { logoutUser } from '$lib/sessions';
import { loadSessionCookiesOrRedirectToLogin } from '$lib/cookies';

export const logout = async ({ cookies }: RequestEvent) => {
	const sessionCookies = loadSessionCookiesOrRedirectToLogin(cookies);

	const logoutResponse = await logoutUser(sessionCookies.apiKey, sessionCookies.apiUser);

	if (logoutResponse.error()) {
		return {
			success: false,
			message: logoutResponse.failureMessage()
		};
	}

	redirect(303, '/login');
};
