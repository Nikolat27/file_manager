package models

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type UserModel struct {
	DB *mongo.Database
}

type User struct {
	Id             primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Username       string             `json:"username" bson:"username"`
	Salt           string             `json:"salt" bson:"salt"`
	HashedPassword string             `json:"hashed_password" bson:"hashed_password"`
	CreatedAt      time.Time          `json:"created_at" bson:"created_at"`
}

const userCollectionName = "users"

func (user *UserModel) CreateUserInstance(username, salt, hashedPassword string) (primitive.ObjectID, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	fmt.Println(username, salt, hashedPassword)
	newUser := &User{
		Username:       username,
		Salt:           salt,
		HashedPassword: hashedPassword,
		CreatedAt:      time.Now(),
	}

	result, err := user.DB.Collection(userCollectionName).InsertOne(ctx, newUser)
	if err != nil {
		return primitive.NilObjectID, fmt.Errorf("ERROR create a new user instance: %s", err)
	}

	id, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return primitive.NilObjectID, fmt.Errorf("InsertedID is not an ObjectID, got: %T", result.InsertedID)
	}

	return id, nil
}

func (user *UserModel) FetchUserByUsername(username string) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{
		"username": username,
	}

	var userInstance User
	if err := user.DB.Collection(userCollectionName).FindOne(ctx, filter).Decode(&userInstance); err != nil {
		return nil, err
	}

	return &userInstance, nil
}
