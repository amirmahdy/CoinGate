package token

import (
	"errors"
	"testing"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/require"
)

func TestNewJWTMaker(t *testing.T) {
	secret := "alongsecrettomeet32bytesrequirement"
	maker, err := NewJWTMaker(secret)
	require.NoError(t, err)
	require.NotNil(t, maker)

	invalidSecret := "short"
	_, err = NewJWTMaker(invalidSecret)
	require.Error(t, err)
}

func TestJWTMakerExpiredToken(t *testing.T) {
	secret := "alongsecrettomeet32bytesrequirement"
	maker, err := NewJWTMaker(secret)
	require.NoError(t, err)

	username := "testuser"
	duration := -time.Minute // negative duration to create an expired token

	token, err := maker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	// Verify that the token is expired
	_, err = maker.VerifyToken(token)
	require.Error(t, err)
	require.True(t, errors.Is(err, ErrExpiredToken))
}

func TestJWTMakerCreateToken(t *testing.T) {
	secret := "alongsecrettomeet32bytesrequirement"
	maker, err := NewJWTMaker(secret)
	require.NoError(t, err)

	username := "testuser"
	duration := time.Minute

	token, err := maker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	// Verify the token can be parsed and contains the expected claims.
	payload, err := maker.VerifyToken(token)
	require.NoError(t, err)
	require.Equal(t, username, payload.Username)
	require.WithinDuration(t, time.Now().Add(duration), payload.ExpireAt, time.Second)
}

func TestJWTMakerInvalidAlgorithm(t *testing.T) {
	secret := "alongsecrettomeet32bytesrequirement"
	maker, err := NewJWTMaker(secret)
	require.NoError(t, err)

	payload := NewPayload("testuser", time.Minute)

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodNone, payload)
	token, err := jwtToken.SignedString(jwt.UnsafeAllowNoneSignatureType)
	require.NoError(t, err)

	payload, err = maker.VerifyToken(token)
	require.Error(t, err)
	require.True(t, errors.Is(err, ErrorInvalidToken))
	require.Nil(t, payload)

}
