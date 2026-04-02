package messageModel

import (
	"ciphera-api/constants"
	"ciphera-api/db"
	"ciphera-api/types"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoMessageRepository struct{}

func New() MessageRepository {
	return &MongoMessageRepository{}
}

func (m *MongoMessageRepository) Insert(ctx context.Context, msg *types.Message) error {
	_, err := db.Collection(constants.MessageCollectionName).InsertOne(ctx, msg)
	return err
}

func (m *MongoMessageRepository) FetchConversation(ctx context.Context, userA, userB string) ([]types.Message, error) {
	filter := bson.M{
		"$or": []bson.M{
			{"from": userA, "to": userB},
			{"from": userB, "to": userA},
		},
		"groupid": bson.M{"$in": []interface{}{"", nil}},
	}

	opts := options.Find().SetSort(bson.D{{Key: "createdat", Value: 1}})

	cursor, err := db.Collection(constants.MessageCollectionName).Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var messages []types.Message
	if err = cursor.All(ctx, &messages); err != nil {
		return nil, err
	}
	return messages, nil
}

func (m *MongoMessageRepository) FetchGroupMessages(ctx context.Context, groupId string) ([]types.Message, error) {
	filter := bson.M{"groupid": groupId}
	opts := options.Find().SetSort(bson.D{{Key: "createdat", Value: 1}})

	cursor, err := db.Collection(constants.MessageCollectionName).Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var messages []types.Message
	if err = cursor.All(ctx, &messages); err != nil {
		return nil, err
	}
	return messages, nil
}
