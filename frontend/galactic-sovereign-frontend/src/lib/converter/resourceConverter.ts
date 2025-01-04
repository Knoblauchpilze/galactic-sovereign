import type { PlanetResponseDto } from '$lib/communication/api/planetResponseDto';
import type { ResourceResponseDto } from '$lib/communication/api/resourceResponseDto';
import type { PlanetResourceUiDto } from '$lib/communication/ui/planetResourceUiDto';
import type { ResourceUiDto } from '$lib/communication/ui/resourceUiDto';

export function resourceResponseDtoToResourceUiDto(apiDto: ResourceResponseDto): ResourceUiDto {
	return {
		id: apiDto.id,
		name: apiDto.name
	};
}

export function planetResourceResponseDtoToPlanetResourceUiDto(
	resourceDto: ResourceResponseDto,
	planetDto: PlanetResponseDto): PlanetResourceUiDto {
		const maybePlanetResource = planetDto.resources.find((r) => r.resource === resourceDto.id);
		const amount = maybePlanetResource === undefined ? 0 : maybePlanetResource.amount;

		console.log("converting resource :", JSON.stringify(resourceDto), " with planet: ", JSON.stringify(planetDto));

		const production = planetDto.productions.reduce((currentProduction, resource) => {
			if (resource.resource === resourceDto.id) {
				return currentProduction + resource.production;
			}
			return currentProduction;
		}, 0);

		const maybePlanetStorage = planetDto.storages.find((s) => s.resource === resourceDto.id);
		const storage = maybePlanetStorage === undefined ? 0 : maybePlanetStorage.storage;

	return {
		name: resourceDto.name,
		amount: amount,
		production: production,
		storage: storage,
	};
}
