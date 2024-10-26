package rest

import (
	"fmt"
	"regexp"
	"strings"
)

var multiSlashRegex = regexp.MustCompile("[/]+")

func sanitizePath(route string) string {
	route = fmt.Sprintf("/%s", route)
	route = multiSlashRegex.ReplaceAllString(route, "/")
	route = strings.TrimSuffix(route, "/")

	if len(route) == 0 {
		return "/"
	}

	return route
}

func ConcatenateEndpoints(basePath string, path string) string {
	concatenated := fmt.Sprintf("/%s/%s", basePath, path)
	return sanitizePath(concatenated)
}
