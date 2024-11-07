import { loadAllCookiesOrRedirectToLogin } from '$lib/cookies';
import { analayzeResponseEnvelopAndRedirectIfNeeded } from '$lib/responseEnvelope.js';

import { HOMEPAGE_TITLE } from '$lib/stores/ui/pageTitle';

import { logout } from '$lib/actions/logout';
import { backToLobby } from '$lib/actions/backToLobby';
import { requestDeleteBuildingAction } from '$lib/actions/buildingAction';

import { Universe, getUniverse } from '$lib/game/universes';
import { Planet, getPlanet } from '$lib/game/planets';
import { mapPlanetResourcesToUiResources } from '$lib/game/resources';
import { mapPlanetBuildingsToUiBuildings } from '$lib/game/buildings';
import { mapBuildingActionsToUiActions } from '$lib/game/actions.js';

export async function load({ params, cookies, depends }) {
	const allCookies = loadAllCookiesOrRedirectToLogin(cookies);

	depends('data:planet');

	const planetResponse = await getPlanet(allCookies.session.apiKey, params.planet);
	analayzeResponseEnvelopAndRedirectIfNeeded(planetResponse);
	// https://www.okupter.com/blog/sveltekit-cannot-stringify-arbitrary-non-pojos-error
	const planet = new Planet(planetResponse.getDetails());

	const universeResponse = await getUniverse(allCookies.game.universeId);
	analayzeResponseEnvelopAndRedirectIfNeeded(universeResponse);
	const universe = new Universe(universeResponse.getDetails());

	const out = {
		wepageTitle: HOMEPAGE_TITLE + ' - ' + planet.name,

		universeName: universe.name,
		playerName: allCookies.game.playerName,
		planetName: planet.name,
		planetCreationTime: planet.createdAt,

		resources: mapPlanetResourcesToUiResources(planet, universe.resources),
		buildings: mapPlanetBuildingsToUiBuildings(planet, universe.buildings, universe.resources),
		buildingActions: mapBuildingActionsToUiActions(planet.buildingActions, universe.buildings)
	};

	return out;
}

export const actions = {
	logout: logout,
	backToLobby: backToLobby,
	deleteBuildingAction: requestDeleteBuildingAction
};
