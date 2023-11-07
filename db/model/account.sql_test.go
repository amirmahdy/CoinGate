package db

import (
	"context"
	"testing"
	"utils"

	"github.com/stretchr/testify/require"
)

func createFakeUser(store Store) (User, error) {
	arg := CreateUserParams{
		Username:       utils.CreateRandomName(),
		HashedPassword: utils.CreateRandomName(),
		FullName:       utils.CreateRandomName(),
		Email:          utils.CreateRandomEmail(),
	}

	user, err := store.CreateUser(context.Background(), arg)
	return user, err
}

func createFakeCoin(store Store) (Coin, error) {
	coin, err := store.CreateCoin(context.Background(), utils.CreateRandomName())
	return coin, err
}

func createFakeAccount(store Store, user User, coin Coin, balance int64) (Account, error) {
	arg := CreateAccountParams{
		Username: user.Username,
		Balance:  balance,
		Coin:     coin.Name,
	}

	account, err := store.CreateAccount(context.Background(), arg)
	return account, err
}

func TestCreateAccount(t *testing.T) {
	user, err := createFakeUser(testStore)
	require.NoError(t, err)

	coin, err := createFakeCoin(testStore)
	require.NoError(t, err)

	balance := int64(utils.CreateRandomInt(1000))
	account, err := testStore.CreateAccount(context.Background(), CreateAccountParams{
		Username: user.Username,
		Balance:  balance,
		Coin:     coin.Name,
	})
	require.NoError(t, err)
	require.NotZero(t, account.ID)
	require.Equal(t, user.Username, account.Username)
	require.Equal(t, balance, account.Balance)
	require.Equal(t, coin, account.Coin)
	require.NotZero(t, account.CreatedAt)
}
