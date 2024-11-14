package dbs

import (
	"context"
	"database/sql"
)

func New(db *sql.DB) *DB {
	return &DB{
		db:      db,
		Queries: &Queries{dbi: db},
	}
}

type DB struct {
	db *sql.DB
	*Queries
}

// Avoid imports to database/sql in code creating transactions
type TxOptions struct {
	Isolation IsolationLevel
	ReadOnly  bool
}

// IsolationLevel is the transaction isolation level used in [TxOptions].
type IsolationLevel int

const (
	LevelDefault         = IsolationLevel(sql.LevelDefault)
	LevelReadUncommitted = IsolationLevel(sql.LevelReadUncommitted)
	LevelReadCommitted   = IsolationLevel(sql.LevelReadCommitted)
	LevelWriteCommitted  = IsolationLevel(sql.LevelWriteCommitted)
	LevelRepeatableRead  = IsolationLevel(sql.LevelRepeatableRead)
	LevelSnapshot        = IsolationLevel(sql.LevelSnapshot)
	LevelSerializable    = IsolationLevel(sql.LevelSerializable)
	LevelLinearizable    = IsolationLevel(sql.LevelLinearizable)
)

func (db *DB) Begin() (*Tx, error) {
	tx, err := db.db.Begin()
	if err != nil {
		return nil, err
	}
	return newTx(tx), nil
}

func (db *DB) BeginTx(ctx context.Context, opts *TxOptions) (*Tx, error) {
	var txopts *sql.TxOptions
	if opts != nil {
		txopts = &sql.TxOptions{
			Isolation: sql.IsolationLevel(opts.Isolation),
			ReadOnly:  opts.ReadOnly,
		}
	}

	tx, err := db.db.BeginTx(ctx, txopts)
	if err != nil {
		return nil, err
	}

	return newTx(tx), nil
}
