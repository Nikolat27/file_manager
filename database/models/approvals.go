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
	DB *mongo.Database
}

type Approval struct {
	Id         primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	FileId     primitive.ObjectID `json:"file_id" bson:"file_id"`
	OwnerId    primitive.ObjectID `json:"owner_id" bson:"owner_id"`
	SenderId   primitive.ObjectID `json:"sender_id" bson:"sender_id"` // the Requester id (user-id)
	Status     string             `json:"status" bson:"status"`       // pending, approved, rejected
	Reason     string             `json:"reason" bson:"reason"`
	CreatedAt  time.Time          `json:"created_at" bson:"created_at"`
	ReviewedAt *time.Time         `json:"reviewed_at,omitempty" bson:"reviewed_at,omitempty"`
}

const ApprovalCollectionName = "approvals"

func (approval *ApprovalModel) CreateApprovalInstance(fileId, ownerId, senderId primitive.ObjectID, reason string) (primitive.ObjectID, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var newApproval = &Approval{
		FileId:    fileId,
		OwnerId:   ownerId,
		SenderId:  senderId,
		Status:    "pending", // default
		Reason:    reason,
		CreatedAt: time.Now(),
	}

	id, err := approval.DB.Collection(ApprovalCollectionName).InsertOne(ctx, newApproval)
	if err != nil {
		return primitive.NilObjectID, err
	}

	return id.InsertedID.(primitive.ObjectID), nil
}

func (approval *ApprovalModel) UpdateApprovalStatus(id primitive.ObjectID, status string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if status != "approved" && status != "rejected" && status != "pending" {
		return errors.New("invalid status parameters. Must be: approved, rejected, pending")
	}

	update := bson.M{
		"$set": bson.M{
			"status":      status,
			"reviewed_at": time.Now(),
		},
	}

	result, err := approval.DB.Collection(ApprovalCollectionName).UpdateByID(ctx, id, update)
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

func (approval *ApprovalModel) CheckUserApprovalStatus(fileId, senderId primitive.ObjectID) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{
		"file_id":   fileId,
		"sender_id": senderId,
	}

	projection := bson.M{
		"status": 1,
	}

	approvalOptions := options.FindOne()
	approvalOptions.SetProjection(projection)

	var approvalInstance Approval
	if err := approval.DB.Collection(ApprovalCollectionName).FindOne(ctx, filter, approvalOptions).Decode(&approvalInstance); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return "", nil
		}

		return "", err
	}

	return approvalInstance.Status, nil
}
