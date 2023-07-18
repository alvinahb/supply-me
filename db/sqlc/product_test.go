package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/alvinahb/supply-me/util"
	"github.com/stretchr/testify/require"
)

func CreateRandomProduct(t *testing.T) Product {
	args := CreateProductParams{
		ProductName: util.RandomString(15),
		Origin:      util.RandomString(10),
	}

	product, err := testQueries.CreateProduct(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, product)

	require.Equal(t, args.ProductName, product.ProductName)
	require.Equal(t, args.Origin, product.Origin)

	require.NotZero(t, product.ID)
	require.NotZero(t, product.CreatedAt)

	return product
}

func TestCreateProduct(t *testing.T) {
	CreateRandomProduct(t)
}

func TestGetProduct(t *testing.T) {
	product1 := CreateRandomProduct(t)
	product2, err := testQueries.GetProduct(context.Background(), product1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, product2)

	require.Equal(t, product1.ID, product2.ID)
	require.Equal(t, product1.ProductName, product2.ProductName)
	require.Equal(t, product1.Origin, product2.Origin)
	require.WithinDuration(t, product1.CreatedAt, product2.CreatedAt, time.Second)
}

func TestUpdateProduct(t *testing.T) {
	product1 := CreateRandomProduct(t)

	args := UpdateProductParams{
		ProductName: util.RandomString(15),
		Origin:      util.RandomString(10),
		ID:          product1.ID,
	}

	product2, err := testQueries.UpdateProduct(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, product2)

	require.Equal(t, product1.ID, product2.ID)
	require.Equal(t, args.ProductName, product2.ProductName)
	require.Equal(t, args.Origin, product2.Origin)
	require.WithinDuration(t, product1.CreatedAt, product2.CreatedAt, time.Second)
}

func TestDeleteProduct(t *testing.T) {
	product1 := CreateRandomProduct(t)
	err := testQueries.DeleteProduct(context.Background(), product1.ID)
	require.NoError(t, err)

	product2, err := testQueries.GetProduct(context.Background(), product1.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, product2)
}

func TestListProducts(t *testing.T) {
	for i := 0; i < 10; i++ {
		CreateRandomProduct(t)
	}

	args := ListProductsParams{
		Limit:  5,
		Offset: 5,
	}

	products, err := testQueries.ListProducts(context.Background(), args)
	require.NoError(t, err)
	require.Len(t, products, 5)

	for _, product := range products {
		require.NotEmpty(t, product)
	}
}
