package models

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type UserModel struct {
	db *mongo.Database
}

type User struct {
	Id              primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Username        string             `json:"username" bson:"username"`
	Plan            string             `json:"plan"`
	AvatarUrl       string             `json:"avatar_url" bson:"avatar_url"`
	TotalUploadSize int64              `json:"total_upload_size" bson:"total_upload_size"`
	Salt            string             `json:"salt" bson:"salt"`
	HashedPassword  string             `json:"hashed_password" bson:"hashed_password"`
	CreatedAt       time.Time          `json:"created_at" bson:"created_at"`
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

func (user *UserModel) GetById(id primitive.ObjectID) (*User, error) {
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

func (user *UserModel) Update(id primitive.ObjectID, updates bson.M) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	update := bson.M{
		"$set": updates,
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

func (user *UserModel) GetUsedStorage(id primitive.ObjectID) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{
		"_id": id,
	}

	projection := bson.M{
		"total_upload_size": 1,
	}

	findOptions := options.FindOne()
	findOptions.SetProjection(projection)

	var userInstance User
	if err := user.db.Collection(userCollectionName).FindOne(ctx, filter, findOptions).Decode(&userInstance); err != nil {
		return 0, err
	}

	return userInstance.TotalUploadSize, nil
}

func (user *UserModel) CheckExist(id primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{
		"_id": id,
	}

	projection := bson.M{
		"_id": 1,
	}

	findOptions := options.FindOne()
	findOptions.SetProjection(projection)

	if err := user.db.Collection(userCollectionName).FindOne(ctx, filter, findOptions).Err(); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return errors.New("user with this id does not exist")
		}

		return err
	}

	return nil
}
