package service

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/vliubezny/gnotify/internal/model"
	"github.com/vliubezny/gnotify/internal/storage"
	"github.com/vliubezny/gnotify/internal/storage/mock"
)

var (
	ctx = context.Background()
)

func TestService_GetUser(t *testing.T) {
	testCases := []struct {
		desc  string
		rUser model.User
		rErr  error
		user  model.User
		err   error
	}{
		{
			desc:  "success",
			rUser: model.User{ID: 1, Language: "en"},
			rErr:  nil,
			user:  model.User{ID: 1, Language: "en"},
			err:   nil,
		},
		{
			desc:  "ErrNotFound",
			rUser: model.User{},
			rErr:  storage.ErrNotFound,
			user:  model.User{},
			err:   ErrNotFound,
		},
		{
			desc:  "unexpected error",
			rUser: model.User{},
			rErr:  assert.AnError,
			user:  model.User{},
			err:   assert.AnError,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			id := int64(1)

			st := mock.NewMockStorage(ctrl)
			st.EXPECT().GetUser(ctx, id).Return(tC.rUser, tC.rErr)

			s := New(st)

			user, err := s.GetUser(ctx, id)
			assert.True(t, errors.Is(err, tC.err), fmt.Sprintf("wanted %s got %s", tC.err, err))
			assert.Equal(t, tC.user, user)
		})
	}
}

func TestService_GetUsers(t *testing.T) {
	testUsers := []model.User{
		{ID: 1, Language: "by"},
		{ID: 2, Language: "ru"},
	}

	testCases := []struct {
		desc   string
		rUsers []model.User
		rErr   error
		users  []model.User
		err    error
	}{
		{
			desc:   "success",
			rUsers: testUsers,
			rErr:   nil,
			users:  testUsers,
			err:    nil,
		},
		{
			desc:   "unexpected error",
			rUsers: nil,
			rErr:   assert.AnError,
			users:  nil,
			err:    assert.AnError,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			st := mock.NewMockStorage(ctrl)
			st.EXPECT().GetUsers(ctx).Return(tC.rUsers, tC.rErr)

			s := New(st)

			users, err := s.GetUsers(ctx)
			assert.True(t, errors.Is(err, tC.err), fmt.Sprintf("wanted %s got %s", tC.err, err))
			assert.Equal(t, tC.users, users)
		})
	}
}

func TestService_UpsertUser(t *testing.T) {
	testCases := []struct {
		desc  string
		rUser model.User
		rErr  error
		user  model.User
		err   error
	}{
		{
			desc:  "success",
			rUser: model.User{ID: 1, Language: "en"},
			rErr:  nil,
			user:  model.User{ID: 1, Language: "en"},
			err:   nil,
		},
		{
			desc:  "unexpected error",
			rUser: model.User{},
			rErr:  assert.AnError,
			user:  model.User{ID: 1, Language: "en"},
			err:   assert.AnError,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			st := mock.NewMockStorage(ctrl)
			st.EXPECT().UpsertUser(ctx, tC.user).Return(tC.rErr)

			s := New(st)

			err := s.UpsertUser(ctx, tC.user)

			assert.True(t, errors.Is(err, tC.err), fmt.Sprintf("wanted %s got %s", tC.err, err))
		})
	}
}

func TestService_DeleteUser(t *testing.T) {
	testCases := []struct {
		desc string
		rErr error
		err  error
	}{
		{
			desc: "success",
			rErr: nil,
			err:  nil,
		},
		{
			desc: "ErrNotFound",
			rErr: storage.ErrNotFound,
			err:  ErrNotFound,
		},
		{
			desc: "unexpected error",
			rErr: assert.AnError,
			err:  assert.AnError,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			id := int64(1)

			st := mock.NewMockStorage(ctrl)
			st.EXPECT().DeleteUser(ctx, id).Return(tC.rErr)

			s := New(st)

			err := s.DeleteUser(ctx, id)
			assert.True(t, errors.Is(err, tC.err), fmt.Sprintf("wanted %s got %s", tC.err, err))
		})
	}
}
