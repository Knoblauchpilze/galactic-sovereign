import { error, redirect } from '@sveltejs/kit';
import { loadAllCookies } from '$lib/cookies';

import { ApiFailureReason } from '$lib/responseEnvelope.js';

import { logout } from '$lib/actions/logout';
import { backToLobby } from '$lib/actions/backToLobby';
import { requestDeleteBuildingAction } from '$lib/actions/buildingAction';

import { Universe, type ApiUniverse, getUniverse } from '$lib/game/universes';
import { Planet, getPlanet } from '$lib/game/planets';

export async function load({ params, cookies, depends }) {
	const [valid, allCookies] = loadAllCookies(cookies);
	if (!valid) {
		redirect(303, '/login');
	}

	// https://learn.svelte.dev/tutorial/custom-dependencies
	depends('data:planet');

	const planetResponse = await getPlanet(allCookies.session.apiKey, params.planet);
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

	const universeResponse = await getUniverse(allCookies.game.universeId);
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
		playerName: allCookies.game.playerName,
		resources: universe.resources.map((r) => r.toJson()),
		buildings: universe.buildings.map((b) => b.toJson()),
		planet: planet.toJson()
	};
}

export const actions = {
	logout: logout,
	backToLobby: backToLobby,
	deleteBuildingAction: requestDeleteBuildingAction
};
