package database

import (
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"os"
)

type DB struct {
	MongoDB *mongo.Database
}

func New(uri string) (*DB, error) {
	cli, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}
	
	dbName := getDBName()
	dbInstance := cli.Database(dbName)
	
	var db = &DB{
		MongoDB: dbInstance,
	}

	return db, nil
}

func getDBName() string {
	dbName := os.Getenv("DATABASE_NAME")
	if dbName == "" {
		dbName = "file_manager_test_db"
	}

	return dbName
}
