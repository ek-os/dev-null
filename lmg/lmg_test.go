package lmg_test

import (
	"bytes"
	"context"
	"database/sql"
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

	err := lmg.Test(
		context.Background(),
		env.get,
		stdout,
	)
	noErr(t, err)

	db, err := sql.Open(driver, dsn)
	noErr(t, err)

	rs, err := db.Query("SELECT type, name, tbl_name, rootpage, sql FROM sqlite_master")
	noErr(t, err)

	type sqliteMasterRow struct {
		Type     string
		Name     string
		TblName  string
		RootPage int
		SQL      sql.Null[string]
	}

	var found bool
	for rs.Next() {
		var r sqliteMasterRow
		err = rs.Scan(&r.Type, &r.Name, &r.TblName, &r.RootPage, &r.SQL)
		noErr(t, err)

		if r.Name == "users" {
			found = true
			break
		}
	}

	if !found {
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

	err := lmg.Test(
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

	err := lmg.Test(
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
