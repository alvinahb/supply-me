package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/alvinahb/supply-me/util"
	"github.com/stretchr/testify/require"
)

func createRandomInventory(t *testing.T) Inventory {
	company := CreateRandomCompany(t)
	product := CreateRandomProduct(t)

	args := CreateInventoryParams{
		CompanyID:       company.ID,
		ProductID:       product.ID,
		AmountAvailable: int32(util.RandomInt(0, 100)),
	}

	inventory, err := testQueries.CreateInventory(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, inventory)

	require.Equal(t, args.CompanyID, inventory.CompanyID)
	require.Equal(t, args.ProductID, inventory.ProductID)
	require.Equal(t, args.AmountAvailable, inventory.AmountAvailable)

	require.NotZero(t, inventory.ID)
	require.NotZero(t, inventory.CreatedAt)
	require.NotZero(t, inventory.UpdatedAt)

	return inventory
}

func TestCreateInventory(t *testing.T) {
	createRandomInventory(t)
}

func TestGetInventory(t *testing.T) {
	inventory1 := createRandomInventory(t)
	inventory2, err := testQueries.GetInventory(context.Background(), inventory1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, inventory2)

	require.Equal(t, inventory1.ID, inventory2.ID)
	require.Equal(t, inventory1.CompanyID, inventory2.CompanyID)
	require.Equal(t, inventory1.ProductID, inventory2.ProductID)
	require.Equal(t, inventory1.AmountAvailable, inventory2.AmountAvailable)
	require.WithinDuration(t, inventory1.CreatedAt, inventory2.CreatedAt, time.Second)
	require.WithinDuration(t, inventory1.UpdatedAt, inventory2.UpdatedAt, time.Second)
}

func TestUpdateInventory(t *testing.T) {
	inventory1 := createRandomInventory(t)

	args := UpdateInventoryParams{
		CompanyID:       inventory1.CompanyID,
		ProductID:       inventory1.ProductID,
		AmountAvailable: int32(util.RandomInt(0, 100)),
		ID:              inventory1.ID,
	}

	inventory2, err := testQueries.UpdateInventory(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, inventory2)

	require.Equal(t, inventory1.ID, inventory2.ID)
	require.Equal(t, args.CompanyID, inventory2.CompanyID)
	require.Equal(t, args.ProductID, inventory2.ProductID)
	require.Equal(t, args.AmountAvailable, inventory2.AmountAvailable)
	require.WithinDuration(t, inventory1.CreatedAt, inventory2.CreatedAt, time.Second)
}

func TestDeleteInventory(t *testing.T) {
	inventory1 := createRandomInventory(t)
	err := testQueries.DeleteInventory(context.Background(), inventory1.ID)
	require.NoError(t, err)

	inventory2, err := testQueries.GetInventory(context.Background(), inventory1.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, inventory2)
}

func TestListInventories(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomInventory(t)
	}

	args := ListInventoriesParams{
		Limit:  5,
		Offset: 5,
	}

	inventories, err := testQueries.ListInventories(context.Background(), args)
	require.NoError(t, err)
	require.Len(t, inventories, 5)

	for _, inventory := range inventories {
		require.NotEmpty(t, inventory)
	}
}

func TestGetCompanyProductInventory(t *testing.T) {
	inventory1 := createRandomInventory(t)

	args := GetCompanyProductInventoryParams{
		CompanyID: inventory1.CompanyID,
		ProductID: inventory1.ProductID,
	}

	inventory2, err := testQueries.GetCompanyProductInventory(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, inventory2)

	require.Equal(t, inventory1.ID, inventory2.ID)
	require.Equal(t, inventory1.CompanyID, inventory2.CompanyID)
	require.Equal(t, inventory1.ProductID, inventory2.ProductID)
	require.Equal(t, inventory1.AmountAvailable, inventory2.AmountAvailable)
	require.WithinDuration(t, inventory1.CreatedAt, inventory2.CreatedAt, time.Second)
	require.WithinDuration(t, inventory1.UpdatedAt, inventory2.UpdatedAt, time.Second)
}
