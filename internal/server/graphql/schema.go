package graphql

import (
	"io/ioutil"

	graphql "github.com/graph-gophers/graphql-go"
	"golang.org/x/text/language"
	"golang.org/x/text/language/display"
)

// Language represents user language of choice.
type Language struct {
	Code string
}

// Name returns language name in plain English.
func (lang Language) Name() string {
	tag := language.Make(lang.Code)
	return display.English.Tags().Name(tag)
}

// RootResolver defines root resolvers.
type RootResolver struct{}

// Language resolves language query.
func (*RootResolver) Language() Language {
	return Language{Code: "ru"}
}

// NewSchema parses and creates new graphql schema.
func NewSchema() (*graphql.Schema, error) {
	schema, err := ioutil.ReadFile("../../../static/schema.graphql")
	if err != nil {
		return nil, err
	}

	return graphql.ParseSchema(string(schema), &RootResolver{}, graphql.UseFieldResolvers())
}
