import { error, redirect } from '@sveltejs/kit';
import { registerPlayer } from '$lib/players';
import Universe, { getUniverses } from '$lib/universes';

/** @type {import('./$types').PageServerLoad} */
export async function load({ cookies }) {
	cookies.set('api-key', '', { path: '/' });
	cookies.set('api-user', '', { path: '/' });

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
	signup: async ({ cookies, request }) => {
		const data = await request.formData();

		const email = data.get('email');
		const password = data.get('password');
		const universeId = data.get('universe');
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
		if (!universeId) {
			return {
				success: false,
				missing: true,
				message: 'Please select a universe',

				email
			};
		}
		console.log('universe: ', universeId);

		const playerResponse = await registerPlayer(
			email as string,
			password as string,
			universeId as string,
			'some-random-name'
		);
		if (playerResponse.error()) {
			return {
				success: false,
				incorrect: true,
				message: playerResponse.failureMessage(),

				email
			};
		}

		cookies.set('api-user', '', { path: '/' });
		cookies.set('api-key', '', { path: '/' });

		redirect(303, '/login');
	}
};
