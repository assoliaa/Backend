package db


import (
	"Backend/db/utils"
	"context"
	"testing"
	"time"
	"github.com/stretchr/testify/assert"
)

func createRandomUser(t *testing.T)User{
	hashedPasswod, err := utils.HashPassword(utils.RandomString(6))
	assert.NoError(t, err)
	params :=InsertUserParams{
		Username: utils.RandomOwner(), 
		HashPassword: hashedPasswod,
		FullName: utils.RandomOwner(),
		Email: utils.RandomEmail(),
	}
	user, err := testQueries.InsertUser(context.Background(), params)
    assert.NoError(t, err)
	assert.NotEmpty(t, user)
	
	assert.Equal(t, params.Username, user.Username)
	assert.Equal(t, params.HashPassword, user.HashPassword)
	assert.Equal(t, params.FullName, user.FullName)
	assert.Equal(t, params.Email, user.Email)

	assert.NotZero(t, user.CreatedAt)
	assert.True(t, user.PasswordChangedAt.IsZero())
	return user
}
func TestCreateUser(t *testing.T){
    createRandomUser(t)
}

func TestGetuser(t *testing.T){
	user1:= createRandomUser(t)
	user2, err := testQueries.GetUser(context.Background(), user1.Username)

	assert.NoError(t, err)
	assert.NotEmpty(t, user2)
	

	assert.Equal(t, user1.Username, user2.Username)
	assert.Equal(t, user1.HashPassword, user2.HashPassword)
	assert.Equal(t, user1.FullName, user2.FullName)
	assert.Equal(t, user1.Email, user2.Email)

	assert.WithinDuration(t, user1.PasswordChangedAt, user2.PasswordChangedAt, time.Second)
	assert.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)
}

