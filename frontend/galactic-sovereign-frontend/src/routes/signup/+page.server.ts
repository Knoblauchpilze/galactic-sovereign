import { redirect } from '@sveltejs/kit';
import { resetCookies } from '$lib/cookies';
import { createUser } from '$lib/users';

export async function load({ cookies }) {
	resetCookies(cookies);
}

export const actions = {
	signup: async ({ cookies, request }) => {
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

		const signupResponse = await createUser(email as string, password as string);
		if (signupResponse.error()) {
			return {
				success: false,
				incorrect: true,
				message: signupResponse.failureMessage(),

				email
			};
		}

		resetCookies(cookies);

		redirect(303, '/login');
	}
};
