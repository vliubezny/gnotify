package graphql

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/graphql-go/graphql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServer_graphqlHandler(t *testing.T) {
	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: graphql.NewObject(graphql.ObjectConfig{Name: "Query", Fields: graphql.Fields{
			"hello": &graphql.Field{
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return "world", nil
				},
			},
		}}),
	})
	require.NoError(t, err)

	srv := &server{
		schema: schema,
	}

	rec := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/?query={hello}", nil)

	srv.graphqlHandler(rec, r)

	body, _ := ioutil.ReadAll(rec.Result().Body)

	assert.Equal(t, http.StatusOK, rec.Result().StatusCode)
	assert.JSONEq(t, `{"data":{"hello":"world"}}`, string(body))
}
