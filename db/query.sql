-- name: GetAccount :one
SELECT lim, balance
FROM accounts
WHERE id = $1
LIMIT 1;

-- name: UpdateAccountBalance :exec
UPDATE accounts
SET balance = $2
WHERE id = $1;

-- name: GetTransactions :many
SELECT amount, type, description, date
FROM transactions
WHERE account_id = $1
ORDER BY id DESC
LIMIT 10;

-- name: CreateTransaction :exec
INSERT INTO transactions (account_id, amount, type, description)
VALUES ($1, $2, $3, $4);

-- name: InsertTransaction :exec
CALL insert_transaction(sqlc.arg('p_account_id'), sqlc.arg('p_amount'), sqlc.arg('p_type'), sqlc.arg('p_description'));