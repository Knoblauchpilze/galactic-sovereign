package routes

import (
	"fmt"
	"strings"
)

type slashPrefix int

const (
	trimPrefix slashPrefix = 0
	addPrefix  slashPrefix = 1
)

func sanitizePath(route string, slashPrefixMode slashPrefix) string {
	route = strings.TrimSuffix(route, "/")

	switch slashPrefixMode {
	case trimPrefix:
		route = strings.TrimPrefix(route, "/")
	case addPrefix:
		if !strings.HasPrefix(route, "/") {
			route = fmt.Sprintf("/%s", route)
		}
	}

	return route
}

func concatenateEndpoints(endpoint string, path string) string {
	endpoint = sanitizePath(endpoint, addPrefix)
	path = sanitizePath(path, trimPrefix)

	return fmt.Sprintf("%s/%s", strings.TrimSuffix(endpoint, "/"), path)
}
