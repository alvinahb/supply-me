package db

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)

	company1 := CreateRandomCompany(t)
	company2 := CreateRandomCompany(t)
	product := CreateRandomProduct(t)

	inventory1, err := testQueries.GetCompanyProductInventory(context.Background(),
		GetCompanyProductInventoryParams{
			CompanyID: company1.ID,
			ProductID: product.ID,
		})
	if err != nil {
		if err == sql.ErrNoRows {
			inventory1, err = testQueries.CreateInventory(context.Background(), CreateInventoryParams{
				CompanyID:       company1.ID,
				ProductID:       product.ID,
				AmountAvailable: int32(40),
			})
		}
	}
	require.NoError(t, err)
	require.NotEmpty(t, inventory1)

	inventory2, err := testQueries.GetCompanyProductInventory(context.Background(),
		GetCompanyProductInventoryParams{
			CompanyID: company2.ID,
			ProductID: product.ID,
		})
	if err != nil {
		if err == sql.ErrNoRows {
			inventory2, err = testQueries.CreateInventory(context.Background(), CreateInventoryParams{
				CompanyID:       company2.ID,
				ProductID:       product.ID,
				AmountAvailable: int32(0),
			})
		}
	}
	require.NoError(t, err)
	require.NotEmpty(t, inventory2)

	// Run n concurrent order transactions
	n := 5
	amount := int32(5)

	errs := make(chan error)
	results := make(chan OrderTxResult)

	for i := 0; i < n; i++ {
		go func() {
			result, err := store.OrderTx(context.Background(), OrderTxParams{
				FromCompanyID: company1.ID,
				ToCompanyID:   company2.ID,
				ProductID:     product.ID,
				Amount:        amount,
			})

			errs <- err
			results <- result
		}()
	}

	// Check results
	existed := make(map[int]bool)

	for i := 0; i < n; i++ {
		err := <-errs
		fmt.Println(err)
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		// Check order
		order := result.Order
		require.NotEmpty(t, order)
		require.Equal(t, company1.ID, order.FromCompanyID)
		require.Equal(t, company2.ID, order.ToCompanyID)
		require.Equal(t, amount, order.Amount)
		require.NotZero(t, order.ID)
		require.NotZero(t, order.CreatedAt)

		_, err = store.GetOrder(context.Background(), order.ID)
		require.NoError(t, err)

		// Check entries
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, company1.ID, fromEntry.CompanyID)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, company2.ID, toEntry.CompanyID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		// Check inventories
		fromInventory := result.FromInventory
		require.NotEmpty(t, fromInventory)
		require.Equal(t, company1.ID, fromInventory.CompanyID)

		toInventory := result.ToInventory
		require.NotEmpty(t, toInventory)
		require.Equal(t, company2.ID, toInventory.CompanyID)

		// Check inventories' amount available
		diff1 := inventory1.AmountAvailable - fromInventory.AmountAvailable
		diff2 := toInventory.AmountAvailable - inventory2.AmountAvailable
		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1%amount == 0)

		k := int(diff1 / amount)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true
	}

	// Check the final updated inventories
	newArgs := GetCompanyProductInventoryParams{
		CompanyID: company1.ID,
		ProductID: product.ID,
	}
	updatedInventory1, err := testQueries.GetCompanyProductInventory(
		context.Background(), newArgs)
	require.NoError(t, err)

	newArgs = GetCompanyProductInventoryParams{
		CompanyID: company2.ID,
		ProductID: product.ID,
	}
	updatedInventory2, err := testQueries.GetCompanyProductInventory(
		context.Background(), newArgs)
	require.NoError(t, err)

	require.Equal(t, inventory1.AmountAvailable-int32(n)*amount,
		updatedInventory1.AmountAvailable)
	require.Equal(t, inventory2.AmountAvailable+int32(n)*amount,
		updatedInventory2.AmountAvailable)
}
