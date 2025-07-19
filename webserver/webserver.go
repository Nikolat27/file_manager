package webserver

import (
	"file_manager/handlers"
	"fmt"
	"net/http"
	"os"
)

type Server struct {
	Server *http.Server
	Port   string
}

func New(handler *handlers.Handler, port string) (*Server, error) {
	var srv = &Server{
		Port: port,
	}

	srv.setupHttpServer(handler)

	return srv, nil
}

func (srv *Server) setupHttpServer(handler *handlers.Handler) *Server {
	router := NewRouter(handler)

	srv.Server = &http.Server{
		Addr:    fmt.Sprintf(":%s", srv.Port),
		Handler: router.CoreRouter,
	}

	return srv
}

func (srv *Server) Run() error {
	fmt.Println("project started")

	if certFile, keyFile := getTlsCert(); certFile != "" || keyFile != "" {
		return srv.runHttpsServer(certFile, keyFile)
	} else {
		return srv.runHttpServer()
	}
}

func (srv *Server) Close() error {
	fmt.Println("project is closing...")
	return srv.Server.Close()
}

func (srv *Server) runHttpServer() error {
	return srv.Server.ListenAndServe()
}

func (srv *Server) runHttpsServer(certFile, keyFile string) error {
	return srv.Server.ListenAndServeTLS(certFile, keyFile)
}

func getTlsCert() (string, string) {
	certFile := os.Getenv("certFile")
	keyFile := os.Getenv("keyFile")

	return certFile, keyFile
}
