package repositories

type Repositories struct {
	Acl       AclRepository
	ApiKey    ApiKeyRepository
	UserLimit UserLimitRepository
	User      UserRepository
	Universe  UniverseRepository
	Player    PlayerRepository
	Planet    PlanetRepository
}
