package userService

import (
	"ciphera-api/db/models/userModel"
	"ciphera-api/types"
	"context"
	"errors"
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

func (s *Service) UsernameExists(
	ctx context.Context,
	username string,
) (bool, error) {
	return s.userRepo.UsernameExists(ctx, username)
}

func (s *Service) RegisterUser(
	ctx context.Context,
	user types.User,
) error {
	exists, err := s.UsernameExists(ctx, user.Username)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("username already taken")
	}
	return s.userRepo.Insert(ctx, &user)
}
