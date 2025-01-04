import { error, fail, redirect } from '@sveltejs/kit';
import {
	resetGameCookies,
	setGameCookies,
	loadSessionCookiesOrRedirectToLogin
} from '$lib/cookies';
import {
	HttpStatus,
	parseApiResponseAsArray,
	parseApiResponseAsSingleValue
} from '@totocorpsoftwareinc/frontend-toolkit';
import { getUniverses } from '$lib/service/universes';
import { createPlayer } from '$lib/service/players';
import {
	getErrorMessageFromApiResponse,
	handleApiError,
	redirectToLoginIfNeeded
} from '$lib/rest/api';
import { UniverseResponseDto } from '$lib/communication/api/universeResponseDto';
import { PlayerResponseDto } from '$lib/communication/api/playerResponseDto';
import { universeResponseDtoToUniverseUiDto } from '$lib/converter/universeConverter';
import { fetchPlayerFromApiUser } from '$lib/service/players';
import {
	tryGetFailureReason,
	getHttpStatusCodeFromApiFailure
} from '@totocorpsoftwareinc/frontend-toolkit';
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
	const players = parseApiResponseAsArray(apiResponse, PlayerResponseDto);

	const universesWithAccount = players.map((p) => p.universe);
	// https://stackoverflow.com/questions/33577868/filter-array-not-in-another-array
	const universesWithoutAccount = universes.filter((u) => !universesWithAccount.includes(u.id));

	return {
		universes: universesWithoutAccount.map((u) => universeResponseDtoToUniverseUiDto(u))
	};
}

export const actions = {
	register: async ({ cookies, request }) => {
		const sessionCookies = loadSessionCookiesOrRedirectToLogin(cookies);

		const data = await request.formData();

		const universeId = data.get('universe');
		const playerName = data.get('player');
		if (!universeId) {
			return fail(HttpStatus.UNPROCESSABLE_ENTITY, {
				message: 'Please select a universe',
				player: playerName
			});
		}
		if (!playerName) {
			return fail(HttpStatus.UNPROCESSABLE_ENTITY, {
				message: 'Please choose a name',
				player: ''
			});
		}

		let apiResponse = await createPlayer(
			sessionCookies.apiUser,
			universeId as string,
			playerName as string,
			sessionCookies.apiKey
		);
		if (apiResponse.isError()) {
			const failure = tryGetFailureReason(apiResponse);
			const code = getHttpStatusCodeFromApiFailure(failure);

			return fail(code, {
				message: getErrorMessageFromApiResponse(apiResponse),
				player: playerName
			});
		}

		const playerDto = parseApiResponseAsSingleValue(apiResponse, PlayerResponseDto);
		if (playerDto === undefined) {
			error(HttpStatus.INTERNAL_SERVER_ERROR, 'Failed to get server data');
		}

		apiResponse = await fetchPlanetsFromPlayer(playerDto.id, sessionCookies.apiKey);
		handleApiError(apiResponse);
		const planets = parseApiResponseAsArray(apiResponse, PlanetResponseDto);

		if (planets.length === 0) {
			error(HttpStatus.INTERNAL_SERVER_ERROR, {
				message: 'Player does not have any planet'
			});
		}

		setGameCookies(cookies, playerDto);

		redirect(HttpStatus.SEE_OTHER, '/planets/' + planets[0].id + '/overview');
	}
};
