import { type Cookies } from '@sveltejs/kit';
import { ApiKey } from '$lib/sessions';
import { Player } from '$lib/game/players';

const DEFAULT_COOKIES_OPT = {
	path: '/'
};

const COOKIE_KEY_API_USER = 'api-user';
const COOKIE_KEY_API_KEY = 'api-key';
const COOKIE_KEY_PLAYER_ID = 'player-id';
const COOKIE_KEY_PLAYER_NAME = 'player-name';
const COOKIE_KEY_UNIVERSE_ID = 'universe-id';

export {
	COOKIE_KEY_API_USER,
	COOKIE_KEY_API_KEY,
	COOKIE_KEY_PLAYER_ID,
	COOKIE_KEY_PLAYER_NAME,
	COOKIE_KEY_UNIVERSE_ID
};

export function setSessionCookies(cookies: Cookies, apiKey: ApiKey) {
	cookies.set(COOKIE_KEY_API_USER, apiKey.user, DEFAULT_COOKIES_OPT);
	cookies.set(COOKIE_KEY_API_KEY, apiKey.key, DEFAULT_COOKIES_OPT);
}

export function resetSessionCookies(cookies: Cookies) {
	cookies.set(COOKIE_KEY_API_USER, '', DEFAULT_COOKIES_OPT);
	cookies.set(COOKIE_KEY_API_KEY, '', DEFAULT_COOKIES_OPT);
}

export function setGameCookies(cookies: Cookies, player: Player) {
	cookies.set(COOKIE_KEY_PLAYER_ID, player.id, DEFAULT_COOKIES_OPT);
	cookies.set(COOKIE_KEY_PLAYER_NAME, player.name, DEFAULT_COOKIES_OPT);
	cookies.set(COOKIE_KEY_UNIVERSE_ID, player.universe, DEFAULT_COOKIES_OPT);
}

export function resetGameCookies(cookies: Cookies) {
	cookies.set(COOKIE_KEY_PLAYER_ID, '', DEFAULT_COOKIES_OPT);
	cookies.set(COOKIE_KEY_PLAYER_NAME, '', DEFAULT_COOKIES_OPT);
	cookies.set(COOKIE_KEY_UNIVERSE_ID, '', DEFAULT_COOKIES_OPT);
}

export function setCookies(cookies: Cookies, apiKey: ApiKey, player: Player) {
	setSessionCookies(cookies, apiKey);
	setGameCookies(cookies, player);
}

export function resetCookies(cookies: Cookies) {
	resetSessionCookies(cookies);
	resetGameCookies(cookies);
}

function validOrEmptyString(maybeValue: string | undefined, valid: boolean): string {
	return valid ? (maybeValue as string) : '';
}

export interface SessionCookies {
	readonly apiUser: string;
	readonly apiKey: string;
}

export function loadSessionCookies(cookies: Cookies): [boolean, SessionCookies] {
	const maybeApiUser = cookies.get(COOKIE_KEY_API_USER);
	const maybeApiKey = cookies.get(COOKIE_KEY_API_KEY);

	const validApiUser = maybeApiUser !== undefined;
	const validApiKey = maybeApiKey !== undefined;
	const valid = validApiUser || validApiKey;

	const out: SessionCookies = {
		apiUser: validOrEmptyString(maybeApiUser, validApiUser),
		apiKey: validOrEmptyString(maybeApiKey, validApiKey)
	};

	return [valid, out];
}

export interface GameCookies {
	readonly session: SessionCookies;
	readonly playerId: string;
	readonly playerName: string;
	readonly universeId: string;
}

export function loadCookies(cookies: Cookies): [boolean, GameCookies] {
	const maybePlayerId = cookies.get(COOKIE_KEY_PLAYER_ID);
	const maybePlayerName = cookies.get(COOKIE_KEY_PLAYER_NAME);
	const maybeUniverseId = cookies.get(COOKIE_KEY_UNIVERSE_ID);

	const [validSession, session] = loadSessionCookies(cookies);

	const validPlayerId = maybePlayerId !== undefined;
	const validPlayerName = maybePlayerName !== undefined;
	const validUniverseId = maybeUniverseId !== undefined;
	const valid = validSession || validPlayerId || validUniverseId;

	const out: GameCookies = {
		session: session,
		playerId: validOrEmptyString(maybePlayerId, validPlayerId),
		playerName: validOrEmptyString(maybePlayerName, validPlayerName),
		universeId: validOrEmptyString(maybeUniverseId, validUniverseId)
	};

	return [valid, out];
}
