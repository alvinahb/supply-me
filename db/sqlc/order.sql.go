// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.19.1
// source: order.sql

package db

import (
	"context"
)

const createOrder = `-- name: CreateOrder :one
INSERT INTO orders (
    from_company_id, to_company_id, product_id, amount
) VALUES (
    $1, $2, $3, $4
) RETURNING id, from_company_id, to_company_id, product_id, amount, created_at
`

type CreateOrderParams struct {
	FromCompanyID int64 `json:"from_company_id"`
	ToCompanyID   int64 `json:"to_company_id"`
	ProductID     int64 `json:"product_id"`
	Amount        int32 `json:"amount"`
}

func (q *Queries) CreateOrder(ctx context.Context, arg CreateOrderParams) (Order, error) {
	row := q.db.QueryRowContext(ctx, createOrder,
		arg.FromCompanyID,
		arg.ToCompanyID,
		arg.ProductID,
		arg.Amount,
	)
	var i Order
	err := row.Scan(
		&i.ID,
		&i.FromCompanyID,
		&i.ToCompanyID,
		&i.ProductID,
		&i.Amount,
		&i.CreatedAt,
	)
	return i, err
}

const deleteOrder = `-- name: DeleteOrder :exec
DELETE FROM orders where id = $1
`

func (q *Queries) DeleteOrder(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteOrder, id)
	return err
}

const getOrder = `-- name: GetOrder :one
SELECT id, from_company_id, to_company_id, product_id, amount, created_at FROM orders WHERE id = $1 LIMIT 1
`

func (q *Queries) GetOrder(ctx context.Context, id int64) (Order, error) {
	row := q.db.QueryRowContext(ctx, getOrder, id)
	var i Order
	err := row.Scan(
		&i.ID,
		&i.FromCompanyID,
		&i.ToCompanyID,
		&i.ProductID,
		&i.Amount,
		&i.CreatedAt,
	)
	return i, err
}

const listOrders = `-- name: ListOrders :many
SELECT id, from_company_id, to_company_id, product_id, amount, created_at FROM orders ORDER BY id LIMIT $1 OFFSET $2
`

type ListOrdersParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) ListOrders(ctx context.Context, arg ListOrdersParams) ([]Order, error) {
	rows, err := q.db.QueryContext(ctx, listOrders, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Order{}
	for rows.Next() {
		var i Order
		if err := rows.Scan(
			&i.ID,
			&i.FromCompanyID,
			&i.ToCompanyID,
			&i.ProductID,
			&i.Amount,
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

const updateOrder = `-- name: UpdateOrder :one
UPDATE orders
SET from_company_id = $1, to_company_id = $2, product_id = $3, amount = $4
WHERE id = $5
RETURNING id, from_company_id, to_company_id, product_id, amount, created_at
`

type UpdateOrderParams struct {
	FromCompanyID int64 `json:"from_company_id"`
	ToCompanyID   int64 `json:"to_company_id"`
	ProductID     int64 `json:"product_id"`
	Amount        int32 `json:"amount"`
	ID            int64 `json:"id"`
}

func (q *Queries) UpdateOrder(ctx context.Context, arg UpdateOrderParams) (Order, error) {
	row := q.db.QueryRowContext(ctx, updateOrder,
		arg.FromCompanyID,
		arg.ToCompanyID,
		arg.ProductID,
		arg.Amount,
		arg.ID,
	)
	var i Order
	err := row.Scan(
		&i.ID,
		&i.FromCompanyID,
		&i.ToCompanyID,
		&i.ProductID,
		&i.Amount,
		&i.CreatedAt,
	)
	return i, err
}
