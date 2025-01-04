import type { BuildingActionResponseDto } from '$lib/communication/api/buildingActionResponseDto';
import type { BuildingResponseDto } from '$lib/communication/api/buildingResponseDto';
import type { BuildingActionUiDto } from '$lib/communication/ui/buildingActionUiDto';

export function buildingActionResponseDtoToBuildingActionUiDto(
	apiDto: BuildingActionResponseDto,
	buildingDto: BuildingResponseDto | undefined
): BuildingActionUiDto {
	let buildingName = 'Unknown building';
	if (buildingDto !== undefined && apiDto.building === buildingDto.id) {
		buildingName = buildingDto.name;
	}

	return {
		id: apiDto.id,
		name: buildingName,
		planet: apiDto.planet,
		nextLevel: apiDto.desiredLevel,
		completedAt: apiDto.completedAt
	};
}
