package usecases

import (
	"errors"
	"testing"

	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/usecases/drivenportstest"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

type deletePlanetTestSuite struct {
	ctrl        *gomock.Controller
	mockDeleter *drivenportstest.MockForDeletingPlanet
	mockClock   *drivenportstest.MockForFetchingTime
	usecase     *DeletePlanetUseCase
}

func TestUnit_DeletePlanet_Delete(t *testing.T) {
	t.Run("deletes existing planet through deleter", func(t *testing.T) {
		suite := setupDeletePlanetTestSuite(t)
		id := uuid.New()

		suite.mockClock.EXPECT().Now(gomock.Any()).Times(1).Return(t2)
		suite.mockDeleter.EXPECT().
			Delete(gomock.Any(), id, gomock.Any()).
			Times(1).
			Return(nil)

		err := suite.usecase.Delete(t.Context(), id)
		require.NoError(t, err, "Actual err: %v", err)
	})

	t.Run("returns error when deleter fails", func(t *testing.T) {
		suite := setupDeletePlanetTestSuite(t)
		id := uuid.New()

		suite.mockClock.EXPECT().Now(gomock.Any()).Times(1).Return(t2)
		expected := errors.New("stubbed error")
		suite.mockDeleter.EXPECT().
			Delete(gomock.Any(), id, gomock.Any()).
			Times(1).
			Return(expected)

		err := suite.usecase.Delete(t.Context(), id)

		assert.ErrorIs(t, err, expected, "Actual err: %v", err)
	})
}

func setupDeletePlanetTestSuite(t *testing.T) *deletePlanetTestSuite {
	t.Helper()

	ctrl := gomock.NewController(t)
	mockDeleter := drivenportstest.NewMockForDeletingPlanet(ctrl)
	mockClock := drivenportstest.NewMockForFetchingTime(ctrl)

	return &deletePlanetTestSuite{
		ctrl:        ctrl,
		mockDeleter: mockDeleter,
		mockClock:   mockClock,
		usecase:     NewDeletePlanetUseCase(mockDeleter, mockClock),
	}
}
