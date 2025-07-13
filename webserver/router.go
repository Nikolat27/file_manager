package webserver

import (
	"file_manager/handlers"
	"github.com/julienschmidt/httprouter"
)

type Router struct {
	Router *httprouter.Router
}

func NewRouter(handler *handlers.Handler) *Router {
	newRouter := httprouter.New()

	var router = &Router{
		Router: newRouter,
	}

	router.initRoutes(handler)

	return router
}

func (router *Router) initRoutes(handler *handlers.Handler) {
	// Auth
	router.Router.HandlerFunc("POST", "/api/auth/register", handler.Register)
	router.Router.HandlerFunc("POST", "/api/auth/login", handler.Login)

	// File
	router.Router.HandlerFunc("POST", "/api/file/create", handler.CreateFile)
	router.Router.HandlerFunc("GET", "/api/file/get", handler.GetUserFiles)
}
