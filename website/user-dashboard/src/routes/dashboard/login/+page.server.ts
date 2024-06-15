
import { redirect } from '@sveltejs/kit';
// https://learn.svelte.dev/tutorial/lib
import { loginUser } from '$lib/sessions';
import ApiKey from '$lib/apiKey.js';

/** @type {import('./$types').PageServerLoad} */
export async function load({ cookies }) {
	cookies.set('api-key', '', {path: '/'});
	cookies.set('api-user', '', {path: '/'});
};

/** @type {import('./$types').Actions} */
export const actions = {
	login: async ({ cookies, request }) => {
		const data = await request.formData();

		const email = data.get('email');
		const password = data.get('password');
		if (!email) {
			return {
				success: false,
				missing: true,
				message: 'Please fill in the email',

				email,
			}
		}
		if (!password) {
			return {
				success: false,
				missing: true,
				message: 'Please fill in the password',

				email,
			}
		}

		const loginResponse = await loginUser(email as string, password as string);

		if (loginResponse.error()) {
			return {
				success: false,
				incorrect: true,
				message: loginResponse.failureMessage(),

				email,
			}
		}

		const apiKey = new ApiKey(loginResponse);

		const opts = {
			path: '/',
		};
		cookies.set('api-user', apiKey.user, opts);
		cookies.set('api-key', apiKey.key, opts);

		redirect(303, '/dashboard/overview');
	},
};