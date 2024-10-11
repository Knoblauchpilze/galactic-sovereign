import { resetCookies } from '$lib/cookies';

export async function load({ cookies }) {
	resetCookies(cookies);
}
