package userModel

import (
	"context"

	"ciphera-api/constants"
	"ciphera-api/db"
	"ciphera-api/types"

	"go.mongodb.org/mongo-driver/bson"
)

type MongoUserRepository struct{}

func New() UserRepository {
	return &MongoUserRepository{}
}

func (m *MongoUserRepository) Insert(ctx context.Context, user *types.User) error {
	collection := db.Collection(constants.UserCollectionName)
	_, err := collection.InsertOne(ctx, user)
	return err
}

func (m *MongoUserRepository) FetchByID(ctx context.Context, id string) (types.User, error) {
	var user types.User
	err := db.Collection(constants.UserCollectionName).
		FindOne(ctx, bson.M{"id": id}).
		Decode(&user)

	return user, err
}

func (m *MongoUserRepository) FetchByUsername(ctx context.Context, username string) (types.User, error) {
	var user types.User
	err := db.Collection(constants.UserCollectionName).
		FindOne(ctx, bson.M{"username": username}).
		Decode(&user)

	return user, err
}

func (m *MongoUserRepository) UsernameExists(ctx context.Context, username string) (bool, error) {
	count, err := db.Collection(constants.UserCollectionName).
		CountDocuments(ctx, bson.M{"username": username})

	return count > 0, err
}

func (m *MongoUserRepository) FetchPublicKey(ctx context.Context, username string) (string, error) {
	user, err := m.FetchByUsername(ctx, username)
	if err != nil {
		return "", err
	}
	return user.PublicKey, nil
}
