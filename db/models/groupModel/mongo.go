package groupModel

import (
	"ciphera-api/constants"
	"ciphera-api/db"
	"ciphera-api/types"
	"context"

	"go.mongodb.org/mongo-driver/bson"
)

type MongoGroupRepository struct{}

func New() GroupRepository {
	return &MongoGroupRepository{}
}

func (m *MongoGroupRepository) Insert(ctx context.Context, group *types.Group) error {
	_, err := db.Collection(constants.GroupsCollectionName).InsertOne(ctx, group)
	return err
}

func (m *MongoGroupRepository) FetchByUUID(ctx context.Context, uuid string) (types.Group, error) {
	var group types.Group
	err := db.Collection(constants.GroupsCollectionName).
		FindOne(ctx, bson.M{"uuid": uuid}).
		Decode(&group)
	return group, err
}

func (m *MongoGroupRepository) FetchByUser(ctx context.Context, userId string) ([]types.Group, error) {
	filter := bson.M{
		"$or": []bson.M{
			{"users": userId},
			{"admins": userId},
		},
	}

	cursor, err := db.Collection(constants.GroupsCollectionName).Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var groups []types.Group
	if err = cursor.All(ctx, &groups); err != nil {
		return nil, err
	}
	return groups, nil
}

func (m *MongoGroupRepository) AddUser(ctx context.Context, groupId, userId string) error {
	_, err := db.Collection(constants.GroupsCollectionName).UpdateOne(
		ctx,
		bson.M{"uuid": groupId},
		bson.M{"$addToSet": bson.M{"users": userId}},
	)
	return err
}
