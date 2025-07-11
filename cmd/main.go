package main

import (
	"errors"
	"file_manager/database"
	"file_manager/handlers"
	"file_manager/webserver"
	"fmt"
	"github.com/joho/godotenv"
	"os"
)

func main() {
	if err := godotenv.Load(); err != nil {
		panic(fmt.Errorf("ERROR loading environmental variables: %s", err))
	}

	uri, err := getMongoUri()
	if err != nil {
		panic(fmt.Errorf("ERROR %s", err))
	}

	db, err := database.New(uri)
	if err != nil {
		panic(fmt.Errorf("ERROR initializing database client: %s", err))
	}

	handler, err := handlers.New(db)
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

	if err := srv.Run(); err != nil {
		panic(fmt.Errorf("ERROR running http server: %s", err))
	}
}

func getMongoUri() (string, error) {
	uri := os.Getenv("MONGO_URI")
	if uri == "" {
		return "", errors.New("MONGO_URI environmental variable does not exist")
	}

	return uri, nil
}
