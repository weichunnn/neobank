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

// neccessary input param for transfer transaction
type TranferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

/*
problem and solution

1. concurrent transaction causing updates to be unequal - solve by having a FOR UPDATE lock
2. locking may cause deadlock (tx 1 depends on tx2 and vice verse) - solve by using NO KEY - allow concurrent transaction that does use the KEY COLUMN in their operations
  - mainly to allow INSERT into TRANSFER (FK contrains block acc table) and SELECT FOR UDPATE needing to lock, now allowing them to run concurrently as we guarantee the key won't be changed
  - LOCK by TX updates are still working

3. Pairwise transaction deadlock (ie tx1 -> tx2, tx2 -> tx1 concurrent), solved by allowing smaller account id to run first (deadlock won't occur)
*/
func (store *SQLStore) TransferTx(ctx context.Context, arg TranferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error
		// txName := ctx.Value(txKey)

		// fmt.Println(txName, "create transfer")
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}

		// fmt.Println(txName, "create entry 1")
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}

		// fmt.Println(txName, "create entry 2")
		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}

		// lock to prevent other transactions interfering with current operations
		// always need to acquire the lock of smaller account ID first
		// fmt.Println(txName, "update account 1")
		if arg.FromAccountID < arg.ToAccountID {
			result.FromAccount, result.ToAccount, err = addMoney(ctx, q, arg.FromAccountID, -arg.Amount, arg.ToAccountID, arg.Amount)
		} else {
			result.ToAccount, result.FromAccount, err = addMoney(ctx, q, arg.ToAccountID, arg.Amount, arg.FromAccountID, -arg.Amount)
		}
		if err != nil {
			return err
		}

		return nil // base case
	})

	return result, err
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
