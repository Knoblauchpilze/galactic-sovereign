import { fail, redirect } from '@sveltejs/kit';
import { resetSessionCookies } from '$lib/cookies';
import { createUser } from '$lib/service/users';
import {
	getHttpStatusCodeFromApiFailure,
	HttpStatus,
	tryGetFailureReason
} from '@totocorpsoftwareinc/frontend-toolkit';
import { getErrorMessageFromApiResponse } from '$lib/rest/api';

export async function load({ cookies }) {
	resetSessionCookies(cookies);
}

export const actions = {
	signup: async ({ cookies, request }) => {
		const data = await request.formData();

		const email = data.get('email');
		const password = data.get('password');
		if (!email) {
			return fail(HttpStatus.UNPROCESSABLE_ENTITY, {
				message: 'Please fill in the email',
				email: email
			});
		}
		if (!password) {
			return fail(HttpStatus.UNPROCESSABLE_ENTITY, {
				message: 'Please fill in the password',
				email: email
			});
		}

		const apiResponse = await createUser(email as string, password as string);
		if (apiResponse.isError()) {
			const failure = tryGetFailureReason(apiResponse);
			const code = getHttpStatusCodeFromApiFailure(failure);

			return fail(code, {
				message: getErrorMessageFromApiResponse(apiResponse),
				email: email
			});
		}

		resetSessionCookies(cookies);

		redirect(HttpStatus.SEE_OTHER, '/login');
	}
};
