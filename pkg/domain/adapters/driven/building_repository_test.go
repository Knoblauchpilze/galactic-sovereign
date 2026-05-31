package driven

import (
	"context"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/Knoblauchpilze/backend-toolkit/pkg/db"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/ports/driven"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	metalResourceId = uuid.MustParse("b4419b6b-b3bf-4576-aa92-055283addbc8")
)

func TestIT_BuildingRepository_List(t *testing.T) {
	repo, conn := newTestBuildingRepository(t)
	defer conn.Close(context.Background())
	b1 := insertTestBuilding(t, conn, true)
	b2 := insertTestBuilding(t, conn, false)

	fmt.Printf("building with costs is %s, name=%s, costs=%d\n", b1.Id, b1.Name, len(b1.BaseCosts))
	fmt.Printf("building without costs is %s, name=%s, costs=%d\n", b2.Id, b2.Name, len(b2.BaseCosts))

	actual, err := repo.List(context.Background())
	require.NoError(t, err, "Actual err: %v", err)

	assert.GreaterOrEqual(t, len(actual), 2)
	assert.Contains(t, actual, b1)
	assert.Contains(t, actual, b2)
}

func newTestBuildingRepository(t *testing.T) (driven.ForListingBuildings, db.Connection) {
	conn := newTestConnection(t)
	return NewBuildingRepository(conn), conn
}

func insertTestBuilding(t *testing.T, conn db.Connection, withCost bool) models.Building {
	t.Helper()

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

	if withCost {
		cost := insertTestBuildingCost(t, conn, building.Id)
		building.BaseCosts = append(building.BaseCosts, cost)
	} else {
		// This is intentional: the base costs are returned as an empty slice by the adapter
		building.BaseCosts = []models.BuildingCost{}
	}

	return building
}

func insertTestBuildingCost(t *testing.T, conn db.Connection, buildingId uuid.UUID) models.BuildingCost {
	t.Helper()

	cost := models.BuildingCost{
		Resource: metalResourceId,
		Cost:     rand.Intn(897),
		// Progress is stored with 5 decimals in the DB
		Progress: randFloat(5),
	}

	sqlQuery := `INSERT INTO building_cost (building, resource, cost, progress) VALUES ($1, $2, $3, $4)`
	_, err := conn.Exec(
		context.Background(),
		sqlQuery,
		buildingId,
		cost.Resource,
		cost.Cost,
		cost.Progress,
	)
	require.Nil(t, err)

	return cost
}
