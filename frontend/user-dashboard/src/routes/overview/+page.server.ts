import { logoutUser } from '$lib/service/sessions';
import { getUser } from '$lib/service/users';
import { getErrorMessageFromApiResponse, redirectToLoginIfNeeded } from '$lib/rest/api';
import { loadSessionCookiesOrRedirectToLogin } from '$lib/cookies';
import { HttpStatus, parseApiResponseAsSingleValue } from '@totocorpsoftwareinc/frontend-toolkit';
import { UserResponseDto } from '$lib/communication/api/userResponseDto';
import { userResponseDtoToUserUiDto } from '$lib/converter/userConverter';
import { error, fail, redirect } from '@sveltejs/kit';

export async function load({ cookies }) {
	const sessionCookies = loadSessionCookiesOrRedirectToLogin(cookies);

	const apiResponse = await getUser(sessionCookies.apiKey, sessionCookies.apiUser);
	redirectToLoginIfNeeded(apiResponse);

	const apiUserDto = parseApiResponseAsSingleValue(apiResponse, UserResponseDto);
	if (apiUserDto === undefined) {
		error(HttpStatus.NOT_FOUND, 'Failed to get user data');
	}

	return {
		user: userResponseDtoToUserUiDto(apiUserDto),
		apiKey: sessionCookies.apiKey
	};
}

export const actions = {
	logout: async ({ cookies }) => {
		const sessionCookies = loadSessionCookiesOrRedirectToLogin(cookies);

		const apiResponse = await logoutUser(sessionCookies.apiKey, sessionCookies.apiUser);
		if (apiResponse.isError()) {
			return fail(HttpStatus.UNPROCESSABLE_ENTITY, {
				message: getErrorMessageFromApiResponse(apiResponse)
			});
		}

		redirect(HttpStatus.SEE_OTHER, '/login');
	}
};
