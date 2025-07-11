package database

import (
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type DB struct {
	MongoCli *mongo.Client
}

func New(uri string) (*DB, error) {
	cli, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	var db = &DB{
		MongoCli: cli,
	}

	return db, nil
}
