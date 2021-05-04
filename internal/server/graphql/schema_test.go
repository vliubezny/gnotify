package graphql

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vliubezny/gnotify/internal/auth"
	"github.com/vliubezny/gnotify/internal/model"
	"github.com/vliubezny/gnotify/internal/service/mock"
)

var ctx = context.Background()

func TestSchema_currentUser(t *testing.T) {
	user := model.User{
		ID:       1,
		Language: "ru",
		Devices: []model.Device{
			{
				ID:   "132323",
				Name: "Chrome",
				Settings: model.NotificationSettings{
					Frequency:    model.Daily,
					PriceChanged: true,
				},
			},
		},
	}

	testCases := []struct {
		desc      string
		principal auth.Principal
		rUser     model.User
		rErr      error
		query     string
		data      string
	}{
		{
			desc:      "query current user language",
			principal: auth.Principal{UserID: 1},
			rUser:     user,
			query: `{
					currentUser {
						settings {
							language {
								name
								code
							}
						}
						devices {
							id
							name
							settings {
								frequency
								priceChanged
							}
						}
					}
				}`,
			data: `{
					"data": {
						"currentUser": {
							"settings": {
								"language": {
									"code": "ru",
									"name": "Russian"
								}
							},
							"devices": [
								{
									"id": "132323",
									"name": "Chrome",
									"settings": {
										"frequency": "DAILY",
										"priceChanged": true
									}
								}
							]
						}
					}
				}`,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			svc := mock.NewMockService(ctrl)
			s, err := NewSchema(svc)
			require.NoError(t, err)

			c := context.WithValue(ctx, principalKey{}, tc.principal)

			svc.EXPECT().GetUser(gomock.Any(), tc.rUser.ID).Return(tc.rUser, tc.rErr)

			result := s.Exec(c, tc.query, "", nil)

			json, err := json.Marshal(result)
			require.NoError(t, err)

			assert.JSONEq(t, tc.data, string(json))
		})
	}
}

func TestSchema_addDeviceForCurrentUser(t *testing.T) {
	device := model.Device{
		ID:   "132323",
		Name: "Chrome",
		Settings: model.NotificationSettings{
			Frequency:    model.Daily,
			PriceChanged: true,
		},
	}

	testCases := []struct {
		desc      string
		principal auth.Principal
		vars      map[string]interface{}
		rDevice   model.Device
		rErr      error
		query     string
		data      string
	}{
		{
			desc:      "addDeviceForCurrentUser",
			principal: auth.Principal{UserID: 1},
			vars: map[string]interface{}{
				"device": map[string]interface{}{
					"name":         "Chrome",
					"priceChanged": true,
					"frequency":    "DAILY",
				},
			},
			rDevice: device,
			query: `mutation ($device: DeviceInput!) {
				addDeviceForCurrentUser(device: $device) {
					id
					name
					settings {
						frequency
						priceChanged
					}
				}
			}`,
			data: `{
				"data": {
					"addDeviceForCurrentUser": {
						"id": "132323",
						"name": "Chrome",
						"settings": {
							"frequency": "DAILY",
							"priceChanged": true
						}
					}
				}
			}`,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			svc := mock.NewMockService(ctrl)
			s, err := NewSchema(svc)
			require.NoError(t, err)

			c := context.WithValue(ctx, principalKey{}, tc.principal)

			svc.EXPECT().AddDevice(gomock.Any(), tc.principal.UserID, gomock.Any()).
				DoAndReturn(func(ctx context.Context, id int64, device model.Device) (model.Device, error) {
					assert.Empty(t, device.ID)
					device.ID = tc.rDevice.ID
					assert.Equal(t, tc.rDevice, device)

					return tc.rDevice, tc.rErr
				})

			result := s.Exec(c, tc.query, "", tc.vars)

			json, err := json.Marshal(result)
			require.NoError(t, err)

			assert.JSONEq(t, tc.data, string(json))
		})
	}
}
