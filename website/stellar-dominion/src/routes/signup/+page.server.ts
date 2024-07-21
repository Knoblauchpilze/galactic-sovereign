import { error, redirect } from '@sveltejs/kit';
import { createUser } from '$lib/users';
import Universe, { getUniverses } from '$lib/universes';
import ApiKey from '$lib/apiKey.js';

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

	const universesData = universes.map(u => ({
		id: u.id,
		name: u.name,
	}));

	return {
		universes: universesData,
	};
}

/** @type {import('./$types').Actions} */
export const actions = {
	signup: async ({ cookies, request }) => {
		const data = await request.formData();

		const email = data.get('email');
		const password = data.get('password');
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

		const signupResponse = await createUser(email as string, password as string);

		if (signupResponse.error()) {
			return {
				success: false,
				incorrect: true,
				message: signupResponse.failureMessage(),

				email
			};
		}

		const apiKey = new ApiKey(signupResponse);

		cookies.set('api-user', apiKey.user, { path: '/' });
		cookies.set('api-key', apiKey.key, { path: '/' });

		redirect(303, '/login');
	}
};
