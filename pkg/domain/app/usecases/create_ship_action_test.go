package usecases

import (
	"testing"
	"time"

	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	domainerrors "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models/errors"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models/request"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/usecases/drivenportstest"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

type createShipActionTestSuite struct {
	ctrl         *gomock.Controller
	mockShipRepo *drivenportstest.MockForFetchingShip
	mockMutator  *drivenportstest.MockForMutatingPlanet
	mockClock    *drivenportstest.MockForFetchingTime
	usecase      *CreateShipActionUseCase
}

func TestUnit_CreateShipAction_Create(t *testing.T) {
	t.Run("persists created ship action", func(t *testing.T) {
		suite := setupCreateShipActionTestSuite(t)

		planet := generateTestPlanetWithShip()
		ship := generateTestShip(planet)
		request := generateTestShipActionRequest(planet)

		suite.mockClock.EXPECT().Now(gomock.Any()).Times(1).Return(t2)
		suite.mockShipRepo.EXPECT().
			Get(gomock.Any(), ship.Id).
			Times(1).
			Return(ship, nil)
		suite.mockMutator.EXPECT().
			Mutate(gomock.Any(), planet.Id, gomock.Any()).
			Times(1).
			DoAndReturn(generateApplyingMutatorMock(&planet))

		actual, err := suite.usecase.Create(t.Context(), request)
		require.NoError(t, err, "Actual err: %v", err)

		completionTime := 168480 * time.Millisecond

		expected := models.ShipAction{
			Id:                 actual.Id,
			Ship:               request.Ship,
			Count:              request.Count,
			CreatedAt:          t2,
			NextCompletionAt:   t2.Add(completionTime),
			UnitCompletionTime: completionTime,
			Costs: []models.ShipActionCost{
				{
					Resource: metalResourceId,
					Amount:   100,
				},
				{
					Resource: crystalResourceId,
					Amount:   134,
				},
			},
		}
		assert.Equal(t, expected, actual)
		assert.Equal(t, []models.ShipAction{expected}, planet.ShipActions)
	})

	t.Run("updates planet to current time", func(t *testing.T) {
		suite := setupCreateShipActionTestSuite(t)

		planet := generateTestPlanetWithShip()
		ship := generateTestShip(planet)
		request := generateTestShipActionRequest(planet)

		initialVersion := planet.Version

		suite.mockClock.EXPECT().Now(gomock.Any()).Times(1).Return(t2)
		suite.mockShipRepo.EXPECT().
			Get(gomock.Any(), ship.Id).
			Times(1).
			Return(ship, nil)
		suite.mockMutator.EXPECT().
			Mutate(gomock.Any(), planet.Id, gomock.Any()).
			Times(1).
			DoAndReturn(generateApplyingMutatorMock(&planet))

		_, err := suite.usecase.Create(t.Context(), request)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, t2, planet.UpdatedAt)
		// One bump due to the update to the current time, one bump for the action
		assert.Equal(t, initialVersion+2, planet.Version)
	})

	t.Run("applies completed action and create a new one", func(t *testing.T) {
		suite := setupCreateShipActionTestSuite(t)

		planet := generateTestPlanetWithShipAction(duration)
		ship := generateTestShip(planet)
		request := generateTestShipActionRequest(planet)

		initialVersion := planet.Version
		initialCount := planet.Ships[0].Count

		suite.mockClock.EXPECT().Now(gomock.Any()).Times(1).Return(t3)
		suite.mockShipRepo.EXPECT().
			Get(gomock.Any(), ship.Id).
			Times(1).
			Return(ship, nil)
		suite.mockMutator.EXPECT().
			Mutate(gomock.Any(), planet.Id, gomock.Any()).
			Times(1).
			DoAndReturn(generateApplyingMutatorMock(&planet))

		actual, err := suite.usecase.Create(t.Context(), request)
		require.NoError(t, err, "Actual err: %v", err)

		completionTime := 168480 * time.Millisecond
		expected := models.Planet{
			Id:        planet.Id,
			Fields:    100,
			CreatedAt: t1,
			UpdatedAt: t3,
			// Update to current time, action completion, update to current time and
			// action creation
			Version: initialVersion + 4,
			Resources: []models.PlanetResource{
				{
					Resource: metalResourceId,
					Amount:   99899,
				},
				{
					Resource: crystalResourceId,
					Amount:   99865,
				},
			},
			Ships: []models.PlanetShip{
				{
					Ship:  request.Ship,
					Count: initialCount,
				},
			},
			ShipActions: []models.ShipAction{
				{
					Id:                 actual.Id,
					Ship:               request.Ship,
					Count:              request.Count,
					CreatedAt:          t3,
					NextCompletionAt:   t3.Add(completionTime),
					UnitCompletionTime: completionTime,
					Costs: []models.ShipActionCost{
						{
							Resource: metalResourceId,
							Amount:   100,
						},
						{
							Resource: crystalResourceId,
							Amount:   134,
						},
					},
				},
			},
		}
		assert.Equal(t, expected, planet)
	})

	t.Run("returns error when planet is deleted", func(t *testing.T) {
		suite := setupCreateShipActionTestSuite(t)

		planet := generateTestPlanetWithShip()
		ship := generateTestShip(planet)
		request := generateTestShipActionRequest(planet)

		suite.mockClock.EXPECT().Now(gomock.Any()).Times(1).Return(t2)
		suite.mockShipRepo.EXPECT().
			Get(gomock.Any(), ship.Id).
			Times(1).
			Return(ship, nil)
		suite.mockMutator.EXPECT().
			Mutate(gomock.Any(), planet.Id, gomock.Any()).
			Times(1).
			Return(models.PlanetMutationResult{Deleted: true}, nil)

		_, err := suite.usecase.Create(t.Context(), request)

		assert.ErrorIs(t, err, domainerrors.ErrNotFound, "Actual err: %v", err)
	})

	t.Run("appends action when planet already has an action running", func(t *testing.T) {
		suite := setupCreateShipActionTestSuite(t)

		planet := generateTestPlanetWithShipAction(duration)
		ship := generateTestShip(planet)
		request := generateTestShipActionRequest(planet)

		initialVersion := planet.Version
		initialCount := planet.Ships[0].Count
		initialAction := planet.ShipActions[0]

		suite.mockClock.EXPECT().Now(gomock.Any()).Times(1).Return(t2)
		suite.mockShipRepo.EXPECT().
			Get(gomock.Any(), ship.Id).
			Times(1).
			Return(ship, nil)
		suite.mockMutator.EXPECT().
			Mutate(gomock.Any(), planet.Id, gomock.Any()).
			Times(1).
			DoAndReturn(generateApplyingMutatorMock(&planet))

		actual, err := suite.usecase.Create(t.Context(), request)
		require.NoError(t, err, "Actual err: %v", err)

		completionTime := 168480 * time.Millisecond
		expected := models.Planet{
			Id:        planet.Id,
			Fields:    100,
			CreatedAt: t1,
			UpdatedAt: t2,
			// Update to current time and action creation
			Version: initialVersion + 2,
			Resources: []models.PlanetResource{
				{
					Resource: metalResourceId,
					Amount:   99899,
				},
				{
					Resource: crystalResourceId,
					Amount:   99865,
				},
			},
			Ships: []models.PlanetShip{
				{
					Ship:  request.Ship,
					Count: initialCount,
				},
			},
			ShipActions: []models.ShipAction{
				initialAction,
				{
					Id:                 actual.Id,
					Ship:               request.Ship,
					Count:              request.Count,
					CreatedAt:          t3,
					NextCompletionAt:   t3.Add(completionTime),
					UnitCompletionTime: completionTime,
					Costs: []models.ShipActionCost{
						{
							Resource: metalResourceId,
							Amount:   100,
						},
						{
							Resource: crystalResourceId,
							Amount:   134,
						},
					},
				},
			},
		}
		assert.Equal(t, expected, planet)
	})

	t.Run("returns error when planet does not contain requested ship", func(t *testing.T) {
		suite := setupCreateShipActionTestSuite(t)

		planet := models.Planet{Id: uuid.New()}
		ship := models.Ship{Id: uuid.New()}
		req := request.ShipActionCreationRequest{
			Planet: planet.Id,
			Ship:   ship.Id,
			Count:  1,
		}

		suite.mockClock.EXPECT().Now(gomock.Any()).Times(1).Return(t2)
		suite.mockShipRepo.EXPECT().
			Get(gomock.Any(), ship.Id).
			Times(1).
			Return(ship, nil)
		suite.mockMutator.EXPECT().
			Mutate(gomock.Any(), planet.Id, gomock.Any()).
			Times(1).
			DoAndReturn(generateApplyingMutatorMock(&planet))

		_, err := suite.usecase.Create(t.Context(), req)

		assert.ErrorIs(t, err, domainerrors.ErrShipNotFound, "Actual err: %v", err)
	})

	t.Run("returns error when ship does not exist", func(t *testing.T) {
		suite := setupCreateShipActionTestSuite(t)

		planet := models.Planet{Id: uuid.New()}
		shipId := uuid.New()
		req := request.ShipActionCreationRequest{
			Planet: planet.Id,
			Ship:   shipId,
			Count:  36,
		}

		suite.mockClock.EXPECT().Now(gomock.Any()).Times(1).Return(t2)
		suite.mockShipRepo.EXPECT().
			Get(gomock.Any(), shipId).
			Times(1).
			Return(models.Ship{}, domainerrors.ErrShipNotFound)

		_, err := suite.usecase.Create(t.Context(), req)

		assert.ErrorIs(t, err, domainerrors.ErrShipNotFound, "Actual err: %v", err)
	})
}

func setupCreateShipActionTestSuite(t *testing.T) *createShipActionTestSuite {
	t.Helper()

	ctrl := gomock.NewController(t)
	mockShipRepo := drivenportstest.NewMockForFetchingShip(ctrl)
	mockMutator := drivenportstest.NewMockForMutatingPlanet(ctrl)
	mockClock := drivenportstest.NewMockForFetchingTime(ctrl)

	return &createShipActionTestSuite{
		ctrl:         ctrl,
		mockShipRepo: mockShipRepo,
		mockMutator:  mockMutator,
		mockClock:    mockClock,
		usecase: NewCreateShipActionUseCase(
			mockShipRepo,
			mockMutator,
			mockClock,
		),
	}
}

func generateTestPlanetWithShip() models.Planet {
	return models.Planet{
		Id:     uuid.New(),
		Fields: 100,
		Resources: []models.PlanetResource{
			{
				Resource: metalResourceId,
				Amount:   99999,
			},
			{
				Resource: crystalResourceId,
				Amount:   99999,
			},
		},
		Ships: []models.PlanetShip{
			{
				Ship:  uuid.New(),
				Count: 3,
			},
		},
		CreatedAt: t1,
		UpdatedAt: t1,
		Version:   2,
	}
}

func generateTestPlanetWithShipAction(completionTime time.Duration) models.Planet {
	p := generateTestPlanetWithShip()
	p.ShipActions = append(p.ShipActions, models.ShipAction{
		Id:                 uuid.New(),
		Ship:               p.Ships[0].Ship,
		Count:              37,
		CreatedAt:          t1,
		UnitCompletionTime: completionTime,
	})

	return p
}

func generateTestShip(planet models.Planet) models.Ship {
	return models.Ship{
		Id: planet.Ships[0].Ship,
		Costs: []models.ShipCost{
			{
				Resource:              metalResourceId,
				Cost:                  50,
				BuildTimeHoursPerUnit: 0.0004,
			},
			{
				Resource:              crystalResourceId,
				Cost:                  67,
				BuildTimeHoursPerUnit: 0.0004,
			},
		},
	}
}

func generateTestShipActionRequest(
	planet models.Planet,
) request.ShipActionCreationRequest {
	return request.ShipActionCreationRequest{
		Planet: planet.Id,
		Ship:   planet.Ships[0].Ship,
		Count:  2,
	}
}
