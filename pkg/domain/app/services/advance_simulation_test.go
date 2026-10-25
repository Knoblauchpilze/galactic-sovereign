package domainservices

import (
	"errors"
	"testing"

	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/services/drivenportstest"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

type simulationStepTestSuite struct {
	ctrl         *gomock.Controller
	mockExecutor *drivenportstest.MockForExecutingSimulationStep
}

func TestUnit_AdvanceSimulationUntil(t *testing.T) {
	t.Run("calls executor with correct parameter", func(t *testing.T) {
		suite := setupSimulationStepTestSuite(t)

		planet := generateSamplePlanet()

		expected := models.SimulationStepResult{
			Resolution: models.PlanetUpdated,
			Planet:     &planet,
		}
		suite.mockExecutor.EXPECT().
			Execute(gomock.Any(), planet.Id, t2, gomock.Any()).
			Times(1).
			Return(expected, nil)

		resolver := generateDummyResolver()
		result, err := AdvanceSimulationUntil(t.Context(), suite.mockExecutor, planet.Id, t2, resolver)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, expected, result)
	})

	t.Run("calls resolver again when it returns FleetEventResolved", func(t *testing.T) {
		suite := setupSimulationStepTestSuite(t)

		planet := generateSamplePlanet()
		fleet := generateSampleFleet(planet)

		result1 := models.SimulationStepResult{
			Resolution: models.FleetEventResolved,
			Fleet:      &fleet,
		}
		suite.mockExecutor.EXPECT().
			Execute(gomock.Any(), planet.Id, t2, gomock.Any()).
			Times(1).
			Return(result1, nil)
		result2 := models.SimulationStepResult{
			Resolution: models.PlanetUpdated,
			Planet:     &planet,
		}
		suite.mockExecutor.EXPECT().
			Execute(gomock.Any(), planet.Id, t2, gomock.Any()).
			Times(1).
			Return(result2, nil)

		resolver := generateDummyResolver()
		result, err := AdvanceSimulationUntil(t.Context(), suite.mockExecutor, planet.Id, t2, resolver)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, result2, result)
	})

	t.Run("calls resolver again when it returns FleetDeleted", func(t *testing.T) {
		suite := setupSimulationStepTestSuite(t)

		planet := generateSamplePlanet()

		result1 := models.SimulationStepResult{
			Resolution: models.FleetDeleted,
		}
		suite.mockExecutor.EXPECT().
			Execute(gomock.Any(), planet.Id, t2, gomock.Any()).
			Times(1).
			Return(result1, nil)
		result2 := models.SimulationStepResult{
			Resolution: models.PlanetUpdated,
			Planet:     &planet,
		}
		suite.mockExecutor.EXPECT().
			Execute(gomock.Any(), planet.Id, t2, gomock.Any()).
			Times(1).
			Return(result2, nil)

		resolver := generateDummyResolver()
		result, err := AdvanceSimulationUntil(t.Context(), suite.mockExecutor, planet.Id, t2, resolver)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, result2, result)
	})

	t.Run("stops calling resolver when it returns FleetCreated", func(t *testing.T) {
		suite := setupSimulationStepTestSuite(t)

		planet := generateSamplePlanet()
		fleet := generateSampleFleet(planet)

		result := models.SimulationStepResult{
			Resolution: models.FleetCreated,
			Planet:     &planet,
			Fleet:      &fleet,
		}
		suite.mockExecutor.EXPECT().
			Execute(gomock.Any(), planet.Id, t2, gomock.Any()).
			Times(1).
			Return(result, nil)

		resolver := generateDummyResolver()
		result, err := AdvanceSimulationUntil(t.Context(), suite.mockExecutor, planet.Id, t2, resolver)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, result, result)
	})

	t.Run("stops calling resolver when it returns PlanetUpdated", func(t *testing.T) {
		suite := setupSimulationStepTestSuite(t)

		planet := generateSamplePlanet()

		result := models.SimulationStepResult{
			Resolution: models.PlanetUpdated,
			Planet:     &planet,
		}
		suite.mockExecutor.EXPECT().
			Execute(gomock.Any(), planet.Id, t2, gomock.Any()).
			Times(1).
			Return(result, nil)

		resolver := generateDummyResolver()
		result, err := AdvanceSimulationUntil(t.Context(), suite.mockExecutor, planet.Id, t2, resolver)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, result, result)
	})

	t.Run("returns error when executor returns an error", func(t *testing.T) {
		suite := setupSimulationStepTestSuite(t)

		planet := generateSamplePlanet()

		errSample := errors.New("stubbed error")
		suite.mockExecutor.EXPECT().
			Execute(gomock.Any(), planet.Id, t2, gomock.Any()).
			Times(1).
			Return(models.SimulationStepResult{}, errSample)

		resolver := generateDummyResolver()
		_, err := AdvanceSimulationUntil(t.Context(), suite.mockExecutor, planet.Id, t2, resolver)

		assert.ErrorIs(t, err, errSample, "Actual err: %v", err)
	})

	t.Run("returns error when subsequent calls to executor return an error", func(t *testing.T) {
		suite := setupSimulationStepTestSuite(t)

		planet := generateSamplePlanet()
		fleet := generateSampleFleet(planet)

		result := models.SimulationStepResult{
			Resolution: models.FleetEventResolved,
			Fleet:      &fleet,
		}
		suite.mockExecutor.EXPECT().
			Execute(gomock.Any(), planet.Id, t2, gomock.Any()).
			Times(1).
			Return(result, nil)
		errSample := errors.New("stubbed error")
		suite.mockExecutor.EXPECT().
			Execute(gomock.Any(), planet.Id, t2, gomock.Any()).
			Times(1).
			Return(models.SimulationStepResult{}, errSample)

		resolver := generateDummyResolver()
		_, err := AdvanceSimulationUntil(t.Context(), suite.mockExecutor, planet.Id, t2, resolver)

		assert.ErrorIs(t, err, errSample, "Actual err: %v", err)
	})
}

func setupSimulationStepTestSuite(t *testing.T) *simulationStepTestSuite {
	t.Helper()

	ctrl := gomock.NewController(t)
	mockExecutor := drivenportstest.NewMockForExecutingSimulationStep(ctrl)

	return &simulationStepTestSuite{
		ctrl:         ctrl,
		mockExecutor: mockExecutor,
	}
}

func generateSamplePlanet() models.Planet {
	return models.Planet{
		Id:        uuid.New(),
		Player:    uuid.New(),
		CreatedAt: t1,
		UpdatedAt: t1,
		Version:   0,
	}
}

func generateSampleFleet(source models.Planet) models.Fleet {
	return models.Fleet{
		Id:      uuid.New(),
		Player:  source.Player,
		Source:  source.Id,
		Mission: models.MissionColonize,
		Schedule: models.FleetFlightSchedule{
			ArrivalAt: t2,
			ReturnAt:  t4,
			Status:    models.StatusOutbound,
		},
		CreatedAt: t1,
		UpdatedAt: t1,
		Version:   0,
	}
}

func generateDummyResolver() models.SimulationStepResolver {
	return func(_ models.SimulationStep) (models.SimulationStepResult, error) {
		return models.SimulationStepResult{}, nil
	}
}
