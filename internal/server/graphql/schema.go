package graphql

import (
	"github.com/graphql-go/graphql"

	"golang.org/x/text/language"
	"golang.org/x/text/language/display"
)

// Language represents user language of choice.
type Language struct {
	Code string `json:"code"`
}

// Settings represents notification settings.
type Settings struct {
	Language Language `json:"language"`
}

// User represents gstore user.
type User struct {
	ID       int64    `json:"id"`
	Settings Settings `json:"settings"`
}

// LanguageType graphql type.
var LanguageType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Language",
	Fields: graphql.Fields{
		"code": &graphql.Field{
			Type: graphql.String,
		},
		"name": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if l, ok := p.Source.(Language); ok {
					tag := language.Make(l.Code)
					return display.English.Tags().Name(tag), nil
				}
				return nil, nil
			},
		},
	},
})

// SettingsType graphql type.
var SettingsType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Settings",
	Fields: graphql.Fields{
		"language": &graphql.Field{
			Type: LanguageType,
		},
	},
})

// NewSchema defines notifications schema.
func NewSchema() (graphql.Schema, error) {
	fields := graphql.Fields{
		"language": &graphql.Field{
			Type: LanguageType,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return Language{Code: "ru"}, nil
			},
		},
	}
	rootQuery := graphql.ObjectConfig{Name: "Query", Fields: fields}
	schemaConfig := graphql.SchemaConfig{Query: graphql.NewObject(rootQuery)}
	return graphql.NewSchema(schemaConfig)
}
