export class ResourceResponseDto {
	readonly id: string = '00000000-0000-0000-0000-000000000000';
	readonly name: string = '';

	constructor(data: object) {
		if ('id' in data && typeof data.id === 'string') {
			this.id = data.id;
		}

		if ('name' in data && typeof data.name === 'string') {
			this.name = data.name;
		}
	}
}
