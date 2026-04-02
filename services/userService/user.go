package userService

import (
	"ciphera-api/db/models/userModel"
	"ciphera-api/types"
	"context"
)

type Service struct {
	userRepo userModel.UserRepository
}

func New(userRepo userModel.UserRepository) *Service {
	return &Service{
		userRepo: userRepo,
	}
}

func (s *Service) GetByUsername(
	ctx context.Context,
	username string,
) (types.User, error) {
	return s.userRepo.FetchByUsername(ctx, username)
}
