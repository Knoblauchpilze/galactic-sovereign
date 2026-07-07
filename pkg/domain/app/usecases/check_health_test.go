package usecases

import (
	"errors"
	"testing"

	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/usecases/drivenportstest"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestUnit_CheckHealth_Healthy(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockChecker := drivenportstest.NewMockForCheckingDatabaseConnection(ctrl)

	t.Run("returns healthy when database connection returns no error", func(t *testing.T) {
		mockChecker.EXPECT().
			Ping(gomock.Any()).
			Times(1).
			Return(nil)

		usecase := NewCheckHealthUseCase(mockChecker)
		actual := usecase.Healthy(t.Context())

		assert.True(t, actual)
	})

	t.Run("returns unhealthy when database connection returns an error", func(t *testing.T) {
		mockChecker.EXPECT().
			Ping(gomock.Any()).
			Times(1).
			Return(errors.New("stubbed error"))

		usecase := NewCheckHealthUseCase(mockChecker)
		actual := usecase.Healthy(t.Context())

		assert.False(t, actual)
	})
}
