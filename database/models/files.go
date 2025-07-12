package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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
	Salt           string             `json:"salt" bson:"salt"`
	HashedPassword string             `json:"hashed_password" bson:"hashed_password"`
	MaxDownloads   uint               `json:"max_downloads" bson:"max_downloads"`
	CreatedAt      time.Time          `json:"created_at" bson:"created_at"`
	ExpireAt       time.Time          `json:"expire_at" bson:"expire_at"`
}

func (file *FileModel) CreateFileInstance(ownerId primitive.ObjectID, fileName, address string, salt, hashedPassword []byte, approvable bool, maxDownloads uint64, expireAt time.Time) {
	
}
