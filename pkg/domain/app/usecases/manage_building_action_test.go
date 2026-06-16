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
)

func TestUnit_ManageBuildingAction_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockActionRepo := drivenportstest.NewMockForManagingBuildingActions(ctrl)
	mockPlanetRepo := drivenportstest.NewMockForManagingPlanets(ctrl)
	mockBuildingRepo := drivenportstest.NewMockForListingBuildings(ctrl)

	request := request.BuildingActionCreationRequest{
		Planet:   uuid.New(),
		Building: uuid.New(),
	}

	planet := models.Planet{
		Id: request.Planet,
		Buildings: []models.PlanetBuilding{
			{
				Building: request.Building,
				Level:    2,
			},
		},
	}

	building := models.Building{
		Id: request.Building,
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

	t.Run("persists created building action", func(t *testing.T) {
		mockPlanetRepo.EXPECT().
			Get(gomock.Any(), gomock.Eq(request.Planet)).
			Times(1).
			Return(planet, nil)

		mockBuildingRepo.EXPECT().
			Get(gomock.Any(), gomock.Eq(request.Building)).
			Times(1).
			Return(building, nil)

		var captured models.BuildingAction
		mockActionRepo.EXPECT().
			Create(gomock.Any(), gomock.AssignableToTypeOf(captured)).
			Times(1).
			DoAndReturn(func(ctx context.Context, action models.BuildingAction) error {
				captured = action
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
			CurrentLevel: planet.Buildings[0].Level,
			DesiredLevel: planet.Buildings[0].Level + 1,
			CreatedAt:    actual.CreatedAt,
			CompletedAt:  actual.CreatedAt.Add(completionTime),
			Version:      0,
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
		assert.Equal(t, expected, captured)
	})

	t.Run("returns error when planet is not found", func(t *testing.T) {
		mockPlanetRepo.EXPECT().
			Get(gomock.Any(), gomock.Eq(request.Planet)).
			Times(1).
			Return(models.Planet{}, domainerrors.ErrNotFound)

		usecase := NewBuildingActionUseCase(mockActionRepo, mockPlanetRepo, mockBuildingRepo)
		_, err := usecase.Create(context.Background(), request)

		assert.Equal(t, domainerrors.ErrNotFound, err, "Actual err: %v", err)
	})

	t.Run("returns error when building is not found", func(t *testing.T) {
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
	mockBuildingRepo := drivenportstest.NewMockForListingBuildings(ctrl)

	t.Run("deletes existing building action", func(t *testing.T) {
		id := uuid.New()

		mockActionRepo.EXPECT().
			Delete(gomock.Any(), gomock.Eq(id)).
			Times(1).
			Return(nil)

		usecase := NewBuildingActionUseCase(mockActionRepo, mockPlanetRepo, mockBuildingRepo)
		err := usecase.Delete(context.Background(), id)
		require.NoError(t, err, "Actual err: %v", err)
	})

	t.Run("returns error when repository fails", func(t *testing.T) {
		expectedErr := errors.New("stubbed error")
		mockActionRepo.EXPECT().
			Delete(gomock.Any(), gomock.Any()).
			Times(1).
			Return(expectedErr)

		usecase := NewBuildingActionUseCase(mockActionRepo, mockPlanetRepo, mockBuildingRepo)
		err := usecase.Delete(context.Background(), uuid.New())

		assert.ErrorIs(t, expectedErr, err, "Actual err: %v", err)
	})
}
