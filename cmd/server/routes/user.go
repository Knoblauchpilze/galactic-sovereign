package routes

import (
	"net/http"

	"github.com/KnoblauchPilze/user-service/pkg/communication"
	"github.com/KnoblauchPilze/user-service/pkg/db"
	"github.com/KnoblauchPilze/user-service/pkg/errors"
	"github.com/KnoblauchPilze/user-service/pkg/repositories"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func UserEndpoint(conn db.Connection) Routes {
	repo := repositories.NewUserRepository(conn)

	path := "/users"

	get := NewRoute("GET", path, wrapWithDb(getUser, repo))
	post := NewRoute("POST", path, wrapWithDb(createUser, repo))
	delete := NewRoute("DELETE", path, wrapWithDb(deleteUser, repo))

	var routes Routes
	routes = append(routes, get)
	routes = append(routes, post)
	routes = append(routes, delete)

	return routes
}

type userHttpHandler func(echo.Context, repositories.UserRepository) error

func wrapWithDb(handler userHttpHandler, repo repositories.UserRepository) echo.HandlerFunc {
	return func(c echo.Context) error {
		return handler(c, repo)
	}
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

func createUser(c echo.Context, repo repositories.UserRepository) error {
	return errors.NewCode(errors.NotImplementedCode)
}

func deleteUser(c echo.Context, repo repositories.UserRepository) error {
	return c.String(http.StatusOK, "{\"message\":\"deleted\"}")
}
