package repositories

import (
	"context"
	"testing"
	"time"

	"github.com/Knoblauchpilze/backend-toolkit/pkg/db"
	"github.com/Knoblauchpilze/backend-toolkit/pkg/db/pgx"
	"github.com/Knoblauchpilze/backend-toolkit/pkg/errors"
	eassert "github.com/Knoblauchpilze/easy-assert/assert"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/persistence"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIT_PlanetResourceRepository_Create(t *testing.T) {
	repo, conn, tx := newTestPlanetResourceRepositoryAndTransaction(t)
	defer conn.Close(context.Background())
	planet, _, _ := insertTestPlanetForPlayer(t, conn)
	resource := insertTestResource(t, conn)

	planetResource := persistence.PlanetResource{
		Planet:    planet.Id,
		Resource:  resource.Id,
		Amount:    29,
		CreatedAt: time.Now(),
	}

	actual, err := repo.Create(context.Background(), tx, planetResource)
	assert.Nil(t, err)
	tx.Close(context.Background())

	assert.True(t, eassert.EqualsIgnoringFields(actual, planetResource, "UpdatedAt"))
	assert.Equal(t, actual.UpdatedAt, actual.CreatedAt)
	assertPlanetResourceExists(t, conn, planet.Id, resource.Id)
}

func TestIT_PlanetResourceRepository_Create_WhenDuplicateResource_ExpectFailure(t *testing.T) {
	repo, conn, tx := newTestPlanetResourceRepositoryAndTransaction(t)
	defer conn.Close(context.Background())
	planet, _, _ := insertTestPlanetForPlayer(t, conn)
	planetResource, resource := insertTestPlanetResource(t, conn, planet.Id)

	newResource := persistence.PlanetResource{
		Planet:    planet.Id,
		Resource:  resource.Id,
		Amount:    planetResource.Amount * 2,
		CreatedAt: time.Now(),
	}

	_, err := repo.Create(context.Background(), tx, newResource)
	tx.Close(context.Background())

	assert.True(t, errors.IsErrorWithCode(err, pgx.UniqueConstraintViolation), "Actual err: %v", err)
	assertPlanetResourceAmount(t, conn, planet.Id, resource.Id, planetResource.Amount)
}

func TestIT_PlanetResourceRepository_ListForPlanet(t *testing.T) {
	repo, conn, tx := newTestPlanetResourceRepositoryAndTransaction(t)
	defer conn.Close(context.Background())
	planet1, _, _ := insertTestPlanetForPlayer(t, conn)
	pr1, _ := insertTestPlanetResource(t, conn, planet1.Id)
	pr2, _ := insertTestPlanetResource(t, conn, planet1.Id)
	planet2, _, _ := insertTestPlanetForPlayer(t, conn)
	_, r3 := insertTestPlanetResource(t, conn, planet2.Id)

	actual, err := repo.ListForPlanet(context.Background(), tx, planet1.Id)
	tx.Close(context.Background())

	assert.Nil(t, err)
	assert.GreaterOrEqual(t, len(actual), 2)
	assert.True(t, eassert.ContainsIgnoringFields(actual, pr1))
	assert.True(t, eassert.ContainsIgnoringFields(actual, pr2))
	for _, planetResource := range actual {
		assert.NotEqual(t, planetResource.Planet, planet2.Id)
		assert.NotEqual(t, planetResource.Resource, r3.Id)
	}
}

func TestIT_PlanetResourceRepository_Update(t *testing.T) {
	repo, conn, tx := newTestPlanetResourceRepositoryAndTransaction(t)
	defer conn.Close(context.Background())
	planet, _, _ := insertTestPlanetForPlayer(t, conn)
	planetResource, resource := insertTestPlanetResource(t, conn, planet.Id)

	updatedPlanetResource := planetResource
	updatedPlanetResource.UpdatedAt = planetResource.UpdatedAt.Add(3 * time.Second)
	updatedPlanetResource.Amount = planetResource.Amount * 3

	actual, err := repo.Update(context.Background(), tx, updatedPlanetResource)
	tx.Close(context.Background())

	assert.Nil(t, err)

	expected := persistence.PlanetResource{
		Planet:    planet.Id,
		Resource:  resource.Id,
		Amount:    planetResource.Amount * 3,
		CreatedAt: planetResource.CreatedAt,
		UpdatedAt: updatedPlanetResource.UpdatedAt,
		Version:   planetResource.Version + 1,
	}
	assert.True(t, eassert.EqualsIgnoringFields(actual, expected))
}

func TestIT_PlanetResourceRepository_Update_WhenVersionIsWrong_ExpectOptimisticLockException(t *testing.T) {
	repo, conn, tx := newTestPlanetResourceRepositoryAndTransaction(t)
	defer conn.Close(context.Background())
	planet, _, _ := insertTestPlanetForPlayer(t, conn)
	planetResource, _ := insertTestPlanetResource(t, conn, planet.Id)

	updatedPlanetResource := planetResource
	updatedPlanetResource.Amount = planetResource.Amount * 4
	updatedPlanetResource.Version = planetResource.Version + 2

	_, err := repo.Update(context.Background(), tx, updatedPlanetResource)
	tx.Close(context.Background())

	assert.True(t, errors.IsErrorWithCode(err, OptimisticLockException), "Actual err: %v", err)
}

func TestIT_PlanetResourceRepository_Update_BumpsUpdatedAt(t *testing.T) {
	repo, conn := newTestPlanetResourceRepository(t)
	defer conn.Close(context.Background())
	planet, _, _ := insertTestPlanetForPlayer(t, conn)
	planetResource, _ := insertTestPlanetResource(t, conn, planet.Id)

	updatedPlanetResource := planetResource
	updatedPlanetResource.UpdatedAt = planetResource.UpdatedAt.Add(2 * time.Second)
	updatedPlanetResource.Amount = planetResource.Amount * 3

	func() {
		tx, err := conn.BeginTx(context.Background())
		require.Nil(t, err)
		defer tx.Close(context.Background())

		_, err = repo.Update(context.Background(), tx, updatedPlanetResource)
		assert.Nil(t, err)
	}()

	var updatedPlanetResourceFromDb persistence.PlanetResource
	func() {
		tx, err := conn.BeginTx(context.Background())
		require.Nil(t, err)
		defer tx.Close(context.Background())

		allResources, err := repo.ListForPlanet(context.Background(), tx, planet.Id)
		require.Nil(t, err)
		assert.Len(t, allResources, 1)

		updatedPlanetResourceFromDb = allResources[0]
	}()

	assert.True(t, updatedPlanetResource.UpdatedAt.Equal(updatedPlanetResourceFromDb.UpdatedAt))
}

func TestIT_PlanetResourceRepository_Update_BumpsVersion(t *testing.T) {
	repo, conn := newTestPlanetResourceRepository(t)
	defer conn.Close(context.Background())
	planet, _, _ := insertTestPlanetForPlayer(t, conn)
	planetResource, _ := insertTestPlanetResource(t, conn, planet.Id)

	updatedPlanetResource := planetResource
	updatedPlanetResource.Amount = planetResource.Amount * 3

	func() {
		tx, err := conn.BeginTx(context.Background())
		require.Nil(t, err)
		defer tx.Close(context.Background())

		_, err = repo.Update(context.Background(), tx, updatedPlanetResource)
		assert.Nil(t, err)
	}()

	var updatedPlanetResourceFromDb persistence.PlanetResource
	func() {
		tx, err := conn.BeginTx(context.Background())
		require.Nil(t, err)
		defer tx.Close(context.Background())

		allResources, err := repo.ListForPlanet(context.Background(), tx, planet.Id)
		require.Nil(t, err)
		assert.Len(t, allResources, 1)

		updatedPlanetResourceFromDb = allResources[0]
	}()

	assert.Equal(t, planetResource.Version+1, updatedPlanetResourceFromDb.Version)
}

func newTestPlanetResourceRepository(t *testing.T) (PlanetResourceRepository, db.Connection) {
	conn := newTestConnection(t)
	return NewPlanetResourceRepository(), conn
}

func newTestPlanetResourceRepositoryAndTransaction(t *testing.T) (PlanetResourceRepository, db.Connection, db.Transaction) {
	repo, conn := newTestPlanetResourceRepository(t)
	tx, err := conn.BeginTx(context.Background())
	require.Nil(t, err)
	return repo, conn, tx
}

func insertTestPlanetResource(t *testing.T, conn db.Connection, planet uuid.UUID) (persistence.PlanetResource, persistence.Resource) {
	someTime := time.Date(2024, 12, 1, 10, 42, 30, 0, time.UTC)

	resource := insertTestResource(t, conn)

	planetResource := persistence.PlanetResource{
		Planet:    planet,
		Resource:  resource.Id,
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

	return planetResource, resource
}

func assertPlanetResourceExists(t *testing.T, conn db.Connection, planet uuid.UUID, resource uuid.UUID) {
	sqlQuery := `SELECT COUNT(*) FROM planet_resource WHERE planet = $1 AND resource = $2`
	value, err := db.QueryOne[int](context.Background(), conn, sqlQuery, planet, resource)
	require.Nil(t, err)
	require.Equal(t, 1, value)
}

func assertPlanetResourceDoesNotExist(t *testing.T, conn db.Connection, planet uuid.UUID) {
	sqlQuery := `SELECT COUNT(resource) FROM planet_resource WHERE planet = $1`
	value, err := db.QueryOne[int](context.Background(), conn, sqlQuery, planet)
	require.Nil(t, err)
	require.Zero(t, value)
}

func assertPlanetResourceAmount(t *testing.T, conn db.Connection, planet uuid.UUID, resource uuid.UUID, amount float64) {
	sqlQuery := `SELECT amount FROM planet_resource WHERE planet = $1 AND resource = $2`
	value, err := db.QueryOne[float64](context.Background(), conn, sqlQuery, planet, resource)
	require.Nil(t, err)
	require.Equal(t, amount, value)
}
