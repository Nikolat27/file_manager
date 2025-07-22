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

type ApprovalModel struct {
	db *mongo.Database
}

type Approval struct {
	Id         primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	FileId     primitive.ObjectID `json:"file_id" bson:"file_id"`
	FileName   string             `json:"file_name" bson:"file_name"`
	OwnerId    primitive.ObjectID `json:"owner_id" bson:"owner_id"`   // the file owner id
	SenderId   primitive.ObjectID `json:"sender_id" bson:"sender_id"` // the Requester id (user-id)
	Status     string             `json:"status" bson:"status"`       // pending, approved, rejected
	Reason     string             `json:"reason" bson:"reason"`
	CreatedAt  time.Time          `json:"created_at" bson:"created_at"`
	ReviewedAt *time.Time         `json:"reviewed_at,omitempty" bson:"reviewed_at,omitempty"`
}

func (approval *ApprovalModel) Create(fileId, ownerId, senderId primitive.ObjectID, fileName, reason string) (primitive.ObjectID, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var newApproval = &Approval{
		FileId:    fileId,
		OwnerId:   ownerId,
		SenderId:  senderId,
		Status:    "pending", // default
		FileName:  fileName,
		Reason:    reason,
		CreatedAt: time.Now(),
	}

	id, err := approval.db.Collection("approvals").InsertOne(ctx, newApproval)
	if err != nil {
		return primitive.NilObjectID, err
	}

	return id.InsertedID.(primitive.ObjectID), nil
}

// GetAll -> Returns All
func (approval *ApprovalModel) GetAll(filter, projection bson.M, page, pageLimit int64) ([]Approval, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	findOptions := options.Find()
	findOptions.SetProjection(projection)
	findOptions.SetSkip((page - 1) * pageLimit)
	findOptions.SetLimit(pageLimit)

	var approvals []Approval
	cursor, err := approval.db.Collection("approvals").Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}

	if err := cursor.All(ctx, &approvals); err != nil {
		return nil, err
	}

	return approvals, nil
}

// Get -> Returns One
func (approval *ApprovalModel) Get(filter, projection bson.M) (*Approval, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	findOptions := options.FindOne()
	findOptions.SetProjection(projection)

	var approvalInstance Approval
	if err := approval.db.Collection("approvals").FindOne(ctx, filter, findOptions).Decode(&approvalInstance); err != nil {
		return nil, err
	}

	return &approvalInstance, nil
}

func (approval *ApprovalModel) Update(id primitive.ObjectID, updates bson.M) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	update := bson.M{
		"$set": updates,
	}

	result, err := approval.db.Collection("approvals").UpdateByID(ctx, id, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return errors.New("this approval does not even exist")
	}

	if result.ModifiedCount == 0 {
		return errors.New("did not detect any changes")
	}

	return nil
}

func (approval *ApprovalModel) DeleteOne(filter bson.M) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := approval.db.Collection("approvals").DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return errors.New("nothing deleted")
	}

	return nil
}
