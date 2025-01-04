import { parseObjectAsArray } from '@totocorpsoftwareinc/frontend-toolkit';
import { BuildingCostResponseDto } from './buildingCostResponseDto';
import { BuildingResourceProductionResponseDto } from './buildingResourceProductionResponseDto';

export class BuildingResponseDto {
	readonly id: string = '00000000-0000-0000-0000-000000000000';
	readonly name: string = '';
	readonly costs: BuildingCostResponseDto[] = [];
	readonly resourceProductions: BuildingResourceProductionResponseDto[] = [];

	constructor(data: object) {
		if ('id' in data && typeof data.id === 'string') {
			this.id = data.id;
		}

		if ('name' in data && typeof data.name === 'string') {
			this.name = data.name;
		}

		if ('costs' in data && Array.isArray(data.costs)) {
			this.costs = parseObjectAsArray(data.costs, BuildingCostResponseDto);
		}

		if ('productions' in data && Array.isArray(data.productions)) {
			this.resourceProductions = parseObjectAsArray(
				data.productions,
				BuildingResourceProductionResponseDto
			);
		}
	}
}
