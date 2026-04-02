package groupService

import (
	"ciphera-api/db/models/groupModel"
	"ciphera-api/db/models/userModel"
	"ciphera-api/types"
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

type Service struct {
	groupRepo groupModel.GroupRepository
	userRepo  userModel.UserRepository
}

func New(groupRepo groupModel.GroupRepository, userRepo userModel.UserRepository) *Service {
	return &Service{
		groupRepo: groupRepo,
		userRepo:  userRepo,
	}
}

func (s *Service) CreateGroup(ctx context.Context, creatorId string, name string, memberUsernames []string) (types.Group, error) {
	group := types.Group{
		UUID:      uuid.NewString(),
		Name:      name,
		Admins:    []string{creatorId},
		Users:     []string{creatorId},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Resolve usernames to UUIDs and add them
	for _, username := range memberUsernames {
		user, err := s.userRepo.FetchByUsername(ctx, username)
		if err != nil {
			continue // skip invalid usernames
		}
		group.Users = append(group.Users, user.UUID)
	}

	err := s.groupRepo.Insert(ctx, &group)
	return group, err
}

func (s *Service) GetMyGroups(ctx context.Context, userId string) ([]types.Group, error) {
	return s.groupRepo.FetchByUser(ctx, userId)
}

func (s *Service) GetGroupMembers(ctx context.Context, groupId string) ([]types.User, error) {
	group, err := s.groupRepo.FetchByUUID(ctx, groupId)
	if err != nil {
		return nil, err
	}

	var members []types.User
	for _, uid := range group.Users {
		user, err := s.userRepo.FetchByID(ctx, uid)
		if err != nil {
			continue
		}
		members = append(members, user)
	}
	return members, nil
}

func (s *Service) AddMember(ctx context.Context, requesterId, groupId, username string) error {
	group, err := s.groupRepo.FetchByUUID(ctx, groupId)
	if err != nil {
		return err
	}

	// Only admins can add members
	isAdmin := false
	for _, admin := range group.Admins {
		if admin == requesterId {
			isAdmin = true
			break
		}
	}
	if !isAdmin {
		return errors.New("only admins can add members")
	}

	user, err := s.userRepo.FetchByUsername(ctx, username)
	if err != nil {
		return errors.New("user not found")
	}

	return s.groupRepo.AddUser(ctx, groupId, user.UUID)
}
