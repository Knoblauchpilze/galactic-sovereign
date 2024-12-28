export class ApiKeyResponseDto {
	readonly user: string = '00000000-0000-0000-0000-000000000000';
	readonly key: string = '00000000-0000-0000-0000-000000000000';
	readonly validUntil: Date = new Date();

	constructor(data: object) {
		if ('user' in data && typeof data.user === 'string') {
			this.user = data.user;
		}

		if ('key' in data && typeof data.key === 'string') {
			this.key = data.key;
		}

		if ('validUntil' in data && typeof data.validUntil === 'string') {
			this.validUntil = new Date(data.validUntil);
		}
	}
}
