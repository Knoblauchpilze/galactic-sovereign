export class PlanetBuildingResponseDto {
	readonly building: string = '00000000-0000-0000-0000-000000000000';
	readonly level: number = 0;

	constructor(data: object) {
		if ('building' in data && typeof data.building === 'string') {
			this.building = data.building;
		}

		if ('level' in data && typeof data.level === 'number') {
			this.level = data.level;
		}
	}
}
