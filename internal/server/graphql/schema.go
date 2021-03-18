package graphql

import (
	"context"
	"io/ioutil"

	graphql "github.com/graph-gophers/graphql-go"
	"golang.org/x/text/language"
	"golang.org/x/text/language/display"
)

// Frequency enum.
const (
	Hourly = "HOURLY"
	Daily  = "DAILY"
	Weekly = "WEEKLY"
	Never  = "NEVER"
)

// User represents user of the service.
type User struct {
	ID       graphql.ID
	Settings Settings
	Devices  []Device
}

// Settings represents user settings.
type Settings struct {
	Language Language
}

// Language represents user language of choice.
type Language struct {
	Code string
}

// Device represents user device.
type Device struct {
	ID       graphql.ID
	Name     string
	Settings NotificationSettings
}

// NotificationSettings represents notification settings for the device.
type NotificationSettings struct {
	PriceChanged bool
	Frequency    string
}

// Name returns language name in plain English.
func (lang Language) Name() string {
	tag := language.Make(lang.Code)
	return display.English.Tags().Name(tag)
}

// RootResolver defines root resolvers.
type RootResolver struct{}

// CurrentUser resolves current user data.
func (*RootResolver) CurrentUser(ctx context.Context) User {
	return User{
		Settings: Settings{
			Language: Language{Code: "ru"},
		},
	}
}

// NewSchema parses and creates new graphql schema.
func NewSchema() (*graphql.Schema, error) {
	schema, err := ioutil.ReadFile("../../../static/schema.graphql")
	if err != nil {
		return nil, err
	}

	return graphql.ParseSchema(string(schema), &RootResolver{}, graphql.UseFieldResolvers())
}
