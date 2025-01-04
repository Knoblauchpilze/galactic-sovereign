import { fail, type RequestEvent } from '@sveltejs/kit';
import { loadSessionCookiesOrRedirectToLogin } from '$lib/cookies';
import { createBuildingAction, deleteBuildingAction } from '$lib/service/actions';
import { HttpStatus } from '@totocorpsoftwareinc/frontend-toolkit';
import { getErrorMessageFromApiResponse } from '$lib/rest/api';

export const requestCreateBuildingAction = async ({ cookies, request }: RequestEvent) => {
	const sessionCookies = loadSessionCookiesOrRedirectToLogin(cookies);

	const data = await request.formData();

	const building = data.get('building');
	const planet = data.get('planet');

	if (!building) {
		return fail(HttpStatus.UNPROCESSABLE_ENTITY, {
			message: 'Please select a building',
			building: '',
			planet: planet
		});
	}

	if (!planet) {
		return fail(HttpStatus.UNPROCESSABLE_ENTITY, {
			message: 'Please select a planet',
			building: building,
			planet: ''
		});
	}

	const apiResponse = await createBuildingAction(
		sessionCookies.apiKey,
		planet as string,
		building as string
	);
	if (apiResponse.isError()) {
		return {
			message: getErrorMessageFromApiResponse(apiResponse)
		};
	}
};

export const requestDeleteBuildingAction = async ({ cookies, request }: RequestEvent) => {
	const sessionCookies = loadSessionCookiesOrRedirectToLogin(cookies);

	const data = await request.formData();

	const actionId = data.get('action');
	if (!actionId) {
		return fail(HttpStatus.UNPROCESSABLE_ENTITY, {
			message: 'Please select an action'
		});
	}

	const apiResponse = await deleteBuildingAction(sessionCookies.apiKey, actionId as string);
	if (apiResponse.isError()) {
		return {
			message: getErrorMessageFromApiResponse(apiResponse)
		};
	}
};
