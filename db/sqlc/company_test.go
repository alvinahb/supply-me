package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/alvinahb/supply-me/util"
	"github.com/stretchr/testify/require"
)

func CreateRandomCompany(t *testing.T) Company {
	args := CreateCompanyParams{
		CompanyType: util.RandomCompanyType(),
		CompanyName: util.RandomString(20),
	}

	company, err := testQueries.CreateCompany(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, company)

	require.Equal(t, args.CompanyType, company.CompanyType)
	require.Equal(t, args.CompanyName, company.CompanyName)

	require.NotZero(t, company.ID)
	require.NotZero(t, company.CreatedAt)

	return company
}

func TestCreateCompany(t *testing.T) {
	CreateRandomCompany(t)
}

func TestGetCompany(t *testing.T) {
	company1 := CreateRandomCompany(t)
	company2, err := testQueries.GetCompany(context.Background(), company1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, company2)

	require.Equal(t, company1.ID, company2.ID)
	require.Equal(t, company1.CompanyType, company2.CompanyType)
	require.Equal(t, company1.CompanyName, company2.CompanyName)
	require.WithinDuration(t, company1.CreatedAt, company2.CreatedAt, time.Second)
}

func TestUpdateCompany(t *testing.T) {
	company1 := CreateRandomCompany(t)

	args := UpdateCompanyParams{
		ID:          company1.ID,
		CompanyName: util.RandomString(20),
	}

	company2, err := testQueries.UpdateCompany(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, company2)

	require.Equal(t, company1.ID, company2.ID)
	require.Equal(t, company1.CompanyType, company2.CompanyType)
	require.Equal(t, args.CompanyName, company2.CompanyName)
	require.WithinDuration(t, company1.CreatedAt, company2.CreatedAt, time.Second)
}

func TestDeleteCompany(t *testing.T) {
	company1 := CreateRandomCompany(t)
	err := testQueries.DeleteCompany(context.Background(), company1.ID)
	require.NoError(t, err)

	company2, err := testQueries.GetCompany(context.Background(), company1.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, company2)
}

func TestListCompanies(t *testing.T) {
	for i := 0; i < 10; i++ {
		CreateRandomCompany(t)
	}

	args := ListCompaniesParams{
		Limit:  5,
		Offset: 5,
	}

	companies, err := testQueries.ListCompanies(context.Background(), args)
	require.NoError(t, err)
	require.Len(t, companies, 5)

	for _, company := range companies {
		require.NotEmpty(t, company)
	}
}
