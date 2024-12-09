package lmg

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
)

const (
	ENV_CHANGELOG = "LMG_CHANGELOG_PATH"
)

func Run() {
	if err := run(
		context.Background(),
		os.Getenv,
		os.Stdout,
	); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run(
	ctx context.Context,
	getenv func(string) string,
	stdout io.Writer,
) error {
	changelogPath := getenv(ENV_CHANGELOG)
	migrations, err := readChangelog(changelogPath)
	if err != nil {
		return fmt.Errorf("read changelog: %w", err)
	}

	for _, migration := range migrations {
		fmt.Fprintln(stdout, migration)
	}
	return nil
}

func readChangelog(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var migrationPaths []string
	s := bufio.NewScanner(f)
	for s.Scan() {
		migrationPaths = append(migrationPaths, s.Text())
	}

	if err := s.Err(); err != nil {
		return nil, err
	}

	return migrationPaths, nil
}
