package graphql

import (
	"encoding/json"
	"net/http"
	"runtime/debug"
	"strings"

	"github.com/go-chi/chi"
	"github.com/graphql-go/graphql"
	"github.com/sirupsen/logrus"
	"github.com/vliubezny/gnotify/internal/auth"
)

type server struct {
	schema graphql.Schema
}

// SetupRouter setups routes and handlers.
func SetupRouter(r chi.Router, authenticator auth.Authenticator) error {
	s, err := NewSchema()
	if err != nil {
		return err
	}

	srv := &server{
		schema: s,
	}

	r.Use(
		loggerMiddleware,
		recoveryMiddleware,
		jwtAuthMiddleware(authenticator),
	)

	r.Get("/graphql", srv.graphqlHandler)

	return nil
}

func getLogger(r *http.Request) logrus.FieldLogger {
	return r.Context().Value(loggerKey{}).(logrus.FieldLogger)
}

func extractBearer(r *http.Request) string {
	auth := r.Header.Get("Authorization")
	if len(auth) > 7 && strings.ToUpper(auth[0:7]) == "BEARER " {
		return auth[7:]
	}
	return ""
}

func writeError(l logrus.FieldLogger, w http.ResponseWriter, code int, message string) {
	l.Error(message)

	body, _ := json.Marshal(errorResponse{
		Error: message,
	})

	w.WriteHeader(code)
	w.Write(body)
}

func writeInternalError(l logrus.FieldLogger, w http.ResponseWriter, message string) {
	l.Errorf("%s\n%s", message, string(debug.Stack()))

	body, _ := json.Marshal(errorResponse{
		Error: "internal error",
	})

	w.WriteHeader(http.StatusInternalServerError)
	w.Write(body)
}

func writeOK(l logrus.FieldLogger, w http.ResponseWriter, payload interface{}) {
	body, err := json.Marshal(payload)
	if err != nil {
		writeInternalError(l.WithError(err), w, "fail to serialize payload")
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(body)
}
