package usecases

import (
	"context"
	"errors"
	"testing"

	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	domainerrors "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models/errors"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models/request"
	drivenports "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/ports/driven"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/usecases/drivenportstest"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

type FleetDispatcherMock func(context.Context, request.FleetCreationRequest, drivenports.FleetCreator) (models.Fleet, error)

type createFleetTestSuite struct {
	ctrl             *gomock.Controller
	mockUniverseRepo *drivenportstest.MockForFetchingUniverses
	mockDispatcher   *drivenportstest.MockForDispatchingFleet
	mockClock        *drivenportstest.MockForFetchingTime
	usecase          *CreateFleetUseCase
}

func TestUnit_CreateFleet_Create(t *testing.T) {
	t.Run("forwards fleet returned by the dispatcher", func(t *testing.T) {
		suite := setupCreateFleetTestSuite(t)

		req := request.FleetCreationRequest{
			Planet: uuid.New(),
			Destination: request.FleetDestinationRequest{
				Galaxy:      1,
				SolarSystem: 3,
				Position:    4,
			},
			Ships: []request.FleetShipRequest{
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

		expected := models.Fleet{
			Id:     uuid.New(),
			Player: uuid.New(),
			Source: req.Planet,
			Destination: models.Coordinate{
				Galaxy:      req.Destination.Galaxy,
				SolarSystem: req.Destination.SolarSystem,
				Position:    req.Destination.Position,
			},
			Ships: []models.FleetShip{
				{Ship: req.Ships[0].Ship, Count: req.Ships[0].Count},
				{Ship: req.Ships[1].Ship, Count: req.Ships[1].Count},
			},
			CreatedAt: t2,
			ArrivalAt: t3,
			ReturnAt:  t4,
			UpdatedAt: t2,
			Version:   0,
		}

		suite.mockUniverseRepo.EXPECT().
			GetByPlanetId(gomock.Any(), req.Planet).
			Times(1).
			Return(universe, nil)
		suite.mockClock.EXPECT().Now(gomock.Any()).Times(1).Return(t2)
		suite.mockDispatcher.EXPECT().
			Dispatch(gomock.Any(), req, gomock.Any()).
			Times(1).
			Return(expected, nil)

		actual, err := suite.usecase.Create(t.Context(), req)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, expected, actual)
	})

	t.Run("returns error when destination coordinates are out of bounds", func(t *testing.T) {
		suite := setupCreateFleetTestSuite(t)

		req := request.FleetCreationRequest{
			Planet: uuid.New(),
			Destination: request.FleetDestinationRequest{
				Galaxy:      0,
				SolarSystem: 1,
				Position:    5,
			},
			Ships: []request.FleetShipRequest{
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
			GetByPlanetId(gomock.Any(), req.Planet).
			Times(1).
			Return(universe, nil)

		_, err := suite.usecase.Create(t.Context(), req)

		assert.ErrorIs(t, err, domainerrors.ErrCoordinatesOutOfBound, "Actual err: %v", err)
	})

	t.Run("returns error when dispatcher fails", func(t *testing.T) {
		suite := setupCreateFleetTestSuite(t)

		req := request.FleetCreationRequest{
			Planet: uuid.New(),
			Destination: request.FleetDestinationRequest{
				Galaxy:      1,
				SolarSystem: 3,
				Position:    4,
			},
			Ships: []request.FleetShipRequest{
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
			GetByPlanetId(gomock.Any(), req.Planet).
			Times(1).
			Return(universe, nil)
		suite.mockClock.EXPECT().Now(gomock.Any()).Times(1).Return(t2)
		errSample := errors.New("stubbed error")
		suite.mockDispatcher.EXPECT().
			Dispatch(gomock.Any(), req, gomock.Any()).
			Times(1).
			Return(models.Fleet{}, errSample)

		_, err := suite.usecase.Create(t.Context(), req)
		assert.ErrorIs(t, err, errSample, "Actual err: %v", err)
	})

	t.Run("returns error when planet does not exist", func(t *testing.T) {
		suite := setupCreateFleetTestSuite(t)

		suite.mockUniverseRepo.EXPECT().
			GetByPlanetId(gomock.Any(), gomock.Any()).
			Times(1).
			Return(models.Universe{}, domainerrors.ErrNotFound)

		req := request.FleetCreationRequest{}
		_, err := suite.usecase.Create(t.Context(), req)

		assert.ErrorIs(t, err, domainerrors.ErrNotFound, "Actual err: %v", err)
	})
}

func setupCreateFleetTestSuite(t *testing.T) *createFleetTestSuite {
	t.Helper()

	ctrl := gomock.NewController(t)
	mockUniverseRepo := drivenportstest.NewMockForFetchingUniverses(ctrl)
	mockDispatcher := drivenportstest.NewMockForDispatchingFleet(ctrl)
	mockClock := drivenportstest.NewMockForFetchingTime(ctrl)

	return &createFleetTestSuite{
		ctrl:             ctrl,
		mockUniverseRepo: mockUniverseRepo,
		mockDispatcher:   mockDispatcher,
		mockClock:        mockClock,
		usecase: NewCreateFleetUseCase(
			mockUniverseRepo,
			mockDispatcher,
			mockClock,
		),
	}
}

// generateApplyingDispatcherMock generates a function mock for the fleet dispatcher
// which applies the provided mutator to a known planet.
func generateApplyingDispatcherMock(p *models.Planet) FleetDispatcherMock {
	return func(
		ctx context.Context, req request.FleetCreationRequest, c drivenports.FleetCreator,
	) (models.Fleet, error) {
		return c(p, req)
	}
}
