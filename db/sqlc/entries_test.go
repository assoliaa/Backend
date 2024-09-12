package db

import (
	"Backend/db/utils"
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"
	"github.com/stretchr/testify/assert"
)



func GetRandomAccountId(ctx context.Context) (int64, error) {
    var accountID int64
    err := testQueries.db.QueryRowContext(ctx, `
        SELECT id 
        FROM accounts 
        ORDER BY RANDOM() 
        LIMIT 1
    `).Scan(&accountID)

    if err != nil {
        return 0, fmt.Errorf("could not get random account ID: %w", err)
    }

    return accountID, nil
}

func getRandomEntryID(ctx context.Context)(int64, error){
	var entryID int64
	accountId, err :=GetRandomAccountId(context.Background())
	if err !=nil{
		return 0, fmt.Errorf("could not get accounts: %w", err)
	}
	err = testQueries.db.QueryRowContext(ctx, `
        SELECT id 
        FROM entries 
        WHERE account_id=$1
    `, accountId).Scan(&entryID)
	if err!=nil{
		return 0, fmt.Errorf("could not get accounts: %w", err)
	}
	return entryID, nil
}


func createRandomEntry(t *testing.T)Entry{
	id, err :=GetRandomAccountId(context.Background())
	if err !=nil{
		fmt.Errorf("could not get accounts: %w", err)
	}
	params:= InsertEntryParams{
		AccountID: id,
		Amount:utils.RandomMoney(),
	    
	}
	entry, err := testQueries.InsertEntry(context.Background(), params) //here

	assert.NoError(t, err)
	assert.NotEmpty(t, entry)

	assert.Equal(t, params.AccountID, entry.AccountID)
	assert.Equal(t, params.Amount, entry.Amount)

	assert.NotZero(t, entry)
	return entry
}


func TestCreateEntry(t *testing.T){
	createRandomEntry(t)
}

func TestGetEntry(t *testing.T){
    entry1 :=createRandomEntry(t)
	entry2, err := testQueries.GetEntry(context.Background(), entry1.ID)
	
	assert.NoError(t, err)
	assert.NotEmpty(t, entry2)

	assert.Equal(t, entry1.ID, entry2.ID)
	assert.Equal(t, entry1.Amount, entry2.Amount)
	assert.Equal(t, entry1.AccountID, entry2.AccountID)

	assert.WithinDuration(t, entry1.CreatedAt, entry2.CreatedAt, time.Second)
}

func TestUpdateEntry(t *testing.T){
	entry := createRandomEntry(t)
    params := UpdateEntryParams{
		ID :entry.ID,
		Amount:utils.RandomMoney(),
	}
    updatedEntry, err := testQueries.UpdateEntry(context.Background(), params)
	assert.NoError(t, err)
	assert.NotEmpty(t, updatedEntry)

	assert.Equal(t, entry.ID, updatedEntry.ID)
	assert.Equal(t, entry.AccountID, updatedEntry.AccountID)
	assert.NotEqual(t, entry.Amount, updatedEntry.Amount)

    assert.WithinDuration(t, entry.CreatedAt, updatedEntry.CreatedAt, time.Second)
}

func TestDeleteEntry(t *testing.T){
	entry := createRandomEntry(t)

	err := testQueries.DeleteEntry(context.Background(), entry.ID)

	assert.NoError(t, err)
	
	deletedEntry, err := testQueries.GetEntry(context.Background(), entry.ID)
   
	assert.Error(t, err)
	assert.EqualError(t, err, sql.ErrNoRows.Error()) // надо
	assert.Empty(t, deletedEntry)

}

func TestGetEntries(t *testing.T){
	for i :=0; i<10; i++{
		createRandomEntry(t)
	}
    params:= GetEntriesParams{
		Limit:5,
		Offset:5,
	}
	entries, err := testQueries.GetEntries(context.Background(), params)
    
	assert.NoError(t, err)
	assert.Len(t, entries, 5)

	for _, entry := range entries{
		assert.NotEmpty(t, entry)
	}

}