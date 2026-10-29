package usecases

import (
	"errors"
	"testing"
	"time"

	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	domainerrors "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models/errors"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/usecases/drivenportstest"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

type createFleetTestSuite struct {
	ctrl             *gomock.Controller
	mockUniverseRepo *drivenportstest.MockForFetchingUniverses
	mockExecutor     *drivenportstest.MockForExecutingSimulationStep
	mockClock        *drivenportstest.MockForFetchingTime
	usecase          *CreateFleetUseCase
}

func TestUnit_CreateFleet_Create(t *testing.T) {
	t.Run("calls executor with planet and timestamp", func(t *testing.T) {
		suite := setupCreateFleetTestSuite(t)

		planet := generateTestPlanetWithShip()
		req := models.FleetOrder{
			Destination: models.Coordinate{
				Galaxy:      1,
				SolarSystem: 3,
				Position:    4,
			},
			Ships: []models.FleetShip{
				{Ship: uuid.New(), Count: 3},
				{Ship: uuid.New(), Count: 4},
			},
		}

		universe := models.Universe{
			Id: uuid.New(),
			Topology: models.UniverseTopology{
				Galaxies:     2,
				SolarSystems: 10,
				Orbits:       14,
			},
			Ships: []models.Ship{
				{Id: req.Ships[0].Ship},
				{Id: req.Ships[1].Ship},
			},
		}

		fleet := models.Fleet{
			Id:      uuid.New(),
			Player:  planet.Player,
			Source:  planet.Id,
			Mission: models.MissionColonize,
			Destination: models.FleetDestination{
				Coordinate: models.Coordinate{
					Galaxy:      req.Destination.Galaxy,
					SolarSystem: req.Destination.SolarSystem,
					Position:    req.Destination.Position,
				},
			},
			Ships: []models.FleetShip{
				{Ship: req.Ships[0].Ship, Count: req.Ships[0].Count},
				{Ship: req.Ships[1].Ship, Count: req.Ships[1].Count},
			},
			Schedule: models.FleetFlightSchedule{
				ArrivalAt: t3,
				ReturnAt:  t4,
			},
			CreatedAt: t2,
			UpdatedAt: t2,
			Version:   0,
		}

		suite.mockUniverseRepo.EXPECT().
			GetByPlanetId(gomock.Any(), planet.Id).
			Times(1).
			Return(universe, nil)
		suite.mockClock.EXPECT().Now(gomock.Any()).Times(1).Return(t2)
		result := models.SimulationStepResult{
			Resolution: models.FleetCreated,
			Planet:     &planet,
			Fleet:      &fleet,
		}
		suite.mockExecutor.EXPECT().
			Execute(gomock.Any(), planet.Id, t2, gomock.Any()).
			Times(1).
			Return(result, nil)

		actual, err := suite.usecase.Create(t.Context(), planet.Id, req)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, fleet, actual)
	})

	t.Run("retries executor when fleet event is resolved", func(t *testing.T) {
		suite := setupCreateFleetTestSuite(t)

		planet := generateTestPlanetWithShip()
		req := models.FleetOrder{
			Destination: models.Coordinate{
				Galaxy:      1,
				SolarSystem: 3,
				Position:    4,
			},
			Ships: []models.FleetShip{
				{Ship: uuid.New(), Count: 3},
				{Ship: uuid.New(), Count: 4},
			},
		}

		universe := models.Universe{
			Id: uuid.New(),
			Topology: models.UniverseTopology{
				Galaxies:     2,
				SolarSystems: 10,
				Orbits:       14,
			},
			Ships: []models.Ship{
				{Id: req.Ships[0].Ship},
				{Id: req.Ships[1].Ship},
			},
		}

		fleet := models.Fleet{
			Id:      uuid.New(),
			Source:  planet.Id,
			Mission: models.MissionColonize,
		}

		suite.mockUniverseRepo.EXPECT().
			GetByPlanetId(gomock.Any(), planet.Id).
			Times(1).
			Return(universe, nil)
		suite.mockClock.EXPECT().Now(gomock.Any()).Times(1).Return(t2)
		fleetResult := models.SimulationStepResult{
			Resolution: models.FleetEventResolved,
			Fleet: &models.Fleet{
				Id: uuid.New(),
			},
		}
		suite.mockExecutor.EXPECT().
			Execute(gomock.Any(), planet.Id, t2, gomock.Any()).
			Times(1).
			Return(fleetResult, nil)
		result := models.SimulationStepResult{
			Resolution: models.FleetCreated,
			Planet:     &planet,
			Fleet:      &fleet,
		}
		suite.mockExecutor.EXPECT().
			Execute(gomock.Any(), planet.Id, t2, gomock.Any()).
			Times(1).
			Return(result, nil)

		actual, err := suite.usecase.Create(t.Context(), planet.Id, req)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, fleet, actual)
	})

	t.Run("deduct requested ships from planet", func(t *testing.T) {
		suite := setupCreateFleetTestSuite(t)

		planet := generateTestPlanetWithShip()
		planet.Coordinate = models.Coordinate{Galaxy: 1, SolarSystem: 2, Position: 3}
		planet.Ships = []models.PlanetShip{
			{Ship: uuid.New(), Count: 7},
			{Ship: uuid.New(), Count: 6},
		}
		req := models.FleetOrder{
			Mission: models.MissionColonize,
			Destination: models.Coordinate{
				Galaxy:      1,
				SolarSystem: 3,
				Position:    4,
			},
			Ships: []models.FleetShip{
				{Ship: planet.Ships[0].Ship, Count: 3},
				{Ship: planet.Ships[1].Ship, Count: 4},
			},
		}

		universe := models.Universe{
			Id: uuid.New(),
			Topology: models.UniverseTopology{
				Galaxies:     2,
				SolarSystems: 10,
				Orbits:       14,
			},
			Ships: []models.Ship{
				{
					Id:        planet.Ships[0].Ship,
					BaseSpeed: 3600,
				},
				{
					Id:        planet.Ships[1].Ship,
					BaseSpeed: 7000,
				},
			},
		}

		suite.mockUniverseRepo.EXPECT().
			GetByPlanetId(gomock.Any(), planet.Id).
			Times(1).
			Return(universe, nil)
		suite.mockClock.EXPECT().Now(gomock.Any()).Times(1).Return(t2)
		suite.mockExecutor.EXPECT().
			Execute(gomock.Any(), planet.Id, t2, gomock.Any()).
			Times(1).
			DoAndReturn(generatePlanetSimulatorRecorderMock(&planet))

		actual, err := suite.usecase.Create(t.Context(), planet.Id, req)
		require.NoError(t, err, "Actual err: %v", err)

		flightTime := 9762 * time.Second

		expected := models.Fleet{
			Id:      actual.Id,
			Player:  planet.Player,
			Source:  planet.Id,
			Mission: models.MissionColonize,
			Destination: models.FleetDestination{
				Coordinate: models.Coordinate{
					Galaxy:      req.Destination.Galaxy,
					SolarSystem: req.Destination.SolarSystem,
					Position:    req.Destination.Position,
				},
			},
			Ships: []models.FleetShip{
				{Ship: req.Ships[0].Ship, Count: req.Ships[0].Count},
				{Ship: req.Ships[1].Ship, Count: req.Ships[1].Count},
			},
			Schedule: models.FleetFlightSchedule{
				ArrivalAt: t2.Add(flightTime),
				ReturnAt:  t2.Add(2 * flightTime),
				Status:    models.StatusOutbound,
			},
			CreatedAt: t2,
			UpdatedAt: t2,
			Version:   0,
		}
		assert.Equal(t, expected, actual)
	})

	t.Run("returns error when destination coordinates are out of bounds", func(t *testing.T) {
		suite := setupCreateFleetTestSuite(t)

		planet := models.Planet{Id: uuid.New()}

		req := models.FleetOrder{
			Destination: models.Coordinate{
				Galaxy:      0,
				SolarSystem: 1,
				Position:    5,
			},
			Ships: []models.FleetShip{
				{Ship: uuid.New(), Count: 3},
			},
		}

		universe := models.Universe{
			Id: uuid.New(),
			Topology: models.UniverseTopology{
				Galaxies:     1,
				SolarSystems: 1,
				Orbits:       10,
			},
			Ships: []models.Ship{
				{Id: req.Ships[0].Ship},
			},
		}

		suite.mockUniverseRepo.EXPECT().
			GetByPlanetId(gomock.Any(), planet.Id).
			Times(1).
			Return(universe, nil)

		_, err := suite.usecase.Create(t.Context(), planet.Id, req)

		assert.ErrorIs(t, err, domainerrors.ErrCoordinatesOutOfBound, "Actual err: %v", err)
	})

	t.Run("returns error when fleet ship does not exist", func(t *testing.T) {
		suite := setupCreateFleetTestSuite(t)

		planet := models.Planet{Id: uuid.New()}

		req := models.FleetOrder{
			Destination: models.Coordinate{
				Galaxy:      0,
				SolarSystem: 0,
				Position:    5,
			},
			Ships: []models.FleetShip{
				{Ship: uuid.New(), Count: 3},
			},
		}

		universe := models.Universe{
			Id: uuid.New(),
			Topology: models.UniverseTopology{
				Galaxies:     1,
				SolarSystems: 1,
				Orbits:       10,
			},
			Ships: []models.Ship{
				{Id: uuid.New(), BaseSpeed: 10},
			},
		}

		suite.mockUniverseRepo.EXPECT().
			GetByPlanetId(gomock.Any(), planet.Id).
			Times(1).
			Return(universe, nil)

		_, err := suite.usecase.Create(t.Context(), planet.Id, req)

		assert.ErrorIs(t, err, domainerrors.ErrShipNotFound, "Actual err: %v", err)
	})

	t.Run("returns error when dispatcher fails", func(t *testing.T) {
		suite := setupCreateFleetTestSuite(t)

		planet := models.Planet{Id: uuid.New()}

		req := models.FleetOrder{
			Destination: models.Coordinate{
				Galaxy:      1,
				SolarSystem: 3,
				Position:    4,
			},
			Ships: []models.FleetShip{
				{Ship: uuid.New(), Count: 3},
				{Ship: uuid.New(), Count: 4},
			},
		}

		universe := models.Universe{
			Id: uuid.New(),
			Topology: models.UniverseTopology{
				Galaxies:     2,
				SolarSystems: 10,
				Orbits:       14,
			},
			Ships: []models.Ship{
				{Id: req.Ships[0].Ship},
				{Id: req.Ships[1].Ship},
			},
		}

		suite.mockUniverseRepo.EXPECT().
			GetByPlanetId(gomock.Any(), planet.Id).
			Times(1).
			Return(universe, nil)
		suite.mockClock.EXPECT().Now(gomock.Any()).Times(1).Return(t2)
		errSample := errors.New("stubbed error")
		suite.mockExecutor.EXPECT().
			Execute(gomock.Any(), planet.Id, t2, gomock.Any()).
			Times(1).
			Return(models.SimulationStepResult{}, errSample)

		_, err := suite.usecase.Create(t.Context(), planet.Id, req)
		assert.ErrorIs(t, err, errSample, "Actual err: %v", err)
	})

	t.Run("returns error when planet does not exist", func(t *testing.T) {
		suite := setupCreateFleetTestSuite(t)

		suite.mockUniverseRepo.EXPECT().
			GetByPlanetId(gomock.Any(), gomock.Any()).
			Times(1).
			Return(models.Universe{}, domainerrors.ErrNotFound)

		req := models.FleetOrder{}
		_, err := suite.usecase.Create(t.Context(), uuid.New(), req)

		assert.ErrorIs(t, err, domainerrors.ErrNotFound, "Actual err: %v", err)
	})
}

func setupCreateFleetTestSuite(t *testing.T) *createFleetTestSuite {
	t.Helper()

	ctrl := gomock.NewController(t)
	mockUniverseRepo := drivenportstest.NewMockForFetchingUniverses(ctrl)
	mockExecutor := drivenportstest.NewMockForExecutingSimulationStep(ctrl)
	mockClock := drivenportstest.NewMockForFetchingTime(ctrl)

	return &createFleetTestSuite{
		ctrl:             ctrl,
		mockUniverseRepo: mockUniverseRepo,
		mockExecutor:     mockExecutor,
		mockClock:        mockClock,
		usecase: NewCreateFleetUseCase(
			mockUniverseRepo,
			mockExecutor,
			mockClock,
		),
	}
}
