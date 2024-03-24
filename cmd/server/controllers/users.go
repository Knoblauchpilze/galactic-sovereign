package controllers

import (
	"net/http"

	"github.com/KnoblauchPilze/user-service/cmd/server/routes"
	"github.com/KnoblauchPilze/user-service/pkg/communication"
	"github.com/KnoblauchPilze/user-service/pkg/db"
	"github.com/KnoblauchPilze/user-service/pkg/errors"
	"github.com/KnoblauchPilze/user-service/pkg/repositories"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func UserEndpoints(conn db.Connection) routes.Routes {
	repo := repositories.NewUserRepository(conn)

	var out routes.Routes

	getHandler := generateEchoHandler(getUser, repo)
	get := routes.NewResourceRoute(http.MethodGet, "/users", getHandler)
	out = append(out, get)

	postHandler := generateEchoHandler(createUser, repo)
	post := routes.NewRoute(http.MethodPost, "/users", postHandler)
	out = append(out, post)

	deleteHandler := generateEchoHandler(deleteUser, repo)
	delete := routes.NewResourceRoute(http.MethodDelete, "/users", deleteHandler)
	out = append(out, delete)

	return out
}

func getUser(c echo.Context, repo repositories.UserRepository) error {
	maybeId := c.Param("id")
	id, err := uuid.Parse(maybeId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid id syntax")
	}

	user, err := repo.Get(id)
	if err != nil {
		if errors.IsErrorWithCode(err, db.NoMatchingSqlRows) {
			return c.JSON(http.StatusNotFound, "No such user")
		}

		return c.JSON(http.StatusInternalServerError, err)
	}

	out := communication.ToUserDtoResponse(user)
	return c.JSON(http.StatusOK, out)
}

func createUser(c echo.Context, repo repositories.UserRepository) error {
	// https://echo.labstack.com/docs/binding
	var userDtoRequest communication.UserDtoRequest
	err := c.Bind(&userDtoRequest)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid user syntax")
	}

	user := communication.FromUserDtoRequest(userDtoRequest)
	err = repo.Create(user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	out := communication.ToUserDtoResponse(user)
	return c.JSON(http.StatusCreated, out)
}

func deleteUser(c echo.Context, repo repositories.UserRepository) error {
	maybeId := c.Param("id")
	id, err := uuid.Parse(maybeId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid id syntax")
	}

	err = repo.Delete(id)
	if err != nil {
		if errors.IsErrorWithCode(err, db.NoMatchingSqlRows) {
			return c.JSON(http.StatusNotFound, "No such user")
		}

		return c.JSON(http.StatusInternalServerError, err)
	}

	c.Response().Status = http.StatusNoContent

	return nil
}
