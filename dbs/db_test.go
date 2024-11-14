package dbs_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/ek-os/dbs"

	_ "github.com/mattn/go-sqlite3" // import sqlite3 driver
)

func TestDBS(t *testing.T) {
	sqldb, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}

	if _, err := sqldb.Exec(`
		CREATE TABLE users (
			id INTEGER PRIMARY KEY,
			first_name TEXT NOT NULL,
			last_name TEXT NOT NULL
		);`); err != nil {
		t.Fatalf("failed to create users table: %s", err)
	}

	var (
		ctx = context.Background()
		db  = dbs.New(sqldb)
	)

	// one-off queries using db
	id, err := db.SaveUser(ctx, "foo", "bar")
	if err != nil {
		t.Fatalf("failed to save user: %s", err)
	}

	u, err := db.FindUser(ctx, id)
	if err != nil {
		t.Fatalf("failed to find user by id %d: %s", id, err)
	}

	t.Log(u)

	// use the Tx and call the same methods as in db
	tx, err := db.BeginTx(ctx, &dbs.TxOptions{Isolation: dbs.LevelReadCommitted})
	if err != nil {
		t.Fatalf("failed to begin tx: %s", err)
	}

	id2, err := tx.SaveUser(ctx, "baz", "qux")
	if err != nil {
		tx.Rollback()
		t.Fatalf("failed to save user: %s", err)
	}

	u, err = tx.FindUser(ctx, id2)
	if err != nil {
		tx.Rollback()
		t.Fatalf("failed to find user by id %d: %s", id2, err)
	}

	t.Log(u)

	if err := tx.Commit(); err != nil {
		t.Fatalf("failed to commit tx: %s", err)
	}
}
