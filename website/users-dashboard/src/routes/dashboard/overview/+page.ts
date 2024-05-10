
import { error } from '@sveltejs/kit';
import User, { getUser,  } from '$lib/users';

/** @type {import('./$types').PageLoad} */
export async function load({ params }) {
	console.log("params: " + JSON.stringify(params));

	// TODO: Replace this by loading data from the cookie?
	const DUMMY_API_KEY = 'your-key';
	const DUMMY_USER_ID = 'your-id';
	const userResponse = await getUser(DUMMY_API_KEY, DUMMY_USER_ID);

	// https://kit.svelte.dev/docs/errors
	if (userResponse.error()) {
		error(404, { message: userResponse.failureReason() });
	}

	return {
		user: new User(userResponse),
	};
}
