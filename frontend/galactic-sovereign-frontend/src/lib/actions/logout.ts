import { type RequestEvent, redirect } from '@sveltejs/kit';
import { loadSessionCookiesOrRedirectToLogin } from '$lib/cookies';
import { logoutUser } from '$lib/service/sessions';
import { HttpStatus } from '@totocorpsoftwareinc/frontend-toolkit';
import { handleApiError } from '$lib/rest/api';

// https://www.reddit.com/r/sveltejs/comments/185585c/how_to_share_pageservertslogic_to_multiple_pages/?share_id=HuVFD5EAH469JAbtSW-mH&utm_content=2&utm_medium=android_app&utm_name=androidcss&utm_source=share&utm_term=1
export const logout = async ({ cookies }: RequestEvent) => {
	const sessionCookies = loadSessionCookiesOrRedirectToLogin(cookies);

	const apiResponse = await logoutUser(sessionCookies.apiKey, sessionCookies.apiUser);
	handleApiError(apiResponse);

	redirect(HttpStatus.SEE_OTHER, '/login');
};
