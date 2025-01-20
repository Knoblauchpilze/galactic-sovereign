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

func TestIT_BuildingResourceProductionRepository_List(t *testing.T) {
	repo, conn, tx := newTestBuildingResourceProductionRepositoryAndTransaction(t)
	defer conn.Close(context.Background())
	defer tx.Close(context.Background())
	brp1, b, _ := insertTestBuildingResourceProductionForBuilding(t, conn)
	brp2, _ := insertTestBuildingResourceProduction(t, conn, b.Id)

	actual, err := repo.ListForBuilding(context.Background(), tx, b.Id)

	assert.Nil(t, err)
	assert.GreaterOrEqual(t, len(actual), 2)
	assert.Contains(t, actual, brp1)
	assert.Contains(t, actual, brp2)
}

func newTestBuildingResourceProductionRepositoryAndTransaction(t *testing.T) (BuildingResourceProductionRepository, db.Connection, db.Transaction) {
	conn := newTestConnection(t)
	tx, err := conn.BeginTx(context.Background())
	require.Nil(t, err)
	return NewBuildingResourceProductionRepository(), conn, tx
}

func insertTestBuildingResourceProduction(t *testing.T, conn db.Connection, building uuid.UUID) (persistence.BuildingResourceProduction, persistence.Resource) {
	resource := insertTestResource(t, conn)

	prod := persistence.BuildingResourceProduction{
		Building: building,
		Resource: resource.Id,
		Base:     17,
		Progress: 1.6,
	}

	sqlQuery := `INSERT INTO building_resource_production (building, resource, base, progress) VALUES ($1, $2, $3, $4)`
	_, err := conn.Exec(
		context.Background(),
		sqlQuery,
		prod.Building,
		prod.Resource,
		prod.Base,
		prod.Progress,
	)
	require.Nil(t, err)

	return prod, resource
}

func insertTestBuildingResourceProductionForBuilding(t *testing.T, conn db.Connection) (persistence.BuildingResourceProduction, persistence.Building, persistence.Resource) {
	building := insertTestBuilding(t, conn)
	prod, resource := insertTestBuildingResourceProduction(t, conn, building.Id)
	return prod, building, resource
}
