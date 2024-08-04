import { error, redirect } from '@sveltejs/kit';
import { loadCookies } from '$lib/cookies';
import { Planet, getPlanet } from '$lib/planets';
import { getResources, responseToResourcesArray } from '$lib/resources';
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

	const resourcesResponse = await getResources(gameCookies.apiKey);
	if (resourcesResponse.error()) {
		const reason = resourcesResponse.failureReason();

		switch (reason) {
			case ApiFailureReason.API_KEY_EXPIRED:
				redirect(303, '/login');
		}

		error(404, { message: resourcesResponse.failureMessage() });
	}

	const resources = responseToResourcesArray(resourcesResponse);

	return {
		resources: resources.map((r) => r.toJson()),
		planet: {
			...planet
		}
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
	}
};
