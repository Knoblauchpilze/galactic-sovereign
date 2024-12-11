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

func TestIT_PlanetResourceStorageRepository_Create(t *testing.T) {
	repo, conn, tx := newTestPlanetResourceStorageRepositoryAndTransaction(t)
	defer conn.Close(context.Background())
	planet, _, _ := insertTestPlanetForPlayer(t, conn)
	resource := insertTestResource(t, conn)

	storage := persistence.PlanetResourceStorage{
		Planet:    planet.Id,
		Resource:  resource.Id,
		Storage:   29,
		CreatedAt: time.Now(),
	}

	actual, err := repo.Create(context.Background(), tx, storage)
	assert.Nil(t, err)
	tx.Close(context.Background())

	assert.True(t, eassert.EqualsIgnoringFields(actual, storage, "UpdatedAt"))
	assert.Equal(t, actual.UpdatedAt, actual.CreatedAt)
	assertPlanetResourceStorageExists(t, conn, planet.Id, resource.Id)
}

func TestIT_PlanetResourceStorageRepository_Create_WhenDuplicateResource_ExpectFailure(t *testing.T) {
	repo, conn, tx := newTestPlanetResourceStorageRepositoryAndTransaction(t)
	defer conn.Close(context.Background())
	planet, _, _ := insertTestPlanetForPlayer(t, conn)
	storage, resource := insertTestPlanetResourceStorage(t, conn, planet.Id)

	newStorage := persistence.PlanetResourceStorage{
		Planet:    planet.Id,
		Resource:  resource.Id,
		Storage:   storage.Storage * 2,
		CreatedAt: time.Now(),
	}

	_, err := repo.Create(context.Background(), tx, newStorage)
	tx.Close(context.Background())

	assert.True(t, errors.IsErrorWithCode(err, pgx.UniqueConstraintViolation), "Actual err: %v", err)
	assertPlanetResourceStorage(t, conn, planet.Id, resource.Id, storage.Storage)
}

func TestIT_PlanetResourceStorageRepository_ListForPlanet(t *testing.T) {
	repo, conn, tx := newTestPlanetResourceStorageRepositoryAndTransaction(t)
	defer conn.Close(context.Background())
	planet1, _, _ := insertTestPlanetForPlayer(t, conn)
	ps1, _ := insertTestPlanetResourceStorage(t, conn, planet1.Id)
	ps2, _ := insertTestPlanetResourceStorage(t, conn, planet1.Id)
	planet2, _, _ := insertTestPlanetForPlayer(t, conn)
	_, r3 := insertTestPlanetResourceStorage(t, conn, planet2.Id)

	actual, err := repo.ListForPlanet(context.Background(), tx, planet1.Id)
	tx.Close(context.Background())

	assert.Nil(t, err)
	assert.GreaterOrEqual(t, len(actual), 2)
	assert.True(t, eassert.ContainsIgnoringFields(actual, ps1))
	assert.True(t, eassert.ContainsIgnoringFields(actual, ps2))
	for _, storage := range actual {
		assert.NotEqual(t, storage.Planet, planet2.Id)
		assert.NotEqual(t, storage.Resource, r3.Id)
	}
}

func TestIT_PlanetResourceStorageRepository_Update(t *testing.T) {
	repo, conn, tx := newTestPlanetResourceStorageRepositoryAndTransaction(t)
	defer conn.Close(context.Background())
	planet, _, _ := insertTestPlanetForPlayer(t, conn)
	storage, resource := insertTestPlanetResourceStorage(t, conn, planet.Id)

	updatedStorage := storage
	updatedStorage.UpdatedAt = storage.UpdatedAt.Add(3 * time.Second)
	updatedStorage.Storage = storage.Storage * 3

	actual, err := repo.Update(context.Background(), tx, updatedStorage)
	tx.Close(context.Background())

	assert.Nil(t, err)

	expected := persistence.PlanetResourceStorage{
		Planet:    planet.Id,
		Resource:  resource.Id,
		Storage:   storage.Storage * 3,
		CreatedAt: storage.CreatedAt,
		UpdatedAt: updatedStorage.UpdatedAt,
		Version:   storage.Version + 1,
	}
	assert.True(t, eassert.EqualsIgnoringFields(actual, expected))
}

func TestIT_PlanetResourceStorageRepository_Update_WhenVersionIsWrong_ExpectOptimisticLockException(t *testing.T) {
	repo, conn, tx := newTestPlanetResourceStorageRepositoryAndTransaction(t)
	defer conn.Close(context.Background())
	planet, _, _ := insertTestPlanetForPlayer(t, conn)
	storage, _ := insertTestPlanetResourceStorage(t, conn, planet.Id)

	updatedStorage := storage
	updatedStorage.Storage = storage.Storage * 4
	updatedStorage.Version = storage.Version + 2

	_, err := repo.Update(context.Background(), tx, updatedStorage)
	tx.Close(context.Background())

	assert.True(t, errors.IsErrorWithCode(err, OptimisticLockException), "Actual err: %v", err)
}

func TestIT_PlanetResourceStorageRepository_Update_BumpsUpdatedAt(t *testing.T) {
	repo, conn := newTestPlanetResourceStorageRepository(t)
	defer conn.Close(context.Background())
	planet, _, _ := insertTestPlanetForPlayer(t, conn)
	storage, _ := insertTestPlanetResourceStorage(t, conn, planet.Id)

	updatedStorage := storage
	updatedStorage.UpdatedAt = storage.UpdatedAt.Add(2 * time.Second)
	updatedStorage.Storage = storage.Storage * 3

	func() {
		tx, err := conn.BeginTx(context.Background())
		require.Nil(t, err)
		defer tx.Close(context.Background())

		_, err = repo.Update(context.Background(), tx, updatedStorage)
		assert.Nil(t, err)
	}()

	var updatedStorageFromDb persistence.PlanetResourceStorage
	func() {
		tx, err := conn.BeginTx(context.Background())
		require.Nil(t, err)
		defer tx.Close(context.Background())

		allStorages, err := repo.ListForPlanet(context.Background(), tx, planet.Id)
		require.Nil(t, err)
		assert.Len(t, allStorages, 1)

		updatedStorageFromDb = allStorages[0]
	}()

	assert.True(t, updatedStorage.UpdatedAt.Equal(updatedStorageFromDb.UpdatedAt))
}

func TestIT_PlanetResourceStorageRepository_Update_BumpsVersion(t *testing.T) {
	repo, conn := newTestPlanetResourceStorageRepository(t)
	defer conn.Close(context.Background())
	planet, _, _ := insertTestPlanetForPlayer(t, conn)
	storage, _ := insertTestPlanetResourceStorage(t, conn, planet.Id)

	updatedStorage := storage
	updatedStorage.Storage = storage.Storage * 3

	func() {
		tx, err := conn.BeginTx(context.Background())
		require.Nil(t, err)
		defer tx.Close(context.Background())

		_, err = repo.Update(context.Background(), tx, updatedStorage)
		assert.Nil(t, err)
	}()

	var updatedStorageFromDb persistence.PlanetResourceStorage
	func() {
		tx, err := conn.BeginTx(context.Background())
		require.Nil(t, err)
		defer tx.Close(context.Background())

		allStorages, err := repo.ListForPlanet(context.Background(), tx, planet.Id)
		require.Nil(t, err)
		assert.Len(t, allStorages, 1)

		updatedStorageFromDb = allStorages[0]
	}()

	assert.Equal(t, storage.Version+1, updatedStorageFromDb.Version)
}

func newTestPlanetResourceStorageRepository(t *testing.T) (PlanetResourceStorageRepository, db.Connection) {
	conn := newTestConnection(t)
	return NewPlanetResourceStorageRepository(), conn
}

func newTestPlanetResourceStorageRepositoryAndTransaction(t *testing.T) (PlanetResourceStorageRepository, db.Connection, db.Transaction) {
	repo, conn := newTestPlanetResourceStorageRepository(t)
	tx, err := conn.BeginTx(context.Background())
	require.Nil(t, err)
	return repo, conn, tx
}

func insertTestPlanetResourceStorage(t *testing.T, conn db.Connection, planet uuid.UUID) (persistence.PlanetResourceStorage, persistence.Resource) {
	someTime := time.Date(2024, 12, 1, 21, 55, 27, 0, time.UTC)

	resource := insertTestResource(t, conn)

	planetResourceStorage := persistence.PlanetResourceStorage{
		Planet:    planet,
		Resource:  resource.Id,
		Storage:   6233,
		CreatedAt: someTime,
		UpdatedAt: someTime,
	}

	sqlQuery := `INSERT INTO planet_resource_storage (planet, resource, storage, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)`
	_, err := conn.Exec(
		context.Background(),
		sqlQuery,
		planetResourceStorage.Planet,
		planetResourceStorage.Resource,
		planetResourceStorage.Storage,
		planetResourceStorage.CreatedAt,
		planetResourceStorage.UpdatedAt,
	)
	require.Nil(t, err)

	return planetResourceStorage, resource
}

func assertPlanetResourceStorageExists(t *testing.T, conn db.Connection, planet uuid.UUID, resource uuid.UUID) {
	sqlQuery := `SELECT COUNT(*) FROM planet_resource_storage WHERE planet = $1 AND resource = $2`
	value, err := db.QueryOne[int](context.Background(), conn, sqlQuery, planet, resource)
	require.Nil(t, err)
	require.Equal(t, 1, value)
}

func assertPlanetResourceStorageDoesNotExist(t *testing.T, conn db.Connection, planet uuid.UUID) {
	sqlQuery := `SELECT COUNT(resource) FROM planet_resource_storage WHERE planet = $1`
	value, err := db.QueryOne[int](context.Background(), conn, sqlQuery, planet)
	require.Nil(t, err)
	require.Zero(t, value)
}

func assertPlanetResourceStorage(t *testing.T, conn db.Connection, planet uuid.UUID, resource uuid.UUID, storage int) {
	sqlQuery := `SELECT storage FROM planet_resource_storage WHERE planet = $1 AND resource = $2`
	value, err := db.QueryOne[int](context.Background(), conn, sqlQuery, planet, resource)
	require.Nil(t, err)
	require.Equal(t, storage, value)
}
