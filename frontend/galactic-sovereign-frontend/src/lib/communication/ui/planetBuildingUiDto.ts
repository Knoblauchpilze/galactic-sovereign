export interface PlanetBuildingCostUiDto {
	readonly resource: string;
	readonly cost: number;
}

export interface PlanetBuildingGainUiDto {
	readonly resource: string;
	readonly nextProduction: number;
	readonly gain: number;
}

export interface PlanetBuildingUiDto {
	readonly id: string;
	readonly name: string;
	readonly level: number;

	readonly planet: string;

	readonly costs: PlanetBuildingCostUiDto[];
	readonly resourcesProduction: PlanetBuildingGainUiDto[];
}
