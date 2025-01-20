package repositories

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/Knoblauchpilze/backend-toolkit/pkg/db"
	eassert "github.com/Knoblauchpilze/easy-assert/assert"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/persistence"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIT_BuildingRepository_List(t *testing.T) {
	repo, conn, tx := newTestBuildingRepositoryAndTransaction(t)
	defer conn.Close(context.Background())
	defer tx.Close(context.Background())
	b1 := insertTestBuilding(t, conn)
	b2 := insertTestBuilding(t, conn)

	actual, err := repo.List(context.Background(), tx)

	assert.Nil(t, err)
	assert.GreaterOrEqual(t, len(actual), 2)
	assert.True(t, eassert.ContainsIgnoringFields(actual, b1))
	assert.True(t, eassert.ContainsIgnoringFields(actual, b2))
}

func newTestBuildingRepositoryAndTransaction(t *testing.T) (BuildingRepository, db.Connection, db.Transaction) {
	conn := newTestConnection(t)
	tx, err := conn.BeginTx(context.Background())
	require.Nil(t, err)
	return NewBuildingRepository(), conn, tx
}

func insertTestBuilding(t *testing.T, conn db.Connection) persistence.Building {
	someTime := time.Date(2024, 11, 30, 9, 12, 03, 0, time.UTC)

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
