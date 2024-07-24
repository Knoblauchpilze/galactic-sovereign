import { error, redirect } from '@sveltejs/kit';
import Planet, { getPlanet } from '$lib/planets';
import { ApiFailureReason } from '$lib/responseEnvelope.js';

/** @type {import('./$types').PageLoad} */
export async function load({ params, cookies }) {
	const apiKey = cookies.get('api-key');
	if (!apiKey) {
		redirect(303, '/login');
	}

	const apiUser = cookies.get('api-user');
	if (!apiUser) {
		redirect(303, '/login');
	}

	const planetResponse = await getPlanet(apiKey, params.planet);	
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

	return {
		planet: {
			...planet
		},
	}
}
