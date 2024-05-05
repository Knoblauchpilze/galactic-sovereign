package controller

import (
	"net/http"

	"github.com/KnoblauchPilze/user-service/internal/users/service"
	"github.com/KnoblauchPilze/user-service/pkg/communication"
	"github.com/KnoblauchPilze/user-service/pkg/db"
	"github.com/KnoblauchPilze/user-service/pkg/errors"
	"github.com/KnoblauchPilze/user-service/pkg/rest"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func UserEndpoints(service service.UserService) rest.Routes {
	var out rest.Routes

	postHandler := fromRepositoriesAwareHttpHandler(createUser, service)
	post := rest.NewRoute(http.MethodPost, false, "/users", postHandler)
	out = append(out, post)

	getHandler := fromRepositoriesAwareHttpHandler(getUser, service)
	get := rest.NewResourceRoute(http.MethodGet, true, "/users", getHandler)
	out = append(out, get)

	listHandler := fromRepositoriesAwareHttpHandler(listUsers, service)
	list := rest.NewRoute(http.MethodGet, true, "/users", listHandler)
	out = append(out, list)

	updateHandler := fromRepositoriesAwareHttpHandler(updateUser, service)
	update := rest.NewResourceRoute(http.MethodPatch, true, "/users", updateHandler)
	out = append(out, update)

	deleteHandler := fromRepositoriesAwareHttpHandler(deleteUser, service)
	delete := rest.NewResourceRoute(http.MethodDelete, true, "/users", deleteHandler)
	out = append(out, delete)

	loginHandler := fromRepositoriesAwareHttpHandler(loginUser, service)
	login := rest.NewResourceRoute(http.MethodPost, false, "/users/sessions", loginHandler)
	out = append(out, login)

	logoutHandler := fromRepositoriesAwareHttpHandler(logoutUser, service)
	logout := rest.NewResourceRoute(http.MethodDelete, true, "/users/sessions", logoutHandler)
	out = append(out, logout)

	return out
}

func createUser(c echo.Context, service service.UserService) error {
	// https://echo.labstack.com/docs/binding
	var userDtoRequest communication.UserDtoRequest
	err := c.Bind(&userDtoRequest)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid user syntax")
	}

	out, err := service.Create(c.Request().Context(), userDtoRequest)
	if err != nil {
		if errors.IsErrorWithCode(err, db.DuplicatedKeySqlKey) {
			return c.JSON(http.StatusConflict, "Email already used")
		}

		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusCreated, out)
}

func getUser(c echo.Context, service service.UserService) error {
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

func listUsers(c echo.Context, service service.UserService) error {
	out, err := service.List(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, out)
}

func updateUser(c echo.Context, service service.UserService) error {
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

	out, err := service.Update(c.Request().Context(), id, userDtoRequest)
	if err != nil {
		if errors.IsErrorWithCode(err, db.NoMatchingSqlRows) {
			return c.JSON(http.StatusNotFound, "No such user")
		}

		if errors.IsErrorWithCode(err, db.OptimisticLockException) {
			return c.JSON(http.StatusConflict, "User is not up to date")
		}

		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, out)
}

func deleteUser(c echo.Context, service service.UserService) error {
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

func loginUser(c echo.Context, service service.UserService) error {
	maybeId := c.Param("id")
	id, err := uuid.Parse(maybeId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid id syntax")
	}

	out, err := service.Login(c.Request().Context(), id)
	if err != nil {
		if errors.IsErrorWithCode(err, db.NoMatchingSqlRows) {
			return c.JSON(http.StatusNotFound, "No such user")
		}

		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusCreated, out)
}

func logoutUser(c echo.Context, service service.UserService) error {
	maybeId := c.Param("id")
	id, err := uuid.Parse(maybeId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid id syntax")
	}

	err = service.Logout(c.Request().Context(), id)
	if err != nil {
		if errors.IsErrorWithCode(err, db.NoMatchingSqlRows) {
			return c.JSON(http.StatusNotFound, "No such user")
		}

		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.NoContent(http.StatusNoContent)
}
