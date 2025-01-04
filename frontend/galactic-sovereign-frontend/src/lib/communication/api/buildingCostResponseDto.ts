export class BuildingCostResponseDto {
	readonly resource: string = '00000000-0000-0000-0000-000000000000';
	readonly cost: number = 0;
	readonly progress: number = 1;

	constructor(data: object) {
		if ('resource' in data && typeof data.resource === 'string') {
			this.resource = data.resource;
		}

		if ('cost' in data && typeof data.cost === 'number') {
			this.cost = data.cost;
		}

		if ('progress' in data && typeof data.progress === 'number') {
			this.progress = data.progress;
		}
	}
}
