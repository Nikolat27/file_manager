package handlers

import (
	"file_manager/database"
	"file_manager/token"
)

type Handler struct {
	PasetoMaker *token.PasetoMaker
	DB          *database.DB
}

func New(db *database.DB) (*Handler, error) {
	paseto, err := token.New()
	if err != nil {
		return nil, err
	}

	var handler = &Handler{
		PasetoMaker: paseto,
		DB:          db,
	}

	return handler, nil
}
