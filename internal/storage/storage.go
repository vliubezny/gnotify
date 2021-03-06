package storage

import (
	"context"
	"errors"

	"github.com/vliubezny/gnotify/internal/model"
)

//go:generate mockgen -destination=./mock/mock.go -package=mock -source=storage.go

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

	// AddDevice add new device for user
	AddDevice(ctx context.Context, userID int64, input model.Device) (model.Device, error)
}
