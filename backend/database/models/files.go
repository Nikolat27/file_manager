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

type FileModel struct {
	db *mongo.Database
}

type File struct {
	Id        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	OwnerId   primitive.ObjectID `json:"owner_id" bson:"owner_id"`
	TeamId    primitive.ObjectID `json:"team_id" bson:"team_id"`
	FolderId  primitive.ObjectID `json:"folder_id" bson:"folder_id"`
	Name      string             `json:"name" bson:"name"`
	Address   string             `json:"address" bson:"address"`
	ExpireAt  time.Time          `json:"expire_at" bson:"expire_at"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
}

func (file *FileModel) Create(ownerId, teamId, folderId primitive.ObjectID, fileName, address string,
	expireAt time.Time) (primitive.ObjectID, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	newFile := &File{
		OwnerId:   ownerId,
		TeamId:    teamId,
		FolderId:  folderId,
		Name:      fileName,
		Address:   address,
		ExpireAt:  expireAt,
		CreatedAt: time.Now(),
	}

	id, err := file.db.Collection("files").InsertOne(ctx, newFile)
	if err != nil {
		return primitive.NilObjectID, err
	}

	return id.InsertedID.(primitive.ObjectID), nil
}

// GetAll -> Returns List
func (file *FileModel) GetAll(filter bson.M, page, pageSize int64) ([]File, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	projection := bson.M{
		"hashed_password": 0,
		"salt":            0,
	}

	findOptions := options.Find()
	findOptions.SetProjection(projection)
	findOptions.SetSkip((page - 1) * pageSize)
	findOptions.SetLimit(pageSize)

	cursor, err := file.db.Collection("files").Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}

	var files []File
	if err := cursor.All(ctx, &files); err != nil {
		return nil, err
	}

	return files, nil
}

func (file *FileModel) Delete(id primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{
		"_id": id,
	}

	deletedCount, err := file.db.Collection("files").DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	if deletedCount.DeletedCount == 0 {
		return errors.New("file with this id does not exist")
	}

	return nil
}

func (file *FileModel) Update(id primitive.ObjectID, updates bson.M) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	update := bson.M{
		"$set": updates,
	}

	updateResult, err := file.db.Collection("files").UpdateByID(ctx, id, update)
	if err != nil {
		return err
	}

	if updateResult.MatchedCount == 0 {
		return errors.New("file not found")
	}

	if updateResult.ModifiedCount == 0 {
		return errors.New("no change detected. Please enter a different name")
	}

	return nil
}

// Get -> Returns One
func (file *FileModel) Get(filter, projection bson.M) (*File, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	findOptions := options.FindOne()
	findOptions.SetProjection(projection)

	var fileInstance File
	if err := file.db.Collection("files").FindOne(ctx, filter, findOptions).Decode(&fileInstance); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("file does not exist")
		}

		return nil, err
	}

	return &fileInstance, nil
}
