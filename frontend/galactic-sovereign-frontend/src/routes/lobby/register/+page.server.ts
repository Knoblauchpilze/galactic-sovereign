import { error, redirect } from '@sveltejs/kit';
import { resetGameCookies, setGameCookies, loadSessionCookies } from '$lib/cookies';
import { Player, createPlayer } from '$lib/game/players';
import { getUniverses, responseToUniverseArray } from '$lib/game/universes';
import { fetchPlanetsFromPlayer, responseToPlanetArray } from '$lib/game/planets';

export async function load({ cookies }) {
	resetGameCookies(cookies);

	const universesResponse = await getUniverses();

	if (universesResponse.error()) {
		error(404, { message: universesResponse.failureMessage() });
	}

	const universes = responseToUniverseArray(universesResponse);

	return {
		universes: universes.map((u) => u.toJson())
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
				success: false,
				missing: true,
				message: 'Please select a universe'
			};
		}
		if (!playerName) {
			return {
				success: false,
				missing: true,
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
				success: false,
				incorrect: true,
				message: planetsResponse.failureMessage()
			};
		}

		const planets = responseToPlanetArray(planetsResponse);

		const [maybePlanet] = planets;

		if (maybePlanet === undefined) {
			return {
				success: false,
				incorrect: true,
				message: 'Player does not have any planet'
			};
		}

		setGameCookies(cookies, player);

		redirect(303, '/planets/' + maybePlanet.id + '/overview');
	}
};
