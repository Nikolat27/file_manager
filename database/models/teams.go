package models

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
	Plan        string               `json:"plan" bson:"plan"`
	OwnerId     primitive.ObjectID   `json:"owner_id" bson:"owner_id"`
	Users       []primitive.ObjectID `json:"users" bson:"users"`
	Admins      []primitive.ObjectID `json:"admins" bson:"admins"`
	Files       []primitive.ObjectID `json:"files" bson:"files"`
	StorageUsed int64                `json:"storage_used" bson:"storage_used"`
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

	// owner is considered admin too
	newTeam := &Team{
		Id:          teamId,
		Name:        name,
		Description: description,
		AvatarUrl:   avatarUrl,
		OwnerId:     ownerId,
		Admins:      []primitive.ObjectID{ownerId},
	}

	newId, err := team.db.Collection("teams").InsertOne(ctx, newTeam)
	if err != nil {
		return primitive.NilObjectID, err
	}

	return newId.InsertedID.(primitive.ObjectID), nil
}

func (team *TeamModel) ValidateAdmin(id, userId primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{
		"_id":    id,
		"admins": userId,
	}

	projection := bson.M{
		"_id": 1,
	}

	findOptions := options.FindOne()
	findOptions.SetProjection(projection)

	if err := team.db.Collection("teams").FindOne(ctx, filter, findOptions).Err(); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return errors.New("either user is not an admin or team with this id does not exist")
		}

		return err
	}

	return nil
}
