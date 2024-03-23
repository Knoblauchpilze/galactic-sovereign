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

	out := communication.FromUser(user)

	return c.JSON(http.StatusOK, out)
}

func createUser(echo.Context, repositories.UserRepository) error {
	return errors.NewCode(errors.NotImplementedCode)
}

func deleteUser(c echo.Context, _ repositories.UserRepository) error {
	return c.String(http.StatusOK, "{\"message\":\"deleted\"}")
}
