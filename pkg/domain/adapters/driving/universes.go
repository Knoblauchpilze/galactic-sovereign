package drivingadapters

import (
	"log/slog"
	"net/http"

	"github.com/Knoblauchpilze/backend-toolkit/pkg/rest"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/adapters/driving/dtos"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/adapters/driving/mappers"
	domainerrors "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models/errors"
	drivingports "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/ports/driving"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func UniverseEndpoints(usecase drivingports.ForManagingUniverse) Routes {
	var out Routes

	handler := generateHandler(createUniverse, usecase)
	post := rest.NewRoute(http.MethodPost, "/universes", handler)
	out = append(out, post)

	handler = generateHandler(getUniverse, usecase)
	get := rest.NewRoute(http.MethodGet, "/universes/:id", handler)
	out = append(out, get)

	handler = generateHandler(listUniverses, usecase)
	list := rest.NewRoute(http.MethodGet, "/universes", handler)
	out = append(out, list)

	handler = generateHandler(deleteUniverse, usecase)
	delete := rest.NewRoute(http.MethodDelete, "/universes/:id", handler)
	out = append(out, delete)

	return out
}

// createUniverse godoc
//
//	@Summary		Create universe
//	@Description	Creates a universe.
//	@Tags			universes
//	@Produce		json
//	@Param			request	body		dtos.UniverseDtoRequest	true	"Universe payload"
//	@Success		201		{object}	rest.ResponseEnvelope[dtos.UniverseDtoResponse]
//	@Failure		400		{object}	rest.ResponseEnvelope[string]
//	@Failure		409		{object}	rest.ResponseEnvelope[string]
//	@Failure		500		{object}	rest.ResponseEnvelope[string]
//	@Router			/universes [post]
func createUniverse(c *gin.Context, usecase drivingports.ForManagingUniverse) {
	var inputDto dtos.UniverseDtoRequest
	err := c.Bind(&inputDto)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, "invalid universe syntax")
		return
	}

	request := mappers.ToUniverseCreationRequest(inputDto)
	universe, err := usecase.Create(c.Request.Context(), request)
	if err != nil {
		if err == domainerrors.ErrNameAlreadyTaken {
			c.AbortWithStatusJSON(http.StatusConflict, "name already used")
			return
		}

		logError(c.Request, "Failed to create universe", slog.Any("error", err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, "failed to create universe")
		return
	}

	out := mappers.ToUniverseResponse(universe)
	c.JSON(http.StatusCreated, out)
}

// getUniverse godoc
//
//	@Summary		Get universe
//	@Description	Returns a universe and related resources/buildings.
//	@Tags			universes
//	@Produce		json
//	@Param			id	path		string	true	"Universe id (UUID)"	Format(uuid)
//	@Success		200	{object}	rest.ResponseEnvelope[dtos.UniverseDtoResponse]
//	@Failure		400	{object}	rest.ResponseEnvelope[string]
//	@Failure		404	{object}	rest.ResponseEnvelope[string]
//	@Failure		500	{object}	rest.ResponseEnvelope[string]
//	@Router			/universes/{id} [get]
func getUniverse(c *gin.Context, usecase drivingports.ForManagingUniverse) {
	maybeId := c.Param("id")
	id, err := uuid.Parse(maybeId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, "invalid id syntax")
		return
	}

	universe, err := usecase.Get(c.Request.Context(), id)
	if err != nil {
		if err == domainerrors.ErrNotFound {
			c.AbortWithStatusJSON(http.StatusNotFound, "no such universe")
			return
		}

		logError(c.Request, "Failed to get universe", slog.Any("error", err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, "failed to get universe")
		return
	}

	out := mappers.ToUniverseResponse(universe)
	c.JSON(http.StatusOK, out)
}

// listUniverses godoc
//
//	@Summary		List universes
//	@Description	Returns all universes.
//	@Tags			universes
//	@Produce		json
//	@Success		200	{object}	rest.ResponseEnvelope[[]dtos.UniverseDtoResponse]
//	@Failure		500	{object}	rest.ResponseEnvelope[string]
//	@Router			/universes [get]
func listUniverses(c *gin.Context, usecase drivingports.ForManagingUniverse) {
	universes, err := usecase.List(c.Request.Context())
	if err != nil {
		logError(c.Request, "Failed to list universes", slog.Any("error", err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, "failed to list universes")
		return
	}

	out := mappers.ToUniversesResponse(universes)

	c.JSON(http.StatusOK, out)
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
//	@Failure		409	{object}	rest.ResponseEnvelope[string]
//	@Failure		500	{object}	rest.ResponseEnvelope[string]
//	@Router			/universes/{id} [delete]
func deleteUniverse(c *gin.Context, usecase drivingports.ForManagingUniverse) {
	maybeId := c.Param("id")
	id, err := uuid.Parse(maybeId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, "invalid id syntax")
		return
	}

	err = usecase.Delete(c.Request.Context(), id)
	if err != nil {
		if err == domainerrors.ErrUniverseIsNotEmpty {
			c.AbortWithStatusJSON(http.StatusConflict, "universe is not empty")
			return
		}

		logError(c.Request, "Failed to delete universe", slog.Any("error", err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, "failed to delete universe")
		return
	}

	c.Status(http.StatusNoContent)
}
