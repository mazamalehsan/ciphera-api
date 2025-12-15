package userModel

import (
	"ciphera-api/constants"
	"ciphera-api/db"
	"ciphera-api/types"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func InsertUser(user *types.User) (*mongo.InsertOneResult, error) {
	user.ID = primitive.NewObjectID().Hex()
	collection := db.Collection(constants.UserCollectionName)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	result, err := collection.InsertOne(ctx, user)

	if err != nil {
		return nil, err
	}

	return result, nil
}
