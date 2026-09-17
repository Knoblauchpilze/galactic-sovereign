package usecases

import (
	"testing"

	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	domainerrors "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models/errors"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/usecases/drivenportstest"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

type solarSystemTestSuite struct {
	ctrl                *gomock.Controller
	mockSolarSystemRepo *drivenportstest.MockForFetchingSolarSystem
	usecase             *FetchSolarSystemUseCase
}

func TestUnit_FetchSolarSystem_GetSolarSystem(t *testing.T) {
	t.Run("gets solar system in existing universe", func(t *testing.T) {
		suite := setupSolarSystemTestSuite(t)

		universeId := uuid.New()

		expected := models.SolarSystem{
			Universe: universeId,
			Galaxy:   1,
			Number:   3,
			Orbits:   5,
			Planets:  []models.SolarSystemPlanet{},
		}

		suite.mockSolarSystemRepo.EXPECT().
			GetSolarSystem(gomock.Any(), universeId, 1, 3).
			Times(1).
			Return(expected, nil)

		actual, err := suite.usecase.GetSolarSystem(t.Context(), universeId, 1, 3)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, expected, actual)
	})

	t.Run("returns error when universe does not exist", func(t *testing.T) {
		suite := setupSolarSystemTestSuite(t)

		suite.mockSolarSystemRepo.EXPECT().
			GetSolarSystem(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Times(1).
			Return(models.SolarSystem{}, domainerrors.ErrNotFound)

		_, err := suite.usecase.GetSolarSystem(t.Context(), uuid.New(), 2, 4)

		assert.ErrorIs(t, err, domainerrors.ErrNotFound, "Actual err: %v", err)
	})

	t.Run("returns error when requested coordinates are out of bound", func(t *testing.T) {
		suite := setupSolarSystemTestSuite(t)

		suite.mockSolarSystemRepo.EXPECT().
			GetSolarSystem(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Times(1).
			// This is the same error as when the universe does not exist
			Return(models.SolarSystem{}, domainerrors.ErrNotFound)

		_, err := suite.usecase.GetSolarSystem(t.Context(), uuid.New(), 2, 4)

		assert.ErrorIs(t, err, domainerrors.ErrNotFound, "Actual err: %v", err)
	})
}

func setupSolarSystemTestSuite(t *testing.T) *solarSystemTestSuite {
	t.Helper()

	ctrl := gomock.NewController(t)
	mockSolarSystemRepo := drivenportstest.NewMockForFetchingSolarSystem(ctrl)

	return &solarSystemTestSuite{
		ctrl:                ctrl,
		mockSolarSystemRepo: mockSolarSystemRepo,
		usecase:             NewFetchSolarSystemUseCase(mockSolarSystemRepo),
	}
}
