package storage

import (
	"context"
	"errors"

	"github.com/vliubezny/gnotify/internal/model"
)

var (
	// ErrNotFound states that record was not found in storage.
	ErrNotFound = errors.New("not found")
)

// Storage saves and loads user notification settings.
type Storage interface {
	// GetUser returns user by ID.
	GetUser(ctx context.Context, id int64) (model.User, error)

	// GetUsers returns list of users.
	GetUsers(ctx context.Context) ([]model.User, error)

	// UpsertUser inserts or updates user setting if record exists.
	UpsertUser(ctx context.Context, user model.User) error

	// DeleteUser deletes user by ID.
	DeleteUser(ctx context.Context, id int64) error
}
