package repositories

import (
	"context"
	"testing"
	"time"

	"github.com/Knoblauchpilze/backend-toolkit/pkg/db"
	"github.com/Knoblauchpilze/backend-toolkit/pkg/errors"
	eassert "github.com/Knoblauchpilze/easy-assert/assert"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/persistence"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIT_PlanetBuildingRepository_GetForPlanetAndBuilding(t *testing.T) {
	repo, conn, tx := newTestPlanetBuildingRepositoryAndTransaction(t)
	defer conn.Close(context.Background())
	planetId, _ := insertTestPlanetForPlayer(t, conn)
	pb, building := insertTestPlanetBuildingForPlanet(t, conn, planetId)

	actual, err := repo.GetForPlanetAndBuilding(context.Background(), tx, planetId, building.Id)
	tx.Close(context.Background())
	assert.Nil(t, err)

	assert.True(t, eassert.EqualsIgnoringFields(actual, pb))
}

func TestIT_PlanetBuildingRepository_ListForPlanet(t *testing.T) {
	repo, conn, tx := newTestPlanetBuildingRepositoryAndTransaction(t)
	defer conn.Close(context.Background())
	planetId1, _ := insertTestPlanetForPlayer(t, conn)
	pb1, _ := insertTestPlanetBuildingForPlanet(t, conn, planetId1)
	pb2, _ := insertTestPlanetBuildingForPlanet(t, conn, planetId1)
	planetId2, _ := insertTestPlanetForPlayer(t, conn)
	_, b3 := insertTestPlanetBuildingForPlanet(t, conn, planetId2)

	actual, err := repo.ListForPlanet(context.Background(), tx, planetId1)
	tx.Close(context.Background())

	assert.Nil(t, err)
	assert.GreaterOrEqual(t, len(actual), 2)
	assert.True(t, eassert.ContainsIgnoringFields(actual, pb1))
	assert.True(t, eassert.ContainsIgnoringFields(actual, pb2))
	for _, planetBuilding := range actual {
		assert.NotEqual(t, planetBuilding.Planet, planetId2)
		assert.NotEqual(t, planetBuilding.Building, b3.Id)
	}
}

func TestIT_PlanetBuildingRepository_Update(t *testing.T) {
	repo, conn, tx := newTestPlanetBuildingRepositoryAndTransaction(t)
	defer conn.Close(context.Background())
	planetId, _ := insertTestPlanetForPlayer(t, conn)
	planetBuilding, building := insertTestPlanetBuildingForPlanet(t, conn, planetId)

	updatedBuilding := planetBuilding
	updatedBuilding.UpdatedAt = planetBuilding.UpdatedAt.Add(24 * time.Minute)
	updatedBuilding.Level = planetBuilding.Level + 4

	actual, err := repo.Update(context.Background(), tx, updatedBuilding)
	tx.Close(context.Background())

	assert.Nil(t, err)

	expected := persistence.PlanetBuilding{
		Planet:    planetId,
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
	planetId, _ := insertTestPlanetForPlayer(t, conn)
	planetBuilding, _ := insertTestPlanetBuildingForPlanet(t, conn, planetId)

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
	planetId, _ := insertTestPlanetForPlayer(t, conn)
	planetBuilding, _ := insertTestPlanetBuildingForPlanet(t, conn, planetId)

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

		allBuildings, err := repo.ListForPlanet(context.Background(), tx, planetId)
		require.Nil(t, err)
		assert.Len(t, allBuildings, 1)

		updatedBuildingFromDb = allBuildings[0]
	}()

	assert.Equal(t, updatedBuildingFromDb.UpdatedAt, updatedBuildingFromDb.UpdatedAt)
}

func TestIT_PlanetBuildingRepository_Update_BumpsVersion(t *testing.T) {
	repo, conn := newTestPlanetBuildingRepository(t)
	defer conn.Close(context.Background())
	planetId, _ := insertTestPlanetForPlayer(t, conn)
	planetBuilding, _ := insertTestPlanetBuildingForPlanet(t, conn, planetId)

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

		allBuildings, err := repo.ListForPlanet(context.Background(), tx, planetId)
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
