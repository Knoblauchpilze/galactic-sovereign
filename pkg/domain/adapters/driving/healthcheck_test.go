package drivingadapters

import (
	"net/http"
	"testing"

	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/adapters/driving/drivingportstest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestUnit_Healthcheck_Healthcheck(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockUsecase := drivingportstest.NewMockForCheckingServiceHealth(ctrl)

	t.Run("returns 200 when service is healthy", func(t *testing.T) {
		req := generateTestRequest(t, http.MethodGet)
		ctx, rw := generateTestContextFromRequest(t, req)

		mockUsecase.EXPECT().
			Healthy(gomock.Any()).
			Times(1).
			Return(true)

		err := healthcheck(ctx, mockUsecase)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, http.StatusOK, rw.Code)
		actual := decodeResponseBody[string](t, rw)
		assert.Equal(t, "OK", actual)
	})

	t.Run("returns 503 when service is not healthy", func(t *testing.T) {
		req := generateTestRequest(t, http.MethodGet)
		ctx, rw := generateTestContextFromRequest(t, req)

		mockUsecase.EXPECT().
			Healthy(gomock.Any()).
			Times(1).
			Return(false)

		err := healthcheck(ctx, mockUsecase)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, http.StatusServiceUnavailable, rw.Code)
		actual := decodeResponseBody[string](t, rw)
		assert.Equal(t, "KO", actual)
	})
}
