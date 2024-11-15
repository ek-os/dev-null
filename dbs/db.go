package dbs

import (
	"context"
	"database/sql"
	"sync"
)

func New(db *sql.DB) *DB {
	res := &DB{
		db:    db,
		stmts: new(sync.Map),
	}
	res.Queries = &Queries{stmts: res}
	return res
}

type DB struct {
	db    *sql.DB
	stmts *sync.Map
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
	return newTx(db, tx), nil
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

	return newTx(db, tx), nil
}

func (db *DB) stmt(ctx context.Context, query string) (*sql.Stmt, error) {
	v, ok := db.stmts.Load(query)
	if ok {
		// Statement already prepared and cached.
		return v.(*sql.Stmt), nil
	}

	// Statement not yet prepared.
	stmt, err := db.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}

	res, loaded := db.stmts.LoadOrStore(query, stmt)
	if loaded {
		// Another goroutine managed to LoadOrStore first, close statement
		// and use loaded value.
		if err := stmt.Close(); err != nil {
			return nil, err
		}
	}

	return res.(*sql.Stmt), nil
}
