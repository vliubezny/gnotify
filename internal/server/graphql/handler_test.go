package graphql

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	graphql "github.com/graph-gophers/graphql-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type resolver struct{}

func (*resolver) Hello() string {
	return "world"
}

func TestServer_graphqlHandler(t *testing.T) {
	s := `
		schema {
			query: Query
		}

		type Query {
			hello: String!
		}
	`

	schema, err := graphql.ParseSchema(s, &resolver{})
	require.NoError(t, err)

	srv := &server{
		schema: schema,
	}

	rec := httptest.NewRecorder()
	b := strings.NewReader(`{"query":"{hello}"}`)
	r := httptest.NewRequest(http.MethodPost, "/", b)

	srv.graphqlHandler(rec, r)

	body, _ := ioutil.ReadAll(rec.Result().Body)

	assert.Equal(t, http.StatusOK, rec.Result().StatusCode)
	assert.JSONEq(t, `{"data":{"hello":"world"}}`, string(body))
}
