package models

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type UserModel struct {
	db *mongo.Database
}

type User struct {
	Id             primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Username       string             `json:"username" bson:"username"`
	Plan           string             `json:"plan"`
	Salt           string             `json:"salt" bson:"salt"`
	HashedPassword string             `json:"hashed_password" bson:"hashed_password"`
	CreatedAt      time.Time          `json:"created_at" bson:"created_at"`
}

const userCollectionName = "users"

func (user *UserModel) Create(username, plan, salt, hashedPassword string) (primitive.ObjectID, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// All plans: Free, Plus, Premium
	newUser := &User{
		Username:       username,
		Plan:           plan,
		Salt:           salt,
		HashedPassword: hashedPassword,
		CreatedAt:      time.Now(),
	}

	result, err := user.db.Collection(userCollectionName).InsertOne(ctx, newUser)
	if err != nil {
		return primitive.NilObjectID, fmt.Errorf("ERROR create a new user instance: %s", err)
	}

	id, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return primitive.NilObjectID, fmt.Errorf("InsertedID is not an ObjectID, got: %T", result.InsertedID)
	}

	return id, nil
}

func (user *UserModel) GetById(id string) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{
		"_id": id,
	}

	var userInstance User
	if err := user.db.Collection(userCollectionName).FindOne(ctx, filter).Decode(&userInstance); err != nil {
		return nil, err
	}

	return &userInstance, nil
}

func (user *UserModel) GetByUsername(username string) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{
		"username": username,
	}

	var userInstance User
	if err := user.db.Collection(userCollectionName).FindOne(ctx, filter).Decode(&userInstance); err != nil {
		return nil, err
	}

	return &userInstance, nil
}

func (user *UserModel) UpdatePlan(id primitive.ObjectID, newPlan string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	update := bson.M{
		"$set": bson.M{
			"plan": newPlan,
		},
	}

	result, err := user.db.Collection(userCollectionName).UpdateByID(ctx, id, update)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return errors.New("user with this Id does not exist")
		}

		return err
	}

	if result.MatchedCount == 0 {
		return errors.New("user with this Id does not exist")
	}

	if result.ModifiedCount == 0 {
		return errors.New("did not detect any change")
	}

	return nil
}
