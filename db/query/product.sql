-- name: CreateProduct :one
INSERT INTO products (
    product_name, description
) VALUES (
    $1, $2
) RETURNING *;

-- name: GetProduct :one
SELECT * FROM products WHERE id = $1 LIMIT 1;

-- name: ListProducts :many
SELECT * FROM products ORDER BY id LIMIT $1 OFFSET $2;

-- name: UpdateProduct :one
UPDATE products SET product_name = $1, description = $2 WHERE id = $3
RETURNING *;

-- name: DeleteProduct :exec
DELETE FROM products where id = $1;
