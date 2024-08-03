package repositories

type Repositories struct {
	Acl            AclRepository
	ApiKey         ApiKeyRepository
	Planet         PlanetRepository
	PlanetResource PlanetResourceRepository
	Player         PlayerRepository
	Resource       ResourceRepository
	UserLimit      UserLimitRepository
	User           UserRepository
	Universe       UniverseRepository
}
