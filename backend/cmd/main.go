package main

import (
	"file_manager/cronjobs"
	"file_manager/database"
	"file_manager/database/models"
	"file_manager/handlers"
	"file_manager/webserver"
	"fmt"
)

func main() {
	db, err := database.New()
	if err != nil {
		panic(fmt.Errorf("ERROR initializing database client: %s", err))
	}

	newModels := models.New(db)

	handler, err := handlers.New(newModels)
	if err != nil {
		panic(fmt.Errorf("ERROR creating the handler: %s", err))
	}

	srv, err := webserver.New(handler, "8000")
	if err != nil {
		panic(fmt.Errorf("ERROR creating the server: %s", err))
	}

	defer func() {
		if err := srv.Close(); err != nil {
			panic(fmt.Errorf("ERROR closing http server: %s", err))
		}
	}()

	// background
	go func() {
		cronInstance := cronjobs.New()
		cronInstance.Start(handler)
	}()

	if err := srv.Run(); err != nil {
		panic(fmt.Errorf("ERROR running http server: %s", err))
	}
}
