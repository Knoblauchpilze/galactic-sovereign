/** @type {import('./$types').PageServerLoad} */
export async function load({ cookies }) {
	cookies.set('api-key', '', { path: '/' });
	cookies.set('api-user', '', { path: '/' });
	cookies.set('player-id', '', { path: '/' });
}
