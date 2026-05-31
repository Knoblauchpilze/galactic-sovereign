package driven

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/Knoblauchpilze/backend-toolkit/pkg/db"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/ports/driven"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIT_ResourceRepository_List(t *testing.T) {
	repo, conn := newTestResourceRepository(t)
	defer conn.Close(context.Background())
	r1 := insertTestResource(t, conn)
	r2 := insertTestResource(t, conn)

	actual, err := repo.List(context.Background())

	assert.Nil(t, err)
	assert.GreaterOrEqual(t, len(actual), 2)
	assert.Contains(t, actual, r1)
	assert.Contains(t, actual, r2)
}

func newTestResourceRepository(t *testing.T) (driven.ForListingResources, db.Connection) {
	conn := newTestConnection(t)
	return NewResourceRepository(conn), conn
}

func insertTestResource(t *testing.T, conn db.Connection) models.Resource {
	someTime := time.Date(2024, 11, 30, 9, 23, 53, 0, time.UTC)

	resource := models.Resource{
		Id:              uuid.New(),
		Name:            fmt.Sprintf("my-resource-%s", uuid.NewString()),
		StartAmount:     456,
		StartProduction: 321,
		StartStorage:    778899,
		CreatedAt:       someTime,
	}

	sqlQuery := `INSERT INTO resource (id, name, start_amount, start_production, start_storage, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := conn.Exec(
		context.Background(),
		sqlQuery,
		resource.Id,
		resource.Name,
		resource.StartAmount,
		resource.StartProduction,
		resource.StartStorage,
		resource.CreatedAt,
	)
	require.Nil(t, err)

	return resource
}
