package repositories

type Repositories struct {
	Acl            AclRepository
	ApiKey         ApiKeyRepository
	Building       BuildingRepository
	Planet         PlanetRepository
	PlanetBuilding PlanetBuildingRepository
	PlanetResource PlanetResourceRepository
	Player         PlayerRepository
	Resource       ResourceRepository
	UserLimit      UserLimitRepository
	User           UserRepository
	Universe       UniverseRepository
}
