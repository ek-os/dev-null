package lmg_test

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"testing"

	"github.com/ek-os/lmg"

	_ "github.com/mattn/go-sqlite3"
)

func TestCorrectlyHandlesTrailingWhitespace(t *testing.T) {
	sys := newTestSystem(map[string]string{
		lmg.ENV_CHANGELOG: "testdata/changelog-whitespace.txt",
		lmg.ENV_DSN:       persistentDSN,
		lmg.ENV_DRIVER:    driver,
	})

	err := lmg.TestRun(context.Background(), sys)
	noErr(t, err)

	db, err := openTestDB(driver, persistentDSN)
	noErr(t, err)

	exists, err := db.tableExists("users")
	noErr(t, err)

	if !exists {
		t.Errorf(`Expected table "users" to be present`)
	}
}

func TestCorrectlyAppliesChangelog(t *testing.T) {
	sys := newTestSystem(map[string]string{
		lmg.ENV_CHANGELOG: "testdata/changelog.txt",
		lmg.ENV_DSN:       persistentDSN,
		lmg.ENV_DRIVER:    driver,
	})

	err := lmg.TestRun(context.Background(), sys)
	noErr(t, err)

	db, err := openTestDB(driver, persistentDSN)
	noErr(t, err)

	exists, err := db.tableExists("users")
	noErr(t, err)

	if !exists {
		t.Errorf(`Expected table "users" to be present`)
	}
}

func TestFailFindMigration(t *testing.T) {
	sys := newTestSystem(map[string]string{
		lmg.ENV_CHANGELOG: "testdata/changelog-dud.txt",
		lmg.ENV_DSN:       transientDSN,
		lmg.ENV_DRIVER:    driver,
	})

	err := lmg.TestRun(context.Background(), sys)

	errIsString(t, err, "execute testdata/migrations/bar.sql: open testdata/migrations/bar.sql: no such file or directory")
}

func TestFailReadChangelog(t *testing.T) {
	sys := newTestSystem(map[string]string{
		lmg.ENV_CHANGELOG: "foo",
		lmg.ENV_DSN:       transientDSN,
		lmg.ENV_DRIVER:    driver,
	})

	err := lmg.TestRun(context.Background(), sys)

	errIsString(t, err, "read changelog: open foo: no such file or directory")
}

const (
	driver = "sqlite3"

	persistentDSN = "file::memory:?cache=shared"
	transientDSN  = ":memory:"
)

func newTestSystem(env map[string]string) testSystem {
	return testSystem{
		env:    env,
		stdout: &bytes.Buffer{},
	}
}

type testSystem struct {
	env    map[string]string
	stdout *bytes.Buffer
}

func (t testSystem) Getenv(key string) string {
	if val, ok := t.env[key]; ok {
		return val
	}
	return ""
}

func (t testSystem) Stdout() io.Writer {
	return t.stdout
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
	assertTableExists(table string) func(t *testing.T)
	assertTableDoesntExist(table string) func(t *testing.T)
}

type sqlite3TestDB struct {
	db *sql.DB
}

func (db *sqlite3TestDB) tableExists(table string) (bool, error) {
	var exists bool
	if err := db.db.QueryRow(
		"SELECT true FROM sqlite_master WHERE name = :name",
		sql.Named("name", table),
	).Scan(&exists); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return exists, nil
}

// assertTableExists implements testDB.
func (db *sqlite3TestDB) assertTableExists(table string) func(t *testing.T) {
	return func(t *testing.T) {
		ok, err := db.tableExists(table)
		noErr(t, err)

		if !ok {
			t.Errorf(`Expected table %q to be present`, table)
		}
	}
}

// assertTableDoesntExist implements testDB.
func (db *sqlite3TestDB) assertTableDoesntExist(table string) func(t *testing.T) {
	return func(t *testing.T) {
		ok, err := db.tableExists(table)
		noErr(t, err)

		if ok {
			t.Errorf(`Expected table %q to not be present`, table)
		}
	}
}

func noErr(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}
}
