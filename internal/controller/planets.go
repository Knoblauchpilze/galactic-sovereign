package controller

import (
	"net/http"

	"github.com/Knoblauchpilze/backend-toolkit/pkg/db"
	"github.com/Knoblauchpilze/backend-toolkit/pkg/errors"
	"github.com/Knoblauchpilze/backend-toolkit/pkg/rest"
	"github.com/Knoblauchpilze/galactic-sovereign/internal/service"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/communication"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/game"
	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
)

func PlanetEndpoints(planetService service.PlanetService,
	actionService game.ActionService,
	planetResourceService game.PlanetResourceService) rest.Routes {
	var out rest.Routes

	postHandler := fromPlanetServiceAwareHttpHandler(createPlanet, planetService)
	post := rest.NewRoute(http.MethodPost, "/planets", postHandler)
	out = append(out, post)

	getHandler := fromPlanetServiceAwareHttpHandler(getPlanet, planetService)
	get := game.NewResourceRoute(http.MethodGet, "/planets", getHandler, actionService, planetResourceService)
	out = append(out, get)

	listHandler := fromPlanetServiceAwareHttpHandler(listPlanets, planetService)
	list := rest.NewRoute(http.MethodGet, "/planets", listHandler)
	out = append(out, list)

	deleteHandler := fromPlanetServiceAwareHttpHandler(deletePlanet, planetService)
	delete := game.NewResourceRoute(http.MethodDelete, "/planets", deleteHandler, actionService, planetResourceService)
	out = append(out, delete)

	return out
}

// createPlanet godoc
//
//	@Summary		Create planet
//	@Description	Creates a planet from the provided payload.
//	@Tags			planets
//	@Accept			json
//	@Produce		json
//	@Param			request	body		PlanetRequestDoc	true	"Planet payload"
//	@Success		201		{object}	PlanetResponseDoc
//	@Failure		400		{string}	string	"Invalid planet syntax"
//	@Failure		500		{object}	ToolkitErrorDoc
//	@Router			/planets [post]
func createPlanet(c *echo.Context, s service.PlanetService) error {
	var planetDtoRequest communication.PlanetDtoRequest
	err := c.Bind(&planetDtoRequest)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid planet syntax")
	}

	out, err := s.Create(c.Request().Context(), planetDtoRequest)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusCreated, out)
}

// getPlanet godoc
//
//	@Summary		Get planet
//	@Description	Returns a planet and all related game data.
//	@Tags			planets
//	@Produce		json
//	@Param			id	path		string	true	"Planet id (UUID)"
//	@Success		200	{object}	FullPlanetResponseDoc
//	@Failure		400	{string}	string	"Invalid id syntax"
//	@Failure		404	{string}	string	"No such planet"
//	@Failure		500	{object}	ToolkitErrorDoc
//	@Router			/planets/{id} [get]
func getPlanet(c *echo.Context, s service.PlanetService) error {
	maybeId := c.Param("id")
	id, err := uuid.Parse(maybeId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid id syntax")
	}

	out, err := s.Get(c.Request().Context(), id)
	if err != nil {
		if errors.IsErrorWithCode(err, db.NoMatchingRows) {
			return c.JSON(http.StatusNotFound, "No such planet")
		}

		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, out)
}

// listPlanets godoc
//
//	@Summary		List planets
//	@Description	Returns planets, optionally filtered by player id.
//	@Tags			planets
//	@Produce		json
//	@Param			player	query		string	false	"Player id (UUID)"
//	@Success		200		{array}		PlanetResponseDoc
//	@Failure		400		{string}	string	"Invalid id syntax"
//	@Failure		500		{object}	ToolkitErrorDoc
//	@Router			/planets [get]
func listPlanets(c *echo.Context, s service.PlanetService) error {
	exists, playerId, err := fetchIdFromQueryParam("player", c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid id syntax")
	}

	var planets []communication.PlanetDtoResponse

	if exists {
		planets, err = s.ListForPlayer(c.Request().Context(), playerId)
	} else {
		planets, err = s.List(c.Request().Context())
	}

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	out, err := marshalNilToEmptySlice(planets)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSONBlob(http.StatusOK, out)
}

// deletePlanet godoc
//
//	@Summary		Delete planet
//	@Description	Deletes a planet by id.
//	@Tags			planets
//	@Produce		json
//	@Param			id	path		string	true	"Planet id (UUID)"
//	@Success		204	{string}	string
//	@Failure		400	{string}	string	"Invalid id syntax"
//	@Failure		404	{string}	string	"No such planet"
//	@Failure		500	{object}	ToolkitErrorDoc
//	@Router			/planets/{id} [delete]
func deletePlanet(c *echo.Context, s service.PlanetService) error {
	maybeId := c.Param("id")
	id, err := uuid.Parse(maybeId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid id syntax")
	}

	err = s.Delete(c.Request().Context(), id)
	if err != nil {
		if errors.IsErrorWithCode(err, db.NoMatchingRows) {
			return c.JSON(http.StatusNotFound, "No such planet")
		}

		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.NoContent(http.StatusNoContent)
}
