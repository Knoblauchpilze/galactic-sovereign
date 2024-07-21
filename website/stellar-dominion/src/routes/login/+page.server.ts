import { error, redirect } from '@sveltejs/kit';
import { fetchPlayerFromApiUser, responseToPlayerArray } from '$lib/players';
import ApiKey, { loginUser } from '$lib/sessions';
import Universe, { getUniverses } from '$lib/universes';

/** @type {import('./$types').PageServerLoad} */
export async function load({ cookies }) {
	cookies.set('api-key', '', { path: '/' });
	cookies.set('api-user', '', { path: '/' });
	cookies.set('player-id', '', { path: '/' });

	const universesResponse = await getUniverses();

	if (universesResponse.error()) {
		error(404, { message: universesResponse.failureMessage() });
	}

	const universesJson = universesResponse.getDetails();

	if (!Array.isArray(universesJson)) {
		error(404, { message: 'Failed to fetch universes' });
	}

	const universes = new Array<Universe>();
	for (const maybeUniverse of universesJson) {
		const universe = new Universe(maybeUniverse);
		universes.push(universe);
	}

	const universesData = universes.map((u) => ({
		id: u.id,
		name: u.name
	}));

	return {
		universes: universesData
	};
}

/** @type {import('./$types').Actions} */
export const actions = {
	login: async ({ cookies, request }) => {
		const data = await request.formData();

		const universeId = data.get('universe');
		const email = data.get('email');
		const password = data.get('password');
		const playerName = data.get('player');
		if (!universeId) {
			return {
				success: false,
				missing: true,
				message: 'Please select a universe',

				universeId
			};
		}
		if (!email) {
			return {
				success: false,
				missing: true,
				message: 'Please fill in the email',

				email
			};
		}
		if (!password) {
			return {
				success: false,
				missing: true,
				message: 'Please fill in the password',

				email
			};
		}
		if (!playerName) {
			return {
				success: false,
				missing: true,
				message: 'Please choose a name',

				email
			};
		}

		const loginResponse = await loginUser(email as string, password as string);
		if (loginResponse.error()) {
			return {
				success: false,
				incorrect: true,
				message: loginResponse.failureMessage(),

				email
			};
		}

		const apiKey = new ApiKey(loginResponse);

		const playerResponse = await fetchPlayerFromApiUser(apiKey.user, apiKey.key);
		if (playerResponse.error()) {
			return {
				success: false,
				incorrect: true,
				message: playerResponse.failureMessage(),

				email
			};
		}

		const players = responseToPlayerArray(playerResponse);
		const maybePlayer = players.find(
			(player) => player.universe === universeId && player.name === playerName
		);
		// https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Array/find
		if (maybePlayer === undefined) {
			return {
				success: false,
				incorrect: true,
				message: 'No such player',

				email
			};
		}

		console.log('player: ', JSON.stringify(maybePlayer));

		const opts = {
			path: '/'
		};
		cookies.set('api-user', apiKey.user, opts);
		cookies.set('api-key', apiKey.key, opts);
		cookies.set('player-id', maybePlayer.id, opts);

		// TODO: Should fetch the planet and redirect to the overview
		redirect(303, '/dashboard/overview');
	}
};
