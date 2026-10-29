package drivenadapters

import (
	"context"
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/Knoblauchpilze/backend-toolkit/pkg/db"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/adapters/driven/database"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/adapters/driven/mappers"
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

		planet, destination := insertTwoTestPlanets(t, conn)
		fleet := insertTestFleet(t, conn, planet, addFleetShip)
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

	t.Run("does not load fleet event when planet is not from the same universe", func(t *testing.T) {
		adapter, conn := newTestSimulationStepExecutor(t)

		universe1 := insertTestUniverse(t, conn)
		universe2 := insertTestUniverse(t, conn)

		player1 := insertTestPlayer(t, conn, universe1.Id)
		player2 := insertTestPlayer(t, conn, universe2.Id)

		planet1 := insertTestPlanet(t, conn, player1.Id)
		planet2 := insertTestPlanet(t, conn, player2.Id)

		fleet := insertTestFleet(t, conn, planet2, addFleetShip)

		var captured models.SimulationStep
		simulate := func(step models.SimulationStep) (models.SimulationStepResult, error) {
			captured = step

			planet := new(*step.Planet)
			planet.Fields++
			planet.Version++

			result := models.SimulationStepResult{
				Resolution: models.PlanetUpdated,
				Planet:     planet,
				Fleet:      nil,
			}

			return result, nil
		}

		afterArrival := fleet.Schedule.ArrivalAt.Add(1 * time.Minute)

		actual, err := adapter.Execute(t.Context(), planet1.Id, afterArrival, simulate)
		require.NoError(t, err, "Actual err: %v", err)

		expectedStep := models.SimulationStep{
			Timestamp:      afterArrival,
			Planet:         &planet1,
			NextFleetEvent: nil,
		}
		assert.Equal(t, expectedStep, captured)
		expectedResult := models.SimulationStepResult{
			Resolution: models.PlanetUpdated,
			Planet:     &planet1,
			Fleet:      nil,
		}
		expectedResult.Planet.Fields++
		expectedResult.Planet.Version++
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

		simulate := generatePlanetSimulator(func(p *models.Planet) {
			p.Fields++
		})

		_, err := adapter.Execute(t.Context(), planet.Id, planet.UpdatedAt, simulate)

		assert.ErrorIs(t, err, domainerrors.ErrMutationWithoutVersionBump, "Actual err: %v", err)

		actual := loadPlanetFromDb(t, conn, planet.Id)
		assert.Equal(t, planet, actual)
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

func TestIT_SimulationStepExecutor_Execute_CreateFleet(t *testing.T) {
	t.Run("returns created fleet", func(t *testing.T) {
		adapter, conn := newTestSimulationStepExecutor(t)

		planet, _, _ := insertTestPlanetForPlayer(t, conn)
		fleet := generateSampleFleet(planet)

		simulate := generateFleetCreator(fleet)

		returned, err := adapter.Execute(t.Context(), planet.Id, planet.UpdatedAt, simulate)
		require.NoError(t, err, "Actual err: %v", err)

		require.NotNil(t, returned.Fleet)
		assert.Equal(t, &fleet, returned.Fleet)
	})

	t.Run("persists created fleet", func(t *testing.T) {
		adapter, conn := newTestSimulationStepExecutor(t)

		planet, _, _ := insertTestPlanetForPlayer(t, conn)
		fleet := generateSampleFleet(planet)

		simulate := generateFleetCreator(fleet)

		_, err := adapter.Execute(t.Context(), planet.Id, planet.UpdatedAt, simulate)
		require.NoError(t, err, "Actual err: %v", err)

		actual := loadFleetFromDb(t, conn, fleet.Id)
		assert.Equal(t, fleet, actual)
	})

	t.Run("persists fleet event", func(t *testing.T) {
		adapter, conn := newTestSimulationStepExecutor(t)

		planet, _, _ := insertTestPlanetForPlayer(t, conn)
		fleet := generateSampleFleet(planet)

		simulate := generateFleetCreator(fleet)

		_, err := adapter.Execute(t.Context(), planet.Id, planet.UpdatedAt, simulate)
		require.NoError(t, err, "Actual err: %v", err)

		assertFleetEventExists(t, conn, fleet.Id, fleet.Schedule.ArrivalAt)
		assertFleetEventTypeMatches(t, conn, fleet.Id, models.ArrivalEvent)
	})

	t.Run("persists fleet with destination", func(t *testing.T) {
		adapter, conn := newTestSimulationStepExecutor(t)

		planetA, planetB := insertTwoTestPlanets(t, conn)

		fleet := generateSampleFleet(planetA)
		fleet.Mission = models.MissionTransport
		fleet.Destination.Target = &planetB.Id

		simulate := generateFleetCreator(fleet)

		_, err := adapter.Execute(t.Context(), planetA.Id, planetA.UpdatedAt, simulate)
		require.NoError(t, err, "Actual err: %v", err)

		actual := loadFleetFromDb(t, conn, fleet.Id)
		assert.Equal(t, fleet, actual)
	})

	t.Run("persists updated planet", func(t *testing.T) {
		adapter, conn := newTestSimulationStepExecutor(t)

		planet, _, _ := insertTestPlanetForPlayer(t, conn, addPlanetShip)
		fleet := generateSampleFleet(planet)

		simulate := func(step models.SimulationStep) (models.SimulationStepResult, error) {
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

		_, err := adapter.Execute(t.Context(), planet.Id, planet.UpdatedAt, simulate)
		require.NoError(t, err, "Actual err: %v", err)

		actual := loadPlanetFromDb(t, conn, planet.Id)
		expected := planet
		expected.Fields++
		expected.Version++
		expected.Ships[0].Count++
		assert.Equal(t, expected, actual)
	})

	t.Run("returns error when fleet source does not exist", func(t *testing.T) {
		adapter, conn := newTestSimulationStepExecutor(t)

		planet, _, _ := insertTestPlanetForPlayer(t, conn)
		fleet := generateSampleFleet(planet)
		fleet.Source = uuid.New()

		simulate := func(step models.SimulationStep) (models.SimulationStepResult, error) {
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

		_, err := adapter.Execute(t.Context(), planet.Id, planet.UpdatedAt, simulate)

		assert.ErrorIs(t, err, domainerrors.ErrPlanetNotFound, "Actual err: %v", err)
		assertFleetDoesNotExist(t, conn, fleet.Id)
	})

	t.Run("returns error when fleet destination does not exist", func(t *testing.T) {
		adapter, conn := newTestSimulationStepExecutor(t)

		planetA, _, _ := insertTestPlanetForPlayer(t, conn)

		fleet := generateSampleFleet(planetA)
		fleet.Mission = models.MissionTransport
		fakePlanetId := uuid.New()
		fleet.Destination.Target = &fakePlanetId

		simulate := generateFleetCreator(fleet)

		_, err := adapter.Execute(t.Context(), planetA.Id, planetA.UpdatedAt, simulate)

		assert.ErrorIs(t, err, domainerrors.ErrFleetDestinationInvalid, "Actual err: %v", err)
		assertFleetDoesNotExist(t, conn, fleet.Id)
	})

	t.Run("returns error when planet version is not bumped", func(t *testing.T) {
		adapter, conn := newTestSimulationStepExecutor(t)

		planet, _, _ := insertTestPlanetForPlayer(t, conn, addPlanetShip)
		fleet := generateSampleFleet(planet)

		simulate := func(step models.SimulationStep) (models.SimulationStepResult, error) {
			result := models.SimulationStepResult{
				Resolution: models.FleetCreated,
				Planet:     step.Planet,
				Fleet:      &fleet,
			}

			return result, nil
		}

		_, err := adapter.Execute(t.Context(), planet.Id, planet.UpdatedAt, simulate)

		assert.ErrorIs(t, err, domainerrors.ErrMutationWithoutVersionBump, "Actual err: %v", err)
		assertFleetDoesNotExist(t, conn, fleet.Id)
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

		planet, destination := insertTwoTestPlanets(t, conn)
		fleet := insertTestFleet(t, conn, planet, addFleetShip)
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

		planetA, planetB := insertTwoTestPlanets(t, conn)
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

		assertFleetVersion(t, conn, fleet.Id, fleet.Version)

		release()
		require.NoError(t, <-doneA)

		assertPlanetVersion(t, conn, planetA.Id, planetA.Version)
		assertPlanetVersion(t, conn, planetB.Id, planetB.Version)
		actual := loadFleetFromDb(t, conn, fleet.Id)
		assert.Equal(t, models.StatusReturning, actual.Schedule.Status)
		assert.Equal(t, fleet.Version+1, actual.Version)
	})

	t.Run("allows planet mutation in another universe while a due fleet event is locked", func(t *testing.T) {
		adapter, conn := newTestSimulationStepExecutor(t)

		planetA, _, universeA := insertTestPlanetForPlayer(t, conn)
		planetB, _, universeB := insertTestPlanetForPlayer(t, conn)
		require.NotEqual(t, universeA.Id, universeB.Id)
		fleet := insertTestFleet(t, conn, planetA)
		timestamp := fleet.Schedule.ArrivalAt.Add(time.Minute)

		executionCtx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
		defer cancel()
		enteredA := make(chan struct{})
		releaseA := make(chan struct{})
		var releaseOnce sync.Once
		release := func() {
			releaseOnce.Do(func() { close(releaseA) })
		}
		t.Cleanup(release)
		doneA := make(chan error, 1)
		resolveFleetEvent := generateFleetSimulator(func(f *models.Fleet) {
			f.Schedule.Status = models.StatusReturning
			f.Version++
		})
		simulateFleetEvent := func(step models.SimulationStep) (models.SimulationStepResult, error) {
			close(enteredA)
			select {
			case <-releaseA:
			case <-executionCtx.Done():
				return models.SimulationStepResult{}, executionCtx.Err()
			}
			return resolveFleetEvent(step)
		}
		go func() {
			_, err := adapter.Execute(executionCtx, planetA.Id, timestamp, simulateFleetEvent)
			doneA <- err
		}()

		select {
		case <-enteredA:
		case err := <-doneA:
			t.Fatalf("fleet execution exited before reaching its callback: %v", err)
		case <-executionCtx.Done():
			t.Fatal("fleet execution did not acquire its event lock")
		}

		planetCtx, cancelPlanet := context.WithTimeout(executionCtx, time.Second)
		defer cancelPlanet()
		simulatePlanetUpdate := generatePlanetSimulator(func(p *models.Planet) { p.Version++ })
		simulateOtherUniverse := func(step models.SimulationStep) (models.SimulationStepResult, error) {
			assert.Nil(t, step.NextFleetEvent, "fleet event from another universe should not be selected")
			return simulatePlanetUpdate(step)
		}
		result, err := adapter.Execute(planetCtx, planetB.Id, timestamp, simulateOtherUniverse)
		require.NoError(t, err, "planet update was blocked by a fleet event in another universe: %v", err)
		assert.Equal(t, models.PlanetUpdated, result.Resolution)

		assertPlanetVersion(t, conn, planetB.Id, planetB.Version+1)
		assertPlanetVersion(t, conn, planetA.Id, planetA.Version)
		assertFleetVersion(t, conn, fleet.Id, fleet.Version)
		assertFleetEventExists(t, conn, fleet.Id, fleet.Schedule.ArrivalAt)

		release()
		require.NoError(t, <-doneA, "fleet execution failed")
		actual := loadFleetFromDb(t, conn, fleet.Id)
		assert.Equal(t, models.StatusReturning, actual.Schedule.Status)
		assert.Equal(t, fleet.Version+1, actual.Version)
		assertFleetEventExists(t, conn, fleet.Id, fleet.Schedule.ReturnAt)
		assertFleetEventTypeMatches(t, conn, fleet.Id, models.ReturnEvent)
	})

	t.Run("performs planet mutation after fleet event lock is released", func(t *testing.T) {
		adapter, conn := newTestSimulationStepExecutor(t)

		planetA, planetB := insertTwoTestPlanets(t, conn)
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

		assertFleetVersion(t, conn, fleet.Id, fleet.Version)
		assertPlanetVersion(t, conn, planetB.Id, planetB.Version)

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

		assertPlanetVersion(t, conn, planetA.Id, planetA.Version)
		assertPlanetVersion(t, conn, planetB.Id, planetB.Version+1)
		actual := loadFleetFromDb(t, conn, fleet.Id)
		assert.Equal(t, models.StatusReturning, actual.Schedule.Status)
		assert.Equal(t, fleet.Version+1, actual.Version)
	})

	t.Run("allows concurrent modification of planets when no fleet event is pending", func(t *testing.T) {
		adapter, conn := newTestSimulationStepExecutor(t)

		planetA, planetB := insertTwoTestPlanets(t, conn)
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

		planetA, planetB := insertTwoTestPlanets(t, conn)
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

	t.Run("allows concurrent reciprocal fleet dispatches", func(t *testing.T) {
		adapter, conn := newTestSimulationStepExecutor(t)

		planetA, _, _ := insertTestPlanetForPlayer(t, conn, addPlanetShip)
		planetB := insertTestPlanet(t, conn, planetA.Player, addPlanetShip)

		fleetA := generateSampleFleet(planetA)
		fleetA.Mission = models.MissionTransport
		fleetA.Destination.Target = &planetB.Id
		fleetA.Destination.Coordinate = planetB.Coordinate
		fleetA.Ships = []models.FleetShip{{Ship: lightFighterId, Count: 1}}

		fleetB := generateSampleFleet(planetB)
		fleetB.Mission = models.MissionTransport
		fleetB.Destination.Target = &planetA.Id
		fleetB.Destination.Coordinate = planetA.Coordinate
		fleetB.Ships = []models.FleetShip{{Ship: lightFighterId, Count: 1}}

		executionCtx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
		defer cancel()
		entered := make(chan struct{}, 2)
		releaseSimulations := make(chan struct{})
		var releaseOnce sync.Once
		release := func() { releaseOnce.Do(func() { close(releaseSimulations) }) }
		t.Cleanup(release)
		doneA := make(chan error, 1)
		doneB := make(chan error, 1)

		createFleetA := generateFleetCreator(fleetA)
		simulateDispatchA := func(step models.SimulationStep) (models.SimulationStepResult, error) {
			entered <- struct{}{}
			select {
			case <-releaseSimulations:
			case <-executionCtx.Done():
				return models.SimulationStepResult{}, executionCtx.Err()
			}
			step.Planet.Ships[0].Count--
			return createFleetA(step)
		}
		createFleetB := generateFleetCreator(fleetB)
		simulateDispatchB := func(step models.SimulationStep) (models.SimulationStepResult, error) {
			entered <- struct{}{}
			select {
			case <-releaseSimulations:
			case <-executionCtx.Done():
				return models.SimulationStepResult{}, executionCtx.Err()
			}
			step.Planet.Ships[0].Count--
			return createFleetB(step)
		}
		executeDispatchA := func() {
			_, err := adapter.Execute(executionCtx, planetA.Id, planetA.UpdatedAt, simulateDispatchA)
			doneA <- err
		}
		executeDispatchB := func() {
			_, err := adapter.Execute(executionCtx, planetB.Id, planetB.UpdatedAt, simulateDispatchB)
			doneB <- err
		}
		go executeDispatchA()
		go executeDispatchB()

		for range 2 {
			select {
			case <-entered:
			case <-executionCtx.Done():
				t.Fatal("both dispatches did not acquire their source planet locks")
			}
		}
		release()
		errA, errB := <-doneA, <-doneB
		require.NoError(t, errA, "Actual err: %v", errA)
		require.NoError(t, errB, "Actual err: %v", errB)

		actualFleetA := loadFleetFromDb(t, conn, fleetA.Id)
		assert.Equal(t, fleetA, actualFleetA)
		assertFleetEventExists(t, conn, fleetA.Id, fleetA.Schedule.ArrivalAt)

		actualFleetB := loadFleetFromDb(t, conn, fleetB.Id)
		assert.Equal(t, fleetB, actualFleetB)
		assertFleetEventExists(t, conn, fleetB.Id, fleetB.Schedule.ArrivalAt)

		actualPlanetA := loadPlanetFromDb(t, conn, planetA.Id)
		assert.Equal(t, planetA.Ships[0].Count-1, actualPlanetA.Ships[0].Count)
		assertPlanetVersion(t, conn, planetA.Id, planetA.Version+1)

		actualPlanetB := loadPlanetFromDb(t, conn, planetB.Id)
		assert.Equal(t, planetB.Ships[0].Count-1, actualPlanetB.Ships[0].Count)
		assertPlanetVersion(t, conn, planetB.Id, planetB.Version+1)
	})

	t.Run("serializes fleet dispatches from the same source planet", func(t *testing.T) {
		adapter, conn := newTestSimulationStepExecutor(t)

		planet, _, _ := insertTestPlanetForPlayer(t, conn, addPlanetShip)
		fleetA := generateSampleFleet(planet)
		fleetA.Ships = []models.FleetShip{{Ship: lightFighterId, Count: 1}}
		fleetB := generateSampleFleet(planet)
		fleetB.Ships = []models.FleetShip{{Ship: smallCargoId, Count: 1}}

		executionCtx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
		defer cancel()
		enteredA := make(chan struct{})
		releaseA := make(chan struct{})
		var releaseOnce sync.Once
		release := func() { releaseOnce.Do(func() { close(releaseA) }) }
		t.Cleanup(release)
		enteredB := make(chan struct{})
		startedB := make(chan struct{})
		doneA := make(chan error, 1)
		doneB := make(chan error, 1)

		createFleetA := generateFleetCreator(fleetA)
		simulateFirstDispatch := func(step models.SimulationStep) (models.SimulationStepResult, error) {
			close(enteredA)
			select {
			case <-releaseA:
			case <-executionCtx.Done():
				return models.SimulationStepResult{}, executionCtx.Err()
			}
			// Decrease the ship count on the first mutation
			step.Planet.Ships[0].Count--
			return createFleetA(step)
		}

		createFleetB := generateFleetCreator(fleetB)
		simulateSecondDispatch := func(step models.SimulationStep) (models.SimulationStepResult, error) {
			close(enteredB)
			// The second mutation should see the updated version and ship count
			assert.Equal(t, planet.Version+1, step.Planet.Version)
			assert.Equal(t, planet.Ships[0].Count-1, step.Planet.Ships[0].Count,
				"second dispatch should see the first dispatch's ship deduction")
			step.Planet.Ships[0].Count--
			return createFleetB(step)
		}

		executeFirstDispatch := func() {
			_, err := adapter.Execute(executionCtx, planet.Id, planet.UpdatedAt, simulateFirstDispatch)
			doneA <- err
		}
		executeSecondDispatch := func() {
			close(startedB)
			_, err := adapter.Execute(executionCtx, planet.Id, planet.UpdatedAt, simulateSecondDispatch)
			doneB <- err
		}
		go executeFirstDispatch()
		select {
		case <-enteredA:
		case <-executionCtx.Done():
			t.Fatal("first dispatch did not reach its simulator")
		}
		go executeSecondDispatch()
		<-startedB
		select {
		case <-enteredB:
			t.Error("second dispatch reached its simulator while the first held the source planet lock")
		case <-time.After(100 * time.Millisecond):
		}

		release()
		errA, errB := <-doneA, <-doneB
		require.NoError(t, errA, "first dispatch failed")
		require.NoError(t, errB, "second dispatch failed")

		assert.Equal(t, fleetA, loadFleetFromDb(t, conn, fleetA.Id))
		assert.Equal(t, fleetB, loadFleetFromDb(t, conn, fleetB.Id))
		assertFleetEventExists(t, conn, fleetA.Id, fleetA.Schedule.ArrivalAt)
		assertFleetEventExists(t, conn, fleetB.Id, fleetB.Schedule.ArrivalAt)
		assert.Equal(t, planet.Ships[0].Count-2, loadPlanetFromDb(t, conn, planet.Id).Ships[0].Count)
		assertPlanetVersion(t, conn, planet.Id, planet.Version+2)
	})

	t.Run("rolls back source planet changes when fleet creation fails", func(t *testing.T) {
		adapter, conn := newTestSimulationStepExecutor(t)

		planet, _, _ := insertTestPlanetForPlayer(t, conn, addPlanetShip, addPlanetResource)
		fleet := generateSampleFleet(planet)
		fleet.Mission = models.MissionTransport
		missingDestination := uuid.New()
		fleet.Destination.Target = &missingDestination
		fleet.Ships = []models.FleetShip{{Ship: lightFighterId, Count: 1}}

		createFleet := generateFleetCreator(fleet)
		simulateDispatch := func(step models.SimulationStep) (models.SimulationStepResult, error) {
			step.Planet.Ships[0].Count--
			step.Planet.Resources[0].Amount++
			return createFleet(step)
		}
		_, err := adapter.Execute(t.Context(), planet.Id, planet.UpdatedAt, simulateDispatch)

		assert.ErrorIs(t, err, domainerrors.ErrFleetDestinationInvalid, "Actual err: %v", err)

		actual := loadPlanetFromDb(t, conn, planet.Id)
		assert.Equal(t, planet, actual)
		assertFleetDoesNotExist(t, conn, fleet.Id)
	})
}

func TestIT_SimulationStepExecutor_Execute_ActionCreationDeletionWorkflow(t *testing.T) {
	conn := newTestConnection(t)
	simulator := NewSimulationStepExecutor(conn)

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

	mutation := func(step models.SimulationStep) (models.SimulationStepResult, error) {
		require.NotNil(t, step.Planet)
		require.Nil(t, step.NextFleetEvent)

		step.Planet.BuildingAction = &action
		step.Planet.Version++

		result := models.SimulationStepResult{
			Resolution: models.PlanetUpdated,
			Planet:     step.Planet,
		}
		return result, nil
	}
	result, err := simulator.Execute(t.Context(), planet.Id, action.CreatedAt, mutation)
	require.NoError(t, err, "Actual err: %v", err)
	assert.Equal(t, models.PlanetUpdated, result.Resolution)
	require.NotNil(t, result.Planet)
	require.NotNil(t, result.Planet.BuildingAction)
	assert.Equal(t, action, *result.Planet.BuildingAction)
	assertBuildingActionExists(t, conn, action.Id)

	func() {
		actual := loadPlanetFromDb(t, conn, planet.Id)

		require.NotNil(t, actual.BuildingAction)
		assert.Equal(t, action, *actual.BuildingAction)
	}()

	mutation = func(step models.SimulationStep) (models.SimulationStepResult, error) {
		require.NotNil(t, step.Planet)
		require.Nil(t, step.NextFleetEvent)

		step.Planet.BuildingAction = nil
		step.Planet.Version++

		result := models.SimulationStepResult{
			Resolution: models.PlanetUpdated,
			Planet:     step.Planet,
		}

		return result, nil
	}
	result, err = simulator.Execute(t.Context(), planet.Id, action.CreatedAt, mutation)
	require.NoError(t, err, "Actual err: %v", err)
	require.NotNil(t, result.Planet)
	assert.Nil(t, result.Planet.BuildingAction)
	assertBuildingActionDoesNotExist(t, conn, action.Id)
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

// insertTwoTestPlanets inserts two planets in the same universe for two different players
func insertTwoTestPlanets(t *testing.T, conn database.Connection) (models.Planet, models.Planet) {
	universe := insertTestUniverse(t, conn)

	player1 := insertTestPlayer(t, conn, universe.Id)
	player2 := insertTestPlayer(t, conn, universe.Id)

	planet1 := insertTestPlanet(t, conn, player1.Id)
	planet2 := insertTestPlanet(t, conn, player2.Id)

	return planet1, planet2
}

func insertTestBuildingActionForPlanet(
	t *testing.T,
	conn database.Connection,
	planetId uuid.UUID,
	modifiers ...func(*testing.T, database.Connection, *models.BuildingAction),
) models.BuildingAction {
	t.Helper()

	action := models.BuildingAction{
		Id:           uuid.New(),
		Building:     metalMineId,
		DesiredLevel: 5,
		CreatedAt:    someTime,
		CompletedAt:  someTime.Add(1*time.Hour + 2*time.Minute),
		// This is intentional: the details (e.g. costs, productions, etc.) are returned as empty
		// slices by the adapter
		Costs:       []models.BuildingActionCost{},
		Storages:    []models.BuildingActionResourceStorage{},
		Productions: []models.BuildingActionResourceProduction{},
	}

	sqlQuery := `INSERT INTO building_action
		(id, planet, building, desired_level, created_at, completed_at)
		VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := conn.Exec(
		t.Context(),
		sqlQuery,
		action.Id,
		planetId,
		action.Building,
		action.DesiredLevel,
		action.CreatedAt,
		action.CompletedAt,
	)
	require.NoError(t, err, "Actual err: %v", err)

	for _, modifier := range modifiers {
		modifier(t, conn, &action)
	}

	return action
}

func insertTestBuildingAction(
	t *testing.T,
	conn database.Connection,
	modifiers ...func(*testing.T, database.Connection, *models.BuildingAction),
) (models.BuildingAction, models.Planet) {
	t.Helper()

	planet, _, _ := insertTestPlanetForPlayer(t, conn)
	action := insertTestBuildingActionForPlanet(t, conn, planet.Id, modifiers...)
	planet.BuildingAction = &action
	return action, planet
}

func addBuildingActionCost(t *testing.T, conn database.Connection, a *models.BuildingAction) {
	t.Helper()

	insertBuildingActionCost(t, conn, metalResourceId, a)
}

func insertBuildingActionCost(t *testing.T, conn database.Connection, resourceId uuid.UUID, a *models.BuildingAction) {
	cost := models.BuildingActionCost{
		Resource: resourceId,
		Amount:   rand.Intn(4589),
	}

	sqlQuery := `INSERT INTO building_action_cost (action, resource, amount)
		VALUES ($1, $2, $3)`
	_, err := conn.Exec(
		t.Context(),
		sqlQuery,
		a.Id,
		cost.Resource,
		cost.Amount,
	)
	require.NoError(t, err, "Actual err: %v", err)

	a.Costs = append(a.Costs, cost)
}

func addBuildingActionStorage(t *testing.T, conn database.Connection, a *models.BuildingAction) {
	t.Helper()

	storage := models.BuildingActionResourceStorage{
		Resource: crystalResourceId,
		Storage:  rand.Intn(65114),
	}

	sqlQuery := `INSERT INTO building_action_resource_storage (action, resource, storage)
		VALUES ($1, $2, $3)`
	_, err := conn.Exec(
		t.Context(),
		sqlQuery,
		a.Id,
		storage.Resource,
		storage.Storage,
	)
	require.NoError(t, err, "Actual err: %v", err)

	a.Storages = append(a.Storages, storage)
}

func addBuildingActionProduction(t *testing.T, conn database.Connection, a *models.BuildingAction) {
	t.Helper()

	production := models.BuildingActionResourceProduction{
		Resource:   crystalResourceId,
		Production: rand.Intn(7451),
	}

	sqlQuery := `INSERT INTO building_action_resource_production (action, resource, production)
		VALUES ($1, $2, $3)`
	_, err := conn.Exec(
		t.Context(),
		sqlQuery,
		a.Id,
		production.Resource,
		production.Production,
	)
	require.NoError(t, err, "Actual err: %v", err)

	a.Productions = append(a.Productions, production)
}

func insertTestShipActionForPlanet(
	t *testing.T,
	conn database.Connection,
	planetId uuid.UUID,
	modifiers ...func(*testing.T, *models.ShipAction),
) models.ShipAction {
	t.Helper()

	action := models.ShipAction{
		Id:                 uuid.New(),
		Ship:               lightFighterId,
		Count:              5,
		CreatedAt:          someTime,
		NextCompletionAt:   someTime.Add(someDuration),
		UnitCompletionTime: someDuration,
	}

	for _, modifier := range modifiers {
		modifier(t, &action)
	}

	sqlQuery := `INSERT INTO ship_action
		(id, planet, ship, count, created_at, next_completion_at, unit_completion_time)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := conn.Exec(
		t.Context(),
		sqlQuery,
		action.Id,
		planetId,
		action.Ship,
		action.Count,
		action.CreatedAt,
		action.NextCompletionAt,
		action.UnitCompletionTime,
	)
	require.NoError(t, err, "Actual err: %v", err)

	return action
}

func insertTestShipAction(
	t *testing.T,
	conn database.Connection,
	modifiers ...func(*testing.T, *models.ShipAction),
) (models.ShipAction, models.Planet) {
	t.Helper()

	planet, _, _ := insertTestPlanetForPlayer(t, conn)
	action := insertTestShipActionForPlanet(t, conn, planet.Id, modifiers...)
	planet.ShipActions = append(planet.ShipActions, action)
	return action, planet
}

func insertTestFleet(
	t *testing.T,
	conn database.Connection,
	source models.Planet,
	modifiers ...func(*testing.T, database.Connection, *models.Fleet),
) models.Fleet {
	t.Helper()

	fleet := models.Fleet{
		Id:      uuid.New(),
		Player:  source.Player,
		Source:  source.Id,
		Mission: models.MissionColonize,
		Destination: models.FleetDestination{
			Coordinate: models.Coordinate{
				Galaxy:      source.Coordinate.Galaxy,
				SolarSystem: source.Coordinate.SolarSystem + 1,
				Position:    source.Coordinate.Position,
			},
			Target: nil,
		},
		Ships: []models.FleetShip{},
		Schedule: models.FleetFlightSchedule{
			ArrivalAt: source.UpdatedAt.Add(1 * time.Hour),
			ReturnAt:  source.UpdatedAt.Add(2 * time.Hour),
			Status:    models.StatusOutbound,
		},
		CreatedAt: source.UpdatedAt,
		UpdatedAt: source.UpdatedAt,
		Version:   0,
	}

	sqlQuery := `INSERT INTO
		fleet (id, player, source, mission, created_at, updated_at, version)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := conn.Exec(
		t.Context(),
		sqlQuery,
		fleet.Id,
		fleet.Player,
		fleet.Source,
		fleet.Mission,
		fleet.CreatedAt,
		fleet.UpdatedAt,
		fleet.Version,
	)
	require.NoError(t, err, "Actual err: %v", err)

	sqlQuery = `INSERT INTO
		fleet_destination (fleet, galaxy, solar_system, position, planet)
		VALUES ($1, $2, $3, $4, $5)`
	_, err = conn.Exec(
		t.Context(),
		sqlQuery,
		fleet.Id,
		fleet.Destination.Coordinate.Galaxy,
		fleet.Destination.Coordinate.SolarSystem,
		fleet.Destination.Coordinate.Position,
		fleet.Destination.Target,
	)
	require.NoError(t, err, "Actual err: %v", err)

	sqlQuery = `INSERT INTO
		fleet_flight (fleet, arrival_at, return_at, status)
		VALUES ($1, $2, $3, $4)`
	_, err = conn.Exec(
		t.Context(),
		sqlQuery,
		fleet.Id,
		fleet.Schedule.ArrivalAt,
		fleet.Schedule.ReturnAt,
		fleet.Schedule.Status,
	)
	require.NoError(t, err, "Actual err: %v", err)

	sqlQuery = `INSERT INTO
		fleet_event (fleet, next_event_at, event_type)
		VALUES ($1, $2, $3)`
	_, err = conn.Exec(
		t.Context(),
		sqlQuery,
		fleet.Id,
		fleet.Schedule.ArrivalAt,
		models.ArrivalEvent,
	)
	require.NoError(t, err, "Actual err: %v", err)

	for _, modifier := range modifiers {
		modifier(t, conn, &fleet)
	}

	return fleet
}

func addFleetShip(t *testing.T, conn database.Connection, f *models.Fleet) {
	t.Helper()

	ship := models.FleetShip{
		Ship:  lightFighterId,
		Count: 18,
	}

	sqlQuery := `INSERT INTO fleet_ship (fleet, ship, count) VALUES ($1, $2, $3)`
	_, err := conn.Exec(
		t.Context(),
		sqlQuery,
		f.Id,
		ship.Ship,
		ship.Count,
	)
	require.NoError(t, err, "Actual err: %v", err)

	f.Ships = append(f.Ships, ship)
}

// upsertFleetDestination inserts the coordinates/target defined for the fleet in the
// database, or updates them if a row already exists for this fleet. It does
// not verify them compared to the universe the planet belongs to and will fail in
// case the target is not nil and does not correspond to an existing planet.
func upsertFleetDestination(t *testing.T, conn database.Connection, fleet models.Fleet) {
	t.Helper()

	sqlQuery := `INSERT INTO fleet_destination (fleet, galaxy, solar_system, position, planet)
		VALUES($1, $2, $3, $4, $5)
		ON CONFLICT (fleet) DO UPDATE SET
			galaxy = excluded.galaxy,
			solar_system = excluded.solar_system,
			position = excluded.position,
			planet = excluded.planet`
	_, err := conn.Exec(
		t.Context(),
		sqlQuery,
		fleet.Id,
		fleet.Destination.Coordinate.Galaxy,
		fleet.Destination.Coordinate.SolarSystem,
		fleet.Destination.Coordinate.Position,
		fleet.Destination.Target,
	)
	require.NoError(t, err, "Actual err: %v", err)
}

func generatePlanetSimulator(
	modifier func(*models.Planet),
) models.SimulationStepResolver {
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
) models.SimulationStepResolver {
	return func(step models.SimulationStep) (models.SimulationStepResult, error) {
		var planet *models.Planet
		if step.NextFleetEvent.Planet != nil {
			planet = new(*step.NextFleetEvent.Planet)
			planet.Version++
		}

		result := models.SimulationStepResult{
			Resolution: models.FleetEventResolved,
			Planet:     planet,
			Fleet:      new(step.NextFleetEvent.Fleet),
		}

		modifier(result.Fleet)

		return result, nil
	}
}

func generateFleetCreator(fleet models.Fleet) models.SimulationStepResolver {
	return func(step models.SimulationStep) (models.SimulationStepResult, error) {
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
}

func generateFleetDeleter() models.SimulationStepResolver {
	return func(step models.SimulationStep) (models.SimulationStepResult, error) {
		var planet *models.Planet
		if step.NextFleetEvent.Planet != nil {
			planet = new(*step.NextFleetEvent.Planet)
			planet.Version++
		}

		result := models.SimulationStepResult{
			Resolution: models.FleetDeleted,
			Planet:     planet,
			Fleet:      nil,
		}

		return result, nil
	}
}

func generateSampleFleet(p models.Planet) models.Fleet {
	return models.Fleet{
		Id:      uuid.New(),
		Player:  p.Player,
		Source:  p.Id,
		Mission: models.MissionColonize,
		Destination: models.FleetDestination{
			Coordinate: models.Coordinate{
				Galaxy:      p.Coordinate.Galaxy,
				SolarSystem: p.Coordinate.SolarSystem + 1,
				Position:    p.Coordinate.Position,
			},
			Target: nil,
		},
		Ships: []models.FleetShip{
			{
				Ship:  lightFighterId,
				Count: 12,
			},
			{
				Ship:  smallCargoId,
				Count: 3,
			},
		},
		Schedule: models.FleetFlightSchedule{
			ArrivalAt: p.UpdatedAt.Add(1 * time.Hour),
			ReturnAt:  p.UpdatedAt.Add(2 * time.Hour),
			Status:    models.StatusOutbound,
		},
		CreatedAt: p.UpdatedAt,
		UpdatedAt: p.UpdatedAt,
		Version:   0,
	}
}

func assertBuildingActionExists(t *testing.T, conn database.Connection, id uuid.UUID) {
	t.Helper()

	sqlQuery := `SELECT COUNT(*) FROM building_action WHERE id = $1`
	value, err := db.QueryOne[int](t.Context(), conn, sqlQuery, id)
	require.NoError(t, err, "Actual err: %v", err)
	require.Equal(t, 1, value)
}

func assertBuildingActionDoesNotExist(t *testing.T, conn database.Connection, action uuid.UUID) {
	t.Helper()

	sqlQuery := `SELECT COUNT(*) FROM building_action WHERE id = $1`
	value, err := db.QueryOne[int](t.Context(), conn, sqlQuery, action)
	require.NoError(t, err, "Actual err: %v", err)
	require.Zero(t, value)
}

func assertBuildingActionCostDoesNotExist(t *testing.T, conn database.Connection, action uuid.UUID) {
	t.Helper()

	sqlQuery := `SELECT COUNT(*) FROM building_action_cost WHERE action = $1`
	value, err := db.QueryOne[int](t.Context(), conn, sqlQuery, action)
	require.NoError(t, err, "Actual err: %v", err)
	require.Zero(t, value)
}

func assertBuildingActionStorageDoesNotExist(t *testing.T, conn database.Connection, action uuid.UUID) {
	t.Helper()

	sqlQuery := `SELECT COUNT(*) FROM building_action_resource_storage WHERE action = $1`
	value, err := db.QueryOne[int](t.Context(), conn, sqlQuery, action)
	require.NoError(t, err, "Actual err: %v", err)
	require.Zero(t, value)
}

func assertBuildingActionProductionDoesNotExist(t *testing.T, conn database.Connection, action uuid.UUID) {
	t.Helper()

	sqlQuery := `SELECT COUNT(*) FROM building_action_resource_production WHERE action = $1`
	value, err := db.QueryOne[int](t.Context(), conn, sqlQuery, action)
	require.NoError(t, err, "Actual err: %v", err)
	require.Zero(t, value)
}

func assertPlanetResourceAmount(t *testing.T, conn database.Connection, planet uuid.UUID, resource uuid.UUID, amount float64) {
	t.Helper()

	sqlQuery := `SELECT amount FROM planet_resource WHERE planet = $1 AND resource = $2`
	value, err := db.QueryOne[float64](t.Context(), conn, sqlQuery, planet, resource)
	require.NoError(t, err, "Actual err: %v", err)
	require.InDelta(t, amount, value, 0.00001)
}

func assertPlanetResourceStorage(t *testing.T, conn database.Connection, planet uuid.UUID, resource uuid.UUID, storage int) {
	t.Helper()

	sqlQuery := `SELECT storage FROM planet_resource_storage WHERE planet = $1 AND resource = $2`
	value, err := db.QueryOne[int](t.Context(), conn, sqlQuery, planet, resource)
	require.NoError(t, err, "Actual err: %v", err)
	require.Equal(t, storage, value)
}

func assertPlanetResourceProduction(
	t *testing.T,
	conn database.Connection,
	planet uuid.UUID,
	resource uuid.UUID,
	building *uuid.UUID,
	production int,
) {
	t.Helper()

	sqlQuery := `SELECT production
		FROM planet_resource_production
		WHERE
			planet = $1
			AND resource = $2
			AND building IS NOT DISTINCT FROM $3`
	value, err := db.QueryOne[int](t.Context(), conn, sqlQuery, planet, resource, building)
	require.NoError(t, err, "Actual err: %v", err)
	require.Equal(t, production, value)
}

func assertPlanetBuildingLevel(t *testing.T, conn database.Connection, planet uuid.UUID, building uuid.UUID, level int) {
	t.Helper()

	sqlQuery := `SELECT level FROM planet_building WHERE planet = $1 AND building = $2`
	value, err := db.QueryOne[int](t.Context(), conn, sqlQuery, planet, building)
	require.NoError(t, err, "Actual err: %v", err)
	require.Equal(t, level, value)
}

func assertPlanetShipCount(t *testing.T, conn database.Connection, planet uuid.UUID, ship uuid.UUID, count int) {
	t.Helper()

	sqlQuery := `SELECT count FROM planet_ship WHERE planet = $1 AND ship = $2`
	value, err := db.QueryOne[int](t.Context(), conn, sqlQuery, planet, ship)
	require.NoError(t, err, "Actual err: %v", err)
	require.Equal(t, count, value)
}

func assertPlanetVersion(t *testing.T, conn database.Connection, planet uuid.UUID, version int) {
	t.Helper()

	sqlQuery := `SELECT version FROM planet WHERE id = $1`
	value, err := db.QueryOne[int](t.Context(), conn, sqlQuery, planet)
	require.NoError(t, err, "Actual err: %v", err)
	require.Equal(t, version, value)
}

func assertFleetVersion(t *testing.T, conn database.Connection, fleet uuid.UUID, version int) {
	t.Helper()

	sqlQuery := `SELECT version FROM fleet WHERE id = $1`
	value, err := db.QueryOne[int](t.Context(), conn, sqlQuery, fleet)
	require.NoError(t, err, "Actual err: %v", err)
	require.Equal(t, version, value)
}

func assertFleetDoesNotExist(t *testing.T, conn database.Connection, id uuid.UUID) {
	t.Helper()

	sqlQuery := `SELECT COUNT(id) FROM fleet WHERE id = $1`
	value, err := db.QueryOne[int](t.Context(), conn, sqlQuery, id)
	require.NoError(t, err, "Actual err: %v", err)
	require.Zero(t, value)
}

func assertFleetEventExists(t *testing.T, conn database.Connection, fleet uuid.UUID, moment time.Time) {
	t.Helper()

	sqlQuery := `SELECT next_event_at FROM fleet_event WHERE fleet = $1`
	value, err := db.QueryOne[time.Time](t.Context(), conn, sqlQuery, fleet)
	require.NoError(t, err, "Actual err: %v", err)
	require.Equal(t, moment, value)
}

func assertFleetEventTypeMatches(t *testing.T, conn database.Connection, fleet uuid.UUID, eventType models.FleetEventType) {
	t.Helper()

	sqlQuery := `SELECT event_type FROM fleet_event WHERE fleet = $1`
	value, err := db.QueryOne[mappers.DbFleetEventType](t.Context(), conn, sqlQuery, fleet)
	require.NoError(t, err, "Actual err: %v", err)
	require.Equal(t, eventType, value.ToDomain())
}

func loadFleetFromDb(t *testing.T, conn database.Connection, id uuid.UUID) models.Fleet {
	t.Helper()

	tx, err := conn.BeginTx(t.Context())
	require.NoError(t, err, "Actual err: %v", err)
	defer tx.Close(t.Context())

	fleet, err := loadFleetAndDetails(t.Context(), tx, id)
	require.NoError(t, err, "Actual err: %v", err)

	return fleet
}
