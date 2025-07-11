package models

import (
	"context"
	"file_manager/database"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
	"time"
)

type UserModel struct {
	DB *database.DB
}

type User struct {
	Id             primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Username       string             `json:"username" bson:"username"`
	Salt           string             `json:"salt" bson:"salt"`
	HashedPassword string             `json:"hashed_password" bson:"hashed_password"`
	CreatedAt      time.Time          `json:"created_at" bson:"created_at"`
}

const userCollectionName = "users"

func (user *UserModel) CreateUserInstance(username, salt, hashedPassword string) (primitive.ObjectID, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	newUser := bson.M{
		"username":        username,
		"salt":            salt,
		"hashed_password": hashedPassword,
	}

	result, err := user.DB.MongoDB.Collection(userCollectionName).InsertOne(ctx, newUser)
	if err != nil {
		return primitive.NilObjectID, fmt.Errorf("ERROR create a new user instance: %s", err)
	}

	return result.InsertedID.(primitive.ObjectID), nil
}
