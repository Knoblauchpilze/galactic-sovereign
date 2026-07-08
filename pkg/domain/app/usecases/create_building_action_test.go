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

type createBuildingActionTestSuite struct {
	ctrl             *gomock.Controller
	mockBuildingRepo *drivenportstest.MockForFetchingBuilding
	mockMutator      *drivenportstest.MockForMutatingPlanet
	mockClock        *drivenportstest.MockForFetchingTime
	usecase          *CreateBuildingActionUseCase
}

func TestUnit_CreateBuildingAction_Create(t *testing.T) {
	t.Run("persists created building action", func(t *testing.T) {
		suite := setupCreateBuildingActionTestSuite(t)

		planet := generateTestPlanet()
		building := generateTestBuilding(planet)
		request := generateTestBuildingActionRequest(planet)

		suite.mockClock.EXPECT().Now(gomock.Any()).Times(1).Return(t2)
		suite.mockBuildingRepo.EXPECT().
			Get(gomock.Any(), building.Id).
			Times(1).
			Return(building, nil)
		suite.mockMutator.EXPECT().
			Mutate(gomock.Any(), planet.Id, gomock.Any()).
			Times(1).
			DoAndReturn(generateApplyingMutatorMock(&planet))

		actual, err := suite.usecase.Create(t.Context(), request)
		require.NoError(t, err, "Actual err: %v", err)

		completionTime := 289440 * time.Millisecond

		expected := models.BuildingAction{
			Id:           actual.Id,
			Building:     request.Building,
			DesiredLevel: planet.Buildings[0].Level + 1,
			CreatedAt:    t2,
			CompletedAt:  t2.Add(completionTime),
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
		assert.Equal(t, &expected, planet.BuildingAction)
	})

	t.Run("updates planet to current time", func(t *testing.T) {
		suite := setupCreateBuildingActionTestSuite(t)

		planet := generateTestPlanet()
		building := generateTestBuilding(planet)
		request := generateTestBuildingActionRequest(planet)

		initialVersion := planet.Version

		suite.mockClock.EXPECT().Now(gomock.Any()).Times(1).Return(t2)
		suite.mockBuildingRepo.EXPECT().
			Get(gomock.Any(), building.Id).
			Times(1).
			Return(building, nil)
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
		suite := setupCreateBuildingActionTestSuite(t)

		planet := generateTestPlanet()
		building := generateTestBuilding(planet)
		planet.BuildingAction = &models.BuildingAction{
			Id:           uuid.New(),
			Building:     building.Id,
			DesiredLevel: planet.Buildings[0].Level + 1,
			CreatedAt:    t1,
			CompletedAt:  t2,
		}
		request := generateTestBuildingActionRequest(planet)

		initialVersion := planet.Version
		initialLevel := planet.Buildings[0].Level

		suite.mockClock.EXPECT().Now(gomock.Any()).Times(1).Return(t3)
		suite.mockBuildingRepo.EXPECT().
			Get(gomock.Any(), building.Id).
			Times(1).
			Return(building, nil)
		suite.mockMutator.EXPECT().
			Mutate(gomock.Any(), planet.Id, gomock.Any()).
			Times(1).
			DoAndReturn(generateApplyingMutatorMock(&planet))

		actual, err := suite.usecase.Create(t.Context(), request)
		require.NoError(t, err, "Actual err: %v", err)

		completionTime := 381600 * time.Millisecond
		expected := models.Planet{
			Id:        planet.Id,
			CreatedAt: t1,
			UpdatedAt: t3,
			// Update to current time, action completion, update to current time and
			// action creation
			Version: initialVersion + 4,
			Resources: []models.PlanetResource{
				{
					Resource: metalResourceId,
					Amount:   99902,
				},
				{
					Resource: crystalResourceId,
					Amount:   99831,
				},
			},
			Buildings: []models.PlanetBuilding{
				{
					Building: request.Building,
					Level:    initialLevel + 1,
				},
			},
			BuildingAction: &models.BuildingAction{
				Id:           actual.Id,
				Building:     request.Building,
				DesiredLevel: 4,
				CreatedAt:    t3,
				CompletedAt:  t3.Add(completionTime),
				Costs: []models.BuildingActionCost{
					{
						Resource: metalResourceId,
						Amount:   97,
					},
					{
						Resource: crystalResourceId,
						Amount:   168,
					},
				},
				Storages:    []models.BuildingActionResourceStorage{},
				Productions: []models.BuildingActionResourceProduction{},
			},
		}
		assert.Equal(t, expected, planet)
	})

	t.Run("returns error when planet is deleted", func(t *testing.T) {
		suite := setupCreateBuildingActionTestSuite(t)

		planet := generateTestPlanet()
		building := generateTestBuilding(planet)
		request := generateTestBuildingActionRequest(planet)

		suite.mockClock.EXPECT().Now(gomock.Any()).Times(1).Return(t2)
		suite.mockBuildingRepo.EXPECT().
			Get(gomock.Any(), building.Id).
			Times(1).
			Return(building, nil)
		suite.mockMutator.EXPECT().
			Mutate(gomock.Any(), planet.Id, gomock.Any()).
			Times(1).
			Return(models.PlanetMutationResult{Deleted: true}, nil)

		_, err := suite.usecase.Create(t.Context(), request)

		assert.ErrorIs(t, err, domainerrors.ErrNotFound, "Actual err: %v", err)
	})

	t.Run("returns error when planet already has an action running", func(t *testing.T) {
		suite := setupCreateBuildingActionTestSuite(t)

		planet := generateTestPlanet()
		planet.UpdatedAt = t1
		planet.Version = 2
		planet.BuildingAction = &models.BuildingAction{
			Id:          uuid.New(),
			Building:    planet.Buildings[0].Building,
			CreatedAt:   t1,
			CompletedAt: t3,
		}
		building := generateTestBuilding(planet)
		request := generateTestBuildingActionRequest(planet)

		suite.mockClock.EXPECT().Now(gomock.Any()).Times(1).Return(t2)
		suite.mockBuildingRepo.EXPECT().
			Get(gomock.Any(), building.Id).
			Times(1).
			Return(building, nil)
		suite.mockMutator.EXPECT().
			Mutate(gomock.Any(), planet.Id, gomock.Any()).
			Times(1).
			DoAndReturn(generateApplyingMutatorMock(&planet))

		_, err := suite.usecase.Create(t.Context(), request)

		assert.ErrorIs(t, err, domainerrors.ErrActionAlreadyInProgress, "Actual err: %v", err)
	})

	t.Run("returns error when planet does not contain requested building", func(t *testing.T) {
		suite := setupCreateBuildingActionTestSuite(t)

		planet := models.Planet{Id: uuid.New()}
		building := models.Building{Id: uuid.New()}
		req := request.BuildingActionCreationRequest{
			Planet:   planet.Id,
			Building: building.Id,
		}

		suite.mockClock.EXPECT().Now(gomock.Any()).Times(1).Return(t2)
		suite.mockBuildingRepo.EXPECT().
			Get(gomock.Any(), building.Id).
			Times(1).
			Return(building, nil)
		suite.mockMutator.EXPECT().
			Mutate(gomock.Any(), planet.Id, gomock.Any()).
			Times(1).
			DoAndReturn(generateApplyingMutatorMock(&planet))

		_, err := suite.usecase.Create(t.Context(), req)

		assert.ErrorIs(t, err, domainerrors.ErrBuildingNotFound, "Actual err: %v", err)
	})

	t.Run("returns error when building does not exist", func(t *testing.T) {
		suite := setupCreateBuildingActionTestSuite(t)

		planet := models.Planet{Id: uuid.New()}
		buildingId := uuid.New()
		req := request.BuildingActionCreationRequest{
			Planet:   planet.Id,
			Building: buildingId,
		}

		suite.mockClock.EXPECT().Now(gomock.Any()).Times(1).Return(t2)
		suite.mockBuildingRepo.EXPECT().
			Get(gomock.Any(), buildingId).
			Times(1).
			Return(models.Building{}, domainerrors.ErrBuildingNotFound)

		_, err := suite.usecase.Create(t.Context(), req)

		assert.ErrorIs(t, err, domainerrors.ErrBuildingNotFound, "Actual err: %v", err)
	})
}

func setupCreateBuildingActionTestSuite(t *testing.T) *createBuildingActionTestSuite {
	t.Helper()

	ctrl := gomock.NewController(t)
	mockBuildingRepo := drivenportstest.NewMockForFetchingBuilding(ctrl)
	mockMutator := drivenportstest.NewMockForMutatingPlanet(ctrl)
	mockClock := drivenportstest.NewMockForFetchingTime(ctrl)

	return &createBuildingActionTestSuite{
		ctrl:             ctrl,
		mockBuildingRepo: mockBuildingRepo,
		mockMutator:      mockMutator,
		mockClock:        mockClock,
		usecase: NewCreateBuildingActionUseCase(
			mockBuildingRepo,
			mockMutator,
			mockClock,
		),
	}
}
