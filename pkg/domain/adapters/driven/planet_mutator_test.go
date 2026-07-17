package drivenadapters

import (
	"context"
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

	t.Run("passes planet to mutator", func(t *testing.T) {
		planet, _, _ := insertTestPlanetForPlayer(t, conn)

		var captured models.Planet
		mutator := func(p *models.Planet) (bool, error) {
			captured = *p
			p.Version++
			return false, nil
		}

		returned, err := adapter.Mutate(t.Context(), planet.Id, mutator)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, planet, captured)
		expected := planet
		expected.Version++
		assert.False(t, returned.Deleted)
		assert.Equal(t, expected, returned.Planet)
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

		assert.False(t, returned.Deleted)
		assert.Equal(t, yetAnotherTime, returned.Planet.UpdatedAt)
	})

	t.Run("persists mutated planet", func(t *testing.T) {
		planet, _, _ := insertTestPlanetForPlayer(t, conn)
		require.NotEqual(t, planet.UpdatedAt, yetAnotherTime)
		require.NotEqual(t, 326, planet.Fields)

		mutator := generateModifyingMutator(func(p *models.Planet) {
			p.Fields = 326
			p.UpdatedAt = yetAnotherTime
			p.Version++
		})

		_, err := adapter.Mutate(t.Context(), planet.Id, mutator)
		require.NoError(t, err, "Actual err: %v", err)

		actual := loadPlanetFromDb(t, conn, planet.Id)
		expected := planet
		expected.Fields = 326
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

		assert.False(t, returned.Deleted)
		expected := []models.PlanetResource{
			{Resource: crystalResourceId, Amount: 5874},
		}
		assert.Equal(t, expected, returned.Planet.Resources)
		assertPlanetResourceAmount(t, conn, planet.Id, crystalResourceId, 5874)
	})

	t.Run("does not delete existing planet resource", func(t *testing.T) {
		planet, _, _ := insertTestPlanetForPlayer(t, conn, addPlanetResource)
		resource := planet.Resources[0].Resource
		amount := planet.Resources[0].Amount

		mutator := generateModifyingMutator(func(p *models.Planet) {
			p.Resources = []models.PlanetResource{}
			p.Version++
		})

		returned, err := adapter.Mutate(t.Context(), planet.Id, mutator)
		require.NoError(t, err, "Actual err: %v", err)

		assert.False(t, returned.Deleted)
		expected := []models.PlanetResource{
			{Resource: resource, Amount: amount},
		}
		assert.Equal(t, expected, returned.Planet.Resources)
		assertPlanetResourceAmount(t, conn, planet.Id, resource, amount)
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

		assert.False(t, returned.Deleted)
		expected := []models.PlanetResourceProduction{
			{Resource: metalResourceId, Production: 39841},
		}
		assert.Equal(t, expected, returned.Planet.Productions)
		assertPlanetResourceProduction(t, conn, planet.Id, metalResourceId, nil, 39841)
	})

	t.Run("persists additional planet productions", func(t *testing.T) {
		planet, _, _ := insertTestPlanetForPlayer(
			t, conn, addPlanetProduction, addPlanetProductionForBuilding,
		)
		require.NotEqual(t, crystalResourceId, planet.Productions[0].Building)
		initial0 := planet.Productions[0].Production
		initial1 := planet.Productions[1].Production

		mutator := generateModifyingMutator(func(p *models.Planet) {
			prod := models.PlanetResourceProduction{
				Resource:   crystalResourceId,
				Production: 354789,
			}
			p.Productions = append(p.Productions, prod)
			p.Version++
		})

		returned, err := adapter.Mutate(t.Context(), planet.Id, mutator)
		require.NoError(t, err, "Actual err: %v", err)

		assert.False(t, returned.Deleted)
		expected := []models.PlanetResourceProduction{
			{Resource: metalResourceId, Production: initial0},
			{Resource: metalResourceId, Building: &crystalMineId, Production: initial1},
			{Resource: crystalResourceId, Building: nil, Production: 354789},
		}
		assert.Equal(t, expected, returned.Planet.Productions)
		assertPlanetResourceProduction(t, conn, planet.Id, metalResourceId, nil, initial0)
		assertPlanetResourceProduction(t, conn, planet.Id, metalResourceId, &crystalMineId, initial1)
		assertPlanetResourceProduction(t, conn, planet.Id, crystalResourceId, nil, 354789)
	})

	t.Run("persists deleted planet production", func(t *testing.T) {
		planet, _, _ := insertTestPlanetForPlayer(t, conn, addPlanetProduction)

		mutator := generateModifyingMutator(func(p *models.Planet) {
			p.Productions = []models.PlanetResourceProduction{}
			p.Version++
		})

		returned, err := adapter.Mutate(t.Context(), planet.Id, mutator)
		require.NoError(t, err, "Actual err: %v", err)

		assert.False(t, returned.Deleted)
		assert.Equal(t, []models.PlanetResourceProduction{}, returned.Planet.Productions)
		assertPlanetProductionDoesNotExist(t, conn, planet.Id)
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

		assert.False(t, returned.Deleted)
		expected := []models.PlanetResourceProduction{
			{Resource: metalResourceId, Building: &crystalMineId, Production: 1235},
		}
		assert.Equal(t, expected, returned.Planet.Productions)
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

		assert.False(t, returned.Deleted)
		expected := []models.PlanetResourceProduction{
			{Resource: metalResourceId, Production: initial0},
			{Resource: metalResourceId, Building: &crystalMineId, Production: initial1},
			{Resource: crystalResourceId, Building: &metalMineId, Production: 354789},
		}
		assert.Equal(t, expected, returned.Planet.Productions)
		assertPlanetResourceProduction(t, conn, planet.Id, metalResourceId, nil, initial0)
		assertPlanetResourceProduction(t, conn, planet.Id, metalResourceId, &crystalMineId, initial1)
		assertPlanetResourceProduction(t, conn, planet.Id, crystalResourceId, &metalMineId, 354789)
	})

	t.Run("persists deleted planet production for building", func(t *testing.T) {
		planet, _, _ := insertTestPlanetForPlayer(t, conn, addPlanetProductionForBuilding)

		mutator := generateModifyingMutator(func(p *models.Planet) {
			p.Productions = []models.PlanetResourceProduction{}
			p.Version++
		})

		returned, err := adapter.Mutate(t.Context(), planet.Id, mutator)
		require.NoError(t, err, "Actual err: %v", err)

		assert.False(t, returned.Deleted)
		assert.Equal(t, []models.PlanetResourceProduction{}, returned.Planet.Productions)
		assertPlanetProductionDoesNotExist(t, conn, planet.Id)
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

		assert.False(t, returned.Deleted)
		expected := []models.PlanetResourceStorage{
			{Resource: crystalResourceId, Storage: 4598},
		}
		assert.Equal(t, expected, returned.Planet.Storages)
		assertPlanetResourceStorage(t, conn, planet.Id, crystalResourceId, 4598)
	})

	t.Run("does not delete existing planet storage", func(t *testing.T) {
		planet, _, _ := insertTestPlanetForPlayer(t, conn, addPlanetStorage)
		resource := planet.Storages[0].Resource
		storage := planet.Storages[0].Storage

		mutator := generateModifyingMutator(func(p *models.Planet) {
			p.Storages = []models.PlanetResourceStorage{}
			p.Version++
		})

		returned, err := adapter.Mutate(t.Context(), planet.Id, mutator)
		require.NoError(t, err, "Actual err: %v", err)

		assert.False(t, returned.Deleted)
		expected := []models.PlanetResourceStorage{
			{Resource: resource, Storage: storage},
		}
		assert.Equal(t, expected, returned.Planet.Storages)
		assertPlanetResourceStorage(t, conn, planet.Id, resource, storage)
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

		assert.False(t, returned.Deleted)
		expected := []models.PlanetBuilding{
			{Building: metalStorageId, Level: 6},
		}
		assert.Equal(t, expected, returned.Planet.Buildings)
		assertPlanetBuildingLevel(t, conn, planet.Id, metalStorageId, 6)
	})

	t.Run("does not delete existing planet building", func(t *testing.T) {
		planet, _, _ := insertTestPlanetForPlayer(t, conn, addPlanetBuilding)
		building := planet.Buildings[0].Building
		level := planet.Buildings[0].Level

		mutator := generateModifyingMutator(func(p *models.Planet) {
			p.Buildings = []models.PlanetBuilding{}
			p.Version++
		})

		returned, err := adapter.Mutate(t.Context(), planet.Id, mutator)
		require.NoError(t, err, "Actual err: %v", err)

		assert.False(t, returned.Deleted)
		expected := []models.PlanetBuilding{
			{Building: building, Level: level},
		}
		assert.Equal(t, expected, returned.Planet.Buildings)
		assertPlanetBuildingLevel(t, conn, planet.Id, building, level)
	})

	t.Run("persists mutated planet with action", func(t *testing.T) {
		planet, _, _ := insertTestPlanetForPlayer(t, conn)
		require.Nil(t, planet.BuildingAction)

		action := models.BuildingAction{
			Id:           uuid.New(),
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

		assert.False(t, returned.Deleted)
		actual := loadPlanetFromDb(t, conn, planet.Id)
		assert.Equal(t, returned.Planet, actual)
		require.NotNil(t, returned.Planet.BuildingAction)
		assert.Equal(t, action, *returned.Planet.BuildingAction)
		require.NotNil(t, actual.BuildingAction)
		assert.Equal(t, action, *actual.BuildingAction)
	})

	t.Run("persists mutated planet with action and costs", func(t *testing.T) {
		planet, _, _ := insertTestPlanetForPlayer(t, conn)
		require.Nil(t, planet.BuildingAction)

		action := models.BuildingAction{
			Id:           uuid.New(),
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

		assert.False(t, returned.Deleted)
		actual := loadPlanetFromDb(t, conn, planet.Id)
		assert.Equal(t, returned.Planet, actual)
		require.NotNil(t, returned.Planet.BuildingAction)
		assert.Equal(t, action, *returned.Planet.BuildingAction)
		require.NotNil(t, actual.BuildingAction)
		assert.Equal(t, action, *actual.BuildingAction)
	})

	t.Run("persists mutated planet with action and storages", func(t *testing.T) {
		planet, _, _ := insertTestPlanetForPlayer(t, conn)
		require.Nil(t, planet.BuildingAction)

		action := models.BuildingAction{
			Id:           uuid.New(),
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

		assert.False(t, returned.Deleted)
		actual := loadPlanetFromDb(t, conn, planet.Id)
		assert.Equal(t, returned.Planet, actual)
		require.NotNil(t, returned.Planet.BuildingAction)
		assert.Equal(t, action, *returned.Planet.BuildingAction)
		require.NotNil(t, actual.BuildingAction)
		assert.Equal(t, action, *actual.BuildingAction)
	})

	t.Run("persists mutated planet with action and productions", func(t *testing.T) {
		planet, _, _ := insertTestPlanetForPlayer(t, conn)
		require.Nil(t, planet.BuildingAction)

		action := models.BuildingAction{
			Id:           uuid.New(),
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

		assert.False(t, returned.Deleted)
		actual := loadPlanetFromDb(t, conn, planet.Id)
		assert.Equal(t, returned.Planet, actual)
		require.NotNil(t, returned.Planet.BuildingAction)
		assert.Equal(t, action, *returned.Planet.BuildingAction)
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

		assert.False(t, returned.Deleted)
		actual := loadPlanetFromDb(t, conn, planet.Id)
		assert.Equal(t, returned.Planet, actual)
		require.NotNil(t, returned.Planet.BuildingAction)
		assert.Equal(t, yetAnotherTime, returned.Planet.BuildingAction.CompletedAt)
		require.NotNil(t, actual.BuildingAction)
		assert.Equal(t, yetAnotherTime, actual.BuildingAction.CompletedAt)
	})

	t.Run("persists mutated planet with update to existing action costs", func(t *testing.T) {
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

		assert.False(t, returned.Deleted)
		actual := loadPlanetFromDb(t, conn, planet.Id)
		require.NotNil(t, returned.Planet.BuildingAction)
		assert.Equal(t, costs, returned.Planet.BuildingAction.Costs)
		require.NotNil(t, actual.BuildingAction)
		assert.Equal(t, costs, actual.BuildingAction.Costs)
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

		assert.False(t, returned.Deleted)
		actual := loadPlanetFromDb(t, conn, planet.Id)
		assert.Equal(t, returned.Planet, actual)
		require.NotNil(t, returned.Planet.BuildingAction)
		assert.Equal(t, costs, returned.Planet.BuildingAction.Costs)
		require.NotNil(t, actual.BuildingAction)
		assert.Equal(t, costs, actual.BuildingAction.Costs)
	})

	t.Run("persists mutated planet with update to existing action storages", func(t *testing.T) {
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

		assert.False(t, returned.Deleted)
		actual := loadPlanetFromDb(t, conn, planet.Id)
		require.NotNil(t, returned.Planet.BuildingAction)
		assert.Equal(t, storages, returned.Planet.BuildingAction.Storages)
		require.NotNil(t, actual.BuildingAction)
		assert.Equal(t, storages, actual.BuildingAction.Storages)
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

		assert.False(t, returned.Deleted)
		actual := loadPlanetFromDb(t, conn, planet.Id)
		assert.Equal(t, returned.Planet, actual)
		require.NotNil(t, returned.Planet.BuildingAction)
		assert.Equal(t, storages, returned.Planet.BuildingAction.Storages)
		require.NotNil(t, actual.BuildingAction)
		assert.Equal(t, storages, actual.BuildingAction.Storages)
	})

	t.Run("persists mutated planet with update to existing action productions", func(t *testing.T) {
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

		assert.False(t, returned.Deleted)
		actual := loadPlanetFromDb(t, conn, planet.Id)
		require.NotNil(t, returned.Planet.BuildingAction)
		assert.Equal(t, productions, returned.Planet.BuildingAction.Productions)
		require.NotNil(t, actual.BuildingAction)
		assert.Equal(t, productions, actual.BuildingAction.Productions)
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

		assert.False(t, returned.Deleted)
		actual := loadPlanetFromDb(t, conn, planet.Id)
		assert.Equal(t, returned.Planet, actual)
		require.NotNil(t, returned.Planet.BuildingAction)
		assert.Equal(t, productions, returned.Planet.BuildingAction.Productions)
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

		assert.False(t, returned.Deleted)
		assertBuildingActionDoesNotExist(t, conn, action.Id)
		assert.Nil(t, returned.Planet.BuildingAction)
	})

	t.Run("persists mutated planet with deleted action with costs", func(t *testing.T) {
		action, planet := insertTestBuildingAction(t, conn, addBuildingActionCost)

		mutator := generateModifyingMutator(func(p *models.Planet) {
			p.BuildingAction = nil
			p.Version++
		})

		returned, err := adapter.Mutate(t.Context(), planet.Id, mutator)
		require.NoError(t, err, "Actual err: %v", err)

		assert.False(t, returned.Deleted)
		assertBuildingActionDoesNotExist(t, conn, action.Id)
		assertBuildingActionCostDoesNotExist(t, conn, action.Id)
		assert.Nil(t, returned.Planet.BuildingAction)
	})

	t.Run("persists mutated planet with deleted action with storages", func(t *testing.T) {
		action, planet := insertTestBuildingAction(t, conn, addBuildingActionStorage)

		mutator := generateModifyingMutator(func(p *models.Planet) {
			p.BuildingAction = nil
			p.Version++
		})

		returned, err := adapter.Mutate(t.Context(), planet.Id, mutator)
		require.NoError(t, err, "Actual err: %v", err)

		assert.False(t, returned.Deleted)
		assertBuildingActionDoesNotExist(t, conn, action.Id)
		assertBuildingActionStorageDoesNotExist(t, conn, action.Id)
		assert.Nil(t, returned.Planet.BuildingAction)
	})

	t.Run("persists mutated planet with deleted action with productions", func(t *testing.T) {
		action, planet := insertTestBuildingAction(t, conn, addBuildingActionProduction)

		mutator := generateModifyingMutator(func(p *models.Planet) {
			p.BuildingAction = nil
			p.Version++
		})

		returned, err := adapter.Mutate(t.Context(), planet.Id, mutator)
		require.NoError(t, err, "Actual err: %v", err)

		assert.False(t, returned.Deleted)
		assertBuildingActionDoesNotExist(t, conn, action.Id)
		assertBuildingActionProductionDoesNotExist(t, conn, action.Id)
		assert.Nil(t, returned.Planet.BuildingAction)
	})

	t.Run("persists mutated planet with new action", func(t *testing.T) {
		action, planet := insertTestBuildingAction(t, conn)
		require.NotEqual(t, crystalMineId, action.Building)

		newAction := models.BuildingAction{
			Id:           uuid.New(),
			Building:     crystalMineId,
			DesiredLevel: 4,
			CreatedAt:    someTime,
			CompletedAt:  someTime.Add(1 * time.Hour),
			Costs:        []models.BuildingActionCost{},
			Storages:     []models.BuildingActionResourceStorage{},
			Productions:  []models.BuildingActionResourceProduction{},
		}

		mutator := generateModifyingMutator(func(p *models.Planet) {
			p.BuildingAction = &newAction
			p.Version++
		})

		returned, err := adapter.Mutate(t.Context(), planet.Id, mutator)
		require.NoError(t, err, "Actual err: %v", err)

		assert.False(t, returned.Deleted)
		actual := loadPlanetFromDb(t, conn, planet.Id)
		assert.Equal(t, returned.Planet, actual)
		require.NotNil(t, returned.Planet.BuildingAction)
		assert.Equal(t, newAction, *returned.Planet.BuildingAction)
		require.NotNil(t, actual.BuildingAction)
		assert.Equal(t, newAction, *actual.BuildingAction)
	})

	t.Run("returns error when planet does not exist", func(t *testing.T) {
		mutator := generateModifyingMutator(func(p *models.Planet) {
			p.UpdatedAt = yetAnotherTime
			p.Version++
		})

		returned, err := adapter.Mutate(t.Context(), uuid.New(), mutator)

		assert.False(t, returned.Deleted)
		assert.ErrorIs(t, err, domainerrors.ErrNotFound, "Actual err: %v", err)
	})

	t.Run("returns error when mutator does not update version", func(t *testing.T) {
		planet, _, _ := insertTestPlanetForPlayer(t, conn)

		mutator := generateModifyingMutator(func(*models.Planet) {})

		_, err := adapter.Mutate(t.Context(), planet.Id, mutator)

		assert.ErrorIs(t, err, domainerrors.ErrMutationWithoutVersionBump, "Actual err: %v", err)
	})

	t.Run("returns error when new building is added", func(t *testing.T) {
		planet, _, _ := insertTestPlanetForPlayer(t, conn)

		mutator := generateModifyingMutator(func(p *models.Planet) {
			p.Buildings = []models.PlanetBuilding{{Building: metalMineId, Level: 5}}
			p.Version++
		})

		_, err := adapter.Mutate(t.Context(), planet.Id, mutator)

		assert.ErrorIs(t, err, domainerrors.ErrBuildingNotFound, "Actual err: %v", err)
	})

	t.Run("returns error when new storage is added", func(t *testing.T) {
		planet, _, _ := insertTestPlanetForPlayer(t, conn)

		mutator := generateModifyingMutator(func(p *models.Planet) {
			p.Storages = []models.PlanetResourceStorage{{Resource: crystalResourceId, Storage: 5478}}
			p.Version++
		})

		_, err := adapter.Mutate(t.Context(), planet.Id, mutator)

		assert.ErrorIs(t, err, domainerrors.ErrResourceNotFound, "Actual err: %v", err)
	})

	t.Run("deletes planet when mutator indicates it", func(t *testing.T) {
		planet, _, _ := insertTestPlanetForPlayer(t, conn)
		require.NotEqual(t, planet.UpdatedAt, yetAnotherTime)

		returned, err := adapter.Mutate(t.Context(), planet.Id, generateDeletingMutator())
		require.NoError(t, err, "Actual err: %v", err)

		assert.True(t, returned.Deleted)
		assertPlanetDoesNotExist(t, conn, planet.Id)
	})

	t.Run("deletes planet with resources when mutator indicates it", func(t *testing.T) {
		planet, _, _ := insertTestPlanetForPlayer(t, conn, addPlanetResource)
		require.NotEqual(t, planet.UpdatedAt, yetAnotherTime)

		returned, err := adapter.Mutate(t.Context(), planet.Id, generateDeletingMutator())
		require.NoError(t, err, "Actual err: %v", err)

		assert.True(t, returned.Deleted)
		assertPlanetDoesNotExist(t, conn, planet.Id)
		assertPlanetResourceDoesNotExist(t, conn, planet.Id)
	})

	t.Run("deletes planet with storages when mutator indicates it", func(t *testing.T) {
		planet, _, _ := insertTestPlanetForPlayer(t, conn, addPlanetStorage)
		require.NotEqual(t, planet.UpdatedAt, yetAnotherTime)

		returned, err := adapter.Mutate(t.Context(), planet.Id, generateDeletingMutator())
		require.NoError(t, err, "Actual err: %v", err)

		assert.True(t, returned.Deleted)
		assertPlanetDoesNotExist(t, conn, planet.Id)
		assertPlanetStorageDoesNotExist(t, conn, planet.Id)
	})

	t.Run("deletes planet with productions when mutator indicates it", func(t *testing.T) {
		planet, _, _ := insertTestPlanetForPlayer(t, conn, addPlanetProduction)
		require.NotEqual(t, planet.UpdatedAt, yetAnotherTime)

		returned, err := adapter.Mutate(t.Context(), planet.Id, generateDeletingMutator())
		require.NoError(t, err, "Actual err: %v", err)

		assert.True(t, returned.Deleted)
		assertPlanetDoesNotExist(t, conn, planet.Id)
		assertPlanetProductionDoesNotExist(t, conn, planet.Id)
	})

	t.Run("deletes planet with production for building when mutator indicates it", func(t *testing.T) {
		planet, _, _ := insertTestPlanetForPlayer(t, conn, addPlanetProductionForBuilding)
		require.NotEqual(t, planet.UpdatedAt, yetAnotherTime)

		returned, err := adapter.Mutate(t.Context(), planet.Id, generateDeletingMutator())
		require.NoError(t, err, "Actual err: %v", err)

		assert.True(t, returned.Deleted)
		assertPlanetDoesNotExist(t, conn, planet.Id)
		assertPlanetProductionDoesNotExist(t, conn, planet.Id)
	})

	t.Run("deletes planet with buildings when mutator indicates it", func(t *testing.T) {
		planet, _, _ := insertTestPlanetForPlayer(t, conn, addPlanetBuilding)
		require.NotEqual(t, planet.UpdatedAt, yetAnotherTime)

		returned, err := adapter.Mutate(t.Context(), planet.Id, generateDeletingMutator())
		require.NoError(t, err, "Actual err: %v", err)

		assert.True(t, returned.Deleted)
		assertPlanetDoesNotExist(t, conn, planet.Id)
		assertPlanetBuildingDoesNotExist(t, conn, planet.Id)
	})

	t.Run("deletes planet with building action when mutator indicates it", func(t *testing.T) {
		planet, _, _ := insertTestPlanetForPlayer(t, conn, addPlanetBuildingAction)
		require.NotEqual(t, planet.UpdatedAt, yetAnotherTime)

		returned, err := adapter.Mutate(t.Context(), planet.Id, generateDeletingMutator())
		require.NoError(t, err, "Actual err: %v", err)

		assert.True(t, returned.Deleted)
		assertPlanetDoesNotExist(t, conn, planet.Id)
		require.NotNil(t, planet.BuildingAction)
		assertBuildingActionDoesNotExist(t, conn, planet.BuildingAction.Id)
	})
}

func TestIT_PlanetMutator_Mutate_Concurrency(t *testing.T) {
	adapter, conn := newTestPlanetMutator(t)

	t.Run("blocks concurrent mutation for same planet", func(t *testing.T) {
		planet, _, _ := insertTestPlanetForPlayer(t, conn, addPlanetResource)
		require.NotEqual(t, 9876.0, planet.Resources[0].Amount)

		enteredA := make(chan struct{})
		releaseA := make(chan struct{})
		enteredB := make(chan struct{})
		doneA := make(chan error, 1)
		doneB := make(chan error, 1)

		go func() {
			_, err := adapter.Mutate(t.Context(), planet.Id, func(p *models.Planet) (bool, error) {
				close(enteredA)
				<-releaseA
				p.UpdatedAt = yetAnotherTime
				p.Version++
				return false, nil
			})
			doneA <- err
		}()

		<-enteredA

		blockingCtx, cancel := context.WithTimeout(t.Context(), 100*time.Millisecond)
		defer cancel()

		go func() {
			_, err := adapter.Mutate(blockingCtx, planet.Id, func(p *models.Planet) (bool, error) {
				close(enteredB)
				p.Resources[0].Amount = 9877.0
				p.Version++
				return false, nil
			})
			doneB <- err
		}()

		err := <-doneB
		require.ErrorIs(
			t, err, context.DeadlineExceeded, "Expected deadline exceeded, got err: %v", err,
		)

		select {
		case <-enteredB:
			t.Fatalf("second mutation reached callback even though first mutation held lock")
		default:
		}

		close(releaseA)

		err = <-doneA
		require.NoError(t, err, "Actual err: %v", err)

		result, err := adapter.Mutate(t.Context(), planet.Id, func(p *models.Planet) (bool, error) {
			p.Version++
			p.Resources[0].Amount = 9878.0
			return false, nil
		})
		require.NoError(t, err, "Actual err: %v", err)

		assert.False(t, result.Deleted)
		assert.Equal(t, yetAnotherTime, result.Planet.UpdatedAt)
		assert.Equal(t, 9878.0, result.Planet.Resources[0].Amount)
		assert.Equal(t, planet.Version+2, result.Planet.Version)
	})

	t.Run("does not block concurrent mutation for different planets", func(t *testing.T) {
		planetA, _, _ := insertTestPlanetForPlayer(t, conn)
		require.NotEqual(t, yetAnotherTime, planetA.UpdatedAt)
		planetB, _, _ := insertTestPlanetForPlayer(t, conn, addPlanetResource)
		require.NotEqual(t, 8765.0, planetB.Resources[0].Amount)

		enteredA := make(chan struct{})
		releaseA := make(chan struct{})
		doneA := make(chan error, 1)
		doneB := make(chan error, 1)

		go func() {
			_, err := adapter.Mutate(t.Context(), planetA.Id, func(p *models.Planet) (bool, error) {
				close(enteredA)
				<-releaseA
				p.UpdatedAt = yetAnotherTime
				p.Version++
				return false, nil
			})
			doneA <- err
		}()

		<-enteredA

		go func() {
			_, err := adapter.Mutate(t.Context(), planetB.Id, func(p *models.Planet) (bool, error) {
				p.Resources[0].Amount = 8765.0
				p.Version++
				return false, nil
			})
			doneB <- err
		}()

		select {
		case err := <-doneB:
			require.NoError(t, err, "Actual err: %v", err)
		case <-time.After(500 * time.Millisecond):
			t.Fatalf("mutation on second planet was blocked while first planet was locked")
		}

		assertPlanetResourceAmount(t, conn, planetB.Id, crystalResourceId, 8765.0)

		close(releaseA)

		err := <-doneA
		require.NoError(t, err, "Actual err: %v", err)

		resultA, err := adapter.Mutate(t.Context(), planetA.Id, func(p *models.Planet) (bool, error) {
			p.Version++
			return false, nil
		})
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, yetAnotherTime, resultA.Planet.UpdatedAt)
		assertPlanetResourceAmount(
			t, conn, planetB.Id, planetB.Resources[0].Resource, 8765.0,
		)
	})

	t.Run("waiting mutation respects context timeout", func(t *testing.T) {
		planet, _, _ := insertTestPlanetForPlayer(t, conn, addPlanetResource)

		enteredA := make(chan struct{})
		releaseA := make(chan struct{})
		enteredB := make(chan struct{})
		doneA := make(chan error, 1)
		doneB := make(chan error, 1)

		go func() {
			_, err := adapter.Mutate(t.Context(), planet.Id, func(p *models.Planet) (bool, error) {
				close(enteredA)
				<-releaseA
				p.UpdatedAt = yetAnotherTime
				p.Version++
				return false, nil
			})
			doneA <- err
		}()

		<-enteredA

		blockingCtx, cancel := context.WithTimeout(t.Context(), 100*time.Millisecond)
		defer cancel()

		go func() {
			_, err := adapter.Mutate(blockingCtx, planet.Id, func(p *models.Planet) (bool, error) {
				close(enteredB)
				p.Resources[0].Amount = 2222.0
				p.Version++
				return false, nil
			})
			doneB <- err
		}()

		errB := <-doneB
		require.ErrorIs(
			t, errB, context.DeadlineExceeded, "Expected deadline exceeded, got err: %v", errB,
		)

		select {
		case <-enteredB:
			t.Fatalf("waiting mutation callback should not be reached before lock is released")
		default:
		}

		close(releaseA)

		errA := <-doneA
		require.NoError(t, errA, "Actual err: %v", errA)

		result, err := adapter.Mutate(t.Context(), planet.Id, func(p *models.Planet) (bool, error) {
			p.Resources[0].Amount = 3333.0
			p.Version++
			return false, nil
		})
		require.NoError(t, err, "Actual err: %v", err)

		assert.False(t, result.Deleted)
		assert.Equal(t, yetAnotherTime, result.Planet.UpdatedAt)
		assert.Equal(t, 3333.0, result.Planet.Resources[0].Amount)
	})
}

func TestIT_PlanetMutator_ActionCreationDeletionWorkflow(t *testing.T) {
	conn := newTestConnection(t)
	planetMutator := NewPlanetMutator(conn)

	action := models.BuildingAction{
		Id:           uuid.New(),
		Building:     metalMineId,
		DesiredLevel: 27,
		CreatedAt:    time.Date(2024, time.December, 7, 20, 26, 47, 0, time.UTC),
		CompletedAt:  time.Date(2024, time.December, 7, 21, 26, 47, 0, time.UTC),
		Costs:        []models.BuildingActionCost{},
		Storages:     []models.BuildingActionResourceStorage{},
		Productions:  []models.BuildingActionResourceProduction{},
	}

	planet, _, _ := insertTestPlanetForPlayer(t, conn)

	mutation := func(p *models.Planet) (bool, error) {
		p.BuildingAction = &action
		p.Version++
		return false, nil
	}
	result, err := planetMutator.Mutate(t.Context(), planet.Id, mutation)
	require.NoError(t, err, "Actual err: %v", err)
	assert.False(t, result.Deleted)
	require.NotNil(t, result.Planet.BuildingAction)
	assert.Equal(t, action, *result.Planet.BuildingAction)
	assertBuildingActionExists(t, conn, action.Id)

	func() {
		actual := loadPlanetFromDb(t, conn, planet.Id)

		require.NotNil(t, actual.BuildingAction)
		assert.Equal(t, action, *actual.BuildingAction)
	}()

	mutation = func(p *models.Planet) (bool, error) {
		p.BuildingAction = nil
		p.Version++
		return false, nil
	}
	result, err = planetMutator.Mutate(t.Context(), planet.Id, mutation)
	require.NoError(t, err, "Actual err: %v", err)
	assert.False(t, result.Deleted)
	assert.Nil(t, result.Planet.BuildingAction)
	assertBuildingActionDoesNotExist(t, conn, action.Id)
}

func newTestPlanetMutator(t *testing.T) (*PlanetMutator, db.Connection) {
	t.Helper()
	conn := newTestConnection(t)
	return NewPlanetMutator(conn), conn
}

func generateModifyingMutator(modifier func(p *models.Planet)) drivenports.PlanetMutator {
	return func(p *models.Planet) (bool, error) {
		modifier(p)
		return false, nil
	}
}

func generateDeletingMutator() drivenports.PlanetMutator {
	return func(*models.Planet) (bool, error) {
		return true, nil
	}
}

func assertPlanetBuildingLevel(t *testing.T, conn db.Connection, planet uuid.UUID, building uuid.UUID, level int) {
	t.Helper()

	sqlQuery := `SELECT level FROM planet_building WHERE planet = $1 AND building = $2`
	value, err := db.QueryOne[int](t.Context(), conn, sqlQuery, planet, building)
	require.NoError(t, err, "Actual err: %v", err)
	require.Equal(t, level, value)
}
