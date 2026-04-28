package controller

import (
	"net/http"

	"github.com/Knoblauchpilze/backend-toolkit/pkg/db"
	"github.com/Knoblauchpilze/backend-toolkit/pkg/db/pgx"
	"github.com/Knoblauchpilze/backend-toolkit/pkg/errors"
	"github.com/Knoblauchpilze/backend-toolkit/pkg/rest"
	"github.com/Knoblauchpilze/galactic-sovereign/internal/service"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/communication"
	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
)

func UniverseEndpoints(service service.UniverseService) rest.Routes {
	var out rest.Routes

	postHandler := fromUniverseServiceAwareHttpHandler(createUniverse, service)
	post := rest.NewRoute(http.MethodPost, "/universes", postHandler)
	out = append(out, post)

	getHandler := fromUniverseServiceAwareHttpHandler(getUniverse, service)
	get := rest.NewRoute(http.MethodGet, "/universes/:id", getHandler)
	out = append(out, get)

	listHandler := fromUniverseServiceAwareHttpHandler(listUniverses, service)
	list := rest.NewRoute(http.MethodGet, "/universes", listHandler)
	out = append(out, list)

	deleteHandler := fromUniverseServiceAwareHttpHandler(deleteUniverse, service)
	delete := rest.NewRoute(http.MethodDelete, "/universes/:id", deleteHandler)
	out = append(out, delete)

	return out
}

// createUniverse godoc
//
//	@Summary		Create universe
//	@Description	Creates a universe.
//	@Tags			universes
//	@Accept			json
//	@Produce		json
//	@Param			request	body		communication.UniverseDtoRequest	true	"Universe payload"
//	@Success		201		{object}	rest.ResponseEnvelope[communication.UniverseDtoResponse]
//	@Failure		400		{object}	rest.ResponseEnvelope[string]
//	@Failure		409		{object}	rest.ResponseEnvelope[string]
//	@Failure		500		{object}	rest.ResponseEnvelope[string]
func createUniverse(c *echo.Context, s service.UniverseService) error {
	var universeDtoRequest communication.UniverseDtoRequest
	err := c.Bind(&universeDtoRequest)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid universe syntax")
	}

	out, err := s.Create(c.Request().Context(), universeDtoRequest)
	if err != nil {
		if errors.IsErrorWithCode(err, pgx.UniqueConstraintViolation) {
			return c.JSON(http.StatusConflict, "Name already used")
		}

		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusCreated, out)
}

// getUniverse godoc
//
//	@Summary		Get universe
//	@Description	Returns a universe and related resources/buildings.
//	@Tags			universes
//	@Produce		json
//	@Param			id	path		string	true	"Universe id (UUID)"	Format(uuid)
//	@Success		200	{object}	rest.ResponseEnvelope[communication.FullUniverseDtoResponse]
//	@Failure		400	{object}	rest.ResponseEnvelope[string]
//	@Failure		404	{object}	rest.ResponseEnvelope[string]
//	@Failure		500	{object}	rest.ResponseEnvelope[string]
//	@Router			/universes/{id} [get]
func getUniverse(c *echo.Context, s service.UniverseService) error {
	maybeId := c.Param("id")
	id, err := uuid.Parse(maybeId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid id syntax")
	}

	out, err := s.Get(c.Request().Context(), id)
	if err != nil {
		if errors.IsErrorWithCode(err, db.NoMatchingRows) {
			return c.JSON(http.StatusNotFound, "No such universe")
		}

		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, out)
}

// listUniverses godoc
//
//	@Summary		List universes
//	@Description	Returns all universes.
//	@Tags			universes
//	@Produce		json
//	@Success		200	{object}	rest.ResponseEnvelope[[]communication.UniverseDtoResponse]
//	@Failure		500	{object}	rest.ResponseEnvelope[string]
//	@Router			/universes [get]
func listUniverses(c *echo.Context, s service.UniverseService) error {
	universes, err := s.List(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	out, err := marshalNilToEmptySlice(universes)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSONBlob(http.StatusOK, out)
}

// deleteUniverse godoc
//
//	@Summary		Delete universe
//	@Description	Deletes a universe by id.
//	@Tags			universes
//	@Produce		json
//	@Param			id	path		string	true	"Universe id (UUID)"	Format(uuid)
//	@Success		204	{string}	string
//	@Failure		400	{object}	rest.ResponseEnvelope[string]
//	@Failure		404	{object}	rest.ResponseEnvelope[string]
//	@Failure		500	{object}	rest.ResponseEnvelope[string]
//	@Router			/universes/{id} [delete]
func deleteUniverse(c *echo.Context, s service.UniverseService) error {
	maybeId := c.Param("id")
	id, err := uuid.Parse(maybeId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid id syntax")
	}

	err = s.Delete(c.Request().Context(), id)
	if err != nil {
		if errors.IsErrorWithCode(err, db.NoMatchingRows) {
			return c.JSON(http.StatusNotFound, "No such universe")
		}

		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.NoContent(http.StatusNoContent)
}
