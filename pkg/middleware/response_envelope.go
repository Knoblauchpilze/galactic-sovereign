package middleware

import (
	"encoding/json"
	"net/http"

	"github.com/KnoblauchPilze/user-service/pkg/logger"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type responseEnvelope struct {
	RequestId uuid.UUID
	Status    string
	Details   json.RawMessage `json:",omitempty"`
}

type envelopeResponseWriter struct {
	response responseEnvelope
	writer   http.ResponseWriter

	logger echo.Logger
}

const RequestIdKey = "requestIdKey"

func ResponseEnvelope() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			id := uuid.Must(uuid.NewRandom())
			log := logger.New(id.String())

			c.Set(RequestIdKey, id)
			c.SetLogger(log)

			w := &envelopeResponseWriter{
				response: responseEnvelope{
					RequestId: id,
				},
				writer: c.Response().Writer,
				logger: log,
			}
			c.Response().Writer = w

			return next(c)
		}
	}
}

func (erw *envelopeResponseWriter) Header() http.Header {
	return erw.writer.Header()
}

func (erw *envelopeResponseWriter) Write(data []byte) (int, error) {
	erw.response.Details = data
	out, err := json.Marshal(erw.response)
	if err != nil {
		erw.logger.Warnf("Failed to write data %s (err: %v), no response envelope", string(data), err)
		return erw.writer.Write(data)
	}

	return erw.writer.Write(out)
}

func (erw *envelopeResponseWriter) WriteHeader(statusCode int) {
	if statusCode < 200 || statusCode > 299 {
		erw.response.Status = "ERROR"
	} else {
		erw.response.Status = "SUCCESS"
	}
	erw.writer.WriteHeader(statusCode)
}
