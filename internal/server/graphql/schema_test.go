package graphql

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var ctx = context.Background()

func Test_Schema(t *testing.T) {
	testCases := []struct {
		desc  string
		query string
		data  string
	}{
		{
			desc:  "query language",
			query: `{language{name,code}}`,
			data: `{
					"data": {
						"language": {
							"code":"ru",
							"name":"Russian"
						}
					}
				}`,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			s, err := NewSchema()
			require.NoError(t, err)

			result := s.Exec(ctx, tc.query, "", nil)

			json, err := json.Marshal(result)
			require.NoError(t, err)

			assert.JSONEq(t, tc.data, string(json))
		})
	}
}
