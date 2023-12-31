// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.19.1
// source: product.sql

package db

import (
	"context"
)

const createProduct = `-- name: CreateProduct :one
INSERT INTO products (
    product_name, description
) VALUES (
    $1, $2
) RETURNING id, product_name, description, created_at
`

type CreateProductParams struct {
	ProductName string `json:"product_name"`
	Description string `json:"description"`
}

func (q *Queries) CreateProduct(ctx context.Context, arg CreateProductParams) (Product, error) {
	row := q.db.QueryRowContext(ctx, createProduct, arg.ProductName, arg.Description)
	var i Product
	err := row.Scan(
		&i.ID,
		&i.ProductName,
		&i.Description,
		&i.CreatedAt,
	)
	return i, err
}

const deleteProduct = `-- name: DeleteProduct :exec
DELETE FROM products where id = $1
`

func (q *Queries) DeleteProduct(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteProduct, id)
	return err
}

const getProduct = `-- name: GetProduct :one
SELECT id, product_name, description, created_at FROM products WHERE id = $1 LIMIT 1
`

func (q *Queries) GetProduct(ctx context.Context, id int64) (Product, error) {
	row := q.db.QueryRowContext(ctx, getProduct, id)
	var i Product
	err := row.Scan(
		&i.ID,
		&i.ProductName,
		&i.Description,
		&i.CreatedAt,
	)
	return i, err
}

const listProducts = `-- name: ListProducts :many
SELECT id, product_name, description, created_at FROM products ORDER BY id LIMIT $1 OFFSET $2
`

type ListProductsParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) ListProducts(ctx context.Context, arg ListProductsParams) ([]Product, error) {
	rows, err := q.db.QueryContext(ctx, listProducts, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Product{}
	for rows.Next() {
		var i Product
		if err := rows.Scan(
			&i.ID,
			&i.ProductName,
			&i.Description,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateProduct = `-- name: UpdateProduct :one
UPDATE products SET product_name = $1, description = $2 WHERE id = $3
RETURNING id, product_name, description, created_at
`

type UpdateProductParams struct {
	ProductName string `json:"product_name"`
	Description string `json:"description"`
	ID          int64  `json:"id"`
}

func (q *Queries) UpdateProduct(ctx context.Context, arg UpdateProductParams) (Product, error) {
	row := q.db.QueryRowContext(ctx, updateProduct, arg.ProductName, arg.Description, arg.ID)
	var i Product
	err := row.Scan(
		&i.ID,
		&i.ProductName,
		&i.Description,
		&i.CreatedAt,
	)
	return i, err
}
