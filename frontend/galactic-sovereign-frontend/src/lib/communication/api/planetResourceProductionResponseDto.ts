export class PlanetResourceProductionResponseDto {
	readonly resource: string = '00000000-0000-0000-0000-000000000000';
	readonly building: string | undefined = undefined;
	readonly production: number = 0;

	constructor(data: object) {
		if ('resource' in data && typeof data.resource === 'string') {
			this.resource = data.resource;
		}

		if ('building' in data && typeof data.building === 'string') {
			this.building = data.building;
		}

		if ('production' in data && typeof data.production === 'number') {
			this.production = data.production;
		}
	}
}
