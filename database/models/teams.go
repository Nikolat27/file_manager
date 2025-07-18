package models

import (
	"context"
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

func (team *TeamModel) Create(id, ownerId primitive.ObjectID, name, description, avatarUrl string) (primitive.ObjectID, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var teamId primitive.ObjectID
	if id != primitive.NilObjectID {
		teamId = id
	}

	newTeam := &Team{
		Id:          teamId,
		Name:        name,
		Description: description,
		AvatarUrl:   avatarUrl,
		OwnerId:     ownerId,
	}

	newId, err := team.db.Collection("teams").InsertOne(ctx, newTeam)
	if err != nil {
		return primitive.NilObjectID, err
	}

	return newId.InsertedID.(primitive.ObjectID), nil
}
