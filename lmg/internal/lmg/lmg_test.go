package lmg_test

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/ek-os/lmg/internal/lmg"
)

func TestRun(t *testing.T) {
	want := "hello"

	getenv := func(key string) string {
		switch key {
		case lmg.ENV_CHANGELOG:
			return want
		default:
			return ""
		}
	}

	stdout := &bytes.Buffer{}

	lmg.Run(context.Background(), getenv, stdout)

	got := strings.TrimSpace(stdout.String())
	if got != want {
		t.Errorf("%s: expected %q, got %q", lmg.ENV_CHANGELOG, want, got)
	}
}
