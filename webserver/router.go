package webserver

import (
	"file_manager/handlers"
	"github.com/julienschmidt/httprouter"
	"net/http"
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
	// Static Files
	staticFilesHandler := getStaticFilesHandler()
	router.Router.Handler("GET", "/static/*filepath", staticFilesHandler)

	// Auth
	router.Router.HandlerFunc("POST", "/api/auth/register", handler.Register)
	router.Router.HandlerFunc("POST", "/api/auth/login", handler.Login)

	// User
	router.Router.HandlerFunc("PUT", "/api/user/plan/change", handler.UpdateUserPlan)

	// File
	router.Router.HandlerFunc("POST", "/api/file/create", handler.UploadUserFile)
	router.Router.HandlerFunc("GET", "/api/file/get", handler.GetFiles)
	router.Router.HandlerFunc("DELETE", "/api/file/delete/:id", handler.DeleteFile)
	router.Router.HandlerFunc("PUT", "/api/file/rename/:id", handler.RenameFile)
	router.Router.HandlerFunc("GET", "/api/file/get/:id", handler.GetFile)
	router.Router.HandlerFunc("POST", "/api/file/get/:id", handler.GetFile)
	router.Router.HandlerFunc("POST", "/api/file/search", handler.SearchFiles)

	// File Settings
	router.Router.HandlerFunc("POST", "/api/file/settings/create/:id", handler.CreateFileSettings)

	// Approval
	router.Router.HandlerFunc("POST", "/api/approval/create", handler.CreateApproval)
	router.Router.HandlerFunc("PUT", "/api/approval/status", handler.UpdateApproval)

	// Team
	router.Router.HandlerFunc("POST", "/api/team/create", handler.CreateTeam)
	router.Router.HandlerFunc("POST", "/api/team/file/upload/:id", handler.UploadTeamFile)
	router.Router.HandlerFunc("GET", "/api/team/get/:id", handler.GetTeam)
}

func getStaticFilesHandler() http.Handler {
	dir := http.Dir("./")
	fileServer := http.FileServer(dir)
	return http.StripPrefix("/static/", fileServer)
}
