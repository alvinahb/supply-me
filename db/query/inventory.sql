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

-- name: DeleteInventory :exec
DELETE FROM inventories where id = $1;
