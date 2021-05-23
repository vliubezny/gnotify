package graphql

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"runtime"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/vliubezny/gnotify/internal/auth"
	authMock "github.com/vliubezny/gnotify/internal/auth/mock"
)

func Test_loggerMiddleware(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
	hook := test.NewGlobal()
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/v1/test", nil)
	req.Header.Set("User-Agent", "curl")

	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(nil)
	})

	loggerMiddleware(h).ServeHTTP(rec, req)

	log := hook.LastEntry()
	require.NotNil(t, log)
	assert.Equal(t, logrus.DebugLevel, log.Level)
	assert.Equal(t, "POST /v1/test", log.Message, "Incorrect request entry")
	assert.Equal(t, "curl", log.Data["agent"], "Incorrect user agent")
}

func Test_recoveryMiddleware(t *testing.T) {
	logger, hook := test.NewNullLogger()
	ctx := context.WithValue(context.Background(), loggerKey{}, logger)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest("", "/", nil).WithContext(ctx)

	_, file, _, ok := runtime.Caller(0)
	require.True(t, ok)

	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("test panic")
	})

	recoveryMiddleware(h).ServeHTTP(rec, req)

	body, _ := ioutil.ReadAll(rec.Result().Body)

	assert.Equal(t, http.StatusInternalServerError, rec.Result().StatusCode)
	assert.Equal(t, `{"error":"internal error"}`, string(body))

	log := hook.LastEntry()
	require.NotNil(t, log)
	assert.Equal(t, logrus.ErrorLevel, log.Level)
	assert.Contains(t, log.Message, "test panic", "Missing panic message")
	assert.Contains(t, log.Message, file, "Missing stacktrace")
}

func Test_jwtAuthMiddleware(t *testing.T) {
	principal := auth.Principal{UserID: 1}
	testCases := []struct {
		desc  string
		token string
		err   error
		rcode int
		rdata string
	}{
		{
			desc:  "allow valid token",
			token: "testtoken",
			err:   nil,
			rcode: http.StatusOK,
			rdata: `{"result":"OK"}`,
		},
		{
			desc:  "missing token",
			token: "",
			err:   errSkip,
			rcode: http.StatusUnauthorized,
			rdata: `{"error":"missing token"}`,
		},
		{
			desc:  "invalid token",
			token: "testtoken",
			err:   auth.ErrInvalidToken,
			rcode: http.StatusUnauthorized,
			rdata: `{"error":"invalid access token"}`,
		},
		{
			desc:  "internal error",
			token: "testtoken",
			err:   assert.AnError,
			rcode: http.StatusInternalServerError,
			rdata: `{"error":"internal error"}`,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			logger, _ := test.NewNullLogger()
			ctx := context.WithValue(context.Background(), loggerKey{}, logger)
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/", nil).WithContext(ctx)

			if tc.token != "" {
				req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tc.token))
			}

			h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				p := auth.FromContext(r.Context())
				assert.Equal(t, principal, p)

				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"result":"OK"}`))
			})

			a := authMock.NewMockAuthenticator(ctrl)
			if tc.err != errSkip {
				a.EXPECT().Authenthicate(tc.token).Return(principal, tc.err)
			}

			jwtAuthMiddleware(a)(h).ServeHTTP(rec, req)

			body, _ := ioutil.ReadAll(rec.Result().Body)

			assert.Equal(t, tc.rcode, rec.Result().StatusCode)
			assert.JSONEq(t, tc.rdata, string(body))
		})
	}
}
