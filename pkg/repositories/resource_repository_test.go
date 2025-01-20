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

func TestIT_ResourceRepository_List(t *testing.T) {
	repo, conn, tx := newTestResourceRepositoryAndTransaction(t)
	defer conn.Close(context.Background())
	defer tx.Close(context.Background())
	r1 := insertTestResource(t, conn)
	r2 := insertTestResource(t, conn)

	actual, err := repo.List(context.Background(), tx)

	assert.Nil(t, err)
	assert.GreaterOrEqual(t, len(actual), 2)
	assert.True(t, eassert.ContainsIgnoringFields(actual, r1))
	assert.True(t, eassert.ContainsIgnoringFields(actual, r2))
}

func newTestResourceRepositoryAndTransaction(t *testing.T) (ResourceRepository, db.Connection, db.Transaction) {
	conn := newTestConnection(t)
	tx, err := conn.BeginTx(context.Background())
	require.Nil(t, err)
	return NewResourceRepository(), conn, tx
}

func insertTestResource(t *testing.T, conn db.Connection) persistence.Resource {
	someTime := time.Date(2024, 11, 30, 9, 23, 53, 0, time.UTC)

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
