package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi"
	"github.com/jessevdk/go-flags"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"

	"github.com/vliubezny/gnotify/internal/auth"
	"github.com/vliubezny/gnotify/internal/server/graphql"
	"github.com/vliubezny/gnotify/internal/service"
	"github.com/vliubezny/gnotify/internal/storage/mongodb"
)

var errTerminated = errors.New("terminated")

var opts = struct {
	Host string `long:"http.host" env:"HTTP_HOST" default:"0.0.0.0" description:"IP address to listen"`
	Port int    `long:"http.port" env:"HTTP_PORT" default:"8080" description:"port to listen"`

	LogLevel string `long:"log.level" env:"LOG_LEVEL" default:"debug" description:"Log level" choice:"debug" choice:"info" choice:"warning" choice:"error"`

	SignKey string `long:"auth.signkey" env:"AUTH_SIGN_KEY" default:"changeme" description:"sign key for JWT"`

	MongoDBURI  string `long:"mongodb.uri" env:"MONGODB_URI" default:"mongodb://localhost:27017"`
	MongoDBName string `long:"mongodb.name" env:"MONGODB_NAME" default:"gnotify"`
}{}

func main() {
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	parser := flags.NewParser(&opts, flags.Default)
	parser.Name = "gnotify"
	parser.LongDescription = "Starts gnotify server."

	if _, err := parser.Parse(); err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		}
		logrus.WithError(err).Fatal("failed to parse flags")
	}

	lvl, _ := logrus.ParseLevel(opts.LogLevel)
	logrus.SetLevel(lvl)

	logrus.Info("starting service")
	logrus.Infof("%+v", opts) // can print secrets!

	stg, err := mongodb.New(opts.MongoDBURI, opts.MongoDBName)
	if err != nil {
		logrus.WithError(err).Fatal("failed to setup storage")
	}

	svc := service.New(stg)

	r := chi.NewMux()
	a := auth.New(opts.SignKey)

	if err := graphql.SetupRouter(r, a, svc); err != nil {
		logrus.WithError(err).Fatal("failed to setup graphql")
	}

	srv := http.Server{
		Addr:    fmt.Sprintf("%s:%d", opts.Host, opts.Port),
		Handler: r,
	}

	gr, _ := errgroup.WithContext(context.Background())
	gr.Go(srv.ListenAndServe)

	gr.Go(func() error {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

		s := <-sigs
		logrus.Infof("terminating by %s signal", s)

		if err := srv.Shutdown(context.Background()); err != nil {
			logrus.WithError(err).Error("failed to gracefully shutdown server")
		}

		return errTerminated
	})

	logrus.Info("service started")

	if err := gr.Wait(); err != nil && !errors.Is(err, errTerminated) && !errors.Is(err, http.ErrServerClosed) {
		logrus.WithError(err).Fatal("service unexpectedly stopped")
	}
}
