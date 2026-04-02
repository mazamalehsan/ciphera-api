package groupModel

import (
	"ciphera-api/types"
	"context"
)

type GroupRepository interface {
	Insert(ctx context.Context, group *types.Group) error
	FetchByUUID(ctx context.Context, uuid string) (types.Group, error)
	FetchByUser(ctx context.Context, userId string) ([]types.Group, error)
	AddUser(ctx context.Context, groupId, userId string) error
}
