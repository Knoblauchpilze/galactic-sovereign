
import { redirect } from '@sveltejs/kit';
// https://learn.svelte.dev/tutorial/lib
import { login } from '$lib/sessions';
import ApiKey from '$lib/apiKey.js';

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

		const loginResponse = await login(email as string, password as string);

		if (loginResponse.error()) {
			return {
				success: false,
				incorrect: true,
				message: loginResponse.failureMessage(),

				email,
			}
		}

		const apiKey = new ApiKey(loginResponse);

		cookies.set('api-user', apiKey.user, { path: '/' });
		cookies.set('api-key', apiKey.key, { path: '/' });

		redirect(302, '/dashboard/overview');

		return {
			success: true,
			message: ''
		};
	},
};