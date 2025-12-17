package loginChallengeModel

import (
	"ciphera-api/constants"
	"ciphera-api/db"
	"ciphera-api/types"
	"context"

	"go.mongodb.org/mongo-driver/bson"
)

type MongoLoginChallengeRepository struct{}

func New() LoginChallengeRepository {
	return &MongoLoginChallengeRepository{}
}

func (m *MongoLoginChallengeRepository) Insert(ctx context.Context, loginChallenge *types.UserLoginChallenge) error {
	collection := db.Collection(constants.LoginChallengeCollectionName)
	_, err := collection.InsertOne(ctx, loginChallenge)
	return err
}

func (m *MongoLoginChallengeRepository) FetchByUUIDAndUserUUID(ctx context.Context, uuid, userUUID string) (types.UserLoginChallenge, error) {
	var loginChallenge types.UserLoginChallenge
	err := db.Collection(constants.LoginChallengeCollectionName).
		FindOne(ctx, bson.M{"uuid": uuid, "useruuid": userUUID}).
		Decode(&loginChallenge)

	return loginChallenge, err
}
