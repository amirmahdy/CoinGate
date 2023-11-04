-- name: GetAccount :one
SELECT * FROM accounts
WHERE username LIKE $1 LIMIT 10;

-- name: CreateAccount :one
INSERT INTO accounts (
  username, balance, coin
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: UpdateAccount :one
UPDATE accounts SET
  balance = balance + sqlc.arg(amount)
WHERE id = $1
RETURNING *;