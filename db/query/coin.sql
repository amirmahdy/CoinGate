-- name: GetCoin :many
SELECT * FROM coin
WHERE name LIKE $1 
LIMIT 10;

-- name: CreateCoin :one
INSERT INTO coin (
  name
) VALUES (
  $1
)
RETURNING *;