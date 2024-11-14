package dbs

import "database/sql"

func newTx(tx *sql.Tx) *Tx {
	return &Tx{
		tx:      tx,
		Queries: &Queries{dbi: tx},
	}
}

type Tx struct {
	tx *sql.Tx
	*Queries
}

func (tx *Tx) Commit() error {
	return tx.tx.Commit()
}

func (tx *Tx) Rollback() error {
	return tx.tx.Rollback()
}
