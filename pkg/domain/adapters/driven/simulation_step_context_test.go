package drivenadapters

import (
	"errors"
	"testing"
	"time"

	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/adapters/driven/database"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	domainerrors "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models/errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	t1 = time.Date(2026, time.October, 31, 16, 21, 23, 0, time.UTC)
	t2 = time.Date(2026, time.October, 31, 17, 21, 23, 0, time.UTC)
)

func TestUnit_SimulationStepContext_NewPlanetStepContext(t *testing.T) {
	t.Run("captures planet version when available", func(t *testing.T) {
		planet := models.Planet{
			Id:      uuid.New(),
			Version: 3,
		}

		ctx := newPlanetStepContext(planet, true, dummyResolver)

		assert.Equal(t, 3, ctx.initialPlanetVersion)
		assert.Equal(t, 0, ctx.initialFleetVersion)
	})

	t.Run("configures planet deletion as configured", func(t *testing.T) {
		planet := models.Planet{
			Id:      uuid.New(),
			Version: 3,
		}

		ctx := newPlanetStepContext(planet, false, dummyResolver)
		assert.False(t, ctx.allowPlanetDeletion)

		ctx = newPlanetStepContext(planet, true, dummyResolver)
		assert.True(t, ctx.allowPlanetDeletion)
	})
}

func TestUnit_SimulationStepContext_NewFleetStepContext(t *testing.T) {
	t.Run("captures fleet version when available", func(t *testing.T) {
		event := models.FleetEventContext{
			Timestamp: someTime,
			Fleet: models.Fleet{
				Id:      uuid.New(),
				Version: 4,
			},
			Planet: nil,
		}

		ctx := newFleetStepContext(event, dummyResolver)

		assert.Equal(t, 0, ctx.initialPlanetVersion)
		assert.Equal(t, 4, ctx.initialFleetVersion)
	})

	t.Run("captures planet and fleet version from fleet event when available", func(t *testing.T) {
		event := models.FleetEventContext{
			Timestamp: someTime,
			Fleet: models.Fleet{
				Id:      uuid.New(),
				Version: 4,
			},
			Planet: &models.Planet{
				Id:      uuid.New(),
				Version: 36,
			},
		}

		ctx := newFleetStepContext(event, dummyResolver)

		assert.Equal(t, 36, ctx.initialPlanetVersion)
		assert.Equal(t, 4, ctx.initialFleetVersion)
	})

	t.Run("disables planet deletion by default", func(t *testing.T) {
		event := models.FleetEventContext{
			Timestamp: someTime,
			Fleet: models.Fleet{
				Id:      uuid.New(),
				Version: 4,
			},
			Planet: nil,
		}

		ctx := newFleetStepContext(event, dummyResolver)

		assert.False(t, ctx.allowPlanetDeletion)
	})
}

func TestUnit_SimulationStepContext_Resolve(t *testing.T) {
	t.Run("forwards resolver error to caller", func(t *testing.T) {
		planet := models.Planet{
			Id:      uuid.New(),
			Version: 3,
		}

		errSample := errors.New("stubbed error")
		resolver := func(_ models.SimulationStep) (models.SimulationStepResult, error) {
			return models.SimulationStepResult{}, errSample
		}

		ctx := newPlanetStepContext(planet, true, resolver)
		_, err := ctx.resolve(t.Context(), nil, t1)

		assert.ErrorIs(t, err, errSample, "Actual err: %v", err)
	})

	t.Run("returns error when resolver requests planet deletion but is not allowed", func(t *testing.T) {
		planet := models.Planet{
			Id:      uuid.New(),
			Version: 3,
		}

		result := models.SimulationStepResult{
			Resolution: models.PlanetDeleted,
			Planet:     nil,
			Fleet:      nil,
		}

		resolver, _ := generateCapturingResolver(result)

		ctx := newPlanetStepContext(planet, false, resolver)

		_, err := ctx.resolve(t.Context(), nil, t2)

		assert.ErrorIs(t, err, domainerrors.ErrInvalidSimulationStepResult, "Actual err: %v", err)
	})

	t.Run("returns error when resolver requests planet deletion for a fleet event", func(t *testing.T) {
		event := models.FleetEventContext{
			Timestamp: someTime,
			Fleet: models.Fleet{
				Id:      uuid.New(),
				Version: 4,
			},
			Planet: nil,
		}

		result := models.SimulationStepResult{
			Resolution: models.PlanetDeleted,
			Planet:     nil,
			Fleet:      nil,
		}

		resolver, _ := generateCapturingResolver(result)

		ctx := newFleetStepContext(event, resolver)

		_, err := ctx.resolve(t.Context(), nil, t2)

		assert.ErrorIs(t, err, domainerrors.ErrInvalidSimulationStepResult, "Actual err: %v", err)
	})

	t.Run("returns error when resolver requests planet update for a fleet event", func(t *testing.T) {
		event := models.FleetEventContext{
			Timestamp: someTime,
			Fleet: models.Fleet{
				Id:      uuid.New(),
				Version: 4,
			},
			Planet: nil,
		}

		result := models.SimulationStepResult{
			Resolution: models.PlanetUpdated,
			Planet:     nil,
			Fleet:      nil,
		}

		resolver, _ := generateCapturingResolver(result)

		ctx := newFleetStepContext(event, resolver)

		_, err := ctx.resolve(t.Context(), nil, t2)

		assert.ErrorIs(t, err, domainerrors.ErrInvalidSimulationStepResult, "Actual err: %v", err)
	})
}

func TestIT_SimulationStepContext_Resolve(t *testing.T) {
	t.Run("calls resolver with correct step for planet update", func(t *testing.T) {
		conn, tx := newTestConnectionAndTransaction(t)
		planet, _, _ := insertTestPlanetForPlayer(t, conn)

		result := models.SimulationStepResult{
			Resolution: models.PlanetUpdated,
			Planet:     new(planet),
			Fleet:      nil,
		}
		result.Planet.Version++

		resolver, captor := generateCapturingResolver(result)

		ctx := newPlanetStepContext(planet, true, resolver)

		actual, err := ctx.resolve(t.Context(), tx, t1)
		require.NoError(t, err, "Actual err: %v", err)

		expected := models.SimulationStep{
			Timestamp:      t1,
			Planet:         new(planet),
			NextFleetEvent: nil,
		}
		assert.Equal(t, expected, *captor)
		assert.Equal(t, result, actual)
	})

	t.Run("calls resolver with correct step for fleet update", func(t *testing.T) {
		conn, tx := newTestConnectionAndTransaction(t)
		planet, _, _ := insertTestPlanetForPlayer(t, conn)
		fleet := insertTestFleet(t, conn, planet)

		event := models.FleetEventContext{
			Timestamp: fleet.Schedule.ArrivalAt,
			Fleet:     fleet,
			Planet:    nil,
		}

		result := models.SimulationStepResult{
			Resolution: models.FleetEventResolved,
			Planet:     nil,
			Fleet:      new(fleet),
		}
		result.Fleet.Version++

		resolver, captor := generateCapturingResolver(result)

		ctx := newFleetStepContext(event, resolver)

		afterArrival := fleet.Schedule.ArrivalAt.Add(1 * time.Minute)
		actual, err := ctx.resolve(t.Context(), tx, afterArrival)
		require.NoError(t, err, "Actual err: %v", err)

		expected := models.SimulationStep{
			Timestamp:      afterArrival,
			Planet:         nil,
			NextFleetEvent: new(event),
		}
		assert.Equal(t, expected, *captor)
		assert.Equal(t, result, actual)
	})

	t.Run("accepts deletion when resolver requests planet deletion", func(t *testing.T) {
		conn, tx := newTestConnectionAndTransaction(t)
		planet, _, _ := insertTestPlanetForPlayer(t, conn)

		result := models.SimulationStepResult{
			Resolution: models.PlanetDeleted,
			Planet:     nil,
			Fleet:      nil,
		}

		resolver, _ := generateCapturingResolver(result)

		ctx := newPlanetStepContext(planet, true, resolver)

		actual, err := ctx.resolve(t.Context(), tx, t2)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, result, actual)
	})
}

func TestIT_SimulationStepContext_Resolve_SavePlanet(t *testing.T) {
	t.Run("returns mutated planet", func(t *testing.T) {
		conn, tx := newTestConnectionAndTransaction(t)

		planet, _, _ := insertTestPlanetForPlayer(t, conn)
		require.NotEqual(t, planet.UpdatedAt, yetAnotherTime)

		resolver := generatePlanetSimulator(func(p *models.Planet) {
			p.UpdatedAt = yetAnotherTime
			p.Version++
		})

		ctx := newPlanetStepContext(planet, false, resolver)

		actual, err := ctx.resolve(t.Context(), tx, yetAnotherTime)
		require.NoError(t, err, "Actual err: %v", err)
		tx.Close(t.Context())

		require.NotNil(t, actual.Planet)
		assert.Equal(t, yetAnotherTime, actual.Planet.UpdatedAt)
		assert.Equal(t, planet.Version+1, actual.Planet.Version)
	})

	t.Run("persists mutated planet", func(t *testing.T) {
		conn, tx := newTestConnectionAndTransaction(t)

		planet, _, _ := insertTestPlanetForPlayer(t, conn)
		require.NotEqual(t, planet.UpdatedAt, yetAnotherTime)
		require.NotEqual(t, 326, planet.Fields)

		resolver := generatePlanetSimulator(func(p *models.Planet) {
			p.Fields = 326
			p.UpdatedAt = yetAnotherTime
			p.Version++
		})

		ctx := newPlanetStepContext(planet, false, resolver)

		_, err := ctx.resolve(t.Context(), tx, yetAnotherTime)
		require.NoError(t, err, "Actual err: %v", err)
		tx.Close(t.Context())

		actual := loadPlanetFromDb(t, conn, planet.Id)
		expected := planet
		expected.Fields = 326
		expected.UpdatedAt = yetAnotherTime
		expected.Version++
		assert.Equal(t, expected, actual)
	})

	t.Run("persists mutated planet resources", func(t *testing.T) {
		conn, tx := newTestConnectionAndTransaction(t)

		planet, _, _ := insertTestPlanetForPlayer(t, conn, addPlanetResource)
		require.NotEqual(t, 5874, planet.Resources[0].Amount)

		resolver := generatePlanetSimulator(func(p *models.Planet) {
			p.Resources[0].Amount = 5874
			p.Version++
		})

		ctx := newPlanetStepContext(planet, false, resolver)

		returned, err := ctx.resolve(t.Context(), tx, planet.UpdatedAt)
		require.NoError(t, err, "Actual err: %v", err)
		tx.Close(t.Context())

		require.NotNil(t, returned.Planet)
		expected := []models.PlanetResource{
			{Resource: crystalResourceId, Amount: 5874},
		}
		assert.Equal(t, expected, returned.Planet.Resources)
		assertPlanetResourceAmount(t, conn, planet.Id, crystalResourceId, 5874)
	})

	t.Run("does not delete existing planet resource", func(t *testing.T) {
		conn, tx := newTestConnectionAndTransaction(t)

		planet, _, _ := insertTestPlanetForPlayer(t, conn, addPlanetResource)
		resource := planet.Resources[0].Resource
		amount := planet.Resources[0].Amount

		resolver := generatePlanetSimulator(func(p *models.Planet) {
			p.Resources = []models.PlanetResource{}
			p.Version++
		})

		ctx := newPlanetStepContext(planet, false, resolver)

		returned, err := ctx.resolve(t.Context(), tx, planet.UpdatedAt)
		require.NoError(t, err, "Actual err: %v", err)
		tx.Close(t.Context())

		require.NotNil(t, returned.Planet)
		expected := []models.PlanetResource{
			{Resource: resource, Amount: amount},
		}
		assert.Equal(t, expected, returned.Planet.Resources)
		assertPlanetResourceAmount(t, conn, planet.Id, resource, amount)
	})

	t.Run("persists mutated planet productions", func(t *testing.T) {
		conn, tx := newTestConnectionAndTransaction(t)

		planet, _, _ := insertTestPlanetForPlayer(t, conn, addPlanetProduction)
		require.NotEqual(t, 39841, planet.Productions[0].Production)

		resolver := generatePlanetSimulator(func(p *models.Planet) {
			p.Productions[0].Production = 39841
			p.Version++
		})

		ctx := newPlanetStepContext(planet, false, resolver)

		returned, err := ctx.resolve(t.Context(), tx, planet.UpdatedAt)
		require.NoError(t, err, "Actual err: %v", err)
		tx.Close(t.Context())

		require.NotNil(t, returned.Planet)
		expected := []models.PlanetResourceProduction{
			{Resource: metalResourceId, Production: 39841},
		}
		assert.Equal(t, expected, returned.Planet.Productions)
		assertPlanetResourceProduction(t, conn, planet.Id, metalResourceId, nil, 39841)
	})

	t.Run("persists additional planet productions", func(t *testing.T) {
		conn, tx := newTestConnectionAndTransaction(t)

		planet, _, _ := insertTestPlanetForPlayer(
			t, conn, addPlanetProduction, addPlanetProductionForBuilding,
		)
		require.NotEqual(t, crystalResourceId, planet.Productions[0].Building)
		initial0 := planet.Productions[0].Production
		initial1 := planet.Productions[1].Production

		resolver := generatePlanetSimulator(func(p *models.Planet) {
			prod := models.PlanetResourceProduction{
				Resource:   crystalResourceId,
				Production: 354789,
			}
			p.Productions = append(p.Productions, prod)
			p.Version++
		})

		ctx := newPlanetStepContext(planet, false, resolver)

		returned, err := ctx.resolve(t.Context(), tx, planet.UpdatedAt)
		require.NoError(t, err, "Actual err: %v", err)
		tx.Close(t.Context())

		require.NotNil(t, returned.Planet)
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
		conn, tx := newTestConnectionAndTransaction(t)

		planet, _, _ := insertTestPlanetForPlayer(t, conn, addPlanetProduction)

		resolver := generatePlanetSimulator(func(p *models.Planet) {
			p.Productions = []models.PlanetResourceProduction{}
			p.Version++
		})

		ctx := newPlanetStepContext(planet, false, resolver)

		returned, err := ctx.resolve(t.Context(), tx, planet.UpdatedAt)
		require.NoError(t, err, "Actual err: %v", err)
		tx.Close(t.Context())

		require.NotNil(t, returned.Planet)
		assert.Equal(t, []models.PlanetResourceProduction{}, returned.Planet.Productions)
		assertPlanetProductionDoesNotExist(t, conn, planet.Id)
	})

	t.Run("persists mutated planet production for building", func(t *testing.T) {
		conn, tx := newTestConnectionAndTransaction(t)

		planet, _, _ := insertTestPlanetForPlayer(t, conn, addPlanetProductionForBuilding)
		require.NotEqual(t, 1235, planet.Productions[0].Production)

		resolver := generatePlanetSimulator(func(p *models.Planet) {
			p.Productions[0].Production = 1235
			p.Version++
		})

		ctx := newPlanetStepContext(planet, false, resolver)

		returned, err := ctx.resolve(t.Context(), tx, planet.UpdatedAt)
		require.NoError(t, err, "Actual err: %v", err)
		tx.Close(t.Context())

		require.NotNil(t, returned.Planet)
		expected := []models.PlanetResourceProduction{
			{Resource: metalResourceId, Building: &crystalMineId, Production: 1235},
		}
		assert.Equal(t, expected, returned.Planet.Productions)
		assertPlanetResourceProduction(
			t, conn, planet.Id, metalResourceId, &crystalMineId, 1235,
		)
	})

	t.Run("persists additional planet productions for building", func(t *testing.T) {
		conn, tx := newTestConnectionAndTransaction(t)

		planet, _, _ := insertTestPlanetForPlayer(
			t, conn, addPlanetProduction, addPlanetProductionForBuilding,
		)
		require.NotEqual(t, metalMineId, planet.Productions[1].Building)
		initial0 := planet.Productions[0].Production
		initial1 := planet.Productions[1].Production

		resolver := generatePlanetSimulator(func(p *models.Planet) {
			prod := models.PlanetResourceProduction{
				Resource:   crystalResourceId,
				Building:   &metalMineId,
				Production: 354789,
			}
			p.Productions = append(p.Productions, prod)
			p.Version++
		})

		ctx := newPlanetStepContext(planet, false, resolver)

		returned, err := ctx.resolve(t.Context(), tx, planet.UpdatedAt)
		require.NoError(t, err, "Actual err: %v", err)
		tx.Close(t.Context())

		require.NotNil(t, returned.Planet)
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
		conn, tx := newTestConnectionAndTransaction(t)

		planet, _, _ := insertTestPlanetForPlayer(t, conn, addPlanetProductionForBuilding)

		resolver := generatePlanetSimulator(func(p *models.Planet) {
			p.Productions = []models.PlanetResourceProduction{}
			p.Version++
		})

		ctx := newPlanetStepContext(planet, false, resolver)

		returned, err := ctx.resolve(t.Context(), tx, planet.UpdatedAt)
		require.NoError(t, err, "Actual err: %v", err)
		tx.Close(t.Context())

		require.NotNil(t, returned.Planet)
		assert.Equal(t, []models.PlanetResourceProduction{}, returned.Planet.Productions)
		assertPlanetProductionDoesNotExist(t, conn, planet.Id)
	})

	t.Run("persists mutated planet storages", func(t *testing.T) {
		conn, tx := newTestConnectionAndTransaction(t)

		planet, _, _ := insertTestPlanetForPlayer(t, conn, addPlanetStorage)
		require.NotEqual(t, 4598, planet.Storages[0].Storage)

		resolver := generatePlanetSimulator(func(p *models.Planet) {
			p.Storages[0].Storage = 4598
			p.Version++
		})

		ctx := newPlanetStepContext(planet, false, resolver)

		returned, err := ctx.resolve(t.Context(), tx, planet.UpdatedAt)
		require.NoError(t, err, "Actual err: %v", err)
		tx.Close(t.Context())

		require.NotNil(t, returned.Planet)
		expected := []models.PlanetResourceStorage{
			{Resource: crystalResourceId, Storage: 4598},
		}
		assert.Equal(t, expected, returned.Planet.Storages)
		assertPlanetResourceStorage(t, conn, planet.Id, crystalResourceId, 4598)
	})

	t.Run("does not delete existing planet storage", func(t *testing.T) {
		conn, tx := newTestConnectionAndTransaction(t)

		planet, _, _ := insertTestPlanetForPlayer(t, conn, addPlanetStorage)
		resource := planet.Storages[0].Resource
		storage := planet.Storages[0].Storage

		resolver := generatePlanetSimulator(func(p *models.Planet) {
			p.Storages = []models.PlanetResourceStorage{}
			p.Version++
		})

		ctx := newPlanetStepContext(planet, false, resolver)

		returned, err := ctx.resolve(t.Context(), tx, planet.UpdatedAt)
		require.NoError(t, err, "Actual err: %v", err)
		tx.Close(t.Context())

		require.NotNil(t, returned.Planet)
		expected := []models.PlanetResourceStorage{
			{Resource: resource, Storage: storage},
		}
		assert.Equal(t, expected, returned.Planet.Storages)
		assertPlanetResourceStorage(t, conn, planet.Id, resource, storage)
	})

	t.Run("persists mutated planet buildings", func(t *testing.T) {
		conn, tx := newTestConnectionAndTransaction(t)

		planet, _, _ := insertTestPlanetForPlayer(t, conn, addPlanetBuilding)
		require.NotEqual(t, 6, planet.Buildings[0].Level)

		resolver := generatePlanetSimulator(func(p *models.Planet) {
			p.Buildings[0].Level = 6
			p.Version++
		})

		ctx := newPlanetStepContext(planet, false, resolver)

		returned, err := ctx.resolve(t.Context(), tx, planet.UpdatedAt)
		require.NoError(t, err, "Actual err: %v", err)
		tx.Close(t.Context())

		require.NotNil(t, returned.Planet)
		expected := []models.PlanetBuilding{
			{Building: metalStorageId, Level: 6},
		}
		assert.Equal(t, expected, returned.Planet.Buildings)
		assertPlanetBuildingLevel(t, conn, planet.Id, metalStorageId, 6)
	})

	t.Run("does not delete existing planet building", func(t *testing.T) {
		conn, tx := newTestConnectionAndTransaction(t)

		planet, _, _ := insertTestPlanetForPlayer(t, conn, addPlanetBuilding)
		building := planet.Buildings[0].Building
		level := planet.Buildings[0].Level

		resolver := generatePlanetSimulator(func(p *models.Planet) {
			p.Buildings = []models.PlanetBuilding{}
			p.Version++
		})

		ctx := newPlanetStepContext(planet, false, resolver)

		returned, err := ctx.resolve(t.Context(), tx, planet.UpdatedAt)
		require.NoError(t, err, "Actual err: %v", err)
		tx.Close(t.Context())

		require.NotNil(t, returned.Planet)
		expected := []models.PlanetBuilding{
			{Building: building, Level: level},
		}
		assert.Equal(t, expected, returned.Planet.Buildings)
		assertPlanetBuildingLevel(t, conn, planet.Id, building, level)
	})

	t.Run("persists mutated planet ships", func(t *testing.T) {
		conn, tx := newTestConnectionAndTransaction(t)

		planet, _, _ := insertTestPlanetForPlayer(t, conn, addPlanetShip)
		require.NotEqual(t, 776, planet.Ships[0].Count)

		resolver := generatePlanetSimulator(func(p *models.Planet) {
			p.Ships[0].Count = 776
			p.Version++
		})

		ctx := newPlanetStepContext(planet, false, resolver)

		returned, err := ctx.resolve(t.Context(), tx, planet.UpdatedAt)
		require.NoError(t, err, "Actual err: %v", err)
		tx.Close(t.Context())

		require.NotNil(t, returned.Planet)
		expected := []models.PlanetShip{
			{Ship: lightFighterId, Count: 776},
		}
		assert.Equal(t, expected, returned.Planet.Ships)
		assertPlanetShipCount(t, conn, planet.Id, lightFighterId, 776)
	})

	t.Run("does not delete existing planet ship", func(t *testing.T) {
		conn, tx := newTestConnectionAndTransaction(t)

		planet, _, _ := insertTestPlanetForPlayer(t, conn, addPlanetShip)
		ship := planet.Ships[0].Ship
		count := planet.Ships[0].Count

		resolver := generatePlanetSimulator(func(p *models.Planet) {
			p.Ships = []models.PlanetShip{}
			p.Version++
		})

		ctx := newPlanetStepContext(planet, false, resolver)

		returned, err := ctx.resolve(t.Context(), tx, planet.UpdatedAt)
		require.NoError(t, err, "Actual err: %v", err)
		tx.Close(t.Context())

		require.NotNil(t, returned.Planet)
		expected := []models.PlanetShip{
			{Ship: ship, Count: count},
		}
		assert.Equal(t, expected, returned.Planet.Ships)
		assertPlanetShipCount(t, conn, planet.Id, ship, count)
	})

	t.Run("persists mutated planet with action", func(t *testing.T) {
		conn, tx := newTestConnectionAndTransaction(t)

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

		resolver := generatePlanetSimulator(func(p *models.Planet) {
			p.BuildingAction = &action
			p.Version++
		})

		ctx := newPlanetStepContext(planet, false, resolver)

		returned, err := ctx.resolve(t.Context(), tx, planet.UpdatedAt)
		require.NoError(t, err, "Actual err: %v", err)
		tx.Close(t.Context())

		actual := loadPlanetFromDb(t, conn, planet.Id)
		assert.Equal(t, &actual, returned.Planet)
		require.NotNil(t, returned.Planet.BuildingAction)
		assert.Equal(t, action, *returned.Planet.BuildingAction)
	})

	t.Run("persists mutated planet with action and costs", func(t *testing.T) {
		conn, tx := newTestConnectionAndTransaction(t)

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

		resolver := generatePlanetSimulator(func(p *models.Planet) {
			p.BuildingAction = &action
			p.Version++
		})

		ctx := newPlanetStepContext(planet, false, resolver)

		returned, err := ctx.resolve(t.Context(), tx, planet.UpdatedAt)
		require.NoError(t, err, "Actual err: %v", err)
		tx.Close(t.Context())

		actual := loadPlanetFromDb(t, conn, planet.Id)
		assert.Equal(t, &actual, returned.Planet)
		require.NotNil(t, returned.Planet.BuildingAction)
		assert.Equal(t, action, *returned.Planet.BuildingAction)
	})

	t.Run("persists mutated planet with action and storages", func(t *testing.T) {
		conn, tx := newTestConnectionAndTransaction(t)

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

		resolver := generatePlanetSimulator(func(p *models.Planet) {
			p.BuildingAction = &action
			p.Version++
		})

		ctx := newPlanetStepContext(planet, false, resolver)

		returned, err := ctx.resolve(t.Context(), tx, planet.UpdatedAt)
		require.NoError(t, err, "Actual err: %v", err)
		tx.Close(t.Context())

		actual := loadPlanetFromDb(t, conn, planet.Id)
		assert.Equal(t, &actual, returned.Planet)
		require.NotNil(t, returned.Planet.BuildingAction)
		assert.Equal(t, action, *returned.Planet.BuildingAction)
	})

	t.Run("persists mutated planet with action and productions", func(t *testing.T) {
		conn, tx := newTestConnectionAndTransaction(t)

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

		resolver := generatePlanetSimulator(func(p *models.Planet) {
			p.BuildingAction = &action
			p.Version++
		})

		ctx := newPlanetStepContext(planet, false, resolver)

		returned, err := ctx.resolve(t.Context(), tx, planet.UpdatedAt)
		require.NoError(t, err, "Actual err: %v", err)
		tx.Close(t.Context())

		actual := loadPlanetFromDb(t, conn, planet.Id)
		assert.Equal(t, &actual, returned.Planet)
		require.NotNil(t, returned.Planet.BuildingAction)
		assert.Equal(t, action, *returned.Planet.BuildingAction)
	})

	t.Run("persists mutated planet with updated action completion time", func(t *testing.T) {
		conn, tx := newTestConnectionAndTransaction(t)

		planet, _, _ := insertTestPlanetForPlayer(t, conn)
		action := insertTestBuildingActionForPlanet(t, conn, planet.Id)
		require.NotEqual(t, yetAnotherTime, action.CompletedAt)
		planet.BuildingAction = &action

		resolver := generatePlanetSimulator(func(p *models.Planet) {
			p.BuildingAction.CompletedAt = yetAnotherTime
			p.Version++
		})

		ctx := newPlanetStepContext(planet, false, resolver)

		returned, err := ctx.resolve(t.Context(), tx, planet.UpdatedAt)
		require.NoError(t, err, "Actual err: %v", err)
		tx.Close(t.Context())

		actual := loadPlanetFromDb(t, conn, planet.Id)
		assert.Equal(t, &actual, returned.Planet)
		require.NotNil(t, returned.Planet.BuildingAction)
		assert.Equal(t, yetAnotherTime, returned.Planet.BuildingAction.CompletedAt)
	})

	t.Run("persists mutated planet with update to existing action costs", func(t *testing.T) {
		conn, tx := newTestConnectionAndTransaction(t)

		planet, _, _ := insertTestPlanetForPlayer(t, conn)
		action := insertTestBuildingActionForPlanet(t, conn, planet.Id, addBuildingActionCost)
		require.Equal(t, metalResourceId, action.Costs[0].Resource)
		require.NotEqual(t, 32, action.Costs[0].Amount)
		planet.BuildingAction = &action

		costs := []models.BuildingActionCost{{Resource: metalResourceId, Amount: 32}}

		resolver := generatePlanetSimulator(func(p *models.Planet) {
			p.BuildingAction.Costs = costs
			p.Version++
		})

		ctx := newPlanetStepContext(planet, false, resolver)

		returned, err := ctx.resolve(t.Context(), tx, planet.UpdatedAt)
		require.NoError(t, err, "Actual err: %v", err)
		tx.Close(t.Context())

		actual := loadPlanetFromDb(t, conn, planet.Id)
		assert.Equal(t, &actual, returned.Planet)
		require.NotNil(t, returned.Planet.BuildingAction)
		assert.Equal(t, costs, returned.Planet.BuildingAction.Costs)
	})

	t.Run("persists mutated planet with updated action and costs", func(t *testing.T) {
		conn, tx := newTestConnectionAndTransaction(t)

		planet, _, _ := insertTestPlanetForPlayer(t, conn)
		action := insertTestBuildingActionForPlanet(t, conn, planet.Id)
		planet.BuildingAction = &action

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

		resolver := generatePlanetSimulator(func(p *models.Planet) {
			p.BuildingAction.Costs = costs
			p.Version++
		})

		ctx := newPlanetStepContext(planet, false, resolver)

		returned, err := ctx.resolve(t.Context(), tx, planet.UpdatedAt)
		require.NoError(t, err, "Actual err: %v", err)
		tx.Close(t.Context())

		actual := loadPlanetFromDb(t, conn, planet.Id)
		assert.Equal(t, &actual, returned.Planet)
		require.NotNil(t, returned.Planet.BuildingAction)
		assert.Equal(t, costs, returned.Planet.BuildingAction.Costs)
	})

	t.Run("persists mutated planet with update to existing action storages", func(t *testing.T) {
		conn, tx := newTestConnectionAndTransaction(t)

		planet, _, _ := insertTestPlanetForPlayer(t, conn)
		action := insertTestBuildingActionForPlanet(t, conn, planet.Id, addBuildingActionStorage)
		require.Equal(t, crystalResourceId, action.Storages[0].Resource)
		require.NotEqual(t, 32, action.Storages[0].Storage)
		planet.BuildingAction = &action

		storages := []models.BuildingActionResourceStorage{{Resource: crystalResourceId, Storage: 32}}

		resolver := generatePlanetSimulator(func(p *models.Planet) {
			p.BuildingAction.Storages = storages
			p.Version++
		})

		ctx := newPlanetStepContext(planet, false, resolver)

		returned, err := ctx.resolve(t.Context(), tx, planet.UpdatedAt)
		require.NoError(t, err, "Actual err: %v", err)
		tx.Close(t.Context())

		actual := loadPlanetFromDb(t, conn, planet.Id)
		assert.Equal(t, &actual, returned.Planet)
		require.NotNil(t, returned.Planet.BuildingAction)
		assert.Equal(t, storages, returned.Planet.BuildingAction.Storages)
	})

	t.Run("persists mutated planet with updated action and storages", func(t *testing.T) {
		conn, tx := newTestConnectionAndTransaction(t)

		planet, _, _ := insertTestPlanetForPlayer(t, conn)
		action := insertTestBuildingActionForPlanet(t, conn, planet.Id)
		planet.BuildingAction = &action

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

		resolver := generatePlanetSimulator(func(p *models.Planet) {
			p.BuildingAction.Storages = storages
			p.Version++
		})

		ctx := newPlanetStepContext(planet, false, resolver)

		returned, err := ctx.resolve(t.Context(), tx, planet.UpdatedAt)
		require.NoError(t, err, "Actual err: %v", err)
		tx.Close(t.Context())

		actual := loadPlanetFromDb(t, conn, planet.Id)
		assert.Equal(t, &actual, returned.Planet)
		require.NotNil(t, returned.Planet.BuildingAction)
		assert.Equal(t, storages, returned.Planet.BuildingAction.Storages)
	})

	t.Run("persists mutated planet with update to existing action productions", func(t *testing.T) {
		conn, tx := newTestConnectionAndTransaction(t)

		planet, _, _ := insertTestPlanetForPlayer(t, conn)
		action := insertTestBuildingActionForPlanet(t, conn, planet.Id, addBuildingActionProduction)
		require.Equal(t, crystalResourceId, action.Productions[0].Resource)
		require.NotEqual(t, 32, action.Productions[0].Production)
		planet.BuildingAction = &action

		productions := []models.BuildingActionResourceProduction{{Resource: crystalResourceId, Production: 32}}

		resolver := generatePlanetSimulator(func(p *models.Planet) {
			p.BuildingAction.Productions = productions
			p.Version++
		})

		ctx := newPlanetStepContext(planet, false, resolver)

		returned, err := ctx.resolve(t.Context(), tx, planet.UpdatedAt)
		require.NoError(t, err, "Actual err: %v", err)
		tx.Close(t.Context())

		actual := loadPlanetFromDb(t, conn, planet.Id)
		assert.Equal(t, &actual, returned.Planet)
		require.NotNil(t, returned.Planet.BuildingAction)
		assert.Equal(t, productions, returned.Planet.BuildingAction.Productions)
	})

	t.Run("persists mutated planet with updated action and productions", func(t *testing.T) {
		conn, tx := newTestConnectionAndTransaction(t)

		planet, _, _ := insertTestPlanetForPlayer(t, conn)
		action := insertTestBuildingActionForPlanet(t, conn, planet.Id)
		planet.BuildingAction = &action

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

		resolver := generatePlanetSimulator(func(p *models.Planet) {
			p.BuildingAction.Productions = productions
			p.Version++
		})

		ctx := newPlanetStepContext(planet, false, resolver)

		returned, err := ctx.resolve(t.Context(), tx, planet.UpdatedAt)
		require.NoError(t, err, "Actual err: %v", err)
		tx.Close(t.Context())

		actual := loadPlanetFromDb(t, conn, planet.Id)
		assert.Equal(t, &actual, returned.Planet)
		require.NotNil(t, returned.Planet.BuildingAction)
		assert.Equal(t, productions, returned.Planet.BuildingAction.Productions)
	})

	t.Run("persists mutated planet with deleted action", func(t *testing.T) {
		conn, tx := newTestConnectionAndTransaction(t)

		action, planet := insertTestBuildingAction(t, conn)

		resolver := generatePlanetSimulator(func(p *models.Planet) {
			p.BuildingAction = nil
			p.Version++
		})

		ctx := newPlanetStepContext(planet, false, resolver)

		returned, err := ctx.resolve(t.Context(), tx, planet.UpdatedAt)
		require.NoError(t, err, "Actual err: %v", err)
		tx.Close(t.Context())

		assertBuildingActionDoesNotExist(t, conn, action.Id)
		require.NotNil(t, returned.Planet)
		assert.Nil(t, returned.Planet.BuildingAction)
	})

	t.Run("persists mutated planet with deleted action with costs", func(t *testing.T) {
		conn, tx := newTestConnectionAndTransaction(t)

		action, planet := insertTestBuildingAction(t, conn, addBuildingActionCost)

		resolver := generatePlanetSimulator(func(p *models.Planet) {
			p.BuildingAction = nil
			p.Version++
		})

		ctx := newPlanetStepContext(planet, false, resolver)

		returned, err := ctx.resolve(t.Context(), tx, planet.UpdatedAt)
		require.NoError(t, err, "Actual err: %v", err)
		tx.Close(t.Context())

		assertBuildingActionDoesNotExist(t, conn, action.Id)
		assertBuildingActionCostDoesNotExist(t, conn, action.Id)
		require.NotNil(t, returned.Planet)
		assert.Nil(t, returned.Planet.BuildingAction)
	})

	t.Run("persists mutated planet with deleted action with storages", func(t *testing.T) {
		conn, tx := newTestConnectionAndTransaction(t)

		action, planet := insertTestBuildingAction(t, conn, addBuildingActionStorage)

		resolver := generatePlanetSimulator(func(p *models.Planet) {
			p.BuildingAction = nil
			p.Version++
		})

		ctx := newPlanetStepContext(planet, false, resolver)

		returned, err := ctx.resolve(t.Context(), tx, planet.UpdatedAt)
		require.NoError(t, err, "Actual err: %v", err)
		tx.Close(t.Context())

		assertBuildingActionDoesNotExist(t, conn, action.Id)
		assertBuildingActionStorageDoesNotExist(t, conn, action.Id)
		require.NotNil(t, returned.Planet)
		assert.Nil(t, returned.Planet.BuildingAction)
	})

	t.Run("persists mutated planet with deleted action with productions", func(t *testing.T) {
		conn, tx := newTestConnectionAndTransaction(t)

		action, planet := insertTestBuildingAction(t, conn, addBuildingActionProduction)

		resolver := generatePlanetSimulator(func(p *models.Planet) {
			p.BuildingAction = nil
			p.Version++
		})

		ctx := newPlanetStepContext(planet, false, resolver)

		returned, err := ctx.resolve(t.Context(), tx, planet.UpdatedAt)
		require.NoError(t, err, "Actual err: %v", err)
		tx.Close(t.Context())

		assertBuildingActionDoesNotExist(t, conn, action.Id)
		assertBuildingActionProductionDoesNotExist(t, conn, action.Id)
		require.NotNil(t, returned.Planet)
		assert.Nil(t, returned.Planet.BuildingAction)
	})

	t.Run("persists mutated planet with new action", func(t *testing.T) {
		conn, tx := newTestConnectionAndTransaction(t)

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

		resolver := generatePlanetSimulator(func(p *models.Planet) {
			p.BuildingAction = &newAction
			p.Version++
		})

		ctx := newPlanetStepContext(planet, false, resolver)

		returned, err := ctx.resolve(t.Context(), tx, planet.UpdatedAt)
		require.NoError(t, err, "Actual err: %v", err)
		tx.Close(t.Context())

		actual := loadPlanetFromDb(t, conn, planet.Id)
		assert.Equal(t, &actual, returned.Planet)
		assert.Equal(t, &newAction, returned.Planet.BuildingAction)
	})

	t.Run("persists mutated planet with new ship action", func(t *testing.T) {
		conn, tx := newTestConnectionAndTransaction(t)

		planet, _, _ := insertTestPlanetForPlayer(t, conn)
		require.Empty(t, planet.ShipActions)

		action := models.ShipAction{
			Id:                 uuid.New(),
			Ship:               lightFighterId,
			Count:              3,
			CreatedAt:          someTime,
			NextCompletionAt:   someTime.Add(1 * time.Hour),
			UnitCompletionTime: 1 * time.Hour,
		}

		resolver := generatePlanetSimulator(func(p *models.Planet) {
			p.ShipActions = append(p.ShipActions, action)
			p.Version++
		})

		ctx := newPlanetStepContext(planet, false, resolver)

		returned, err := ctx.resolve(t.Context(), tx, planet.UpdatedAt)
		require.NoError(t, err, "Actual err: %v", err)
		tx.Close(t.Context())

		actual := loadPlanetFromDb(t, conn, planet.Id)
		assert.Equal(t, &actual, returned.Planet)
		assert.Equal(t, []models.ShipAction{action}, returned.Planet.ShipActions)
	})

	t.Run("persists mutated planet with updated ship action", func(t *testing.T) {
		conn, tx := newTestConnectionAndTransaction(t)

		action, planet := insertTestShipAction(t, conn)

		resolver := generatePlanetSimulator(func(p *models.Planet) {
			p.ShipActions[0].Count++
			p.Version++
		})

		ctx := newPlanetStepContext(planet, false, resolver)

		returned, err := ctx.resolve(t.Context(), tx, planet.UpdatedAt)
		require.NoError(t, err, "Actual err: %v", err)
		tx.Close(t.Context())

		require.NotNil(t, returned.Planet)

		actual := loadPlanetFromDb(t, conn, planet.Id)
		assert.Equal(t, &actual, returned.Planet)
		expected := action
		expected.Count++
		assert.Equal(t, []models.ShipAction{expected}, returned.Planet.ShipActions)
	})

	t.Run("persists mutated planet with deleted ship action", func(t *testing.T) {
		conn, tx := newTestConnectionAndTransaction(t)

		_, planet := insertTestShipAction(t, conn)

		resolver := generatePlanetSimulator(func(p *models.Planet) {
			p.ShipActions = nil
			p.Version++
		})

		ctx := newPlanetStepContext(planet, false, resolver)

		returned, err := ctx.resolve(t.Context(), tx, planet.UpdatedAt)
		require.NoError(t, err, "Actual err: %v", err)
		tx.Close(t.Context())

		assertShipActionDoesNotExist(t, conn, planet.Id)
		require.NotNil(t, returned.Planet)
		assert.Equal(t, []models.ShipAction{}, returned.Planet.ShipActions)
	})

	t.Run("returns error when mutator does not update version", func(t *testing.T) {
		conn, tx := newTestConnectionAndTransaction(t)

		planet, _, _ := insertTestPlanetForPlayer(t, conn)

		resolver := generatePlanetSimulator(func(p *models.Planet) {
			p.Fields++
		})

		ctx := newPlanetStepContext(planet, false, resolver)

		_, err := ctx.resolve(t.Context(), tx, planet.UpdatedAt)

		assert.ErrorIs(t, err, domainerrors.ErrMutationWithoutVersionBump, "Actual err: %v", err)

		actual := loadPlanetFromDb(t, conn, planet.Id)
		assert.Equal(t, planet, actual)
	})

	t.Run("returns error when new storage is added", func(t *testing.T) {
		conn, tx := newTestConnectionAndTransaction(t)

		planet, _, _ := insertTestPlanetForPlayer(t, conn)

		resolver := generatePlanetSimulator(func(p *models.Planet) {
			p.Storages = []models.PlanetResourceStorage{{Resource: crystalResourceId, Storage: 5478}}
			p.Version++
		})

		ctx := newPlanetStepContext(planet, false, resolver)

		_, err := ctx.resolve(t.Context(), tx, planet.UpdatedAt)

		assert.ErrorIs(t, err, domainerrors.ErrResourceNotFound, "Actual err: %v", err)
	})

	t.Run("returns error when new building is added", func(t *testing.T) {
		conn, tx := newTestConnectionAndTransaction(t)

		planet, _, _ := insertTestPlanetForPlayer(t, conn)

		resolver := generatePlanetSimulator(func(p *models.Planet) {
			p.Buildings = []models.PlanetBuilding{{Building: metalMineId, Level: 5}}
			p.Version++
		})

		ctx := newPlanetStepContext(planet, false, resolver)

		_, err := ctx.resolve(t.Context(), tx, planet.UpdatedAt)

		assert.ErrorIs(t, err, domainerrors.ErrBuildingNotFound, "Actual err: %v", err)
	})

	t.Run("returns error when new ship is added", func(t *testing.T) {
		conn, tx := newTestConnectionAndTransaction(t)

		planet, _, _ := insertTestPlanetForPlayer(t, conn)

		resolver := generatePlanetSimulator(func(p *models.Planet) {
			p.Ships = []models.PlanetShip{{Ship: lightFighterId, Count: 34}}
			p.Version++
		})

		ctx := newPlanetStepContext(planet, false, resolver)

		_, err := ctx.resolve(t.Context(), tx, planet.UpdatedAt)

		assert.ErrorIs(t, err, domainerrors.ErrShipNotFound, "Actual err: %v", err)
	})
}

func TestIT_SimulationStepContext_Resolve_SaveFleet(t *testing.T) {
	t.Run("returns mutated fleet", func(t *testing.T) {
		conn, tx := newTestConnectionAndTransaction(t)

		planet, _, _ := insertTestPlanetForPlayer(t, conn)
		fleet := insertTestFleet(t, conn, planet, addFleetShip)
		require.NotEqual(t, fleet.UpdatedAt, yetAnotherTime)

		resolver := generateFleetSimulator(func(f *models.Fleet) {
			f.UpdatedAt = yetAnotherTime
			f.Version++
		})

		event := models.FleetEventContext{
			Timestamp: fleet.Schedule.ArrivalAt,
			Fleet:     fleet,
			Planet:    nil,
		}
		ctx := newFleetStepContext(event, resolver)

		afterArrival := fleet.Schedule.ArrivalAt.Add(1 * time.Minute)
		returned, err := ctx.resolve(t.Context(), tx, afterArrival)
		require.NoError(t, err, "Actual err: %v", err)

		require.NotNil(t, returned.Fleet)
		assert.Equal(t, yetAnotherTime, returned.Fleet.UpdatedAt)
		assert.Equal(t, fleet.Version+1, returned.Fleet.Version)
	})

	t.Run("persists mutated fleet", func(t *testing.T) {
		conn, tx := newTestConnectionAndTransaction(t)

		planet, _, _ := insertTestPlanetForPlayer(t, conn)
		fleet := insertTestFleet(t, conn, planet, addFleetShip)
		require.NotEqual(t, fleet.UpdatedAt, yetAnotherTime)

		resolver := generateFleetSimulator(func(f *models.Fleet) {
			f.UpdatedAt = yetAnotherTime
			f.Version++
		})

		event := models.FleetEventContext{
			Timestamp: fleet.Schedule.ArrivalAt,
			Fleet:     fleet,
			Planet:    nil,
		}
		ctx := newFleetStepContext(event, resolver)

		afterArrival := fleet.Schedule.ArrivalAt.Add(1 * time.Minute)
		returned, err := ctx.resolve(t.Context(), tx, afterArrival)
		require.NoError(t, err, "Actual err: %v", err)
		tx.Close(t.Context())

		actual := loadFleetFromDb(t, conn, fleet.Id)
		assert.Equal(t, &actual, returned.Fleet)
		expected := fleet
		expected.UpdatedAt = yetAnotherTime
		expected.Version++
		assert.Equal(t, expected, actual)
	})

	t.Run("persists updated fleet ships", func(t *testing.T) {
		conn, tx := newTestConnectionAndTransaction(t)

		planet, _, _ := insertTestPlanetForPlayer(t, conn)
		fleet := insertTestFleet(t, conn, planet, addFleetShip)
		initialCount := fleet.Ships[0].Count

		resolver := generateFleetSimulator(func(f *models.Fleet) {
			f.Ships[0].Count--
			f.Version++
		})

		event := models.FleetEventContext{
			Timestamp: fleet.Schedule.ArrivalAt,
			Fleet:     fleet,
			Planet:    nil,
		}
		ctx := newFleetStepContext(event, resolver)

		afterArrival := fleet.Schedule.ArrivalAt.Add(1 * time.Minute)
		returned, err := ctx.resolve(t.Context(), tx, afterArrival)
		require.NoError(t, err, "Actual err: %v", err)
		tx.Close(t.Context())

		actual := loadFleetFromDb(t, conn, fleet.Id)
		expected := []models.FleetShip{
			{Ship: fleet.Ships[0].Ship, Count: initialCount - 1},
		}
		assert.Equal(t, expected, actual.Ships)
		assert.Equal(t, actual.Ships, returned.Fleet.Ships)
	})

	t.Run("persists deleted fleet ships", func(t *testing.T) {
		conn, tx := newTestConnectionAndTransaction(t)

		planet, _, _ := insertTestPlanetForPlayer(t, conn)
		fleet := insertTestFleet(t, conn, planet, addFleetShip)

		resolver := generateFleetSimulator(func(f *models.Fleet) {
			f.Ships = []models.FleetShip{}
			f.Version++
		})

		event := models.FleetEventContext{
			Timestamp: fleet.Schedule.ArrivalAt,
			Fleet:     fleet,
			Planet:    nil,
		}
		ctx := newFleetStepContext(event, resolver)

		afterArrival := fleet.Schedule.ArrivalAt.Add(1 * time.Minute)
		returned, err := ctx.resolve(t.Context(), tx, afterArrival)
		require.NoError(t, err, "Actual err: %v", err)
		tx.Close(t.Context())

		actual := loadFleetFromDb(t, conn, fleet.Id)
		assert.Equal(t, actual.Ships, returned.Fleet.Ships)
		assert.Equal(t, []models.FleetShip{}, actual.Ships)
	})

	t.Run("persists updated fleet status", func(t *testing.T) {
		conn, tx := newTestConnectionAndTransaction(t)

		planet, _, _ := insertTestPlanetForPlayer(t, conn)
		fleet := insertTestFleet(t, conn, planet, addFleetShip)

		resolver := generateFleetSimulator(func(f *models.Fleet) {
			f.Schedule.Status = models.StatusReturning
			f.Version++
		})

		event := models.FleetEventContext{
			Timestamp: fleet.Schedule.ArrivalAt,
			Fleet:     fleet,
			Planet:    nil,
		}
		ctx := newFleetStepContext(event, resolver)

		afterArrival := fleet.Schedule.ArrivalAt.Add(1 * time.Minute)
		returned, err := ctx.resolve(t.Context(), tx, afterArrival)
		require.NoError(t, err, "Actual err: %v", err)
		tx.Close(t.Context())

		actual := loadFleetFromDb(t, conn, fleet.Id)
		assert.Equal(t, actual.Schedule.Status, returned.Fleet.Schedule.Status)
		assert.Equal(t, models.StatusReturning, actual.Schedule.Status)
	})

	t.Run("persists updated fleet event", func(t *testing.T) {
		conn, tx := newTestConnectionAndTransaction(t)

		planet, _, _ := insertTestPlanetForPlayer(t, conn)
		fleet := insertTestFleet(t, conn, planet, addFleetShip)
		assertFleetEventTypeMatches(t, conn, fleet.Id, models.ArrivalEvent)

		resolver := generateFleetSimulator(func(f *models.Fleet) {
			f.Schedule.Status = models.StatusReturning
			f.Version++
		})

		event := models.FleetEventContext{
			Timestamp: fleet.Schedule.ArrivalAt,
			Fleet:     fleet,
			Planet:    nil,
		}
		ctx := newFleetStepContext(event, resolver)

		afterArrival := fleet.Schedule.ArrivalAt.Add(1 * time.Minute)
		returned, err := ctx.resolve(t.Context(), tx, afterArrival)
		require.NoError(t, err, "Actual err: %v", err)
		tx.Close(t.Context())

		assert.Equal(t, models.StatusReturning, returned.Fleet.Schedule.Status)
		assertFleetEventTypeMatches(t, conn, fleet.Id, models.ReturnEvent)
	})
}

func TestIT_SimulationStepContext_Resolve_CreateFleet(t *testing.T) {
	t.Run("returns created fleet", func(t *testing.T) {
		conn, tx := newTestConnectionAndTransaction(t)

		planet, _, _ := insertTestPlanetForPlayer(t, conn)
		fleet := generateSampleFleet(planet)

		resolver := generateFleetCreator(fleet)

		ctx := newPlanetStepContext(planet, false, resolver)
		returned, err := ctx.resolve(t.Context(), tx, planet.UpdatedAt)
		require.NoError(t, err, "Actual err: %v", err)
		tx.Close(t.Context())

		assert.Equal(t, &fleet, returned.Fleet)
	})

	t.Run("returns updated planet", func(t *testing.T) {
		conn, tx := newTestConnectionAndTransaction(t)

		planet, _, _ := insertTestPlanetForPlayer(t, conn)
		fleet := generateSampleFleet(planet)

		resolver := generateFleetCreator(fleet)

		ctx := newPlanetStepContext(planet, false, resolver)
		returned, err := ctx.resolve(t.Context(), tx, planet.UpdatedAt)
		require.NoError(t, err, "Actual err: %v", err)
		tx.Close(t.Context())

		expected := planet
		expected.Version++
		assert.Equal(t, &expected, returned.Planet)
	})

	t.Run("persists created fleet", func(t *testing.T) {
		conn, tx := newTestConnectionAndTransaction(t)

		planet, _, _ := insertTestPlanetForPlayer(t, conn)
		fleet := generateSampleFleet(planet)

		resolver := generateFleetCreator(fleet)

		ctx := newPlanetStepContext(planet, false, resolver)
		_, err := ctx.resolve(t.Context(), tx, planet.UpdatedAt)
		require.NoError(t, err, "Actual err: %v", err)
		tx.Close(t.Context())

		actual := loadFleetFromDb(t, conn, fleet.Id)
		assert.Equal(t, fleet, actual)
	})

	t.Run("persists fleet event", func(t *testing.T) {
		conn, tx := newTestConnectionAndTransaction(t)

		planet, _, _ := insertTestPlanetForPlayer(t, conn)
		fleet := generateSampleFleet(planet)

		resolver := generateFleetCreator(fleet)

		ctx := newPlanetStepContext(planet, false, resolver)
		_, err := ctx.resolve(t.Context(), tx, planet.UpdatedAt)
		require.NoError(t, err, "Actual err: %v", err)
		tx.Close(t.Context())

		assertFleetEventExists(t, conn, fleet.Id, fleet.Schedule.ArrivalAt)
		assertFleetEventTypeMatches(t, conn, fleet.Id, models.ArrivalEvent)
	})

	t.Run("persists fleet with destination", func(t *testing.T) {
		conn, tx := newTestConnectionAndTransaction(t)

		planetA, planetB := insertTwoTestPlanets(t, conn)

		fleet := generateSampleFleet(planetA)
		fleet.Mission = models.MissionTransport
		fleet.Destination.Target = &planetB.Id

		resolver := generateFleetCreator(fleet)

		ctx := newPlanetStepContext(planetA, false, resolver)
		_, err := ctx.resolve(t.Context(), tx, planetA.UpdatedAt)
		require.NoError(t, err, "Actual err: %v", err)
		tx.Close(t.Context())

		actual := loadFleetFromDb(t, conn, fleet.Id)
		assert.Equal(t, fleet, actual)
	})

	t.Run("persists updated planet", func(t *testing.T) {
		conn, tx := newTestConnectionAndTransaction(t)

		planet, _, _ := insertTestPlanetForPlayer(t, conn, addPlanetShip)
		fleet := generateSampleFleet(planet)

		// Copying those as the slices in `newPlanetStepContext` are not copied.
		// Using them after the resolve function has been called would lead to
		// double increment.
		initialFields := planet.Fields
		initialVersion := planet.Version
		initialShipCount := planet.Ships[0].Count

		resolver := func(step models.SimulationStep) (models.SimulationStepResult, error) {
			step.Planet.Fields++
			step.Planet.Version++
			step.Planet.Ships[0].Count++

			result := models.SimulationStepResult{
				Resolution: models.FleetCreated,
				Planet:     step.Planet,
				Fleet:      &fleet,
			}

			return result, nil
		}

		ctx := newPlanetStepContext(planet, false, resolver)
		_, err := ctx.resolve(t.Context(), tx, planet.UpdatedAt)
		require.NoError(t, err, "Actual err: %v", err)
		tx.Close(t.Context())

		actual := loadPlanetFromDb(t, conn, planet.Id)
		expected := planet
		expected.Fields = initialFields + 1
		expected.Version = initialVersion + 1
		expected.Ships[0].Count = initialShipCount + 1
		assert.Equal(t, expected, actual)
	})

	t.Run("returns error when fleet source does not exist", func(t *testing.T) {
		conn, tx := newTestConnectionAndTransaction(t)

		planet, _, _ := insertTestPlanetForPlayer(t, conn)
		fleet := generateSampleFleet(planet)
		fleet.Source = uuid.New()

		resolver := func(step models.SimulationStep) (models.SimulationStepResult, error) {
			var planet *models.Planet
			if step.Planet != nil {
				planet = new(*step.Planet)
				planet.Version++
			}

			result := models.SimulationStepResult{
				Resolution: models.FleetCreated,
				Planet:     planet,
				Fleet:      &fleet,
			}

			return result, nil
		}

		ctx := newPlanetStepContext(planet, false, resolver)
		_, err := ctx.resolve(t.Context(), tx, planet.UpdatedAt)
		tx.Close(t.Context())

		assert.ErrorIs(t, err, domainerrors.ErrPlanetNotFound, "Actual err: %v", err)
		assertFleetDoesNotExist(t, conn, fleet.Id)
	})

	t.Run("returns error when fleet destination does not exist", func(t *testing.T) {
		conn, tx := newTestConnectionAndTransaction(t)

		planetA, _, _ := insertTestPlanetForPlayer(t, conn)

		fleet := generateSampleFleet(planetA)
		fleet.Mission = models.MissionTransport
		fakePlanetId := uuid.New()
		fleet.Destination.Target = &fakePlanetId

		resolver := generateFleetCreator(fleet)

		ctx := newPlanetStepContext(planetA, false, resolver)
		_, err := ctx.resolve(t.Context(), tx, planetA.UpdatedAt)
		tx.Close(t.Context())

		assert.ErrorIs(t, err, domainerrors.ErrFleetDestinationInvalid, "Actual err: %v", err)
		assertFleetDoesNotExist(t, conn, fleet.Id)
	})

	t.Run("returns error when planet version is not bumped", func(t *testing.T) {
		conn, tx := newTestConnectionAndTransaction(t)

		planet, _, _ := insertTestPlanetForPlayer(t, conn, addPlanetShip)
		fleet := generateSampleFleet(planet)

		resolver := func(step models.SimulationStep) (models.SimulationStepResult, error) {
			result := models.SimulationStepResult{
				Resolution: models.FleetCreated,
				Planet:     step.Planet,
				Fleet:      &fleet,
			}

			return result, nil
		}

		ctx := newPlanetStepContext(planet, false, resolver)
		_, err := ctx.resolve(t.Context(), tx, planet.UpdatedAt)
		tx.Close(t.Context())

		assert.ErrorIs(t, err, domainerrors.ErrMutationWithoutVersionBump, "Actual err: %v", err)
		assertFleetDoesNotExist(t, conn, fleet.Id)
	})
}

func TestIT_SimulationStepContext_Resolve_DeleteFleet(t *testing.T) {
	t.Run("deletes fleet", func(t *testing.T) {
		conn, tx := newTestConnectionAndTransaction(t)

		planet, _, _ := insertTestPlanetForPlayer(t, conn)
		fleet := insertTestFleet(t, conn, planet)

		resolver := generateFleetDeleter()

		event := models.FleetEventContext{
			Timestamp: fleet.Schedule.ArrivalAt,
			Fleet:     fleet,
			Planet:    nil,
		}
		ctx := newFleetStepContext(event, resolver)

		afterArrival := fleet.Schedule.ArrivalAt.Add(1 * time.Minute)
		returned, err := ctx.resolve(t.Context(), tx, afterArrival)
		require.NoError(t, err, "Actual err: %v", err)
		tx.Close(t.Context())

		require.Nil(t, returned.Fleet)
		assertFleetDoesNotExist(t, conn, fleet.Id)
	})

	t.Run("deletes fleet with ships", func(t *testing.T) {
		conn, tx := newTestConnectionAndTransaction(t)

		planet, _, _ := insertTestPlanetForPlayer(t, conn)
		fleet := insertTestFleet(t, conn, planet, addFleetShip)

		resolver := generateFleetDeleter()

		event := models.FleetEventContext{
			Timestamp: fleet.Schedule.ArrivalAt,
			Fleet:     fleet,
			Planet:    nil,
		}
		ctx := newFleetStepContext(event, resolver)

		afterArrival := fleet.Schedule.ArrivalAt.Add(1 * time.Minute)
		returned, err := ctx.resolve(t.Context(), tx, afterArrival)
		require.NoError(t, err, "Actual err: %v", err)
		tx.Close(t.Context())

		require.Nil(t, returned.Fleet)
		assertFleetDoesNotExist(t, conn, fleet.Id)
	})

	t.Run("deletes fleet with target", func(t *testing.T) {
		conn, tx := newTestConnectionAndTransaction(t)

		planet, destination := insertTwoTestPlanets(t, conn)
		fleet := insertTestFleet(t, conn, planet, addFleetShip)
		fleet.Destination.Target = &destination.Id
		upsertFleetDestination(t, conn, fleet)

		resolver := generateFleetDeleter()

		event := models.FleetEventContext{
			Timestamp: fleet.Schedule.ArrivalAt,
			Fleet:     fleet,
			Planet:    nil,
		}
		ctx := newFleetStepContext(event, resolver)

		afterArrival := fleet.Schedule.ArrivalAt.Add(1 * time.Minute)
		returned, err := ctx.resolve(t.Context(), tx, afterArrival)
		require.NoError(t, err, "Actual err: %v", err)
		tx.Close(t.Context())

		require.Nil(t, returned.Fleet)
		assertFleetDoesNotExist(t, conn, fleet.Id)
	})
}

func newTestConnectionAndTransaction(t *testing.T) (database.Connection, database.Transaction) {
	t.Helper()

	conn := newTestConnection(t)
	tx, err := conn.BeginTx(t.Context())
	require.NoError(t, err, "Actual err: %v", err)

	t.Cleanup(func() {
		tx.Close(t.Context())
	})

	return conn, tx
}

func generateCapturingResolver(result models.SimulationStepResult) (models.SimulationStepResolver, *models.SimulationStep) {
	var captor models.SimulationStep

	resolver := func(step models.SimulationStep) (models.SimulationStepResult, error) {
		captor = step
		return result, nil
	}

	return resolver, &captor
}

func dummyResolver(_ models.SimulationStep) (models.SimulationStepResult, error) {
	return models.SimulationStepResult{
		Resolution: models.PlanetUpdated,
		Planet:     &models.Planet{},
		Fleet:      nil,
	}, nil
}
