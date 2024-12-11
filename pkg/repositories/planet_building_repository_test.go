package repositories

import (
	"context"
	"testing"
	"time"

	"github.com/KnoblauchPilze/backend-toolkit/pkg/db"
	"github.com/KnoblauchPilze/backend-toolkit/pkg/errors"
	eassert "github.com/KnoblauchPilze/easy-assert/assert"
	"github.com/KnoblauchPilze/galactic-sovereign/pkg/persistence"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIT_PlanetBuildingRepository_GetForPlanetAndBuilding(t *testing.T) {
	repo, conn, tx := newTestPlanetBuildingRepositoryAndTransaction(t)
	defer conn.Close(context.Background())
	planet, _, _ := insertTestPlanetForPlayer(t, conn)
	pb, building := insertTestPlanetBuildingForPlanet(t, conn, planet.Id)

	actual, err := repo.GetForPlanetAndBuilding(context.Background(), tx, planet.Id, building.Id)
	tx.Close(context.Background())
	assert.Nil(t, err)

	assert.True(t, eassert.EqualsIgnoringFields(actual, pb))
}

func TestIT_PlanetBuildingRepository_ListForPlanet(t *testing.T) {
	repo, conn, tx := newTestPlanetBuildingRepositoryAndTransaction(t)
	defer conn.Close(context.Background())
	planet1, _, _ := insertTestPlanetForPlayer(t, conn)
	pb1, _ := insertTestPlanetBuildingForPlanet(t, conn, planet1.Id)
	pb2, _ := insertTestPlanetBuildingForPlanet(t, conn, planet1.Id)
	planet2, _, _ := insertTestPlanetForPlayer(t, conn)
	_, b3 := insertTestPlanetBuildingForPlanet(t, conn, planet2.Id)

	actual, err := repo.ListForPlanet(context.Background(), tx, planet1.Id)
	tx.Close(context.Background())

	assert.Nil(t, err)
	assert.GreaterOrEqual(t, len(actual), 2)
	assert.True(t, eassert.ContainsIgnoringFields(actual, pb1))
	assert.True(t, eassert.ContainsIgnoringFields(actual, pb2))
	for _, planetBuilding := range actual {
		assert.NotEqual(t, planetBuilding.Planet, planet2.Id)
		assert.NotEqual(t, planetBuilding.Building, b3.Id)
	}
}

func TestIT_PlanetBuildingRepository_Update(t *testing.T) {
	repo, conn, tx := newTestPlanetBuildingRepositoryAndTransaction(t)
	defer conn.Close(context.Background())
	planet, _, _ := insertTestPlanetForPlayer(t, conn)
	planetBuilding, building := insertTestPlanetBuildingForPlanet(t, conn, planet.Id)

	updatedBuilding := planetBuilding
	updatedBuilding.UpdatedAt = planetBuilding.UpdatedAt.Add(24 * time.Minute)
	updatedBuilding.Level = planetBuilding.Level + 4

	actual, err := repo.Update(context.Background(), tx, updatedBuilding)
	tx.Close(context.Background())

	assert.Nil(t, err)

	expected := persistence.PlanetBuilding{
		Planet:    planet.Id,
		Building:  building.Id,
		Level:     planetBuilding.Level + 4,
		CreatedAt: planetBuilding.CreatedAt,
		UpdatedAt: updatedBuilding.UpdatedAt,
		Version:   planetBuilding.Version + 1,
	}

	assert.True(t, eassert.EqualsIgnoringFields(actual, expected))
}

func TestIT_PlanetBuildingRepository_Update_WhenVersionIsWrong_ExpectOptimisticLockException(t *testing.T) {
	repo, conn, tx := newTestPlanetBuildingRepositoryAndTransaction(t)
	defer conn.Close(context.Background())
	planet, _, _ := insertTestPlanetForPlayer(t, conn)
	planetBuilding, _ := insertTestPlanetBuildingForPlanet(t, conn, planet.Id)

	updatedBuilding := planetBuilding
	updatedBuilding.Level = planetBuilding.Level * 4
	updatedBuilding.Version = planetBuilding.Version + 2

	_, err := repo.Update(context.Background(), tx, updatedBuilding)
	tx.Close(context.Background())

	assert.True(t, errors.IsErrorWithCode(err, OptimisticLockException), "Actual err: %v", err)
}

func TestIT_PlanetBuildingRepository_Update_BumpsUpdatedAt(t *testing.T) {
	repo, conn := newTestPlanetBuildingRepository(t)
	defer conn.Close(context.Background())
	planet, _, _ := insertTestPlanetForPlayer(t, conn)
	planetBuilding, _ := insertTestPlanetBuildingForPlanet(t, conn, planet.Id)

	updatedBuilding := planetBuilding
	updatedBuilding.UpdatedAt = planetBuilding.UpdatedAt.Add(1 * time.Hour)
	updatedBuilding.Level = planetBuilding.Level + 2

	func() {
		tx, err := conn.BeginTx(context.Background())
		require.Nil(t, err)
		defer tx.Close(context.Background())

		_, err = repo.Update(context.Background(), tx, updatedBuilding)
		assert.Nil(t, err)
	}()

	var updatedBuildingFromDb persistence.PlanetBuilding
	func() {
		tx, err := conn.BeginTx(context.Background())
		require.Nil(t, err)
		defer tx.Close(context.Background())

		allBuildings, err := repo.ListForPlanet(context.Background(), tx, planet.Id)
		require.Nil(t, err)
		assert.Len(t, allBuildings, 1)

		updatedBuildingFromDb = allBuildings[0]
	}()

	assert.Equal(t, updatedBuildingFromDb.UpdatedAt, updatedBuildingFromDb.UpdatedAt)
}

func TestIT_PlanetBuildingRepository_Update_BumpsVersion(t *testing.T) {
	repo, conn := newTestPlanetBuildingRepository(t)
	defer conn.Close(context.Background())
	planet, _, _ := insertTestPlanetForPlayer(t, conn)
	planetBuilding, _ := insertTestPlanetBuildingForPlanet(t, conn, planet.Id)

	updatedBuilding := planetBuilding
	updatedBuilding.UpdatedAt = planetBuilding.UpdatedAt.Add(1 * time.Hour)
	updatedBuilding.Level = planetBuilding.Level + 2

	func() {
		tx, err := conn.BeginTx(context.Background())
		require.Nil(t, err)
		defer tx.Close(context.Background())

		_, err = repo.Update(context.Background(), tx, updatedBuilding)
		assert.Nil(t, err)
	}()

	var updatedBuildingFromDb persistence.PlanetBuilding
	func() {
		tx, err := conn.BeginTx(context.Background())
		require.Nil(t, err)
		defer tx.Close(context.Background())

		allBuildings, err := repo.ListForPlanet(context.Background(), tx, planet.Id)
		require.Nil(t, err)
		assert.Len(t, allBuildings, 1)

		updatedBuildingFromDb = allBuildings[0]
	}()

	assert.Equal(t, planetBuilding.Version+1, updatedBuildingFromDb.Version)
}

func newTestPlanetBuildingRepository(t *testing.T) (PlanetBuildingRepository, db.Connection) {
	conn := newTestConnection(t)
	return NewPlanetBuildingRepository(), conn
}

func newTestPlanetBuildingRepositoryAndTransaction(t *testing.T) (PlanetBuildingRepository, db.Connection, db.Transaction) {
	repo, conn := newTestPlanetBuildingRepository(t)
	tx, err := conn.BeginTx(context.Background())
	require.Nil(t, err)
	return repo, conn, tx
}

func insertTestPlanetBuildingForPlanet(t *testing.T, conn db.Connection, planet uuid.UUID) (persistence.PlanetBuilding, persistence.Building) {
	someTime := time.Date(2024, 12, 4, 21, 51, 15, 0, time.UTC)

	building := insertTestBuilding(t, conn)

	planetBuilding := persistence.PlanetBuilding{
		Planet:    planet,
		Building:  building.Id,
		Level:     0,
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

func assertPlanetBuildingExists(t *testing.T, conn db.Connection, planet uuid.UUID, building uuid.UUID) {
	sqlQuery := `SELECT COUNT(*) FROM planet_building WHERE planet = $1 AND building = $2`
	value, err := db.QueryOne[int](context.Background(), conn, sqlQuery, planet, building)
	require.Nil(t, err)
	require.Equal(t, 1, value)
}

func assertPlanetBuildingDoesNotExist(t *testing.T, conn db.Connection, planet uuid.UUID) {
	sqlQuery := `SELECT COUNT(building) FROM planet_building WHERE planet = $1`
	value, err := db.QueryOne[int](context.Background(), conn, sqlQuery, planet)
	require.Nil(t, err)
	require.Zero(t, value)
}

func assertPlanetBuildingLevel(t *testing.T, conn db.Connection, planet uuid.UUID, building uuid.UUID, level int) {
	sqlQuery := `SELECT level FROM planet_building WHERE planet = $1 AND building = $2`
	value, err := db.QueryOne[int](context.Background(), conn, sqlQuery, planet, building)
	require.Nil(t, err)
	require.Equal(t, level, value)
}
