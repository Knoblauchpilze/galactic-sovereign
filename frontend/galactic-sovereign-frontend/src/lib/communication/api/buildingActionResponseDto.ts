export class BuildingActionResponseDto {
	readonly id: string = '00000000-0000-0000-0000-000000000000';
	readonly building: string = '00000000-0000-0000-0000-000000000000';
	readonly planet: string = '00000000-0000-0000-0000-000000000000';
	readonly desiredLevel: number = 0;
	readonly completedAt: Date = new Date();

	constructor(data: object) {
		if ('id' in data && typeof data.id === 'string') {
			this.id = data.id;
		}

		if ('planet' in data && typeof data.planet === 'string') {
			this.planet = data.planet;
		}

		if ('building' in data && typeof data.building === 'string') {
			this.building = data.building;
		}

		if ('desiredLevel' in data && typeof data.desiredLevel === 'number') {
			this.desiredLevel = data.desiredLevel;
		}

		if ('completedAt' in data && typeof data.completedAt === 'string') {
			this.completedAt = new Date(data.completedAt);
		}
	}
}
