import { redirect, fail } from '@sveltejs/kit';
// https://learn.svelte.dev/tutorial/lib
import { resetSessionCookies, setSessionCookies } from '$lib/cookies';
import { loginUser } from '$lib/service/sessions';
import { ApiKeyResponseDto } from '$lib/communication/api/apiKeyResponseDto';
import { HttpStatus, parseApiResponseAsSingleValue } from '@totocorpsoftwareinc/frontend-toolkit';
import { getErrorMessageFromApiResponse } from '$lib/rest/api';

export async function load({ cookies }) {
	resetSessionCookies(cookies);
}

export const actions = {
	login: async ({ cookies, request }) => {
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

		const apiResponse = await loginUser(email as string, password as string);
		if (apiResponse.isError()) {
			// TODO: Should probably be something else than this code: ideally the
			// code coming from the API response.
			return fail(HttpStatus.UNPROCESSABLE_ENTITY, {
				message: getErrorMessageFromApiResponse(apiResponse),
				email: email
			});
		}

		const apiKeyDto = parseApiResponseAsSingleValue(apiResponse, ApiKeyResponseDto);
		if (apiKeyDto === undefined) {
			fail(HttpStatus.NOT_FOUND, {
				message: 'Failed to get login data',
				email: email
			});
		} else {
			setSessionCookies(cookies, apiKeyDto);
		}

		redirect(HttpStatus.SEE_OTHER, '/overview');
	}
};
