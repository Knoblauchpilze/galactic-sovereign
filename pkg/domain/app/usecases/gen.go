//go:generate go run go.uber.org/mock/mockgen -source=../ports/driven/for_managing_planets.go -destination=drivenportstest/planets_mocks.go -package=drivenportstest
//go:generate go run go.uber.org/mock/mockgen -source=../ports/driven/for_managing_players.go -destination=drivenportstest/players_mocks.go -package=drivenportstest
//go:generate go run go.uber.org/mock/mockgen -source=../ports/driven/for_managing_universes.go -destination=drivenportstest/universes_mocks.go -package=drivenportstest

package usecases
