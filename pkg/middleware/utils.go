package middleware

import (
	"net/http"

	"github.com/KnoblauchPilze/user-service/pkg/logger"
)

func pathFromRequest(req *http.Request) string {
	out := req.Host
	if req.URL.Path != "" {
		out += req.URL.Path
	}
	return out
}

func formatHttpStatusCode(status int) string {
	switch {
	case status >= 500:
		return logger.FormatWithColor(status, logger.Red)
	case status >= 400:
		return logger.FormatWithColor(status, logger.Yellow)
	case status >= 300:
		return logger.FormatWithColor(status, logger.Cyan)
	default:
		return logger.FormatWithColor(status, logger.Green)
	}
}
