package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/vliubezny/gnotify/internal/model"
	"github.com/vliubezny/gnotify/internal/storage"
)

//go:generate mockgen -destination=./mock/mock.go -package=mock -source=service.go

var (
	// ErrNotFound states that record was not found in storage.
	ErrNotFound = errors.New("not found")
)

type Service interface {
	// GetUser returns user by ID.
	GetUser(ctx context.Context, id int64) (model.User, error)

	// GetUsers returns list of users.
	GetUsers(ctx context.Context) ([]model.User, error)

	// UpsertUser inserts or updates user setting if record exists.
	UpsertUser(ctx context.Context, user model.User) error

	// DeleteUser deletes user by ID.
	DeleteUser(ctx context.Context, id int64) error
}

type service struct {
	s storage.Storage
}

func New(s storage.Storage) Service {
	return &service{
		s: s,
	}
}

func (s *service) GetUser(ctx context.Context, id int64) (model.User, error) {
	user, err := s.s.GetUser(ctx, id)
	if err != nil {
		if err == storage.ErrNotFound {
			return model.User{}, ErrNotFound
		}
		return model.User{}, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
}

func (s *service) GetUsers(ctx context.Context) ([]model.User, error) {
	users, err := s.s.GetUsers(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}
	return users, nil
}

func (s *service) UpsertUser(ctx context.Context, user model.User) error {
	if err := s.s.UpsertUser(ctx, user); err != nil {
		return fmt.Errorf("failed to upsert user: %w", err)
	}
	return nil
}

func (s *service) DeleteUser(ctx context.Context, id int64) error {
	if err := s.s.DeleteUser(ctx, id); err != nil {
		if err == storage.ErrNotFound {
			return ErrNotFound
		}
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}
