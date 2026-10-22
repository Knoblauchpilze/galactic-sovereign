package drivenadapters

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/adapters/driven/database"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	domainerrors "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models/errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIT_SimulationStepExecutor_Execute_LoadPlanetBehavior(t *testing.T) {
	testCases := []struct {
		name      string
		generator func(t *testing.T, conn database.Connection) models.Planet
	}{
		{
			name: "planet",
			generator: func(t *testing.T, conn database.Connection) models.Planet {
				planet, _, _ := insertTestPlanetForPlayer(t, conn)
				return planet
			},
		},
		{
			name: "planet with resources",
			generator: func(t *testing.T, conn database.Connection) models.Planet {
				planet, _, _ := insertTestPlanetForPlayer(t, conn, addPlanetResource)
				return planet
			},
		},
		{
			name: "planet with resource productions",
			generator: func(t *testing.T, conn database.Connection) models.Planet {
				planet, _, _ := insertTestPlanetForPlayer(t, conn, addPlanetProduction)
				return planet
			},
		},
		{
			name: "planet with resource productions for building",
			generator: func(t *testing.T, conn database.Connection) models.Planet {
				planet, _, _ := insertTestPlanetForPlayer(t, conn, addPlanetProductionForBuilding)
				return planet
			},
		},
		{
			name: "planet with resource storages",
			generator: func(t *testing.T, conn database.Connection) models.Planet {
				planet, _, _ := insertTestPlanetForPlayer(t, conn, addPlanetStorage)
				return planet
			},
		},
		{
			name: "planet with buildings",
			generator: func(t *testing.T, conn database.Connection) models.Planet {
				planet, _, _ := insertTestPlanetForPlayer(t, conn, addPlanetBuilding)
				return planet
			},
		},
		{
			name: "planet with buildings with speedup",
			generator: func(t *testing.T, conn database.Connection) models.Planet {
				planet, _, _ := insertTestPlanetForPlayer(t, conn, addPlanetBuildingWithShipSpeedup)
				return planet
			},
		},
		{
			name: "planet with ships",
			generator: func(t *testing.T, conn database.Connection) models.Planet {
				planet, _, _ := insertTestPlanetForPlayer(t, conn, addPlanetShip)
				return planet
			},
		},
		{
			name: "planet with building action",
			generator: func(t *testing.T, conn database.Connection) models.Planet {
				planet, _, _ := insertTestPlanetForPlayer(t, conn, addPlanetBuildingAction)
				return planet
			},
		},
		{
			name: "planet with ship actions",
			generator: func(t *testing.T, conn database.Connection) models.Planet {
				planet, _, _ := insertTestPlanetForPlayer(t, conn, addPlanetShipAction)
				return planet
			},
		},
		{
			name: "planet with multiple ship actions are sorted",
			generator: func(t *testing.T, conn database.Connection) models.Planet {
				planet, _, _ := insertTestPlanetForPlayer(t, conn)
				action1 := insertTestShipActionForPlanet(t, conn, planet.Id, func(t *testing.T, a *models.ShipAction) {
					a.CreatedAt = planet.CreatedAt.Add(2 * time.Hour)
				})
				action2 := insertTestShipActionForPlanet(t, conn, planet.Id, func(t *testing.T, a *models.ShipAction) {
					a.CreatedAt = planet.CreatedAt.Add(1 * time.Hour)
				})

				planet.ShipActions = []models.ShipAction{action2, action1}

				return planet
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			adapter, conn := newTestSimulationStepExecutor(t)

			planet := tc.generator(t, conn)

			var captured models.SimulationStep
			simulate := func(step models.SimulationStep) (models.SimulationStepResult, error) {
				captured = step

				planet := new(*step.Planet)
				planet.Version++

				result := models.SimulationStepResult{
					Resolution: models.PlanetUpdated,
					Planet:     planet,
					Fleet:      nil,
				}

				return result, nil
			}

			actual, err := adapter.Execute(t.Context(), planet.Id, planet.UpdatedAt, simulate)
			require.NoError(t, err, "Actual err: %v", err)

			expectedStep := models.SimulationStep{
				Timestamp:      planet.UpdatedAt,
				Planet:         &planet,
				NextFleetEvent: nil,
			}
			assert.Equal(t, expectedStep, captured)
			expectedResult := models.SimulationStepResult{
				Resolution: models.PlanetUpdated,
				Planet:     &planet,
				Fleet:      nil,
			}
			expectedResult.Planet.Version++
			assert.Equal(t, expectedResult, actual)
		})
	}
}

func TestIT_SimulationStepExecutor_Execute_LoadFleetBehavior(t *testing.T) {
	t.Run("loads fleet event", func(t *testing.T) {
		adapter, conn := newTestSimulationStepExecutor(t)

		planet, _, _ := insertTestPlanetForPlayer(t, conn)
		fleet := insertTestFleet(t, conn, planet, addFleetShip)

		var captured models.SimulationStep
		simulate := func(step models.SimulationStep) (models.SimulationStepResult, error) {
			captured = step

			fleet := new(step.NextFleetEvent.Fleet)
			fleet.Version++

			result := models.SimulationStepResult{
				Resolution: models.FleetEventResolved,
				Planet:     nil,
				Fleet:      fleet,
			}

			return result, nil
		}

		afterArrival := fleet.Schedule.ArrivalAt.Add(1 * time.Minute)

		actual, err := adapter.Execute(t.Context(), planet.Id, afterArrival, simulate)
		require.NoError(t, err, "Actual err: %v", err)

		expectedStep := models.SimulationStep{
			Timestamp: afterArrival,
			Planet:    nil,
			NextFleetEvent: &models.FleetEventContext{
				Timestamp: fleet.Schedule.ArrivalAt,
				Fleet:     fleet,
				Planet:    nil,
			},
		}
		assert.Equal(t, expectedStep, captured)
		expectedResult := models.SimulationStepResult{
			Resolution: models.FleetEventResolved,
			Planet:     nil,
			Fleet:      &fleet,
		}
		expectedResult.Fleet.Version++
		assert.Equal(t, expectedResult, actual)
	})

	t.Run("loads fleet event when fleet has a destination", func(t *testing.T) {
		adapter, conn := newTestSimulationStepExecutor(t)

		planet, _, _ := insertTestPlanetForPlayer(t, conn)
		fleet := insertTestFleet(t, conn, planet, addFleetShip)

		destination, _, _ := insertTestPlanetForPlayer(t, conn)
		fleet.Destination.Target = &destination.Id
		upsertFleetDestination(t, conn, fleet)

		var captured models.SimulationStep
		simulate := func(step models.SimulationStep) (models.SimulationStepResult, error) {
			captured = step

			planet := new(*step.NextFleetEvent.Planet)
			planet.Version++

			fleet := new(step.NextFleetEvent.Fleet)
			fleet.Version++

			result := models.SimulationStepResult{
				Resolution: models.FleetEventResolved,
				Planet:     planet,
				Fleet:      fleet,
			}

			return result, nil
		}

		afterArrival := fleet.Schedule.ArrivalAt.Add(1 * time.Minute)

		actual, err := adapter.Execute(t.Context(), planet.Id, afterArrival, simulate)
		require.NoError(t, err, "Actual err: %v", err)

		expectedStep := models.SimulationStep{
			Timestamp: afterArrival,
			Planet:    nil,
			NextFleetEvent: &models.FleetEventContext{
				Timestamp: fleet.Schedule.ArrivalAt,
				Fleet:     fleet,
				Planet:    &destination,
			},
		}
		assert.Equal(t, expectedStep, captured)
		expectedResult := models.SimulationStepResult{
			Resolution: models.FleetEventResolved,
			Planet:     &destination,
			Fleet:      &fleet,
		}
		expectedResult.Planet.Version++
		expectedResult.Fleet.Version++
		assert.Equal(t, expectedResult, actual)
	})

	t.Run("loads fleet event for returning fleet", func(t *testing.T) {
		adapter, conn := newTestSimulationStepExecutor(t)

		planet, _, _ := insertTestPlanetForPlayer(t, conn)
		fleet := insertTestFleet(t, conn, planet, addFleetShip)
		fleet.UpdatedAt = fleet.Schedule.ArrivalAt
		fleet.Schedule.Status = models.StatusReturning
		upsertFleetSchedule(t, conn, fleet)

		var captured models.SimulationStep
		simulate := func(step models.SimulationStep) (models.SimulationStepResult, error) {
			captured = step

			planet := new(*step.NextFleetEvent.Planet)
			planet.Version++

			fleet := new(step.NextFleetEvent.Fleet)
			fleet.Version++

			result := models.SimulationStepResult{
				Resolution: models.FleetEventResolved,
				Planet:     planet,
				Fleet:      fleet,
			}

			return result, nil
		}

		afterReturn := fleet.Schedule.ReturnAt.Add(1 * time.Minute)

		actual, err := adapter.Execute(t.Context(), planet.Id, afterReturn, simulate)
		require.NoError(t, err, "Actual err: %v", err)

		expectedStep := models.SimulationStep{
			Timestamp: afterReturn,
			Planet:    nil,
			NextFleetEvent: &models.FleetEventContext{
				Timestamp: fleet.Schedule.ReturnAt,
				Fleet:     fleet,
				Planet:    &planet,
			},
		}
		assert.Equal(t, expectedStep, captured)
		expectedResult := models.SimulationStepResult{
			Resolution: models.FleetEventResolved,
			Planet:     &planet,
			Fleet:      &fleet,
		}
		expectedResult.Planet.Version++
		expectedResult.Fleet.Version++
		assert.Equal(t, expectedResult, actual)
	})
}

func TestIT_SimulationStepExecutor_Execute_SavePlanet(t *testing.T) {
	t.Run("returns mutated planet", func(t *testing.T) {
		adapter, conn := newTestSimulationStepExecutor(t)

		planet, _, _ := insertTestPlanetForPlayer(t, conn)
		require.NotEqual(t, planet.UpdatedAt, yetAnotherTime)

		simulate := generatePlanetSimulator(func(p *models.Planet) {
			p.UpdatedAt = yetAnotherTime
			p.Version++
		})

		returned, err := adapter.Execute(t.Context(), planet.Id, planet.UpdatedAt, simulate)
		require.NoError(t, err, "Actual err: %v", err)

		require.NotNil(t, returned.Planet)
		assert.Equal(t, yetAnotherTime, returned.Planet.UpdatedAt)
		assert.Equal(t, planet.Version+1, returned.Planet.Version)
	})

	t.Run("persists mutated planet", func(t *testing.T) {
		adapter, conn := newTestSimulationStepExecutor(t)

		planet, _, _ := insertTestPlanetForPlayer(t, conn)
		require.NotEqual(t, planet.UpdatedAt, yetAnotherTime)
		require.NotEqual(t, 326, planet.Fields)

		simulate := generatePlanetSimulator(func(p *models.Planet) {
			p.Fields = 326
			p.UpdatedAt = yetAnotherTime
			p.Version++
		})

		_, err := adapter.Execute(t.Context(), planet.Id, planet.UpdatedAt, simulate)
		require.NoError(t, err, "Actual err: %v", err)

		actual := loadPlanetFromDb(t, conn, planet.Id)
		expected := planet
		expected.Fields = 326
		expected.UpdatedAt = yetAnotherTime
		expected.Version++
		assert.Equal(t, expected, actual)
	})

	t.Run("persists mutated planet resources", func(t *testing.T) {
		adapter, conn := newTestSimulationStepExecutor(t)

		planet, _, _ := insertTestPlanetForPlayer(t, conn, addPlanetResource)
		require.NotEqual(t, 5874, planet.Resources[0].Amount)

		simulate := generatePlanetSimulator(func(p *models.Planet) {
			p.Resources[0].Amount = 5874
			p.Version++
		})

		returned, err := adapter.Execute(t.Context(), planet.Id, planet.UpdatedAt, simulate)
		require.NoError(t, err, "Actual err: %v", err)

		require.NotNil(t, returned.Planet)
		expected := []models.PlanetResource{
			{Resource: crystalResourceId, Amount: 5874},
		}
		assert.Equal(t, expected, returned.Planet.Resources)
		assertPlanetResourceAmount(t, conn, planet.Id, crystalResourceId, 5874)
	})

	t.Run("does not delete existing planet resource", func(t *testing.T) {
		adapter, conn := newTestSimulationStepExecutor(t)

		planet, _, _ := insertTestPlanetForPlayer(t, conn, addPlanetResource)
		resource := planet.Resources[0].Resource
		amount := planet.Resources[0].Amount

		simulate := generatePlanetSimulator(func(p *models.Planet) {
			p.Resources = []models.PlanetResource{}
			p.Version++
		})

		returned, err := adapter.Execute(t.Context(), planet.Id, planet.UpdatedAt, simulate)
		require.NoError(t, err, "Actual err: %v", err)

		require.NotNil(t, returned.Planet)
		expected := []models.PlanetResource{
			{Resource: resource, Amount: amount},
		}
		assert.Equal(t, expected, returned.Planet.Resources)
		assertPlanetResourceAmount(t, conn, planet.Id, resource, amount)
	})

	t.Run("persists mutated planet productions", func(t *testing.T) {
		adapter, conn := newTestSimulationStepExecutor(t)

		planet, _, _ := insertTestPlanetForPlayer(t, conn, addPlanetProduction)
		require.NotEqual(t, 39841, planet.Productions[0].Production)

		simulate := generatePlanetSimulator(func(p *models.Planet) {
			p.Productions[0].Production = 39841
			p.Version++
		})

		returned, err := adapter.Execute(t.Context(), planet.Id, planet.UpdatedAt, simulate)
		require.NoError(t, err, "Actual err: %v", err)

		require.NotNil(t, returned.Planet)
		expected := []models.PlanetResourceProduction{
			{Resource: metalResourceId, Production: 39841},
		}
		assert.Equal(t, expected, returned.Planet.Productions)
		assertPlanetResourceProduction(t, conn, planet.Id, metalResourceId, nil, 39841)
	})

	t.Run("persists additional planet productions", func(t *testing.T) {
		adapter, conn := newTestSimulationStepExecutor(t)

		planet, _, _ := insertTestPlanetForPlayer(
			t, conn, addPlanetProduction, addPlanetProductionForBuilding,
		)
		require.NotEqual(t, crystalResourceId, planet.Productions[0].Building)
		initial0 := planet.Productions[0].Production
		initial1 := planet.Productions[1].Production

		simulate := generatePlanetSimulator(func(p *models.Planet) {
			prod := models.PlanetResourceProduction{
				Resource:   crystalResourceId,
				Production: 354789,
			}
			p.Productions = append(p.Productions, prod)
			p.Version++
		})

		returned, err := adapter.Execute(t.Context(), planet.Id, planet.UpdatedAt, simulate)
		require.NoError(t, err, "Actual err: %v", err)

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
		adapter, conn := newTestSimulationStepExecutor(t)

		planet, _, _ := insertTestPlanetForPlayer(t, conn, addPlanetProduction)

		simulate := generatePlanetSimulator(func(p *models.Planet) {
			p.Productions = []models.PlanetResourceProduction{}
			p.Version++
		})

		returned, err := adapter.Execute(t.Context(), planet.Id, planet.UpdatedAt, simulate)
		require.NoError(t, err, "Actual err: %v", err)

		require.NotNil(t, returned.Planet)
		assert.Equal(t, []models.PlanetResourceProduction{}, returned.Planet.Productions)
		assertPlanetProductionDoesNotExist(t, conn, planet.Id)
	})

	t.Run("persists mutated planet production for building", func(t *testing.T) {
		adapter, conn := newTestSimulationStepExecutor(t)

		planet, _, _ := insertTestPlanetForPlayer(t, conn, addPlanetProductionForBuilding)
		require.NotEqual(t, 1235, planet.Productions[0].Production)

		simulate := generatePlanetSimulator(func(p *models.Planet) {
			p.Productions[0].Production = 1235
			p.Version++
		})

		returned, err := adapter.Execute(t.Context(), planet.Id, planet.UpdatedAt, simulate)
		require.NoError(t, err, "Actual err: %v", err)

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
		adapter, conn := newTestSimulationStepExecutor(t)

		planet, _, _ := insertTestPlanetForPlayer(
			t, conn, addPlanetProduction, addPlanetProductionForBuilding,
		)
		require.NotEqual(t, metalMineId, planet.Productions[1].Building)
		initial0 := planet.Productions[0].Production
		initial1 := planet.Productions[1].Production

		simulate := generatePlanetSimulator(func(p *models.Planet) {
			prod := models.PlanetResourceProduction{
				Resource:   crystalResourceId,
				Building:   &metalMineId,
				Production: 354789,
			}
			p.Productions = append(p.Productions, prod)
			p.Version++
		})

		returned, err := adapter.Execute(t.Context(), planet.Id, planet.UpdatedAt, simulate)
		require.NoError(t, err, "Actual err: %v", err)

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
		adapter, conn := newTestSimulationStepExecutor(t)

		planet, _, _ := insertTestPlanetForPlayer(t, conn, addPlanetProductionForBuilding)

		simulate := generatePlanetSimulator(func(p *models.Planet) {
			p.Productions = []models.PlanetResourceProduction{}
			p.Version++
		})

		returned, err := adapter.Execute(t.Context(), planet.Id, planet.UpdatedAt, simulate)
		require.NoError(t, err, "Actual err: %v", err)

		require.NotNil(t, returned.Planet)
		assert.Equal(t, []models.PlanetResourceProduction{}, returned.Planet.Productions)
		assertPlanetProductionDoesNotExist(t, conn, planet.Id)
	})

	t.Run("persists mutated planet storages", func(t *testing.T) {
		adapter, conn := newTestSimulationStepExecutor(t)

		planet, _, _ := insertTestPlanetForPlayer(t, conn, addPlanetStorage)
		require.NotEqual(t, 4598, planet.Storages[0].Storage)

		simulate := generatePlanetSimulator(func(p *models.Planet) {
			p.Storages[0].Storage = 4598
			p.Version++
		})

		returned, err := adapter.Execute(t.Context(), planet.Id, planet.UpdatedAt, simulate)
		require.NoError(t, err, "Actual err: %v", err)

		require.NotNil(t, returned.Planet)
		expected := []models.PlanetResourceStorage{
			{Resource: crystalResourceId, Storage: 4598},
		}
		assert.Equal(t, expected, returned.Planet.Storages)
		assertPlanetResourceStorage(t, conn, planet.Id, crystalResourceId, 4598)
	})

	t.Run("does not delete existing planet storage", func(t *testing.T) {
		adapter, conn := newTestSimulationStepExecutor(t)

		planet, _, _ := insertTestPlanetForPlayer(t, conn, addPlanetStorage)
		resource := planet.Storages[0].Resource
		storage := planet.Storages[0].Storage

		simulate := generatePlanetSimulator(func(p *models.Planet) {
			p.Storages = []models.PlanetResourceStorage{}
			p.Version++
		})

		returned, err := adapter.Execute(t.Context(), planet.Id, planet.UpdatedAt, simulate)
		require.NoError(t, err, "Actual err: %v", err)

		require.NotNil(t, returned.Planet)
		expected := []models.PlanetResourceStorage{
			{Resource: resource, Storage: storage},
		}
		assert.Equal(t, expected, returned.Planet.Storages)
		assertPlanetResourceStorage(t, conn, planet.Id, resource, storage)
	})

	t.Run("persists mutated planet buildings", func(t *testing.T) {
		adapter, conn := newTestSimulationStepExecutor(t)

		planet, _, _ := insertTestPlanetForPlayer(t, conn, addPlanetBuilding)
		require.NotEqual(t, 6, planet.Buildings[0].Level)

		simulate := generatePlanetSimulator(func(p *models.Planet) {
			p.Buildings[0].Level = 6
			p.Version++
		})

		returned, err := adapter.Execute(t.Context(), planet.Id, planet.UpdatedAt, simulate)
		require.NoError(t, err, "Actual err: %v", err)

		require.NotNil(t, returned.Planet)
		expected := []models.PlanetBuilding{
			{Building: metalStorageId, Level: 6},
		}
		assert.Equal(t, expected, returned.Planet.Buildings)
		assertPlanetBuildingLevel(t, conn, planet.Id, metalStorageId, 6)
	})

	t.Run("does not delete existing planet building", func(t *testing.T) {
		adapter, conn := newTestSimulationStepExecutor(t)

		planet, _, _ := insertTestPlanetForPlayer(t, conn, addPlanetBuilding)
		building := planet.Buildings[0].Building
		level := planet.Buildings[0].Level

		simulate := generatePlanetSimulator(func(p *models.Planet) {
			p.Buildings = []models.PlanetBuilding{}
			p.Version++
		})

		returned, err := adapter.Execute(t.Context(), planet.Id, planet.UpdatedAt, simulate)
		require.NoError(t, err, "Actual err: %v", err)

		require.NotNil(t, returned.Planet)
		expected := []models.PlanetBuilding{
			{Building: building, Level: level},
		}
		assert.Equal(t, expected, returned.Planet.Buildings)
		assertPlanetBuildingLevel(t, conn, planet.Id, building, level)
	})

	t.Run("persists mutated planet ships", func(t *testing.T) {
		adapter, conn := newTestSimulationStepExecutor(t)

		planet, _, _ := insertTestPlanetForPlayer(t, conn, addPlanetShip)
		require.NotEqual(t, 776, planet.Ships[0].Count)

		simulate := generatePlanetSimulator(func(p *models.Planet) {
			p.Ships[0].Count = 776
			p.Version++
		})

		returned, err := adapter.Execute(t.Context(), planet.Id, planet.UpdatedAt, simulate)
		require.NoError(t, err, "Actual err: %v", err)

		require.NotNil(t, returned.Planet)
		expected := []models.PlanetShip{
			{Ship: lightFighterId, Count: 776},
		}
		assert.Equal(t, expected, returned.Planet.Ships)
		assertPlanetShipCount(t, conn, planet.Id, lightFighterId, 776)
	})

	t.Run("does not delete existing planet ship", func(t *testing.T) {
		adapter, conn := newTestSimulationStepExecutor(t)

		planet, _, _ := insertTestPlanetForPlayer(t, conn, addPlanetShip)
		ship := planet.Ships[0].Ship
		count := planet.Ships[0].Count

		simulate := generatePlanetSimulator(func(p *models.Planet) {
			p.Ships = []models.PlanetShip{}
			p.Version++
		})

		returned, err := adapter.Execute(t.Context(), planet.Id, planet.UpdatedAt, simulate)
		require.NoError(t, err, "Actual err: %v", err)

		require.NotNil(t, returned.Planet)
		expected := []models.PlanetShip{
			{Ship: ship, Count: count},
		}
		assert.Equal(t, expected, returned.Planet.Ships)
		assertPlanetShipCount(t, conn, planet.Id, ship, count)
	})

	t.Run("persists mutated planet with action", func(t *testing.T) {
		adapter, conn := newTestSimulationStepExecutor(t)

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

		simulate := generatePlanetSimulator(func(p *models.Planet) {
			p.BuildingAction = &action
			p.Version++
		})

		returned, err := adapter.Execute(t.Context(), planet.Id, planet.UpdatedAt, simulate)
		require.NoError(t, err, "Actual err: %v", err)

		actual := loadPlanetFromDb(t, conn, planet.Id)
		assert.Equal(t, returned.Planet, &actual)
		require.NotNil(t, returned.Planet.BuildingAction)
		assert.Equal(t, action, *returned.Planet.BuildingAction)
		require.NotNil(t, actual.BuildingAction)
		assert.Equal(t, action, *actual.BuildingAction)
	})

	t.Run("persists mutated planet with action and costs", func(t *testing.T) {
		adapter, conn := newTestSimulationStepExecutor(t)

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

		simulate := generatePlanetSimulator(func(p *models.Planet) {
			p.BuildingAction = &action
			p.Version++
		})

		returned, err := adapter.Execute(t.Context(), planet.Id, planet.UpdatedAt, simulate)
		require.NoError(t, err, "Actual err: %v", err)

		actual := loadPlanetFromDb(t, conn, planet.Id)
		assert.Equal(t, returned.Planet, &actual)
		require.NotNil(t, returned.Planet.BuildingAction)
		assert.Equal(t, action, *returned.Planet.BuildingAction)
		require.NotNil(t, actual.BuildingAction)
		assert.Equal(t, action, *actual.BuildingAction)
	})

	t.Run("persists mutated planet with action and storages", func(t *testing.T) {
		adapter, conn := newTestSimulationStepExecutor(t)

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

		simulate := generatePlanetSimulator(func(p *models.Planet) {
			p.BuildingAction = &action
			p.Version++
		})

		returned, err := adapter.Execute(t.Context(), planet.Id, planet.UpdatedAt, simulate)
		require.NoError(t, err, "Actual err: %v", err)

		actual := loadPlanetFromDb(t, conn, planet.Id)
		assert.Equal(t, returned.Planet, &actual)
		require.NotNil(t, returned.Planet.BuildingAction)
		assert.Equal(t, action, *returned.Planet.BuildingAction)
		require.NotNil(t, actual.BuildingAction)
		assert.Equal(t, action, *actual.BuildingAction)
	})

	t.Run("persists mutated planet with action and productions", func(t *testing.T) {
		adapter, conn := newTestSimulationStepExecutor(t)

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

		simulate := generatePlanetSimulator(func(p *models.Planet) {
			p.BuildingAction = &action
			p.Version++
		})

		returned, err := adapter.Execute(t.Context(), planet.Id, planet.UpdatedAt, simulate)
		require.NoError(t, err, "Actual err: %v", err)

		actual := loadPlanetFromDb(t, conn, planet.Id)
		assert.Equal(t, returned.Planet, &actual)
		require.NotNil(t, returned.Planet.BuildingAction)
		assert.Equal(t, action, *returned.Planet.BuildingAction)
		require.NotNil(t, actual.BuildingAction)
		assert.Equal(t, action, *actual.BuildingAction)
	})

	t.Run("persists mutated planet with updated completion time", func(t *testing.T) {
		adapter, conn := newTestSimulationStepExecutor(t)

		planet, _, _ := insertTestPlanetForPlayer(t, conn)
		action := insertTestBuildingActionForPlanet(t, conn, planet.Id)
		require.NotEqual(t, yetAnotherTime, action.CompletedAt)

		simulate := generatePlanetSimulator(func(p *models.Planet) {
			p.BuildingAction.CompletedAt = yetAnotherTime
			p.Version++
		})

		returned, err := adapter.Execute(t.Context(), planet.Id, planet.UpdatedAt, simulate)
		require.NoError(t, err, "Actual err: %v", err)

		actual := loadPlanetFromDb(t, conn, planet.Id)
		assert.Equal(t, returned.Planet, &actual)
		require.NotNil(t, returned.Planet.BuildingAction)
		assert.Equal(t, yetAnotherTime, returned.Planet.BuildingAction.CompletedAt)
		require.NotNil(t, actual.BuildingAction)
		assert.Equal(t, yetAnotherTime, actual.BuildingAction.CompletedAt)
	})

	t.Run("persists mutated planet with update to existing action costs", func(t *testing.T) {
		adapter, conn := newTestSimulationStepExecutor(t)

		planet, _, _ := insertTestPlanetForPlayer(t, conn)
		action := insertTestBuildingActionForPlanet(t, conn, planet.Id, addBuildingActionCost)
		require.Equal(t, metalResourceId, action.Costs[0].Resource)
		require.NotEqual(t, 32, action.Costs[0].Amount)

		costs := []models.BuildingActionCost{{Resource: metalResourceId, Amount: 32}}

		simulate := generatePlanetSimulator(func(p *models.Planet) {
			p.BuildingAction.Costs = costs
			p.Version++
		})

		returned, err := adapter.Execute(t.Context(), planet.Id, planet.UpdatedAt, simulate)
		require.NoError(t, err, "Actual err: %v", err)

		actual := loadPlanetFromDb(t, conn, planet.Id)
		require.NotNil(t, returned.Planet)
		require.NotNil(t, returned.Planet.BuildingAction)
		assert.Equal(t, costs, returned.Planet.BuildingAction.Costs)
		require.NotNil(t, actual.BuildingAction)
		assert.Equal(t, costs, actual.BuildingAction.Costs)
	})

	t.Run("persists mutated planet with updated action and costs", func(t *testing.T) {
		adapter, conn := newTestSimulationStepExecutor(t)

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

		simulate := generatePlanetSimulator(func(p *models.Planet) {
			p.BuildingAction.Costs = costs
			p.Version++
		})

		returned, err := adapter.Execute(t.Context(), planet.Id, planet.UpdatedAt, simulate)
		require.NoError(t, err, "Actual err: %v", err)

		actual := loadPlanetFromDb(t, conn, planet.Id)
		assert.Equal(t, returned.Planet, &actual)
		require.NotNil(t, returned.Planet.BuildingAction)
		assert.Equal(t, costs, returned.Planet.BuildingAction.Costs)
		require.NotNil(t, actual.BuildingAction)
		assert.Equal(t, costs, actual.BuildingAction.Costs)
	})

	t.Run("persists mutated planet with update to existing action storages", func(t *testing.T) {
		adapter, conn := newTestSimulationStepExecutor(t)

		planet, _, _ := insertTestPlanetForPlayer(t, conn)
		action := insertTestBuildingActionForPlanet(t, conn, planet.Id, addBuildingActionStorage)
		require.Equal(t, crystalResourceId, action.Storages[0].Resource)
		require.NotEqual(t, 32, action.Storages[0].Storage)

		storages := []models.BuildingActionResourceStorage{{Resource: crystalResourceId, Storage: 32}}

		simulate := generatePlanetSimulator(func(p *models.Planet) {
			p.BuildingAction.Storages = storages
			p.Version++
		})

		returned, err := adapter.Execute(t.Context(), planet.Id, planet.UpdatedAt, simulate)
		require.NoError(t, err, "Actual err: %v", err)

		actual := loadPlanetFromDb(t, conn, planet.Id)
		require.NotNil(t, returned.Planet)
		require.NotNil(t, returned.Planet.BuildingAction)
		assert.Equal(t, storages, returned.Planet.BuildingAction.Storages)
		require.NotNil(t, actual.BuildingAction)
		assert.Equal(t, storages, actual.BuildingAction.Storages)
	})

	t.Run("persists mutated planet with updated action and storages", func(t *testing.T) {
		adapter, conn := newTestSimulationStepExecutor(t)

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

		simulate := generatePlanetSimulator(func(p *models.Planet) {
			p.BuildingAction.Storages = storages
			p.Version++
		})

		returned, err := adapter.Execute(t.Context(), planet.Id, planet.UpdatedAt, simulate)
		require.NoError(t, err, "Actual err: %v", err)

		actual := loadPlanetFromDb(t, conn, planet.Id)
		assert.Equal(t, returned.Planet, &actual)
		require.NotNil(t, returned.Planet.BuildingAction)
		assert.Equal(t, storages, returned.Planet.BuildingAction.Storages)
		require.NotNil(t, actual.BuildingAction)
		assert.Equal(t, storages, actual.BuildingAction.Storages)
	})

	t.Run("persists mutated planet with update to existing action productions", func(t *testing.T) {
		adapter, conn := newTestSimulationStepExecutor(t)

		planet, _, _ := insertTestPlanetForPlayer(t, conn)
		action := insertTestBuildingActionForPlanet(t, conn, planet.Id, addBuildingActionProduction)
		require.Equal(t, crystalResourceId, action.Productions[0].Resource)
		require.NotEqual(t, 32, action.Productions[0].Production)

		productions := []models.BuildingActionResourceProduction{{Resource: crystalResourceId, Production: 32}}

		simulate := generatePlanetSimulator(func(p *models.Planet) {
			p.BuildingAction.Productions = productions
			p.Version++
		})

		returned, err := adapter.Execute(t.Context(), planet.Id, planet.UpdatedAt, simulate)
		require.NoError(t, err, "Actual err: %v", err)

		actual := loadPlanetFromDb(t, conn, planet.Id)
		require.NotNil(t, returned.Planet)
		require.NotNil(t, returned.Planet.BuildingAction)
		assert.Equal(t, productions, returned.Planet.BuildingAction.Productions)
		require.NotNil(t, actual.BuildingAction)
		assert.Equal(t, productions, actual.BuildingAction.Productions)
	})

	t.Run("persists mutated planet with updated action and productions", func(t *testing.T) {
		adapter, conn := newTestSimulationStepExecutor(t)

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

		simulate := generatePlanetSimulator(func(p *models.Planet) {
			p.BuildingAction.Productions = productions
			p.Version++
		})

		returned, err := adapter.Execute(t.Context(), planet.Id, planet.UpdatedAt, simulate)
		require.NoError(t, err, "Actual err: %v", err)

		actual := loadPlanetFromDb(t, conn, planet.Id)
		assert.Equal(t, returned.Planet, &actual)
		require.NotNil(t, returned.Planet.BuildingAction)
		assert.Equal(t, productions, returned.Planet.BuildingAction.Productions)
		require.NotNil(t, actual.BuildingAction)
		assert.Equal(t, productions, actual.BuildingAction.Productions)
	})

	t.Run("persists mutated planet with deleted action", func(t *testing.T) {
		adapter, conn := newTestSimulationStepExecutor(t)

		action, planet := insertTestBuildingAction(t, conn)

		simulate := generatePlanetSimulator(func(p *models.Planet) {
			p.BuildingAction = nil
			p.Version++
		})

		returned, err := adapter.Execute(t.Context(), planet.Id, planet.UpdatedAt, simulate)
		require.NoError(t, err, "Actual err: %v", err)

		assertBuildingActionDoesNotExist(t, conn, action.Id)
		require.NotNil(t, returned.Planet)
		assert.Nil(t, returned.Planet.BuildingAction)
	})

	t.Run("persists mutated planet with deleted action with costs", func(t *testing.T) {
		adapter, conn := newTestSimulationStepExecutor(t)

		action, planet := insertTestBuildingAction(t, conn, addBuildingActionCost)

		simulate := generatePlanetSimulator(func(p *models.Planet) {
			p.BuildingAction = nil
			p.Version++
		})

		returned, err := adapter.Execute(t.Context(), planet.Id, planet.UpdatedAt, simulate)
		require.NoError(t, err, "Actual err: %v", err)

		assertBuildingActionDoesNotExist(t, conn, action.Id)
		assertBuildingActionCostDoesNotExist(t, conn, action.Id)
		require.NotNil(t, returned.Planet)
		assert.Nil(t, returned.Planet.BuildingAction)
	})

	t.Run("persists mutated planet with deleted action with storages", func(t *testing.T) {
		adapter, conn := newTestSimulationStepExecutor(t)

		action, planet := insertTestBuildingAction(t, conn, addBuildingActionStorage)

		simulate := generatePlanetSimulator(func(p *models.Planet) {
			p.BuildingAction = nil
			p.Version++
		})

		returned, err := adapter.Execute(t.Context(), planet.Id, planet.UpdatedAt, simulate)
		require.NoError(t, err, "Actual err: %v", err)

		assertBuildingActionDoesNotExist(t, conn, action.Id)
		assertBuildingActionStorageDoesNotExist(t, conn, action.Id)
		require.NotNil(t, returned.Planet)
		assert.Nil(t, returned.Planet.BuildingAction)
	})

	t.Run("persists mutated planet with deleted action with productions", func(t *testing.T) {
		adapter, conn := newTestSimulationStepExecutor(t)

		action, planet := insertTestBuildingAction(t, conn, addBuildingActionProduction)

		simulate := generatePlanetSimulator(func(p *models.Planet) {
			p.BuildingAction = nil
			p.Version++
		})

		returned, err := adapter.Execute(t.Context(), planet.Id, planet.UpdatedAt, simulate)
		require.NoError(t, err, "Actual err: %v", err)

		assertBuildingActionDoesNotExist(t, conn, action.Id)
		assertBuildingActionProductionDoesNotExist(t, conn, action.Id)
		require.NotNil(t, returned.Planet)
		assert.Nil(t, returned.Planet.BuildingAction)
	})

	t.Run("persists mutated planet with new action", func(t *testing.T) {
		adapter, conn := newTestSimulationStepExecutor(t)

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

		simulate := generatePlanetSimulator(func(p *models.Planet) {
			p.BuildingAction = &newAction
			p.Version++
		})

		returned, err := adapter.Execute(t.Context(), planet.Id, planet.UpdatedAt, simulate)
		require.NoError(t, err, "Actual err: %v", err)

		actual := loadPlanetFromDb(t, conn, planet.Id)
		assert.Equal(t, returned.Planet, &actual)
		require.NotNil(t, returned.Planet.BuildingAction)
		assert.Equal(t, newAction, *returned.Planet.BuildingAction)
		require.NotNil(t, actual.BuildingAction)
		assert.Equal(t, newAction, *actual.BuildingAction)
	})

	t.Run("persists mutated planet with new ship action", func(t *testing.T) {
		adapter, conn := newTestSimulationStepExecutor(t)

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

		simulate := generatePlanetSimulator(func(p *models.Planet) {
			p.ShipActions = append(p.ShipActions, action)
			p.Version++
		})

		returned, err := adapter.Execute(t.Context(), planet.Id, planet.UpdatedAt, simulate)
		require.NoError(t, err, "Actual err: %v", err)

		actual := loadPlanetFromDb(t, conn, planet.Id)
		assert.Equal(t, returned.Planet, &actual)
		require.NotEmpty(t, returned.Planet.ShipActions)
		assert.Equal(t, action, returned.Planet.ShipActions[0])
		require.NotEmpty(t, actual.ShipActions)
		assert.Equal(t, action, actual.ShipActions[0])
	})

	t.Run("persists mutated planet with updated completion time", func(t *testing.T) {
		adapter, conn := newTestSimulationStepExecutor(t)

		planet, _, _ := insertTestPlanetForPlayer(t, conn)
		action := insertTestBuildingActionForPlanet(t, conn, planet.Id)
		require.NotEqual(t, yetAnotherTime, action.CompletedAt)

		simulate := generatePlanetSimulator(func(p *models.Planet) {
			p.BuildingAction.CompletedAt = yetAnotherTime
			p.Version++
		})

		returned, err := adapter.Execute(t.Context(), planet.Id, planet.UpdatedAt, simulate)
		require.NoError(t, err, "Actual err: %v", err)

		actual := loadPlanetFromDb(t, conn, planet.Id)
		assert.Equal(t, returned.Planet, &actual)
		require.NotNil(t, returned.Planet.BuildingAction)
		assert.Equal(t, yetAnotherTime, returned.Planet.BuildingAction.CompletedAt)
		require.NotNil(t, actual.BuildingAction)
		assert.Equal(t, yetAnotherTime, actual.BuildingAction.CompletedAt)
	})

	t.Run("persists mutated planet with deleted ship action", func(t *testing.T) {
		adapter, conn := newTestSimulationStepExecutor(t)

		_, planet := insertTestShipAction(t, conn)

		simulate := generatePlanetSimulator(func(p *models.Planet) {
			p.ShipActions = nil
			p.Version++
		})

		returned, err := adapter.Execute(t.Context(), planet.Id, planet.UpdatedAt, simulate)
		require.NoError(t, err, "Actual err: %v", err)

		assertShipActionDoesNotExist(t, conn, planet.Id)
		require.NotNil(t, returned.Planet)
		assert.Equal(t, []models.ShipAction{}, returned.Planet.ShipActions)
	})

	t.Run("persists mutated planet with new ship action", func(t *testing.T) {
		adapter, conn := newTestSimulationStepExecutor(t)

		action, planet := insertTestShipAction(t, conn)
		require.NotEqual(t, smallCargoId, action.Ship)

		newAction := models.ShipAction{
			Id:                 uuid.New(),
			Ship:               smallCargoId,
			Count:              30,
			CreatedAt:          someTime,
			NextCompletionAt:   someTime.Add(2 * time.Minute),
			UnitCompletionTime: 2 * time.Minute,
		}

		simulate := generatePlanetSimulator(func(p *models.Planet) {
			p.ShipActions = []models.ShipAction{newAction}
			p.Version++
		})

		returned, err := adapter.Execute(t.Context(), planet.Id, planet.UpdatedAt, simulate)
		require.NoError(t, err, "Actual err: %v", err)

		actual := loadPlanetFromDb(t, conn, planet.Id)
		assert.Equal(t, returned.Planet, &actual)
		require.NotEmpty(t, returned.Planet.ShipActions)
		assert.Equal(t, newAction, returned.Planet.ShipActions[0])
		require.NotEmpty(t, actual.ShipActions)
		assert.Equal(t, newAction, actual.ShipActions[0])
	})

	t.Run("returns error when planet does not exist", func(t *testing.T) {
		adapter, _ := newTestSimulationStepExecutor(t)

		simulate := generatePlanetSimulator(func(p *models.Planet) {
			p.UpdatedAt = yetAnotherTime
			p.Version++
		})

		_, err := adapter.Execute(t.Context(), uuid.New(), time.Now(), simulate)

		assert.ErrorIs(t, err, domainerrors.ErrNotFound, "Actual err: %v", err)
	})

	t.Run("returns error when mutator does not update version", func(t *testing.T) {
		adapter, conn := newTestSimulationStepExecutor(t)

		planet, _, _ := insertTestPlanetForPlayer(t, conn)

		simulate := generatePlanetSimulator(func(*models.Planet) {})

		_, err := adapter.Execute(t.Context(), planet.Id, planet.UpdatedAt, simulate)

		assert.ErrorIs(t, err, domainerrors.ErrMutationWithoutVersionBump, "Actual err: %v", err)
	})

	t.Run("returns error when new storage is added", func(t *testing.T) {
		adapter, conn := newTestSimulationStepExecutor(t)

		planet, _, _ := insertTestPlanetForPlayer(t, conn)

		simulate := generatePlanetSimulator(func(p *models.Planet) {
			p.Storages = []models.PlanetResourceStorage{{Resource: crystalResourceId, Storage: 5478}}
			p.Version++
		})

		_, err := adapter.Execute(t.Context(), planet.Id, planet.UpdatedAt, simulate)

		assert.ErrorIs(t, err, domainerrors.ErrResourceNotFound, "Actual err: %v", err)
	})

	t.Run("returns error when new building is added", func(t *testing.T) {
		adapter, conn := newTestSimulationStepExecutor(t)

		planet, _, _ := insertTestPlanetForPlayer(t, conn)

		simulate := generatePlanetSimulator(func(p *models.Planet) {
			p.Buildings = []models.PlanetBuilding{{Building: metalMineId, Level: 5}}
			p.Version++
		})

		_, err := adapter.Execute(t.Context(), planet.Id, planet.UpdatedAt, simulate)

		assert.ErrorIs(t, err, domainerrors.ErrBuildingNotFound, "Actual err: %v", err)
	})

	t.Run("returns error when new ship is added", func(t *testing.T) {
		adapter, conn := newTestSimulationStepExecutor(t)

		planet, _, _ := insertTestPlanetForPlayer(t, conn)

		simulate := generatePlanetSimulator(func(p *models.Planet) {
			p.Ships = []models.PlanetShip{{Ship: lightFighterId, Count: 34}}
			p.Version++
		})

		_, err := adapter.Execute(t.Context(), planet.Id, planet.UpdatedAt, simulate)

		assert.ErrorIs(t, err, domainerrors.ErrShipNotFound, "Actual err: %v", err)
	})
}
func TestIT_SimulationStepExecutor_Execute_SaveFleet(t *testing.T) {
	t.Run("returns mutated fleet", func(t *testing.T) {
		adapter, conn := newTestSimulationStepExecutor(t)

		planet, _, _ := insertTestPlanetForPlayer(t, conn)
		fleet := insertTestFleet(t, conn, planet, addFleetShip)
		require.NotEqual(t, fleet.UpdatedAt, yetAnotherTime)

		simulate := generateFleetSimulator(func(f *models.Fleet) {
			f.UpdatedAt = yetAnotherTime
			f.Version++
		})

		afterArrival := fleet.Schedule.ArrivalAt.Add(1 * time.Minute)
		returned, err := adapter.Execute(t.Context(), planet.Id, afterArrival, simulate)
		require.NoError(t, err, "Actual err: %v", err)

		require.NotNil(t, returned.Fleet)
		assert.Equal(t, yetAnotherTime, returned.Fleet.UpdatedAt)
		assert.Equal(t, fleet.Version+1, returned.Fleet.Version)
	})

	t.Run("persists mutated fleet", func(t *testing.T) {
		adapter, conn := newTestSimulationStepExecutor(t)

		planet, _, _ := insertTestPlanetForPlayer(t, conn)
		fleet := insertTestFleet(t, conn, planet, addFleetShip)
		require.NotEqual(t, fleet.UpdatedAt, yetAnotherTime)

		simulate := generateFleetSimulator(func(f *models.Fleet) {
			f.UpdatedAt = yetAnotherTime
			f.Version++
		})

		afterArrival := fleet.Schedule.ArrivalAt.Add(1 * time.Minute)
		_, err := adapter.Execute(t.Context(), planet.Id, afterArrival, simulate)
		require.NoError(t, err, "Actual err: %v", err)

		actual := loadFleetFromDb(t, conn, fleet.Id)
		expected := fleet
		expected.UpdatedAt = yetAnotherTime
		expected.Version++
		assert.Equal(t, expected, actual)
	})

	t.Run("persists updated fleet ships", func(t *testing.T) {
		adapter, conn := newTestSimulationStepExecutor(t)

		planet, _, _ := insertTestPlanetForPlayer(t, conn)
		fleet := insertTestFleet(t, conn, planet, addFleetShip)

		simulate := generateFleetSimulator(func(f *models.Fleet) {
			f.Ships[0].Count--
			f.Version++
		})

		afterArrival := fleet.Schedule.ArrivalAt.Add(1 * time.Minute)
		_, err := adapter.Execute(t.Context(), planet.Id, afterArrival, simulate)
		require.NoError(t, err, "Actual err: %v", err)

		actual := loadFleetFromDb(t, conn, fleet.Id)
		expected := []models.FleetShip{
			{Ship: fleet.Ships[0].Ship, Count: fleet.Ships[0].Count - 1},
		}
		assert.Equal(t, expected, actual.Ships)
	})

	t.Run("persists deleted fleet ships", func(t *testing.T) {
		adapter, conn := newTestSimulationStepExecutor(t)

		planet, _, _ := insertTestPlanetForPlayer(t, conn)
		fleet := insertTestFleet(t, conn, planet, addFleetShip)

		simulate := generateFleetSimulator(func(f *models.Fleet) {
			f.Ships = []models.FleetShip{}
			f.Version++
		})

		afterArrival := fleet.Schedule.ArrivalAt.Add(1 * time.Minute)
		_, err := adapter.Execute(t.Context(), planet.Id, afterArrival, simulate)
		require.NoError(t, err, "Actual err: %v", err)

		actual := loadFleetFromDb(t, conn, fleet.Id)
		assert.Equal(t, []models.FleetShip{}, actual.Ships)
	})

	t.Run("persists updated fleet status", func(t *testing.T) {
		adapter, conn := newTestSimulationStepExecutor(t)

		planet, _, _ := insertTestPlanetForPlayer(t, conn)
		fleet := insertTestFleet(t, conn, planet, addFleetShip)

		simulate := generateFleetSimulator(func(f *models.Fleet) {
			f.Schedule.Status = models.StatusReturning
			f.Version++
		})

		afterArrival := fleet.Schedule.ArrivalAt.Add(1 * time.Minute)
		_, err := adapter.Execute(t.Context(), planet.Id, afterArrival, simulate)
		require.NoError(t, err, "Actual err: %v", err)

		actual := loadFleetFromDb(t, conn, fleet.Id)
		assert.Equal(t, models.StatusReturning, actual.Schedule.Status)
	})

	t.Run("persists updated fleet event", func(t *testing.T) {
		adapter, conn := newTestSimulationStepExecutor(t)

		planet, _, _ := insertTestPlanetForPlayer(t, conn)
		fleet := insertTestFleet(t, conn, planet, addFleetShip)
		assertFleetEventTypeMatches(t, conn, fleet.Id, models.ArrivalEvent)

		simulate := generateFleetSimulator(func(f *models.Fleet) {
			f.Schedule.Status = models.StatusReturning
			f.Version++
		})

		afterArrival := fleet.Schedule.ArrivalAt.Add(1 * time.Minute)
		_, err := adapter.Execute(t.Context(), planet.Id, afterArrival, simulate)
		require.NoError(t, err, "Actual err: %v", err)

		assertFleetEventTypeMatches(t, conn, fleet.Id, models.ReturnEvent)
	})
}

func TestIT_SimulationStepExecutor_Execute_DeleteFleet(t *testing.T) {
	t.Run("deletes fleet", func(t *testing.T) {
		adapter, conn := newTestSimulationStepExecutor(t)

		planet, _, _ := insertTestPlanetForPlayer(t, conn)
		fleet := insertTestFleet(t, conn, planet)

		simulate := generateFleetDeleter()

		afterArrival := fleet.Schedule.ArrivalAt.Add(1 * time.Minute)
		returned, err := adapter.Execute(t.Context(), planet.Id, afterArrival, simulate)
		require.NoError(t, err, "Actual err: %v", err)

		require.Nil(t, returned.Fleet)
		assertFleetDoesNotExist(t, conn, fleet.Id)
	})

	t.Run("deletes fleet with ships", func(t *testing.T) {
		adapter, conn := newTestSimulationStepExecutor(t)

		planet, _, _ := insertTestPlanetForPlayer(t, conn)
		fleet := insertTestFleet(t, conn, planet, addFleetShip)

		simulate := generateFleetDeleter()

		afterArrival := fleet.Schedule.ArrivalAt.Add(1 * time.Minute)
		returned, err := adapter.Execute(t.Context(), planet.Id, afterArrival, simulate)
		require.NoError(t, err, "Actual err: %v", err)

		require.Nil(t, returned.Fleet)
		assertFleetDoesNotExist(t, conn, fleet.Id)
	})

	t.Run("deletes fleet with target", func(t *testing.T) {
		adapter, conn := newTestSimulationStepExecutor(t)

		planet, _, _ := insertTestPlanetForPlayer(t, conn)
		fleet := insertTestFleet(t, conn, planet, addFleetShip)

		destination, _, _ := insertTestPlanetForPlayer(t, conn)
		fleet.Destination.Target = &destination.Id
		upsertFleetDestination(t, conn, fleet)

		simulate := generateFleetDeleter()

		afterArrival := fleet.Schedule.ArrivalAt.Add(1 * time.Minute)
		returned, err := adapter.Execute(t.Context(), planet.Id, afterArrival, simulate)
		require.NoError(t, err, "Actual err: %v", err)

		require.Nil(t, returned.Fleet)
		assertFleetDoesNotExist(t, conn, fleet.Id)
	})
}

func TestIT_SimulationStepExecutor_Execute_Concurrency(t *testing.T) {
	t.Run("prevents concurrent modification of different planets when fleet event is resolved", func(t *testing.T) {
		adapter, conn := newTestSimulationStepExecutor(t)

		planetA, _, _ := insertTestPlanetForPlayer(t, conn)
		planetB, _, _ := insertTestPlanetForPlayer(t, conn)
		fleet := insertTestFleet(t, conn, planetA)
		timestamp := fleet.Schedule.ArrivalAt.Add(time.Minute)

		enteredA := make(chan struct{})
		releaseA := make(chan struct{})
		var releaseOnce sync.Once
		release := func() {
			releaseOnce.Do(func() { close(releaseA) })
		}
		t.Cleanup(release)

		enteredB := make(chan struct{})
		doneA := make(chan error, 1)
		doneB := make(chan error, 1)
		resolveFleetEvent := generateFleetSimulator(func(f *models.Fleet) {
			f.Schedule.Status = models.StatusReturning
			f.Version++
		})

		simulateFleetEvent := func(step models.SimulationStep) (models.SimulationStepResult, error) {
			close(enteredA)
			<-releaseA
			return resolveFleetEvent(step)
		}
		executeFleetEvent := func() {
			_, err := adapter.Execute(t.Context(), planetA.Id, timestamp, simulateFleetEvent)
			doneA <- err
		}
		go executeFleetEvent()

		<-enteredA

		blockingCtx, cancel := context.WithTimeout(t.Context(), 100*time.Millisecond)
		defer cancel()
		simulateOtherPlanet := generatePlanetSimulator(func(p *models.Planet) { p.Version++ })
		observeBlockedExecution := func(step models.SimulationStep) (models.SimulationStepResult, error) {
			close(enteredB)
			return simulateOtherPlanet(step)
		}
		executeOtherPlanet := func() {
			_, err := adapter.Execute(blockingCtx, planetB.Id, timestamp, observeBlockedExecution)
			doneB <- err
		}
		go executeOtherPlanet()

		err := <-doneB
		require.ErrorIs(t, err, context.DeadlineExceeded, "Expected deadline exceeded, got err: %v", err)
		select {
		case <-enteredB:
			t.Fatalf("second execution reached its callback while a fleet event was being resolved")
		default:
		}

		assertPlanetVersion(t, conn, planetA.Id, planetA.Version)
		assertPlanetVersion(t, conn, planetB.Id, planetB.Version)

		release()
		require.NoError(t, <-doneA)

		assertPlanetVersion(t, conn, planetA.Id, planetA.Version+1)
		assertPlanetVersion(t, conn, planetB.Id, planetB.Version)
		actual := loadFleetFromDb(t, conn, fleet.Id)
		assert.Equal(t, models.StatusReturning, actual.Schedule.Status)
		assert.Equal(t, fleet.Version+1, actual.Version)
	})

	t.Run("performs planet mutation after fleet event lock is released", func(t *testing.T) {
		adapter, conn := newTestSimulationStepExecutor(t)

		planetA, _, _ := insertTestPlanetForPlayer(t, conn)
		planetB, _, _ := insertTestPlanetForPlayer(t, conn)
		fleet := insertTestFleet(t, conn, planetA)
		timestamp := fleet.Schedule.ArrivalAt.Add(time.Minute)

		enteredA := make(chan struct{})
		releaseA := make(chan struct{})
		var releaseOnce sync.Once
		release := func() {
			releaseOnce.Do(func() { close(releaseA) })
		}
		t.Cleanup(release)

		doneA := make(chan error, 1)
		doneB := make(chan error, 1)
		resolveFleetEvent := generateFleetSimulator(func(f *models.Fleet) {
			f.Schedule.Status = models.StatusReturning
			f.Version++
		})
		simulateFleetEvent := func(step models.SimulationStep) (models.SimulationStepResult, error) {
			close(enteredA)
			<-releaseA
			return resolveFleetEvent(step)
		}
		executeFleetEvent := func() {
			_, err := adapter.Execute(t.Context(), planetA.Id, timestamp, simulateFleetEvent)
			doneA <- err
		}
		go executeFleetEvent()
		<-enteredA

		executionCtx, cancel := context.WithTimeout(t.Context(), 2*time.Second)
		defer cancel()
		startedB := make(chan struct{})
		enteredB := make(chan struct{})
		simulateOtherPlanet := generatePlanetSimulator(func(p *models.Planet) { p.Version++ })
		simulateAfterFleetEvent := func(step models.SimulationStep) (models.SimulationStepResult, error) {
			assert.Nil(t, step.NextFleetEvent, "resolved fleet event should not be selected again")
			close(enteredB)
			return simulateOtherPlanet(step)
		}
		executeOtherPlanet := func() {
			close(startedB)
			_, err := adapter.Execute(executionCtx, planetB.Id, timestamp, simulateAfterFleetEvent)
			doneB <- err
		}
		go executeOtherPlanet()
		<-startedB

		select {
		case <-enteredB:
			t.Fatalf("second execution reached its callback before the fleet event lock was released")
		case <-time.After(100 * time.Millisecond):
		}

		release()
		require.NoError(t, <-doneA, "first execution failed")
		require.NoError(t, <-doneB, "waiting execution failed")
		require.Eventually(t, func() bool {
			select {
			case <-enteredB:
				return true
			default:
				return false
			}
		}, time.Second, 10*time.Millisecond, "waiting execution did not reach its callback")

		assertPlanetVersion(t, conn, planetA.Id, planetA.Version+1)
		assertPlanetVersion(t, conn, planetB.Id, planetB.Version+1)
		actual := loadFleetFromDb(t, conn, fleet.Id)
		assert.Equal(t, models.StatusReturning, actual.Schedule.Status)
		assert.Equal(t, fleet.Version+1, actual.Version)
	})

	t.Run("allows concurrent modification of planets when no fleet event is pending", func(t *testing.T) {
		adapter, conn := newTestSimulationStepExecutor(t)

		planetA, _, _ := insertTestPlanetForPlayer(t, conn)
		planetB, _, _ := insertTestPlanetForPlayer(t, conn)
		enteredA := make(chan struct{})
		releaseA := make(chan struct{})
		var releaseOnce sync.Once
		release := func() {
			releaseOnce.Do(func() { close(releaseA) })
		}
		t.Cleanup(release)
		doneA := make(chan error, 1)
		doneB := make(chan error, 1)
		simulatePlanetUpdate := generatePlanetSimulator(func(p *models.Planet) { p.Version++ })
		simulatePlanetA := func(step models.SimulationStep) (models.SimulationStepResult, error) {
			close(enteredA)
			<-releaseA
			return simulatePlanetUpdate(step)
		}
		executePlanetA := func() {
			_, err := adapter.Execute(t.Context(), planetA.Id, planetA.UpdatedAt, simulatePlanetA)
			doneA <- err
		}
		go executePlanetA()

		<-enteredA

		executePlanetB := func() {
			_, err := adapter.Execute(t.Context(), planetB.Id, planetB.UpdatedAt, simulatePlanetUpdate)
			doneB <- err
		}
		go executePlanetB()

		select {
		case err := <-doneB:
			require.NoError(t, err, "Actual err: %v", err)
		case <-time.After(500 * time.Millisecond):
			t.Fatalf("execution for the second planet was blocked while no fleet event was pending")
		}

		assertPlanetVersion(t, conn, planetA.Id, planetA.Version)
		assertPlanetVersion(t, conn, planetB.Id, planetB.Version+1)

		release()
		require.NoError(t, <-doneA)
		assertPlanetVersion(t, conn, planetA.Id, planetA.Version+1)
	})

	t.Run("serializes concurrent executions on the same planet", func(t *testing.T) {
		adapter, conn := newTestSimulationStepExecutor(t)

		planet, _, _ := insertTestPlanetForPlayer(t, conn, addPlanetResource)
		require.NotEqual(t, 9878.0, planet.Resources[0].Amount)

		executionCtx, cancel := context.WithTimeout(t.Context(), 2*time.Second)
		defer cancel()
		enteredA := make(chan struct{})
		releaseA := make(chan struct{})
		var releaseOnce sync.Once
		release := func() { releaseOnce.Do(func() { close(releaseA) }) }
		t.Cleanup(release)

		doneA := make(chan error, 1)
		doneB := make(chan error, 1)
		simulateFirstUpdate := func(step models.SimulationStep) (models.SimulationStepResult, error) {
			close(enteredA)
			<-releaseA
			step.Planet.Resources[0].Amount = 9878.0
			step.Planet.Version++
			return models.SimulationStepResult{
				Resolution: models.PlanetUpdated,
				Planet:     step.Planet,
			}, nil
		}

		executeFirstUpdate := func() {
			_, err := adapter.Execute(executionCtx, planet.Id, planet.UpdatedAt, simulateFirstUpdate)
			doneA <- err
		}
		go executeFirstUpdate()

		select {
		case <-enteredA:
		case <-time.After(2 * time.Second):
			t.Fatal("first execution did not reach its simulator")
		}

		startedB := make(chan struct{})
		enteredB := make(chan struct{})
		simulateSecondUpdate := func(step models.SimulationStep) (models.SimulationStepResult, error) {
			close(enteredB)
			assert.Equal(t, planet.Version+1, step.Planet.Version, "second execution should load the committed version")
			assert.Equal(t, 9878.0, step.Planet.Resources[0].Amount, "second execution should see the first update")
			step.Planet.Resources[0].Amount++
			step.Planet.Version++
			return models.SimulationStepResult{
				Resolution: models.PlanetUpdated,
				Planet:     step.Planet,
			}, nil
		}
		executeSecondUpdate := func() {
			close(startedB)
			_, err := adapter.Execute(executionCtx, planet.Id, planet.UpdatedAt, simulateSecondUpdate)
			doneB <- err
		}
		go executeSecondUpdate()
		<-startedB

		select {
		case <-enteredB:
			t.Error("second execution reached its simulator while the first execution held the planet")
		case <-time.After(100 * time.Millisecond):
		}

		release()
		assert.NoError(t, <-doneA, "first execution failed")
		assert.NoError(t, <-doneB, "second execution failed")

		assertPlanetResourceAmount(t, conn, planet.Id, planet.Resources[0].Resource, 9879.0)
		assertPlanetVersion(t, conn, planet.Id, planet.Version+2)
	})

	t.Run("allows concurrent modification of planets when a fleet event is not yet due", func(t *testing.T) {
		adapter, conn := newTestSimulationStepExecutor(t)

		planetA, _, _ := insertTestPlanetForPlayer(t, conn)
		planetB, _, _ := insertTestPlanetForPlayer(t, conn)
		fleet := insertTestFleet(t, conn, planetA)
		timestamp := fleet.Schedule.ArrivalAt.Add(-time.Minute)
		require.True(t, timestamp.Before(fleet.Schedule.ArrivalAt))

		enteredA := make(chan struct{})
		releaseA := make(chan struct{})
		var releaseOnce sync.Once
		release := func() {
			releaseOnce.Do(func() { close(releaseA) })
		}
		t.Cleanup(release)

		doneA := make(chan error, 1)
		doneB := make(chan error, 1)
		simulatePlanetUpdate := generatePlanetSimulator(func(p *models.Planet) { p.Version++ })
		simulatePlanetA := func(step models.SimulationStep) (models.SimulationStepResult, error) {
			assert.Nil(t, step.NextFleetEvent, "future fleet event should not be included in the simulation step")
			close(enteredA)
			<-releaseA
			return simulatePlanetUpdate(step)
		}
		executePlanetA := func() {
			_, err := adapter.Execute(t.Context(), planetA.Id, timestamp, simulatePlanetA)
			doneA <- err
		}
		go executePlanetA()

		<-enteredA

		simulatePlanetB := func(step models.SimulationStep) (models.SimulationStepResult, error) {
			assert.Nil(t, step.NextFleetEvent, "future fleet event should not be included in the simulation step")
			return simulatePlanetUpdate(step)
		}
		executePlanetB := func() {
			_, err := adapter.Execute(t.Context(), planetB.Id, timestamp, simulatePlanetB)
			doneB <- err
		}
		go executePlanetB()

		select {
		case err := <-doneB:
			require.NoError(t, err, "Actual err: %v", err)
		case <-time.After(500 * time.Millisecond):
			t.Fatalf("execution for the second planet was blocked by a fleet event scheduled in the future")
		}

		assertPlanetVersion(t, conn, planetA.Id, planetA.Version)
		assertPlanetVersion(t, conn, planetB.Id, planetB.Version+1)
		assertFleetEventExists(t, conn, fleet.Id, fleet.Schedule.ArrivalAt)
		assertFleetEventTypeMatches(t, conn, fleet.Id, models.ArrivalEvent)

		release()
		require.NoError(t, <-doneA)
		assertPlanetVersion(t, conn, planetA.Id, planetA.Version+1)
		assertFleetEventExists(t, conn, fleet.Id, fleet.Schedule.ArrivalAt)
		assertFleetEventTypeMatches(t, conn, fleet.Id, models.ArrivalEvent)
	})

	t.Run("resolves multiple fleet events before resolving planet event", func(t *testing.T) {
		adapter, conn := newTestSimulationStepExecutor(t)

		planet, _, _ := insertTestPlanetForPlayer(t, conn)
		fleetA := insertTestFleet(t, conn, planet)
		fleetB := insertTestFleet(t, conn, planet)
		fleetB.Schedule.ArrivalAt = fleetA.Schedule.ArrivalAt.Add(time.Minute)
		upsertFleetSchedule(t, conn, fleetB)

		timestamp := fleetB.Schedule.ArrivalAt.Add(2 * time.Minute)

		var resolvedFleetIds []uuid.UUID
		planetEventResolved := false
		simulate := func(step models.SimulationStep) (models.SimulationStepResult, error) {
			if step.NextFleetEvent == nil {
				planetEventResolved = true
				step.Planet.Version++
				return models.SimulationStepResult{
					Resolution: models.PlanetUpdated,
					Planet:     step.Planet,
				}, nil
			}

			fleet := step.NextFleetEvent.Fleet
			resolvedFleetIds = append(resolvedFleetIds, fleet.Id)
			fleet.Schedule.Status = models.StatusReturning
			fleet.Version++
			return models.SimulationStepResult{
				Resolution: models.FleetEventResolved,
				Fleet:      &fleet,
			}, nil
		}

		for range 2 {
			result, err := adapter.Execute(t.Context(), planet.Id, timestamp, simulate)
			require.NoError(t, err, "Actual err: %v", err)
			assert.Equal(t, models.FleetEventResolved, result.Resolution)
		}

		require.False(t, planetEventResolved, "planet event was resolved while fleet events were still pending")
		assert.Equal(t, []uuid.UUID{fleetA.Id, fleetB.Id}, resolvedFleetIds)

		result, err := adapter.Execute(t.Context(), planet.Id, timestamp, simulate)
		require.NoError(t, err, "Actual err: %v", err)
		assert.Equal(t, models.PlanetUpdated, result.Resolution)
		require.True(t, planetEventResolved)
		assertPlanetVersion(t, conn, planet.Id, planet.Version+1)
	})
}

func newTestSimulationStepExecutor(t *testing.T) (*SimulationStepExecutor, database.Connection) {
	t.Helper()
	conn := newTestConnection(t)
	return NewSimulationStepExecutor(conn), conn
}

// upsertFleetSchedule persists the fleet's last update time, schedule, and next event.
func upsertFleetSchedule(t *testing.T, conn database.Connection, fleet models.Fleet) {
	t.Helper()

	sqlQuery := `UPDATE fleet SET updated_at = $1 WHERE id = $2`
	affected, err := conn.Exec(t.Context(), sqlQuery, fleet.UpdatedAt, fleet.Id)
	require.NoError(t, err, "Actual err: %v", err)
	require.Equal(t, int64(1), affected)

	sqlQuery = `INSERT INTO fleet_flight (fleet, arrival_at, return_at, status)
		VALUES($1, $2, $3, $4)
		ON CONFLICT (fleet) DO UPDATE SET
			arrival_at = excluded.arrival_at,
			return_at = excluded.return_at,
			status = excluded.status`
	_, err = conn.Exec(
		t.Context(),
		sqlQuery,
		fleet.Id,
		fleet.Schedule.ArrivalAt,
		fleet.Schedule.ReturnAt,
		fleet.Schedule.Status,
	)
	require.NoError(t, err, "Actual err: %v", err)

	nextEvent, err := fleet.NextEvent()
	require.NoError(t, err, "Actual err: %v", err)

	sqlQuery = `INSERT INTO fleet_event (fleet, next_event_at, event_type)
		VALUES($1, $2, $3)
		ON CONFLICT (fleet) DO UPDATE SET
			next_event_at = excluded.next_event_at,
			event_type = excluded.event_type`
	_, err = conn.Exec(
		t.Context(),
		sqlQuery,
		fleet.Id,
		nextEvent.Timestamp,
		nextEvent.Event,
	)
	require.NoError(t, err, "Actual err: %v", err)
}

func generatePlanetSimulator(
	modifier func(*models.Planet),
) func(models.SimulationStep) (models.SimulationStepResult, error) {
	return func(step models.SimulationStep) (models.SimulationStepResult, error) {
		result := models.SimulationStepResult{
			Resolution: models.PlanetUpdated,
			Planet:     step.Planet,
			Fleet:      nil,
		}

		modifier(result.Planet)

		return result, nil
	}
}

func generateFleetSimulator(
	modifier func(*models.Fleet),
) func(models.SimulationStep) (models.SimulationStepResult, error) {
	return func(step models.SimulationStep) (models.SimulationStepResult, error) {
		result := models.SimulationStepResult{
			Resolution: models.FleetEventResolved,
			Planet:     step.Planet,
			Fleet:      new(step.NextFleetEvent.Fleet),
		}

		modifier(result.Fleet)
		result.Planet.Version++

		return result, nil
	}
}

func generateFleetDeleter() func(models.SimulationStep) (models.SimulationStepResult, error) {
	return func(step models.SimulationStep) (models.SimulationStepResult, error) {
		result := models.SimulationStepResult{
			Resolution: models.FleetDeleted,
			Planet:     step.Planet,
			Fleet:      nil,
		}

		result.Planet.Version++

		return result, nil
	}
}
