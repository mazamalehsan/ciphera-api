package loginChallengeModel

import (
	"ciphera-api/types"
	"context"
)

type LoginChallengeRepository interface {
	Insert(ctx context.Context, user *types.UserLoginChallenge) error
	FetchByUUIDAndUserUUID(ctx context.Context, uuid, userUuid string) (types.UserLoginChallenge, error)
}
