package service

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/KnoblauchPilze/backend-toolkit/pkg/db"
	"github.com/KnoblauchPilze/backend-toolkit/pkg/db/postgresql"
	"github.com/KnoblauchPilze/galactic-sovereign/pkg/persistence"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

var dbTestConfig = postgresql.NewConfigForLocalhost("db_galactic_sovereign", "galactic_sovereign_manager", "manager_password")

func newTestConnection(t *testing.T) db.Connection {
	conn, err := db.New(context.Background(), dbTestConfig)
	require.Nil(t, err)
	return conn
}

func insertTestUniverse(t *testing.T, conn db.Connection) persistence.Universe {
	someTime := time.Date(2024, 12, 8, 10, 10, 22, 0, time.UTC)

	universe := persistence.Universe{
		Id:        uuid.New(),
		Name:      fmt.Sprintf("my-universe-%s", uuid.NewString()),
		CreatedAt: someTime,
	}

	sqlQuery := `INSERT INTO universe (id, name, created_at) VALUES ($1, $2, $3) RETURNING updated_at`
	updatedAt, err := db.QueryOne[time.Time](
		context.Background(),
		conn,
		sqlQuery,
		universe.Id,
		universe.Name,
		universe.CreatedAt,
	)
	require.Nil(t, err)

	universe.UpdatedAt = updatedAt

	return universe
}

func insertTestResource(t *testing.T, conn db.Connection) persistence.Resource {
	someTime := time.Date(2024, 12, 8, 10, 26, 57, 0, time.UTC)

	resource := persistence.Resource{
		Id:              uuid.New(),
		Name:            fmt.Sprintf("my-resource-%s", uuid.NewString()),
		StartAmount:     456,
		StartProduction: 321,
		StartStorage:    778899,
		CreatedAt:       someTime,
	}

	sqlQuery := `INSERT INTO resource (id, name, start_amount, start_production, start_storage, created_at) VALUES ($1, $2, $3, $4, $5, $6) RETURNING updated_at`
	updatedAt, err := db.QueryOne[time.Time](
		context.Background(),
		conn,
		sqlQuery,
		resource.Id,
		resource.Name,
		resource.StartAmount,
		resource.StartProduction,
		resource.StartStorage,
		resource.CreatedAt,
	)
	require.Nil(t, err)

	resource.UpdatedAt = updatedAt

	return resource
}

func insertTestBuilding(t *testing.T, conn db.Connection) persistence.Building {
	someTime := time.Date(2024, 12, 8, 10, 12, 15, 0, time.UTC)

	building := persistence.Building{
		Id:        uuid.New(),
		Name:      fmt.Sprintf("my-building-%s", uuid.NewString()),
		CreatedAt: someTime,
	}

	sqlQuery := `INSERT INTO building (id, name, created_at) VALUES ($1, $2, $3) RETURNING updated_at`
	updatedAt, err := db.QueryOne[time.Time](
		context.Background(),
		conn,
		sqlQuery,
		building.Id,
		building.Name,
		building.CreatedAt,
	)
	require.Nil(t, err)

	building.UpdatedAt = updatedAt

	return building
}

func insertTestBuildingCost(t *testing.T, conn db.Connection, building uuid.UUID) (persistence.BuildingCost, persistence.Resource) {
	resource := insertTestResource(t, conn)

	cost := persistence.BuildingCost{
		Building: building,
		Resource: resource.Id,
		Cost:     41,
		Progress: 1.6,
	}

	sqlQuery := `INSERT INTO building_cost (building, resource, cost, progress) VALUES ($1, $2, $3, $4)`
	_, err := conn.Exec(
		context.Background(),
		sqlQuery,
		cost.Building,
		cost.Resource,
		cost.Cost,
		cost.Progress,
	)
	require.Nil(t, err)

	return cost, resource
}

func insertTestPlayer(t *testing.T, conn db.Connection, universe uuid.UUID) persistence.Player {
	someTime := time.Date(2024, 12, 8, 10, 9, 48, 0, time.UTC)

	player := persistence.Player{
		Id:        uuid.New(),
		ApiUser:   uuid.New(),
		Universe:  universe,
		Name:      fmt.Sprintf("my-player-%s", uuid.NewString()),
		CreatedAt: someTime,
	}

	sqlQuery := `INSERT INTO player (id, api_user, universe, name, created_at) VALUES ($1, $2, $3, $4, $5) RETURNING updated_at`
	updatedAt, err := db.QueryOne[time.Time](
		context.Background(),
		conn,
		sqlQuery,
		player.Id,
		player.ApiUser,
		player.Universe,
		player.Name,
		player.CreatedAt,
	)
	require.Nil(t, err)

	player.UpdatedAt = updatedAt

	return player
}

func insertTestPlayerInUniverse(t *testing.T, conn db.Connection) (persistence.Player, persistence.Universe) {
	universe := insertTestUniverse(t, conn)
	player := insertTestPlayer(t, conn, universe.Id)
	return player, universe
}

func insertTestPlanet(t *testing.T, conn db.Connection, player uuid.UUID) persistence.Planet {
	someTime := time.Date(2024, 12, 8, 10, 9, 58, 0, time.UTC)

	planet := persistence.Planet{
		Id:        uuid.New(),
		Player:    player,
		Name:      fmt.Sprintf("my-planet-%s", uuid.NewString()),
		Homeworld: true,
		CreatedAt: someTime,
		UpdatedAt: someTime,
	}

	sqlQuery := `INSERT INTO planet (id, player, name, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)`
	_, err := conn.Exec(
		context.Background(),
		sqlQuery,
		planet.Id,
		planet.Player,
		planet.Name,
		planet.CreatedAt,
		planet.UpdatedAt,
	)
	require.Nil(t, err)

	sqlQuery = `INSERT INTO homeworld (player, planet) VALUES ($1, $2)`
	_, err = conn.Exec(context.Background(), sqlQuery, planet.Player, planet.Id)
	require.Nil(t, err)

	return planet
}

func insertTestPlanetForPlayer(t *testing.T, conn db.Connection) (persistence.Planet, persistence.Player, persistence.Universe) {
	player, universe := insertTestPlayerInUniverse(t, conn)
	planet := insertTestPlanet(t, conn, player.Id)
	return planet, player, universe
}

func insertTestPlanetBuildingForPlanet(t *testing.T, conn db.Connection, planet uuid.UUID) (persistence.PlanetBuilding, persistence.Building) {
	someTime := time.Date(2024, 12, 8, 10, 22, 31, 0, time.UTC)

	building := insertTestBuilding(t, conn)

	planetBuilding := persistence.PlanetBuilding{
		Planet:    planet,
		Building:  building.Id,
		Level:     4,
		CreatedAt: someTime,
		UpdatedAt: someTime,
	}

	sqlQuery := `INSERT INTO planet_building (planet, building, level, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)`
	_, err := conn.Exec(
		context.Background(),
		sqlQuery,
		planetBuilding.Planet,
		planetBuilding.Building,
		planetBuilding.Level,
		planetBuilding.CreatedAt,
		planetBuilding.UpdatedAt,
	)
	require.Nil(t, err)

	return planetBuilding, building
}

func insertTestPlanetResourceForResource(t *testing.T, conn db.Connection, planet uuid.UUID, resource uuid.UUID) persistence.PlanetResource {
	someTime := time.Date(2024, 12, 8, 10, 28, 20, 0, time.UTC)

	planetResource := persistence.PlanetResource{
		Planet:    planet,
		Resource:  resource,
		Amount:    1011,
		CreatedAt: someTime,
		UpdatedAt: someTime,
	}

	sqlQuery := `INSERT INTO planet_resource (planet, resource, amount, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)`
	_, err := conn.Exec(
		context.Background(),
		sqlQuery,
		planetResource.Planet,
		planetResource.Resource,
		planetResource.Amount,
		planetResource.CreatedAt,
		planetResource.UpdatedAt,
	)
	require.Nil(t, err)

	return planetResource
}

func insertTestBuildingActionForPlanet(t *testing.T, conn db.Connection, planet uuid.UUID) (persistence.BuildingAction, persistence.Building) {
	someTime := time.Date(2024, 12, 10, 19, 38, 15, 0, time.UTC)

	building := insertTestBuilding(t, conn)

	action := persistence.BuildingAction{
		Id:           uuid.New(),
		Planet:       planet,
		Building:     building.Id,
		CurrentLevel: 4,
		DesiredLevel: 5,
		CreatedAt:    someTime,
		CompletedAt:  someTime.Add(1*time.Hour + 2*time.Minute),
	}

	sqlQuery := `INSERT INTO building_action (id, planet, building, current_level, desired_level, created_at, completed_at) VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := conn.Exec(
		context.Background(),
		sqlQuery,
		action.Id,
		action.Planet,
		action.Building,
		action.CurrentLevel,
		action.DesiredLevel,
		action.CreatedAt,
		action.CompletedAt,
	)
	require.Nil(t, err)

	return action, building
}
