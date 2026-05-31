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
	b1 := insertTestBuilding(t, conn, addBuildingCost)
	b2 := insertTestBuilding(t, conn, addBuildingProduction)
	b3 := insertTestBuilding(t, conn, addBuildingCost, addBuildingProduction)

	actual, err := repo.List(context.Background())
	require.NoError(t, err, "Actual err: %v", err)

	assert.GreaterOrEqual(t, len(actual), 2)
	assert.Contains(t, actual, b1)
	assert.Contains(t, actual, b2)
	assert.Contains(t, actual, b3)
}

func newTestBuildingRepository(t *testing.T) (driven.ForListingBuildings, db.Connection) {
	conn := newTestConnection(t)
	return NewBuildingRepository(conn), conn
}

func insertTestBuilding(t *testing.T, conn db.Connection, modifiers ...func(*testing.T, db.Connection, *models.Building)) models.Building {
	t.Helper()

	someTime := time.Date(2024, 11, 30, 9, 12, 03, 0, time.UTC)

	building := models.Building{
		Id:        uuid.New(),
		Name:      fmt.Sprintf("my-building-%s", uuid.NewString()),
		CreatedAt: someTime,
		// This is intentional: the details (e.g. costs, productions, etc.) are returned as empty
		// slices by the adapter
		Costs:       []models.BuildingCost{},
		Productions: []models.BuildingResourceProduction{},
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

	for _, modifier := range modifiers {
		modifier(t, conn, &building)
	}

	return building
}

func addBuildingCost(t *testing.T, conn db.Connection, b *models.Building) {
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
		b.Id,
		cost.Resource,
		cost.Cost,
		cost.Progress,
	)
	require.Nil(t, err)

	b.Costs = append(b.Costs, cost)
}

func addBuildingProduction(t *testing.T, conn db.Connection, b *models.Building) {
	t.Helper()

	production := models.BuildingResourceProduction{
		Resource: metalResourceId,
		Base:     rand.Intn(1748),
		// Progress is stored with 5 decimals in the DB
		Progress: randFloat(5),
	}

	sqlQuery := `INSERT INTO building_resource_production (building, resource, base, progress) VALUES ($1, $2, $3, $4)`
	_, err := conn.Exec(
		context.Background(),
		sqlQuery,
		b.Id,
		production.Resource,
		production.Base,
		production.Progress,
	)
	require.Nil(t, err)

	b.Productions = append(b.Productions, production)
}
