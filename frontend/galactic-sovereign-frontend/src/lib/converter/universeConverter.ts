import type { UniverseResponseDto } from '$lib/communication/api/universeResponseDto';
import type { UniverseUiDto } from '$lib/communication/ui/universeUiDto';

export function universeResponseDtoToUniverseUiDto(apiDto: UniverseResponseDto): UniverseUiDto {
	return {
		id: apiDto.id,
		name: apiDto.name
	};
}
