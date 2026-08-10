package drivingadapters

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/adapters/driving/drivingportstest"
	drivingports "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/ports/driving"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestUnit_Healthcheck_Healthcheck(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ctrl := gomock.NewController(t)
	mockUsecase := drivingportstest.NewMockForCheckingServiceHealth(ctrl)

	t.Run("returns 200 when service is healthy", func(t *testing.T) {
		mockUsecase.EXPECT().
			Healthy(gomock.Any()).
			Times(1).
			Return(true)

		handler := generateHandler[drivingports.ForCheckingServiceHealth](
			healthcheck,
			mockUsecase,
		)
		r := createTestGinRouterWithHandler(t, handler)

		req := httptest.NewRequest(http.MethodGet, "http://example.com/", nil)
		rw := httptest.NewRecorder()
		r.ServeHTTP(rw, req)

		assert.Equal(t, http.StatusOK, rw.Code)
		actual := decodeResponseBody[string](t, rw)
		assert.Equal(t, "OK", actual)
	})

	t.Run("returns 503 when service is not healthy", func(t *testing.T) {
		mockUsecase.EXPECT().
			Healthy(gomock.Any()).
			Times(1).
			Return(false)

		handler := generateHandler[drivingports.ForCheckingServiceHealth](
			healthcheck,
			mockUsecase,
		)
		r := createTestGinRouterWithHandler(t, handler)

		req := httptest.NewRequest(http.MethodGet, "http://example.com/", nil)
		rw := httptest.NewRecorder()
		r.ServeHTTP(rw, req)

		assert.Equal(t, http.StatusServiceUnavailable, rw.Code)
		actual := decodeResponseBody[string](t, rw)
		assert.Equal(t, "KO", actual)
	})
}
