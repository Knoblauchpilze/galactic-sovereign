import { type RequestEvent, redirect } from '@sveltejs/kit';
import { loadSessionCookies } from '$lib/cookies';
import { createBuildingAction, deleteBuildingAction } from '$lib/game/planets';

export const requestCreateBuildingAction = async ({ cookies, request }: RequestEvent) => {
	const [valid, sessionCookies] = loadSessionCookies(cookies);
	if (!valid) {
		redirect(303, '/login');
	}

	const data = await request.formData();

	const buildingId = data.get('building');
	if (!buildingId) {
		return {
			success: false,
			missing: true,
			message: 'Please select a building',

			buildingId
		};
	}

	const planetId = data.get('planet');
	if (!planetId) {
		return {
			success: false,
			missing: true,
			message: 'Please select a planet',

			planetId
		};
	}

	const actionResponse = await createBuildingAction(
		sessionCookies.apiKey,
		planetId as string,
		buildingId as string
	);
	if (actionResponse.error()) {
		return {
			success: false,
			message: actionResponse.failureMessage()
		};
	}
};

export const requestDeleteBuildingAction = async ({ cookies, request }: RequestEvent) => {
	const [valid, sessionCookies] = loadSessionCookies(cookies);
	if (!valid) {
		redirect(303, '/login');
	}

	const data = await request.formData();

	const actionId = data.get('action');
	if (!actionId) {
		return {
			success: false,
			missing: true,
			message: 'Please select an action',

			actionId
		};
	}

	const actionResponse = await deleteBuildingAction(sessionCookies.apiKey, actionId as string);
	if (actionResponse.error()) {
		return {
			success: false,
			message: actionResponse.failureMessage()
		};
	}
};
