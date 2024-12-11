package lmg_test

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/ek-os/lmg"

	_ "github.com/mattn/go-sqlite3"
)

func TestRun(t *testing.T) {
	t.Run("correctly applies changelog migrations", testCorrectlyAppliesChangelog)
	t.Run("returns error when changelog refers to non-existing migration", testFailFindMigration)
	t.Run("returns error when changelog doesn't exist", testFailReadChangelog)
}

func testCorrectlyAppliesChangelog(t *testing.T) {
	var (
		driver = "sqlite3"
		dsn    = "file::memory:?cache=shared"
		env    = mapenv{
			lmg.ENV_CHANGELOG: "testdata/changelog.txt",
			lmg.ENV_DRIVER:    driver,
			lmg.ENV_DSN:       dsn,
		}
		stdout = &bytes.Buffer{}
	)

	err := lmg.TestRun(
		context.Background(),
		env.get,
		stdout,
	)
	noErr(t, err)

	db, err := openTestDB(driver, dsn)
	noErr(t, err)

	exists, err := db.tableExists("users")
	noErr(t, err)

	if !exists {
		t.Errorf(`Expected table "users" to be present`)
	}
}

func testFailFindMigration(t *testing.T) {
	var (
		driver = "sqlite3"
		dsn    = ":memory:"
		env    = mapenv{
			lmg.ENV_CHANGELOG: "testdata/changelog-dud.txt",
			lmg.ENV_DRIVER:    driver,
			lmg.ENV_DSN:       dsn,
		}
		stdout = &bytes.Buffer{}
	)

	err := lmg.TestRun(
		context.Background(),
		env.get,
		stdout,
	)

	errIsString(t, err, "execute testdata/migrations/bar.sql: open testdata/migrations/bar.sql: no such file or directory")
}

func testFailReadChangelog(t *testing.T) {
	env := mapenv{
		lmg.ENV_CHANGELOG: "foo",
	}

	stdout := &bytes.Buffer{}

	err := lmg.TestRun(
		context.Background(),
		env.get,
		stdout,
	)

	errIsString(t, err, "read changelog: open foo: no such file or directory")
}

type mapenv map[string]string

func (m mapenv) get(key string) string {
	if val, ok := m[key]; ok {
		return val
	}
	return ""
}

func noErr(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}
}

func errIsString(t *testing.T, err error, want string) {
	t.Helper()
	got := err.Error()
	if got != want {
		t.Fatalf("Error doesn't match.\nwant: %q\ngot:  %q\n", want, got)
	}
}

func openTestDB(driver, dsn string) (testDB, error) {
	switch driver {
	case "sqlite3":
		db, err := sql.Open(driver, dsn)
		if err != nil {
			return nil, err
		}
		return &sqlite3TestDB{db: db}, nil
	default:
		return nil, fmt.Errorf("Unknown driver: %s", driver)
	}
}

type testDB interface {
	tableExists(table string) (bool, error)
}

type sqlite3TestDB struct {
	db *sql.DB
}

func (t *sqlite3TestDB) tableExists(table string) (bool, error) {
	var exists bool
	if err := t.db.QueryRow(
		"SELECT true FROM sqlite_master WHERE name = :name",
		sql.Named("name", table),
	).Scan(&exists); err != nil {
		return false, err
	}
	return exists, nil
}
