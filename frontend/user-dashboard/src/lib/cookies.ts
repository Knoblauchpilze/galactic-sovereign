import { type Cookies, redirect } from '@sveltejs/kit';
import { ApiKey } from '$lib/sessions';

const DEFAULT_COOKIES_OPT = {
	path: '/'
};

const COOKIE_KEY_API_USER = 'api-user';
const COOKIE_KEY_API_KEY = 'api-key';

export { COOKIE_KEY_API_USER, COOKIE_KEY_API_KEY };

function validOrEmptyString(maybeValue: string | undefined, valid: boolean): string {
	return valid ? (maybeValue as string) : '';
}

export interface SessionCookies {
	readonly apiUser: string;
	readonly apiKey: string;
}

export function setSessionCookies(cookies: Cookies, apiKey: ApiKey) {
	cookies.set(COOKIE_KEY_API_USER, apiKey.user, DEFAULT_COOKIES_OPT);
	cookies.set(COOKIE_KEY_API_KEY, apiKey.key, DEFAULT_COOKIES_OPT);
}

export function resetSessionCookies(cookies: Cookies) {
	cookies.set(COOKIE_KEY_API_USER, '', DEFAULT_COOKIES_OPT);
	cookies.set(COOKIE_KEY_API_KEY, '', DEFAULT_COOKIES_OPT);
}

export function loadSessionCookies(cookies: Cookies): [boolean, SessionCookies] {
	const maybeApiUser = cookies.get(COOKIE_KEY_API_USER);
	const maybeApiKey = cookies.get(COOKIE_KEY_API_KEY);

	const validApiUser = maybeApiUser !== undefined && maybeApiUser !== '';
	const validApiKey = maybeApiKey !== undefined && maybeApiKey !== '';
	const valid = validApiUser && validApiKey;

	const out: SessionCookies = {
		apiUser: validOrEmptyString(maybeApiUser, validApiUser),
		apiKey: validOrEmptyString(maybeApiKey, validApiKey)
	};

	return [valid, out];
}

export function loadSessionCookiesOrRedirectToLogin(cookies: Cookies): SessionCookies {
	const [valid, sessionCookies] = loadSessionCookies(cookies);
	if (!valid) {
		redirect(303, '/login');
	}

	return sessionCookies;
}
