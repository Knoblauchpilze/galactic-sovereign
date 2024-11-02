import { redirect } from '@sveltejs/kit';
import {
	resetGameCookies,
	setGameCookies,
	loadSessionCookiesOrRedirectToLogin
} from '$lib/cookies';
import { analayzeResponseEnvelopAndRedirectIfNeeded } from '$lib/responseEnvelope';
import {
	fetchPlayerFromApiUser,
	responseToPlayerArray,
	mapPlayersToUiPlayers
} from '$lib/game/players';
import { getUniverses, responseToUniverseArray } from '$lib/game/universes';
import { fetchPlanetsFromPlayer, responseToPlanetArray } from '$lib/game/planets';

export async function load({ cookies }) {
	resetGameCookies(cookies);

	const sessionCookies = loadSessionCookiesOrRedirectToLogin(cookies);

	const universesResponse = await getUniverses(sessionCookies.apiKey);
	analayzeResponseEnvelopAndRedirectIfNeeded(universesResponse);
	const universes = responseToUniverseArray(universesResponse);

	const playerResponse = await fetchPlayerFromApiUser(
		sessionCookies.apiUser,
		sessionCookies.apiKey
	);
	analayzeResponseEnvelopAndRedirectIfNeeded(playerResponse);
	const players = responseToPlayerArray(playerResponse);

	return {
		players: mapPlayersToUiPlayers(players, universes)
	};
}

export const actions = {
	login: async ({ cookies, request }) => {
		const sessionCookies = loadSessionCookiesOrRedirectToLogin(cookies);

		const data = await request.formData();

		const universeId = data.get('universe');
		const playerName = data.get('player');
		if (!universeId) {
			return {
				message: 'Please select a universe'
			};
		}
		if (!playerName) {
			return {
				message: 'Please choose a name'
			};
		}

		const playerResponse = await fetchPlayerFromApiUser(
			sessionCookies.apiUser,
			sessionCookies.apiKey
		);
		analayzeResponseEnvelopAndRedirectIfNeeded(playerResponse);
		const players = responseToPlayerArray(playerResponse);

		const maybePlayer = players.find(
			(player) => player.universe === universeId && player.name === playerName
		);
		// https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Array/find
		if (maybePlayer === undefined) {
			return {
				message: 'No such player'
			};
		}

		const planetsResponse = await fetchPlanetsFromPlayer(maybePlayer.id, sessionCookies.apiKey);
		analayzeResponseEnvelopAndRedirectIfNeeded(planetsResponse);
		const planets = responseToPlanetArray(planetsResponse);

		// https://stackoverflow.com/questions/35605548/get-first-object-from-array
		const [maybePlanet] = planets;

		if (maybePlanet === undefined) {
			return {
				message: 'Player does not have any planet'
			};
		}

		setGameCookies(cookies, maybePlayer);

		redirect(303, '/planets/' + maybePlanet.id + '/overview');
	}
};
