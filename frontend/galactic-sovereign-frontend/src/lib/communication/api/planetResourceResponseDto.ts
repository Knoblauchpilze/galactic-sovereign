export class PlanetResourceResponseDto {
	readonly resource: string = '00000000-0000-0000-0000-000000000000';
	readonly amount: number = 0;

	constructor(data: object) {
		if ('resource' in data && typeof data.resource === 'string') {
			this.resource = data.resource;
		}

		if ('amount' in data && typeof data.amount === 'number') {
			this.amount = data.amount;
		}
	}
}
