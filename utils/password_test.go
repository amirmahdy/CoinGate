package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateHashPassword(t *testing.T) {
	password := "password123"
	hash, err := CreateHashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hash)
}

func TestVerifyHashPassword(t *testing.T) {
	password := "password123"
	hash, err := CreateHashPassword(password)
	require.NoError(t, err)

	err = VerifyHashPassword(password, hash)
	require.NoError(t, err)
	err = VerifyHashPassword("wrongpassword", hash)
	require.NoError(t, err)
}
