-- name: CreateOrder :one
INSERT INTO orders (
    from_company_id, to_company_id, product_id, amount
) VALUES (
    $1, $2, $3, $4
) RETURNING *;

-- name: GetOrder :one
SELECT * FROM orders WHERE id = $1 LIMIT 1;

-- name: ListOrders :many
SELECT * FROM orders ORDER BY id LIMIT $1 OFFSET $2;

-- name: UpdateOrder :one
UPDATE orders
SET from_company_id = $1, to_company_id = $2, product_id = $3, amount = $4
WHERE id = $5
RETURNING *;

-- name: DeleteOrder :exec
DELETE FROM orders where id = $1;
