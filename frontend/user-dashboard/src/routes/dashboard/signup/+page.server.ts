import { redirect } from '@sveltejs/kit';
import { createUser } from '$lib/users';

export async function load({ cookies }) {
	cookies.set('api-key', '', { path: '/' });
	cookies.set('api-user', '', { path: '/' });
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

		cookies.set('api-user', '', { path: '/' });
		cookies.set('api-key', '', { path: '/' });

		redirect(303, '/dashboard/login');
	}
};
