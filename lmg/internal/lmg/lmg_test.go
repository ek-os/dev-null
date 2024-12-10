package lmg_test

import (
	"bytes"
	"context"
	"database/sql"
	"testing"

	"github.com/ek-os/lmg/internal/lmg"
)

func TestRun(t *testing.T) {
	t.Run("correctly applies changelog migrations", testCorrectlyReadsChangelog)
	t.Run("returns error when changelog doesn't exist", testFailReadChangelog)
}

func testCorrectlyReadsChangelog(t *testing.T) {
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

	noErr(
		t,
		lmg.Test(
			context.Background(),
			env.get,
			stdout,
		),
	)

	db, err := sql.Open(driver, dsn)
	noErr(t, err)

	rs, err := db.Query("SELECT type, name, tbl_name, rootpage, sql FROM sqlite_master")
	noErr(t, err)

	var found bool
	for rs.Next() {
		var r sqliteMasterRow
		noErr(t, rs.Scan(&r.Type, &r.Name, &r.TblName, &r.RootPage, &r.SQL))
		if r.Name == "users" {
			found = true
			break
		}
	}

	if !found {
		t.Errorf(`Expected table "users" to be present`)
	}
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

	var (
		gotErr  = err.Error()
		wantErr = "read changelog: open foo: no such file or directory"
	)

	if gotErr != wantErr {
		t.Errorf("Error doesn't match.\nwant: %q\ngot:  %q\n", wantErr, gotErr)
	}
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

type sqliteMasterRow struct {
	Type     string
	Name     string
	TblName  string
	RootPage int
	SQL      sql.Null[string]
}
