package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/alvinahb/supply-me/util"
	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) User {
	company := CreateRandomCompany(t)

	args := CreateUserParams{
		FirstName: util.RandomString(8),
		LastName:  util.RandomString(8),
		Email:     util.RandomString(20),
		Password:  util.RandomString(12),
		CompanyID: company.ID,
		Role:      util.RandomRole(),
	}

	user, err := testQueries.CreateUser(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, args.FirstName, user.FirstName)
	require.Equal(t, args.LastName, user.LastName)
	require.Equal(t, args.Email, user.Email)
	require.Equal(t, args.Password, user.Password)
	require.Equal(t, args.CompanyID, user.CompanyID)
	require.Equal(t, args.Role, user.Role)

	require.NotZero(t, user.ID)
	require.NotZero(t, user.CreatedAt)
	require.NotZero(t, user.UpdatedAt)

	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	user1 := createRandomUser(t)
	user2, err := testQueries.GetUser(context.Background(), user1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.ID, user2.ID)
	require.Equal(t, user1.FirstName, user2.FirstName)
	require.Equal(t, user1.LastName, user2.LastName)
	require.Equal(t, user1.Email, user2.Email)
	require.Equal(t, user1.Password, user2.Password)
	require.Equal(t, user1.CompanyID, user2.CompanyID)
	require.Equal(t, user1.Role, user2.Role)
	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)
	require.WithinDuration(t, user1.UpdatedAt, user2.UpdatedAt, time.Second)
}

func TestUpdateUser(t *testing.T) {
	user1 := createRandomUser(t)

	args := UpdateUserParams{
		FirstName: util.RandomString(8),
		LastName:  util.RandomString(8),
		Email:     util.RandomString(20),
		Password:  util.RandomString(12),
		CompanyID: user1.CompanyID,
		Role:      util.RandomRole(),
		ID:        user1.ID,
	}

	user2, err := testQueries.UpdateUser(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.ID, user2.ID)
	require.Equal(t, args.FirstName, user2.FirstName)
	require.Equal(t, args.LastName, user2.LastName)
	require.Equal(t, args.Email, user2.Email)
	require.Equal(t, args.Password, user2.Password)
	require.Equal(t, args.CompanyID, user2.CompanyID)
	require.Equal(t, args.Role, user2.Role)
	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)
}

func TestDeleteUser(t *testing.T) {
	user1 := createRandomUser(t)
	err := testQueries.DeleteUser(context.Background(), user1.ID)
	require.NoError(t, err)

	user2, err := testQueries.GetUser(context.Background(), user1.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, user2)
}

func TestListUsers(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomUser(t)
	}

	args := ListUsersParams{
		Limit:  5,
		Offset: 5,
	}

	users, err := testQueries.ListUsers(context.Background(), args)
	require.NoError(t, err)
	require.Len(t, users, 5)

	for _, user := range users {
		require.NotEmpty(t, user)
	}
}
