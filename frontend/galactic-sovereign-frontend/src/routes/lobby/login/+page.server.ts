import { error, fail, redirect } from '@sveltejs/kit';
import {
	resetGameCookies,
	setGameCookies,
	loadSessionCookiesOrRedirectToLogin
} from '$lib/cookies';
import { HttpStatus, parseApiResponseAsArray } from '@totocorpsoftwareinc/frontend-toolkit';
import { getUniverses } from '$lib/service/universes';
import { handleApiError, redirectToLoginIfNeeded } from '$lib/rest/api';
import { UniverseResponseDto } from '$lib/communication/api/universeResponseDto';
import { fetchPlayerFromApiUser } from '$lib/service/players';
import { PlayerResponseDto } from '$lib/communication/api/playerResponseDto';
import { playerResponseDtoToPlayerUiDto } from '$lib/converter/playerConverter';
import { fetchPlanetsFromPlayer } from '$lib/service/planets';
import { PlanetResponseDto } from '$lib/communication/api/planetResponseDto';

export async function load({ cookies }) {
	resetGameCookies(cookies);

	const sessionCookies = loadSessionCookiesOrRedirectToLogin(cookies);

	let apiResponse = await getUniverses(sessionCookies.apiKey);
	redirectToLoginIfNeeded(apiResponse);
	handleApiError(apiResponse);

	const universes = parseApiResponseAsArray(apiResponse, UniverseResponseDto);
	if (universes.length === 0) {
		error(HttpStatus.INTERNAL_SERVER_ERROR, 'Failed to get server data');
	}

	apiResponse = await fetchPlayerFromApiUser(sessionCookies.apiUser, sessionCookies.apiKey);
	handleApiError(apiResponse);
	const apiPlayers = parseApiResponseAsArray(apiResponse, PlayerResponseDto);

	const uiPlayers = apiPlayers.map((p) => {
		const maybeUniverse = universes.find((u) => u.id === p.universe);
		return playerResponseDtoToPlayerUiDto(p, maybeUniverse);
	});

	return {
		players: uiPlayers
	};
}

export const actions = {
	login: async ({ cookies, request }) => {
		const sessionCookies = loadSessionCookiesOrRedirectToLogin(cookies);

		const data = await request.formData();

		const universeId = data.get('universe');
		const playerName = data.get('player');
		if (!universeId) {
			return fail(HttpStatus.UNPROCESSABLE_ENTITY, {
				message: 'Please select a universe'
			});
		}
		if (!playerName) {
			return fail(HttpStatus.UNPROCESSABLE_ENTITY, {
				message: 'Please choose a name'
			});
		}

		let apiResponse = await fetchPlayerFromApiUser(sessionCookies.apiUser, sessionCookies.apiKey);
		redirectToLoginIfNeeded(apiResponse);
		handleApiError(apiResponse);

		const players = parseApiResponseAsArray(apiResponse, PlayerResponseDto);

		const maybePlayer = players.find(
			(player) => player.universe === universeId && player.name === playerName
		);
		// https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Array/find
		if (maybePlayer === undefined) {
			error(HttpStatus.NOT_FOUND, {
				message: 'Player does not seem to exist'
			});
		}

		apiResponse = await fetchPlanetsFromPlayer(maybePlayer.id, sessionCookies.apiKey);
		handleApiError(apiResponse);
		const planets = parseApiResponseAsArray(apiResponse, PlanetResponseDto);

		if (planets.length === 0) {
			error(HttpStatus.INTERNAL_SERVER_ERROR, {
				message: 'Player does not have any planet'
			});
		}

		setGameCookies(cookies, maybePlayer);

		redirect(HttpStatus.SEE_OTHER, '/planets/' + planets[0].id + '/overview');
	}
};
