package graphql

import (
	"testing"

	"github.com/graphql-go/graphql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Schema(t *testing.T) {
	testCases := []struct {
		desc  string
		query string
		data  interface{}
	}{
		{
			desc:  "query language",
			query: `{language{name,code}}`,
			data: map[string]interface{}{
				"language": map[string]interface{}{
					"code": "ru",
					"name": "Russian",
				},
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			s, err := NewSchema()
			require.NoError(t, err)

			result := graphql.Do(graphql.Params{
				Schema:        s,
				RequestString: tc.query,
			})

			assert.Empty(t, result.Errors)
			assert.Equal(t, tc.data, result.Data)
		})
	}
}
