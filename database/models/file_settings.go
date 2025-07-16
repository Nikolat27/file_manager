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

type FileSettingModel struct {
	db *mongo.Database
}

type FileSettings struct {
	Id             primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	FileId         primitive.ObjectID `json:"file_id" bson:"file_id"`
	Salt           string             `json:"salt" bson:"salt"`
	HashedPassword string             `json:"hashed_password" bson:"hashed_password"`
	MaxDownloads   uint64             `json:"max_downloads" bson:"max_downloads"`
	ViewOnly       bool               `json:"view_only" bson:"view_only"`
	Approvable     bool               `json:"approvable" bson:"approvable"`
}

const FileSettingsCollectionName = "file_settings"

func (file *FileSettingModel) Create(fileId primitive.ObjectID, salt, hashedPassword string, maxDownloads uint64,
	viewOnly, approvable bool) error {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	newFile := &FileSettings{
		FileId:         fileId,
		Salt:           salt,
		HashedPassword: hashedPassword,
		MaxDownloads:   maxDownloads,
		ViewOnly:       viewOnly,
		Approvable:     approvable,
	}

	if _, err := file.db.Collection(FileSettingsCollectionName).InsertOne(ctx, newFile); err != nil {
		return err
	}

	return nil
}

func (file *FileSettingModel) IsApprovalRequired(fileId primitive.ObjectID) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{
		"file_id": fileId,
	}

	projection := bson.M{
		"approvable": 1,
	}

	fileOptions := options.FindOne()
	fileOptions.SetProjection(projection)

	var fileInstance FileSettings
	if err := file.db.Collection(FileSettingsCollectionName).FindOne(ctx, filter, fileOptions).Decode(&fileInstance); err != nil {
		return false, err
	}

	if fileInstance.Approvable {
		return true, nil
	}

	return false, nil
}

func (file *FileSettingModel) IsPasswordRequired(fileId primitive.ObjectID) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{
		"file_id": fileId,
	}

	var fileInstance FileSettings
	if err := file.db.Collection(FileSettingsCollectionName).FindOne(ctx, filter).Decode(&fileInstance); err != nil {
		return false, err
	}

	if fileInstance.HashedPassword == "" {
		return false, nil
	}

	return true, nil
}

func (file *FileSettingModel) GetOne(fileId primitive.ObjectID) (*FileSettings, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{
		"file_id": fileId,
	}

	var fileInstance FileSettings
	if err := file.db.Collection(FileSettingsCollectionName).FindOne(ctx, filter).Decode(&fileInstance); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("file with this Id does not exist")
		}

		return nil, err
	}

	return &fileInstance, nil
}
