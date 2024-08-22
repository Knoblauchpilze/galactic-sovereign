import { error, redirect } from '@sveltejs/kit';
import { loadCookies } from '$lib/cookies';
import { Universe, type ApiUniverse, getUniverse } from '$lib/universes';
import { Planet, getPlanet, createBuildingAction } from '$lib/planets';
import { ApiFailureReason } from '$lib/responseEnvelope.js';
import { logoutUser } from '$lib/sessions';

/** @type {import('./$types').PageLoad} */
export async function load({ params, cookies }) {
	const [valid, gameCookies] = loadCookies(cookies);
	if (!valid) {
		redirect(303, '/login');
	}

	const planetResponse = await getPlanet(gameCookies.apiKey, params.planet);
	if (planetResponse.error()) {
		const reason = planetResponse.failureReason();

		switch (reason) {
			case ApiFailureReason.API_KEY_EXPIRED:
				redirect(303, '/login');
		}

		error(404, { message: planetResponse.failureMessage() });
	}

	// https://www.okupter.com/blog/sveltekit-cannot-stringify-arbitrary-non-pojos-error
	const planet = new Planet(planetResponse.getDetails());

	const universeResponse = await getUniverse(gameCookies.universeId);
	if (universeResponse.error()) {
		error(404, { message: universeResponse.failureMessage() });
	}

	const universe = new Universe(universeResponse.getDetails());
	const universeApi: ApiUniverse = {
		id: universe.id,
		name: universe.name
	};

	return {
		universe: universeApi,
		resources: universe.resources.map((r) => r.toJson()),
		buildings: universe.buildings.map((b) => b.toJson()),
		planet: planet.toJson()
	};
}

/** @type {import('./$types').Actions} */
export const actions = {
	logout: async ({ cookies }) => {
		const apiKey = cookies.get('api-key');
		if (!apiKey) {
			redirect(303, '/login');
		}

		const apiUser = cookies.get('api-user');
		if (!apiUser) {
			redirect(303, '/login');
		}

		const logoutResponse = await logoutUser(apiKey, apiUser);

		if (logoutResponse.error()) {
			return {
				success: false,
				message: logoutResponse.failureMessage()
			};
		}

		redirect(303, '/login');
	},

	createBuildingAction: async ({ cookies, params, request }) => {
		const apiKey = cookies.get('api-key');
		if (!apiKey) {
			redirect(303, '/login');
		}

		const apiUser = cookies.get('api-user');
		if (!apiUser) {
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

		const actionResponse = await createBuildingAction(apiKey, params.planet, buildingId as string);
		if (actionResponse.error()) {
			return {
				success: false,
				message: actionResponse.failureMessage()
			};
		}
	}
};
