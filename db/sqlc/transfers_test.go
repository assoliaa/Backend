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

func createRandomTransfer(t *testing.T) Transfer {
	fromAccountID, err := GetRandomAccountId(context.Background())
	if err != nil {
		fmt.Errorf("could not get accounts: %w", err)// дрегие ошибки использовать
	}
	toAccountID, err := GetRandomAccountId(context.Background())
	if err != nil {
		fmt.Errorf("could not get accounts: %w", err)
	}
	params := InsertTransferParams{
		FromAccountID: fromAccountID,
		ToAccountID:   toAccountID,
		Amount:        utils.RandomMoney(),
	}
	transfer, err := testQueries.InsertTransfer(context.Background(), params)

	assert.NoError(t, err)
	assert.NotEmpty(t, transfer)

	assert.Equal(t, params.FromAccountID, transfer.FromAccountID)
	assert.Equal(t, params.ToAccountID, transfer.ToAccountID)
	assert.Equal(t, params.Amount, transfer.Amount)

	assert.NotZero(t, transfer.ID)
	return transfer
}
func TestCreateTransfer(t *testing.T) {
	createRandomTransfer(t)
}

func TestGetTransfer(t *testing.T) {
	transfer1 := createRandomTransfer(t)
	transfer2, err := testQueries.GetTransfer(context.Background(), transfer1.ID)

	assert.NotEmpty(t, transfer2)
	assert.NoError(t, err)

	assert.Equal(t, transfer1.ID, transfer2.ID)
	assert.Equal(t, transfer1.FromAccountID, transfer2.FromAccountID)
	assert.Equal(t, transfer1.ToAccountID, transfer2.ToAccountID)
	assert.Equal(t, transfer1.Amount, transfer2.Amount)

	assert.WithinDuration(t, transfer1.CreatedAt, transfer2.CreatedAt, time.Second)
}

func TestUpdateTransfer(t *testing.T) {
	transfer1 := createRandomTransfer(t)

	params := UpdateTransferParams{
		ID:     transfer1.ID,
		Amount: utils.RandomMoney(),
	}

	transfer2, err := testQueries.UpdateTransfer(context.Background(), params)
	assert.NoError(t, err)
	assert.NotEmpty(t, transfer2)

	assert.Equal(t, transfer1.ID, transfer2.ID)
	assert.Equal(t, transfer1.FromAccountID, transfer2.FromAccountID)
	assert.Equal(t, transfer1.ToAccountID, transfer2.ToAccountID)
	assert.NotEqual(t, transfer1.Amount, transfer2.Amount)

    assert.WithinDuration(t, transfer1.CreatedAt, transfer2.CreatedAt, time.Second)
}

func TestDeleteTransfer(t *testing.T) {
	transfer1 := createRandomTransfer(t)
	err := testQueries.DeleteTransfer(context.Background(), transfer1.ID)

	assert.NoError(t, err)

	transfer2, err := testQueries.GetTransfer(context.Background(), transfer1.ID)

	assert.Error(t, err)
	assert.EqualError(t, err, sql.ErrNoRows.Error())
	assert.Empty(t, transfer2)
}

func TestGetTransfers(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomTransfer(t)
	}
	params := GetTransfersParams{
		Limit:  5,
		Offset: 5,
	}

	transfers, err := testQueries.GetTransfers(context.Background(), params)
	assert.NoError(t, err)
	assert.Len(t, transfers, 5)
}
