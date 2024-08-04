import { resetCookies } from '$lib/cookies';

/** @type {import('./$types').PageServerLoad} */
export async function load({ cookies }) {
	resetCookies(cookies);
}
