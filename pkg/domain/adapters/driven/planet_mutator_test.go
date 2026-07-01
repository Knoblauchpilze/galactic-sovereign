package drivenadapters

import (
	"testing"
	"time"

	"github.com/Knoblauchpilze/backend-toolkit/pkg/db"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	domainerrors "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models/errors"
	drivenports "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/ports/driven"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	yetAnotherTime = time.Date(2026, time.June, 30, 8, 43, 1, 0, time.UTC)
)

func TestIT_PlanetMutator_Mutate(t *testing.T) {
	adapter, conn := newTestPlanetMutator(t)
	planetRepo := NewPlanetRepository(conn)

	t.Run("passes planet to mutator", func(t *testing.T) {
		planet, _, _ := insertTestPlanetForPlayer(t, conn)

		var captured models.Planet
		mutator := func(p *models.Planet, t time.Time) error {
			captured = *p
			p.Version++
			return nil
		}

		returned, err := adapter.Mutate(t.Context(), planet.Id, mutator)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, planet, captured)
		expected := planet
		expected.Version++
		assert.Equal(t, expected, returned)
	})

	t.Run("returns mutated planet", func(t *testing.T) {
		planet, _, _ := insertTestPlanetForPlayer(t, conn)
		require.NotEqual(t, planet.UpdatedAt, yetAnotherTime)

		mutator := generateModifyingMutator(func(p *models.Planet) {
			p.UpdatedAt = yetAnotherTime
			p.Version++
		})

		returned, err := adapter.Mutate(t.Context(), planet.Id, mutator)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, yetAnotherTime, returned.UpdatedAt)
	})

	t.Run("persists mutated planet", func(t *testing.T) {
		planet, _, _ := insertTestPlanetForPlayer(t, conn)
		require.NotEqual(t, planet.UpdatedAt, yetAnotherTime)

		mutator := generateModifyingMutator(func(p *models.Planet) {
			p.UpdatedAt = yetAnotherTime
			p.Version++
		})

		_, err := adapter.Mutate(t.Context(), planet.Id, mutator)
		require.NoError(t, err, "Actual err: %v", err)

		actual, err := planetRepo.Get(t.Context(), planet.Id)
		require.NoError(t, err, "Actual err: %v", err)
		expected := planet
		expected.UpdatedAt = yetAnotherTime
		expected.Version++
		assert.Equal(t, expected, actual)
	})

	t.Run("persists mutated planet resources", func(t *testing.T) {
		planet, _, _ := insertTestPlanetForPlayer(t, conn, addPlanetResource)
		require.NotEqual(t, 5874, planet.Resources[0].Amount)

		mutator := generateModifyingMutator(func(p *models.Planet) {
			p.Resources[0].Amount = 5874
			p.Version++
		})

		returned, err := adapter.Mutate(t.Context(), planet.Id, mutator)
		require.NoError(t, err, "Actual err: %v", err)

		expected := []models.PlanetResource{
			{Resource: crystalResourceId, Amount: 5874},
		}
		assert.Equal(t, expected, returned.Resources)
		assertPlanetResourceAmount(t, conn, planet.Id, crystalResourceId, 5874)
	})

	t.Run("persists mutated planet productions", func(t *testing.T) {
		planet, _, _ := insertTestPlanetForPlayer(t, conn, addPlanetProduction)
		require.NotEqual(t, 39841, planet.Productions[0].Production)

		mutator := generateModifyingMutator(func(p *models.Planet) {
			p.Productions[0].Production = 39841
			p.Version++
		})

		returned, err := adapter.Mutate(t.Context(), planet.Id, mutator)
		require.NoError(t, err, "Actual err: %v", err)

		expected := []models.PlanetResourceProduction{
			{Resource: metalResourceId, Production: 39841},
		}
		assert.Equal(t, expected, returned.Productions)
		assertPlanetResourceProduction(t, conn, planet.Id, metalResourceId, nil, 39841)
	})

	t.Run("persists mutated planet production for building", func(t *testing.T) {
		planet, _, _ := insertTestPlanetForPlayer(t, conn, addPlanetProductionForBuilding)
		require.NotEqual(t, 1235, planet.Productions[0].Production)

		mutator := generateModifyingMutator(func(p *models.Planet) {
			p.Productions[0].Production = 1235
			p.Version++
		})

		returned, err := adapter.Mutate(t.Context(), planet.Id, mutator)
		require.NoError(t, err, "Actual err: %v", err)

		expected := []models.PlanetResourceProduction{
			{Resource: metalResourceId, Building: &crystalMineId, Production: 1235},
		}
		assert.Equal(t, expected, returned.Productions)
		assertPlanetResourceProduction(
			t, conn, planet.Id, metalResourceId, &crystalMineId, 1235,
		)
	})

	t.Run("persists additional planet productions for building", func(t *testing.T) {
		planet, _, _ := insertTestPlanetForPlayer(
			t, conn, addPlanetProduction, addPlanetProductionForBuilding,
		)
		require.NotEqual(t, metalMineId, planet.Productions[1].Building)
		initial0 := planet.Productions[0].Production
		initial1 := planet.Productions[1].Production

		mutator := generateModifyingMutator(func(p *models.Planet) {
			prod := models.PlanetResourceProduction{
				Resource:   crystalResourceId,
				Building:   &metalMineId,
				Production: 354789,
			}
			p.Productions = append(p.Productions, prod)
			p.Version++
		})

		returned, err := adapter.Mutate(t.Context(), planet.Id, mutator)
		require.NoError(t, err, "Actual err: %v", err)

		expected := []models.PlanetResourceProduction{
			{Resource: metalResourceId, Production: initial0},
			{Resource: metalResourceId, Building: &crystalMineId, Production: initial1},
			{Resource: crystalResourceId, Building: &metalMineId, Production: 354789},
		}
		assert.Equal(t, expected, returned.Productions)
		assertPlanetResourceProduction(t, conn, planet.Id, metalResourceId, nil, initial0)
		assertPlanetResourceProduction(t, conn, planet.Id, metalResourceId, &crystalMineId, initial1)
		assertPlanetResourceProduction(t, conn, planet.Id, crystalResourceId, &metalMineId, 354789)
	})

	t.Run("persists mutated planet storages", func(t *testing.T) {
		planet, _, _ := insertTestPlanetForPlayer(t, conn, addPlanetStorage)
		require.NotEqual(t, 4598, planet.Storages[0].Storage)

		mutator := generateModifyingMutator(func(p *models.Planet) {
			p.Storages[0].Storage = 4598
			p.Version++
		})

		returned, err := adapter.Mutate(t.Context(), planet.Id, mutator)
		require.NoError(t, err, "Actual err: %v", err)

		expected := []models.PlanetResourceStorage{
			{Resource: crystalResourceId, Storage: 4598},
		}
		assert.Equal(t, expected, returned.Storages)
		assertPlanetResourceStorage(t, conn, planet.Id, crystalResourceId, 4598)
	})

	t.Run("persists mutated planet buildings", func(t *testing.T) {
		planet, _, _ := insertTestPlanetForPlayer(t, conn, addPlanetBuilding)
		require.NotEqual(t, 6, planet.Buildings[0].Level)

		mutator := generateModifyingMutator(func(p *models.Planet) {
			p.Buildings[0].Level = 6
			p.Version++
		})

		returned, err := adapter.Mutate(t.Context(), planet.Id, mutator)
		require.NoError(t, err, "Actual err: %v", err)

		expected := []models.PlanetBuilding{
			{Building: metalStorageId, Level: 6},
		}
		assert.Equal(t, expected, returned.Buildings)
		assertPlanetBuildingLevel(t, conn, planet.Id, metalStorageId, 6)
	})

	t.Run("persists mutated planet with action", func(t *testing.T) {
		planet, _, _ := insertTestPlanetForPlayer(t, conn)
		require.Nil(t, planet.BuildingAction)

		action := models.BuildingAction{
			Id:           uuid.New(),
			Planet:       planet.Id,
			Building:     metalMineId,
			DesiredLevel: 3,
			CreatedAt:    someTime,
			CompletedAt:  someTime.Add(1 * time.Hour),
			Costs:        []models.BuildingActionCost{},
			Storages:     []models.BuildingActionResourceStorage{},
			Productions:  []models.BuildingActionResourceProduction{},
		}

		mutator := generateModifyingMutator(func(p *models.Planet) {
			p.BuildingAction = &action
			p.Version++
		})

		returned, err := adapter.Mutate(t.Context(), planet.Id, mutator)
		require.NoError(t, err, "Actual err: %v", err)

		actual, err := planetRepo.Get(t.Context(), planet.Id)
		require.NoError(t, err, "Actual err: %v", err)
		assert.Equal(t, returned, actual)
		require.NotNil(t, returned.BuildingAction)
		assert.Equal(t, action, *returned.BuildingAction)
		require.NotNil(t, actual.BuildingAction)
		assert.Equal(t, action, *actual.BuildingAction)
	})

	t.Run("persists mutated planet with action and costs", func(t *testing.T) {
		planet, _, _ := insertTestPlanetForPlayer(t, conn)
		require.Nil(t, planet.BuildingAction)

		action := models.BuildingAction{
			Id:           uuid.New(),
			Planet:       planet.Id,
			Building:     metalMineId,
			DesiredLevel: 3,
			CreatedAt:    someTime,
			CompletedAt:  someTime.Add(1 * time.Hour),
			Costs: []models.BuildingActionCost{
				{
					Resource: metalResourceId,
					Amount:   36,
				},
				{
					Resource: crystalResourceId,
					Amount:   798,
				},
			},
			Storages:    []models.BuildingActionResourceStorage{},
			Productions: []models.BuildingActionResourceProduction{},
		}

		mutator := generateModifyingMutator(func(p *models.Planet) {
			p.BuildingAction = &action
			p.Version++
		})

		returned, err := adapter.Mutate(t.Context(), planet.Id, mutator)
		require.NoError(t, err, "Actual err: %v", err)

		actual, err := planetRepo.Get(t.Context(), planet.Id)
		require.NoError(t, err, "Actual err: %v", err)
		assert.Equal(t, returned, actual)
		require.NotNil(t, returned.BuildingAction)
		assert.Equal(t, action, *returned.BuildingAction)
		require.NotNil(t, actual.BuildingAction)
		assert.Equal(t, action, *actual.BuildingAction)
	})

	t.Run("persists mutated planet with action and storages", func(t *testing.T) {
		planet, _, _ := insertTestPlanetForPlayer(t, conn)
		require.Nil(t, planet.BuildingAction)

		action := models.BuildingAction{
			Id:           uuid.New(),
			Planet:       planet.Id,
			Building:     metalMineId,
			DesiredLevel: 3,
			CreatedAt:    someTime,
			CompletedAt:  someTime.Add(1 * time.Hour),
			Costs:        []models.BuildingActionCost{},
			Storages: []models.BuildingActionResourceStorage{
				{
					Resource: metalResourceId,
					Storage:  321417,
				},
				{
					Resource: crystalResourceId,
					Storage:  65478,
				},
			},
			Productions: []models.BuildingActionResourceProduction{},
		}

		mutator := generateModifyingMutator(func(p *models.Planet) {
			p.BuildingAction = &action
			p.Version++
		})

		returned, err := adapter.Mutate(t.Context(), planet.Id, mutator)
		require.NoError(t, err, "Actual err: %v", err)

		actual, err := planetRepo.Get(t.Context(), planet.Id)
		require.NoError(t, err, "Actual err: %v", err)
		assert.Equal(t, returned, actual)
		require.NotNil(t, returned.BuildingAction)
		assert.Equal(t, action, *returned.BuildingAction)
		require.NotNil(t, actual.BuildingAction)
		assert.Equal(t, action, *actual.BuildingAction)
	})

	t.Run("persists mutated planet with action and productions", func(t *testing.T) {
		planet, _, _ := insertTestPlanetForPlayer(t, conn)
		require.Nil(t, planet.BuildingAction)

		action := models.BuildingAction{
			Id:           uuid.New(),
			Planet:       planet.Id,
			Building:     metalMineId,
			DesiredLevel: 3,
			CreatedAt:    someTime,
			CompletedAt:  someTime.Add(1 * time.Hour),
			Costs:        []models.BuildingActionCost{},
			Storages:     []models.BuildingActionResourceStorage{},
			Productions: []models.BuildingActionResourceProduction{
				{
					Resource:   metalResourceId,
					Production: 147,
				},
				{
					Resource:   crystalResourceId,
					Production: 3254,
				},
			},
		}

		mutator := generateModifyingMutator(func(p *models.Planet) {
			p.BuildingAction = &action
			p.Version++
		})

		returned, err := adapter.Mutate(t.Context(), planet.Id, mutator)
		require.NoError(t, err, "Actual err: %v", err)

		actual, err := planetRepo.Get(t.Context(), planet.Id)
		require.NoError(t, err, "Actual err: %v", err)
		assert.Equal(t, returned, actual)
		require.NotNil(t, returned.BuildingAction)
		assert.Equal(t, action, *returned.BuildingAction)
		require.NotNil(t, actual.BuildingAction)
		assert.Equal(t, action, *actual.BuildingAction)
	})

	t.Run("persists mutated planet with updated completion time", func(t *testing.T) {
		planet, _, _ := insertTestPlanetForPlayer(t, conn)
		action := insertTestBuildingActionForPlanet(t, conn, planet.Id)
		require.NotEqual(t, yetAnotherTime, action.CompletedAt)

		mutator := generateModifyingMutator(func(p *models.Planet) {
			p.BuildingAction.CompletedAt = yetAnotherTime
			p.Version++
		})

		returned, err := adapter.Mutate(t.Context(), planet.Id, mutator)
		require.NoError(t, err, "Actual err: %v", err)

		actual, err := planetRepo.Get(t.Context(), planet.Id)
		require.NoError(t, err, "Actual err: %v", err)
		assert.Equal(t, returned, actual)
		require.NotNil(t, returned.BuildingAction)
		assert.Equal(t, yetAnotherTime, returned.BuildingAction.CompletedAt)
		require.NotNil(t, actual.BuildingAction)
		assert.Equal(t, yetAnotherTime, actual.BuildingAction.CompletedAt)
	})

	t.Run("persists mutated planet without update to existing action costs", func(t *testing.T) {
		planet, _, _ := insertTestPlanetForPlayer(t, conn)
		action := insertTestBuildingActionForPlanet(t, conn, planet.Id, addBuildingActionCost)
		require.Equal(t, metalResourceId, action.Costs[0].Resource)
		require.NotEqual(t, 32, action.Costs[0].Amount)

		costs := []models.BuildingActionCost{{Resource: metalResourceId, Amount: 32}}

		mutator := generateModifyingMutator(func(p *models.Planet) {
			p.BuildingAction.Costs = costs
			p.Version++
		})

		returned, err := adapter.Mutate(t.Context(), planet.Id, mutator)
		require.NoError(t, err, "Actual err: %v", err)

		actual, err := planetRepo.Get(t.Context(), planet.Id)
		require.NoError(t, err, "Actual err: %v", err)
		require.NotNil(t, returned.BuildingAction)
		assert.Equal(t, costs, returned.BuildingAction.Costs)
		require.NotNil(t, actual.BuildingAction)
		assert.Equal(t, action.Costs, actual.BuildingAction.Costs)
	})

	t.Run("persists mutated planet with updated action and costs", func(t *testing.T) {
		planet, _, _ := insertTestPlanetForPlayer(t, conn)
		insertTestBuildingActionForPlanet(t, conn, planet.Id)

		costs := []models.BuildingActionCost{
			{
				Resource: metalResourceId,
				Amount:   45698,
			},
			{
				Resource: crystalResourceId,
				Amount:   120305,
			},
		}

		mutator := generateModifyingMutator(func(p *models.Planet) {
			p.BuildingAction.Costs = costs
			p.Version++
		})

		returned, err := adapter.Mutate(t.Context(), planet.Id, mutator)
		require.NoError(t, err, "Actual err: %v", err)

		actual, err := planetRepo.Get(t.Context(), planet.Id)
		require.NoError(t, err, "Actual err: %v", err)
		assert.Equal(t, returned, actual)
		require.NotNil(t, returned.BuildingAction)
		assert.Equal(t, costs, returned.BuildingAction.Costs)
		require.NotNil(t, actual.BuildingAction)
		assert.Equal(t, costs, actual.BuildingAction.Costs)
	})

	t.Run("persists mutated planet without update to existing action storages", func(t *testing.T) {
		planet, _, _ := insertTestPlanetForPlayer(t, conn)
		action := insertTestBuildingActionForPlanet(t, conn, planet.Id, addBuildingActionStorage)
		require.Equal(t, crystalResourceId, action.Storages[0].Resource)
		require.NotEqual(t, 32, action.Storages[0].Storage)

		storages := []models.BuildingActionResourceStorage{{Resource: crystalResourceId, Storage: 32}}

		mutator := generateModifyingMutator(func(p *models.Planet) {
			p.BuildingAction.Storages = storages
			p.Version++
		})

		returned, err := adapter.Mutate(t.Context(), planet.Id, mutator)
		require.NoError(t, err, "Actual err: %v", err)

		actual, err := planetRepo.Get(t.Context(), planet.Id)
		require.NoError(t, err, "Actual err: %v", err)
		require.NotNil(t, returned.BuildingAction)
		assert.Equal(t, storages, returned.BuildingAction.Storages)
		require.NotNil(t, actual.BuildingAction)
		assert.Equal(t, action.Storages, actual.BuildingAction.Storages)
	})

	t.Run("persists mutated planet with updated action and storages", func(t *testing.T) {
		planet, _, _ := insertTestPlanetForPlayer(t, conn)
		insertTestBuildingActionForPlanet(t, conn, planet.Id)

		storages := []models.BuildingActionResourceStorage{
			{
				Resource: metalResourceId,
				Storage:  123456,
			},
			{
				Resource: crystalResourceId,
				Storage:  987654,
			},
		}

		mutator := generateModifyingMutator(func(p *models.Planet) {
			p.BuildingAction.Storages = storages
			p.Version++
		})

		returned, err := adapter.Mutate(t.Context(), planet.Id, mutator)
		require.NoError(t, err, "Actual err: %v", err)

		actual, err := planetRepo.Get(t.Context(), planet.Id)
		require.NoError(t, err, "Actual err: %v", err)
		assert.Equal(t, returned, actual)
		require.NotNil(t, returned.BuildingAction)
		assert.Equal(t, storages, returned.BuildingAction.Storages)
		require.NotNil(t, actual.BuildingAction)
		assert.Equal(t, storages, actual.BuildingAction.Storages)
	})

	t.Run("persists mutated planet without update to existing action productions", func(t *testing.T) {
		planet, _, _ := insertTestPlanetForPlayer(t, conn)
		action := insertTestBuildingActionForPlanet(t, conn, planet.Id, addBuildingActionProduction)
		require.Equal(t, crystalResourceId, action.Productions[0].Resource)
		require.NotEqual(t, 32, action.Productions[0].Production)

		productions := []models.BuildingActionResourceProduction{{Resource: crystalResourceId, Production: 32}}

		mutator := generateModifyingMutator(func(p *models.Planet) {
			p.BuildingAction.Productions = productions
			p.Version++
		})

		returned, err := adapter.Mutate(t.Context(), planet.Id, mutator)
		require.NoError(t, err, "Actual err: %v", err)

		actual, err := planetRepo.Get(t.Context(), planet.Id)
		require.NoError(t, err, "Actual err: %v", err)
		require.NotNil(t, returned.BuildingAction)
		assert.Equal(t, productions, returned.BuildingAction.Productions)
		require.NotNil(t, actual.BuildingAction)
		assert.Equal(t, action.Productions, actual.BuildingAction.Productions)
	})

	t.Run("persists mutated planet with updated action and productions", func(t *testing.T) {
		planet, _, _ := insertTestPlanetForPlayer(t, conn)
		insertTestBuildingActionForPlanet(t, conn, planet.Id)

		productions := []models.BuildingActionResourceProduction{
			{
				Resource:   metalResourceId,
				Production: 45698,
			},
			{
				Resource:   crystalResourceId,
				Production: 120305,
			},
		}

		mutator := generateModifyingMutator(func(p *models.Planet) {
			p.BuildingAction.Productions = productions
			p.Version++
		})

		returned, err := adapter.Mutate(t.Context(), planet.Id, mutator)
		require.NoError(t, err, "Actual err: %v", err)

		actual, err := planetRepo.Get(t.Context(), planet.Id)
		require.NoError(t, err, "Actual err: %v", err)
		assert.Equal(t, returned, actual)
		require.NotNil(t, returned.BuildingAction)
		assert.Equal(t, productions, returned.BuildingAction.Productions)
		require.NotNil(t, actual.BuildingAction)
		assert.Equal(t, productions, actual.BuildingAction.Productions)
	})

	t.Run("persists mutated planet with deleted action", func(t *testing.T) {
		action, planet := insertTestBuildingAction(t, conn)

		mutator := generateModifyingMutator(func(p *models.Planet) {
			p.BuildingAction = nil
			p.Version++
		})

		returned, err := adapter.Mutate(t.Context(), planet.Id, mutator)
		require.NoError(t, err, "Actual err: %v", err)

		assertBuildingActionDoesNotExist(t, conn, action.Id)
		assert.Nil(t, returned.BuildingAction)
	})

	t.Run("persists mutated planet with deleted action with costs", func(t *testing.T) {
		action, planet := insertTestBuildingAction(t, conn, addBuildingActionCost)

		mutator := generateModifyingMutator(func(p *models.Planet) {
			p.BuildingAction = nil
			p.Version++
		})

		returned, err := adapter.Mutate(t.Context(), planet.Id, mutator)
		require.NoError(t, err, "Actual err: %v", err)

		assertBuildingActionDoesNotExist(t, conn, action.Id)
		assertBuildingActionCostDoesNotExist(t, conn, action.Id)
		assert.Nil(t, returned.BuildingAction)
	})

	t.Run("persists mutated planet with deleted action with storages", func(t *testing.T) {
		action, planet := insertTestBuildingAction(t, conn, addBuildingActionStorage)

		mutator := generateModifyingMutator(func(p *models.Planet) {
			p.BuildingAction = nil
			p.Version++
		})

		returned, err := adapter.Mutate(t.Context(), planet.Id, mutator)
		require.NoError(t, err, "Actual err: %v", err)

		assertBuildingActionDoesNotExist(t, conn, action.Id)
		assertBuildingActionStorageDoesNotExist(t, conn, action.Id)
		assert.Nil(t, returned.BuildingAction)
	})

	t.Run("persists mutated planet with deleted action with productions", func(t *testing.T) {
		action, planet := insertTestBuildingAction(t, conn, addBuildingActionProduction)

		mutator := generateModifyingMutator(func(p *models.Planet) {
			p.BuildingAction = nil
			p.Version++
		})

		returned, err := adapter.Mutate(t.Context(), planet.Id, mutator)
		require.NoError(t, err, "Actual err: %v", err)

		assertBuildingActionDoesNotExist(t, conn, action.Id)
		assertBuildingActionProductionDoesNotExist(t, conn, action.Id)
		assert.Nil(t, returned.BuildingAction)
	})

	t.Run("returns error when mutator does not update version", func(t *testing.T) {
		planet, _, _ := insertTestPlanetForPlayer(t, conn)

		mutator := generateModifyingMutator(func(*models.Planet) {})

		_, err := adapter.Mutate(t.Context(), planet.Id, mutator)

		assert.ErrorIs(t, domainerrors.ErrOptimisticLocking, err)
	})

	t.Run("returns error when new building is added", func(t *testing.T) {
		planet, _, _ := insertTestPlanetForPlayer(t, conn)

		mutator := generateModifyingMutator(func(p *models.Planet) {
			p.Buildings = []models.PlanetBuilding{{Building: metalMineId, Level: 5}}
			p.Version++
		})

		_, err := adapter.Mutate(t.Context(), planet.Id, mutator)

		assert.ErrorIs(t, domainerrors.ErrBuildingNotFound, err)
	})

	t.Run("returns error when new storage is added", func(t *testing.T) {
		planet, _, _ := insertTestPlanetForPlayer(t, conn)

		mutator := generateModifyingMutator(func(p *models.Planet) {
			p.Storages = []models.PlanetResourceStorage{{Resource: crystalResourceId, Storage: 5478}}
			p.Version++
		})

		_, err := adapter.Mutate(t.Context(), planet.Id, mutator)

		assert.ErrorIs(t, domainerrors.ErrResourceNotFound, err)
	})
}

func newTestPlanetMutator(t *testing.T) (*PlanetMutator, db.Connection) {
	t.Helper()
	conn := newTestConnection(t)
	return NewPlanetMutator(conn), conn
}

func generateModifyingMutator(modifier func(p *models.Planet)) drivenports.PlanetMutator {
	return func(p *models.Planet, t time.Time) error {
		modifier(p)
		return nil
	}
}

func assertPlanetBuildingLevel(t *testing.T, conn db.Connection, planet uuid.UUID, building uuid.UUID, level int) {
	t.Helper()

	sqlQuery := `SELECT level FROM planet_building WHERE planet = $1 AND building = $2`
	value, err := db.QueryOne[int](t.Context(), conn, sqlQuery, planet, building)
	require.NoError(t, err, "Actual err: %v", err)
	require.Equal(t, level, value)
}
