package mongodb

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"github.com/vliubezny/gnotify/internal/model"
	"github.com/vliubezny/gnotify/internal/storage"
)

const (
	db    = "gnotify"
	users = "users"
)

type mongoStorage struct {
	client *mongo.Client
}

// New creates mongodb storage.
func New(uri string) (storage.Storage, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, err
	}

	return &mongoStorage{
		client: client,
	}, nil
}

func (s *mongoStorage) GetUser(ctx context.Context, id int64) (model.User, error) {
	r := s.client.Database("gnotify").Collection("users").FindOne(ctx, bson.M{"id": id})
	if r.Err() != nil {
		return model.User{}, fmt.Errorf("failed to get user: %w", r.Err())
	}

	var u struct {
		ID   int64
		Lang string
	}
	if err := r.Decode(&u); err != nil {
		return model.User{}, fmt.Errorf("failed to get user: %w", err)
	}

	return model.User{
		ID:       u.ID,
		Language: u.Lang,
	}, nil
}
