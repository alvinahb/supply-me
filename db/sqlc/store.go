package db

import (
	"context"
	"database/sql"
	"fmt"
)

// Store provides all functions to execute db queries and transactions
type Store struct {
	*Queries
	db *sql.DB
}

// NewStore creates a now Store
func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

// execTx executes a function within a database transaction
func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}

// OrderTxParams contains the input parameters of the order transaction
type OrderTxParams struct {
	FromCompanyID int64 `json:"from_company_id"`
	ToCompanyID   int64 `json:"to_company_id"`
	ProductID     int64 `json:"product_id"`
	Amount        int32 `json:"amount"`
}

// OrderTxResult is the result of the order transaction
type OrderTxResult struct {
	Order         Order     `json:"order"`
	FromInventory Inventory `json:"from_inventory"`
	ToInventory   Inventory `json:"to_Inventory"`
	FromEntry     Entry     `json:"from_entry"`
	ToEntry       Entry     `json:"to_entry"`
}

// OrderTx performs a product order from a provider to a restaurant.
// It creates an order record, add company entries and update companies inventories within a single database transaction
func (store *Store) OrderTx(ctx context.Context, args OrderTxParams) (OrderTxResult, error) {
	var result OrderTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		// Create an order record
		result.Order, err = q.CreateOrder(ctx, CreateOrderParams{
			FromCompanyID: args.FromCompanyID,
			ToCompanyID:   args.ToCompanyID,
			ProductID:     args.ProductID,
			Amount:        args.Amount,
		})
		if err != nil {
			return err
		}

		// Create an entry for the provider
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			CompanyID: args.FromCompanyID,
			ProductID: args.ProductID,
			Amount:    -args.Amount,
		})
		if err != nil {
			return err
		}

		// Create an entry for the restaurant
		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			CompanyID: args.ToCompanyID,
			ProductID: args.ProductID,
			Amount:    args.Amount,
		})
		if err != nil {
			return err
		}

		// Update provider inventory (should be existing)
		result.FromInventory, err = q.AddInventoryAmount(ctx, AddInventoryAmountParams{
			Amount:    -args.Amount,
			CompanyID: args.FromCompanyID,
			ProductID: args.ProductID,
		})
		if err != nil {
			return err
		}

		// Create restaurant inventory if not existing yet for this product
		// OR update restaurant inventory if already existing
		result.ToInventory, err = q.AddInventoryAmount(ctx, AddInventoryAmountParams{
			Amount:    args.Amount,
			CompanyID: args.ToCompanyID,
			ProductID: args.ProductID,
		})
		if err != nil {
			if err == sql.ErrNoRows {
				_, err = q.CreateInventory(ctx, CreateInventoryParams{
					CompanyID:       args.ToCompanyID,
					ProductID:       args.ProductID,
					AmountAvailable: int32(args.Amount),
				})
				if err != nil {
					return err
				}
			} else {
				return err
			}
		}

		return nil
	})

	return result, err
}
