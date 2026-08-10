package drivingadapters

import (
	"github.com/Knoblauchpilze/backend-toolkit/pkg/rest"
	"github.com/gin-gonic/gin"
)

type Routes []*rest.Route

type drivingAdapter[T any] = func(*gin.Context, T)

func generateHandler[T any](handler drivingAdapter[T], usecase T) gin.HandlerFunc {
	return func(c *gin.Context) {
		handler(c, usecase)
	}
}
