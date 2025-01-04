import { parseObjectAsArray } from "@totocorpsoftwareinc/frontend-toolkit";
import { BuildingResponseDto } from "./buildingResponseDto";
import { ResourceResponseDto } from "./resourceResponseDto";

export class UniverseResponseDto {
  readonly id: string = '00000000-0000-0000-0000-000000000000';
  readonly name: string = '';

  readonly resources: ResourceResponseDto[] = [];
  readonly buildings: BuildingResponseDto[] = [];

  constructor(data: object) {
    if ('id' in data && typeof data.id === 'string') {
      this.id = data.id;
    }

    if ('name' in data && typeof data.name === 'string') {
      this.name = data.name;
    }

    if ('resources' in data && Array.isArray(data.resources)) {
      this.resources = parseObjectAsArray(data.resources, ResourceResponseDto);
    }

    if ('buildings' in data && Array.isArray(data.buildings)) {
      this.buildings = parseObjectAsArray(data.buildings, BuildingResponseDto);
    }
  }
}