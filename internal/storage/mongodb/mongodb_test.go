package mongodb

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"github.com/vliubezny/gnotify/internal/model"
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

	s, err := New(uri)
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

func TestMongoStorage_GetUser(t *testing.T) {
	u := model.User{
		ID:       1,
		Language: "en",
	}

	_, err := ms.client.Database(db).Collection(users).InsertOne(ctx, bson.D{
		{Key: "id", Value: u.ID},
		{Key: "lang", Value: u.Language},
	})
	require.NoError(t, err)

	user, err := ms.GetUser(ctx, u.ID)

	require.NoError(t, err)
	assert.Equal(t, u, user)
}
