//+build integration

package mongodb

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"github.com/vliubezny/gnotify/internal/model"
	"github.com/vliubezny/gnotify/internal/storage"
	"go.mongodb.org/mongo-driver/bson"
)

var (
	ctx = context.Background()
	ms  *mongoStorage
)

func TestMain(m *testing.M) {
	shutdown := setup()

	code := m.Run()

	shutdown()
	os.Exit(code)
}

func setup() func() {
	req := testcontainers.ContainerRequest{
		Image:        "mongo:4",
		ExposedPorts: []string{"27017/tcp"},
		WaitingFor:   wait.ForListeningPort("27017/tcp"),
	}
	c, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
	})
	if err != nil {
		logrus.WithError(err).Fatalf("failed to create container")
	}

	if err := c.Start(ctx); err != nil {
		logrus.WithError(err).Fatal("failed to start container")
	}

	host, err := c.Host(ctx)
	if err != nil {
		logrus.WithError(err).Fatal("failed to get host")
	}

	port, err := c.MappedPort(ctx, "27017")
	if err != nil {
		logrus.WithError(err).Fatal("failed to map port")
	}

	uri := fmt.Sprintf("mongodb://%s:%s", host, port)

	s, err := New(uri, "gnotify")
	if err != nil {
		logrus.WithError(err).Fatal("failed to setup mongodb storage")
	}
	ms = s.(*mongoStorage)

	shutdownFn := func() {
		if c != nil {
			c.Terminate(ctx)
		}
	}

	return shutdownFn
}

func cleanup(t *testing.T) {
	_, err := ms.db.Collection(users).DeleteMany(ctx, bson.D{})
	require.NoError(t, err)
}

func TestMongoStorage_GetUser(t *testing.T) {
	defer cleanup(t)

	u := model.User{
		ID:       1,
		Language: "en",
	}

	_, err := ms.db.Collection(users).InsertOne(ctx, bson.D{
		{Key: "id", Value: u.ID},
		{Key: "lang", Value: u.Language},
	})
	require.NoError(t, err)

	user, err := ms.GetUser(ctx, u.ID)

	require.NoError(t, err)
	assert.Equal(t, u, user)

	_, err = ms.GetUser(ctx, 100500)
	assert.True(t, errors.Is(err, storage.ErrNotFound), fmt.Sprintf("wanted %s got %s", storage.ErrNotFound, err))
}

func TestMongoStorage_UpsertUser(t *testing.T) {
	defer cleanup(t)

	newUser := model.User{
		ID:       1,
		Language: "en",
	}

	require.NoError(t, ms.UpsertUser(ctx, newUser))

	user, err := ms.GetUser(ctx, newUser.ID)
	require.NoError(t, err)
	assert.Equal(t, newUser, user)

	newUser.Language = "by"
	require.NoError(t, ms.UpsertUser(ctx, newUser))

	user, err = ms.GetUser(ctx, newUser.ID)
	require.NoError(t, err)
	assert.Equal(t, newUser, user)
}

func TestMongoStorage_GetUsers(t *testing.T) {
	defer cleanup(t)

	us, err := ms.GetUsers(ctx)
	require.NoError(t, err)
	assert.Empty(t, us)

	users := []model.User{
		{ID: 1, Language: "en"},
		{ID: 2, Language: "by"},
	}

	require.NoError(t, ms.UpsertUser(ctx, users[0]))
	require.NoError(t, ms.UpsertUser(ctx, users[1]))

	us, err = ms.GetUsers(ctx)
	require.NoError(t, err)
	assert.Equal(t, users, us)
}

func TestMongoStorage_DeleteUser(t *testing.T) {
	defer cleanup(t)

	u := model.User{
		ID:       1,
		Language: "en",
	}

	require.NoError(t, ms.UpsertUser(ctx, u))

	err := ms.DeleteUser(ctx, u.ID)
	require.NoError(t, err)

	_, err = ms.GetUser(ctx, u.ID)
	assert.True(t, errors.Is(err, storage.ErrNotFound), fmt.Sprintf("wanted %s got %s", storage.ErrNotFound, err))

	err = ms.DeleteUser(ctx, 100500)
	assert.True(t, errors.Is(err, storage.ErrNotFound), fmt.Sprintf("wanted %s got %s", storage.ErrNotFound, err))
}
