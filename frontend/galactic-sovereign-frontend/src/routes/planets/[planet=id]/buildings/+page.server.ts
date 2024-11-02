import { redirect } from '@sveltejs/kit';
import { loadAllCookies } from '$lib/cookies';

import { analayzeResponseEnvelopAndRedirectIfNeeded } from '$lib/responseEnvelope.js';

import { logout } from '$lib/actions/logout';
import { backToLobby } from '$lib/actions/backToLobby';
import {
	requestCreateBuildingAction,
	requestDeleteBuildingAction
} from '$lib/actions/buildingAction';

import { Universe, getUniverse } from '$lib/game/universes';
import { Planet, getPlanet } from '$lib/game/planets';

export async function load({ params, cookies, depends }) {
	const [valid, allCookies] = loadAllCookies(cookies);
	if (!valid) {
		redirect(303, '/login');
	}

	// https://learn.svelte.dev/tutorial/custom-dependencies
	depends('data:planet');

	const planetResponse = await getPlanet(allCookies.session.apiKey, params.planet);
	analayzeResponseEnvelopAndRedirectIfNeeded(planetResponse);
	const planet = new Planet(planetResponse.getDetails());

	const universeResponse = await getUniverse(allCookies.game.universeId);
	analayzeResponseEnvelopAndRedirectIfNeeded(universeResponse);
	const universe = new Universe(universeResponse.getDetails());

	return {
		universe: universe.toJson(),
		playerName: allCookies.game.playerName,
		resources: universe.resources.map((r) => r.toJson()),
		buildings: universe.buildings.map((b) => b.toJson()),
		planet: planet.toJson()
	};
}

export const actions = {
	logout: logout,
	backToLobby: backToLobby,
	createBuildingAction: requestCreateBuildingAction,
	deleteBuildingAction: requestDeleteBuildingAction
};
