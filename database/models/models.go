package models

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type Models struct {
	User         UserModel
	File         FileModel
	FileSettings FileSettingModel
	Approval     ApprovalModel
	Teams        TeamModel
}

func New(db *mongo.Database) *Models {
	return &Models{
		User:         UserModel{db: db},
		File:         FileModel{db: db},
		FileSettings: FileSettingModel{db: db},
		Approval:     ApprovalModel{db: db},
		Teams:        TeamModel{db: db},
	}
}
