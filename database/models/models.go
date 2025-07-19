package models

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type Models struct {
	User         UserModel
	File         FileModel
	Folder       FolderModel
	FileSettings FileSettingModel
	Approval     ApprovalModel
	Team         TeamModel
}

func New(db *mongo.Database) *Models {
	return &Models{
		User:         UserModel{db: db},
		File:         FileModel{db: db},
		Folder:       FolderModel{db: db},
		FileSettings: FileSettingModel{db: db},
		Approval:     ApprovalModel{db: db},
		Team:         TeamModel{db: db},
	}
}
