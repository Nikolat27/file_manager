package utils

import (
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ConvertStringToObjectID(s string) (primitive.ObjectID, error) {
	if s == "" {
		return primitive.NilObjectID, errors.New("given string is empty")
	}

	objId, err := primitive.ObjectIDFromHex(s)
	if err != nil {
		return primitive.NilObjectID, err
	}

	return objId, err
}
