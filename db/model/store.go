package db

import (
	"database/sql"
	"errors"
	"fmt"
)

type Store interface {
	SendTransferTX(trans SendTransferTXParams) (SendTransactionResultType, error)
	Querier
}

type SQLStore struct {
	db *sql.DB
	*Queries
}

var (
	ErrUniqueViolation = errors.New("unique_violation")
)

func NewSQLStore(db *sql.DB) Store {
	return &SQLStore{
		db:      db,
		Queries: New(db),
	}
}

func (store *SQLStore) execTX(fn func(*Queries) error) error {
	tx, err := store.db.Begin()
	if err != nil {
		return err
	}

	qTX := store.WithTx(tx)
	err = fn(qTX)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx error: %v, rb error: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}
