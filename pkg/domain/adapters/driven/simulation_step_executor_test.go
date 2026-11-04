package drivenadapters

import (
	"context"
	"errors"
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

func TestIT_SimulationStepExecutor_LoadPlanetBehavior(t *testing.T) {
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

	t.Run("Execute", func(t *testing.T) {
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
	})

	t.Run("ExecuteForDelete", func(t *testing.T) {
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
						Resolution: models.PlanetDeleted,
						Planet:     nil,
						Fleet:      nil,
					}

					return result, nil
				}

				actual, err := adapter.ExecuteForDelete(t.Context(), planet.Id, planet.UpdatedAt, simulate)
				require.NoError(t, err, "Actual err: %v", err)

				expectedStep := models.SimulationStep{
					Timestamp:      planet.UpdatedAt,
					Planet:         &planet,
					NextFleetEvent: nil,
				}
				assert.Equal(t, expectedStep, captured)
				expectedResult := models.SimulationStepResult{
					Resolution: models.PlanetDeleted,
					Planet:     nil,
					Fleet:      nil,
				}
				assert.Equal(t, expectedResult, actual)
			})
		}
	})
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

		planet, destination, _ := insertTwoTestPlanets(t, conn)
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

	t.Run("loads fleet event for returning fleet with destination", func(t *testing.T) {
		adapter, conn := newTestSimulationStepExecutor(t)

		planet, destination, _ := insertTwoTestPlanets(t, conn)
		fleet := insertTestFleet(t, conn, planet, addFleetShip)
		fleet.Schedule.Status = models.StatusReturning
		upsertFleetSchedule(t, conn, fleet)
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

func TestIT_SimulationStepExecutor_ExecuteForDelete_LoadFleetBehavior(t *testing.T) {
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

		actual, err := adapter.ExecuteForDelete(t.Context(), planet.Id, afterArrival, simulate)
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

		planet, destination, _ := insertTwoTestPlanets(t, conn)
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

		actual, err := adapter.ExecuteForDelete(t.Context(), planet.Id, afterArrival, simulate)
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

		actual, err := adapter.ExecuteForDelete(t.Context(), planet.Id, afterReturn, simulate)
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

	t.Run("loads fleet event for returning fleet with destination", func(t *testing.T) {
		adapter, conn := newTestSimulationStepExecutor(t)

		planet, destination, _ := insertTwoTestPlanets(t, conn)
		fleet := insertTestFleet(t, conn, planet, addFleetShip)
		fleet.Schedule.Status = models.StatusReturning
		upsertFleetSchedule(t, conn, fleet)
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

		afterReturn := fleet.Schedule.ReturnAt.Add(1 * time.Minute)

		actual, err := adapter.ExecuteForDelete(t.Context(), planet.Id, afterReturn, simulate)
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

			result := models.SimulationStepResult{
				Resolution: models.PlanetDeleted,
				Planet:     nil,
				Fleet:      nil,
			}

			return result, nil
		}

		afterArrival := fleet.Schedule.ArrivalAt.Add(1 * time.Minute)

		actual, err := adapter.ExecuteForDelete(t.Context(), planet1.Id, afterArrival, simulate)
		require.NoError(t, err, "Actual err: %v", err)

		expectedStep := models.SimulationStep{
			Timestamp:      afterArrival,
			Planet:         &planet1,
			NextFleetEvent: nil,
		}
		assert.Equal(t, expectedStep, captured)
		expectedResult := models.SimulationStepResult{
			Resolution: models.PlanetDeleted,
			Planet:     nil,
			Fleet:      nil,
		}
		assert.Equal(t, expectedResult, actual)
	})
}

func TestIT_SimulationStepExecutor_FleetBehavior(t *testing.T) {
	executors := []struct {
		name    string
		execute func(
			*SimulationStepExecutor,
			context.Context,
			uuid.UUID,
			time.Time,
			models.SimulationStepResolver,
		) (models.SimulationStepResult, error)
	}{
		{"Execute", (*SimulationStepExecutor).Execute},
		{"ExecuteForDelete", (*SimulationStepExecutor).ExecuteForDelete},
	}

	for _, executor := range executors {
		t.Run(executor.name, func(t *testing.T) {
			t.Run("resolves fleet event when available", func(t *testing.T) {
				adapter, conn := newTestSimulationStepExecutor(t)

				planet, _, _ := insertTestPlanetForPlayer(t, conn)
				fleet := insertTestFleet(t, conn, planet, addFleetShip)

				var captured models.SimulationStep
				simulate := func(step models.SimulationStep) (models.SimulationStepResult, error) {
					captured = step

					fleet := new(step.NextFleetEvent.Fleet)
					fleet.Schedule.Status = models.StatusReturning
					fleet.UpdatedAt = step.NextFleetEvent.Timestamp
					fleet.Version++

					result := models.SimulationStepResult{
						Resolution: models.FleetEventResolved,
						Planet:     nil,
						Fleet:      fleet,
					}

					return result, nil
				}

				afterArrival := fleet.Schedule.ArrivalAt.Add(1 * time.Minute)

				actual, err := executor.execute(adapter, t.Context(), planet.Id, afterArrival, simulate)
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
					Fleet:      new(fleet),
				}
				expectedResult.Fleet.UpdatedAt = fleet.Schedule.ArrivalAt
				expectedResult.Fleet.Schedule.Status = models.StatusReturning
				expectedResult.Fleet.Version++
				assert.Equal(t, expectedResult, actual)

				fleetFromDb := loadFleetFromDb(t, conn, fleet.Id)
				expectedFleet := fleet
				expectedFleet.UpdatedAt = fleet.Schedule.ArrivalAt
				expectedFleet.Schedule.Status = models.StatusReturning
				expectedFleet.Version++
				assert.Equal(t, expectedFleet, fleetFromDb)
				assertFleetEventExists(t, conn, fleet.Id, fleet.Schedule.ReturnAt)
				assertFleetEventTypeMatches(t, conn, fleet.Id, models.ReturnEvent)
			})

			t.Run("resolves fleet event with destination planet when available", func(t *testing.T) {
				adapter, conn := newTestSimulationStepExecutor(t)

				planet, destination, _ := insertTwoTestPlanets(t, conn)
				fleet := insertTestFleet(t, conn, planet, addFleetShip)
				fleet.Destination.Target = &destination.Id
				upsertFleetDestination(t, conn, fleet)

				var captured models.SimulationStep
				simulate := func(step models.SimulationStep) (models.SimulationStepResult, error) {
					captured = step

					fleet := new(step.NextFleetEvent.Fleet)
					fleet.UpdatedAt = step.NextFleetEvent.Timestamp
					fleet.Version++

					planet := new(*step.NextFleetEvent.Planet)
					planet.UpdatedAt = step.NextFleetEvent.Timestamp
					planet.Fields++
					planet.Version++

					result := models.SimulationStepResult{
						Resolution: models.FleetEventResolved,
						Planet:     planet,
						Fleet:      fleet,
					}

					return result, nil
				}

				afterArrival := fleet.Schedule.ArrivalAt.Add(1 * time.Minute)

				actual, err := executor.execute(adapter, t.Context(), planet.Id, afterArrival, simulate)
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
				expectedResult.Fleet.UpdatedAt = fleet.Schedule.ArrivalAt
				expectedResult.Fleet.Version++
				expectedResult.Planet.UpdatedAt = fleet.Schedule.ArrivalAt
				expectedResult.Planet.Fields++
				expectedResult.Planet.Version++
				assert.Equal(t, expectedResult, actual)

				planetFromDb := loadPlanetFromDb(t, conn, planet.Id)
				assert.Equal(t, planet, planetFromDb)
				destinationFromDb := loadPlanetFromDb(t, conn, destination.Id)
				assert.Equal(t, *expectedResult.Planet, destinationFromDb)
				fleetFromDb := loadFleetFromDb(t, conn, fleet.Id)
				assert.Equal(t, *expectedResult.Fleet, fleetFromDb)
			})

			t.Run("resolves earliest fleet event when multiple are available", func(t *testing.T) {
				adapter, conn := newTestSimulationStepExecutor(t)

				planet, _, _ := insertTestPlanetForPlayer(t, conn)
				fleetA := insertTestFleet(t, conn, planet)
				fleetB := insertTestFleet(t, conn, planet, addFleetShip)
				fleetA.Schedule.ArrivalAt = fleetB.Schedule.ArrivalAt.Add(time.Minute)
				upsertFleetSchedule(t, conn, fleetA)

				var captured models.SimulationStep
				simulate := func(step models.SimulationStep) (models.SimulationStepResult, error) {
					captured = step

					fleet := new(step.NextFleetEvent.Fleet)
					// Deep copy: otherwise the captured step will also be impacted by the
					// change to the slice as new does not deep copy slices.
					fleet.Ships = make([]models.FleetShip, len(step.NextFleetEvent.Fleet.Ships))
					copy(fleet.Ships, step.NextFleetEvent.Fleet.Ships)
					fleet.Ships[0].Count++
					fleet.Version++

					result := models.SimulationStepResult{
						Resolution: models.FleetEventResolved,
						Planet:     nil,
						Fleet:      fleet,
					}

					return result, nil
				}

				afterArrival := fleetA.Schedule.ArrivalAt.Add(1 * time.Minute)

				actual, err := executor.execute(adapter, t.Context(), planet.Id, afterArrival, simulate)
				require.NoError(t, err, "Actual err: %v", err)

				expectedStep := models.SimulationStep{
					Timestamp: afterArrival,
					Planet:    nil,
					NextFleetEvent: &models.FleetEventContext{
						Timestamp: fleetB.Schedule.ArrivalAt,
						Fleet:     fleetB,
						Planet:    nil,
					},
				}
				assert.Equal(t, expectedStep, captured)
				expectedResult := models.SimulationStepResult{
					Resolution: models.FleetEventResolved,
					Planet:     nil,
					Fleet:      &fleetB,
				}
				expectedResult.Fleet.Ships[0].Count++
				expectedResult.Fleet.Version++
				assert.Equal(t, expectedResult, actual)

				fleetBFromDb := loadFleetFromDb(t, conn, fleetB.Id)
				assert.Equal(t, *expectedResult.Fleet, fleetBFromDb)
				fleetAFromDb := loadFleetFromDb(t, conn, fleetA.Id)
				assert.Equal(t, fleetA, fleetAFromDb)
			})

			t.Run("resolves fleet event exactly at selected timestamp", func(t *testing.T) {
				adapter, conn := newTestSimulationStepExecutor(t)

				planet, destination, _ := insertTwoTestPlanets(t, conn)
				addPlanetResource(t, conn, &destination)
				fleet := insertTestFleet(t, conn, planet, addFleetShip)
				fleet.Destination.Target = &destination.Id
				upsertFleetDestination(t, conn, fleet)

				var captured models.SimulationStep
				simulate := func(step models.SimulationStep) (models.SimulationStepResult, error) {
					captured = step

					fleet := new(step.NextFleetEvent.Fleet)
					fleet.Version++

					planet := new(*step.NextFleetEvent.Planet)
					// Deep copy: otherwise the captured step will also be impacted by the
					// change to the slice as new does not deep copy slices.
					planet.Resources = make([]models.PlanetResource, len(step.NextFleetEvent.Planet.Resources))
					copy(planet.Resources, step.NextFleetEvent.Planet.Resources)
					planet.Resources[0].Amount++
					planet.Version++

					result := models.SimulationStepResult{
						Resolution: models.FleetEventResolved,
						Planet:     planet,
						Fleet:      fleet,
					}

					return result, nil
				}

				atArrival := fleet.Schedule.ArrivalAt

				actual, err := executor.execute(adapter, t.Context(), planet.Id, fleet.Schedule.ArrivalAt, simulate)
				require.NoError(t, err, "Actual err: %v", err)

				expectedStep := models.SimulationStep{
					Timestamp: atArrival,
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
				expectedResult.Planet.Resources[0].Amount++
				expectedResult.Planet.Version++
				expectedResult.Fleet.Version++
				assert.Equal(t, expectedResult, actual)

				fleetFromDb := loadFleetFromDb(t, conn, fleet.Id)
				assert.Equal(t, *expectedResult.Fleet, fleetFromDb)
				planetFromDb := loadPlanetFromDb(t, conn, planet.Id)
				assert.Equal(t, planet, planetFromDb)
				destinationFromDb := loadPlanetFromDb(t, conn, destination.Id)
				assert.Equal(t, *expectedResult.Planet, destinationFromDb)
			})

			t.Run("resolves fleet event with fleet identifier in case timestamp tie", func(t *testing.T) {
				adapter, conn := newTestSimulationStepExecutor(t)

				planetA, planetB, universe := insertTwoTestPlanets(t, conn)
				planetC, _ := insertTestPlanetInUniverse(t, conn, universe)
				fleetA := insertTestFleetWithId(t, conn, uuid.MustParse("7f8e8adc-5957-417b-931b-64afc528c3d8"), planetA, addFleetShip)
				fleetB := insertTestFleetWithId(t, conn, uuid.MustParse("6e25a02b-625f-4161-8bb4-b30d46d94f29"), planetB, addFleetShip)
				fleetB.Schedule.ArrivalAt = fleetA.Schedule.ArrivalAt
				upsertFleetSchedule(t, conn, fleetB)

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

				afterArrival := fleetA.Schedule.ArrivalAt.Add(1 * time.Minute)

				actual, err := executor.execute(adapter, t.Context(), planetC.Id, afterArrival, simulate)
				require.NoError(t, err, "Actual err: %v", err)

				expectedStep := models.SimulationStep{
					Timestamp: afterArrival,
					Planet:    nil,
					NextFleetEvent: &models.FleetEventContext{
						Timestamp: fleetB.Schedule.ArrivalAt,
						Fleet:     fleetB,
						Planet:    nil,
					},
				}
				assert.Equal(t, expectedStep, captured)
				expectedResult := models.SimulationStepResult{
					Resolution: models.FleetEventResolved,
					Planet:     nil,
					Fleet:      &fleetB,
				}
				expectedResult.Fleet.Version++
				assert.Equal(t, expectedResult, actual)

				fleetAFromDb := loadFleetFromDb(t, conn, fleetA.Id)
				assert.Equal(t, fleetA, fleetAFromDb)
				fleetBFromDb := loadFleetFromDb(t, conn, fleetB.Id)
				assert.Equal(t, *expectedResult.Fleet, fleetBFromDb)
				planetAFromDb := loadPlanetFromDb(t, conn, planetA.Id)
				assert.Equal(t, planetA, planetAFromDb)
				planetBFromDb := loadPlanetFromDb(t, conn, planetB.Id)
				assert.Equal(t, planetB, planetBFromDb)
				planetCFromDb := loadPlanetFromDb(t, conn, planetC.Id)
				assert.Equal(t, planetC, planetCFromDb)
			})

			t.Run("resolves fleet event affecting other planet in same universe", func(t *testing.T) {
				adapter, conn := newTestSimulationStepExecutor(t)

				planetA, planetB, _ := insertTwoTestPlanets(t, conn)
				fleet := insertTestFleet(t, conn, planetA, addFleetShip)

				var captured models.SimulationStep
				simulate := func(step models.SimulationStep) (models.SimulationStepResult, error) {
					captured = step

					fleet := new(step.NextFleetEvent.Fleet)
					fleet.UpdatedAt = step.NextFleetEvent.Timestamp
					fleet.Version++

					result := models.SimulationStepResult{
						Resolution: models.FleetEventResolved,
						Planet:     nil,
						Fleet:      fleet,
					}

					return result, nil
				}

				afterArrival := fleet.Schedule.ArrivalAt.Add(1 * time.Minute)

				actual, err := executor.execute(adapter, t.Context(), planetB.Id, afterArrival, simulate)
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
				expectedResult.Fleet.UpdatedAt = fleet.Schedule.ArrivalAt
				expectedResult.Fleet.Version++
				assert.Equal(t, expectedResult, actual)

				planetAFromDb := loadPlanetFromDb(t, conn, planetA.Id)
				assert.Equal(t, planetA, planetAFromDb)
				planetBFromDb := loadPlanetFromDb(t, conn, planetB.Id)
				assert.Equal(t, planetB, planetBFromDb)
				fleetFromDb := loadFleetFromDb(t, conn, fleet.Id)
				assert.Equal(t, *expectedResult.Fleet, fleetFromDb)
			})

			t.Run("resolves fleet event and performs deletion when requested", func(t *testing.T) {
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
					planet.Fields++
					planet.Version++

					result := models.SimulationStepResult{
						Resolution: models.FleetDeleted,
						Planet:     planet,
						Fleet:      nil,
					}

					return result, nil
				}

				afterReturn := fleet.Schedule.ReturnAt.Add(1 * time.Minute)

				actual, err := executor.execute(adapter, t.Context(), planet.Id, afterReturn, simulate)
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
					Resolution: models.FleetDeleted,
					Planet:     &planet,
					Fleet:      nil,
				}
				expectedResult.Planet.Fields++
				expectedResult.Planet.Version++
				assert.Equal(t, expectedResult, actual)

				assertFleetDoesNotExist(t, conn, fleet.Id)
				planetFromDb := loadPlanetFromDb(t, conn, planet.Id)
				assert.Equal(t, expectedResult.Planet, &planetFromDb)
			})

			t.Run("returns error when planet does not exist", func(t *testing.T) {
				adapter, _ := newTestSimulationStepExecutor(t)

				var called int
				simulate := func(step models.SimulationStep) (models.SimulationStepResult, error) {
					called++
					return models.SimulationStepResult{}, errors.New("should not happen!!")
				}

				_, err := executor.execute(adapter, t.Context(), uuid.New(), t2, simulate)
				assert.ErrorIs(t, err, domainerrors.ErrNotFound, "Actual err: %v", err)

				assert.Zero(t, called)
			})

			t.Run("does not persist changes to fleet when resolver returns error", func(t *testing.T) {
				adapter, conn := newTestSimulationStepExecutor(t)

				planet, _, _ := insertTestPlanetForPlayer(t, conn)
				fleet := insertTestFleet(t, conn, planet)
				initialVersion := fleet.Version

				errSample := errors.New("stubbed error")
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

					return result, errSample
				}

				afterArrival := fleet.Schedule.ArrivalAt.Add(1 * time.Minute)

				_, err := executor.execute(adapter, t.Context(), planet.Id, afterArrival, simulate)
				assert.ErrorIs(t, err, errSample, "Actual err: %v", err)

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
				actual := loadFleetFromDb(t, conn, fleet.Id)
				assert.Equal(t, initialVersion, actual.Version)
			})

			t.Run("returns error when resolver reports planet updated with a fleet event", func(t *testing.T) {
				adapter, conn := newTestSimulationStepExecutor(t)

				planet, destination, _ := insertTwoTestPlanets(t, conn)
				fleet := insertTestFleet(t, conn, planet, addFleetShip)
				fleet.Destination.Target = &destination.Id
				upsertFleetDestination(t, conn, fleet)

				var captured models.SimulationStep
				simulate := func(step models.SimulationStep) (models.SimulationStepResult, error) {
					captured = step

					planet := new(*step.NextFleetEvent.Planet)
					planet.Version++

					result := models.SimulationStepResult{
						Resolution: models.PlanetUpdated,
						Planet:     planet,
						Fleet:      nil,
					}

					return result, nil
				}

				afterArrival := fleet.Schedule.ArrivalAt.Add(1 * time.Minute)

				_, err := executor.execute(adapter, t.Context(), planet.Id, afterArrival, simulate)

				assert.ErrorIs(t, err, domainerrors.ErrInvalidSimulationStepResult, "Actual err: %v", err)
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

				planetFromDb := loadPlanetFromDb(t, conn, planet.Id)
				assert.Equal(t, planet, planetFromDb)
			})

			t.Run("returns error when resolver reports planet deleted with a fleet event", func(t *testing.T) {
				adapter, conn := newTestSimulationStepExecutor(t)

				planet, destination, _ := insertTwoTestPlanets(t, conn)
				fleet := insertTestFleet(t, conn, planet, addFleetShip)
				fleet.Destination.Target = &destination.Id
				upsertFleetDestination(t, conn, fleet)

				var captured models.SimulationStep
				simulate := func(step models.SimulationStep) (models.SimulationStepResult, error) {
					captured = step

					planet := new(*step.NextFleetEvent.Planet)
					planet.Version++

					result := models.SimulationStepResult{
						Resolution: models.PlanetDeleted,
						Planet:     planet,
						Fleet:      nil,
					}

					return result, nil
				}

				afterArrival := fleet.Schedule.ArrivalAt.Add(1 * time.Minute)

				_, err := executor.execute(adapter, t.Context(), planet.Id, afterArrival, simulate)

				assert.ErrorIs(t, err, domainerrors.ErrInvalidSimulationStepResult, "Actual err: %v", err)
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

				planetFromDb := loadPlanetFromDb(t, conn, planet.Id)
				assert.Equal(t, planet, planetFromDb)
			})
		})
	}

	t.Run("does not resolve fleet event in another universe when execute is called", func(t *testing.T) {
		adapter, conn := newTestSimulationStepExecutor(t)

		planetA, _, universeA := insertTestPlanetForPlayer(t, conn)
		planetB, _, universeB := insertTestPlanetForPlayer(t, conn)
		require.NotEqual(t, universeA.Id, universeB.Id)
		fleetA := insertTestFleet(t, conn, planetA)

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

		afterArrival := fleetA.Schedule.ArrivalAt.Add(time.Minute)

		actual, err := adapter.Execute(t.Context(), planetB.Id, afterArrival, simulate)
		require.NoError(t, err, "Actual err: %v", err)

		expectedStep := models.SimulationStep{
			Timestamp:      afterArrival,
			Planet:         &planetB,
			NextFleetEvent: nil,
		}
		assert.Equal(t, expectedStep, captured)
		expectedResult := models.SimulationStepResult{
			Resolution: models.PlanetUpdated,
			Planet:     &planetB,
			Fleet:      nil,
		}
		expectedResult.Planet.Version++
		assert.Equal(t, expectedResult, actual)

		fleetFromDb := loadFleetFromDb(t, conn, fleetA.Id)
		assert.Equal(t, fleetA, fleetFromDb)
		planetAFromDb := loadPlanetFromDb(t, conn, planetA.Id)
		assert.Equal(t, planetA, planetAFromDb)
		planetBFromDb := loadPlanetFromDb(t, conn, planetB.Id)
		assert.Equal(t, *expectedResult.Planet, planetBFromDb)
	})

	t.Run("does not resolve fleet event in another universe when execute for delete is called", func(t *testing.T) {
		adapter, conn := newTestSimulationStepExecutor(t)

		planetA, _, universeA := insertTestPlanetForPlayer(t, conn)
		planetB, _, universeB := insertTestPlanetForPlayer(t, conn)
		require.NotEqual(t, universeA.Id, universeB.Id)
		fleetA := insertTestFleet(t, conn, planetA)

		var captured models.SimulationStep
		simulate := func(step models.SimulationStep) (models.SimulationStepResult, error) {
			captured = step

			result := models.SimulationStepResult{
				Resolution: models.PlanetDeleted,
				Planet:     nil,
				Fleet:      nil,
			}

			return result, nil
		}

		afterArrival := fleetA.Schedule.ArrivalAt.Add(time.Minute)

		actual, err := adapter.ExecuteForDelete(t.Context(), planetB.Id, afterArrival, simulate)
		require.NoError(t, err, "Actual err: %v", err)

		expectedStep := models.SimulationStep{
			Timestamp:      afterArrival,
			Planet:         &planetB,
			NextFleetEvent: nil,
		}
		assert.Equal(t, expectedStep, captured)
		expectedResult := models.SimulationStepResult{
			Resolution: models.PlanetDeleted,
			Planet:     nil,
			Fleet:      nil,
		}
		assert.Equal(t, expectedResult, actual)

		fleetFromDb := loadFleetFromDb(t, conn, fleetA.Id)
		assert.Equal(t, fleetA, fleetFromDb)
		planetAFromDb := loadPlanetFromDb(t, conn, planetA.Id)
		assert.Equal(t, planetA, planetAFromDb)
		// TODO: Should be deleted
		planetBFromDb := loadPlanetFromDb(t, conn, planetB.Id)
		assert.Equal(t, planetB, planetBFromDb)
	})

	t.Run("does not resolve future fleet event when execute is called", func(t *testing.T) {
		adapter, conn := newTestSimulationStepExecutor(t)

		planet, _, _ := insertTestPlanetForPlayer(t, conn)
		fleet := insertTestFleet(t, conn, planet)

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

		beforeArrival := fleet.Schedule.ArrivalAt.Add(-1 * time.Minute)

		actual, err := adapter.Execute(t.Context(), planet.Id, beforeArrival, simulate)
		require.NoError(t, err, "Actual err: %v", err)

		expectedStep := models.SimulationStep{
			Timestamp:      beforeArrival,
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

		fleetFromDb := loadFleetFromDb(t, conn, fleet.Id)
		assert.Equal(t, fleet, fleetFromDb)
		planetFromDb := loadPlanetFromDb(t, conn, planet.Id)
		assert.Equal(t, *expectedResult.Planet, planetFromDb)
	})

	t.Run("does not resolve future fleet event when execute for delete is called", func(t *testing.T) {
		adapter, conn := newTestSimulationStepExecutor(t)

		planet, _, _ := insertTestPlanetForPlayer(t, conn)
		fleet := insertTestFleet(t, conn, planet)

		var captured models.SimulationStep
		simulate := func(step models.SimulationStep) (models.SimulationStepResult, error) {
			captured = step

			result := models.SimulationStepResult{
				Resolution: models.PlanetDeleted,
				Planet:     nil,
				Fleet:      nil,
			}

			return result, nil
		}

		beforeArrival := fleet.Schedule.ArrivalAt.Add(-1 * time.Minute)

		actual, err := adapter.ExecuteForDelete(t.Context(), planet.Id, beforeArrival, simulate)
		require.NoError(t, err, "Actual err: %v", err)

		expectedStep := models.SimulationStep{
			Timestamp:      beforeArrival,
			Planet:         &planet,
			NextFleetEvent: nil,
		}
		assert.Equal(t, expectedStep, captured)
		expectedResult := models.SimulationStepResult{
			Resolution: models.PlanetDeleted,
			Planet:     nil,
			Fleet:      nil,
		}
		assert.Equal(t, expectedResult, actual)

		fleetFromDb := loadFleetFromDb(t, conn, fleet.Id)
		assert.Equal(t, fleet, fleetFromDb)
		// TODO: Should be deleted
		planetFromDb := loadPlanetFromDb(t, conn, planet.Id)
		assert.Equal(t, planet, planetFromDb)
	})
}

func TestIT_SimulationStepExecutor_Execute_PlanetBehavior(t *testing.T) {
	// TODO: Implement tests
}

func TestIT_SimulationStepExecutor_ExecuteForDelete_PlanetBehavior(t *testing.T) {
	// TODO: Implement tests
}

func TestIT_SimulationStepExecutor_Execute_Concurrency(t *testing.T) {
	t.Run("prevents concurrent modification of different planets when fleet event is resolved", func(t *testing.T) {
		adapter, conn := newTestSimulationStepExecutor(t)

		planetA, planetB, _ := insertTwoTestPlanets(t, conn)
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

		planetA, planetB, _ := insertTwoTestPlanets(t, conn)
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

		planetA, planetB, _ := insertTwoTestPlanets(t, conn)
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

		planetA, planetB, _ := insertTwoTestPlanets(t, conn)
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

func TestIT_SimulationStepExecutor_Execute_Atomicity(t *testing.T) {
	t.Run("rolls back planet update when fleet creation fails", func(t *testing.T) {
		adapter, conn := newTestSimulationStepExecutor(t)
		planet, _, _ := insertTestPlanetForPlayer(t, conn)

		fleet := generateSampleFleet(planet)
		fleet.Source = uuid.New() // Nonexistent source.

		resolver := func(step models.SimulationStep) (models.SimulationStepResult, error) {
			step.Planet.Fields++
			step.Planet.Version++

			return models.SimulationStepResult{
				Resolution: models.FleetCreated,
				Planet:     step.Planet,
				Fleet:      &fleet,
			}, nil
		}

		result, err := adapter.Execute(
			t.Context(), planet.Id, planet.UpdatedAt, resolver,
		)

		assert.ErrorIs(t, err, domainerrors.ErrPlanetNotFound)
		assert.Equal(t, models.SimulationStepResult{}, result)

		actual := loadPlanetFromDb(t, conn, planet.Id)
		assert.Equal(t, planet, actual)
		assertFleetDoesNotExist(t, conn, fleet.Id)
	})
}

func newTestSimulationStepExecutor(t *testing.T) (*SimulationStepExecutor, database.Connection) {
	t.Helper()
	conn := newTestConnection(t)
	return NewSimulationStepExecutor(conn), conn
}

// insertTwoTestPlanets inserts two planets in the same universe for two different players
// returning both planets and the identifier of the universe they belong to
func insertTwoTestPlanets(t *testing.T, conn database.Connection) (models.Planet, models.Planet, uuid.UUID) {
	universe := insertTestUniverse(t, conn)

	player1 := insertTestPlayer(t, conn, universe.Id)
	player2 := insertTestPlayer(t, conn, universe.Id)

	planet1 := insertTestPlanet(t, conn, player1.Id)
	planet2 := insertTestPlanet(t, conn, player2.Id)

	return planet1, planet2, universe.Id
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
	return insertTestFleetWithId(t, conn, uuid.New(), source, modifiers...)
}

func insertTestFleetWithId(
	t *testing.T,
	conn database.Connection,
	id uuid.UUID,
	source models.Planet,
	modifiers ...func(*testing.T, database.Connection, *models.Fleet),
) models.Fleet {
	t.Helper()

	fleet := models.Fleet{
		Id:      id,
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
