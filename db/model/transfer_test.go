package db

import (
	"context"
	"testing"
	"utils"

	"github.com/stretchr/testify/require"
)

func TestTransferTX(t *testing.T) {
	user1, _ := createFakeUser(testStore)
	user2, _ := createFakeUser(testStore)
	coin, _ := createFakeCoin(testStore)

	account1, err := createFakeAccount(testStore, user1, coin, 100)
	require.NoError(t, err)
	account2, err := createFakeAccount(testStore, user2, coin, 100)
	require.NoError(t, err)

	acc1UpdArg := UpdateAccountParams{
		ID:     account1.ID,
		Amount: int64(150),
	}
	acc2UpdArg := UpdateAccountParams{
		ID:     account2.ID,
		Amount: int64(200),
	}

	acc1New, err := testStore.UpdateAccount(context.Background(), acc1UpdArg)
	require.NoError(t, err)
	acc2New, err := testStore.UpdateAccount(context.Background(), acc2UpdArg)
	require.NoError(t, err)

	arg := SendTransferTXParams{
		FromAccountID: account1.ID,
		ToAccountID:   account2.ID,
		Amount:        int64(utils.CreateRandomInt(100)),
		Coin:          coin.Name,
	}
	trans, err := testStore.SendTransferTX(arg)
	require.NoError(t, err)

	require.Equal(t, acc1New.Balance, account1.Balance+acc1UpdArg.Amount)
	require.Equal(t, acc2New.Balance, account2.Balance+acc2UpdArg.Amount)

	require.Equal(t, trans.FromAccount.Balance, acc1New.Balance-arg.Amount)
	require.Equal(t, trans.ToAccount.Balance, acc2New.Balance+arg.Amount)

}
