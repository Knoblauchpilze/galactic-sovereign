//go:generate go run go.uber.org/mock/mockgen -source=../../app/ports/driving/for_checking_service_health.go -destination=drivingportstest/health_mocks.go -package=drivingportstest
//go:generate go run go.uber.org/mock/mockgen -source=../../app/ports/driving/for_creating_building_action.go -destination=drivingportstest/create_building_action_mocks.go -package=drivingportstest
//go:generate go run go.uber.org/mock/mockgen -source=../../app/ports/driving/for_creating_planet.go -destination=drivingportstest/create_planet_mocks.go -package=drivingportstest
//go:generate go run go.uber.org/mock/mockgen -source=../../app/ports/driving/for_deleting_building_action.go -destination=drivingportstest/deleting_building_action_mocks.go -package=drivingportstest
//go:generate go run go.uber.org/mock/mockgen -source=../../app/ports/driving/for_managing_building_action.go -destination=drivingportstest/building_action_mocks.go -package=drivingportstest
//go:generate go run go.uber.org/mock/mockgen -source=../../app/ports/driving/for_managing_planet.go -destination=drivingportstest/planet_mocks.go -package=drivingportstest
//go:generate go run go.uber.org/mock/mockgen -source=../../app/ports/driving/for_managing_player.go -destination=drivingportstest/player_mocks.go -package=drivingportstest
//go:generate go run go.uber.org/mock/mockgen -source=../../app/ports/driving/for_managing_universe.go -destination=drivingportstest/universe_mocks.go -package=drivingportstest

package drivingadapters
