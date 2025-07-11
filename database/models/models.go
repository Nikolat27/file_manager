package models

import "file_manager/database"

type Models struct {
	User UserModel
	File FileModel
}

func New(db *database.DB) *Models {
	return &Models{
		User: UserModel{DB: db},
		File: FileModel{DB: db},
	}
}
