package database

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"time"
)

func New() (*mongo.Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	uri, err := getMongoUri()
	if err != nil {
		return nil, err
	}
	
	cliOptions := options.Client().ApplyURI(uri)

	mongoClient, err := mongo.Connect(ctx, cliOptions)
	if err != nil {
		return nil, err
	}

	if err := mongoClient.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("ERROR pining mongodb Client: %s", err)
	}

	dbName := getDBName()
	dbInstance := mongoClient.Database(dbName)

	return dbInstance, nil
}

func getDBName() string {
	dbName := os.Getenv("DATABASE_NAME")
	if dbName == "" {
		dbName = "file_manager_test_db"
	}

	return dbName
}

func getMongoUri() (string, error) {
	uri := os.Getenv("MONGO_URI")
	if uri == "" {
		return "", errors.New("MONGO_URI environmental variable does not exist")
	}

	return uri, nil
}
