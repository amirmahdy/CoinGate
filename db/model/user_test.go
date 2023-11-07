package db

import (
	"context"
	"testing"
	"utils"

	"github.com/stretchr/testify/require"
)

func TestCreateUser(t *testing.T) {
	hash, err := utils.CreateHashPassword("secret")
	require.NoError(t, err)

	arg := CreateUserParams{
		Username:       utils.CreateRandomName(),
		FullName:       utils.CreateRandomName(),
		Email:          utils.CreateRandomEmail(),
		HashedPassword: hash,
	}
	user, err := testStore.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)
	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.Email, user.Email)
	require.Equal(t, true, utils.VerifyHashPassword("secret", user.HashedPassword))
}
