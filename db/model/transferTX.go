package db

import "context"

type SendTransferTXParams struct {
	FromAccountID int64  `json:"from_account_id"`
	ToAccountID   int64  `json:"to_account_id"`
	Amount        int64  `json:"amount"`
	Coin          string `json:"coin"`
}
type SendTransactionResultType struct {
	Transfer    Transfer
	FromAccount Account
	ToAccount   Account
}

func (store *SQLStore) SendTransferTX(trans SendTransferTXParams) (SendTransactionResultType, error) {
	res := SendTransactionResultType{}
	err := store.execTX(func(q *Queries) error {
		ctx := context.Background()

		transfer, err := q.CreateTransfer(ctx, CreateTransferParams(trans))
		if err != nil {
			return err
		}

		arg2 := UpdateAccountParams{
			ID:     trans.FromAccountID,
			Amount: -trans.Amount,
		}
		fromAccount, err := q.UpdateAccount(ctx, arg2)
		if err != nil {
			return err
		}

		arg3 := UpdateAccountParams{
			ID:     trans.ToAccountID,
			Amount: trans.Amount,
		}
		toAccount, err := q.UpdateAccount(ctx, arg3)
		if err != nil {
			return err
		}
		res.Transfer = transfer
		res.FromAccount = fromAccount
		res.ToAccount = toAccount
		return nil
	})

	return res, err
}
