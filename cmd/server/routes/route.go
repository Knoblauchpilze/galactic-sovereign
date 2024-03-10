package routes

import "github.com/labstack/echo/v4"

type Route struct {
	Path        string
	GetRoute    echo.HandlerFunc
	PostRoute   echo.HandlerFunc
	DeleteRoute echo.HandlerFunc
}
