import { type Cookies } from '@sveltejs/kit';
import { ApiKey } from '$lib/sessions';
import { Player } from '$lib/players';

const DEFAULT_COOKIES_OPT = {
	path: '/'
};

const COOKIE_KEY_API_USER = 'api-user';
const COOKIE_KEY_API_KEY = 'api-key';
const COOKIE_KEY_PLAYER_ID = 'player-id';
const COOKIE_KEY_UNIVERSE_ID = 'universe-id';

export function resetCookies(cookies: Cookies) {
	cookies.set(COOKIE_KEY_API_USER, '', DEFAULT_COOKIES_OPT);
	cookies.set(COOKIE_KEY_API_KEY, '', DEFAULT_COOKIES_OPT);
	cookies.set(COOKIE_KEY_PLAYER_ID, '', DEFAULT_COOKIES_OPT);
	cookies.set(COOKIE_KEY_UNIVERSE_ID, '', DEFAULT_COOKIES_OPT);
}

export function setCookies(cookies: Cookies, apiKey: ApiKey, player: Player) {
	cookies.set(COOKIE_KEY_API_USER, apiKey.user, DEFAULT_COOKIES_OPT);
	cookies.set(COOKIE_KEY_API_KEY, apiKey.key, DEFAULT_COOKIES_OPT);
	cookies.set(COOKIE_KEY_PLAYER_ID, player.id, DEFAULT_COOKIES_OPT);
	cookies.set(COOKIE_KEY_UNIVERSE_ID, player.universe, DEFAULT_COOKIES_OPT);
}
