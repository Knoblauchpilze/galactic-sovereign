export interface ApiBuildingAction {
	readonly id: string;
	readonly building: string;
	readonly desiredLevel: number;
	readonly completedAt: Date;
}

export class BuildingAction {
	readonly id: string = '00000000-0000-0000-0000-000000000000';
	readonly building: string = '00000000-0000-0000-0000-000000000000';
	readonly desiredLevel: number = 0;
	readonly completedAt: Date = new Date();

	constructor(response: object) {
		if ('id' in response && typeof response.id === 'string') {
			this.id = response.id;
		}

		if ('building' in response && typeof response.building === 'string') {
			this.building = response.building;
		}

		if ('desiredLevel' in response && typeof response.desiredLevel === 'number') {
			this.desiredLevel = response.desiredLevel;
		}

		if ('completedAt' in response && typeof response.completedAt === 'string') {
			this.completedAt = new Date(response.completedAt);
		}
	}

	public toJson(): ApiBuildingAction {
		return {
			id: this.id,
			building: this.building,
			desiredLevel: this.desiredLevel,
			completedAt: this.completedAt
		};
	}
}

export function parseBuildingActions(data: object[]): BuildingAction[] {
	const out: BuildingAction[] = [];

	for (const maybeAction of data) {
		const hasId = 'id' in maybeAction && typeof maybeAction.id === 'string';
		const hasBuilding = 'building' in maybeAction && typeof maybeAction.building === 'string';
		const hasDesiredLevel =
			'desiredLevel' in maybeAction && typeof maybeAction.desiredLevel === 'number';
		const hasCompletedAt =
			'completedAt' in maybeAction && typeof maybeAction.completedAt === 'string';

		if (hasId && hasBuilding && hasDesiredLevel && hasCompletedAt) {
			out.push(new BuildingAction(maybeAction));
		}
	}

	return out;
}
