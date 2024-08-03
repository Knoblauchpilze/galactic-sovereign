package controller

import (
	"net/http"

	"github.com/KnoblauchPilze/user-service/internal/service"
	"github.com/KnoblauchPilze/user-service/pkg/rest"
	"github.com/labstack/echo/v4"
)

func ResourceEndpoints(service service.ResourceService) rest.Routes {
	var out rest.Routes

	listHandler := fromResourceServiceAwareHttpHandler(listResources, service)
	list := rest.NewRoute(http.MethodGet, false, "/resources", listHandler)
	out = append(out, list)

	return out
}

func listResources(c echo.Context, s service.ResourceService) error {
	resources, err := s.List(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	out, err := marshalNilToEmptySlice(resources)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSONBlob(http.StatusOK, out)
}
