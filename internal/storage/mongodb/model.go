package mongodb

import (
	"github.com/vliubezny/gnotify/internal/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type user struct {
	ID      int64    `bson:"id"`
	Lang    string   `bson:"lang"`
	Devices []device `bson:"devices,omitempty"`
}

func (u user) toModel() model.User {
	mUser := model.User{
		ID:       u.ID,
		Language: u.Lang,
	}

	if len(u.Devices) > 0 {
		mUser.Devices = make([]model.Device, len(u.Devices))

		for i, d := range u.Devices {
			mUser.Devices[i] = d.toModel()
		}
	}

	return mUser
}

type device struct {
	ID           primitive.ObjectID `bson:"_id"`
	Name         string             `bson:"name"`
	PriceChanged bool               `bson:"priceChanged,omitempty"`
	Frequency    string             `bson:"frequency,omitempty"`
}

func (d device) toModel() model.Device {
	return model.Device{
		ID:   d.ID.String(),
		Name: d.Name,
		Settings: model.NotificationSettings{
			PriceChanged: d.PriceChanged,
			Frequency:    d.Frequency,
		},
	}
}
