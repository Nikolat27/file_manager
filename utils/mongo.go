package utils

import "go.mongodb.org/mongo-driver/bson/primitive"

func ConvertStringToObjectID(s string) (primitive.ObjectID, error) {
	objId, err := primitive.ObjectIDFromHex(s)
	if err != nil {
		return primitive.NilObjectID, err
	}

	return objId, err
}
