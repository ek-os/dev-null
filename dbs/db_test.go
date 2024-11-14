package dbs_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/ek-os/dbs"
)

func TestDBS(t *testing.T) {
	sqldb, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}

	var (
		ctx = context.Background()
		db  = dbs.New(sqldb)
	)

	// one-off queries using db
	if err := db.SaveUser(ctx, "foo", "bar"); err != nil {
		t.Fatalf("failed to save user: %s", err)
	}

	u, err := db.FindUser(ctx, 1)
	if err != nil {
		t.Fatalf("failed to find user by id %d: %s", 1, err)
	}

	t.Log(u)

	// use the Tx and call the same methods as in db
	tx, err := db.BeginTx(ctx, &dbs.TxOptions{Isolation: dbs.LevelReadCommitted})
	if err != nil {
		t.Fatalf("failed to begin tx: %s", err)
	}

	if err := tx.SaveUser(ctx, "foo", "bar"); err != nil {
		tx.Rollback()
		t.Fatalf("failed to save user: %s", err)
	}

	u, err = tx.FindUser(ctx, 1)
	if err != nil {
		tx.Rollback()
		t.Fatalf("failed to find user by id %d: %s", 1, err)
	}

	t.Log(u)

	if err := tx.Commit(); err != nil {
		t.Fatalf("failed to commit tx: %s", err)
	}
}
