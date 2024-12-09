package repositories

import (
	"context"
	"testing"
	"time"

	"github.com/KnoblauchPilze/backend-toolkit/pkg/db"
	"github.com/KnoblauchPilze/backend-toolkit/pkg/db/pgx"
	"github.com/KnoblauchPilze/backend-toolkit/pkg/errors"
	eassert "github.com/KnoblauchPilze/easy-assert/assert"
	"github.com/KnoblauchPilze/galactic-sovereign/pkg/persistence"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIT_PlanetResourceProductionRepository_Create(t *testing.T) {
	repo, conn, tx := newTestPlanetResourceProductionRepositoryAndTransaction(t)
	defer conn.Close(context.Background())
	planet, _, _ := insertTestPlanetForPlayer(t, conn)
	resource := insertTestResource(t, conn)
	building := insertTestBuilding(t, conn)

	prod := persistence.PlanetResourceProduction{
		Planet:     planet.Id,
		Building:   &building.Id,
		Resource:   resource.Id,
		Production: 56,
		CreatedAt:  time.Now(),
	}

	actual, err := repo.Create(context.Background(), tx, prod)
	assert.Nil(t, err)
	tx.Close(context.Background())

	assert.True(t, eassert.EqualsIgnoringFields(actual, prod, "UpdatedAt"))
	assert.Equal(t, actual.UpdatedAt, actual.CreatedAt)
	assertPlanetResourceProductionExists(t, conn, planet.Id, resource.Id)
	assertPlanetResourceProductionForBuilding(t, conn, planet.Id, resource.Id, building.Id, 56)
}

func TestIT_PlanetResourceProductionRepository_Create_WithoutBuilding(t *testing.T) {
	repo, conn, tx := newTestPlanetResourceProductionRepositoryAndTransaction(t)
	defer conn.Close(context.Background())
	planet, _, _ := insertTestPlanetForPlayer(t, conn)
	resource := insertTestResource(t, conn)

	prod := persistence.PlanetResourceProduction{
		Planet:     planet.Id,
		Resource:   resource.Id,
		Production: 56,
		CreatedAt:  time.Now(),
	}

	actual, err := repo.Create(context.Background(), tx, prod)
	assert.Nil(t, err)
	tx.Close(context.Background())

	assert.True(t, eassert.EqualsIgnoringFields(actual, prod, "UpdatedAt"))
	assert.Equal(t, actual.UpdatedAt, actual.CreatedAt)
	assertPlanetResourceProductionExists(t, conn, planet.Id, resource.Id)
	assertPlanetResourceProductionWithoutBuilding(t, conn, planet.Id, resource.Id, 56)
}

func TestIT_PlanetResourceProductionRepository_Create_WhenDuplicateResourceProductionWithoutBuilding_ExpectSuccess(t *testing.T) {
	repo, conn := newTestPlanetResourceProductionRepository(t)
	defer conn.Close(context.Background())
	planet, _, _ := insertTestPlanetForPlayer(t, conn)
	production, resource := insertTestPlanetResourceProductionForBuilding(t, conn, planet.Id, nil)

	newProd := persistence.PlanetResourceProduction{
		Planet:     planet.Id,
		Resource:   resource.Id,
		Production: production.Production * 2,
		CreatedAt:  time.Date(2024, 12, 1, 17, 52, 21, 0, time.UTC),
	}

	func() {
		tx, err := conn.BeginTx(context.Background())
		require.Nil(t, err)
		defer tx.Close(context.Background())

		_, err = repo.Create(context.Background(), tx, newProd)
		assert.Nil(t, err)
	}()

	var allProductions []persistence.PlanetResourceProduction
	func() {
		tx, err := conn.BeginTx(context.Background())
		require.Nil(t, err)
		defer tx.Close(context.Background())

		allProductions, err = repo.ListForPlanet(context.Background(), tx, planet.Id)
		require.Nil(t, err)
	}()

	assert.True(t, eassert.ContainsIgnoringFields(allProductions, production, "UpdatedAt"))
	assert.True(t, eassert.ContainsIgnoringFields(allProductions, newProd, "UpdatedAt"))
}

func TestIT_PlanetResourceProductionRepository_Create_WhenDuplicateResourceProductionForBuilding_ExpectFailure(t *testing.T) {
	repo, conn, tx := newTestPlanetResourceProductionRepositoryAndTransaction(t)
	defer conn.Close(context.Background())
	planet, _, _ := insertTestPlanetForPlayer(t, conn)
	production, building, resource := insertTestPlanetResourceProduction(t, conn, planet.Id)

	newProd := persistence.PlanetResourceProduction{
		Planet:     planet.Id,
		Resource:   resource.Id,
		Building:   &building.Id,
		Production: production.Production * 2,
		CreatedAt:  time.Now(),
	}

	_, err := repo.Create(context.Background(), tx, newProd)
	tx.Close(context.Background())

	assert.True(t, errors.IsErrorWithCode(err, pgx.UniqueConstraintViolation), "Actual err: %v", err)
	assertPlanetResourceProductionForBuilding(t, conn, planet.Id, resource.Id, building.Id, production.Production)
}

func TestIT_PlanetResourceProductionRepository_GetForPlanetAndBuilding(t *testing.T) {
	repo, conn, tx := newTestPlanetResourceProductionRepositoryAndTransaction(t)
	defer conn.Close(context.Background())
	planet1, _, _ := insertTestPlanetForPlayer(t, conn)
	building1 := insertTestBuilding(t, conn)
	building2 := insertTestBuilding(t, conn)
	insertTestPlanetResourceProductionForBuilding(t, conn, planet1.Id, nil)
	prp2, _ := insertTestPlanetResourceProductionForBuilding(t, conn, planet1.Id, &building1.Id)
	insertTestPlanetResourceProductionForBuilding(t, conn, planet1.Id, &building2.Id)

	planet2, _, _ := insertTestPlanetForPlayer(t, conn)
	insertTestPlanetResourceProduction(t, conn, planet2.Id)

	actual, err := repo.GetForPlanetAndBuilding(context.Background(), tx, planet1.Id, &building1.Id)
	tx.Close(context.Background())

	assert.Nil(t, err)
	assert.True(t, eassert.EqualsIgnoringFields(actual, prp2))
}

func TestIT_PlanetResourceProductionRepository_GetForPlanetAndBuilding_WhenBuildingIsNull_ExpectSuccess(t *testing.T) {
	repo, conn, tx := newTestPlanetResourceProductionRepositoryAndTransaction(t)
	defer conn.Close(context.Background())
	planet1, _, _ := insertTestPlanetForPlayer(t, conn)
	building1 := insertTestBuilding(t, conn)
	building2 := insertTestBuilding(t, conn)
	prp1, _ := insertTestPlanetResourceProductionForBuilding(t, conn, planet1.Id, nil)
	insertTestPlanetResourceProductionForBuilding(t, conn, planet1.Id, &building1.Id)
	insertTestPlanetResourceProductionForBuilding(t, conn, planet1.Id, &building2.Id)

	planet2, _, _ := insertTestPlanetForPlayer(t, conn)
	insertTestPlanetResourceProduction(t, conn, planet2.Id)

	actual, err := repo.GetForPlanetAndBuilding(context.Background(), tx, planet1.Id, nil)
	tx.Close(context.Background())

	assert.Nil(t, err)
	assert.True(t, eassert.EqualsIgnoringFields(actual, prp1))
}

func TestIT_PlanetResourceProductionRepository_ListForPlanet(t *testing.T) {
	repo, conn, tx := newTestPlanetResourceProductionRepositoryAndTransaction(t)
	defer conn.Close(context.Background())
	planet1, _, _ := insertTestPlanetForPlayer(t, conn)
	prp1, _ := insertTestPlanetResourceProductionForBuilding(t, conn, planet1.Id, nil)
	prp2, _, _ := insertTestPlanetResourceProduction(t, conn, planet1.Id)
	planet2, _, _ := insertTestPlanetForPlayer(t, conn)
	_, _, r3 := insertTestPlanetResourceProduction(t, conn, planet2.Id)

	actual, err := repo.ListForPlanet(context.Background(), tx, planet1.Id)
	tx.Close(context.Background())

	assert.Nil(t, err)
	assert.GreaterOrEqual(t, len(actual), 2)
	assert.True(t, eassert.ContainsIgnoringFields(actual, prp1))
	assert.True(t, eassert.ContainsIgnoringFields(actual, prp2))
	for _, planetResource := range actual {
		assert.NotEqual(t, planetResource.Planet, planet2.Id)
		assert.NotEqual(t, planetResource.Resource, r3.Id)
	}
}

func TestIT_PlanetResourceProductionRepository_Update(t *testing.T) {
	repo, conn, tx := newTestPlanetResourceProductionRepositoryAndTransaction(t)
	defer conn.Close(context.Background())
	planet, _, _ := insertTestPlanetForPlayer(t, conn)
	production, building, resource := insertTestPlanetResourceProduction(t, conn, planet.Id)

	updatedProduction := production
	updatedProduction.UpdatedAt = production.UpdatedAt.Add(18 * time.Second)
	updatedProduction.Production = production.Production * 7

	actual, err := repo.Update(context.Background(), tx, updatedProduction)
	tx.Close(context.Background())

	assert.Nil(t, err)

	expected := persistence.PlanetResourceProduction{
		Planet:     planet.Id,
		Resource:   resource.Id,
		Building:   &building.Id,
		Production: production.Production * 7,
		CreatedAt:  production.CreatedAt,
		UpdatedAt:  updatedProduction.UpdatedAt,
		Version:    production.Version + 1,
	}

	assert.True(t, eassert.EqualsIgnoringFields(actual, expected))
}

func TestIT_PlanetResourceProductionRepository_Update_WhenVersionIsWrong_ExpectOptimisticLockException(t *testing.T) {
	repo, conn, tx := newTestPlanetResourceProductionRepositoryAndTransaction(t)
	defer conn.Close(context.Background())
	planet, _, _ := insertTestPlanetForPlayer(t, conn)
	production, _, _ := insertTestPlanetResourceProduction(t, conn, planet.Id)

	updatedProduction := production
	updatedProduction.Production = production.Production * 4
	updatedProduction.Version = production.Version + 2

	_, err := repo.Update(context.Background(), tx, updatedProduction)
	tx.Close(context.Background())

	assert.True(t, errors.IsErrorWithCode(err, OptimisticLockException), "Actual err: %v", err)
}

func TestIT_PlanetResourceProductionRepository_Update_WithoutBuilding(t *testing.T) {
	repo, conn, tx := newTestPlanetResourceProductionRepositoryAndTransaction(t)
	defer conn.Close(context.Background())
	planet, _, _ := insertTestPlanetForPlayer(t, conn)
	production, resource := insertTestPlanetResourceProductionForBuilding(t, conn, planet.Id, nil)

	updatedProduction := production
	updatedProduction.UpdatedAt = production.UpdatedAt.Add(18 * time.Second)
	updatedProduction.Production = production.Production * 7

	actual, err := repo.Update(context.Background(), tx, updatedProduction)
	tx.Close(context.Background())

	assert.Nil(t, err)

	expected := persistence.PlanetResourceProduction{
		Planet:     planet.Id,
		Resource:   resource.Id,
		Building:   nil,
		Production: production.Production * 7,
		CreatedAt:  production.CreatedAt,
		UpdatedAt:  updatedProduction.UpdatedAt,
		Version:    production.Version + 1,
	}

	assert.True(t, eassert.EqualsIgnoringFields(actual, expected))
}

func TestIT_PlanetResourceProductionRepository_Update_WithoutBuilding_WhenVersionIsWrong_ExpectOptimisticLockException(t *testing.T) {
	repo, conn, tx := newTestPlanetResourceProductionRepositoryAndTransaction(t)
	defer conn.Close(context.Background())
	planet, _, _ := insertTestPlanetForPlayer(t, conn)
	production, _ := insertTestPlanetResourceProductionForBuilding(t, conn, planet.Id, nil)

	updatedProduction := production
	updatedProduction.Production = production.Production * 4
	updatedProduction.Version = production.Version + 2

	_, err := repo.Update(context.Background(), tx, updatedProduction)
	tx.Close(context.Background())

	assert.True(t, errors.IsErrorWithCode(err, OptimisticLockException), "Actual err: %v", err)
}

func TestIT_PlanetResourceProductionRepository_Update_BumpsUpdatedAt(t *testing.T) {
	repo, conn := newTestPlanetResourceProductionRepository(t)
	defer conn.Close(context.Background())
	planet, _, _ := insertTestPlanetForPlayer(t, conn)
	production, _, _ := insertTestPlanetResourceProduction(t, conn, planet.Id)

	updatedProduction := production
	updatedProduction.UpdatedAt = production.UpdatedAt.Add(2 * time.Second)
	updatedProduction.Production = production.Production * 3

	func() {
		tx, err := conn.BeginTx(context.Background())
		require.Nil(t, err)
		defer tx.Close(context.Background())

		_, err = repo.Update(context.Background(), tx, updatedProduction)
		assert.Nil(t, err)
	}()

	var updatedProductionFromDb persistence.PlanetResourceProduction
	func() {
		tx, err := conn.BeginTx(context.Background())
		require.Nil(t, err)
		defer tx.Close(context.Background())

		allProductions, err := repo.ListForPlanet(context.Background(), tx, planet.Id)
		require.Nil(t, err)
		assert.Len(t, allProductions, 1)

		updatedProductionFromDb = allProductions[0]
	}()

	assert.Equal(t, updatedProductionFromDb.UpdatedAt, updatedProductionFromDb.UpdatedAt)
}

func TestIT_PlanetResourceProductionRepository_Update_BumpsVersion(t *testing.T) {
	repo, conn := newTestPlanetResourceProductionRepository(t)
	defer conn.Close(context.Background())
	planet, _, _ := insertTestPlanetForPlayer(t, conn)
	production, _, _ := insertTestPlanetResourceProduction(t, conn, planet.Id)

	updatedProduction := production
	updatedProduction.Production = production.Production * 3

	func() {
		tx, err := conn.BeginTx(context.Background())
		require.Nil(t, err)
		defer tx.Close(context.Background())

		_, err = repo.Update(context.Background(), tx, updatedProduction)
		assert.Nil(t, err)
	}()

	var updatedProductionFromDb persistence.PlanetResourceProduction
	func() {
		tx, err := conn.BeginTx(context.Background())
		require.Nil(t, err)
		defer tx.Close(context.Background())

		allProductions, err := repo.ListForPlanet(context.Background(), tx, planet.Id)
		require.Nil(t, err)
		assert.Len(t, allProductions, 1)

		updatedProductionFromDb = allProductions[0]
	}()

	assert.Equal(t, production.Version+1, updatedProductionFromDb.Version)
}

func newTestPlanetResourceProductionRepository(t *testing.T) (PlanetResourceProductionRepository, db.Connection) {
	conn := newTestConnection(t)
	return NewPlanetResourceProductionRepository(), conn
}

func newTestPlanetResourceProductionRepositoryAndTransaction(t *testing.T) (PlanetResourceProductionRepository, db.Connection, db.Transaction) {
	repo, conn := newTestPlanetResourceProductionRepository(t)
	tx, err := conn.BeginTx(context.Background())
	require.Nil(t, err)
	return repo, conn, tx
}

func insertTestPlanetResourceProductionForBuilding(t *testing.T, conn db.Connection, planet uuid.UUID, building *uuid.UUID) (persistence.PlanetResourceProduction, persistence.Resource) {
	someTime := time.Date(2024, 12, 1, 17, 18, 25, 0, time.UTC)

	resource := insertTestResource(t, conn)

	production := persistence.PlanetResourceProduction{
		Planet:     planet,
		Building:   building,
		Resource:   resource.Id,
		Production: 7432,
		CreatedAt:  someTime,
		UpdatedAt:  someTime,
	}

	sqlQuery := `INSERT INTO planet_resource_production (planet, building, resource, production, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := conn.Exec(
		context.Background(),
		sqlQuery,
		production.Planet,
		production.Building,
		production.Resource,
		production.Production,
		production.CreatedAt,
		production.UpdatedAt,
	)
	require.Nil(t, err)

	return production, resource
}

func insertTestPlanetResourceProduction(t *testing.T, conn db.Connection, planet uuid.UUID) (persistence.PlanetResourceProduction, persistence.Building, persistence.Resource) {
	building := insertTestBuilding(t, conn)

	production, resource := insertTestPlanetResourceProductionForBuilding(t, conn, planet, &building.Id)

	return production, building, resource
}

func assertPlanetResourceProductionExists(t *testing.T, conn db.Connection, planet uuid.UUID, resource uuid.UUID) {
	sqlQuery := `SELECT COUNT(*) FROM planet_resource_production WHERE planet = $1 AND resource = $2`
	value, err := db.QueryOne[int](context.Background(), conn, sqlQuery, planet, resource)
	require.Nil(t, err)
	require.Equal(t, 1, value)
}

func assertPlanetResourceProductionDoesNotExist(t *testing.T, conn db.Connection, planet uuid.UUID) {
	sqlQuery := `SELECT COUNT(resource) FROM planet_resource_production WHERE planet = $1`
	value, err := db.QueryOne[int](context.Background(), conn, sqlQuery, planet)
	require.Nil(t, err)
	require.Zero(t, value)
}

func assertPlanetResourceProduction(t *testing.T, conn db.Connection, planet uuid.UUID, resource uuid.UUID, prod int) {
	type productionData struct {
		Production int
		Building   *string
	}

	sqlQuery := `SELECT production, building FROM planet_resource_production WHERE planet = $1 AND resource = $2`
	value, err := db.QueryOne[productionData](context.Background(), conn, sqlQuery, planet, resource)
	require.Nil(t, err)
	require.Equal(t, prod, value.Production)
	require.Nil(t, value.Building)
}

func assertPlanetResourceProductionWithoutBuilding(t *testing.T, conn db.Connection, planet uuid.UUID, resource uuid.UUID, prod int) {
	sqlQuery := `SELECT production FROM planet_resource_production WHERE planet = $1 AND resource = $2 AND building IS NULL`
	value, err := db.QueryOne[int](context.Background(), conn, sqlQuery, planet, resource)
	require.Nil(t, err)
	require.Equal(t, prod, value)
}

func assertPlanetResourceProductionForBuilding(t *testing.T, conn db.Connection, planet uuid.UUID, resource uuid.UUID, building uuid.UUID, prod int) {
	sqlQuery := `SELECT production FROM planet_resource_production WHERE planet = $1 AND resource = $2 AND building = $3`
	value, err := db.QueryOne[int](context.Background(), conn, sqlQuery, planet, resource, building)
	require.Nil(t, err)
	require.Equal(t, prod, value)
}
