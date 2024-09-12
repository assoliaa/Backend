// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package db

import (
	"context"
	"database/sql"
	"fmt"
)

type DBTX interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

func New(db DBTX) *Queries {
	return &Queries{db: db}
}

func Prepare(ctx context.Context, db DBTX) (*Queries, error) {
	q := Queries{db: db}
	var err error
	if q.addAccountBalanceStmt, err = db.PrepareContext(ctx, addAccountBalance); err != nil {
		return nil, fmt.Errorf("error preparing query AddAccountBalance: %w", err)
	}
	if q.deleteAccountStmt, err = db.PrepareContext(ctx, deleteAccount); err != nil {
		return nil, fmt.Errorf("error preparing query DeleteAccount: %w", err)
	}
	if q.deleteEntryStmt, err = db.PrepareContext(ctx, deleteEntry); err != nil {
		return nil, fmt.Errorf("error preparing query DeleteEntry: %w", err)
	}
	if q.deleteTransferStmt, err = db.PrepareContext(ctx, deleteTransfer); err != nil {
		return nil, fmt.Errorf("error preparing query DeleteTransfer: %w", err)
	}
	if q.getAccountStmt, err = db.PrepareContext(ctx, getAccount); err != nil {
		return nil, fmt.Errorf("error preparing query GetAccount: %w", err)
	}
	if q.getAccountForUpdateStmt, err = db.PrepareContext(ctx, getAccountForUpdate); err != nil {
		return nil, fmt.Errorf("error preparing query GetAccountForUpdate: %w", err)
	}
	if q.getAccountsStmt, err = db.PrepareContext(ctx, getAccounts); err != nil {
		return nil, fmt.Errorf("error preparing query GetAccounts: %w", err)
	}
	if q.getEntriesStmt, err = db.PrepareContext(ctx, getEntries); err != nil {
		return nil, fmt.Errorf("error preparing query GetEntries: %w", err)
	}
	if q.getEntryStmt, err = db.PrepareContext(ctx, getEntry); err != nil {
		return nil, fmt.Errorf("error preparing query GetEntry: %w", err)
	}
	if q.getTransferStmt, err = db.PrepareContext(ctx, getTransfer); err != nil {
		return nil, fmt.Errorf("error preparing query GetTransfer: %w", err)
	}
	if q.getTransfersStmt, err = db.PrepareContext(ctx, getTransfers); err != nil {
		return nil, fmt.Errorf("error preparing query GetTransfers: %w", err)
	}
	if q.getUserStmt, err = db.PrepareContext(ctx, getUser); err != nil {
		return nil, fmt.Errorf("error preparing query GetUser: %w", err)
	}
	if q.insertAccountStmt, err = db.PrepareContext(ctx, insertAccount); err != nil {
		return nil, fmt.Errorf("error preparing query InsertAccount: %w", err)
	}
	if q.insertEntryStmt, err = db.PrepareContext(ctx, insertEntry); err != nil {
		return nil, fmt.Errorf("error preparing query InsertEntry: %w", err)
	}
	if q.insertTransferStmt, err = db.PrepareContext(ctx, insertTransfer); err != nil {
		return nil, fmt.Errorf("error preparing query InsertTransfer: %w", err)
	}
	if q.insertUserStmt, err = db.PrepareContext(ctx, insertUser); err != nil {
		return nil, fmt.Errorf("error preparing query InsertUser: %w", err)
	}
	if q.updateAccountStmt, err = db.PrepareContext(ctx, updateAccount); err != nil {
		return nil, fmt.Errorf("error preparing query UpdateAccount: %w", err)
	}
	if q.updateEntryStmt, err = db.PrepareContext(ctx, updateEntry); err != nil {
		return nil, fmt.Errorf("error preparing query UpdateEntry: %w", err)
	}
	if q.updateTransferStmt, err = db.PrepareContext(ctx, updateTransfer); err != nil {
		return nil, fmt.Errorf("error preparing query UpdateTransfer: %w", err)
	}
	return &q, nil
}

func (q *Queries) Close() error {
	var err error
	if q.addAccountBalanceStmt != nil {
		if cerr := q.addAccountBalanceStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing addAccountBalanceStmt: %w", cerr)
		}
	}
	if q.deleteAccountStmt != nil {
		if cerr := q.deleteAccountStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing deleteAccountStmt: %w", cerr)
		}
	}
	if q.deleteEntryStmt != nil {
		if cerr := q.deleteEntryStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing deleteEntryStmt: %w", cerr)
		}
	}
	if q.deleteTransferStmt != nil {
		if cerr := q.deleteTransferStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing deleteTransferStmt: %w", cerr)
		}
	}
	if q.getAccountStmt != nil {
		if cerr := q.getAccountStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getAccountStmt: %w", cerr)
		}
	}
	if q.getAccountForUpdateStmt != nil {
		if cerr := q.getAccountForUpdateStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getAccountForUpdateStmt: %w", cerr)
		}
	}
	if q.getAccountsStmt != nil {
		if cerr := q.getAccountsStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getAccountsStmt: %w", cerr)
		}
	}
	if q.getEntriesStmt != nil {
		if cerr := q.getEntriesStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getEntriesStmt: %w", cerr)
		}
	}
	if q.getEntryStmt != nil {
		if cerr := q.getEntryStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getEntryStmt: %w", cerr)
		}
	}
	if q.getTransferStmt != nil {
		if cerr := q.getTransferStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getTransferStmt: %w", cerr)
		}
	}
	if q.getTransfersStmt != nil {
		if cerr := q.getTransfersStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getTransfersStmt: %w", cerr)
		}
	}
	if q.getUserStmt != nil {
		if cerr := q.getUserStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getUserStmt: %w", cerr)
		}
	}
	if q.insertAccountStmt != nil {
		if cerr := q.insertAccountStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing insertAccountStmt: %w", cerr)
		}
	}
	if q.insertEntryStmt != nil {
		if cerr := q.insertEntryStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing insertEntryStmt: %w", cerr)
		}
	}
	if q.insertTransferStmt != nil {
		if cerr := q.insertTransferStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing insertTransferStmt: %w", cerr)
		}
	}
	if q.insertUserStmt != nil {
		if cerr := q.insertUserStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing insertUserStmt: %w", cerr)
		}
	}
	if q.updateAccountStmt != nil {
		if cerr := q.updateAccountStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing updateAccountStmt: %w", cerr)
		}
	}
	if q.updateEntryStmt != nil {
		if cerr := q.updateEntryStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing updateEntryStmt: %w", cerr)
		}
	}
	if q.updateTransferStmt != nil {
		if cerr := q.updateTransferStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing updateTransferStmt: %w", cerr)
		}
	}
	return err
}

func (q *Queries) exec(ctx context.Context, stmt *sql.Stmt, query string, args ...interface{}) (sql.Result, error) {
	switch {
	case stmt != nil && q.tx != nil:
		return q.tx.StmtContext(ctx, stmt).ExecContext(ctx, args...)
	case stmt != nil:
		return stmt.ExecContext(ctx, args...)
	default:
		return q.db.ExecContext(ctx, query, args...)
	}
}

func (q *Queries) query(ctx context.Context, stmt *sql.Stmt, query string, args ...interface{}) (*sql.Rows, error) {
	switch {
	case stmt != nil && q.tx != nil:
		return q.tx.StmtContext(ctx, stmt).QueryContext(ctx, args...)
	case stmt != nil:
		return stmt.QueryContext(ctx, args...)
	default:
		return q.db.QueryContext(ctx, query, args...)
	}
}

func (q *Queries) queryRow(ctx context.Context, stmt *sql.Stmt, query string, args ...interface{}) *sql.Row {
	switch {
	case stmt != nil && q.tx != nil:
		return q.tx.StmtContext(ctx, stmt).QueryRowContext(ctx, args...)
	case stmt != nil:
		return stmt.QueryRowContext(ctx, args...)
	default:
		return q.db.QueryRowContext(ctx, query, args...)
	}
}

type Queries struct {
	db                      DBTX
	tx                      *sql.Tx
	addAccountBalanceStmt   *sql.Stmt
	deleteAccountStmt       *sql.Stmt
	deleteEntryStmt         *sql.Stmt
	deleteTransferStmt      *sql.Stmt
	getAccountStmt          *sql.Stmt
	getAccountForUpdateStmt *sql.Stmt
	getAccountsStmt         *sql.Stmt
	getEntriesStmt          *sql.Stmt
	getEntryStmt            *sql.Stmt
	getTransferStmt         *sql.Stmt
	getTransfersStmt        *sql.Stmt
	getUserStmt             *sql.Stmt
	insertAccountStmt       *sql.Stmt
	insertEntryStmt         *sql.Stmt
	insertTransferStmt      *sql.Stmt
	insertUserStmt          *sql.Stmt
	updateAccountStmt       *sql.Stmt
	updateEntryStmt         *sql.Stmt
	updateTransferStmt      *sql.Stmt
}

func (q *Queries) WithTx(tx *sql.Tx) *Queries {
	return &Queries{
		db:                      tx,
		tx:                      tx,
		addAccountBalanceStmt:   q.addAccountBalanceStmt,
		deleteAccountStmt:       q.deleteAccountStmt,
		deleteEntryStmt:         q.deleteEntryStmt,
		deleteTransferStmt:      q.deleteTransferStmt,
		getAccountStmt:          q.getAccountStmt,
		getAccountForUpdateStmt: q.getAccountForUpdateStmt,
		getAccountsStmt:         q.getAccountsStmt,
		getEntriesStmt:          q.getEntriesStmt,
		getEntryStmt:            q.getEntryStmt,
		getTransferStmt:         q.getTransferStmt,
		getTransfersStmt:        q.getTransfersStmt,
		getUserStmt:             q.getUserStmt,
		insertAccountStmt:       q.insertAccountStmt,
		insertEntryStmt:         q.insertEntryStmt,
		insertTransferStmt:      q.insertTransferStmt,
		insertUserStmt:          q.insertUserStmt,
		updateAccountStmt:       q.updateAccountStmt,
		updateEntryStmt:         q.updateEntryStmt,
		updateTransferStmt:      q.updateTransferStmt,
	}
}