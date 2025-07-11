package models

import "file_manager/database"

type Models struct {
	User UserModel
}

func New(db *database.DB) *Models {
	return &Models{
		User: UserModel{DB: db},
	}
}
