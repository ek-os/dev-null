package lmg

import (
	"context"
	"fmt"
	"io"
)

const (
	ENV_CHANGELOG = "LMG_CHANGELOG_PATH"
)

func Run(
	ctx context.Context,
	getenv func(string) string,
	stdout io.Writer,
) {
	fmt.Fprintln(stdout, getenv(ENV_CHANGELOG))
}
