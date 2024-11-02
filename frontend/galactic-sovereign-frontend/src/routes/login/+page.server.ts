import { redirect } from '@sveltejs/kit';
import { resetAllCookies, setSessionCookies } from '$lib/cookies';
import { ApiKey, loginUser } from '$lib/sessions';

export async function load({ cookies }) {
	resetAllCookies(cookies);
}

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

		const loginResponse = await loginUser(email as string, password as string);
		if (loginResponse.error()) {
			return {
				success: false,
				incorrect: true,
				message: loginResponse.failureMessage(),

				email
			};
		}

		const apiKey = new ApiKey(loginResponse);

		setSessionCookies(cookies, apiKey);

		redirect(303, '/lobby');
	}
};
