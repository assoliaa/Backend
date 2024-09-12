package db

import (
	"Backend/db/utils"
	"context"
	"database/sql"
	"testing"
	"time"
	"github.com/stretchr/testify/assert"
)

func createRandomAccount(t *testing.T)Account{
	user:=createRandomUser(t)
	arg :=InsertAccountParams{
		Owner: user.Username, 
		Balance: utils.RandomMoney(),
		Currency: utils.RandomCurrency(),
	}
	account, err := testQueries.InsertAccount(context.Background(), arg)
    assert.NoError(t, err)
	assert.NotEmpty(t, account)
	
	assert.Equal(t, arg.Owner, account.Owner)
	assert.Equal(t, arg.Balance, account.Balance)
	assert.Equal(t, arg.Currency, account.Currency)

	assert.NotZero(t, account.ID)
	return account
}

func TestCreateAccount(t *testing.T){
    createRandomAccount(t)
}

func TestGetAccount(t *testing.T){
	account1:= createRandomAccount(t)
	account2, err := testQueries.GetAccount(context.Background(), account1.ID)

	assert.NoError(t, err)
	assert.NotEmpty(t, account2)
	
	assert.Equal(t, account1.ID, account2.ID)
	assert.Equal(t, account1.Owner, account2.Owner)
	assert.Equal(t, account1.Balance, account2.Balance)
	assert.Equal(t, account1.Currency, account2.Currency)

	assert.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)
}

func TestUpdateAccount(t *testing.T){
	account1 := createRandomAccount(t)

	arg := UpdateAccountParams{
		ID: account1.ID,
		Balance : utils.RandomMoney(),
	}

	account2, err :=testQueries.UpdateAccount(context.Background(), arg)
    assert.NoError(t, err)
	assert.NotEmpty(t, account2)
	
	assert.Equal(t, account1.ID, account2.ID)
	assert.Equal(t, account1.Owner, account2.Owner)
	assert.Equal(t, arg.Balance, account2.Balance)
	assert.Equal(t, account1.Currency, account2.Currency)

	assert.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)

}

func TestDeleteAccount(t *testing.T){
	account1 := createRandomAccount(t)
	err:= testQueries.DeleteAccount(context.Background(), account1.ID)

	assert.NoError(t,err)

	account2, err := testQueries.GetAccount(context.Background(), account1.ID)
	
	assert.Error(t, err)
	assert.EqualError(t, err, sql.ErrNoRows.Error())	
	assert.Empty(t, account2)
}

func TestListAccounts(t *testing.T) {
	var lastAccount Account
	for i := 0; i < 10; i++ {
		lastAccount = createRandomAccount(t)
	}

	arg := GetAccountsParams{
		Owner:  lastAccount.Owner,
		Limit:  5,
		Offset: 0,
	}

	accounts, err := testQueries.GetAccounts(context.Background(), arg)
	assert.NoError(t, err)
	assert.NotEmpty(t, accounts)

	for _, account := range accounts {
		assert.NotEmpty(t, account)
		assert.Equal(t, lastAccount.Owner, account.Owner)
	}
}
	