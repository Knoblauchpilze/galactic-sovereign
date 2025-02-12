package repositories

type Repositories struct {
	Building                         BuildingRepository
	BuildingAction                   BuildingActionRepository
	BuildingActionCost               BuildingActionCostRepository
	BuildingActionResourceProduction BuildingActionResourceProductionRepository
	BuildingActionResourceStorage    BuildingActionResourceStorageRepository
	BuildingCost                     BuildingCostRepository
	BuildingResourceProduction       BuildingResourceProductionRepository
	BuildingResourceStorage          BuildingResourceStorageRepository
	Planet                           PlanetRepository
	PlanetBuilding                   PlanetBuildingRepository
	PlanetResource                   PlanetResourceRepository
	PlanetResourceProduction         PlanetResourceProductionRepository
	PlanetResourceStorage            PlanetResourceStorageRepository
	Player                           PlayerRepository
	Resource                         ResourceRepository
	Universe                         UniverseRepository
}
