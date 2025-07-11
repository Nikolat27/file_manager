package models

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type Models struct {
	User UserModel
}

func New(db *mongo.Database) *Models {
	return &Models{
		User: UserModel{DB: db},
	}
}
