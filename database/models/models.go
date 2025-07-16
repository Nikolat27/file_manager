package models

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type Models struct {
	User     UserModel
	File     FileModel
	Approval ApprovalModel
}

func New(db *mongo.Database) *Models {
	return &Models{
		User:     UserModel{db: db},
		File:     FileModel{db: db},
		Approval: ApprovalModel{db: db},
	}
}
