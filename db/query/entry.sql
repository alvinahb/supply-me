-- name: CreateEntry :one
INSERT INTO entries (
    company_id, product_id, amount
) VALUES (
    $1, $2, $3
) RETURNING *;

-- name: GetEntry :one
SELECT * FROM entries WHERE id = $1 LIMIT 1;

-- name: ListEntries :many
SELECT * FROM entries ORDER BY id LIMIT $1 OFFSET $2;

-- name: UpdateEntry :one
UPDATE entries SET company_id = $1, product_id = $2, amount = $3 WHERE id = $4
RETURNING *;

-- name: DeleteEntry :exec
DELETE FROM entries where id = $1;
