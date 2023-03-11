package db

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/golang/mock/mockgen/model"
)

// create general method signature for different kind of store to connect to ie mock db / prod db
type Store interface {
	TransferTx(ctx context.Context, arg TranferTxParams) (TransferTxResult, error)
	CreateUserTx(ctx context.Context, arg CreateUserTxParams) (CreateUserTxResult, error)
	Querier
}

// provide functions to exec db query and transactions
type SQLStore struct {
	// https://stackoverflow.com/a/36706189
	// if we embed the Queries struct as a pointer, we need to initialize the pointer separately when creating a new SQLStore instance,
	// while if we embed the Queries struct directly, it is automatically initialized when creating a new SQLStore instance.
	*Queries // embedding - all func will be available inside struct (queries stuct only support table operations)
	db       *sql.DB
}

func NewStore(db *sql.DB) Store {
	return &SQLStore{
		db:      db,
		Queries: New(db), // initialize and embed queries
	}
}

// pointer receiver
func (store *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	// default isolation level - nil == read committed
	// replace with &sql.TxOptions{}
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}

// add money to different account as part of a transaction
func addMoney(ctx context.Context, q *Queries, accountID1 int64, amount1 int64, accountID2 int64, amount2 int64) (account1 Account, account2 Account, err error) {
	account1, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accountID1,
		Amount: amount1,
	})
	if err != nil {
		return
	}

	account2, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accountID2,
		Amount: amount2,
	})
	if err != nil {
		return
	}

	return
}
