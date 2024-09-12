package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store interface {
	Querier
	TransferTx(ctx context.Context, params TransferTxParams) (TransferTxResult, error)
}

// for transactions since sqlc is for single query
type SQLStore struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB)Store{
	return &SQLStore{
		db:db,
		Queries:New(db),
	}
}

func(store *SQLStore)execTx(ctx context.Context, fn func(*Queries)error)error{
	tx, err := store.db.BeginTx(ctx,nil)
	if err!=nil{
		return nil
	}
	q:=New(tx)
	err =fn(q)

	if err!=nil{
		if rbErr:=tx.Rollback(); rbErr!=nil{
			return fmt.Errorf("could not rol back %s %s", err, rbErr)
		}
		return err
	}
	return tx.Commit()
}
type TransferTxParams struct {
	FromAccountID int64 `json:"from_accont_id"`
	ToAccountID int64 `json:"to_account_id"`
	Amount int64 `json:"amount"`
}

type TransferTxResult struct {
	Transfer Transfer `json:"transfer"`
	FromAccount Account `json:"from_account"`
	ToAccount Account `json:"to_account"`
	FromEntry Entry `json:"from_entry"`
	ToEntry Entry `json:"to_entry"`
}

var txKey = struct{}{}
// transfer money from acc1 to acc2
//creates a transfer, add account's 1 and 2 entry, update acc1 balance, acc2 balance
func (store *SQLStore) TransferTx(ctx context.Context, params TransferTxParams) (TransferTxResult, error) {
    var result TransferTxResult
    err := store.execTx(ctx, func(q *Queries) error {
        var err error
        txName := ctx.Value(txKey)
        fmt.Println(txName, "create transfer")

        result.Transfer, err = q.InsertTransfer(ctx, InsertTransferParams{
            FromAccountID: params.FromAccountID,
            ToAccountID:   params.ToAccountID,
            Amount:        params.Amount,
        })
        if err != nil {
            return err
        }
        fmt.Println(txName, "create entry1")

        result.FromEntry, err = q.InsertEntry(ctx, InsertEntryParams{
            AccountID: params.FromAccountID,
            Amount:    -params.Amount,
        })
        if err != nil {
            return err
        }
        fmt.Println(txName, "create entry2")

        result.ToEntry, err = q.InsertEntry(ctx, InsertEntryParams{
            AccountID: params.ToAccountID,
            Amount:    params.Amount,
        })
        if err != nil {
            return err
        }

        if params.FromAccountID < params.ToAccountID {
			result.FromAccount, result.ToAccount, err =addMoney(context.Background(), q,params.FromAccountID, -params.Amount, params.ToAccountID, params.Amount)
			if err != nil {
				return err
			}
		} else {
			result.ToAccount, result.FromAccount, err =addMoney(context.Background(), q, params.ToAccountID, params.Amount,params.FromAccountID, -params.Amount)
			if err != nil {
				return err
			}
		}
		return nil
    })
    return result, err
}

func addMoney(ctx context.Context, q *Queries, account1ID int64, amount1 int64, account2ID int64,amount2 int64) (account1 Account, account2 Account, err error) {
	// Добавление суммы к первому аккаунту
	account1, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     account1ID,
		Amount: amount1,
	})
	if err != nil {
		return 
	}
	account2, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     account2ID,
		Amount: amount2,
	})
	if err != nil {
		return 
	}
	return account1, account2, nil
}

func doIt()(){
	a:=5
	defer func() {
		fmt.Println(a)
	}()
	a = 11
	fmt.Println(1)
	return
}