import type { PlayerResponseDto } from '$lib/communication/api/playerResponseDto';
import type { UniverseResponseDto } from '$lib/communication/api/universeResponseDto';
import type { PlayerUiDto } from '$lib/communication/ui/playerUiDto';

export function playerResponseDtoToPlayerUiDto(
	apiDto: PlayerResponseDto,
	universeDto: UniverseResponseDto | undefined
): PlayerUiDto {
	let universeName = 'Unknown universe';
	if (universeDto !== undefined && apiDto.universe === universeDto.id) {
		universeName = universeDto.name;
	}

	return {
		id: apiDto.id,
		name: apiDto.name,
		universeId: apiDto.universe,
		universeName: universeName
	};
}
