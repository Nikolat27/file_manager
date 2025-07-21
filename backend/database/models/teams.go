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
		Plan:        "free",
		Admins:      []primitive.ObjectID{ownerId},
		Users:       []primitive.ObjectID{ownerId},
		CreatedAt:   time.Now(),
	}

	newId, err := team.db.Collection("teams").InsertOne(ctx, newTeam)
	if err != nil {
		return primitive.NilObjectID, err
	}

	return newId.InsertedID.(primitive.ObjectID), nil
}

// Get -> Returns One
func (team *TeamModel) Get(filter, projection bson.M) (*Team, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	findOptions := options.FindOne()
	findOptions.SetProjection(projection)

	var teamInstance Team
	if err := team.db.Collection("teams").FindOne(ctx, filter, findOptions).Decode(&teamInstance); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("team with this id does not exist")
		}

		return nil, err
	}

	return &teamInstance, nil
}

func (team *TeamModel) Update(id primitive.ObjectID, updates bson.M) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	update := bson.M{
		"$set": updates,
	}

	result, err := team.db.Collection("teams").UpdateByID(ctx, id, update)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return errors.New("team with this Id does not exist")
		}

		return err
	}

	if result.MatchedCount == 0 {
		return errors.New("team with this Id does not exist")
	}

	if result.ModifiedCount == 0 {
		return errors.New("did not detect any change")
	}

	return nil
}

func (team *TeamModel) Delete(id primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{
		"_id": id,
	}

	result, err := team.db.Collection("teams").DeleteOne(ctx, filter)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return errors.New("team with this Id does not exist")
		}

		return err
	}

	if result.DeletedCount == 0 {
		return errors.New("can`t delete this team...")
	}

	return nil
}
