//go:generate go run go.uber.org/mock/mockgen -source=../ports/driven/for_checking_database_connection.go -destination=drivenportstest/database_mocks.go -package=drivenportstest
//go:generate go run go.uber.org/mock/mockgen -source=../ports/driven/for_fetching_time.go -destination=drivenportstest/time_mocks.go -package=drivenportstest
//go:generate go run go.uber.org/mock/mockgen -source=../ports/driven/for_listing_buildings.go -destination=drivenportstest/buildings_mocks.go -package=drivenportstest
//go:generate go run go.uber.org/mock/mockgen -source=../ports/driven/for_managing_planets.go -destination=drivenportstest/planets_mocks.go -package=drivenportstest
//go:generate go run go.uber.org/mock/mockgen -source=../ports/driven/for_managing_players.go -destination=drivenportstest/players_mocks.go -package=drivenportstest
//go:generate go run go.uber.org/mock/mockgen -source=../ports/driven/for_managing_universes.go -destination=drivenportstest/universes_mocks.go -package=drivenportstest
//go:generate go run go.uber.org/mock/mockgen -source=../ports/driven/for_mutating_planet.go -destination=drivenportstest/planet_mutator_mocks.go -package=drivenportstest

package usecases
