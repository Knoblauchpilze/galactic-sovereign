package routes

import "net/http"

type Route struct {
	Path        string
	GetRoute    http.HandlerFunc
	PostRoute   http.HandlerFunc
	DeleteRoute http.HandlerFunc
}
