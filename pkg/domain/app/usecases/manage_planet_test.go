package usecases

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
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

	metalMineId = uuid.MustParse("d176e82d-f2ca-4611-996b-c4804096caef")
)

type MutatorMock func(context.Context, uuid.UUID, drivenports.PlanetMutator) (models.Planet, error)

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
			Return(expected, nil)

		actual, err := suite.usecase.Get(context.Background(), expected.Id)
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

		actual, err := suite.usecase.Get(context.Background(), planet.Id)
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

		actual, err := suite.usecase.Get(context.Background(), planet.Id)
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

		actual, err := suite.usecase.Get(context.Background(), planet.Id)
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
			Return(models.Planet{}, expectedErr)

		_, err := suite.usecase.Get(context.Background(), uuid.New())

		assert.ErrorIs(t, expectedErr, err, "Actual err: %v", err)
	})
}

func TestUnit_ManagePlanet_List(t *testing.T) {
	t.Run("lists existing planets", func(t *testing.T) {
		suite := setupPlanetTestSuite(t)
		expected := []models.Planet{
			{
				Id:   uuid.New(),
				Name: "planet-1",
			},
			{
				Id:   uuid.New(),
				Name: "planet-2",
			},
		}

		suite.mockPlanetRepo.EXPECT().
			List(gomock.Any()).
			Times(1).
			Return(expected, nil)

		actual, err := suite.usecase.List(context.Background())
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, expected, actual)
	})

	t.Run("returns error when repository fails", func(t *testing.T) {
		suite := setupPlanetTestSuite(t)
		expectedErr := errors.New("stubbed error")

		suite.mockPlanetRepo.EXPECT().
			List(gomock.Any()).
			Times(1).
			Return(nil, expectedErr)

		_, err := suite.usecase.List(context.Background())

		assert.ErrorIs(t, expectedErr, err, "Actual err: %v", err)
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
			Return(p1, nil)
		suite.mockPlanetMutator.EXPECT().
			Mutate(gomock.Any(), gomock.Eq(p2.Id), gomock.Any()).
			Times(1).
			Return(p2, nil)

		actual, err := suite.usecase.ListForPlayer(context.Background(), player)
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

		actual, err := suite.usecase.ListForPlayer(context.Background(), player)
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

		actual, err := suite.usecase.ListForPlayer(context.Background(), player)
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

		actual, err := suite.usecase.ListForPlayer(context.Background(), player)
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
			Return(models.Planet{}, expectedErr)

		_, err := suite.usecase.ListForPlayer(context.Background(), uuid.New())

		assert.ErrorIs(t, expectedErr, err, "Actual err: %v", err)
	})
}

func TestUnit_ManagePlanet_Delete(t *testing.T) {
	t.Run("deletes existing planet", func(t *testing.T) {
		suite := setupPlanetTestSuite(t)
		id := uuid.New()

		suite.mockPlanetRepo.EXPECT().
			Delete(gomock.Any(), gomock.Eq(id)).
			Times(1).
			Return(nil)

		err := suite.usecase.Delete(context.Background(), id)
		require.NoError(t, err, "Actual err: %v", err)
	})

	t.Run("returns error when repository fails", func(t *testing.T) {
		suite := setupPlanetTestSuite(t)
		expectedErr := errors.New("stubbed error")
		suite.mockPlanetRepo.EXPECT().
			Delete(gomock.Any(), gomock.Any()).
			Times(1).
			Return(expectedErr)

		err := suite.usecase.Delete(context.Background(), uuid.New())

		assert.ErrorIs(t, expectedErr, err, "Actual err: %v", err)
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
	return func(ctx context.Context, id uuid.UUID, m drivenports.PlanetMutator) (models.Planet, error) {
		err := m(p)
		return *p, err
	}
}
