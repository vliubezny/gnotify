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
	users = "users"
)

type mongoStorage struct {
	db *mongo.Database
}

// New creates mongodb storage.
func New(uri, db string) (storage.Storage, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, err
	}

	ms := &mongoStorage{
		db: client.Database(db),
	}

	if err = createScheme(ms.db); err != nil {
		return nil, err
	}

	return ms, nil
}

func createScheme(db *mongo.Database) error {
	ctx := context.Background()

	_, err := db.Collection(users).Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "id", Value: 1}},
		Options: options.Index().SetName("id").SetUnique(true),
	})

	return err
}

func (s *mongoStorage) GetUser(ctx context.Context, id int64) (model.User, error) {
	r := s.db.Collection(users).FindOne(ctx, bson.M{"id": id})
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
