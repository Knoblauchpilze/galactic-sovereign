package repositories

import (
	"context"
	"fmt"
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

func TestIT_UniverseRepository_Create(t *testing.T) {
	repo, conn := newTestUniverseRepository(t)

	universe := persistence.Universe{
		Id:        uuid.New(),
		Name:      fmt.Sprintf("universe-%s", uuid.NewString()),
		CreatedAt: time.Now(),
	}

	actual, err := repo.Create(context.Background(), universe)
	assert.Nil(t, err)

	assert.True(t, eassert.EqualsIgnoringFields(actual, universe, "UpdatedAt"))
	assert.True(t, eassert.AreTimeCloserThan(actual.UpdatedAt, actual.CreatedAt, 1*time.Second))
	assertUniverseExists(t, conn, universe.Id)
}

func TestIT_UniverseRepository_Create_WhenDuplicateName_ExpectFailure(t *testing.T) {
	repo, conn := newTestUniverseRepository(t)
	universe := insertTestUniverse(t, conn)

	newUniverse := persistence.Universe{
		Id:        uuid.New(),
		Name:      universe.Name,
		CreatedAt: time.Now(),
	}

	_, err := repo.Create(context.Background(), newUniverse)

	assert.True(t, errors.IsErrorWithCode(err, pgx.UniqueConstraintViolation), "Actual err: %v", err)
	assertUniverseDoesNotExist(t, conn, newUniverse.Id)
}

func TestIT_UniverseRepository_Get(t *testing.T) {
	repo, conn, tx := newTestUniverseRepositoryAndTransaction(t)
	defer conn.Close(context.Background())
	defer tx.Close(context.Background())
	universe := insertTestUniverse(t, conn)

	actual, err := repo.Get(context.Background(), tx, universe.Id)
	assert.Nil(t, err)

	assert.True(t, eassert.EqualsIgnoringFields(actual, universe))
}

func TestIT_UniverseRepository_Get_WhenNotFound_ExpectFailure(t *testing.T) {
	repo, conn, tx := newTestUniverseRepositoryAndTransaction(t)
	defer conn.Close(context.Background())
	defer tx.Close(context.Background())

	// Non-existent id
	id := uuid.MustParse("00000000-1111-2222-1111-000000000000")
	_, err := repo.Get(context.Background(), tx, id)
	assert.True(t, errors.IsErrorWithCode(err, db.NoMatchingRows), "Actual err: %v", err)
}

func TestIT_UniverseRepository_List(t *testing.T) {
	repo, conn := newTestUniverseRepository(t)
	u1 := insertTestUniverse(t, conn)
	u2 := insertTestUniverse(t, conn)

	actual, err := repo.List(context.Background())

	assert.Nil(t, err)
	assert.GreaterOrEqual(t, len(actual), 2)
	assert.True(t, eassert.ContainsIgnoringFields(actual, u1))
	assert.True(t, eassert.ContainsIgnoringFields(actual, u2))
}

func TestIT_UniverseRepository_Delete(t *testing.T) {
	repo, conn, tx := newTestUniverseRepositoryAndTransaction(t)
	defer conn.Close(context.Background())
	universe := insertTestUniverse(t, conn)

	err := repo.Delete(context.Background(), tx, universe.Id)
	tx.Close(context.Background())

	assert.Nil(t, err)
	assertUniverseDoesNotExist(t, conn, universe.Id)
}

func TestIT_UniverseRepository_Delete_WhenNotFound_ExpectSuccess(t *testing.T) {
	repo, conn, tx := newTestUniverseRepositoryAndTransaction(t)
	defer conn.Close(context.Background())
	nonExistingId := uuid.MustParse("00000000-0000-1221-0000-000000000000")

	err := repo.Delete(context.Background(), tx, nonExistingId)
	tx.Close(context.Background())

	assert.Nil(t, err)
}

func newTestUniverseRepository(t *testing.T) (UniverseRepository, db.Connection) {
	conn := newTestConnection(t)
	return NewUniverseRepository(conn), conn
}

func newTestUniverseRepositoryAndTransaction(t *testing.T) (UniverseRepository, db.Connection, db.Transaction) {
	repo, conn := newTestUniverseRepository(t)
	tx, err := conn.BeginTx(context.Background())
	require.Nil(t, err)
	return repo, conn, tx
}

func insertTestUniverse(t *testing.T, conn db.Connection) persistence.Universe {
	someTime := time.Date(2024, 11, 29, 17, 53, 29, 0, time.UTC)

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

func assertUniverseExists(t *testing.T, conn db.Connection, id uuid.UUID) {
	sqlQuery := `SELECT id FROM universe WHERE id = $1`
	value, err := db.QueryOne[uuid.UUID](context.Background(), conn, sqlQuery, id)
	require.Nil(t, err)
	require.Equal(t, id, value)
}

func assertUniverseDoesNotExist(t *testing.T, conn db.Connection, id uuid.UUID) {
	sqlQuery := `SELECT COUNT(id) FROM universe WHERE id = $1`
	value, err := db.QueryOne[int](context.Background(), conn, sqlQuery, id)
	require.Nil(t, err)
	require.Zero(t, value)
}
