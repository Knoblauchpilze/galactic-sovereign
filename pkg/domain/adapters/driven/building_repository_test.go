package drivenadapters

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/Knoblauchpilze/backend-toolkit/pkg/db"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	domainerrors "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models/errors"
	drivenports "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/ports/driven"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	metalResourceId = uuid.MustParse("b4419b6b-b3bf-4576-aa92-055283addbc8")
)

func TestIT_BuildingRepository_Get(t *testing.T) {
	repo, conn := newTestBuildingRepository(t)

	t.Run("gets a building", func(t *testing.T) {
		building := insertTestBuilding(t, conn)

		actual, err := repo.Get(t.Context(), building.Id)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, building, actual)
	})

	t.Run("returns error when building does not exist", func(t *testing.T) {
		id := uuid.MustParse("00000000-1111-2222-1111-000000000000")
		_, err := repo.Get(t.Context(), id)

		assert.ErrorIs(t, err, domainerrors.ErrNotFound, "Actual err: %v", err)
	})
}

func TestIT_BuildingRepository_List(t *testing.T) {
	repo, conn := newTestBuildingRepository(t)
	b1 := insertTestBuilding(t, conn, addBuildingCost)
	b2 := insertTestBuilding(t, conn, addBuildingProduction)
	b3 := insertTestBuilding(t, conn, addBuildingStorage)

	mineLikeBuilding := insertTestBuilding(t, conn, addBuildingCost, addBuildingProduction)
	hangarLikeBuilding := insertTestBuilding(t, conn, addBuildingCost, addBuildingStorage)
	fullCheckBuilding := insertTestBuilding(t, conn, addBuildingCost, addBuildingProduction, addBuildingStorage)

	actual, err := repo.List(t.Context())
	require.NoError(t, err, "Actual err: %v", err)

	// The additional resources are buildings from the seed data
	assert.Contains(t, actual, b1)
	assert.Contains(t, actual, b2)
	assert.Contains(t, actual, b3)
	assert.Contains(t, actual, mineLikeBuilding)
	assert.Contains(t, actual, hangarLikeBuilding)
	assert.Contains(t, actual, fullCheckBuilding)
}

func newTestBuildingRepository(t *testing.T) (drivenports.ForListingBuildings, db.Connection) {
	t.Helper()
	conn := newTestConnection(t)
	return NewBuildingRepository(conn), conn
}

func insertTestBuilding(
	t *testing.T,
	conn db.Connection,
	modifiers ...func(*testing.T, db.Connection, *models.Building),
) models.Building {
	t.Helper()

	building := models.Building{
		Id:        uuid.New(),
		Name:      fmt.Sprintf("my-building-%s", uuid.NewString()),
		CreatedAt: someTime,
		// This is intentional: the details (e.g. costs, productions, etc.) are returned as empty
		// slices by the adapter
		Costs:       []models.BuildingCost{},
		Productions: []models.BuildingResourceProduction{},
		Storages:    []models.BuildingResourceStorage{},
	}

	sqlQuery := `INSERT INTO building (id, name, created_at) VALUES ($1, $2, $3)`
	_, err := conn.Exec(
		t.Context(),
		sqlQuery,
		building.Id,
		building.Name,
		building.CreatedAt,
	)
	require.NoError(t, err, "Actual err: %v", err)

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
		Progress: randFloat(10, 100, 5),
	}

	sqlQuery := `INSERT INTO building_cost (building, resource, cost, progress)
		VALUES ($1, $2, $3, $4)`
	_, err := conn.Exec(
		t.Context(),
		sqlQuery,
		b.Id,
		cost.Resource,
		cost.Cost,
		cost.Progress,
	)
	require.NoError(t, err, "Actual err: %v", err)

	b.Costs = append(b.Costs, cost)
}

func addBuildingProduction(t *testing.T, conn db.Connection, b *models.Building) {
	t.Helper()

	production := models.BuildingResourceProduction{
		Resource: metalResourceId,
		Base:     rand.Intn(1748),
		// Progress is stored with 5 decimals in the DB
		Progress: randFloat(11, 500, 5),
	}

	sqlQuery := `INSERT INTO building_resource_production (building, resource, base, progress)
		VALUES ($1, $2, $3, $4)`
	_, err := conn.Exec(
		t.Context(),
		sqlQuery,
		b.Id,
		production.Resource,
		production.Base,
		production.Progress,
	)
	require.NoError(t, err, "Actual err: %v", err)

	b.Productions = append(b.Productions, production)
}

func addBuildingStorage(t *testing.T, conn db.Connection, b *models.Building) {
	t.Helper()

	storage := models.BuildingResourceStorage{
		Resource: metalResourceId,
		Base:     rand.Intn(1748),
		// Scale and progress are stored with 5 decimals in the DB
		// The min value is arbitrary but ideally should be represented exactly as
		// a float, otherwise some comparison with assert.Contains might fail.
		Scale:    randFloat(0.0625, 1, 5),
		Progress: randFloat(0.0625, 1, 5),
	}

	sqlQuery := `INSERT INTO building_resource_storage (building, resource, base, scale, progress)
		VALUES ($1, $2, $3, $4, $5)`
	_, err := conn.Exec(
		t.Context(),
		sqlQuery,
		b.Id,
		storage.Resource,
		storage.Base,
		storage.Scale,
		storage.Progress,
	)
	require.NoError(t, err, "Actual err: %v", err)

	b.Storages = append(b.Storages, storage)
}
