package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/alvinahb/supply-me/util"
	"github.com/stretchr/testify/require"
)

func createRandomOrder(t *testing.T) Order {
	company1 := CreateRandomCompany(t)
	company2 := CreateRandomCompany(t)
	product := CreateRandomProduct(t)

	args := CreateOrderParams{
		FromCompanyID: company1.ID,
		ToCompanyID:   company2.ID,
		ProductID:     product.ID,
		Amount:        int32(util.RandomInt(1, 100)),
	}

	order, err := testQueries.CreateOrder(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, order)

	require.Equal(t, args.FromCompanyID, order.FromCompanyID)
	require.Equal(t, args.ToCompanyID, order.ToCompanyID)
	require.Equal(t, args.ProductID, order.ProductID)
	require.Equal(t, args.Amount, order.Amount)

	require.NotZero(t, order.ID)
	require.NotZero(t, order.CreatedAt)

	return order
}

func TestCreateOrder(t *testing.T) {
	createRandomOrder(t)
}

func TestGetOrder(t *testing.T) {
	order1 := createRandomOrder(t)
	order2, err := testQueries.GetOrder(context.Background(), order1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, order2)

	require.Equal(t, order1.ID, order2.ID)
	require.Equal(t, order1.FromCompanyID, order2.FromCompanyID)
	require.Equal(t, order1.ToCompanyID, order2.ToCompanyID)
	require.Equal(t, order1.ProductID, order2.ProductID)
	require.Equal(t, order1.Amount, order2.Amount)
	require.WithinDuration(t, order1.CreatedAt, order2.CreatedAt, time.Second)
}

func TestUpdateOrder(t *testing.T) {
	order1 := createRandomOrder(t)

	args := UpdateOrderParams{
		FromCompanyID: order1.FromCompanyID,
		ToCompanyID:   order1.ToCompanyID,
		ProductID:     order1.ProductID,
		Amount:        int32(util.RandomInt(1, 100)),
		ID:            order1.ID,
	}

	order2, err := testQueries.UpdateOrder(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, order2)

	require.Equal(t, order1.ID, order2.ID)
	require.Equal(t, args.FromCompanyID, order2.FromCompanyID)
	require.Equal(t, args.ToCompanyID, order2.ToCompanyID)
	require.Equal(t, args.ProductID, order2.ProductID)
	require.Equal(t, args.Amount, order2.Amount)
	require.WithinDuration(t, order1.CreatedAt, order2.CreatedAt, time.Second)
}

func TestDeleteOrder(t *testing.T) {
	order1 := createRandomOrder(t)
	err := testQueries.DeleteOrder(context.Background(), order1.ID)
	require.NoError(t, err)

	order2, err := testQueries.GetOrder(context.Background(), order1.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, order2)
}

func TestListOrders(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomOrder(t)
	}

	args := ListOrdersParams{
		Limit:  5,
		Offset: 5,
	}

	orders, err := testQueries.ListOrders(context.Background(), args)
	require.NoError(t, err)
	require.Len(t, orders, 5)

	for _, order := range orders {
		require.NotEmpty(t, order)
	}
}
