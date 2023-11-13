-- name: GetCoin :one
SELECT * FROM coin
WHERE name = $1;

-- name: CreateCoin :one
INSERT INTO coin (
  name
) VALUES (
  $1
)
RETURNING *;