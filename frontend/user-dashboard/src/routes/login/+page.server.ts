import { error, redirect } from '@sveltejs/kit';
// https://learn.svelte.dev/tutorial/lib
import { resetSessionCookies, setSessionCookies } from '$lib/cookies';
import { loginUser } from '$lib/service/sessions';
import { ApiKeyResponseDto } from '$lib/communication/api/apiKeyResponseDto';
import {
	HttpStatus,
	parseApiResponseAsSingleValue,
	tryGetFailureReason
} from '@totocorpsoftwareinc/frontend-toolkit';

export async function load({ cookies }) {
	resetSessionCookies(cookies);
}

export const actions = {
	login: async ({ cookies, request }) => {
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

		const apiResponse = await loginUser(email as string, password as string);
		if (apiResponse.isError()) {
			return {
				message: tryGetFailureReason(apiResponse),
				email: email
			};
		}

		const apiKeyDto = parseApiResponseAsSingleValue(apiResponse, ApiKeyResponseDto);
		if (apiKeyDto === undefined) {
			error(HttpStatus.NOT_FOUND, 'Failed to get login data');
		}

		setSessionCookies(cookies, apiKeyDto);

		redirect(HttpStatus.SEE_OTHER, '/overview');
	}
};
