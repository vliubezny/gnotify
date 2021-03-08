package graphql

import (
	"github.com/go-chi/chi"
	"github.com/graphql-go/graphql"
)

type server struct {
	schema graphql.Schema
}

// SetupRouter setups routes and handlers.
func SetupRouter(r chi.Router) error {
	s, err := NewSchema()
	if err != nil {
		return err
	}

	srv := &server{
		schema: s,
	}

	r.Get("/graphql", srv.graphqlHandler)

	return nil
}
