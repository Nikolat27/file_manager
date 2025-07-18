package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type TeamModel struct {
	db *mongo.Database
}

type Team struct {
	Id          primitive.ObjectID   `json:"id,omitempty" bson:"_id,omitempty"`
	Name        string               `json:"name" bson:"name"`
	Description string               `json:"description" bson:"description"`
	AvatarUrl   string               `json:"avatar_url" bson:"avatar_url"`
	OwnerId     primitive.ObjectID   `json:"owner_id" bson:"owner_id"`
	Users       []primitive.ObjectID `json:"users" bson:"users"`
	Admins      []primitive.ObjectID `json:"admins" bson:"admins"`
	Files       []primitive.ObjectID `json:"files" bson:"files"`
	CreatedAt   time.Time            `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time            `json:"updated_at" bson:"updated_at"`
}
