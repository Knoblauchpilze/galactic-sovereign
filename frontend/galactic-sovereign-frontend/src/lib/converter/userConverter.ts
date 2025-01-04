import type { UserResponseDto } from '$lib/communication/api/userResponseDto';
import type { UserUiDto } from '$lib/communication/ui/userUiDto';
import { formatDate } from '$lib/time';

export function userResponseDtoToUserUiDto(apiDto: UserResponseDto): UserUiDto {
	return {
		id: apiDto.id,
		email: apiDto.email,
		password: apiDto.password,
		createdAt: formatDate(apiDto.createdAt)
	};
}
