package repositories

import (
	"context"
	"testing"

	"github.com/Knoblauchpilze/backend-toolkit/pkg/db"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/persistence"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIT_BuildingCostRepository_List(t *testing.T) {
	repo, conn, tx := newTestBuildingCostRepositoryAndTransaction(t)
	defer conn.Close(context.Background())
	defer tx.Close(context.Background())
	bc1, b, _ := insertTestBuildingCostForBuilding(t, conn)
	bc2, _ := insertTestBuildingCost(t, conn, b.Id)

	actual, err := repo.ListForBuilding(context.Background(), tx, b.Id)

	assert.Nil(t, err)
	assert.GreaterOrEqual(t, len(actual), 2)
	assert.Contains(t, actual, bc1)
	assert.Contains(t, actual, bc2)
}

func newTestBuildingCostRepositoryAndTransaction(t *testing.T) (BuildingCostRepository, db.Connection, db.Transaction) {
	conn := newTestConnection(t)
	tx, err := conn.BeginTx(context.Background())
	require.Nil(t, err)
	return NewBuildingCostRepository(), conn, tx
}

func insertTestBuildingCost(t *testing.T, conn db.Connection, building uuid.UUID) (persistence.BuildingCost, persistence.Resource) {
	resource := insertTestResource(t, conn)

	cost := persistence.BuildingCost{
		Building: building,
		Resource: resource.Id,
		Cost:     41,
		Progress: 1.6,
	}

	sqlQuery := `INSERT INTO building_cost (building, resource, cost, progress) VALUES ($1, $2, $3, $4)`
	_, err := conn.Exec(
		context.Background(),
		sqlQuery,
		cost.Building,
		cost.Resource,
		cost.Cost,
		cost.Progress,
	)
	require.Nil(t, err)

	return cost, resource
}

func insertTestBuildingCostForBuilding(t *testing.T, conn db.Connection) (persistence.BuildingCost, persistence.Building, persistence.Resource) {
	building := insertTestBuilding(t, conn)
	cost, resource := insertTestBuildingCost(t, conn, building.Id)
	return cost, building, resource
}
