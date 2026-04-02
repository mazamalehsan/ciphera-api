package authService

import (
	"ciphera-api/db/models/loginChallengeModel"
	"ciphera-api/db/models/userModel"
	"ciphera-api/encryption"
	"ciphera-api/env"
	"ciphera-api/jwt"
	"ciphera-api/types"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Service struct {
	userRepo           userModel.UserRepository
	loginChallengeRepo loginChallengeModel.LoginChallengeRepository
}

func New(userRepo userModel.UserRepository, loginChallengeRepo loginChallengeModel.LoginChallengeRepository) *Service {
	return &Service{
		userRepo:           userRepo,
		loginChallengeRepo: loginChallengeRepo,
	}
}

func (s *Service) CheckUsernameDuplication(ctx context.Context, username string) (bool, error) {
	return s.userRepo.UsernameExists(ctx, username)
}

func (s *Service) RegisterUser(
	ctx context.Context,
	user types.User,
) error {
	exists, err := s.userRepo.UsernameExists(ctx, user.Username)

	if err != nil {
		return err
	}
	if exists {
		return errors.New("username already taken")
	}

	user.UUID = uuid.NewString()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	return s.userRepo.Insert(ctx, &user)
}

func (s *Service) LoginUser(
	ctx context.Context,
	body types.LoginBody,
) (string, error) {
	user, err := s.userRepo.FetchByUsername(ctx, body.Username)
	if err != nil {
		return "", err
	}

	challenge, err := s.loginChallengeRepo.FetchByUUIDAndUserUUID(ctx, body.LoginId, user.UUID)

	if err != nil || challenge.UUID == "" {
		return "", errors.New("challenge cannot be verified")
	}

	if time.Now().After(challenge.ExpiresAt) || challenge.Used {
		return "", errors.New("challenge expired")
	}

	err = encryption.VerifySignature(user.PublicKey, challenge.Challenge, body.SignedChallenge)

	if err != nil {
		fmt.Print(err)
		return "", errors.New("challenge cannot be verified")
	}

	privateKey, err := jwt.LoadRSAPrivateKey(env.GetJWTPrivateKey())

	if err != nil {
		return "", err
	}

	token, err := jwt.GenerateToken(privateKey, user.UUID, time.Now().Add(15*time.Minute))

	if err != nil {
		return "", err
	}

	return token, nil

}

func (s *Service) GenerateLoginChallenge(
	ctx context.Context,
	body types.LoginBody,
) (string, string, error) {
	user, err := s.userRepo.FetchByUsername(ctx, body.Username)
	if err != nil || user.Username == "" {
		return "", "", errors.New("access denied")
	}
	challenge, err := encryption.GenerateChallenge()

	var challengeRecord types.UserLoginChallenge
	challengeRecord.Challenge = challenge
	loginId := uuid.NewString()
	challengeRecord.UUID = loginId
	challengeRecord.UserUUID = user.UUID
	challengeRecord.ExpiresAt = time.Now().Add(1 * time.Minute)
	challengeRecord.Used = false
	err = s.loginChallengeRepo.Insert(ctx, &challengeRecord)
	return challenge, loginId, err
}
