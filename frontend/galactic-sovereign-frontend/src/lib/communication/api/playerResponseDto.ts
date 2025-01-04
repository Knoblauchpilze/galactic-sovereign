export class PlayerResponseDto {
	readonly id: string = '00000000-0000-0000-0000-000000000000';
	readonly apiUser: string = '00000000-0000-0000-0000-000000000000';
	readonly universe: string = '00000000-0000-0000-0000-000000000000';
	readonly name: string = '';
	readonly createdAt: Date = new Date();

	constructor(data: object) {
		if ('id' in data && typeof data.id === 'string') {
			this.id = data.id;
		}

		if ('apiUser' in data && typeof data.apiUser === 'string') {
			this.apiUser = data.apiUser;
		}

		if ('universe' in data && typeof data.universe === 'string') {
			this.universe = data.universe;
		}

		if ('name' in data && typeof data.name === 'string') {
			this.name = data.name;
		}

		if ('createdAt' in data && typeof data.createdAt === 'string') {
			this.createdAt = new Date(data.createdAt);
		}
	}
}
