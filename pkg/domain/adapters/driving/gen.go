//go:generate go run go.uber.org/mock/mockgen -source=../../app/ports/driving/for_managing_planet.go -destination=drivingportstest/planet_mocks.go -package=drivingportstest
//go:generate go run go.uber.org/mock/mockgen -source=../../app/ports/driving/for_managing_player.go -destination=drivingportstest/player_mocks.go -package=drivingportstest
//go:generate go run go.uber.org/mock/mockgen -source=../../app/ports/driving/for_managing_universe.go -destination=drivingportstest/universe_mocks.go -package=drivingportstest

package drivingadapters
