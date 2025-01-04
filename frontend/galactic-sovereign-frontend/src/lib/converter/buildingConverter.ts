import type { BuildingCostResponseDto } from '$lib/communication/api/buildingCostResponseDto';
import type { BuildingResourceProductionResponseDto } from '$lib/communication/api/buildingResourceProductionResponseDto';
import type { BuildingResponseDto } from '$lib/communication/api/buildingResponseDto';
import type { PlanetResponseDto } from '$lib/communication/api/planetResponseDto';
import type { ResourceResponseDto } from '$lib/communication/api/resourceResponseDto';
import type {
	PlanetBuildingCostUiDto,
	PlanetBuildingGainUiDto,
	PlanetBuildingUiDto
} from '$lib/communication/ui/planetBuildingUiDto';
import { computeCostForLevel } from '$lib/game/buildingCost';
import { computeProductionForLevel } from '$lib/game/buildingProduction';

function buildingCostResponseDtoToPlanetBuildingCostUiDto(
	buildingCostDto: BuildingCostResponseDto,
	level: number,
	resourcesDto: ResourceResponseDto[]
): PlanetBuildingCostUiDto {
	const maybeResource = resourcesDto.find((r) => r.id === buildingCostDto.resource);
	const resource = maybeResource === undefined ? 'Unknown resource' : maybeResource.name;

	return {
		resource: resource,
		cost: computeCostForLevel(buildingCostDto.cost, buildingCostDto.progress, level + 1)
	};
}

function buildingProductionResponseDtoToPlanetBuildingGainUiDto(
	buildingResourceProductionDto: BuildingResourceProductionResponseDto,
	level: number,
	resourcesDto: ResourceResponseDto[]
): PlanetBuildingGainUiDto {
	const maybeResource = resourcesDto.find((r) => r.id === buildingResourceProductionDto.resource);
	const resource = maybeResource === undefined ? 'Unknown resource' : maybeResource.name;

	const currentProduction = computeProductionForLevel(
		buildingResourceProductionDto.base,
		buildingResourceProductionDto.progress,
		level
	);
	const nextProduction = computeProductionForLevel(
		buildingResourceProductionDto.base,
		buildingResourceProductionDto.progress,
		level + 1
	);

	return {
		resource: resource,
		nextProduction: nextProduction,
		gain: nextProduction - currentProduction
	};
}

export function buildingResponseDtoToPlanetResourceUiDto(
	buildingDto: BuildingResponseDto,
	resourcesDto: ResourceResponseDto[],
	planetDto: PlanetResponseDto
): PlanetBuildingUiDto {
	const maybePlanetBuilding = planetDto.buildings.find((b) => b.building === buildingDto.id);
	const level = maybePlanetBuilding === undefined ? 0 : maybePlanetBuilding.level;

	const costs = buildingDto.costs.map((c) =>
		buildingCostResponseDtoToPlanetBuildingCostUiDto(c, level, resourcesDto)
	);

	const productions = buildingDto.resourceProductions.map((p) =>
		buildingProductionResponseDtoToPlanetBuildingGainUiDto(p, level, resourcesDto)
	);

	return {
		id: buildingDto.id,
		name: buildingDto.name,
		level: level,

		planet: planetDto.id,

		costs: costs,
		resourcesProduction: productions
	};
}
