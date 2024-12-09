package db

import (
	"context"
	"database/sql"
)

func newTx(db *DB, tx *sql.Tx) *Tx {
	res := &Tx{
		tx:    tx,
		stmts: db,
	}
	res.Queries = &Queries{stmts: res}
	return res
}

type Tx struct {
	tx    *sql.Tx
	stmts stmts
	*Queries
}

func (tx *Tx) Commit() error {
	return tx.tx.Commit()
}

func (tx *Tx) Rollback() error {
	return tx.tx.Rollback()
}

func (tx *Tx) stmt(ctx context.Context, query string) (*sql.Stmt, error) {
	stmt, err := tx.stmts.stmt(ctx, query)
	if err != nil {
		return nil, err
	}

	return tx.tx.StmtContext(ctx, stmt), nil
}
