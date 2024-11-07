import { redirect } from '@sveltejs/kit';
import { resetSessionCookies } from '$lib/cookies';
import { createUser } from '$lib/users';

export async function load({ cookies }) {
	resetSessionCookies(cookies);
}

export const actions = {
	signup: async ({ cookies, request }) => {
		const data = await request.formData();

		const email = data.get('email');
		const password = data.get('password');
		if (!email) {
			return {
				message: 'Please fill in the email',
				email: email
			};
		}
		if (!password) {
			return {
				message: 'Please fill in the password',
				email: email
			};
		}

		const signupResponse = await createUser(email as string, password as string);
		if (signupResponse.error()) {
			return {
				message: signupResponse.failureMessage(),
				email: email
			};
		}

		resetSessionCookies(cookies);

		redirect(303, '/login');
	}
};
