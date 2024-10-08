import { error, redirect } from '@sveltejs/kit';
import { resetCookies } from '$lib/cookies';
import { registerPlayer } from '$lib/game/players';
import { Universe, getUniverses } from '$lib/game/universes';

/** @type {import('./$types').PageServerLoad} */
export async function load({ cookies }) {
	resetCookies(cookies);

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

	return {
		universes: universes.map((u) => u.toJson())
	};
}

export const actions = {
	signup: async ({ cookies, request }) => {
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

		const playerResponse = await registerPlayer(
			email as string,
			password as string,
			universeId as string,
			playerName as string
		);
		if (playerResponse.error()) {
			return {
				success: false,
				incorrect: true,
				message: playerResponse.failureMessage(),

				email
			};
		}

		resetCookies(cookies);

		redirect(303, '/login');
	}
};
