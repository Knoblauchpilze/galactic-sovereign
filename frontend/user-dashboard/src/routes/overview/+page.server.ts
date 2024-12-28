import { logout } from '$lib/actions/logout';

import { getUser } from '$lib/service/users';
import { analyzeApiResponseAndRedirectIfNeeded } from '$lib/rest/api';
import { loadSessionCookiesOrRedirectToLogin } from '$lib/cookies';
import { HttpStatus, parseApiResponseAsSingleValue } from '@totocorpsoftwareinc/frontend-toolkit';
import { UserResponseDto } from '$lib/communication/api/userResponseDto';
import { userResponseDtoToUserUiDto } from '$lib/converter/userConverter';
import { error } from '@sveltejs/kit';

export async function load({ cookies }) {
	const sessionCookies = loadSessionCookiesOrRedirectToLogin(cookies);

	const apiResponse = await getUser(sessionCookies.apiKey, sessionCookies.apiUser);
	analyzeApiResponseAndRedirectIfNeeded(apiResponse);

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
	logout: logout
};
