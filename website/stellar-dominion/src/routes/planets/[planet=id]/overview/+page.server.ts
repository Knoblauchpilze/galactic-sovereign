import { error, redirect } from '@sveltejs/kit';
import { loadCookies } from '$lib/cookies';

import { ApiFailureReason } from '$lib/responseEnvelope.js';

import { logout } from '$lib/actions/logout';
import { requestDeleteBuildingAction } from '$lib/actions/buildingAction';

import { Universe, type ApiUniverse, getUniverse } from '$lib/game/universes';
import { Planet, getPlanet } from '$lib/game/planets';

export async function load({ params, cookies, depends }) {
	const [valid, gameCookies] = loadCookies(cookies);
	if (!valid) {
		redirect(303, '/login');
	}

	// https://learn.svelte.dev/tutorial/custom-dependencies
	depends('data:planet');

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
		playerName: gameCookies.playerName,
		resources: universe.resources.map((r) => r.toJson()),
		buildings: universe.buildings.map((b) => b.toJson()),
		planet: planet.toJson()
	};
}

export const actions = {
	logout: logout,
	deleteBuildingAction: requestDeleteBuildingAction
};
