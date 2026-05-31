package driven

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/Knoblauchpilze/backend-toolkit/pkg/db"
	eassert "github.com/Knoblauchpilze/easy-assert/assert"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/ports/driven"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIT_BuildingRepository_List(t *testing.T) {
	repo, conn := newTestBuildingRepository(t)
	defer conn.Close(context.Background())
	b1 := insertTestBuilding(t, conn)
	b2 := insertTestBuilding(t, conn)

	actual, err := repo.List(context.Background())

	assert.Nil(t, err)
	assert.GreaterOrEqual(t, len(actual), 2)
	assert.True(t, eassert.ContainsIgnoringFields(actual, b1))
	assert.True(t, eassert.ContainsIgnoringFields(actual, b2))
}

func newTestBuildingRepository(t *testing.T) (driven.ForManagingBuildings, db.Connection) {
	conn := newTestConnection(t)
	return NewBuildingRepository(conn), conn
}

func insertTestBuilding(t *testing.T, conn db.Connection) models.Building {
	someTime := time.Date(2024, 11, 30, 9, 12, 03, 0, time.UTC)

	building := models.Building{
		Id:        uuid.New(),
		Name:      fmt.Sprintf("my-building-%s", uuid.NewString()),
		CreatedAt: someTime,
	}

	sqlQuery := `INSERT INTO building (id, name, created_at) VALUES ($1, $2, $3)`
	_, err := conn.Exec(
		context.Background(),
		sqlQuery,
		building.Id,
		building.Name,
		building.CreatedAt,
	)
	require.Nil(t, err)

	return building
}
