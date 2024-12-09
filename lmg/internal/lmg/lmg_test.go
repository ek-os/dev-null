package lmg_test

import (
	"bytes"
	"context"
	"testing"

	"github.com/ek-os/lmg/internal/lmg"
)

func TestRun(t *testing.T) {
	t.Run("correctly reads changelog", testCorrectlyReadsChangelog)
	t.Run("returns error when changelog doesn't exist", testFailReadChangelog)
}

func testCorrectlyReadsChangelog(t *testing.T) {
	getenv := envFromMap(map[string]string{
		lmg.ENV_CHANGELOG: "testdata/changelog.txt",
	})

	stdout := &bytes.Buffer{}

	err := lmg.Test(
		context.Background(),
		getenv,
		stdout,
	)

	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	var (
		wantStdout = "migrations/foo.sql\n"
		gotStdout  = stdout.String()
	)

	if gotStdout != wantStdout {
		t.Errorf("Stdout doesn't match.\nwant: %q\ngot:  %q\n", wantStdout, gotStdout)
	}
}

func testFailReadChangelog(t *testing.T) {
	getenv := envFromMap(map[string]string{
		lmg.ENV_CHANGELOG: "foo",
	})

	stdout := &bytes.Buffer{}

	err := lmg.Test(
		context.Background(),
		getenv,
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

func envFromMap(m map[string]string) func(string) string {
	return func(key string) string {
		if val, ok := m[key]; ok {
			return val
		}
		return ""
	}
}
