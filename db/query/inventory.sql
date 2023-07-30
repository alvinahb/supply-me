-- name: CreateInventory :one
INSERT INTO inventories (
    company_id, product_id, amount_available
) VALUES (
    $1, $2, $3
) RETURNING *;

-- name: GetInventory :one
SELECT * FROM inventories WHERE id = $1 LIMIT 1;

-- name: ListInventories :many
SELECT * FROM inventories ORDER BY id LIMIT $1 OFFSET $2;

-- name: UpdateInventory :one
UPDATE inventories SET company_id = $1, product_id = $2, amount_available = $3
WHERE id = $4
RETURNING *;

-- name: AddInventoryAmount :one
UPDATE inventories
SET amount_available = amount_available + sqlc.arg(amount)
WHERE id = sqlc.arg(id) RETURNING *;

-- name: DeleteInventory :exec
DELETE FROM inventories where id = $1;

-- name: GetCompanyProductInventory :one
SELECT * FROM inventories WHERE company_id = $1 AND product_id = $2 LIMIT 1;

-- name: ListCompanyInventories :many
SELECT * FROM inventories WHERE company_id = $1 ORDER BY id LIMIT $2 OFFSET $3;
