package usecases

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models/request"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/usecases/drivenportstest"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

type universeTestSuite struct {
	ctrl          *gomock.Controller
	mockFetchRepo *drivenportstest.MockForFetchingUniverses
	mockRepo      *drivenportstest.MockForManagingUniverses
	usecase       *UniverseUseCase
}

func TestUnit_ManageUniverse_Create(t *testing.T) {
	request := request.UniverseCreationRequest{
		Name: "the-best-universe",
	}

	t.Run("persists created universe", func(t *testing.T) {
		suite := setupUniverseTestSuite(t)

		// https://pkg.go.dev/go.uber.org/mock/gomock#example-Call.DoAndReturn-CaptureArguments
		var captured models.Universe
		suite.mockRepo.EXPECT().
			Create(gomock.Any(), gomock.AssignableToTypeOf(captured)).
			Times(1).
			DoAndReturn(func(ctx context.Context, universe models.Universe) error {
				captured = universe
				return nil
			})

		beforeInsertion := time.Now()

		actual, err := suite.usecase.Create(t.Context(), request)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, request.Name, captured.Name)
		assert.True(t, beforeInsertion.Before(captured.CreatedAt))
		assert.Equal(t, 0, captured.Version)
		assert.Equal(t, captured, actual)
	})

	t.Run("returns error when repository fails", func(t *testing.T) {
		suite := setupUniverseTestSuite(t)

		expectedErr := errors.New("stubbed error")
		suite.mockRepo.EXPECT().
			Create(gomock.Any(), gomock.Any()).
			Times(1).
			Return(expectedErr)

		_, err := suite.usecase.Create(t.Context(), request)

		assert.ErrorIs(t, err, expectedErr, "Actual err: %v", err)
	})
}

func TestUnit_ManageUniverse_Get(t *testing.T) {
	t.Run("gets existing universe", func(t *testing.T) {
		suite := setupUniverseTestSuite(t)

		expected := models.Universe{
			Id:   uuid.New(),
			Name: "my-universe",
		}

		suite.mockFetchRepo.EXPECT().
			Get(gomock.Any(), gomock.Eq(expected.Id)).
			Times(1).
			Return(expected, nil)

		actual, err := suite.usecase.Get(t.Context(), expected.Id)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, expected, actual)
	})

	t.Run("returns error when repository fails", func(t *testing.T) {
		suite := setupUniverseTestSuite(t)

		expectedErr := errors.New("stubbed error")
		suite.mockFetchRepo.EXPECT().
			Get(gomock.Any(), gomock.Any()).
			Times(1).
			Return(models.Universe{}, expectedErr)

		_, err := suite.usecase.Get(t.Context(), uuid.New())

		assert.ErrorIs(t, err, expectedErr, "Actual err: %v", err)
	})
}

func TestUnit_ManageUniverse_List(t *testing.T) {
	t.Run("lists existing universes", func(t *testing.T) {
		suite := setupUniverseTestSuite(t)

		expected := []models.Universe{
			{
				Id:   uuid.New(),
				Name: "universe-1",
			},
			{
				Id:   uuid.New(),
				Name: "universe-1",
			},
		}

		suite.mockFetchRepo.EXPECT().
			List(gomock.Any()).
			Times(1).
			Return(expected, nil)

		actual, err := suite.usecase.List(t.Context())
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, expected, actual)
	})

	t.Run("returns error when repository fails", func(t *testing.T) {
		suite := setupUniverseTestSuite(t)

		expectedErr := errors.New("stubbed error")

		suite.mockFetchRepo.EXPECT().
			List(gomock.Any()).
			Times(1).
			Return(nil, expectedErr)

		_, err := suite.usecase.List(t.Context())

		assert.ErrorIs(t, err, expectedErr, "Actual err: %v", err)
	})
}

func TestUnit_ManageUniverse_Delete(t *testing.T) {
	t.Run("deletes existing universe", func(t *testing.T) {
		suite := setupUniverseTestSuite(t)

		id := uuid.New()

		suite.mockRepo.EXPECT().
			Delete(gomock.Any(), gomock.Eq(id)).
			Times(1).
			Return(nil)

		err := suite.usecase.Delete(t.Context(), id)
		require.NoError(t, err, "Actual err: %v", err)
	})

	t.Run("returns error when repository fails", func(t *testing.T) {
		suite := setupUniverseTestSuite(t)

		expectedErr := errors.New("stubbed error")
		suite.mockRepo.EXPECT().
			Delete(gomock.Any(), gomock.Any()).
			Times(1).
			Return(expectedErr)

		err := suite.usecase.Delete(t.Context(), uuid.New())

		assert.ErrorIs(t, err, expectedErr, "Actual err: %v", err)
	})
}

func setupUniverseTestSuite(t *testing.T) *universeTestSuite {
	t.Helper()

	ctrl := gomock.NewController(t)
	mockFetchRepo := drivenportstest.NewMockForFetchingUniverses(ctrl)
	mockRepo := drivenportstest.NewMockForManagingUniverses(ctrl)

	return &universeTestSuite{
		ctrl:          ctrl,
		mockFetchRepo: mockFetchRepo,
		mockRepo:      mockRepo,
		usecase:       NewUniverseUseCase(mockFetchRepo, mockRepo),
	}
}
