package graphql

import (
	"encoding/json"
	"net/http"

	"github.com/graphql-go/graphql"
)

// graphqlHandler handles graphql requests.
func (s *server) graphqlHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")
	result := graphql.Do(graphql.Params{
		Schema:        s.schema,
		RequestString: query,
	})
	json.NewEncoder(w).Encode(result)
}
