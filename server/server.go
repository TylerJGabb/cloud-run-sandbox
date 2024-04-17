package server

import (
	"cloud-run-sandbox/middleware"
	"net/http"
)

type App interface {
	Start(string)
	Use(middleware.Middleware)
	Handle(string, http.Handler)
}

