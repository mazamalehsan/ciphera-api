package messageModel

import (
	"ciphera-api/types"
	"context"
)

type MessageRepository interface {
	Insert(ctx context.Context, msg *types.Message) error
	FetchConversation(ctx context.Context, userA, userB string) ([]types.Message, error)
	FetchGroupMessages(ctx context.Context, groupId string) ([]types.Message, error)
}
