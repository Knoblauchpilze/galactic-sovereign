package usecases

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	domainerrors "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models/errors"
	drivenports "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/ports/driven"
	drivingports "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/ports/driving"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/usecases/drivenportstest"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

var (
	t1 = time.Date(2026, time.July, 3, 15, 58, 31, 0, time.UTC)
	t2 = time.Date(2026, time.July, 3, 16, 58, 31, 0, time.UTC)
	t3 = time.Date(2026, time.July, 3, 17, 58, 31, 0, time.UTC)
	t4 = time.Date(2026, time.July, 3, 18, 58, 31, 0, time.UTC)

	someTime = time.Date(2026, time.June, 25, 22, 22, 49, 0, time.UTC)

	metalResourceId   = uuid.MustParse("b4419b6b-b3bf-4576-aa92-055283addbc8")
	crystalResourceId = uuid.MustParse("cd2ac9aa-9968-4ff5-b746-88f1f810fbb3")

	metalMineId = uuid.MustParse("d176e82d-f2ca-4611-996b-c4804096caef")
)

type MutatorMock func(context.Context, uuid.UUID, drivenports.PlanetMutator) (models.PlanetMutationResult, error)

type planetTestSuite struct {
	ctrl              *gomock.Controller
	mockPlanetRepo    *drivenportstest.MockForManagingPlanets
	mockPlanetMutator *drivenportstest.MockForMutatingPlanet
	mockClock         *drivenportstest.MockForFetchingTime
	usecase           drivingports.ForManagingPlanet
}

func TestUnit_ManagePlanet_Get(t *testing.T) {
	t.Run("gets existing planet through mutator", func(t *testing.T) {
		suite := setupPlanetTestSuite(t)
		expected := models.Planet{
			Id:   uuid.New(),
			Name: "my-planet",
		}

		suite.mockClock.EXPECT().Now(gomock.Any()).Times(1).Return(t2)
		suite.mockPlanetMutator.EXPECT().
			Mutate(gomock.Any(), gomock.Eq(expected.Id), gomock.Any()).
			Times(1).
			Return(generateMutationResult(expected), nil)

		actual, err := suite.usecase.Get(t.Context(), expected.Id)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, expected, actual)
	})

	t.Run("updates planet to current time", func(t *testing.T) {
		suite := setupPlanetTestSuite(t)
		planet := models.Planet{
			Id:        uuid.New(),
			Player:    uuid.New(),
			Name:      "my-planet",
			CreatedAt: t1,
			UpdatedAt: t1,
			Version:   2,
		}

		suite.mockClock.EXPECT().Now(gomock.Any()).Times(1).Return(t2)
		suite.mockPlanetMutator.EXPECT().
			Mutate(gomock.Any(), gomock.Eq(planet.Id), gomock.Any()).
			Times(1).
			DoAndReturn(generateApplyingMutatorMock(&planet))

		actual, err := suite.usecase.Get(t.Context(), planet.Id)
		require.NoError(t, err, "Actual err: %v", err)

		expected := models.Planet{
			Id:        planet.Id,
			Player:    planet.Player,
			Name:      planet.Name,
			CreatedAt: t1,
			UpdatedAt: t2,
			Version:   3,
		}
		assert.Equal(t, expected, actual)
	})

	t.Run("does not apply action when current time is before completion time", func(t *testing.T) {
		suite := setupPlanetTestSuite(t)
		planet := models.Planet{
			Id:        uuid.New(),
			Player:    uuid.New(),
			Name:      "my-planet",
			CreatedAt: t1,
			UpdatedAt: t1,
			Version:   2,
			BuildingAction: &models.BuildingAction{
				Id:          uuid.New(),
				CreatedAt:   t1,
				CompletedAt: t3,
			},
		}

		suite.mockClock.EXPECT().Now(gomock.Any()).Times(1).Return(t2)
		suite.mockPlanetMutator.EXPECT().
			Mutate(gomock.Any(), gomock.Eq(planet.Id), gomock.Any()).
			Times(1).
			DoAndReturn(generateApplyingMutatorMock(&planet))

		actual, err := suite.usecase.Get(t.Context(), planet.Id)
		require.NoError(t, err, "Actual err: %v", err)

		expected := models.Planet{
			Id:        planet.Id,
			Player:    planet.Player,
			Name:      planet.Name,
			CreatedAt: t1,
			UpdatedAt: t2,
			Version:   3,
			BuildingAction: &models.BuildingAction{
				Id:          planet.BuildingAction.Id,
				CreatedAt:   t1,
				CompletedAt: t3,
			},
		}
		assert.Equal(t, expected, actual)
	})

	t.Run("apply action when current time is after completion time", func(t *testing.T) {
		suite := setupPlanetTestSuite(t)
		planet := models.Planet{
			Id:        uuid.New(),
			Player:    uuid.New(),
			Name:      "my-planet",
			CreatedAt: t1,
			UpdatedAt: t1,
			Version:   2,
			Buildings: []models.PlanetBuilding{
				{Building: metalMineId, Level: 5},
			},
			BuildingAction: &models.BuildingAction{
				Id:           uuid.New(),
				Building:     metalMineId,
				DesiredLevel: 6,
				CreatedAt:    t1,
				CompletedAt:  t3,
			},
		}

		suite.mockClock.EXPECT().Now(gomock.Any()).Times(1).Return(t4)
		suite.mockPlanetMutator.EXPECT().
			Mutate(gomock.Any(), gomock.Eq(planet.Id), gomock.Any()).
			Times(1).
			DoAndReturn(generateApplyingMutatorMock(&planet))

		actual, err := suite.usecase.Get(t.Context(), planet.Id)
		require.NoError(t, err, "Actual err: %v", err)

		expected := models.Planet{
			Id:        planet.Id,
			Player:    planet.Player,
			Name:      planet.Name,
			CreatedAt: t1,
			UpdatedAt: t4,
			Version:   5,
			Buildings: []models.PlanetBuilding{
				{Building: metalMineId, Level: 6},
			},
			BuildingAction: nil,
		}
		assert.Equal(t, expected, actual)
	})

	t.Run("returns error when mutator fails", func(t *testing.T) {
		suite := setupPlanetTestSuite(t)

		suite.mockClock.EXPECT().Now(gomock.Any()).Times(1).Return(t2)
		expectedErr := errors.New("stubbed error")
		suite.mockPlanetMutator.EXPECT().
			Mutate(gomock.Any(), gomock.Any(), gomock.Any()).
			Times(1).
			Return(models.PlanetMutationResult{}, expectedErr)

		_, err := suite.usecase.Get(t.Context(), uuid.New())

		assert.ErrorIs(t, err, expectedErr, "Actual err: %v", err)
	})

	t.Run("returns error when planet is deleted during mutation", func(t *testing.T) {
		suite := setupPlanetTestSuite(t)
		planetId := uuid.New()

		suite.mockClock.EXPECT().Now(gomock.Any()).Times(1).Return(t2)
		suite.mockPlanetMutator.EXPECT().
			Mutate(gomock.Any(), gomock.Eq(planetId), gomock.Any()).
			Times(1).
			Return(models.PlanetMutationResult{Deleted: true}, nil)

		_, err := suite.usecase.Get(t.Context(), planetId)

		assert.ErrorIs(t, err, domainerrors.ErrNotFound, "Actual err: %v", err)
	})
}

func TestUnit_ManagePlanet_ListForPlayer(t *testing.T) {
	t.Run("lists existing planets through mutator", func(t *testing.T) {
		suite := setupPlanetTestSuite(t)
		player := uuid.New()
		p1 := models.Planet{Id: uuid.New(), Player: player, Name: "planet-1", CreatedAt: t1, UpdatedAt: t1}
		p2 := models.Planet{Id: uuid.New(), Player: player, Name: "planet-2", CreatedAt: t1, UpdatedAt: t1}

		suite.mockClock.EXPECT().Now(gomock.Any()).Times(1).Return(t2)
		suite.mockPlanetRepo.EXPECT().
			ListForPlayer(gomock.Any(), gomock.Eq(player)).
			Times(1).
			Return([]uuid.UUID{p1.Id, p2.Id}, nil)
		suite.mockPlanetMutator.EXPECT().
			Mutate(gomock.Any(), gomock.Eq(p1.Id), gomock.Any()).
			Times(1).
			Return(generateMutationResult(p1), nil)
		suite.mockPlanetMutator.EXPECT().
			Mutate(gomock.Any(), gomock.Eq(p2.Id), gomock.Any()).
			Times(1).
			Return(generateMutationResult(p2), nil)

		actual, err := suite.usecase.ListForPlayer(t.Context(), player)
		require.NoError(t, err, "Actual err: %v", err)

		expected := []models.Planet{p1, p2}
		assert.Equal(t, expected, actual)
	})

	t.Run("updates all planet to same time", func(t *testing.T) {
		suite := setupPlanetTestSuite(t)
		player := uuid.New()
		p1 := models.Planet{
			Id:        uuid.New(),
			Player:    player,
			Name:      "planet-1",
			CreatedAt: t1,
			UpdatedAt: t1,
			Version:   2,
		}
		p2 := models.Planet{
			Id:        uuid.New(),
			Player:    player,
			Name:      "planet-2",
			CreatedAt: t1,
			UpdatedAt: t1,
			Version:   3,
		}

		suite.mockClock.EXPECT().Now(gomock.Any()).Times(1).Return(t2)
		suite.mockPlanetRepo.EXPECT().
			ListForPlayer(gomock.Any(), gomock.Eq(player)).
			Times(1).
			Return([]uuid.UUID{p1.Id, p2.Id}, nil)
		suite.mockPlanetMutator.EXPECT().
			Mutate(gomock.Any(), gomock.Eq(p1.Id), gomock.Any()).
			Times(1).
			DoAndReturn(generateApplyingMutatorMock(&p1))
		suite.mockPlanetMutator.EXPECT().
			Mutate(gomock.Any(), gomock.Eq(p2.Id), gomock.Any()).
			Times(1).
			DoAndReturn(generateApplyingMutatorMock(&p2))

		actual, err := suite.usecase.ListForPlayer(t.Context(), player)
		require.NoError(t, err, "Actual err: %v", err)

		expected := []models.Planet{
			{
				Id:        p1.Id,
				Player:    player,
				Name:      "planet-1",
				Version:   3,
				CreatedAt: t1,
				UpdatedAt: t2,
			},
			{
				Id:        p2.Id,
				Player:    player,
				Name:      "planet-2",
				Version:   4,
				CreatedAt: t1,
				UpdatedAt: t2,
			},
		}
		assert.Equal(t, expected, actual)
	})

	t.Run("does not apply action when current time is before completion time", func(t *testing.T) {
		suite := setupPlanetTestSuite(t)
		player := uuid.New()
		p1 := models.Planet{
			Id:        uuid.New(),
			Player:    player,
			Name:      "planet-1",
			CreatedAt: t1,
			UpdatedAt: t1,
			Version:   2,
			Buildings: []models.PlanetBuilding{
				{Building: metalMineId, Level: 5},
			},
			BuildingAction: &models.BuildingAction{
				Id:          uuid.New(),
				CreatedAt:   t1,
				CompletedAt: t3,
			},
		}
		p2 := models.Planet{
			Id:        uuid.New(),
			Player:    player,
			Name:      "planet-2",
			CreatedAt: t1,
			UpdatedAt: t1,
			Version:   3,
		}

		suite.mockClock.EXPECT().Now(gomock.Any()).Times(1).Return(t2)
		suite.mockPlanetRepo.EXPECT().
			ListForPlayer(gomock.Any(), gomock.Eq(player)).
			Times(1).
			Return([]uuid.UUID{p1.Id, p2.Id}, nil)
		suite.mockPlanetMutator.EXPECT().
			Mutate(gomock.Any(), gomock.Eq(p1.Id), gomock.Any()).
			Times(1).
			DoAndReturn(generateApplyingMutatorMock(&p1))
		suite.mockPlanetMutator.EXPECT().
			Mutate(gomock.Any(), gomock.Eq(p2.Id), gomock.Any()).
			Times(1).
			DoAndReturn(generateApplyingMutatorMock(&p2))

		actual, err := suite.usecase.ListForPlayer(t.Context(), player)
		require.NoError(t, err, "Actual err: %v", err)

		expected := []models.Planet{
			{
				Id:        p1.Id,
				Player:    player,
				Name:      "planet-1",
				Version:   3,
				CreatedAt: t1,
				UpdatedAt: t2,
				Buildings: []models.PlanetBuilding{
					{Building: metalMineId, Level: 5},
				},
				BuildingAction: &models.BuildingAction{
					Id:          p1.BuildingAction.Id,
					CreatedAt:   t1,
					CompletedAt: t3,
				},
			},
			{
				Id:        p2.Id,
				Player:    player,
				Name:      "planet-2",
				Version:   4,
				CreatedAt: t1,
				UpdatedAt: t2,
			},
		}
		assert.Equal(t, expected, actual)
	})

	t.Run("apply action when current time is after completion time", func(t *testing.T) {
		suite := setupPlanetTestSuite(t)
		player := uuid.New()
		p1 := models.Planet{
			Id:        uuid.New(),
			Player:    player,
			Name:      "planet-1",
			CreatedAt: t1,
			UpdatedAt: t1,
			Version:   2,
			Buildings: []models.PlanetBuilding{
				{Building: metalMineId, Level: 5},
			},
			BuildingAction: &models.BuildingAction{
				Id:           uuid.New(),
				Building:     metalMineId,
				DesiredLevel: 6,
				CreatedAt:    t1,
				CompletedAt:  t3,
			},
		}
		p2 := models.Planet{
			Id:        uuid.New(),
			Player:    player,
			Name:      "planet-2",
			CreatedAt: t1,
			UpdatedAt: t1,
			Version:   3,
		}

		suite.mockClock.EXPECT().Now(gomock.Any()).Times(1).Return(t4)
		suite.mockPlanetRepo.EXPECT().
			ListForPlayer(gomock.Any(), gomock.Eq(player)).
			Times(1).
			Return([]uuid.UUID{p1.Id, p2.Id}, nil)
		suite.mockPlanetMutator.EXPECT().
			Mutate(gomock.Any(), gomock.Eq(p1.Id), gomock.Any()).
			Times(1).
			DoAndReturn(generateApplyingMutatorMock(&p1))
		suite.mockPlanetMutator.EXPECT().
			Mutate(gomock.Any(), gomock.Eq(p2.Id), gomock.Any()).
			Times(1).
			DoAndReturn(generateApplyingMutatorMock(&p2))

		actual, err := suite.usecase.ListForPlayer(t.Context(), player)
		require.NoError(t, err, "Actual err: %v", err)

		expected := []models.Planet{
			{
				Id:        p1.Id,
				Player:    player,
				Name:      "planet-1",
				Version:   5,
				CreatedAt: t1,
				UpdatedAt: t4,
				Buildings: []models.PlanetBuilding{
					{Building: metalMineId, Level: 6},
				},
				BuildingAction: nil,
			},
			{
				Id:        p2.Id,
				Player:    player,
				Name:      "planet-2",
				Version:   4,
				CreatedAt: t1,
				UpdatedAt: t4,
			},
		}
		assert.Equal(t, expected, actual)
	})

	t.Run("returns error when repository fails", func(t *testing.T) {
		suite := setupPlanetTestSuite(t)
		expectedErr := errors.New("stubbed error")

		suite.mockClock.EXPECT().Now(gomock.Any()).Times(1).Return(t4)
		suite.mockPlanetRepo.EXPECT().
			ListForPlayer(gomock.Any(), gomock.Any()).
			Times(1).
			Return([]uuid.UUID{uuid.New(), uuid.New()}, nil)
		suite.mockPlanetMutator.EXPECT().
			Mutate(gomock.Any(), gomock.Any(), gomock.Any()).
			Times(1).
			Return(models.PlanetMutationResult{}, expectedErr)

		_, err := suite.usecase.ListForPlayer(t.Context(), uuid.New())

		assert.ErrorIs(t, err, expectedErr, "Actual err: %v", err)
	})

	t.Run("does not return planet when it is deleted during mutation", func(t *testing.T) {
		suite := setupPlanetTestSuite(t)
		player := uuid.New()
		p1 := models.Planet{
			Id:        uuid.New(),
			Player:    player,
			Name:      "planet-1",
			CreatedAt: t1,
			UpdatedAt: t1,
			Version:   2,
		}
		p2 := uuid.New()

		suite.mockClock.EXPECT().Now(gomock.Any()).Times(1).Return(t2)
		suite.mockPlanetRepo.EXPECT().
			ListForPlayer(gomock.Any(), gomock.Eq(player)).
			Times(1).
			Return([]uuid.UUID{p1.Id, p2}, nil)
		suite.mockPlanetMutator.EXPECT().
			Mutate(gomock.Any(), gomock.Eq(p1.Id), gomock.Any()).
			Times(1).
			DoAndReturn(generateApplyingMutatorMock(&p1))
		suite.mockPlanetMutator.EXPECT().
			Mutate(gomock.Any(), gomock.Eq(p2), gomock.Any()).
			Times(1).
			Return(models.PlanetMutationResult{Deleted: true}, nil)

		actual, err := suite.usecase.ListForPlayer(t.Context(), player)
		require.NoError(t, err, "Actual err: %v", err)

		expected := []models.Planet{
			{
				Id:        p1.Id,
				Player:    player,
				Name:      "planet-1",
				Version:   3,
				CreatedAt: t1,
				UpdatedAt: t2,
			},
		}
		assert.Equal(t, expected, actual)
	})
}

func TestUnit_ManagePlanet_Delete(t *testing.T) {
	t.Run("deletes existing planet through mutator", func(t *testing.T) {
		suite := setupPlanetTestSuite(t)
		id := uuid.New()

		suite.mockClock.EXPECT().Now(gomock.Any()).Times(1).Return(t2)
		suite.mockPlanetMutator.EXPECT().
			Mutate(gomock.Any(), gomock.Eq(id), gomock.Any()).
			Times(1).
			Return(models.PlanetMutationResult{Deleted: true}, nil)

		err := suite.usecase.Delete(t.Context(), id)
		require.NoError(t, err, "Actual err: %v", err)
	})

	t.Run("returns error when planet has a building action", func(t *testing.T) {
		suite := setupPlanetTestSuite(t)
		id := uuid.New()

		suite.mockClock.EXPECT().Now(gomock.Any()).Times(1).Return(t2)
		suite.mockPlanetMutator.EXPECT().
			Mutate(gomock.Any(), gomock.Eq(id), gomock.Any()).
			Times(1).
			Return(models.PlanetMutationResult{}, domainerrors.ErrActionNotCompleted)

		err := suite.usecase.Delete(t.Context(), id)

		assert.ErrorIs(t, err, domainerrors.ErrActionNotCompleted, "Actual err: %v", err)
	})

	t.Run("returns error when mutator returns no error but does not mark the planet as deleted", func(t *testing.T) {
		suite := setupPlanetTestSuite(t)
		id := uuid.New()

		suite.mockClock.EXPECT().Now(gomock.Any()).Times(1).Return(t2)
		suite.mockPlanetMutator.EXPECT().
			Mutate(gomock.Any(), gomock.Eq(id), gomock.Any()).
			Times(1).
			Return(models.PlanetMutationResult{Deleted: false}, nil)

		err := suite.usecase.Delete(t.Context(), id)

		assert.ErrorIs(t, err, domainerrors.ErrPlanetDeletionFailed, "Actual err: %v", err)
	})
}

func setupPlanetTestSuite(t *testing.T) *planetTestSuite {
	t.Helper()

	ctrl := gomock.NewController(t)
	mockPlanetRepo := drivenportstest.NewMockForManagingPlanets(ctrl)
	mockPlanetMutator := drivenportstest.NewMockForMutatingPlanet(ctrl)
	mockClock := drivenportstest.NewMockForFetchingTime(ctrl)

	return &planetTestSuite{
		ctrl:              ctrl,
		mockPlanetRepo:    mockPlanetRepo,
		mockPlanetMutator: mockPlanetMutator,
		mockClock:         mockClock,
		usecase:           NewPlanetUseCase(mockPlanetRepo, mockPlanetMutator, mockClock),
	}
}

// generateApplyingMutatorMock generates a function mock for the planet mutator
// which applies the provided mutator to a known planet.
func generateApplyingMutatorMock(p *models.Planet) MutatorMock {
	return func(
		ctx context.Context, id uuid.UUID, m drivenports.PlanetMutator,
	) (models.PlanetMutationResult, error) {
		deleted, err := m(p)
		result := models.PlanetMutationResult{
			Deleted: deleted,
			Planet:  *p,
		}

		return result, err
	}
}

func generateMutationResult(planet models.Planet) models.PlanetMutationResult {
	return models.PlanetMutationResult{
		Deleted: false,
		Planet:  planet,
	}
}
