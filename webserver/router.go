package webserver

import (
	"file_manager/handlers"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type AppRouter struct {
	CoreRouter *httprouter.Router
}

func NewRouter(handler *handlers.Handler) *AppRouter {
	routerInstance := &AppRouter{
		CoreRouter: httprouter.New(),
	}

	routerInstance.registerRoutes(handler)

	return routerInstance
}

// do not use OPTIONS method. Allowed Methods: GET, POST, PUT, DELETE
func (router *AppRouter) registerRoutes(handler *handlers.Handler) {
	router.registerStaticRoutes()

	router.registerAuthRoutes(handler)
	router.registerUserRoutes(handler)

	router.registerFileRoutes(handler)
	router.registerFileSettingsRoutes(handler)

	router.registerFolderRoutes(handler)

	router.registerApprovalRoutes(handler)

	router.registerTeamRoutes(handler)
}

// registerStaticRoutes -> Static Files
func (router *AppRouter) registerStaticRoutes() {
	staticFilesHandler := getStaticFilesHandler()
	router.CoreRouter.Handler("GET", "/static/*filepath", staticFilesHandler)
}

// registerAuthRoutes -> Auth
func (router *AppRouter) registerAuthRoutes(handler *handlers.Handler) {
	router.CoreRouter.HandlerFunc("POST", "/api/auth/register", handler.Register)
	router.CoreRouter.HandlerFunc("POST", "/api/auth/login", handler.Login)
}

// registerUserRoutes -> Users
func (router *AppRouter) registerUserRoutes(handler *handlers.Handler) {
	router.CoreRouter.HandlerFunc("PUT", "/api/user/plan/change", handler.UpdateUserPlan)
}

// registerFileRoutes -> Files
func (router *AppRouter) registerFileRoutes(handler *handlers.Handler) {
	router.CoreRouter.HandlerFunc("POST", "/api/file/create", handler.UploadUserFile)
	router.CoreRouter.HandlerFunc("GET", "/api/file/get", handler.GetFiles)
	router.CoreRouter.HandlerFunc("DELETE", "/api/file/delete/:id", handler.DeleteFile)
	router.CoreRouter.HandlerFunc("PUT", "/api/file/rename/:id", handler.RenameFile)
	router.CoreRouter.HandlerFunc("POST", "/api/file/search", handler.SearchFiles)
	router.CoreRouter.HandlerFunc("GET", "/api/file/download/:id", handler.DownloadFile)

	// GET method (for password-less files)
	router.CoreRouter.HandlerFunc("GET", "/api/file/get/:id", handler.GetFile)
	// POST method (for password requirable files)
	router.CoreRouter.HandlerFunc("POST", "/api/file/get/:id", handler.GetFile)
}

// registerFileRoutes -> Folder
func (router *AppRouter) registerFolderRoutes(handler *handlers.Handler) {
	router.CoreRouter.HandlerFunc("GET", "/api/folder/get", handler.GetFoldersList)
	router.CoreRouter.HandlerFunc("POST", "/api/folder/create", handler.CreateFolder)
	router.CoreRouter.HandlerFunc("GET", "/api/folder/get/:id", handler.GetFolderContents)
	router.CoreRouter.HandlerFunc("PUT", "/api/folder/rename/:id", handler.RenameFolder)
	router.CoreRouter.HandlerFunc("DELETE", "/api/folder/delete/:id", handler.DeleteFolder)
}

// registerFileSettingsRoutes -> File Settings
func (router *AppRouter) registerFileSettingsRoutes(handler *handlers.Handler) {
	router.CoreRouter.HandlerFunc("POST", "/api/file/settings/create/:id", handler.CreateFileSettings)
}

// registerApprovalRoutes -> Approvals
func (router *AppRouter) registerApprovalRoutes(handler *handlers.Handler) {
	router.CoreRouter.HandlerFunc("POST", "/api/approval/create", handler.CreateApproval)
	router.CoreRouter.HandlerFunc("PUT", "/api/approval/status", handler.UpdateApproval)
}

// registerTeamRoutes -> Teams
func (router *AppRouter) registerTeamRoutes(handler *handlers.Handler) {
	router.CoreRouter.HandlerFunc("POST", "/api/team/create", handler.CreateTeam)
	router.CoreRouter.HandlerFunc("POST", "/api/team/file/upload/:id", handler.UploadTeamFile)
	router.CoreRouter.HandlerFunc("GET", "/api/team/get/:id", handler.GetTeam)
	router.CoreRouter.HandlerFunc("DELETE", "/api/team/delete/:id", handler.DeleteTeam)
	router.CoreRouter.HandlerFunc("POST", "/api/team/user/add/:id", handler.AddUserToTeam)
	router.CoreRouter.HandlerFunc("PUT", "/api/team/plan/update/:id", handler.UpdateTeamPlan)
}

func getStaticFilesHandler() http.Handler {
	dir := http.Dir("./")
	fileServer := http.FileServer(dir)
	return http.StripPrefix("/static/", fileServer)
}
