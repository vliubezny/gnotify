package graphql

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/davecgh/go-spew/spew"
	"github.com/sirupsen/logrus"
	"github.com/vliubezny/gnotify/internal/auth"
)

type loggerKey struct{}

type principalKey struct{}

// loggerMiddleware populates request context with logger and logs request entry.
func loggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := logrus.WithFields(logrus.Fields{
			"agent": r.UserAgent(),
		})
		ctx := context.WithValue(r.Context(), loggerKey{}, logger)
		logger.Debugf("%s %s", r.Method, r.RequestURI)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// recoveryMiddleware recovers after panic.
func recoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if e := recover(); e != nil {
				writeInternalError(getLogger(r), w, fmt.Sprintf("recover from panic: %s\n", spew.Sdump(e)))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

// jwtAuthMiddleware authenticates user with JWT.
func jwtAuthMiddleware(authenticator auth.Authenticator) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			l := getLogger(r)

			token := extractBearer(r)
			if token == "" {
				writeError(l, w, http.StatusUnauthorized, "missing token")
				return
			}

			principal, err := authenticator.Authenthicate(token)
			if err != nil {
				if errors.Is(err, auth.ErrInvalidToken) {
					writeError(l.WithError(err), w, http.StatusUnauthorized, "invalid access token")
					return
				}

				writeInternalError(l.WithError(err), w, "failed to validate access token")
				return
			}

			ctx := context.WithValue(r.Context(), principalKey{}, principal)
			ctx = context.WithValue(ctx, loggerKey{}, l.WithField("userID", principal.UserID))

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
