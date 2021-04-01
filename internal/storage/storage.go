package storage

import (
	"context"

	"github.com/vliubezny/gnotify/internal/model"
)

// Storage saves and loads user notification settings.
type Storage interface {
	// GetUser returns user by ID.
	GetUser(ctx context.Context, id int64) (model.User, error)
}
