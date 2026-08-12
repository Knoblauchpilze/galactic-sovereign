package drivingadapters

import (
	"net/http"

	"github.com/Knoblauchpilze/backend-toolkit/pkg/rest"
)

func logError(req *http.Request, msg string, args ...any) {
	log := rest.GetContextLogger(req.Context())
	log.Error(msg, args...)
}
