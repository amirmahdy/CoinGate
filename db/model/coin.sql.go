// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.22.0
// source: coin.sql

package db

import (
	"context"
)

const createCoin = `-- name: CreateCoin :one
INSERT INTO coin (
  name
) VALUES (
  $1
)
RETURNING name, created_at
`

func (q *Queries) CreateCoin(ctx context.Context, name string) (Coin, error) {
	row := q.queryRow(ctx, q.createCoinStmt, createCoin, name)
	var i Coin
	err := row.Scan(&i.Name, &i.CreatedAt)
	return i, err
}

const getCoin = `-- name: GetCoin :one
SELECT name, created_at FROM coin
WHERE name = $1
`

func (q *Queries) GetCoin(ctx context.Context, name string) (Coin, error) {
	row := q.queryRow(ctx, q.getCoinStmt, getCoin, name)
	var i Coin
	err := row.Scan(&i.Name, &i.CreatedAt)
	return i, err
}
