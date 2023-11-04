-- name: GetTransfers :many
SELECT * FROM transfers
WHERE from_account_id = $1 
LIMIT 100;

-- name: CreateTransfer :one
INSERT INTO transfers (
  from_account_id, to_account_id, amount, coin
) VALUES (
  $1, $2, $3, $4
)
RETURNING *;