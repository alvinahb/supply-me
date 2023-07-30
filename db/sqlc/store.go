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
// It creates an order record, add company entries and update companies
// inventories within a single database transaction
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

		// Get provider inventory for the product
		result.FromInventory, err = q.GetCompanyProductInventory(ctx,
			GetCompanyProductInventoryParams{
				CompanyID: args.FromCompanyID,
				ProductID: args.ProductID,
			})
		if err != nil {
			return err
		}

		// Get restaurant inventory for the product
		// OR create it if not existing yet
		result.ToInventory, err = q.GetCompanyProductInventory(ctx,
			GetCompanyProductInventoryParams{
				CompanyID: args.ToCompanyID,
				ProductID: args.ProductID,
			})
		if err != nil {
			if err == sql.ErrNoRows {
				result.ToInventory, err = q.CreateInventory(ctx,
					CreateInventoryParams{
						CompanyID:       args.ToCompanyID,
						ProductID:       args.ProductID,
						AmountAvailable: int32(0),
					})
				if err != nil {
					return err
				}
			} else {
				return err
			}
		}

		// Ordering inventories transactions to avoid deadlocks
		if result.FromInventory.ID < result.ToInventory.ID {
			result.FromInventory, result.ToInventory, err = transferProduct(
				ctx, q, result.FromInventory.ID, -args.Amount, result.ToInventory.ID, args.Amount)
			if err != nil {
				return err
			}
		} else {
			result.ToInventory, result.FromInventory, err = transferProduct(
				ctx, q, result.ToInventory.ID, args.Amount, result.FromInventory.ID, -args.Amount)
			if err != nil {
				return err
			}
		}

		return nil
	})

	return result, err
}

func transferProduct(
	ctx context.Context,
	q *Queries,
	inventoryID1 int64,
	amount1 int32,
	inventoryID2 int64,
	amount2 int32,
) (inventory1 Inventory, inventory2 Inventory, err error) {
	inventory1, err = q.AddInventoryAmount(ctx, AddInventoryAmountParams{
		Amount: amount1,
		ID:     inventoryID1,
	})
	if err != nil {
		return
	}

	inventory2, err = q.AddInventoryAmount(ctx, AddInventoryAmountParams{
		Amount: amount2,
		ID:     inventoryID2,
	})
	if err != nil {
		return
	}

	return
}
