package drivenadapters

import (
	"context"
	"fmt"
	"testing"

	"github.com/Knoblauchpilze/backend-toolkit/pkg/db"
	"github.com/Knoblauchpilze/backend-toolkit/pkg/db/pgx"
	"github.com/Knoblauchpilze/backend-toolkit/pkg/errors"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	drivenports "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/ports/driven"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIT_UniverseRepository_Create(t *testing.T) {
	repo, conn := newTestUniverseRepository(t)
	defer conn.Close(context.Background())

	t.Run("creates a universe", func(t *testing.T) {
		universe := models.Universe{
			Id:        uuid.New(),
			Name:      fmt.Sprintf("universe-%s", uuid.NewString()),
			CreatedAt: someTime,
		}

		err := repo.Create(context.Background(), universe)
		require.NoError(t, err, "Actual err: %v", err)
		assertUniverseExists(t, conn, universe.Id)

		actual, err := repo.Get(context.Background(), universe.Id)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, universe, actual)
	})

	t.Run("returns error when universe with same name already exists", func(t *testing.T) {
		universe := insertTestUniverse(t, conn)

		newUniverse := models.Universe{
			Id:        uuid.New(),
			Name:      universe.Name,
			CreatedAt: someTime,
		}

		err := repo.Create(context.Background(), newUniverse)

		assert.True(t, errors.IsErrorWithCode(err, pgx.UniqueConstraintViolation), "Actual err: %v", err)
		assertUniverseDoesNotExist(t, conn, newUniverse.Id)
	})

}

func TestIT_UniverseRepository_Get(t *testing.T) {
	repo, conn := newTestUniverseRepository(t)
	defer conn.Close(context.Background())

	t.Run("gets a universe", func(t *testing.T) {
		universe := insertTestUniverse(t, conn)

		actual, err := repo.Get(context.Background(), universe.Id)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, actual, universe)
	})

	t.Run("returns error when universe does not exist", func(t *testing.T) {
		id := uuid.MustParse("00000000-1111-2222-1111-000000000000")
		_, err := repo.Get(context.Background(), id)

		assert.True(t, errors.IsErrorWithCode(err, db.NoMatchingRows), "Actual err: %v", err)
	})
}

func TestIT_UniverseRepository_List(t *testing.T) {
	repo, conn := newTestUniverseRepository(t)
	defer conn.Close(context.Background())
	u1 := insertTestUniverse(t, conn)
	u2 := insertTestUniverse(t, conn)

	actual, err := repo.List(context.Background())
	require.NoError(t, err, "Actual err: %v", err)

	assert.GreaterOrEqual(t, len(actual), 2)
	assert.Contains(t, actual, u1)
	assert.Contains(t, actual, u2)
}

func TestIT_UniverseRepository_Delete(t *testing.T) {
	repo, conn := newTestUniverseRepository(t)
	defer conn.Close(context.Background())

	t.Run("deletes universe", func(t *testing.T) {
		universe := insertTestUniverse(t, conn)

		err := repo.Delete(context.Background(), universe.Id)
		require.NoError(t, err, "Actual err: %v", err)

		assertUniverseDoesNotExist(t, conn, universe.Id)
	})

	t.Run("succeeds when the universe does not exist", func(t *testing.T) {
		nonExistingId := uuid.MustParse("00000000-0000-1221-0000-000000000000")

		err := repo.Delete(context.Background(), nonExistingId)
		require.NoError(t, err, "Actual err: %v", err)
	})
}

func newTestUniverseRepository(t *testing.T) (drivenports.ForManagingUniverses, db.Connection) {
	t.Helper()
	conn := newTestConnection(t)
	return NewUniverseRepository(conn), conn
}

func insertTestUniverse(t *testing.T, conn db.Connection) models.Universe {
	t.Helper()

	universe := models.Universe{
		Id:        uuid.New(),
		Name:      fmt.Sprintf("my-universe-%s", uuid.NewString()),
		CreatedAt: someTime,
	}

	sqlQuery := `INSERT INTO universe (id, name, created_at) VALUES ($1, $2, $3)`
	_, err := conn.Exec(
		context.Background(),
		sqlQuery,
		universe.Id,
		universe.Name,
		universe.CreatedAt,
	)
	require.NoError(t, err, "Actual err: %v", err)

	return universe
}

func assertUniverseExists(t *testing.T, conn db.Connection, id uuid.UUID) {
	t.Helper()

	sqlQuery := `SELECT id FROM universe WHERE id = $1`
	value, err := db.QueryOne[uuid.UUID](context.Background(), conn, sqlQuery, id)
	require.NoError(t, err, "Actual err: %v", err)
	require.Equal(t, id, value)
}

func assertUniverseDoesNotExist(t *testing.T, conn db.Connection, id uuid.UUID) {
	t.Helper()

	sqlQuery := `SELECT COUNT(id) FROM universe WHERE id = $1`
	value, err := db.QueryOne[int](context.Background(), conn, sqlQuery, id)
	require.NoError(t, err, "Actual err: %v", err)
	require.Zero(t, value)
}
