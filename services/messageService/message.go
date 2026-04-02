package messageService

import (
	"ciphera-api/db/models/messageModel"
	"ciphera-api/db/models/userModel"
	"ciphera-api/types"
	"context"
	"time"

	"github.com/google/uuid"
)

type Service struct {
	messageRepo messageModel.MessageRepository
	userRepo    userModel.UserRepository
}

func New(messageRepo messageModel.MessageRepository, userRepo userModel.UserRepository) *Service {
	return &Service{
		messageRepo: messageRepo,
		userRepo:    userRepo,
	}
}

func (s *Service) SendMessage(ctx context.Context, senderUserId string, msg types.Message) (types.Message, error) {
	msg.UUID = uuid.NewString()
	msg.From = senderUserId
	msg.CreatedAt = time.Now()
	msg.UpdatedAt = time.Now()

	// Attach sender's public key so recipients can derive shared secret
	sender, err := s.userRepo.FetchByID(ctx, senderUserId)
	if err != nil {
		return msg, err
	}
	msg.SenderPublicKey = sender.PublicKey
	msg.SenderEncryptionKey = sender.EncryptionKey
	msg.FromUsername = sender.Username

	err = s.messageRepo.Insert(ctx, &msg)
	return msg, err
}

func (s *Service) GetConversation(ctx context.Context, userA, userB string) ([]types.Message, error) {
	return s.messageRepo.FetchConversation(ctx, userA, userB)
}

func (s *Service) GetGroupMessages(ctx context.Context, groupId string) ([]types.Message, error) {
	return s.messageRepo.FetchGroupMessages(ctx, groupId)
}
