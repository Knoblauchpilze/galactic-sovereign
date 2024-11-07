import { logout } from '$lib/actions/logout';

import { User, getUser } from '$lib/users';
import { analayzeResponseEnvelopAndRedirectIfNeeded } from '$lib/responseEnvelope.js';
import { loadSessionCookiesOrRedirectToLogin } from '$lib/cookies.js';

export async function load({ cookies }) {
	const sessionCookies = loadSessionCookiesOrRedirectToLogin(cookies);

	const userResponse = await getUser(sessionCookies.apiKey, sessionCookies.apiUser);
	analayzeResponseEnvelopAndRedirectIfNeeded(userResponse);

	const user = new User(userResponse);
	return {
		user: user.toJson(),
		apiKey: sessionCookies.apiKey
	};
}

export const actions = {
	logout: logout
};
