
import { error, redirect } from '@sveltejs/kit';
import User, { getUser,  } from '$lib/users';
import { ApiFailureReason } from '$lib/responseEnvelope.js';

/** @type {import('./$types').PageServerLoad} */
export async function load({ cookies }) {
	const apiKey = cookies.get('api-key');
	if (!apiKey) {
		redirect(302, '/dashboard/login');
	}

	const apiUser = cookies.get('api-user');
	if (!apiUser) {
		redirect(302, '/dashboard/login');
	}

	const userResponse = await getUser(apiKey, apiUser);

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
		apiKey: apiKey,
	};
};
