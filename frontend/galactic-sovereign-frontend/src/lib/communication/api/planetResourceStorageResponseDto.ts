export class PlanetResourceStorageResponseDto {
	readonly resource: string = '00000000-0000-0000-0000-000000000000';
	readonly storage: number = 0;

	constructor(data: object) {
		if ('resource' in data && typeof data.resource === 'string') {
			this.resource = data.resource;
		}

		if ('storage' in data && typeof data.storage === 'number') {
			this.storage = data.storage;
		}
	}
}
