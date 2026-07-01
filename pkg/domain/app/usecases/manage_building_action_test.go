package usecases

import (
	"context"
	"testing"
	"time"

	"github.com/Knoblauchpilze/backend-toolkit/pkg/errors"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	domainerrors "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models/errors"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models/request"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/usecases/drivenportstest"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

var (
	metalResourceId   = uuid.MustParse("b4419b6b-b3bf-4576-aa92-055283addbc8")
	crystalResourceId = uuid.MustParse("cd2ac9aa-9968-4ff5-b746-88f1f810fbb3")

	someTime = time.Date(2026, time.June, 25, 22, 22, 49, 0, time.UTC)
)

func TestUnit_ManageBuildingAction_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockActionRepo := drivenportstest.NewMockForManagingBuildingActions(ctrl)
	mockPlanetRepo := drivenportstest.NewMockForManagingPlanets(ctrl)
	mockBuildingRepo := drivenportstest.NewMockForFetchingBuilding(ctrl)

	t.Run("persists created building action", func(t *testing.T) {
		planet := generateTestPlanet()
		building := generateTestBuilding(planet)
		request := generateTestBuildingActionRequest(planet)

		mockPlanetRepo.EXPECT().
			Get(gomock.Any(), gomock.Eq(request.Planet)).
			Times(1).
			Return(planet, nil)

		mockBuildingRepo.EXPECT().
			Get(gomock.Any(), gomock.Eq(request.Building)).
			Times(1).
			Return(building, nil)

		var captured models.Planet
		mockActionRepo.EXPECT().
			Create(gomock.Any(), gomock.AssignableToTypeOf(captured)).
			Times(1).
			DoAndReturn(func(ctx context.Context, planet models.Planet) error {
				captured = planet
				return nil
			})

		completionTime := 289440 * time.Millisecond

		usecase := NewBuildingActionUseCase(mockActionRepo, mockPlanetRepo, mockBuildingRepo)
		actual, err := usecase.Create(context.Background(), request)
		require.NoError(t, err, "Actual err: %v", err)

		expected := models.BuildingAction{
			Id:           actual.Id,
			Planet:       request.Planet,
			Building:     request.Building,
			DesiredLevel: planet.Buildings[0].Level + 1,
			CreatedAt:    actual.CreatedAt,
			CompletedAt:  actual.CreatedAt.Add(completionTime),
			Costs: []models.BuildingActionCost{
				{
					Resource: metalResourceId,
					Amount:   78,
				},
				{
					Resource: crystalResourceId,
					Amount:   123,
				},
			},
			Storages:    []models.BuildingActionResourceStorage{},
			Productions: []models.BuildingActionResourceProduction{},
		}
		assert.Equal(t, expected, actual)
		assert.Equal(t, &expected, captured.BuildingAction)
	})

	t.Run("persists modified planet", func(t *testing.T) {
		planet := generateTestPlanet()
		building := generateTestBuilding(planet)
		request := generateTestBuildingActionRequest(planet)

		mockPlanetRepo.EXPECT().
			Get(gomock.Any(), gomock.Eq(request.Planet)).
			Times(1).
			Return(planet, nil)

		mockBuildingRepo.EXPECT().
			Get(gomock.Any(), gomock.Eq(request.Building)).
			Times(1).
			Return(building, nil)

		var captured models.Planet
		mockActionRepo.EXPECT().
			Create(gomock.Any(), gomock.AssignableToTypeOf(captured)).
			Times(1).
			DoAndReturn(func(ctx context.Context, planet models.Planet) error {
				captured = planet
				return nil
			})

		usecase := NewBuildingActionUseCase(mockActionRepo, mockPlanetRepo, mockBuildingRepo)
		actual, err := usecase.Create(context.Background(), request)
		require.NoError(t, err, "Actual err: %v", err)

		expected := models.Planet{
			Id:        planet.Id,
			Player:    planet.Player,
			Name:      planet.Name,
			Homeworld: planet.Homeworld,
			CreatedAt: planet.CreatedAt,
			UpdatedAt: actual.CreatedAt,
			Version:   planet.Version + 1,
			Resources: []models.PlanetResource{
				{
					Resource: metalResourceId,
					// Corresponds to the initial amount (99999) minus the cost
					// for each resource (78 for metal and 123 for crystal)
					Amount: 99921,
				},
				{
					Resource: crystalResourceId,
					Amount:   99876,
				},
			},
			Storages:       planet.Storages,
			Productions:    planet.Productions,
			Buildings:      planet.Buildings,
			BuildingAction: &actual,
		}
		assert.Equal(t, planet.Id, actual.Planet)
		assert.Equal(t, expected, captured)
	})

	t.Run("returns error when planet is not found", func(t *testing.T) {
		planet := generateTestPlanet()
		request := generateTestBuildingActionRequest(planet)

		mockPlanetRepo.EXPECT().
			Get(gomock.Any(), gomock.Eq(request.Planet)).
			Times(1).
			Return(models.Planet{}, domainerrors.ErrNotFound)

		usecase := NewBuildingActionUseCase(mockActionRepo, mockPlanetRepo, mockBuildingRepo)
		_, err := usecase.Create(context.Background(), request)

		assert.Equal(t, domainerrors.ErrNotFound, err, "Actual err: %v", err)
	})

	t.Run("returns error when building is not found", func(t *testing.T) {
		planet := generateTestPlanet()
		request := generateTestBuildingActionRequest(planet)

		mockPlanetRepo.EXPECT().
			Get(gomock.Any(), gomock.Eq(request.Planet)).
			Times(1).
			Return(planet, nil)

		mockBuildingRepo.EXPECT().
			Get(gomock.Any(), gomock.Eq(request.Building)).
			Times(1).
			Return(models.Building{}, domainerrors.ErrNotFound)

		usecase := NewBuildingActionUseCase(mockActionRepo, mockPlanetRepo, mockBuildingRepo)
		_, err := usecase.Create(context.Background(), request)

		assert.Equal(t, domainerrors.ErrBuildingNotFound, err, "Actual err: %v", err)
	})

	t.Run("returns error when repository fails", func(t *testing.T) {
		planet := generateTestPlanet()
		building := generateTestBuilding(planet)
		request := generateTestBuildingActionRequest(planet)

		mockPlanetRepo.EXPECT().
			Get(gomock.Any(), gomock.Eq(request.Planet)).
			Times(1).
			Return(planet, nil)

		mockBuildingRepo.EXPECT().
			Get(gomock.Any(), gomock.Eq(request.Building)).
			Times(1).
			Return(building, nil)

		expectedErr := errors.New("stubbed error")
		mockActionRepo.EXPECT().
			Create(gomock.Any(), gomock.Any()).
			Times(1).
			Return(expectedErr)

		usecase := NewBuildingActionUseCase(mockActionRepo, mockPlanetRepo, mockBuildingRepo)
		_, err := usecase.Create(context.Background(), request)

		assert.ErrorIs(t, expectedErr, err, "Actual err: %v", err)
	})
}

func TestUnit_ManageBuildingAction_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockActionRepo := drivenportstest.NewMockForManagingBuildingActions(ctrl)
	mockPlanetRepo := drivenportstest.NewMockForManagingPlanets(ctrl)
	mockBuildingRepo := drivenportstest.NewMockForFetchingBuilding(ctrl)

	t.Run("deletes existing building action", func(t *testing.T) {
		actionId := uuid.New()
		planet := models.Planet{
			Id: uuid.New(),
			BuildingAction: &models.BuildingAction{
				Id: actionId,
			},
		}

		mockPlanetRepo.EXPECT().
			GetByAction(gomock.Any(), gomock.Eq(actionId)).
			Times(1).
			Return(planet, nil)
		var captured models.Planet
		mockActionRepo.EXPECT().
			Delete(gomock.Any(), gomock.AssignableToTypeOf(captured)).
			Times(1).
			DoAndReturn(func(ctx context.Context, planet models.Planet) error {
				captured = planet
				return nil
			})

		usecase := NewBuildingActionUseCase(mockActionRepo, mockPlanetRepo, mockBuildingRepo)
		err := usecase.Delete(context.Background(), actionId)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Nil(t, captured.BuildingAction)
	})

	t.Run("persists modified planet", func(t *testing.T) {
		actionId := uuid.New()
		planet := models.Planet{
			Id: uuid.New(),
			Resources: []models.PlanetResource{
				{
					Resource: metalResourceId,
					Amount:   100,
				},
				{
					Resource: crystalResourceId,
					Amount:   6000,
				},
			},
			BuildingAction: &models.BuildingAction{
				Id: actionId,
				Costs: []models.BuildingActionCost{
					{
						Resource: metalResourceId,
						Amount:   47891,
					},
					{
						Resource: crystalResourceId,
						Amount:   134876,
					},
				},
			},
			UpdatedAt: someTime,
		}

		mockPlanetRepo.EXPECT().
			GetByAction(gomock.Any(), gomock.Eq(actionId)).
			Times(1).
			Return(planet, nil)
		var captured models.Planet
		mockActionRepo.EXPECT().
			Delete(gomock.Any(), gomock.AssignableToTypeOf(captured)).
			Times(1).
			DoAndReturn(func(ctx context.Context, planet models.Planet) error {
				captured = planet
				return nil
			})

		usecase := NewBuildingActionUseCase(mockActionRepo, mockPlanetRepo, mockBuildingRepo)
		err := usecase.Delete(context.Background(), actionId)
		require.NoError(t, err, "Actual err: %v", err)

		expectedPlanet := models.Planet{
			Id:        planet.Id,
			Version:   planet.Version + 1,
			UpdatedAt: someTime,
			Resources: []models.PlanetResource{
				{
					Resource: metalResourceId,
					Amount:   47991,
				},
				{
					Resource: crystalResourceId,
					Amount:   140876,
				},
			},
			BuildingAction: nil,
		}
		assert.Equal(t, expectedPlanet, captured)
	})

	t.Run("succeeds when planet has no building action", func(t *testing.T) {
		planet := models.Planet{
			Id:             uuid.New(),
			BuildingAction: nil,
		}
		actionId := uuid.New()

		mockPlanetRepo.EXPECT().
			GetByAction(gomock.Any(), gomock.Eq(actionId)).
			Times(1).
			Return(planet, nil)

		usecase := NewBuildingActionUseCase(mockActionRepo, mockPlanetRepo, mockBuildingRepo)
		err := usecase.Delete(context.Background(), actionId)
		require.NoError(t, err, "Actual err: %v", err)
	})

	t.Run("succeeds when building action is not found", func(t *testing.T) {
		actionId := uuid.New()

		mockPlanetRepo.EXPECT().
			GetByAction(gomock.Any(), gomock.Eq(actionId)).
			Times(1).
			Return(models.Planet{}, domainerrors.ErrNotFound)

		usecase := NewBuildingActionUseCase(mockActionRepo, mockPlanetRepo, mockBuildingRepo)
		err := usecase.Delete(context.Background(), actionId)
		require.NoError(t, err, "Actual err: %v", err)
	})

	t.Run("returns error when repository fails", func(t *testing.T) {
		actionId := uuid.New()
		planet := models.Planet{
			Id: uuid.New(),
			BuildingAction: &models.BuildingAction{
				Id: actionId,
			},
		}

		mockPlanetRepo.EXPECT().
			GetByAction(gomock.Any(), gomock.Eq(actionId)).
			Times(1).
			Return(planet, nil)
		expectedErr := errors.New("stubbed error")
		mockActionRepo.EXPECT().
			Delete(gomock.Any(), gomock.Any()).
			Times(1).
			Return(expectedErr)

		usecase := NewBuildingActionUseCase(mockActionRepo, mockPlanetRepo, mockBuildingRepo)
		err := usecase.Delete(context.Background(), actionId)

		assert.ErrorIs(t, expectedErr, err, "Actual err: %v", err)
	})
}

func generateTestPlanet() models.Planet {
	return models.Planet{
		Id: uuid.New(),
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
		Buildings: []models.PlanetBuilding{
			{
				Building: uuid.New(),
				Level:    2,
			},
		},
	}
}

func generateTestBuilding(planet models.Planet) models.Building {
	return models.Building{
		Id: planet.Buildings[0].Building,
		Costs: []models.BuildingCost{
			{
				Resource: metalResourceId,
				Cost:     50,
				Progress: 1.25,
			},
			{
				Resource: crystalResourceId,
				Cost:     67,
				Progress: 1.36,
			},
		},
	}
}

func generateTestBuildingActionRequest(
	planet models.Planet,
) request.BuildingActionCreationRequest {
	return request.BuildingActionCreationRequest{
		Planet:   planet.Id,
		Building: planet.Buildings[0].Building,
	}
}
