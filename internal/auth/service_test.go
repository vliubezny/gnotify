package auth

import (
	"context"
	"testing"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const signKey = "testsecret"

func TestService_ValidateAccessToken(t *testing.T) {
	s := New(signKey)
	p := Principal{
		UserID:  1,
		IsAdmin: true,
	}

	token := mustCreateAccessToken(p)

	principal, err := s.Authenthicate(token)
	require.NoError(t, err)

	assert.Equal(t, p, principal)
}

func mustCreateAccessToken(p Principal) string {
	c := accessTokenClaims{
		TokenType: typeAccess,
		UserID:    p.UserID,
		IsAdmin:   p.IsAdmin,
		StandardClaims: jwt.StandardClaims{
			Id:        "1234",
			Issuer:    issuer,
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(10 * time.Minute).Unix(),
		},
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, c).SignedString([]byte(signKey))
	if err != nil {
		panic(err)
	}
	return token
}

func TestPrincipal_Propagate(t *testing.T) {
	p := Principal{UserID: 120}

	ctx := p.Propagate(context.Background())

	extracted := FromContext(ctx)

	assert.Equal(t, p, extracted)
}
