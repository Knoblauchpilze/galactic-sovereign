import { error, redirect } from '@sveltejs/kit';
import { resetGameCookies, setGameCookies, loadSessionCookies } from '$lib/cookies';
import {
	Player,
	createPlayer,
	fetchPlayerFromApiUser,
	responseToPlayerArray
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
	const universesWithAccount = players.map((p) => p.universe);

	// https://stackoverflow.com/questions/33577868/filter-array-not-in-another-array
	const universesWithoutAccount = universes.filter((u) => !universesWithAccount.includes(u.id));

	return {
		universes: universesWithoutAccount.map((u) => u.toJson())
	};
}

export const actions = {
	register: async ({ cookies, request }) => {
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

		const playerResponse = await createPlayer(
			sessionCookies.apiUser,
			universeId as string,
			playerName as string,
			sessionCookies.apiKey
		);
		if (playerResponse.error()) {
			return playerResponse;
		}

		const player = new Player(playerResponse.details);

		const planetsResponse = await fetchPlanetsFromPlayer(player.id, sessionCookies.apiKey);
		if (planetsResponse.error()) {
			return {
				message: planetsResponse.failureMessage()
			};
		}

		const planets = responseToPlanetArray(planetsResponse);

		const [maybePlanet] = planets;

		if (maybePlanet === undefined) {
			return {
				message: 'Player does not have any planet'
			};
		}

		setGameCookies(cookies, player);

		redirect(303, '/planets/' + maybePlanet.id + '/overview');
	}
};
