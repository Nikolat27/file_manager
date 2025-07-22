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
	Id                    primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserId                primitive.ObjectID `json:"user_id" bson:"user_id"`
	FileId                primitive.ObjectID `json:"file_id" bson:"file_id"` // 1 to 1 relationship
	ShortUrl              string             `json:"short_url" bson:"short_url"`
	Salt                  string             `json:"salt" bson:"salt"`
	HashedPassword        string             `json:"hashed_password" bson:"hashed_password"`
	MaxDownloads          int64              `json:"max_downloads" bson:"max_downloads"`
	CurrentDownloadAmount int64              `json:"current_download_amount" bson:"current_download_amount"`
	ViewOnly              bool               `json:"view_only" bson:"view_only"`
	Approvable            bool               `json:"approvable" bson:"approvable"`
	ExpireAt              time.Time          `json:"expiration_at" bson:"expiration_at"`
	CreatedAt             time.Time          `json:"created_at" bson:"created_at"`
}

const FileSettingsCollectionName = "file_settings"

func (file *FileSettingModel) Create(fileId, userId primitive.ObjectID, shortUrl, salt, hashedPassword string, maxDownloads int64,
	viewOnly, approvable bool, expireAt time.Time) error {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	newFile := &FileSettings{
		FileId:         fileId,
		UserId:         userId,
		ShortUrl:       shortUrl,
		Salt:           salt,
		HashedPassword: hashedPassword,
		MaxDownloads:   maxDownloads,
		ViewOnly:       viewOnly,
		Approvable:     approvable,
		ExpireAt:       expireAt,
		CreatedAt:      time.Now(),
	}

	if _, err := file.db.Collection(FileSettingsCollectionName).InsertOne(ctx, newFile); err != nil {
		return err
	}

	return nil
}

// Get -> Returns one
func (file *FileSettingModel) Get(filter, projection bson.M) (*FileSettings, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	findOptions := options.FindOne()
	findOptions.SetProjection(projection)

	var fileInstance FileSettings
	if err := file.db.Collection(FileSettingsCollectionName).FindOne(ctx, filter, findOptions).Decode(&fileInstance); err != nil {
		return nil, err
	}

	return &fileInstance, nil
}

func (file *FileSettingModel) GetAll(filter bson.M) ([]FileSettings, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var settings []FileSettings
	cursor, err := file.db.Collection(FileSettingsCollectionName).Find(ctx, filter)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
	}

	if err := cursor.All(ctx, &settings); err != nil {
		return nil, err
	}

	return settings, nil
}

func (file *FileSettingModel) Delete(filter bson.M) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := file.db.Collection("file_settings").DeleteOne(ctx, filter); err != nil {
		return err
	}

	return nil
}

func (file *FileSettingModel) Update(id primitive.ObjectID, updates bson.M) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	update := bson.M{
		"$set": updates,
	}
	
	result, err := file.db.Collection(FileSettingsCollectionName).UpdateByID(ctx, id, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return errors.New("setting with this id does not exist")
	}

	if result.ModifiedCount == 0 {
		return errors.New("no change detected")
	}

	return nil
}
