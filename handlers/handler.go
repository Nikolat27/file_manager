package handlers

import (
	"file_manager/database/models"
	"file_manager/token"
)

type envelope map[string]any

type Handler struct {
	PasetoMaker *token.PasetoMaker
	Models      *models.Models
}

func New(models *models.Models) (*Handler, error) {
	paseto, err := token.New()
	if err != nil {
		return nil, err
	}

	var handler = &Handler{
		PasetoMaker: paseto,
		Models:      models,
	}

	return handler, nil
}
