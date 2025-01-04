export class BuildingResourceProductionResponseDto {
	readonly resource: string = '00000000-0000-0000-0000-000000000000';
	readonly base: number = 0;
	readonly progress: number = 1;

	constructor(data: object) {
		if ('resource' in data && typeof data.resource === 'string') {
			this.resource = data.resource;
		}

		if ('base' in data && typeof data.base === 'number') {
			this.base = data.base;
		}

		if ('progress' in data && typeof data.progress === 'number') {
			this.progress = data.progress;
		}
	}
}
