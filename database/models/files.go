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
	DB *mongo.Database
}

type File struct {
	Id             primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	OwnerId        primitive.ObjectID `json:"owner_id" bson:"owner_id"`
	Name           string             `json:"name" bson:"name"`
	Address        string             `json:"address" bson:"address"`
	Approvable     bool               `json:"approvable" bson:"approvable"`
	ShortUrl       string             `json:"short_url" bson:"short_url"`
	Salt           string             `json:"salt" bson:"salt"`
	HashedPassword string             `json:"hashed_password" bson:"hashed_password"`
	MaxDownloads   uint64             `json:"max_downloads" bson:"max_downloads"`
	CreatedAt      time.Time          `json:"created_at" bson:"created_at"`
	ExpireAt       time.Time          `json:"expire_at" bson:"expire_at"`
}

const FileCollectionName = "files"

func (file *FileModel) CreateFileInstance(ownerId primitive.ObjectID,
	fileName, address, shortUrl, salt, hashedPassword string, approvable bool,
	maxDownloads uint64, expireAt time.Time) (primitive.ObjectID, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	newFile := &File{
		OwnerId:        ownerId,
		Name:           fileName,
		Address:        address,
		ShortUrl:       shortUrl,
		Salt:           salt,
		HashedPassword: hashedPassword,
		Approvable:     approvable,
		MaxDownloads:   maxDownloads,
		CreatedAt:      time.Now(),
		ExpireAt:       expireAt,
	}

	id, err := file.DB.Collection(FileCollectionName).InsertOne(ctx, newFile)
	if err != nil {
		return primitive.NilObjectID, err
	}

	return id.InsertedID.(primitive.ObjectID), nil
}

// GetFilesInstances -> Returns List
func (file *FileModel) GetFilesInstances(ownerId primitive.ObjectID, page, pageSize int64) ([]File, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{
		"owner_id": ownerId,
	}

	projection := bson.M{
		"hashed_password": 0,
		"salt":            0,
	}

	findOptions := options.Find()
	findOptions.SetProjection(projection)
	findOptions.SetSkip((page - 1) * pageSize)
	findOptions.SetLimit(pageSize)

	cursor, err := file.DB.Collection(FileCollectionName).Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}

	var files []File
	if err := cursor.All(ctx, &files); err != nil {
		return nil, err
	}

	return files, nil
}

// GetFileInstance -> Returns One
func (file *FileModel) GetFileInstance(shortUrl []byte) (*File, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{
		"short_url": string(shortUrl),
	}

	var fileInstance File
	if err := file.DB.Collection(FileCollectionName).FindOne(ctx, filter).Decode(&fileInstance); err != nil {
		return nil, err
	}

	return &fileInstance, nil
}

func (file *FileModel) DeleteFileInstance(id primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{
		"_id": id,
	}

	deletedCount, err := file.DB.Collection(FileCollectionName).DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	if deletedCount.DeletedCount == 0 {
		return errors.New("file with this id does not exist")
	}

	return nil
}

func (file *FileModel) GetFileAddress(id primitive.ObjectID) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{
		"_id": id,
	}

	projection := bson.M{
		"address": 1,
	}

	fileOptions := options.FindOne()
	fileOptions.SetProjection(projection)

	var fileInstance File
	if err := file.DB.Collection(FileCollectionName).FindOne(ctx, filter, fileOptions).Decode(&fileInstance); err != nil {
		return nil, err
	}

	return []byte(fileInstance.Address), nil
}

func (file *FileModel) RenameFileInstance(id primitive.ObjectID, newName []byte) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{
		"_id": id,
	}

	update := bson.M{
		"$set": bson.M{
			"name": string(newName),
		},
	}

	updateResult, err := file.DB.Collection(FileCollectionName).UpdateOne(ctx, filter, update)
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

func (file *FileModel) IsFileExpired(id primitive.ObjectID) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{
		"_id": id,
	}

	projection := bson.M{
		"expire_at": 1,
	}

	fileOptions := options.FindOne()
	fileOptions.SetProjection(projection)

	var fileInstance File
	if err := file.DB.Collection(FileCollectionName).FindOne(ctx, filter, fileOptions).Decode(&fileInstance); err != nil {
		return false, err
	}

	if fileInstance.ExpireAt.IsZero() {
		return false, nil
	}

	if time.Now().After(fileInstance.ExpireAt) {
		return true, nil
	}

	return false, nil
}

func (file *FileModel) RequirePassword(shortUrl []byte) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{
		"short_url": string(shortUrl),
	}

	var fileInstance File
	if err := file.DB.Collection(FileCollectionName).FindOne(ctx, filter).Decode(&fileInstance); err != nil {
		return false, err
	}
	
	if fileInstance.HashedPassword == "" {
		return false, nil
	}
	
	return true, nil
}

func (file *FileModel) GetFileOwnerId(id primitive.ObjectID) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{
		"_id": id,
	}

	projection := bson.M{
		"owner_id": 1,
	}

	fileOptions := options.FindOne()
	fileOptions.SetProjection(projection)

	var fileInstance File
	if err := file.DB.Collection(FileCollectionName).FindOne(ctx, filter, fileOptions).Decode(&fileInstance); err != nil {
		return nil, err
	}

	return []byte(fileInstance.Address), nil
}

func (file *FileModel) CheckFileRequiresApproval(id primitive.ObjectID) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{
		"_id": id,
	}

	projection := bson.M{
		"approvable": 1,
	}

	fileOptions := options.FindOne()
	fileOptions.SetProjection(projection)

	var fileInstance File
	if err := file.DB.Collection(FileCollectionName).FindOne(ctx, filter, fileOptions).Decode(&fileInstance); err != nil {
		return false, err
	}

	if fileInstance.Approvable {
		return true, nil
	}
	
	return false, nil
}
