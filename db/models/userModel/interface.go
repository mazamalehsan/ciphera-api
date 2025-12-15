package userModel

import (
	"ciphera-api/types"
	"context"
)

type UserRepository interface {
	Insert(ctx context.Context, user *types.User) error
	FetchByID(ctx context.Context, id string) (types.User, error)
	FetchByUsername(ctx context.Context, username string) (types.User, error)
	UsernameExists(ctx context.Context, username string) (bool, error)
	FetchPublicKey(ctx context.Context, username string) (string, error)
}
