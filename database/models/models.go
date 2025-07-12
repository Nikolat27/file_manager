package models

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type Models struct {
	User UserModel
	File FileModel
}

func New(db *mongo.Database) *Models {
	return &Models{
		User: UserModel{DB: db},
		File: FileModel{DB: db},
	}
}
