package models

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type Models struct {
	User     userModel
	File     fileModel
	Approval approvalModel
}

func New(db *mongo.Database) *Models {
	return &Models{
		User:     userModel{db: db},
		File:     fileModel{db: db},
		Approval: approvalModel{db: db},
	}
}
