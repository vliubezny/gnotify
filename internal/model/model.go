package model

// Frequency enum.
const (
	Hourly = "HOURLY"
	Daily  = "DAILY"
	Weekly = "WEEKLY"
	Never  = "NEVER"
)

// User represents user notifications preferences.
type User struct {
	ID       int64
	Language string
	Devices  []Device
}

type Device struct {
	ID       string
	Name     string
	Settings NotificationSettings
}

// NotificationSettings represents notification settings for the device.
type NotificationSettings struct {
	PriceChanged bool
	Frequency    string
}
