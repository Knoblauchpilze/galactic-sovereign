import { error, redirect } from '@sveltejs/kit';
import { resetGameCookies, setGameCookies, loadSessionCookies } from '$lib/cookies';
import {
	fetchPlayerFromApiUser,
	responseToPlayerArray,
	mapPlayersToUiPlayers
} from '$lib/game/players';
import { getUniverses, responseToUniverseArray } from '$lib/game/universes';
import { fetchPlanetsFromPlayer, responseToPlanetArray } from '$lib/game/planets';

export async function load({ cookies }) {
	resetGameCookies(cookies);

	const [valid, sessionCookies] = loadSessionCookies(cookies);
	if (!valid) {
		redirect(303, '/login');
	}

	const universesResponse = await getUniverses(sessionCookies.apiKey);
	if (universesResponse.error()) {
		error(404, { message: universesResponse.failureMessage() });
	}
	const universes = responseToUniverseArray(universesResponse);

	const playerResponse = await fetchPlayerFromApiUser(
		sessionCookies.apiUser,
		sessionCookies.apiKey
	);
	if (playerResponse.error()) {
		error(404, { message: playerResponse.failureMessage() });
	}

	const players = responseToPlayerArray(playerResponse);

	return {
		players: mapPlayersToUiPlayers(players, universes)
	};
}

export const actions = {
	login: async ({ cookies, request }) => {
		const [valid, sessionCookies] = loadSessionCookies(cookies);
		if (!valid) {
			redirect(303, '/login');
		}

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
		if (playerResponse.error()) {
			return {
				message: playerResponse.failureMessage()
			};
		}

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
		if (planetsResponse.error()) {
			return {
				message: planetsResponse.failureMessage()
			};
		}

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
