package auth

import (
	"context"
	"errors"
	"fmt"

	jwt "github.com/dgrijalva/jwt-go"
)

//go:generate mockgen -destination=./mock/mock.go -package=mock -source=service.go

const (
	typeAccess = "access"
	issuer     = "gstore.auth"
)

var (
	// ErrInvalidToken states that provided token is invalid.
	ErrInvalidToken = errors.New("invalid token")
)

type principalKey struct{}

// Principal represents authenticated user and its roles.
type Principal struct {
	UserID  int64
	IsAdmin bool
}

// Propagate returns copy of parent context with principal value.
func (p Principal) Propagate(ctx context.Context) context.Context {
	return context.WithValue(ctx, principalKey{}, p)
}

// FromContext extracts principal form context or panics.
func FromContext(ctx context.Context) Principal {
	return ctx.Value(principalKey{}).(Principal)
}

// Authenticator authenthicates user by token.
type Authenticator interface {
	Authenthicate(token string) (Principal, error)
}

type authService struct {
	signKey []byte
}

type accessTokenClaims struct {
	TokenType string `json:"type,omitempty"`
	UserID    int64  `json:"userId,omitempty"`
	IsAdmin   bool   `json:"admin,omitempty"`
	jwt.StandardClaims
}

// New creates instance of Authenticator.
func New(signKey string) Authenticator {
	return &authService{
		signKey: []byte(signKey),
	}
}

func (s *authService) Authenthicate(token string) (Principal, error) {
	at, err := jwt.ParseWithClaims(token, &accessTokenClaims{}, func(t *jwt.Token) (interface{}, error) {
		if t.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, errors.New("token must be signed with HS256 alg")
		}
		return s.signKey, nil
	})

	if err != nil {
		return Principal{}, fmt.Errorf("unable to parse claims: %w", err)
	}

	claims, ok := at.Claims.(*accessTokenClaims)
	if !ok || !at.Valid {
		return Principal{}, errors.New("invalid token")
	}
	if claims.TokenType != typeAccess {
		return Principal{}, fmt.Errorf("invalid access token: type %s", claims.TokenType)
	}
	return Principal{UserID: claims.UserID, IsAdmin: claims.IsAdmin}, nil
}
