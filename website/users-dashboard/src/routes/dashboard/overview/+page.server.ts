
import { error, redirect } from '@sveltejs/kit';
import User, { getUser,  } from '$lib/users';
import { ApiFailureReason } from '$lib/responseEnvelope.js';

/** @type {import('./$types').PageServerLoad} */
export async function load({ params }) {
	console.log("params: " + JSON.stringify(params));

	// TODO: Replace this by loading data from the cookie?
	const DUMMY_API_KEY = 'your-key';
	const DUMMY_USER_ID = 'your-id';
	const userResponse = await getUser(DUMMY_API_KEY, DUMMY_USER_ID);

	// https://kit.svelte.dev/docs/errors
	if (userResponse.error()) {
		const reason = userResponse.failureReason();

		switch (reason) {
			case ApiFailureReason.API_KEY_EXPIRED:
				redirect(302, '/dashboard/login');
		}

		error(404, { message: userResponse.failureMessage() });
	}

	// https://www.okupter.com/blog/sveltekit-cannot-stringify-arbitrary-non-pojos-error
	const user = new User(userResponse);
	return {
		...user,
		apiKey: DUMMY_API_KEY,
	};
}
