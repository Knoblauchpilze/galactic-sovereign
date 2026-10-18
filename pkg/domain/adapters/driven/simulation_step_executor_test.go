package drivenadapters

import (
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

				step.Planet.Version++

				result := models.SimulationStepResult{
					Resolution: models.PlanetUpdated,
					Planet:     &step.Planet,
					Fleet:      nil,
				}

				return result, nil
			}

			actual, err := adapter.Execute(t.Context(), planet.Id, planet.UpdatedAt, simulate)
			require.NoError(t, err, "Actual err: %v", err)

			expectedStep := models.SimulationStep{
				Timestamp:      planet.UpdatedAt,
				Planet:         planet,
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

			result := models.SimulationStepResult{
				Resolution: models.FleetEventResolved,
				Planet:     new(step.Planet),
				Fleet:      new(step.NextFleetEvent.Fleet),
			}

			result.Planet.Version++
			result.Fleet.Version++

			return result, nil
		}

		afterArrival := fleet.Schedule.ArrivalAt.Add(1 * time.Minute)

		actual, err := adapter.Execute(t.Context(), planet.Id, afterArrival, simulate)
		require.NoError(t, err, "Actual err: %v", err)

		expectedStep := models.SimulationStep{
			Timestamp: afterArrival,
			Planet:    planet,
			NextFleetEvent: &models.FleetEventContext{
				Timestamp: fleet.Schedule.ArrivalAt,
				Fleet:     fleet,
				Planet:    nil,
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

			result := models.SimulationStepResult{
				Resolution: models.FleetEventResolved,
				Planet:     new(step.Planet),
				Fleet:      new(step.NextFleetEvent.Fleet),
			}

			result.Planet.Version++
			result.Fleet.Version++

			return result, nil
		}

		afterArrival := fleet.Schedule.ArrivalAt.Add(1 * time.Minute)

		actual, err := adapter.Execute(t.Context(), planet.Id, afterArrival, simulate)
		require.NoError(t, err, "Actual err: %v", err)

		expectedStep := models.SimulationStep{
			Timestamp: afterArrival,
			Planet:    planet,
			NextFleetEvent: &models.FleetEventContext{
				Timestamp: fleet.Schedule.ArrivalAt,
				Fleet:     fleet,
				Planet:    &destination,
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

// TODO: Implement tests
func TestIT_SimulationStepExecutor_Execute_SaveFleet(t *testing.T) {}

func TestIT_SimulationStepExecutor_Execute_DeleteFleet(t *testing.T) {}

func TestIT_SimulationStepExecutor_Execute_Concurrency(t *testing.T) {}

func TestIT_SimulationStepExecutor_Execute_ProcessFleetEvents(t *testing.T) {}

func newTestSimulationStepExecutor(t *testing.T) (*SimulationStepExecutor, database.Connection) {
	t.Helper()
	conn := newTestConnection(t)
	return NewSimulationStepExecutor(conn), conn
}

// upsertFleetSchedule inserts the arrival/return times and status defined for the fleet
// in the database, or updates them if a row already exists for this fleet.
func upsertFleetSchedule(t *testing.T, conn database.Connection, fleet models.Fleet) {
	t.Helper()

	sqlQuery := `INSERT INTO fleet_flight (fleet, arrival_at, return_at, status)
		VALUES($1, $2, $3, $4)
		ON CONFLICT (fleet) DO UPDATE SET
			arrival_at = excluded.arrival_at,
			return_at = excluded.return_at,
			status = excluded.status`
	_, err := conn.Exec(
		t.Context(),
		sqlQuery,
		fleet.Id,
		fleet.Schedule.ArrivalAt,
		fleet.Schedule.ReturnAt,
		fleet.Schedule.Status,
	)
	require.NoError(t, err, "Actual err: %v", err)
}

func generatePlanetSimulator(
	modifier func(*models.Planet),
) func(models.SimulationStep) (models.SimulationStepResult, error) {
	return func(step models.SimulationStep) (models.SimulationStepResult, error) {
		result := models.SimulationStepResult{
			Resolution: models.PlanetUpdated,
			Planet:     new(step.Planet),
			Fleet:      nil,
		}

		modifier(result.Planet)

		return result, nil
	}
}
