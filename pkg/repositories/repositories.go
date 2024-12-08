package repositories

type Repositories struct {
	Building                         BuildingRepository
	BuildingAction                   BuildingActionRepository
	BuildingActionCost               BuildingActionCostRepository
	BuildingActionResourceProduction BuildingActionResourceProductionRepository
	BuildingCost                     BuildingCostRepository
	BuildingResourceProduction       BuildingResourceProductionRepository
	Planet                           PlanetRepository
	PlanetBuilding                   PlanetBuildingRepository
	PlanetResource                   PlanetResourceRepository
	PlanetResourceProduction         PlanetResourceProductionRepository
	PlanetResourceStorage            PlanetResourceStorageRepository
	Player                           PlayerRepository
	Resource                         ResourceRepository
	Universe                         UniverseRepository
}
