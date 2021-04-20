package graphql

import (
	"context"
	"fmt"
	"io/ioutil"
	"strconv"

	graphql "github.com/graph-gophers/graphql-go"
	"github.com/vliubezny/gnotify/internal/auth"
	"github.com/vliubezny/gnotify/internal/service"
	"golang.org/x/text/language"
	"golang.org/x/text/language/display"
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
type RootResolver struct {
	svc service.Service
}

// CurrentUser resolves current user data.
func (r *RootResolver) CurrentUser(ctx context.Context) (User, error) {
	p := ctx.Value(principalKey{}).(auth.Principal)
	u, err := r.svc.GetUser(ctx, p.UserID)
	if err != nil {
		return User{}, fmt.Errorf("failed to resolve current user: %w", err)
	}

	var user User
	user.ID = graphql.ID(strconv.FormatInt(u.ID, 10))
	user.Settings.Language.Code = u.Language

	return user, nil
}

// NewSchema parses and creates new graphql schema.
func NewSchema(svc service.Service) (*graphql.Schema, error) {
	schema, err := ioutil.ReadFile("../../../static/schema.graphql")
	if err != nil {
		return nil, err
	}

	return graphql.ParseSchema(string(schema), &RootResolver{svc: svc}, graphql.UseFieldResolvers())
}
