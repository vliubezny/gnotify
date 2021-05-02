package graphql

import (
	"context"
	"fmt"
	"io/ioutil"
	"strconv"

	graphql "github.com/graph-gophers/graphql-go"
	"github.com/vliubezny/gnotify/internal/auth"
	"github.com/vliubezny/gnotify/internal/model"
	"github.com/vliubezny/gnotify/internal/service"
	"golang.org/x/text/language"
	"golang.org/x/text/language/display"
)

type userResolver struct {
	user model.User
}

func (r *userResolver) ID() graphql.ID {
	return graphql.ID(strconv.FormatInt(r.user.ID, 10))
}

func (r *userResolver) Settings() settingsResolver {
	return settingsResolver{
		Language: languageResolver{Code: r.user.Language},
	}
}

func (r *userResolver) Devices() []deviceResolver {
	dr := make([]deviceResolver, len(r.user.Devices))

	for i, d := range r.user.Devices {
		dr[i] = deviceResolver{d}
	}

	return dr
}

type settingsResolver struct {
	Language languageResolver
}

type languageResolver struct {
	Code string
}

// Name returns language name in plain English.
func (r languageResolver) Name() string {
	tag := language.Make(r.Code)
	return display.English.Tags().Name(tag)
}

type deviceResolver struct {
	device model.Device
}

func (r deviceResolver) ID() graphql.ID {
	return graphql.ID(r.device.ID)
}

func (r deviceResolver) Name() string {
	return r.device.Name
}

func (r deviceResolver) Settings() notificationSettingsResolver {
	return notificationSettingsResolver{settings: r.device.Settings}
}

type notificationSettingsResolver struct {
	settings model.NotificationSettings
}

func (r notificationSettingsResolver) PriceChanged() bool {
	return r.settings.PriceChanged
}

func (r notificationSettingsResolver) Frequency() string {
	return r.settings.Frequency
}

// RootResolver defines root resolvers.
type RootResolver struct {
	svc service.Service
}

// CurrentUser resolves current user data.
func (r *RootResolver) CurrentUser(ctx context.Context) (*userResolver, error) {
	p := ctx.Value(principalKey{}).(auth.Principal)
	u, err := r.svc.GetUser(ctx, p.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve current user: %w", err)
	}

	return &userResolver{user: u}, nil
}

// NewSchema parses and creates new graphql schema.
func NewSchema(svc service.Service) (*graphql.Schema, error) {
	schema, err := ioutil.ReadFile("../../../static/schema.graphql")
	if err != nil {
		return nil, err
	}

	return graphql.ParseSchema(string(schema), &RootResolver{svc: svc}, graphql.UseFieldResolvers())
}
