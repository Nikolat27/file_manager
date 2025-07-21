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

type FolderModel struct {
	db *mongo.Database
}

type Folder struct {
	Id        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	OwnerId   primitive.ObjectID `json:"owner_id" bson:"owner_id"`
	Name      string             `json:"name" bson:"name"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

func (folder *FolderModel) Create(ownerId primitive.ObjectID, name string) (primitive.ObjectID, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	newFolder := &Folder{
		OwnerId:   ownerId,
		Name:      name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	result, err := folder.db.Collection("folders").InsertOne(ctx, newFolder)
	if err != nil {
		return primitive.NilObjectID, err
	}

	return result.InsertedID.(primitive.ObjectID), nil
}

func (folder *FolderModel) Delete(id primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{
		"_id": id,
	}

	result, err := folder.db.Collection("folders").DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return errors.New("no folder found with this id")
	}

	return nil
}

func (folder *FolderModel) Rename(id primitive.ObjectID, updates any) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	update := bson.M{
		"$set": updates,
	}

	result, err := folder.db.Collection("folders").UpdateByID(ctx, id, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return errors.New("no folder found with this id")
	}

	if result.ModifiedCount == 0 {
		return errors.New("did not detect any change")
	}

	return nil
}

func (folder *FolderModel) GetAll(ownerId primitive.ObjectID, page, pageSize int64) ([]Folder, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{
		"owner_id": ownerId,
	}

	findOptions := options.Find()
	findOptions.SetSkip((page - 1) * pageSize)
	findOptions.SetLimit(pageSize)

	cursor, err := folder.db.Collection("folders").Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}

	var folders []Folder
	if err := cursor.All(ctx, &folders); err != nil {
		return nil, err
	}

	return folders, nil
}

// Get -> Returns One
func (folder *FolderModel) Get(filter, projection bson.M) (*Folder, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	findOptions := options.FindOne()
	findOptions.SetProjection(projection)

	var folderInstance Folder

	if err := folder.db.Collection("folders").FindOne(ctx, filter, findOptions).Decode(&folderInstance); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("folder with this filter does not exist")
		}

		return nil, err
	}

	return &folderInstance, nil
}
