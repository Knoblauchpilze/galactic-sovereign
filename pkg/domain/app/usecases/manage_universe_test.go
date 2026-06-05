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

func TestUnit_ManageUniverse_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRepo := drivenportstest.NewMockForManagingUniverses(ctrl)

	request := request.UniverseCreationRequest{
		Name: "the-best-universe",
	}

	t.Run("persists created universe", func(t *testing.T) {
		// https://pkg.go.dev/go.uber.org/mock/gomock#example-Call.DoAndReturn-CaptureArguments
		var captured models.Universe
		mockRepo.EXPECT().
			Create(gomock.Any(), gomock.AssignableToTypeOf(captured)).
			Times(1).
			DoAndReturn(func(ctx context.Context, universe models.Universe) error {
				captured = universe
				return nil
			})

		beforeInsertion := time.Now()

		usecase := NewUniverseUseCase(mockRepo)
		actual, err := usecase.Create(context.Background(), request)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, request.Name, captured.Name)
		assert.True(t, beforeInsertion.Before(captured.CreatedAt))
		assert.Equal(t, 0, captured.Version)
		assert.Equal(t, captured, actual)
	})

	t.Run("returns error when repository fails", func(t *testing.T) {
		// https://pkg.go.dev/go.uber.org/mock/gomock#example-Call.DoAndReturn-CaptureArguments
		expectedErr := errors.New("stubbed error")
		mockRepo.EXPECT().
			Create(gomock.Any(), gomock.Any()).
			Times(1).
			Return(expectedErr)

		usecase := NewUniverseUseCase(mockRepo)
		_, err := usecase.Create(context.Background(), request)

		assert.ErrorIs(t, expectedErr, err, "Actual err: %v", err)
	})
}

func TestUnit_ManageUniverse_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRepo := drivenportstest.NewMockForManagingUniverses(ctrl)

	t.Run("deletes existing universe", func(t *testing.T) {
		id := uuid.New()

		mockRepo.EXPECT().
			Delete(gomock.Any(), gomock.Eq(id)).
			Times(1).
			Return(nil)

		usecase := NewUniverseUseCase(mockRepo)
		err := usecase.Delete(context.Background(), id)
		require.NoError(t, err, "Actual err: %v", err)
	})

	t.Run("returns error when repository fails", func(t *testing.T) {
		expectedErr := errors.New("stubbed error")
		mockRepo.EXPECT().
			Delete(gomock.Any(), gomock.Any()).
			Times(1).
			Return(expectedErr)

		usecase := NewUniverseUseCase(mockRepo)
		err := usecase.Delete(context.Background(), uuid.New())

		assert.ErrorIs(t, expectedErr, err, "Actual err: %v", err)
	})
}
