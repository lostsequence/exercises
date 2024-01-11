package services

import (
	"context"
	"errors"
	"fmt"
	"movies-auth/users/internal/domain"
)

type UsersStorage interface {
	GetUserByID(ctx context.Context, login string) (domain.User, error)
	Insert(ctx context.Context, user domain.User) (domain.User, error)
	IsUserExist(ctx context.Context, login string) (bool, error)
}

type UsersService struct {
	Storage UsersStorage
}

func NewUsersService(storage UsersStorage) *UsersService {
	return &UsersService{
		Storage: storage,
	}
}

func (s *UsersService) Create(ctx context.Context, user domain.User) (domain.User, error) {
	isExist, err := s.Storage.IsUserExist(ctx, user.Login)
	if err != nil && !errors.Is(err, domain.ErrNotFound) {
		return domain.User{}, fmt.Errorf("failed to get user from storage: %w", err)
	}
	if isExist {
		return domain.User{}, fmt.Errorf("user create error: %w", domain.ErrConflict)
	}

	createdUser, err := s.Storage.Insert(ctx, user)
	if err != nil {
		return domain.User{}, err
	}

	return createdUser, nil
}

func (s *UsersService) Login(ctx context.Context, user domain.User) (domain.User, error) {
	existingUser, err := s.Storage.GetUserByID(ctx, user.Login)
	if err != nil {
		return domain.User{}, fmt.Errorf("failed to get user from storage: %w", err)
	}

	if existingUser.Password != user.Password {
		return domain.User{}, fmt.Errorf("password incorrect: %w", err)
	}

	return existingUser, nil
}
