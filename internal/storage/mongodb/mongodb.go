package mongodb

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
		if r.Err() == mongo.ErrNoDocuments {
			return model.User{}, storage.ErrNotFound
		}
		return model.User{}, fmt.Errorf("failed to get user: %w", r.Err())
	}

	var u user
	if err := r.Decode(&u); err != nil {
		return model.User{}, fmt.Errorf("failed to get user: %w", err)
	}

	return u.toModel(), nil
}

func (s *mongoStorage) UpsertUser(ctx context.Context, user model.User) error {
	_, err := s.db.Collection(users).UpdateOne(ctx, bson.M{"id": user.ID},
		bson.M{
			"$setOnInsert": bson.D{
				{Key: "id", Value: user.ID},
			},
			"$set": bson.D{
				{Key: "lang", Value: user.Language},
			},
		}, options.Update().SetUpsert(true))

	if err != nil {
		return fmt.Errorf("failed to upsert user: %w", err)
	}

	return nil
}

func (s *mongoStorage) GetUsers(ctx context.Context) ([]model.User, error) {
	cursor, err := s.db.Collection(users).Find(ctx, bson.D{})
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}
	defer cursor.Close(ctx)

	var users []user
	if err := cursor.All(ctx, &users); err != nil {
		return nil, fmt.Errorf("failed to read users: %w", err)
	}

	mUsers := make([]model.User, len(users))
	for i := range users {
		mUsers[i] = users[i].toModel()
	}

	return mUsers, nil
}

func (s *mongoStorage) DeleteUser(ctx context.Context, id int64) error {
	r, err := s.db.Collection(users).DeleteOne(ctx, bson.M{"id": id})
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	if r.DeletedCount == 0 {
		return storage.ErrNotFound
	}

	return nil
}

func (s *mongoStorage) AddDevice(ctx context.Context, userID int64, input model.Device) (model.Device, error) {
	d := device{
		ID:           primitive.NewObjectID(),
		Name:         input.Name,
		PriceChanged: input.Settings.PriceChanged,
		Frequency:    input.Settings.Frequency,
	}

	r := s.db.Collection(users).FindOneAndUpdate(ctx, bson.M{"id": userID},
		bson.M{
			"$push": bson.D{
				{Key: "devices", Value: d},
			},
		}, options.FindOneAndUpdate().SetReturnDocument(options.After).
			SetProjection(bson.M{
				"devices": bson.D{
					{Key: "$slice", Value: -1},
				},
			}))

	if r.Err() != nil {
		if r.Err() == mongo.ErrNoDocuments {
			return model.Device{}, storage.ErrNotFound
		}
		return model.Device{}, fmt.Errorf("failed to add device: %w", r.Err())
	}

	var u user
	if err := r.Decode(&u); err != nil || len(u.Devices) == 0 {
		return model.Device{}, fmt.Errorf("failed to add device: %w", err)
	}

	return u.Devices[0].toModel(), nil
}
