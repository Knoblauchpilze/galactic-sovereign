package rest

import (
	"fmt"
	"strings"
)

func sanitizePath(route string) string {
	route = strings.TrimSuffix(route, "/")
	if !strings.HasPrefix(route, "/") {
		route = fmt.Sprintf("/%s", route)
	}

	return route
}

func concatenateEndpoints(basePath string, path string) string {
	if len(basePath) == 0 && len(path) == 0 {
		return "/"
	}

	if len(basePath) != 0 {
		basePath = sanitizePath(basePath)
	}

	if len(path) != 0 {
		path = sanitizePath(path)
	}

	return fmt.Sprintf("%s%s", basePath, path)
}
