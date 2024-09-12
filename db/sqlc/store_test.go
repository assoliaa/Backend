package db

import (
	"Backend/db/utils"
	"context"
	"fmt"
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestTrsnasferTx(t *testing.T) {
	store :=NewStore(testDB)

	account1:=createRandomAccount(t)
	account2:=createRandomAccount(t)

	//run concurrent go routines

	n:=5
	amount:=utils.RandomMoney()

	errs :=make(chan error)
	results := make(chan TransferTxResult)
	
	for i:=0; i<n; i++{
		txName:= fmt.Sprintf("tx %d", i+1)
		go func(){
			ctx :=context.WithValue(context.Background(),txKey, txName)
			result, err:= store.TransferTx(ctx, TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID: account2.ID,
				Amount:amount ,
			})
			errs<-err
			results<-result


		}()
	}
	exists :=make(map[int]bool)
	for i:=0; i<n; i++{
		err:=<-errs
		assert.NoError(t, err)

		result :=<-results
		assert.NotEmpty(t, result)

		transfer:=result.Transfer
		assert.NotEmpty(t,transfer)
		assert.Equal(t, account1.ID, transfer.FromAccountID)
		assert.Equal(t, account2.ID, transfer.ToAccountID)
		assert.Equal(t, amount, transfer.Amount)
		assert.NotZero(t, transfer.ID)
		assert.NotZero(t, transfer.CreatedAt)
		
       _, err = store.GetTransfer(context.Background(), transfer.ID)
       assert.NoError(t,err)

	   fromEntry := result.FromEntry
	   assert.NotEmpty(t,fromEntry)
	   assert.Equal(t, account1.ID, fromEntry.AccountID)
	   assert.Equal(t, -amount, fromEntry.Amount)
       assert.NotZero(t, fromEntry.ID)
	   assert.NotZero(t, fromEntry.CreatedAt)
	  
	   _, err = store.GetEntry(context.Background(), fromEntry.ID)
       assert.NoError(t,err)

	   
	   toEntry := result.ToEntry
	   assert.NotEmpty(t,toEntry)
	   assert.Equal(t, account2.ID, toEntry.AccountID)
	   assert.Equal(t, amount, toEntry.Amount)
       assert.NotZero(t, toEntry.ID)
	   assert.NotZero(t, toEntry.CreatedAt)
	  
	   _, err = store.GetEntry(context.Background(), fromEntry.ID)
       assert.NoError(t,err)
    
	//check accounts
	fromAccount :=result.FromAccount
	assert.NotEmpty(t, fromAccount)
	assert.Equal(t, account1.ID, fromAccount.ID)//туу

	toAccount := result.ToAccount
	assert.NotEmpty(t, toAccount)//nhere
	assert.Equal(t, account2.ID, toAccount.ID)//here

	// провеить разницу
	diff :=account1.Balance - fromAccount.Balance
	diff2 :=toAccount.Balance -account2.Balance

	assert.Equal(t, diff, diff2)//here
	assert.True(t, diff>0)//here
	assert.True(t, diff%amount==0)// here

	k := int(diff/amount)
	assert.True(t, k>=1 && k<=n)//here
    assert.NotContains(t, exists, k)
	exists[k] =true
	}
   
	updatedAccount1, err := testQueries.GetAccount(context.Background(), account1.ID)
	assert.NoError(t, err)

	updatedAccount2, err := testQueries.GetAccount(context.Background(), account2.ID)
	assert.NoError(t, err)

	assert.Equal(t, account1.Balance- int64(n)*amount,updatedAccount1.Balance)
	
	assert.Equal(t, account2.Balance+int64(n)*amount,updatedAccount2.Balance)

}

func TestTrsnasferTxDeadlock(t *testing.T) {
	store :=NewStore(testDB)

	account1:=createRandomAccount(t)
	account2:=createRandomAccount(t)

	

	n:=10
	amount:=utils.RandomMoney()

	errs := make(chan error)
	
	for i:=0; i<n; i++{
		fromAccountID := account1.ID
		toAccountID :=account2.ID

		if i%2 ==1{
			fromAccountID = account2.ID
			toAccountID =account1.ID
		}
	    txName:= fmt.Sprintf("tx %d", i+1)
		go func(){
			ctx :=context.WithValue(context.Background(),txKey, txName)
			_, err:= store.TransferTx(ctx, TransferTxParams{
				FromAccountID: fromAccountID,
				ToAccountID: toAccountID,
				Amount:amount ,
			})
			errs<-err
			
		}()
	}
	for i:=0; i<n; i++{
		err :=<-errs
		assert.NoError(t, err)
	}
   
	updatedAccount1, err := testQueries.GetAccount(context.Background(), account1.ID)
	assert.NoError(t, err)

	updatedAccount2, err := testQueries.GetAccount(context.Background(), account2.ID)
	assert.NoError(t, err)

	assert.Equal(t, account1.Balance,updatedAccount1.Balance)
	
	assert.Equal(t, account2.Balance,updatedAccount2.Balance)

}