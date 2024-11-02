import { resetAllCookies } from '$lib/cookies';

export async function load({ cookies }) {
	resetAllCookies(cookies);
}
