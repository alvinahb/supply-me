// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.19.1
// source: inventory.sql

package db

import (
	"context"
)

const createInventory = `-- name: CreateInventory :one
INSERT INTO inventories (
    company_id, product_id, amount_available
) VALUES (
    $1, $2, $3
) RETURNING id, company_id, product_id, amount_available, created_at, updated_at
`

type CreateInventoryParams struct {
	CompanyID       int64
	ProductID       int64
	AmountAvailable int32
}

func (q *Queries) CreateInventory(ctx context.Context, arg CreateInventoryParams) (Inventory, error) {
	row := q.db.QueryRowContext(ctx, createInventory, arg.CompanyID, arg.ProductID, arg.AmountAvailable)
	var i Inventory
	err := row.Scan(
		&i.ID,
		&i.CompanyID,
		&i.ProductID,
		&i.AmountAvailable,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deleteInventory = `-- name: DeleteInventory :exec
DELETE FROM inventories where id = $1
`

func (q *Queries) DeleteInventory(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteInventory, id)
	return err
}

const getInventory = `-- name: GetInventory :one
SELECT id, company_id, product_id, amount_available, created_at, updated_at FROM inventories WHERE id = $1 LIMIT 1
`

func (q *Queries) GetInventory(ctx context.Context, id int64) (Inventory, error) {
	row := q.db.QueryRowContext(ctx, getInventory, id)
	var i Inventory
	err := row.Scan(
		&i.ID,
		&i.CompanyID,
		&i.ProductID,
		&i.AmountAvailable,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const listInventories = `-- name: ListInventories :many
SELECT id, company_id, product_id, amount_available, created_at, updated_at FROM inventories ORDER BY id LIMIT $1 OFFSET $2
`

type ListInventoriesParams struct {
	Limit  int32
	Offset int32
}

func (q *Queries) ListInventories(ctx context.Context, arg ListInventoriesParams) ([]Inventory, error) {
	rows, err := q.db.QueryContext(ctx, listInventories, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Inventory
	for rows.Next() {
		var i Inventory
		if err := rows.Scan(
			&i.ID,
			&i.CompanyID,
			&i.ProductID,
			&i.AmountAvailable,
			&i.CreatedAt,
			&i.UpdatedAt,
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

const updateInventory = `-- name: UpdateInventory :one
UPDATE inventories SET company_id = $1, product_id = $2, amount_available = $3
WHERE id = $4
RETURNING id, company_id, product_id, amount_available, created_at, updated_at
`

type UpdateInventoryParams struct {
	CompanyID       int64
	ProductID       int64
	AmountAvailable int32
	ID              int64
}

func (q *Queries) UpdateInventory(ctx context.Context, arg UpdateInventoryParams) (Inventory, error) {
	row := q.db.QueryRowContext(ctx, updateInventory,
		arg.CompanyID,
		arg.ProductID,
		arg.AmountAvailable,
		arg.ID,
	)
	var i Inventory
	err := row.Scan(
		&i.ID,
		&i.CompanyID,
		&i.ProductID,
		&i.AmountAvailable,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}