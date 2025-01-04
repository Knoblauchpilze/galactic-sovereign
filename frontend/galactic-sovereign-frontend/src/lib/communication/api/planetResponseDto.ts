import { parseObjectAsArray } from '@totocorpsoftwareinc/frontend-toolkit';
import { PlanetResourceResponseDto } from './planetResourceResponseDto';
import { PlanetResourceProductionResponseDto } from './planetResourceProductionResponseDto';
import { PlanetResourceStorageResponseDto } from './planetResourceStorageResponseDto';
import { PlanetBuildingResponseDto } from './planetBuildingResponseDto';
import { BuildingActionResponseDto } from './buildingActionResponseDto';

export class PlanetResponseDto {
	readonly id: string = '00000000-0000-0000-0000-000000000000';
	readonly player: string = '00000000-0000-0000-0000-000000000000';
	readonly name: string = '';

	readonly createdAt: Date = new Date();

	readonly resources: PlanetResourceResponseDto[] = [];
	readonly productions: PlanetResourceProductionResponseDto[] = [];
	readonly storages: PlanetResourceStorageResponseDto[] = [];
	readonly buildings: PlanetBuildingResponseDto[] = [];
	readonly buildingActions: BuildingActionResponseDto[] = [];

	constructor(data: object) {
		if ('id' in data && typeof data.id === 'string') {
			this.id = data.id;
		}

		if ('player' in data && typeof data.player === 'string') {
			this.player = data.player;
		}

		if ('name' in data && typeof data.name === 'string') {
			this.name = data.name;
		}

		if ('createdAt' in data && typeof data.createdAt === 'string') {
			this.createdAt = new Date(data.createdAt);
		}

		if ('resources' in data && Array.isArray(data.resources)) {
			this.resources = parseObjectAsArray(data.resources, PlanetResourceResponseDto);
		}

		if ('productions' in data && Array.isArray(data.productions)) {
			this.productions = parseObjectAsArray(data.productions, PlanetResourceProductionResponseDto);
			console.log(
				'generated productions: ',
				JSON.stringify(this.productions),
				' from input: ',
				JSON.stringify(data.productions)
			);
		}

		if ('storages' in data && Array.isArray(data.storages)) {
			this.storages = parseObjectAsArray(data.storages, PlanetResourceStorageResponseDto);
		}

		if ('buildings' in data && Array.isArray(data.buildings)) {
			this.buildings = parseObjectAsArray(data.buildings, PlanetBuildingResponseDto);
		}

		if ('buildingActions' in data && Array.isArray(data.buildingActions)) {
			this.buildingActions = parseObjectAsArray(data.buildingActions, BuildingActionResponseDto);
		}
	}
}
