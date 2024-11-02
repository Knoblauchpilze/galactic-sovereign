import { type RequestEvent, redirect } from '@sveltejs/kit';
import { logoutUser } from '$lib/sessions';
import { loadSessionCookiesOrRedirectToLogin } from '$lib/cookies';

// https://www.reddit.com/r/sveltejs/comments/185585c/how_to_share_pageservertslogic_to_multiple_pages/?share_id=HuVFD5EAH469JAbtSW-mH&utm_content=2&utm_medium=android_app&utm_name=androidcss&utm_source=share&utm_term=1
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
