-- name: CreateAccount :one
INSERT INTO accounts (
  owner,
  balance,
  currency
) VALUES (
  $1, $2, $3
) RETURNING *;

-- name: GetAccount :one
SELECT * FROM accounts
WHERE id = $1 LIMIT 1;

-- name: GetAccountForUpdate :one
-- for update -> to lock the rows returned as well as get a transaction lock on all the rows that reference the parent table (ensures that another transaction won't update they other key at the saem time breaking cosistency )
-- no key -> weaker lock that tell psql to allow lock even when held by another contraint as we are not updating the id
SELECT * FROM accounts
WHERE id = $1 LIMIT 1 FOR NO KEY UPDATE;

-- name: ListAccounts :many
SELECT * FROM accounts
WHERE owner = $1
ORDER BY id
LIMIT $2
OFFSET $3;

-- name: UpdateAccount :one
UPDATE accounts SET balance = $2 
WHERE id = $1 RETURNING *;

-- name: AddAccountBalance :one
UPDATE accounts SET balance = balance + sqlc.arg(amount) 
WHERE id = sqlc.arg(id) RETURNING *;

-- name: DeleteAccount :exec
DELETE FROM accounts WHERE id = $1;