package controller

import (
	"net/http"

	"github.com/KnoblauchPilze/user-service/internal/users/service"
	"github.com/KnoblauchPilze/user-service/pkg/communication"
	"github.com/KnoblauchPilze/user-service/pkg/db"
	"github.com/KnoblauchPilze/user-service/pkg/errors"
	"github.com/KnoblauchPilze/user-service/pkg/repositories"
	"github.com/KnoblauchPilze/user-service/pkg/rest"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func UserEndpoints(conn db.Connection, service service.UserService) rest.Routes {
	repo := repositories.NewUserRepository(conn)

	var out rest.Routes

	postHandler := generateEchoHandler(createUser, repo, service)
	post := rest.NewRoute(http.MethodPost, "/users", postHandler)
	out = append(out, post)

	getHandler := generateEchoHandler(getUser, repo, service)
	get := rest.NewResourceRoute(http.MethodGet, "/users", getHandler)
	out = append(out, get)

	listHandler := generateEchoHandler(listUsers, repo, service)
	list := rest.NewRoute(http.MethodGet, "/users", listHandler)
	out = append(out, list)

	updateHandler := generateEchoHandler(updateUser, repo, service)
	update := rest.NewResourceRoute(http.MethodPatch, "/users", updateHandler)
	out = append(out, update)

	deleteHandler := generateEchoHandler(deleteUser, repo, service)
	delete := rest.NewResourceRoute(http.MethodDelete, "/users", deleteHandler)
	out = append(out, delete)

	return out
}

func createUser(c echo.Context, repo repositories.UserRepository, service service.UserService) error {
	// https://echo.labstack.com/docs/binding
	var userDtoRequest communication.UserDtoRequest
	err := c.Bind(&userDtoRequest)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid user syntax")
	}

	user := communication.FromUserDtoRequest(userDtoRequest)
	created, err := repo.Create(c.Request().Context(), user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	out := communication.ToUserDtoResponse(created, true)
	return c.JSON(http.StatusCreated, out)
}

func getUser(c echo.Context, repo repositories.UserRepository, service service.UserService) error {
	maybeId := c.Param("id")
	id, err := uuid.Parse(maybeId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid id syntax")
	}

	out, err := service.Get(c.Request().Context(), id)
	if err != nil {
		if errors.IsErrorWithCode(err, db.NoMatchingSqlRows) {
			return c.JSON(http.StatusNotFound, "No such user")
		}

		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, out)
}

func listUsers(c echo.Context, repo repositories.UserRepository, service service.UserService) error {
	out, err := service.List(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, out)
}

func updateUser(c echo.Context, repo repositories.UserRepository, service service.UserService) error {
	maybeId := c.Param("id")
	id, err := uuid.Parse(maybeId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid id syntax")
	}

	var userDtoRequest communication.UserDtoRequest
	err = c.Bind(&userDtoRequest)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid user syntax")
	}

	user, err := repo.Get(c.Request().Context(), id)
	if err != nil {
		if errors.IsErrorWithCode(err, db.NoMatchingSqlRows) {
			return c.JSON(http.StatusNotFound, "No such user")
		}

		return c.JSON(http.StatusInternalServerError, err)
	}

	user.Email = userDtoRequest.Email
	user.Password = userDtoRequest.Password

	user, err = repo.Update(c.Request().Context(), user)
	if err != nil {
		if errors.IsErrorWithCode(err, db.OptimisticLockException) {
			return c.JSON(http.StatusConflict, "User is not up to date")
		}

		return c.JSON(http.StatusInternalServerError, err)
	}

	out := communication.ToUserDtoResponse(user, false)
	return c.JSON(http.StatusOK, out)
}

func deleteUser(c echo.Context, repo repositories.UserRepository, service service.UserService) error {
	maybeId := c.Param("id")
	id, err := uuid.Parse(maybeId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid id syntax")
	}

	err = service.Delete(c.Request().Context(), id)
	if err != nil {
		if errors.IsErrorWithCode(err, db.NoMatchingSqlRows) {
			return c.JSON(http.StatusNotFound, "No such user")
		}

		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.NoContent(http.StatusNoContent)
}
