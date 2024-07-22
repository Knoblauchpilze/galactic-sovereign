import { error, redirect } from '@sveltejs/kit';
import Planet, { getPlanet } from '$lib/planets';
import { ApiFailureReason } from '$lib/responseEnvelope.js';

// https://stackoverflow.com/questions/7905929/how-to-test-valid-uuid-guid
const UUID_REGEX = /^[0-9a-f]{8}-[0-9a-f]{4}-[0-5][0-9a-f]{3}-[089ab][0-9a-f]{3}-[0-9a-f]{12}$/i;

/** @type {import('./$types').PageLoad} */
export async function load({ params, cookies }) {
	if (!(typeof params.planet === 'string') || !UUID_REGEX.test(params.planet)) {
		error(404, 'Planet not found');
	}

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
